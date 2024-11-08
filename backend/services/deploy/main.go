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
	socketio "github.com/googollee/go-socket.io"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

var (
	redisClient  *redis.Client
	ecsClient    *ecs.Client
	socketServer *socketio.Server
)

type ProjectRequest struct {
	GitURL      string `json:"gitURL"`
	ProjectSlug string `json:"projectslug"`
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
	go startSocketIOServer()
	go initRedisSubscribe()

	// Enable CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3001"}, // Adjust this to your client's origin
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	log.Printf("API Server Running on port %s", PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, handler))
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

	fmt.Println("Request: ", req)
	projectSlug := req.ProjectSlug
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

func startSocketIOServer() {
	SOCKET_PORT := os.Getenv("SOCKET_PORT")
	server := socketio.NewServer(nil)

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		return nil
	})

	server.OnEvent("/", "subscribe", func(s socketio.Conn, channel string) {
		fmt.Println("subscribe:", channel)
		s.Join(channel)
		s.Emit("message", fmt.Sprintf("Joined %s", channel))
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", reason)
	})

	go server.Serve()
	defer server.Close()

	http.Handle("/socket.io/", server)
	log.Printf("Socket.IO server running on port %s", SOCKET_PORT)
	log.Fatal(http.ListenAndServe(":"+SOCKET_PORT, nil))
}

func initRedisSubscribe() {
	pubsub := redisClient.PSubscribe(context.Background(), "logs:*")
	defer pubsub.Close()

	log.Println("Subscribed to logs....")

	ch := pubsub.Channel()
	for msg := range ch {
		log.Printf("Received message: %s", msg.Payload)
		socketServer.BroadcastToRoom("/", msg.Channel, "message", msg.Payload)
	}
}
