# MongoDB Configuration

This directory contains the configuration files and initialization scripts for MongoDB used in the SUSBot system, primarily for conversation management.

## Directory Structure

```
mongodb/
├── init-mongo.js    # Database and user initialization script
├── mongod.conf      # MongoDB configuration file
└── README.md        # This file
```

## Configuration Files

### init-mongo.js

This script runs when MongoDB starts for the first time:

- Creates the root user
- Initializes the conversations database
- Sets up collection schemas and validation
- Creates necessary indexes

### mongod.conf

Contains MongoDB server configuration:

- Network settings
- Security configuration
- Storage settings
- Logging preferences

## Database Schema

### Conversations Collection

```javascript
{
  _id: ObjectId,
  user_id: String,         // User identifier
  start_time: Date,        // Conversation start
  end_time: Date,          // Conversation end (null if active)
  messages: [{
    user_id: String,       // User identifier
    sender: String,        // "user" or "system"
    text: String,          // Message content
    timestamp: Date        // Message timestamp
  }]
}
```

## Validation Rules

The collection enforces the following validation rules:

- Required fields: `user_id`, `start_time`, `messages`
- `messages` must be an array with required fields:
  - `user_id` (string)
  - `sender` (string)
  - `text` (string)
  - `timestamp` (date)

## Indexes

The following indexes are created automatically:

- `user_id`: For quick user conversation lookup
- `start_time`: For chronological sorting and queries

## Using Custom Configuration

To use a custom configuration:

1. Modify the configuration files as needed
2. Mount the configuration in docker-compose.yaml:

```yaml
services:
  mongo:
    volumes:
      - ./docker/databases/mongodb/mongod.conf:/etc/mongod.conf
      - ./docker/databases/mongodb/init-mongo.js:/docker-entrypoint-initdb.d/init-mongo.js:ro
    command: ["mongod", "--config", "/etc/mongod.conf"]
```

## Maintenance

### Backup

To backup the database:

```bash
docker exec -t mongo mongodump --username root --password example --authenticationDatabase admin --db my_database --out /dump
```

### Restore

To restore from a backup:

```bash
docker exec -t mongo mongorestore --username root --password example --authenticationDatabase admin /dump
```

### Access via Mongo Express

Mongo Express is available at `http://localhost:8085`:

- Username: root
- Password: example

## Security Notes

1. Default Credentials:
   - The default root username and password are set in the environment files
   - Change these in production environments

2. Authentication:
   - Authentication is enabled by default
   - All connections require valid credentials

3. Network:
   - MongoDB listens on port 27017
   - Only container network access is allowed by default

## Troubleshooting

Common issues and solutions:

1. Connection Failed

```bash
# Check MongoDB status
docker logs mongo

# Verify network connectivity
docker exec -it mongo mongosh --eval "db.runCommand({ping: 1})"
```

2. Authentication Failed

```bash
# Verify credentials
docker exec -it mongo mongosh --eval "db.auth('root', 'example')"
```

3. Data Persistence

```bash
# Check volume mounting
docker volume ls
docker volume inspect susbot-backend_mongodb_data
```

## Monitoring

1. Via Mongo Express:
   - Database size
   - Collection statistics
   - Index usage

2. Via MongoDB Shell:

```javascript
// Check database stats
db.stats()

// Check collection stats
db.conversations.stats()

// Monitor connections
db.serverStatus().connections
```
