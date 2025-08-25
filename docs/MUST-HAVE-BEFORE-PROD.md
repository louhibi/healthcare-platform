# Must-Haves Before Production

This ticket list is based on a code audit of the current repository (api-gateway, user-service, patient-service, appointment-service). It focuses on must-fix items to ensure multi-tenant safety, auth integrity, data correctness, timezone discipline, and basic operational readiness.

## 0) High-Priority Summary (Do These First)
- [ ] Enforce tenant isolation in every data access path (service-layer fixes now; DB RLS next).
- [ ] Validate JWTs in all services; stop trusting `X-*` headers without verification.
- [ ] Remove timezone fallbacks; persist and serve only UTC; require entity timezone everywhere.
- [ ] Prevent cross-tenant data access by ID (scope every read/update/delete by `healthcare_entity_id`).
- [ ] Add DB integrity for scheduling (no-overlap constraints) and implement idempotency for writes.
- [ ] Lock down internal endpoints and CORS; add per-entity/user rate limiting.
- [ ] Stop returning PHI in error responses; centralize safe error format.

---

## 1) Multi-Tenancy Enforcement

- [ ] Scope-by-entity for all ID-based operations
  - Risk: Cross-tenant reads/updates/deletes possible when using only `id`.
  - Examples to fix now (include `healthcare_entity_id = ?` in WHERE):
    - `services/appointment-service/appointment_service.go`
      - `GetAppointmentByID` (currently `WHERE id = $1 AND is_active = true`)
      - `UpdateAppointment` / `UpdateAppointmentStatus` / `DeleteAppointment` (id-only filters)
    - `services/patient-service/patient_service.go`
      - `GetPatientByID` (currently `WHERE id = $1 AND is_active = true`)
      - `UpdatePatient` / `Delete` (id-only filters; ensure entity constraint)
  - Acceptance: All reads/updates/deletes that take a resource `id` also require and bind the caller’s `X-Healthcare-Entity-ID` and filter on it.

- [ ] Add DB-level Row-Level Security (next iteration; can ship right after app fixes)
  - Enable RLS and policies per table (`patients`, `appointments`, `doctor_availability`, `rooms`, `users`).
  - Per-request: `SET LOCAL app.entity_id = $entity_id`; policies filter by `current_setting('app.entity_id')`.
  - Acceptance: Queries without explicit entity filters still cannot read/write across tenants.

## 2) Authentication & Header Trust

- [ ] Validate JWTs in patient- and appointment-service
  - Today: services trust `X-User-ID`/`X-Healthcare-Entity-ID` headers (appointment `main.go`, patient `handlers.go`).
  - Change: forward `Authorization` (gateway already does) and validate in each service using the shared secret (short-term) or RS256 JWKS (preferred).
  - Files:
    - Add a small `auth_middleware.go` to patient- and appointment-service mirroring gateway’s JWT checks.
  - Acceptance: A direct call to services without a valid JWT is rejected; headers alone are not sufficient.

- [ ] Protect internal endpoints
  - User-service: `/api/internal/...` routes are public. Appointment-service calls user-service and patient-service with no auth.
  - Add an internal service token (shared secret via env) or mTLS; require `Authorization: Bearer <INTERNAL_TOKEN>` on internal routes.
  - Files: `services/user-service/main.go` (internal groups), appointment/patient service internal call sites.
  - Acceptance: Internal endpoints reject unauthenticated requests; all inter-service calls supply and validate the token.

- [ ] CORS hardening
  - Replace `*` with `ALLOWED_ORIGINS` from env across services (gateway, patient, appointment, user), aligning with frontend domain(s).
  - Acceptance: Only configured origins are allowed in production build.

## 3) Timezone Discipline (No Fallbacks)

- [ ] Remove UTC fallback in timezone conversion
  - Appointment-service `TimezoneConverter` defaults to `UTC`; `GetTimezoneConverter` falls back to `UTC` on user-service failures.
  - Change: Fail closed. If entity timezone cannot be fetched, return 5xx and surface actionable error; never silently fall back.
  - Files: `services/appointment-service/models.go` (TimezoneConverter), `appointment_service.go` (GetTimezoneConverter).
  - Acceptance: No implicit timezone fallback; all times are stored in UTC (`TIMESTAMPTZ`) and presented only after successful entity timezone resolution.

- [ ] Convert remaining TIMESTAMP columns to TIMESTAMPTZ
  - User- and patient-service schemas use `TIMESTAMP` for `created_at/updated_at`.
  - Add migrations to convert to `TIMESTAMPTZ` and ensure UTC semantics.
  - Files: `services/user-service/migrations.go`, `services/patient-service/migrations.go`.
  - Acceptance: All datetime columns across services are `TIMESTAMPTZ` and serialized as ISO8601 UTC (`...Z`).

## 4) Scheduling Integrity & Idempotency

- [ ] Enforce no-overlap at the database level
  - Add Postgres exclusion constraints for appointments (per doctor, per entity; and per room if `room_id` present):
    - `EXCLUDE USING gist (healthcare_entity_id WITH =, doctor_id WITH =, tsrange(date_time, date_time + duration * interval '1 min') WITH &&)`
  - Files: `services/appointment-service/database.go` new migration.
  - Acceptance: Overlapping appointments are rejected atomically by the DB.

- [ ] Add Idempotency-Key support for writes
  - Apply to: appointment booking/create, patient create, doctor availability create.
  - Store dedupe keys (e.g., in DB table with TTL column). Reject duplicates or return first result.
  - Files: gateway middleware + targeted service handlers.
  - Acceptance: Repeated POSTs with same key are safe and return identical results.

## 5) API Safety & Consistency

- [ ] Fix entity-agnostic stats queries
  - Appointment-service `GetAppointmentStats` counts all entities.
  - Add required `healthcare_entity_id` filter throughout stats and aggregate endpoints.
  - Files: `services/appointment-service/appointment_service.go`.
  - Acceptance: All analytics/metrics endpoints are strictly scoped to caller’s entity.

- [ ] Standardize error responses and remove PHI from responses
  - Patient create currently returns `internal_error` and `patient_data` on failures.
  - Adopt a consistent error envelope (code, message, timestamp); do not echo request bodies or sensitive fields.
  - Files: patient `handlers.go`, appointment `handlers.go`, gateway `ErrorResponse` usage.
  - Acceptance: No PHI in error payloads; consistent, minimal error format across services.

- [ ] Rate limiting per user/entity in gateway
  - Current limiter uses IP only.
  - Extend key derivation to include `user_id` and `healthcare_entity_id` where present.
  - Files: `services/api-gateway/rate_limiter.go`.
  - Acceptance: Configurable per-user+entity limits; heavy endpoints protected.

## 6) Security & Operations

- [ ] Secrets hygiene and deployment
  - Enforce non-default `JWT_SECRET` and internal tokens in production; wire via container secrets.
  - Add startup checks that refuse to run with default secrets.
  - Files: service `main.go` or config loaders.
  - Acceptance: Services exit on default/insecure secrets in prod mode.

- [ ] Service-to-service timeouts, retries, and auth
  - Appointment-service calls to user/patient services use `http.Get` without auth or standardized retries.
  - Use a shared HTTP client with timeouts/retries and add internal auth header.
  - Files: appointment `appointment_service.go` inter-service methods.
  - Acceptance: All inter-service calls authenticate and have resilient client config.

- [ ] Observability baseline (tracing + structured logs)
  - Add request IDs/correlation IDs across gateway and services; log structured JSON with entity/user identifiers (anonymized where needed).
  - Acceptance: Each request traceable end-to-end; logs do not contain PHI by default.

## 7) Configuration & CORS

- [ ] Tighten CORS in all services
  - Gateway and services currently set `*` and broad headers. Restrict to known origins and required headers only.
  - Files: `services/*/main.go`, `services/api-gateway/config.go`.
  - Acceptance: Only production frontend origin(s) allowed; preflight correct; unnecessary headers removed.

## 8) Documentation & Validation

- [ ] Document JWT and internal auth model
  - Short-term: HMAC shared secret validation in all services; long-term: RS256 + JWKS.
  - Acceptance: README sections updated (gateway, services) describing auth verification and header propagation.

- [ ] Add OpenAPI specs for each service (minimum critical endpoints)
  - At least for: auth/login, patients CRUD, appointments CRUD/booking.
  - Acceptance: Specs exist and are checked in; CI/lint optional follow-up.

---

## Pointers to Fixes (by file)

- api-gateway
  - `auth_middleware.go`: HMAC JWT; OK if secrets are strong; ensure token forwarded; consider adding `aud`/`scope` checks.
  - `rate_limiter.go`: Add composite key (IP+user+entity) and route-sensitive quotas.
  - `config.go`: Restrict `ALLOWED_ORIGINS`; ensure headers minimal.

- user-service
  - `auth_service.go`: HMAC today; ensure strong secret in prod. Plan RS256/JWKS later.
  - `migrations.go`: Convert timestamps to `TIMESTAMPTZ`; ensure UTC.
  - `main.go`: Protect `/api/internal/...` with internal token.

- patient-service
  - `handlers.go`: Implement JWT validation middleware; remove PHI from error payloads.
  - `patient_service.go`: Scope by entity for ID-based ops; add entity to WHERE clauses; convert timestamps to `TIMESTAMPTZ` via migration.
  - `main.go`: CORS hardening.

- appointment-service
  - `handlers.go`: Enforce entity scoping on all operations (reads/updates/deletes/status changes).
  - `appointment_service.go`: Fix entity-agnostic stats; add inter-service auth and resilient HTTP client; remove timezone fallback usage; add idempotency hooks.
  - `database.go`: New migration for exclusion constraints; ensure TIMESTAMPTZ everywhere (many already addressed).
  - `models.go`: Remove `UTC` default fallback from `TimezoneConverter`.

---

## Nice-to-Have (Post-MVP hardening)
- RS256 + JWKS with key rotation; gateway and services verify; drop shared secret.
- Postgres RLS fully enabled; middleware sets `app.entity_id` per request.
- Problem Details (RFC 9457) error format; versioned APIs `/api/v1/...`.
- Cursor-based pagination for large lists.
- Outbox pattern for domain events; eventually Kafka.

---

## Verification Checklist (Go/No-Go)
- [ ] Direct service calls without JWT are rejected.
- [ ] Attempt to read/update/delete a resource from another tenant is rejected.
- [ ] Booking API cannot create overlapping appointments (DB enforced).
- [ ] All datetimes in DB are `TIMESTAMPTZ` and API returns ISO8601 UTC.
- [ ] CORS is restricted to production origins.
- [ ] Logs contain no PHI; errors are minimal and consistent.
- [ ] Internal endpoints require internal auth token.
- [ ] Idempotency prevents duplicate create/book calls.

