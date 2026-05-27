package amqpclient

import (
    "log"
    amqp "github.com/rabbitmq/amqp091-go"
)

// MustConnect подключается к RabbitMQ, паникует при ошибке
func MustConnect(url string) *amqp.Connection {
    conn, err := amqp.Dial(url)
    if err != nil {
        log.Fatalf("Failed to connect to RabbitMQ: %v", err)
    }
    return conn
}

// DeclareQueue объявляет durable очередь (переживает перезапуск брокера)
func DeclareQueue(ch *amqp.Channel, queueName string) (amqp.Queue, error) {
    return ch.QueueDeclare(
        queueName, // name
        true,      // durable
        false,     // delete when unused
        false,     // exclusive
        false,     // no-wait
        nil,       // arguments
    )
}
