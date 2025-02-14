-- /postgres/init/03-address.sql
\c addresses;

-- Create UBS table
CREATE TABLE IF NOT EXISTS ubs (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    address VARCHAR(200) NOT NULL,
    city VARCHAR(100) NOT NULL,
    state VARCHAR(2) NOT NULL,
    cep VARCHAR(8) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create teams table
CREATE TABLE IF NOT EXISTS teams (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    ubs_id INTEGER NOT NULL REFERENCES ubs(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create street_segments table
CREATE TABLE IF NOT EXISTS street_segments (
    id SERIAL PRIMARY KEY,
    street_name VARCHAR(200) NOT NULL,
    original_street_name VARCHAR(200) NOT NULL,
    street_type VARCHAR(50) NOT NULL,
    neighborhood VARCHAR(100) NOT NULL,
    city VARCHAR(100) NOT NULL,
    state VARCHAR(2) NOT NULL,
    start_number INTEGER,
    end_number INTEGER,
    cep_prefix VARCHAR(5),
    even_odd VARCHAR(4),
    team_id INTEGER NOT NULL REFERENCES teams(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create trigram index for street name searches
CREATE INDEX IF NOT EXISTS idx_street_segments_street_name_trgm 
ON street_segments USING gin (street_name gin_trgm_ops);

-- Create update triggers
CREATE TRIGGER update_ubs_updated_at
    BEFORE UPDATE ON ubs
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_teams_updated_at
    BEFORE UPDATE ON teams
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_street_segments_updated_at
    BEFORE UPDATE ON street_segments
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
