CREATE TABLE metrics (
    id SERIAL PRIMARY KEY,
    m_id VARCHAR(255) NOT NULL,
    m_type VARCHAR(255) NOT NULL,
    delta INTEGER,
    value double precision
);

CREATE INDEX idx_id ON metrics(id);
CREATE INDEX idx_m_id ON metrics(m_id);
