# Troubleshooting Guide

Comprehensive troubleshooting guide for the Healthcare Management Platform.

## Quick Diagnosis

### System Health Check

```bash
# Quick health check for all services
curl http://localhost:8080/health  # API Gateway
curl http://localhost:8081/health  # User Service
curl http://localhost:8082/health  # Patient Service
curl http://localhost:8083/health  # Appointment Service
curl http://localhost:3000         # Frontend

# Docker container status
docker compose ps

# Check logs for errors
docker compose logs --tail=50
```

### Common Symptoms & Quick Fixes

| Symptom | Quick Check | Quick Fix |
|---------|-------------|-----------|
| 502 Bad Gateway | Service down | `docker compose restart <service>` |
| 401 Unauthorized | JWT token issue | Check JWT_SECRET consistency |
| 429 Too Many Requests | Rate limiting | Wait or increase rate limits |
| 500 Internal Server Error | Service error | Check service logs |
| Database connection failed | DB not running | `docker compose up -d <db-service>` |

## Service-Specific Issues

### API Gateway Issues

#### Issue: 502 Bad Gateway
**Symptoms**: API Gateway returns 502 for all requests

**Diagnosis**:
```bash
# Check API Gateway logs
docker compose logs api-gateway

# Check backend service health
curl http://localhost:8081/health
curl http://localhost:8082/health
curl http://localhost:8083/health

# Check service URLs in gateway config
docker compose exec api-gateway env | grep SERVICE_URL
```

**Solutions**:
```bash
# 1. Restart backend services
docker compose restart user-service patient-service appointment-service

# 2. Check network connectivity
docker compose exec api-gateway ping user-service
docker compose exec api-gateway ping patient-service

# 3. Verify environment variables
docker compose exec api-gateway env | grep -E "(USER|PATIENT|APPOINTMENT)_SERVICE_URL"
```

#### Issue: Rate Limiting Too Aggressive
**Symptoms**: Legitimate requests getting 429 errors

**Diagnosis**:
```bash
# Check current rate limit settings
curl http://localhost:8080/stats

# Check gateway configuration
docker compose exec api-gateway env | grep RATE_LIMIT
```

**Solutions**:
```bash
# Temporarily increase rate limits
export RATE_LIMIT_RPM=500
export RATE_LIMIT_BURST=50
docker compose restart api-gateway

# Or modify docker compose.yml and restart
```

### User Service Issues

#### Issue: JWT Token Validation Fails
**Symptoms**: All authenticated requests return 401

**Diagnosis**:
```bash
# Check JWT secret consistency
docker compose exec user-service env | grep JWT_SECRET
docker compose exec api-gateway env | grep JWT_SECRET

# Test token generation
curl -X POST http://localhost:8081/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@healthcare.local","password":"admin123"}'
```

**Solutions**:
```bash
# 1. Ensure JWT_SECRET is identical across services
# Edit .env files to have same JWT_SECRET
# Restart affected services
docker compose restart user-service api-gateway

# 2. Check token format
# Tokens should start with "Bearer " in Authorization header
```

#### Issue: Database Connection Failed
**Symptoms**: User service can't connect to database

**Diagnosis**:
```bash
# Check database status
docker compose ps user-db

# Check database logs
docker compose logs user-db

# Test database connection
docker compose exec user-service nc -zv user-db 5432

# Check database credentials
docker compose exec user-service env | grep DB_
```

**Solutions**:
```bash
# 1. Restart database
docker compose restart user-db

# 2. Check database initialization
docker compose exec user-db psql -U postgres -d user_service_db -c "\dt"

# 3. Reset database (CAUTION: Data loss)
docker compose stop user-db
docker volume rm healthcare-platform_user_db_data
docker compose up -d user-db
```

### Patient Service Issues

#### Issue: Search Performance Slow
**Symptoms**: Patient search takes >5 seconds

**Diagnosis**:
```bash
# Check database query performance
docker compose exec patient-db psql -U postgres -d patient_service_db -c "
EXPLAIN ANALYZE 
SELECT * FROM patients 
WHERE LOWER(first_name) LIKE '%alice%' 
   OR LOWER(last_name) LIKE '%alice%';"

# Check index usage
docker compose exec patient-db psql -U postgres -d patient_service_db -c "
SELECT schemaname, tablename, indexname, idx_scan 
FROM pg_stat_user_indexes 
WHERE schemaname = 'public';"
```

**Solutions**:
```bash
# 1. Add missing indexes
docker compose exec patient-db psql -U postgres -d patient_service_db -c "
CREATE INDEX IF NOT EXISTS idx_patients_search 
ON patients USING gin(to_tsvector('english', first_name || ' ' || last_name));"

# 2. Update table statistics
docker compose exec patient-db psql -U postgres -d patient_service_db -c "ANALYZE patients;"

# 3. Check for slow queries
docker compose logs patient-service | grep -i "slow"
```

#### Issue: Email Uniqueness Validation Fails
**Symptoms**: Duplicate emails allowed or false conflicts

**Diagnosis**:
```bash
# Check for duplicate emails
docker compose exec patient-db psql -U postgres -d patient_service_db -c "
SELECT email, COUNT(*) 
FROM patients 
WHERE is_active = true 
GROUP BY email 
HAVING COUNT(*) > 1;"

# Check email validation logic
docker compose logs patient-service | grep -i "email"
```

**Solutions**:
```bash
# 1. Clean up duplicate emails
docker compose exec patient-db psql -U postgres -d patient_service_db -c "
UPDATE patients 
SET is_active = false 
WHERE id NOT IN (
    SELECT MIN(id) 
    FROM patients 
    GROUP BY email
) AND is_active = true;"

# 2. Add unique constraint if missing
docker compose exec patient-db psql -U postgres -d patient_service_db -c "
ALTER TABLE patients ADD CONSTRAINT unique_active_email 
EXCLUDE (email WITH =) WHERE (is_active = true);"
```

### Appointment Service Issues

#### Issue: Conflict Detection Not Working
**Symptoms**: Overlapping appointments allowed

**Diagnosis**:
```bash
# Check for conflicting appointments
docker compose exec appointment-db psql -U postgres -d appointment_service_db -c "
SELECT a1.id, a1.doctor_id, a1.date_time, a1.duration,
       a2.id, a2.date_time, a2.duration
FROM appointments a1, appointments a2
WHERE a1.id < a2.id
  AND a1.doctor_id = a2.doctor_id
  AND a1.is_active = true
  AND a2.is_active = true
  AND a1.status IN ('scheduled', 'confirmed', 'in-progress')
  AND a2.status IN ('scheduled', 'confirmed', 'in-progress')
  AND a1.date_time < a2.date_time + (a2.duration || ' minutes')::interval
  AND a2.date_time < a1.date_time + (a1.duration || ' minutes')::interval;"

# Check conflict detection index
docker compose exec appointment-db psql -U postgres -d appointment_service_db -c "
\d+ appointments"
```

**Solutions**:
```bash
# 1. Add conflict detection index if missing
docker compose exec appointment-db psql -U postgres -d appointment_service_db -c "
CREATE INDEX IF NOT EXISTS idx_appointments_conflict_check 
ON appointments(doctor_id, date_time, duration, status) 
WHERE is_active = true AND status IN ('scheduled', 'confirmed', 'in-progress');"

# 2. Fix existing conflicts by cancelling duplicates
docker compose exec appointment-db psql -U postgres -d appointment_service_db -c "
UPDATE appointments 
SET status = 'cancelled' 
WHERE id IN (
    SELECT a2.id
    FROM appointments a1, appointments a2
    WHERE a1.id < a2.id
      AND a1.doctor_id = a2.doctor_id
      AND a1.is_active = true
      AND a2.is_active = true
      AND a1.status IN ('scheduled', 'confirmed')
      AND a2.status IN ('scheduled', 'confirmed')
      AND a1.date_time < a2.date_time + (a2.duration || ' minutes')::interval
      AND a2.date_time < a1.date_time + (a1.duration || ' minutes')::interval
);"
```

#### Issue: Timezone Issues
**Symptoms**: Appointments show wrong times

**Diagnosis**:
```bash
# Check timezone settings
docker compose exec appointment-service env | grep TZ
docker compose exec appointment-db psql -U postgres -c "SHOW timezone;"

# Check appointment times
docker compose exec appointment-db psql -U postgres -d appointment_service_db -c "
SELECT id, date_time, 
       date_time AT TIME ZONE 'UTC' as utc_time,
       date_time AT TIME ZONE 'America/New_York' as local_time
FROM appointments 
ORDER BY date_time DESC LIMIT 5;"
```

**Solutions**:
```bash
# 1. Set consistent timezone
export TZ=UTC
docker compose restart appointment-service appointment-db

# 2. Convert existing times to UTC if needed
docker compose exec appointment-db psql -U postgres -d appointment_service_db -c "
UPDATE appointments 
SET date_time = date_time AT TIME ZONE 'America/New_York' AT TIME ZONE 'UTC'
WHERE date_time > '2024-01-01';"
```

## Database Issues

### PostgreSQL Connection Issues

#### Issue: Too Many Connections
**Symptoms**: "remaining connection slots are reserved"

**Diagnosis**:
```bash
# Check current connections
docker compose exec user-db psql -U postgres -c "
SELECT datname, usename, count(*) 
FROM pg_stat_activity 
GROUP BY datname, usename;"

# Check max connections
docker compose exec user-db psql -U postgres -c "SHOW max_connections;"
```

**Solutions**:
```bash
# 1. Kill idle connections
docker compose exec user-db psql -U postgres -c "
SELECT pg_terminate_backend(pid) 
FROM pg_stat_activity 
WHERE state = 'idle' AND state_change < NOW() - INTERVAL '10 minutes';"

# 2. Increase max connections (restart required)
# Edit postgresql.conf: max_connections = 200
docker compose restart user-db patient-db appointment-db

# 3. Implement connection pooling in application
```

#### Issue: Database Locks
**Symptoms**: Operations hang or timeout

**Diagnosis**:
```bash
# Check for locks
docker compose exec user-db psql -U postgres -c "
SELECT blocked_locks.pid AS blocked_pid,
       blocked_activity.usename AS blocked_user,
       blocking_locks.pid AS blocking_pid,
       blocking_activity.usename AS blocking_user,
       blocked_activity.query AS blocked_statement,
       blocking_activity.query AS current_statement_in_blocking_process
FROM pg_catalog.pg_locks blocked_locks
JOIN pg_catalog.pg_stat_activity blocked_activity ON blocked_activity.pid = blocked_locks.pid
JOIN pg_catalog.pg_locks blocking_locks 
    ON blocking_locks.locktype = blocked_locks.locktype
    AND blocking_locks.database IS NOT DISTINCT FROM blocked_locks.database
    AND blocking_locks.relation IS NOT DISTINCT FROM blocked_locks.relation
    AND blocking_locks.page IS NOT DISTINCT FROM blocked_locks.page
    AND blocking_locks.tuple IS NOT DISTINCT FROM blocked_locks.tuple
    AND blocking_locks.virtualxid IS NOT DISTINCT FROM blocked_locks.virtualxid
    AND blocking_locks.transactionid IS NOT DISTINCT FROM blocked_locks.transactionid
    AND blocking_locks.classid IS NOT DISTINCT FROM blocked_locks.classid
    AND blocking_locks.objid IS NOT DISTINCT FROM blocked_locks.objid
    AND blocking_locks.objsubid IS NOT DISTINCT FROM blocked_locks.objsubid
    AND blocking_locks.pid != blocked_locks.pid
JOIN pg_catalog.pg_stat_activity blocking_activity ON blocking_activity.pid = blocking_locks.pid
WHERE NOT blocked_locks.granted;"
```

**Solutions**:
```bash
# 1. Kill blocking queries (use PID from above query)
docker compose exec user-db psql -U postgres -c "SELECT pg_terminate_backend(12345);"

# 2. Restart database if needed
docker compose restart user-db

# 3. Optimize queries to reduce lock time
```

## Frontend Issues

### Vue.js Application Issues

#### Issue: API Calls Failing with CORS
**Symptoms**: CORS errors in browser console

**Diagnosis**:
```bash
# Check browser console for CORS errors
# Check API Gateway CORS configuration
docker compose logs api-gateway | grep -i cors

# Test CORS with curl
curl -H "Origin: http://localhost:3000" \
     -H "Access-Control-Request-Method: POST" \
     -H "Access-Control-Request-Headers: X-Requested-With" \
     -X OPTIONS \
     http://localhost:8080/api/users/profile
```

**Solutions**:
```bash
# 1. Check CORS configuration in API Gateway
# Ensure frontend URL is in allowed origins

# 2. Restart API Gateway
docker compose restart api-gateway

# 3. For development, ensure frontend proxy is configured correctly
# Check vite.config.js proxy settings
```

#### Issue: Frontend Not Loading
**Symptoms**: Blank page or build errors

**Diagnosis**:
```bash
# Check frontend logs
docker compose logs frontend

# Check if frontend is running
curl http://localhost:3000

# Check build process
docker compose exec frontend npm run build
```

**Solutions**:
```bash
# 1. Restart frontend
docker compose restart frontend

# 2. Rebuild frontend container
docker compose build frontend
docker compose up -d frontend

# 3. Check for JavaScript errors in browser console
```

## Network and Connectivity Issues

### Docker Network Issues

#### Issue: Services Can't Communicate
**Symptoms**: Services return connection refused errors

**Diagnosis**:
```bash
# Check Docker networks
docker network ls
docker network inspect healthcare-platform_healthcare-network

# Test connectivity between containers
docker compose exec api-gateway ping user-service
docker compose exec user-service ping user-db

# Check service discovery
docker compose exec api-gateway nslookup user-service
```

**Solutions**:
```bash
# 1. Restart Docker Compose
docker compose down
docker compose up -d

# 2. Remove and recreate network
docker compose down
docker network prune
docker compose up -d

# 3. Check port conflicts
netstat -tulpn | grep -E "(8080|8081|8082|8083)"
```

### DNS Resolution Issues

#### Issue: Service Names Not Resolving
**Symptoms**: "name resolution failed" errors

**Diagnosis**:
```bash
# Check Docker DNS
docker compose exec api-gateway cat /etc/resolv.conf
docker compose exec api-gateway nslookup user-service

# Check service definitions
docker compose config
```

**Solutions**:
```bash
# 1. Use container names from docker compose.yml
# Ensure service names match in configuration

# 2. Restart Docker daemon if needed
sudo systemctl restart docker
docker compose up -d

# 3. Use IP addresses as temporary workaround
docker inspect $(docker compose ps -q user-service) | grep IPAddress
```

## Performance Issues

### High Memory Usage

#### Issue: Container Memory Limits Exceeded
**Symptoms**: Containers being killed (OOMKilled)

**Diagnosis**:
```bash
# Check memory usage
docker stats

# Check container limits
docker inspect $(docker compose ps -q api-gateway) | grep -i memory

# Check system memory
free -h
```

**Solutions**:
```bash
# 1. Increase memory limits in docker compose.yml
mem_limit: 1g

# 2. Optimize application memory usage
# Check for memory leaks in application logs

# 3. Add swap space if needed
sudo fallocate -l 2G /swapfile
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile
```

### High CPU Usage

#### Issue: Services Using 100% CPU
**Symptoms**: Slow response times, high load average

**Diagnosis**:
```bash
# Check CPU usage
docker stats
top -p $(docker inspect --format '{{.State.Pid}}' $(docker compose ps -q api-gateway))

# Check for infinite loops in logs
docker compose logs api-gateway | grep -i error
```

**Solutions**:
```bash
# 1. Restart affected services
docker compose restart api-gateway

# 2. Scale horizontally if needed
docker compose up -d --scale api-gateway=2

# 3. Profile application for CPU hotspots
# Use Go profiling tools: go tool pprof
```

## Security Issues

### Authentication Problems

#### Issue: Brute Force Attacks
**Symptoms**: Many failed login attempts

**Diagnosis**:
```bash
# Check failed authentication attempts
docker compose logs user-service | grep -i "invalid credentials"

# Check rate limiting effectiveness
curl -v http://localhost:8080/stats
```

**Solutions**:
```bash
# 1. Implement account lockout
# Add lockout logic to user service

# 2. Reduce rate limits temporarily
export RATE_LIMIT_RPM=10
docker compose restart api-gateway

# 3. Block suspicious IPs at network level
# Use iptables or cloud provider security groups
```

### Data Access Issues

#### Issue: Unauthorized Data Access
**Symptoms**: Users accessing data they shouldn't

**Diagnosis**:
```bash
# Check access logs
docker compose logs | grep -E "(403|unauthorized|forbidden)"

# Verify user roles and permissions
docker compose exec user-db psql -U postgres -d user_service_db -c "
SELECT email, role, is_active FROM users WHERE role != 'admin';"
```

**Solutions**:
```bash
# 1. Review and fix authorization logic
# Check role-based access control implementation

# 2. Audit user permissions
# Ensure users have minimum required permissions

# 3. Enable detailed audit logging
export LOG_LEVEL=debug
docker compose restart user-service patient-service appointment-service
```

## Backup and Recovery Issues

### Backup Failures

#### Issue: Database Backup Not Working
**Symptoms**: Backup scripts failing

**Diagnosis**:
```bash
# Test backup manually
docker compose exec user-db pg_dump -U postgres user_service_db > test_backup.sql

# Check backup script logs
crontab -l
tail /var/log/backup.log
```

**Solutions**:
```bash
# 1. Fix backup script permissions
chmod +x /path/to/backup-script.sh

# 2. Ensure backup directory exists
mkdir -p /backup/healthcare

# 3. Test database connectivity in backup script
pg_isready -h localhost -p 5432 -U postgres
```

## Emergency Procedures

### Complete System Recovery

#### When Everything is Down
```bash
# 1. Stop all services
docker compose down

# 2. Check system resources
df -h
free -h
docker system df

# 3. Clean up if needed
docker system prune -a
docker volume prune

# 4. Restart from clean state
docker compose up -d

# 5. Verify each service
curl http://localhost:8080/health
curl http://localhost:8081/health
curl http://localhost:8082/health
curl http://localhost:8083/health
```

### Data Corruption Recovery

#### If Database is Corrupted
```bash
# 1. Stop all services
docker compose stop

# 2. Backup current state (even if corrupted)
docker compose exec user-db pg_dump -U postgres user_service_db > corrupted_backup.sql

# 3. Restore from last known good backup
cat last_good_backup.sql | docker compose exec -T user-db psql -U postgres -d user_service_db

# 4. Restart services
docker compose start

# 5. Verify data integrity
docker compose exec user-db psql -U postgres -d user_service_db -c "SELECT COUNT(*) FROM users;"
```

## Monitoring and Alerting

### Set Up Basic Monitoring

```bash
# Create monitoring script
cat > monitor.sh << 'EOF'
#!/bin/bash
while true; do
    echo "=== $(date) ==="
    curl -s http://localhost:8080/health | jq .
    curl -s http://localhost:8080/stats | jq .total_requests
    docker stats --no-stream --format "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}"
    echo ""
    sleep 60
done
EOF

chmod +x monitor.sh
./monitor.sh
```

### Log Analysis

```bash
# Aggregate error analysis
docker compose logs | grep -i error | sort | uniq -c | sort -nr

# Monitor response times
docker compose logs api-gateway | grep -o 'took [0-9]*ms' | sort -nr

# Check authentication failures
docker compose logs user-service | grep -i "authentication failed" | wc -l
```

---

**For immediate help, check the service logs first:**
```bash
docker compose logs <service-name> --tail=100
```

**If issues persist, contact the development team with:**
- Error messages from logs
- Steps to reproduce the issue
- Current system status (docker compose ps)
- Recent changes made to the system