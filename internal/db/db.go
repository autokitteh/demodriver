package db

import (
	"context"

	_ "modernc.org/sqlite"

	"go.autokitteh.dev/demodriver/internal/types"
)

type DB interface {
	Setup(context.Context) error
	Close() error

	CreateWorkflow(context.Context, *types.Workflow) error
	GetWorkflow(context.Context, types.WorkflowID) (*types.Workflow, error) // nil, nil if not found.
	ListWorkflows(context.Context, types.TriggerName) ([]*types.Workflow, error)

	CreateSignal(context.Context, *types.Signal) error
	DeactivateSignals(context.Context, types.WorkflowID) error
	ListSignals(_ context.Context, src string, wid types.WorkflowID) ([]*types.Signal, error)
}
