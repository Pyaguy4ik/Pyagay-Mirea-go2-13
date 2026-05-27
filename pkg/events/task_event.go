package events

import "time"

// TaskEvent описывает событие, связанное с задачей
type TaskEvent struct {
    Event   string `json:"event"`    // например "task.created"
    TaskID  string `json:"task_id"`
    TS      string `json:"ts"`       // временная метка в формате RFC3339
    // опциональные поля
    RequestID string `json:"request_id,omitempty"`
    Producer  string `json:"producer,omitempty"`
    Version   string `json:"version,omitempty"`
}

// NewTaskCreated создаёт событие о создании задачи
func NewTaskCreated(taskID string) TaskEvent {
    return TaskEvent{
        Event:  "task.created",
        TaskID: taskID,
        TS:     time.Now().UTC().Format(time.RFC3339),
    }
}
