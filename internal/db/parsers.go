package db

import (
	"go.autokitteh.dev/demodriver/internal/types"
	"go.autokitteh.dev/demodriver/sqlc/sqlcgen"
)

func parseWorkflow(w sqlcgen.Workflow) (*types.Workflow, error) {
	return &types.Workflow{
		TriggerName: types.TriggerName(w.Tname),
		WorkflowID:  types.WorkflowID(w.Wid),
	}, nil
}

func parseSignal(s sqlcgen.Signal) (*types.Signal, error) {
	return &types.Signal{
		WorkflowID: types.WorkflowID(s.Wid),
		Src:        s.Src,
		Filter:     s.Filter.String,
		Active:     s.Active,
		Name:       s.Name,
	}, nil
}
