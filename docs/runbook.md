# Mitran Operational Runbook

## 1. Engine Won't Start

### Symptoms
- `mitran start` hangs or exits immediately
- Dashboard shows "Engine offline"
- Port 7780 not listening

### Diagnosis
```bash
lsof -i :7780                          # Check port conflict
echo $MITRAN_SESSION_SECRET            # Verify env vars set
which mitran-engine || ls ./build/     # Check binary exists
journalctl -u mitran-engine --since "5m ago"  # Check logs
```

### Resolution
1. Kill conflicting process: `kill $(lsof -t -i :7780)`
2. Export required env vars: `export MITRAN_SESSION_SECRET=<secret>`
3. Rebuild if binary missing: `go build -o build/mitran-engine ./cmd/engine`
4. Start with verbose logging: `MITRAN_LOG_LEVEL=debug mitran-engine`

---

## 2. Worker Unreachable

### Symptoms
- Tasks stuck in "pending" state
- Engine logs: `connection refused` to worker
- `/api/v1/health` returns worker status "down"

### Diagnosis
```bash
curl http://localhost:8888/health       # Direct worker health check
echo $WORKER_URL                        # Verify URL config
ps aux | grep gunicorn                  # Check worker process
pip check                               # Verify Python dependencies
```

### Resolution
1. Set correct URL: `export WORKER_URL=http://localhost:8888`
2. Install deps: `cd agent-worker && pip install -r requirements.txt`
3. Restart worker: `gunicorn -w 4 -b 0.0.0.0:8888 worker.app:app`
4. Check firewall: `sudo ufw allow 8888/tcp`

---

## 3. ChromaDB Connection Failed

### Symptoms
- Memory/search operations return 500 errors
- Logs: `ConnectionError: ChromaDB unreachable`
- Semantic search returns empty results

### Diagnosis
```bash
echo $CHROMADB_URL                      # Default: http://localhost:8000
curl http://localhost:8000/api/v1/heartbeat  # ChromaDB health
docker ps | grep chroma                 # Check container running
docker logs mitran-chromadb --tail 20   # Container logs
```

### Resolution
1. Start ChromaDB: `docker compose up -d chromadb`
2. Set URL: `export CHROMADB_URL=http://localhost:8000`
3. Reset if corrupted: `docker compose down chromadb && docker volume rm mitran_chroma_data && docker compose up -d chromadb`
4. Verify: `curl http://localhost:8000/api/v1/collections`

---

## 4. High Memory Usage

### Symptoms
- OOM kills in container environments
- Engine response times > 5s
- `/metrics` shows memory > 80% threshold

### Diagnosis
```bash
curl http://localhost:7780/metrics | grep memory  # Prometheus metrics
curl http://localhost:7780/api/v1/queue/stats     # Queue depth
ps aux --sort=-%mem | head -5                     # Top memory processes
docker stats mitran-engine mitran-worker          # Container memory
```

### Resolution
1. Check queue backlog: if queue > 100 tasks, scale workers
2. Restart workers to release memory: `docker compose restart worker`
3. Reduce worker concurrency: set `MITRAN_WORKER_CONCURRENCY=2`
4. Add memory limit: `docker compose up -d --memory=2g`
5. Enable GC tuning: `export GOGC=50` (engine), reduce batch sizes

---

## 5. Authentication Failures

### Symptoms
- 401/403 on API calls
- Dashboard login loop
- Logs: `invalid session` or `token expired`

### Diagnosis
```bash
echo $MITRAN_SESSION_SECRET            # Must be set, min 32 chars
echo $MITRAN_JWT_SECRET                # Required for API auth
curl -v http://localhost:7780/api/v1/auth/health  # Auth endpoint
# Check cookie domain matches request origin
```

### Resolution
1. Set secrets: `export MITRAN_SESSION_SECRET=$(openssl rand -hex 32)`
2. Set JWT secret: `export MITRAN_JWT_SECRET=$(openssl rand -hex 32)`
3. Clear browser cookies and retry login
4. For SSO: verify `MITRAN_SSO_ISSUER` and `MITRAN_SSO_CLIENT_ID` are correct
5. Check session TTL (default 24h): `export MITRAN_SESSION_TTL=24h`

---

## 6. Database Corruption

### Symptoms
- JSON parse errors in logs
- Tasks/projects return partial or empty data
- Engine crashes on startup with "malformed data"

### Diagnosis
```bash
ls -la .mitran/data/                    # Check file sizes (0 = corrupt)
python3 -c "import json; json.load(open('.mitran/data/tasks.json'))"  # Validate JSON
find .mitran/data -name "*.json" -empty  # Find empty files
cat .mitran/data/*.json | python3 -m json.tool > /dev/null  # Batch validate
```

### Resolution
1. **Backup immediately**: `cp -r .mitran/data/ .mitran/data-backup-$(date +%s)/`
2. Restore from last good backup: `cp .mitran/data-backup-<timestamp>/* .mitran/data/`
3. Run migration: `mitran-engine migrate --data-dir .mitran/data/`
4. If no backup exists, reconstruct from WAL: `mitran-engine recover --wal .mitran/wal/`
5. Prevent future corruption: ensure `MITRAN_ATOMIC_WRITES=true` (default)
