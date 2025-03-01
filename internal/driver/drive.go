package driver

import (
	"context"
	"errors"
	"fmt"

	"go.temporal.io/api/serviceerror"
	"go.temporal.io/sdk/client"

	"go.autokitteh.dev/demodriver/internal/types"
)

func (d *driver) Drive(ctx context.Context, src string, data any) error {
	return errors.Join(
		d.driveTriggers(ctx, src, data),
		d.driveSignals(ctx, src, data),
	)
}

func (d *driver) driveTriggers(ctx context.Context, src string, data any) error {
	l := d.l.With("src", src)

	var errs []error

	for _, t := range d.triggers {
		if t.trigger.Src == src {
			l := l.With("trigger", t.trigger)

			l.Debug("considering trigger")

			ok, err := evalFilter(t.filter, t.trigger.Src, data)
			if err != nil {
				l.Error("evaluate filter error", "err", err)
				errs = append(errs, err)
				continue
			}

			if !ok {
				l.Debug("filter did not match")
				continue
			}

			l.Info("driving matching trigger")

			wr, err := d.startWorkflow(ctx, t.trigger, data)
			if err != nil {
				l.Error("start workflow error", "err", err)
				errs = append(errs, err)
			}

			l.Info("started workflow", "wr", wr)

			if err := d.db.CreateWorkflow(ctx, &types.Workflow{
				TriggerName: t.trigger.Name,
				WorkflowID:  types.WorkflowID(wr.GetID()),
			}); err != nil {
				l.Error("create workflow error", "err", err)
				errs = append(errs, err)
			}
		}
	}

	return errors.Join(errs...)
}

func (d *driver) startWorkflow(ctx context.Context, t types.Trigger, data any) (client.WorkflowRun, error) {
	return d.client.ExecuteWorkflow(
		ctx,
		client.StartWorkflowOptions{
			TaskQueue: t.QueueName,
			Memo: map[string]any{
				"trigger_name": t.Name,
			},
		},
		t.WorkflowType,
		t.Name,
		data,
	)
}

func (d *driver) driveSignals(ctx context.Context, src string, data any) error {
	l := d.l.With("src", src)

	sigs, err := d.db.ListSignals(ctx, src, "")
	if err != nil {
		return fmt.Errorf("db.ListSignals: %w", err)
	}

	var errs []error

	for _, sig := range sigs {
		l := l.With("signal", sig)

		f, err := parseFilter(sig.Filter)
		if err != nil {
			l.Error("parse filter error", "err", err)
			continue
		}

		l.Debug("considering signal")

		ok, err := evalFilter(f, sig.Src, data)
		if err != nil {
			l.Error("evaluate filter error", "err", err)
			errs = append(errs, err)
			continue
		}

		if !ok {
			l.Debug("filter did not match")
			continue
		}

		l.Info("driving matching signal")

		if err := d.client.SignalWorkflow(
			ctx,
			string(sig.WorkflowID),
			"",
			sig.Name,
			data,
		); err != nil {
			if errors.Is(err, &serviceerror.NotFound{}) {
				l.Info("signal workflow not found, deactivating all signals for workflow")
				if err := d.db.DeactivateSignals(ctx, sig.WorkflowID); err != nil {
					l.Error("deactivate signals error", "err", err)
				}
				continue
			}

			l.Error("signal workflow error", "err", err)
		}
	}

	return errors.Join(errs...)
}
