package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"

    "github.com/go-redis/redis/v8"
)

type Message struct {
    Message   string `json:"message"`
    Timestamp string `json:"timestamp"`
}

func main() {
    ctx := context.Background()
    rdb := redis.NewClient(&redis.Options{
        Addr: "redis:6379", // Use the Redis service name
    })

    pubsub := rdb.Subscribe(ctx, "mychannel")
    defer pubsub.Close()

    // Goroutine to handle incoming messages
    go func() {
        for {
            msg, err := pubsub.ReceiveMessage(ctx)
            if err != nil {
                log.Println("Error receiving message:", err)
                return
            }
            
            var receivedMessage Message
            if err := json.Unmarshal([]byte(msg.Payload), &receivedMessage); err != nil {
                log.Println("Error unmarshalling JSON:", err)
                continue
            }

            fmt.Printf("Received: %+v\n", receivedMessage)
        }
    }()

    // Keep the main function running
    select {}
}
