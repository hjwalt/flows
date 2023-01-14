DROP TABLE IF EXISTS public.flows_state;

CREATE TABLE IF NOT EXISTS public.flows_state
(
    persistence_id        VARCHAR(255)  NOT NULL,
    content               BYTEA         NULL,
    results               BYTEA         NULL,
    internal              BYTEA         NULL,
    created_timestamp_ms  BIGINT,
    updated_timestamp_ms  BIGINT,
    PRIMARY KEY(persistence_id)
);
