package models

import (
	"encoding/json"
)

type CeleryMessage struct {
	Body    string                 `json:"body"`
	Headers map[string]interface{} `json:"headers,omitempty"`
}

type TaskMessage struct {
	Id     string                 `json:"id"`
	Task   string                 `json:"task"`
	Args   []interface{}          `json:"args"`
	Kwargs map[string]interface{} `json:"kwargs"`
}

type ResultMessage struct {
	Id     string      `json:"task_id"`
	Status string      `json:"status"`
	Result interface{} `json:"result"`
}

func (cm *CeleryMessage) ExtractTaskMessage() (*TaskMessage, error) {
	var taskMessage TaskMessage
	err := json.Unmarshal([]byte(cm.Body), &taskMessage)
	if err != nil {
		return nil, err
	}

	return &taskMessage, nil
}
