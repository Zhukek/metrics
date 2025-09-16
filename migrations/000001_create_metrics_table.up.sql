CREATE TABLE metrics (
    id VARCHAR(255) NOT NULL,
    m_type VARCHAR(255) NOT NULL,
    delta INTEGER,
    value double precision
);

CREATE INDEX idx_id ON metrics(id);
