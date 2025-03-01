package db

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"go.autokitteh.dev/demodriver/internal/types"
	"go.autokitteh.dev/demodriver/sqlc"
	"go.autokitteh.dev/demodriver/sqlc/sqlcgen"
)

type db struct {
	cfg *config
	l   *slog.Logger
	db  *sql.DB
	q   *sqlcgen.Queries
}

func ignoreNotFound(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}

	return err
}

func (db *db) Setup(ctx context.Context) error {
	_, err := db.db.ExecContext(ctx, sqlc.Schema)
	return err
}

func (db *db) Close() error { return db.db.Close() }

func (db *db) CreateWorkflow(ctx context.Context, w *types.Workflow) error {
	return db.q.CreateWorkflow(ctx, sqlcgen.CreateWorkflowParams{
		Tname: string(w.TriggerName),
		Wid:   string(w.WorkflowID),
	})
}

func (db *db) GetWorkflow(ctx context.Context, wid types.WorkflowID) (*types.Workflow, error) {
	w, err := db.q.GetWorkflow(ctx, string(wid))
	if err != nil {
		return nil, ignoreNotFound(err)
	}

	return parseWorkflow(w)
}

func (db *db) ListWorkflows(ctx context.Context, tname types.TriggerName) ([]*types.Workflow, error) {
	var (
		rs  []sqlcgen.Workflow
		err error
	)

	if tname == "" {
		rs, err = db.q.ListAllWorkflows(ctx)
	} else {
		rs, err = db.q.ListWorkflows(ctx, string(tname))
	}

	if err != nil {
		return nil, err
	}

	ws := make([]*types.Workflow, len(rs))
	for i, r := range rs {
		if ws[i], err = parseWorkflow(r); err != nil {
			return nil, err
		}
	}

	return ws, nil
}

func (db *db) CreateSignal(ctx context.Context, s *types.Signal) error {
	return db.q.CreateSignal(ctx, sqlcgen.CreateSignalParams{
		Wid:    string(s.WorkflowID),
		Src:    s.Src,
		Filter: sql.NullString{String: s.Filter, Valid: s.Filter != ""},
		Active: true,
		Name:   s.Name,
	})
}

func (db *db) ListSignals(ctx context.Context, src string, wid types.WorkflowID) ([]*types.Signal, error) {
	var (
		rs  []sqlcgen.Signal
		err error
	)

	if wid != "" {
		rs, err = db.q.ListSignalsForWorkflow(ctx, string(wid))
	} else if src != "" {
		rs, err = db.q.ListSignalsForSource(ctx, src)
	} else {
		rs, err = db.q.ListAllSignals(ctx)
	}

	if err != nil {
		return nil, err
	}

	sigs := make([]*types.Signal, 0, len(rs))
	for _, r := range rs {
		if src != "" && r.Src != src {
			continue
		}

		sig, err := parseSignal(r)
		if err != nil {
			return nil, err
		}

		sigs = append(sigs, sig)
	}

	return sigs, nil
}

func (db *db) DeactivateSignals(ctx context.Context, wid types.WorkflowID) error {
	return db.q.DeactivateSignals(ctx, sqlcgen.DeactivateSignalsParams{
		Wid: string(wid),
	})
}
