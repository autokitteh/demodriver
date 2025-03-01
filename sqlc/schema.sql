CREATE TABLE workflows (
    wid         text        PRIMARY KEY,
    tname       text        NOT NULL
);

CREATE INDEX idx_workflows_trigger_id on workflows(tname);

CREATE TABLE signals (
    name        text        NOT NULL,
    wid         text        NOT NULL,
    src         text        NOT NULL,
    filter      text,
    active      boolean     NOT NULL
);

CREATE INDEX idx_signals_wid on signals(wid);
CREATE INDEX idx_signals_wid_src on signals(wid, src);
CREATE INDEX idx_signals_wid_active on signals(wid, active);
CREATE INDEX idx_signals_wid_src_active on signals(wid, src, active);
