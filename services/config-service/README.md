# Config Service

Provides application configuration and feature flags. Exposes public bootstrap data for the frontend to initialize before authentication and protected admin endpoints for managing settings and flags.

## Public Endpoints (no auth)
- GET /api/config/bootstrap – minimal bootstrap (currently only environment)
- GET /api/config/public/settings – list of public settings (future; not required for initial bootstrap)
- GET /api/config/public/flags – list of public feature flags (future; not required for initial bootstrap)

## Admin Endpoints (auth via API Gateway; admin role)
- GET /api/config/admin/settings – list all settings (including private)
- POST /api/config/admin/settings – upsert a setting { key, value, is_public, description }
- GET /api/config/admin/flags – list all feature flags
- POST /api/config/admin/flags – upsert a flag { name, enabled, is_public, description }

## Environment Variables
See `.env.example` for available configuration. Database defaults to `config_service_db`.

## Running Locally
```
go mod tidy
go run .
```
