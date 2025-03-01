package driver

import (
	"context"
	"fmt"
	"log/slog"

	"go.temporal.io/sdk/client"

	"go.autokitteh.dev/demodriver/internal/db"
	"go.autokitteh.dev/demodriver/internal/types"
)

type DriveFunc func(ctx context.Context, src string, data any) error

type Driver interface {
	Drive(ctx context.Context, src string, data any) error

	RegisterSignal(context.Context, types.Signal) error

	ListTriggers(context.Context) []*types.Trigger
	ListWorkflows(_ context.Context, tname types.TriggerName) ([]*types.Workflow, error)
	ListSignals(_ context.Context, src string, wid types.WorkflowID) ([]*types.Signal, error)
}

type driver struct {
	cfg    *config
	l      *slog.Logger
	client client.Client
	db     db.DB

	triggers map[types.TriggerName]*trigger
}

func (d *driver) RegisterSignal(ctx context.Context, s types.Signal) error {
	d.l.Info("registering signal", "s", s)

	if _, err := parseFilter(s.Filter); err != nil {
		return fmt.Errorf("filter: %w", err)
	}

	return d.db.CreateSignal(ctx, &s)
}

func (d *driver) ListTriggers(ctx context.Context) (ts []*types.Trigger) {
	ts = make([]*types.Trigger, 0, len(d.triggers))
	for _, t := range d.triggers {
		ts = append(ts, &t.trigger)
	}
	return ts
}

func (d *driver) ListWorkflows(ctx context.Context, tname types.TriggerName) ([]*types.Workflow, error) {
	return d.db.ListWorkflows(ctx, tname)
}

func (d *driver) ListSignals(ctx context.Context, src string, wid types.WorkflowID) ([]*types.Signal, error) {
	return d.db.ListSignals(ctx, src, wid)
}
