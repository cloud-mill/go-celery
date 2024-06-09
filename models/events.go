package models

import "time"

const (
	EventTaskReceived  = "task-received"
	EventTaskStarted   = "task-started"
	EventTaskSucceeded = "task-succeeded"
	EventTaskFailed    = "task-failed"
	EventTaskRevoked   = "task-revoked"
	EventTaskRetry     = "task-retry"
)

type BaseEvent struct {
	Event     string    `json:"event"`
	UUID      string    `json:"uuid"`
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type"`
	Hostname  string    `json:"hostname"`
}

type TaskReceivedEvent struct {
	BaseEvent
	Args   []interface{}          `json:"args"`
	Kwargs map[string]interface{} `json:"kwargs"`
}

type TaskStartedEvent struct {
	BaseEvent
}

type TaskSucceededEvent struct {
	BaseEvent
	Result interface{} `json:"result"`
}

type TaskFailedEvent struct {
	BaseEvent
	Exception string `json:"exception"`
	Traceback string `json:"traceback"`
}

type TaskRevokedEvent struct {
	BaseEvent
}

type TaskRetryEvent struct {
	BaseEvent
	Exception string    `json:"exception"`
	Traceback string    `json:"traceback"`
	ETA       time.Time `json:"eta"`
}
