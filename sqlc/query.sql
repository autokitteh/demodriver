-- name: CreateWorkflow :exec
INSERT INTO workflows(wid, tname) VALUES (?, ?);

-- name: GetWorkflow :one
SELECT * FROM workflows WHERE wid = ?;

-- name: ListWorkflows :many
SELECT * FROM workflows WHERE tname = ?;

-- name: ListAllWorkflows :many
SELECT * FROM workflows;

-- name: CreateSignal :exec
INSERT INTO signals(name, wid, src, filter, active) VALUES (?, ?, ?, ?, ?);

-- name: ListAllSignals :many
SELECT * FROM signals WHERE active = True;

-- name: ListSignalsForSource :many
SELECT * FROM signals WHERE src = ? AND active = True;

-- name: ListSignalsForWorkflow :many
SELECT * FROM signals WHERE wid = ? AND active = True;

-- name: DeactivateSignals :exec
UPDATE signals SET active = false WHERE wid = ? AND src = ?;
