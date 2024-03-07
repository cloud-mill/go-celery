package models

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
