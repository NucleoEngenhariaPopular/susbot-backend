# PostgreSQL Configuration

This directory contains the configuration files and initialization scripts for PostgreSQL databases used in the SUSBot system.

## Directory Structure

```
postgres/
├── init/
│   ├── 01-init.sql       # Database creation and extensions
│   ├── 02-users.sql      # User API tables
│   └── 03-address.sql    # Address API tables
├── postgresql.conf       # PostgreSQL configuration
└── README.md            # This file
```

## Initialization Scripts

The initialization scripts run in alphabetical order when the container starts for the first time:

1. `01-init.sql`: Creates the databases and required extensions
   - Creates `users` and `addresses` databases
   - Installs the `pg_trgm` extension for fuzzy text search

2. `02-users.sql`: Sets up the User API schema
   - Creates the `users` table
   - Sets up indexes for CPF searches
   - Creates triggers for `updated_at` timestamps

3. `03-address.sql`: Sets up the Address API schema
   - Creates `ubs`, `teams`, and `street_segments` tables
   - Sets up foreign key relationships
   - Creates trigram indexes for street name searches

## Configuration

The `postgresql.conf` file contains optimized settings for the SUSBot use case:

- Connection settings
- Memory allocation
- Logging configuration
- Localization settings

## Using Custom Configuration

To use a custom configuration:

1. Modify the configuration files as needed
2. Mount the configuration in docker-compose.yaml:

```yaml
services:
  postgres:
    volumes:
      - ./docker/databases/postgres/postgresql.conf:/etc/postgresql/postgresql.conf
      - ./docker/databases/postgres/init:/docker-entrypoint-initdb.d
    command: ["postgres", "-c", "config_file=/etc/postgresql/postgresql.conf"]
```

## Database Schema

### Users Database

```sql
users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    cpf VARCHAR(11) UNIQUE NOT NULL,
    date_of_birth DATE NOT NULL,
    phone_number VARCHAR(20),
    street_name VARCHAR(200) NOT NULL,
    street_number VARCHAR(20) NOT NULL,
    complement VARCHAR(100),
    neighborhood VARCHAR(100) NOT NULL,
    city VARCHAR(100) NOT NULL,
    state VARCHAR(2) NOT NULL,
    cep VARCHAR(8) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE
)
```

### Address Database

```sql
ubs (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    address VARCHAR(200) NOT NULL,
    city VARCHAR(100) NOT NULL,
    state VARCHAR(2) NOT NULL,
    cep VARCHAR(8) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE
)

teams (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    ubs_id INTEGER NOT NULL REFERENCES ubs(id),
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE
)

street_segments (
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
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE
)
```

## Maintenance

### Backup

To backup the databases:

```bash
docker exec -t postgres pg_dumpall -c -U postgres > dump_`date +%d-%m-%Y"_"%H_%M_%S`.sql
```

### Restore

To restore from a backup:

```bash
cat your_dump.sql | docker exec -i postgres psql -U postgres
```

### Access via Adminer

Adminer is available at `http://localhost:8084`:

- System: PostgreSQL
- Server: postgres
- Username: postgres
- Password: postgres
- Database: (select needed database)
