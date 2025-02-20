-- Enable the pg_trgm extension for fuzzy search
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Create UBS (Basic Health Unit) table
CREATE TABLE IF NOT EXISTS ubs (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    address VARCHAR(200) NOT NULL,
    city VARCHAR(100) NOT NULL,
    state CHAR(2) NOT NULL,
    cep CHAR(8) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create Teams table
CREATE TABLE IF NOT EXISTS teams (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    ubs_id INTEGER NOT NULL REFERENCES ubs(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create Street Segments table
CREATE TABLE IF NOT EXISTS street_segments (
    id SERIAL PRIMARY KEY,
    street_name VARCHAR(200) NOT NULL,
    original_street_name VARCHAR(200) NOT NULL,
    street_type VARCHAR(50) NOT NULL,
    neighborhood VARCHAR(100) NOT NULL,
    city VARCHAR(100) NOT NULL,
    state CHAR(2) NOT NULL,
    start_number INTEGER NOT NULL,
    end_number INTEGER NOT NULL,
    cep_prefix CHAR(5),
    even_odd VARCHAR(4) NOT NULL CHECK (even_odd IN ('even', 'odd', 'all')),
    team_id INTEGER NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_ubs_city_state ON ubs(city, state);
CREATE INDEX IF NOT EXISTS idx_teams_ubs_id ON teams(ubs_id);
CREATE INDEX IF NOT EXISTS idx_street_segments_team_id ON street_segments(team_id);
CREATE INDEX IF NOT EXISTS idx_street_segments_city_state ON street_segments(city, state);

-- Create trigram index for street name fuzzy search
CREATE INDEX IF NOT EXISTS idx_street_segments_street_name_trgm 
ON street_segments USING gin (street_name gin_trgm_ops);

-- Create function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for updated_at
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

-- Add check constraints
ALTER TABLE street_segments 
    ADD CONSTRAINT check_start_end_number 
    CHECK (start_number <= end_number);

ALTER TABLE ubs
    ADD CONSTRAINT check_state_length
    CHECK (LENGTH(state) = 2);

ALTER TABLE street_segments
    ADD CONSTRAINT check_state_length_segments
    CHECK (LENGTH(state) = 2);

-- Add unique constraints where appropriate
ALTER TABLE teams
    ADD CONSTRAINT unique_team_name_per_ubs
    UNIQUE (name, ubs_id);

-- Add comment documentation
COMMENT ON TABLE ubs IS 'Basic Health Units (Unidades Básicas de Saúde)';
COMMENT ON TABLE teams IS 'Healthcare teams associated with UBS';
COMMENT ON TABLE street_segments IS 'Street segments assigned to healthcare teams';

COMMENT ON COLUMN street_segments.even_odd IS 'Indicates whether the segment covers even numbers, odd numbers, or all numbers (values: even, odd, all)';
COMMENT ON COLUMN street_segments.cep_prefix IS 'First 5 digits of the postal code';
COMMENT ON COLUMN street_segments.original_street_name IS 'Original street name before normalization';
