# Docker Compose V2 Migration

This project uses the modern `docker compose` command (Compose V2) instead of the legacy `docker-compose` standalone tool.

## What Changed

### Old Command (Legacy)
```bash
docker-compose up -d
docker-compose logs -f
docker-compose down
```

### New Command (Integrated)
```bash
docker compose up -d
docker compose logs -f
docker compose down
```

## Benefits of Docker Compose V2

- **Integrated**: Built into Docker CLI, no separate installation needed
- **Performance**: Faster startup and better resource utilization
- **Consistency**: Single Docker CLI for all operations
- **Future-proof**: Official direction for Docker Compose

## Migration Guide

### 1. Check Your Docker Version
```bash
# Check if you have Docker Compose V2
docker compose version

# Should show something like:
# Docker Compose version v2.x.x
```

### 2. Update Docker (if needed)
```bash
# For Docker Desktop users - update through the application

# For Linux users with package manager:
sudo apt update && sudo apt upgrade docker-ce docker-ce-cli

# Or install Docker Desktop for Linux
```

### 3. Verify Installation
```bash
# Test the new command
docker compose --help

# Start the healthcare platform
docker compose up -d
```

## Compatibility

### File Compatibility
- **compose.yml** files are identical between V1 and V2
- **docker-compose.yml** files work with both versions
- No changes needed to existing configuration files

### Command Compatibility
Most commands are identical, just replace `docker-compose` with `docker compose`:

| Legacy Command | Modern Command |
|----------------|----------------|
| `docker-compose up` | `docker compose up` |
| `docker-compose down` | `docker compose down` |
| `docker-compose logs` | `docker compose logs` |
| `docker-compose exec` | `docker compose exec` |
| `docker-compose build` | `docker compose build` |
| `docker-compose ps` | `docker compose ps` |

## Troubleshooting

### "docker compose: command not found"
Your Docker installation doesn't include Compose V2. Options:

1. **Update Docker Desktop** (recommended)
2. **Install Docker Engine with Compose plugin**:
   ```bash
   # Ubuntu/Debian
   sudo apt install docker-compose-plugin
   ```
3. **Use legacy docker-compose** as fallback (not recommended)

### "docker-compose" still works
If you have both versions installed:
- Use `docker compose` (recommended)
- Legacy `docker-compose` may still work but is deprecated

### Performance Issues
- Compose V2 should be faster than V1
- If experiencing issues, try clearing Docker cache:
  ```bash
  docker system prune -a
  docker compose up -d --build
  ```

## Project-Specific Updates

All documentation in this project has been updated to use `docker compose`:

- ✅ Main README.md
- ✅ Development setup guide
- ✅ Service documentation
- ✅ Troubleshooting guide
- ✅ CLAUDE.md (AI context)

## Benefits for Healthcare Platform

Using Docker Compose V2 provides:

1. **Faster Startup**: Quicker service initialization
2. **Better Resource Management**: More efficient container handling
3. **Improved Logging**: Better log aggregation and formatting
4. **Enhanced Networking**: More reliable service communication
5. **Future Compatibility**: Aligned with Docker's roadmap

## Migration Checklist

- [ ] Update Docker to version 20.10+
- [ ] Verify `docker compose` command works
- [ ] Update any custom scripts to use `docker compose`
- [ ] Update CI/CD pipelines if using docker-compose
- [ ] Train team members on new command syntax

---

**Note**: This healthcare platform is fully compatible with Docker Compose V2 and all documentation reflects the modern command syntax.