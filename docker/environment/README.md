# Environment Variables Documentation

This directory contains environment configuration files for different deployment scenarios.

## Files

- `.env.example`: Template file with all possible variables and their descriptions
- `.env.development`: Configuration for development environment
- `.env.test`: Configuration for testing environment

## Usage

1. For development:

```bash
cp .env.example .env.development
# Edit .env.development with your values
docker compose --env-file docker/environment/.env.development up
```

2. For testing:

```bash
cp .env.example .env.test
# Edit .env.test with your values
docker compose --env-file docker/environment/.env.test up
```

## Required Variables

### API Configuration

GATEWAY_PORT: Port for the main API gateway
USER_API_PORT: Port for the user management API
CONVERSATION_API_PORT: Port for the conversation management API
ADDRESS_API_PORT: Port for the address management API

### Database Configuration

POSTGRES_*: PostgreSQL connection settings
MONGO_*: MongoDB connection settings

### External Services

TWILIO_SID: Your Twilio account SID
TWILIO_AUTH_TOKEN: Your Twilio authentication token

#### Optional Variables

Development

DEBUG: Enable debug mode
LOG_LEVEL: Logging detail level
ENABLE_SWAGGER: Enable Swagger documentation

#### Testing

TEST_MODE: Enable test mode
ENABLE_MOCK_SERVICES: Use mock services instead of real ones

You would then update your docker-compose.yaml to use these environment files:

```yaml
services:
  gateway:
    env_file:
      - ./docker/environment/.env.development
    # ... rest of the configuration

  user-api:
    env_file:
      - ./docker/environment/.env.development
    # ... rest of the configuration

  # ... other services
```

For testing, you would run:

```bash
docker compose --env-file docker/environment/.env.test up
```
