CREATE TABLE metrics (
    id VARCHAR(255) NOT NULL,
    m_type VARCHAR(255) NOT NULL,
    delta INTEGER,
    value double precision,
    PRIMARY KEY (id, m_type)
);

CREATE INDEX idx_id ON metrics(id);
