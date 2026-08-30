CREATE TABLE task_statistics (
    id VARCHAR(50) PRIMARY KEY DEFAULT concat('tst_', gen_random_uuid()),
    customer_id VARCHAR(50) NOT NULL,
    task_id VARCHAR(50) NOT NULL,
    event_type VARCHAR(20) NOT NULL CHECK (event_type IN ('created', 'updated', 'completed', 'deleted')),
    occurred_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_task_statistics_customer_occurred_at
    ON task_statistics (customer_id, occurred_at DESC, id DESC);
