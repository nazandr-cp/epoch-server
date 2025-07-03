-- Initialize database for epoch-server
-- This script runs when the PostgreSQL container starts for the first time

-- Create database if it doesn't exist (handled by POSTGRES_DB env var)
-- The database is automatically created by the postgres image

-- Create schema for epoch management
CREATE SCHEMA IF NOT EXISTS epochs;

-- Create table for storing epoch states
CREATE TABLE IF NOT EXISTS epochs.epoch_states (
    id SERIAL PRIMARY KEY,
    epoch_id BIGINT NOT NULL UNIQUE,
    state VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create table for storing subsidy snapshots
CREATE TABLE IF NOT EXISTS epochs.subsidy_snapshots (
    id SERIAL PRIMARY KEY,
    epoch_id BIGINT NOT NULL,
    user_address VARCHAR(42) NOT NULL,
    collection_address VARCHAR(42) NOT NULL,
    subsidy_amount DECIMAL(36, 18) NOT NULL,
    merkle_proof TEXT[],
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (epoch_id) REFERENCES epochs.epoch_states(epoch_id)
);

-- Create table for tracking distribution events
CREATE TABLE IF NOT EXISTS epochs.distribution_events (
    id SERIAL PRIMARY KEY,
    epoch_id BIGINT NOT NULL,
    transaction_hash VARCHAR(66),
    block_number BIGINT,
    gas_used BIGINT,
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (epoch_id) REFERENCES epochs.epoch_states(epoch_id)
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_epoch_states_epoch_id ON epochs.epoch_states(epoch_id);
CREATE INDEX IF NOT EXISTS idx_subsidy_snapshots_epoch_id ON epochs.subsidy_snapshots(epoch_id);
CREATE INDEX IF NOT EXISTS idx_subsidy_snapshots_user_address ON epochs.subsidy_snapshots(user_address);
CREATE INDEX IF NOT EXISTS idx_distribution_events_epoch_id ON epochs.distribution_events(epoch_id);
CREATE INDEX IF NOT EXISTS idx_distribution_events_tx_hash ON epochs.distribution_events(transaction_hash);

-- Create function to update the updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger to automatically update the updated_at column
CREATE TRIGGER update_epoch_states_updated_at 
    BEFORE UPDATE ON epochs.epoch_states 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Grant permissions to the epoch_user
GRANT USAGE ON SCHEMA epochs TO epoch_user;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA epochs TO epoch_user;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA epochs TO epoch_user;