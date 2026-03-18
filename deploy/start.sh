#!/bin/sh
set -e

# Set up daily scraper cron job (default: 2 AM UTC).
SCRAPE_CRON="${SCRAPE_CRON:-0 2 * * *}"
mkdir -p /etc/crontabs

# Export env vars for cron (crond runs with a clean environment).
printenv | grep -E '^(DATABASE_URL|OPENAI_API_KEY|GOOGLE_MAPS_API_KEY)=' > /tmp/scrape.env
echo "${SCRAPE_CRON} . /tmp/scrape.env && /app/scrape -config /app/roasters.yaml >> /tmp/scrape.log 2>&1" > /etc/crontabs/root
crond -l 8

# Start Go backend on internal port 3001
PORT=3001 ./server &
GO_PID=$!

# Wait briefly for the Go backend to start
sleep 1

# Start Next.js frontend on Dokku's assigned PORT (default 5000)
export API_URL="http://localhost:3001"
export HOSTNAME="0.0.0.0"

node server.js &
NODE_PID=$!

# Handle shutdown: forward signals to all processes
trap 'kill $GO_PID $NODE_PID 2>/dev/null; killall crond 2>/dev/null; wait' SIGTERM SIGINT

# Wait for either process to exit
wait -n $GO_PID $NODE_PID 2>/dev/null || true

# If one exits, kill the other
kill $GO_PID $NODE_PID 2>/dev/null
killall crond 2>/dev/null
wait
