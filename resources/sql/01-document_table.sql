\c docu_db docu_user;

CREATE TABLE docu_schema.documents (
    id SERIAL PRIMARY KEY,
    document_id VARCHAR(255) NOT NULL UNIQUE,
    hash VARCHAR(255) NOT NULL UNIQUE,
    tag VARCHAR(255) NOT NULL,
    status INTEGER NOT NULL,
    analyzed_at TIMESTAMP WITHOUT TIME ZONE,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

-- Indexes
CREATE INDEX idx_document_id ON docu_schema.documents(document_id);
CREATE INDEX idx_hash ON docu_schema.documents(hash);
CREATE INDEX idx_status ON docu_schema.documents(status);
CREATE INDEX idx_analyzed_at ON docu_schema.documents(analyzed_at);

-- Check Constraints 
ALTER TABLE docu_schema.documents ADD CONSTRAINT chk_status CHECK (status IN (0, 1, 2));