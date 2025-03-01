package types

type (
	TriggerName string
	WorkflowID  string
)

type Trigger struct {
	Name         TriggerName `json:"name"`
	Src          string      `json:"src"`
	Filter       string      `json:"filter"`
	QueueName    string      `json:"qname"`
	WorkflowType string      `json:"wtype"`
}

type Workflow struct {
	TriggerName TriggerName `json:"name"`
	WorkflowID  WorkflowID  `json:"wid"`
}

type Signal struct {
	Name       string     `json:"name"`
	WorkflowID WorkflowID `json:"wid"`
	Src        string     `json:"src"`
	Filter     string     `json:"filter"`
	Active     bool       `json:"active"`
}
