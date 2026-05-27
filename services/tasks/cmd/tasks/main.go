package main

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "os"
    "sync"

    amqp "github.com/rabbitmq/amqp091-go"
    "example.com/pz13-rabbit/pkg/amqpclient"
    "example.com/pz13-rabbit/pkg/events"
)

type Task struct {
    ID          string  `json:"id"`
    Title       string  `json:"title"`
    Description *string `json:"description,omitempty"`
    Done        bool    `json:"done"`
}

var (
    tasks   = []Task{}
    tasksMu sync.RWMutex
    nextID  = 1
)

func main() {
    // Чтение переменных окружения
    rabbitURL := os.Getenv("RABBIT_URL")
    if rabbitURL == "" {
        rabbitURL = "amqp://guest:guest@localhost:5672/"
    }
    queueName := os.Getenv("QUEUE_NAME")
    if queueName == "" {
        queueName = "task_events"
    }
    port := os.Getenv("PORT")
    if port == "" {
        port = "8082"
    }

    // Подключение к RabbitMQ
    conn := amqpclient.MustConnect(rabbitURL)
    defer conn.Close()
    ch, err := conn.Channel()
    if err != nil {
        log.Fatal("Failed to open channel:", err)
    }
    defer ch.Close()

    // Объявляем очередь (durable)
    _, err = amqpclient.DeclareQueue(ch, queueName)
    if err != nil {
        log.Fatal("Failed to declare queue:", err)
    }

    // HTTP обработчики
    mux := http.NewServeMux()
    mux.HandleFunc("/v1/tasks", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
            return
        }

        var input struct {
            Title       string  `json:"title"`
            Description *string `json:"description"`
        }
        if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
            http.Error(w, "bad request", http.StatusBadRequest)
            return
        }

        // Создание задачи в памяти
        tasksMu.Lock()
        taskID := nextID
        nextID++
        newTask := Task{
            ID:          string(rune(taskID)), // упрощённо, лучше использовать fmt.Sprintf
            Title:       input.Title,
            Description: input.Description,
            Done:        false,
        }
        tasks = append(tasks, newTask)
        tasksMu.Unlock()

        // Публикация события (best effort – ошибка только логируется)
        event := events.NewTaskCreated(newTask.ID)
        body, _ := json.Marshal(event)
        err = ch.PublishWithContext(
            context.Background(),
            "",      // exchange
            queueName,
            false,   // mandatory
            false,   // immediate
            amqp.Publishing{
                ContentType:  "application/json",
                DeliveryMode: amqp.Persistent, // persistent message
                Body:         body,
            },
        )
        if err != nil {
            log.Printf("WARNING: failed to publish event: %v", err)
        } else {
            log.Printf("Published event for task %s", newTask.ID)
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(newTask)
    })

    addr := ":" + port
    log.Printf("Tasks service started on %s", addr)
    if err := http.ListenAndServe(addr, mux); err != nil {
        log.Fatal(err)
    }
}
