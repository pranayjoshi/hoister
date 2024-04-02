FROM golang:1.17

WORKDIR /home/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

ENTRYPOINT [ "./main" ]