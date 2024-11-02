package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

// const (

// )

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	redisClient *redis.Client
	ecsClient   *ecs.Client
)

type ProjectRequest struct {
	GitURL string `json:"gitURL"`
	Slug   string `json:"slug"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}
	PORT := os.Getenv("RUN_PORT")
	REDIS_ADDR := os.Getenv("REDIS_ADDR")
	AWS_REGION := os.Getenv("AWS_REGION")
	AWS_ACCESS_KEY := os.Getenv("AWS_ACCESS_KEY_ID")
	AWS_SECRET_KEY := os.Getenv("AWS_SECRET_ACCESS_KEY")

	log.Println("PORT: ", REDIS_ADDR)

	router := mux.NewRouter()
	redisClient = redis.NewClient(&redis.Options{
		Addr: REDIS_ADDR,
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(AWS_REGION),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(AWS_ACCESS_KEY, AWS_SECRET_KEY, "")),
	)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	ecsClient = ecs.NewFromConfig(cfg)

	router.HandleFunc("/project", handleProjectCreation).Methods("POST")
	go startWebSocketServer()
	go initRedisSubscribe()

	log.Printf("API Server Running on port %s", PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, router))
}

func handleProjectCreation(w http.ResponseWriter, r *http.Request) {
	ECS_CLUSTER := os.Getenv("ECS_CLUSTER")
	ECS_TASK := os.Getenv("ECS_TASK")
	var req ProjectRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	projectSlug := req.Slug
	if projectSlug == "" {
		projectSlug = uuid.New().String()
	}

	_, err = ecsClient.RunTask(context.TODO(), &ecs.RunTaskInput{
		Cluster:        aws.String(ECS_CLUSTER),
		TaskDefinition: aws.String(ECS_TASK),
		LaunchType:     types.LaunchTypeFargate,
		Count:          aws.Int32(1),
		NetworkConfiguration: &types.NetworkConfiguration{
			AwsvpcConfiguration: &types.AwsVpcConfiguration{
				AssignPublicIp: types.AssignPublicIpEnabled,
				Subnets:        []string{os.Getenv("SUBNET_1"), os.Getenv("SUBNET_2"), os.Getenv("SUBNET_3")},
				SecurityGroups: []string{os.Getenv("SECURITY_GROUPS")},
			},
		},
		Overrides: &types.TaskOverride{
			ContainerOverrides: []types.ContainerOverride{
				{
					Name: aws.String("hoister-upload-service-image"),
					Environment: []types.KeyValuePair{
						{Name: aws.String("GIT_REPOSITORY__URL"), Value: aws.String(req.GitURL)},
						{Name: aws.String("PROJECT_ID"), Value: aws.String(projectSlug)},
					},
				},
			},
		},
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"status": "queued",
		"data": map[string]string{
			"projectSlug": projectSlug,
			"url":         fmt.Sprintf("http://%s.localhost:8000", projectSlug),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func startWebSocketServer() {
	SOCKET_PORT := os.Getenv("SOCKET_PORT")
	router := mux.NewRouter()
	router.HandleFunc("/ws", handleWebSocket)
	log.Printf("WebSocket Server Running on port %s", SOCKET_PORT)
	log.Fatal(http.ListenAndServe(":"+SOCKET_PORT, router))
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v\n", err)
			return
		}
		var message struct {
			Action  string `json:"action"`
			Channel string `json:"channel"`
		}
		err = json.Unmarshal(p, &message)
		if err != nil {
			log.Println(err)
			continue
		}
		if message.Action == "subscribe" {
			err = conn.WriteMessage(messageType, []byte(fmt.Sprintf("Joined %s", message.Channel)))
			if err != nil {
				log.Println("wad", err)
				return
			}
		}
	}
}

func initRedisSubscribe() {
	pubsub := redisClient.PSubscribe(context.Background(), "*")
	// log.Println()
	defer pubsub.Close()

	log.Println("Subscribed to logs....")

	ch := pubsub.Channel()
	for msg := range ch {
		// Broadcast message to all connected WebSocket clients
		// This part is simplified and would need to be implemented
		// based on your WebSocket management strategy
		log.Printf("Received message: %s", msg.Payload)
	}

}
