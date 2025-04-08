FROM golang:1.24

WORKDIR /posts_service

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main ./cmd

EXPOSE 50001
CMD ["./main"]
