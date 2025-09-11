CREATE TABLE metrics (
    id SERIAL PRIMARY KEY,
    MType VARCHAR(255) NOT NULL,
    Delta INTEGER,
    Value double precision
);
