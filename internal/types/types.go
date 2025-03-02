package types

type (
	TriggerName string
	WorkflowID  string
)

type Trigger struct {
	Name         TriggerName `json:"name" yaml:"name"`
	Src          string      `json:"src" yaml:"src"`
	Filter       string      `json:"filter" yaml:"filter"`
	QueueName    string      `json:"qname" yaml:"qname"`
	WorkflowType string      `json:"wtype" yaml:"wtype"`
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
