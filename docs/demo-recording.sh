#!/bin/bash
# Record with: asciinema rec demo.cast --command bash docs/demo-recording.sh
set -e
echo "🚀 Mitran Demo"
mitran init --name "Demo Corp"
cat mitran.yaml
mitran demo &
sleep 2
curl -s localhost:7780/api/v1/health | python3 -m json.tool
curl -s localhost:7780/api/v1/agents | python3 -m json.tool
kill %1 2>/dev/null
echo "✅ Demo complete!"
