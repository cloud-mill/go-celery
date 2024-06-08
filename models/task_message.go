package models

type TaskMessage struct {
	Properties Properties `json:"properties"`
	Headers    Headers    `json:"headers"`
	Body       Body       `json:"body"`
}

type Properties struct {
	CorrelationID   string `json:"correlation_id"`
	ContentType     string `json:"content_type"`
	ContentEncoding string `json:"content_encoding"`
	ReplyTo         string `json:"reply_to,omitempty"`
}

type Headers struct {
	Lang                string `json:"lang"`
	Task                string `json:"task"`
	ID                  string `json:"id"`
	RootID              string `json:"root_id"`
	ParentID            string `json:"parent_id,omitempty"`
	Group               string `json:"group,omitempty"`
	Meth                string `json:"meth,omitempty"`
	Shadow              string `json:"shadow,omitempty"`
	Eta                 string `json:"eta,omitempty"`
	Expires             string `json:"expires,omitempty"`
	Retries             int    `json:"retries"`
	TimeLimit           [2]int `json:"timelimit"`
	ArgsRepr            string `json:"argsrepr"`
	KwargsRepr          string `json:"kwargsrepr"`
	Origin              string `json:"origin"`
	ReplacedTaskNesting int    `json:"replaced_task_nesting,omitempty"`
}

type Body struct {
	Args   []interface{}          `json:"args"`
	Kwargs map[string]interface{} `json:"kwargs"`
	Embed  Embed                  `json:"embed"`
}

type Embed struct {
	Callbacks []Signature `json:"callbacks"`
	Errbacks  []Signature `json:"errbacks"`
	Chain     []Signature `json:"chain"`
	Chord     *Signature  `json:"chord,omitempty"`
}

type Signature struct {
	Task   string                 `json:"task"`
	Args   []interface{}          `json:"args"`
	Kwargs map[string]interface{} `json:"kwargs"`
}
