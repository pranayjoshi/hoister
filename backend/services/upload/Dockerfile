FROM golang:1.22

WORKDIR /home/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN chmod +x main.sh

RUN go build -o main .

ENTRYPOINT [ "/home/app/main.sh" ]