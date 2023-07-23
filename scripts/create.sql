CREATE TABLE events(
    id uuid PRIMARY KEY UNIQUE DEFAULT gen_random_uuid() NOT NULL,
    name VARCHAR NOT NULL,
    description VARCHAR NOT NULL,
    type VARCHAR NOT NULL ,
    event VARCHAR NOT NULL,
    read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX updated_at ON events (updated_at);
