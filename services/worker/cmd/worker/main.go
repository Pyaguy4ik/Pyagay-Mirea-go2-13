package main

import (
    "encoding/json"
    "log"
    "os"

    "example.com/pz13-rabbit/pkg/amqpclient"
    "example.com/pz13-rabbit/pkg/events"
)

func main() {
    rabbitURL := os.Getenv("RABBIT_URL")
    if rabbitURL == "" {
        rabbitURL = "amqp://guest:guest@localhost:5672/"
    }
    queueName := os.Getenv("QUEUE_NAME")
    if queueName == "" {
        queueName = "task_events"
    }

    conn := amqpclient.MustConnect(rabbitURL)
    defer conn.Close()

    ch, err := conn.Channel()
    if err != nil {
        log.Fatal("Failed to open channel:", err)
    }
    defer ch.Close()

    // Объявляем очередь (должна быть та же, что и в producer)
    _, err = amqpclient.DeclareQueue(ch, queueName)
    if err != nil {
        log.Fatal("Failed to declare queue:", err)
    }

    // Устанавливаем prefetch = 1
    if err := ch.Qos(1, 0, false); err != nil {
        log.Fatal("Failed to set QoS:", err)
    }

    msgs, err := ch.Consume(
        queueName,
        "",
        false,
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        log.Fatal("Failed to register consumer:", err)
    }

    log.Println("Worker started, waiting for messages...")

    for msg := range msgs {
        var event events.TaskEvent
        if err := json.Unmarshal(msg.Body, &event); err != nil {
            log.Printf("Failed to unmarshal message: %v", err)
            msg.Nack(false, false)
            continue
        }

        log.Printf("Received event: %s, task_id: %s, ts: %s", event.Event, event.TaskID, event.TS)

        if err := msg.Ack(false); err != nil {
            log.Printf("Failed to ack message: %v", err)
        }
    }
}
