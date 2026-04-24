#!/usr/bin/env bash
set -euo pipefail

wait_for() {
  local url="$1"
  local name="$2"

  for _ in $(seq 1 60); do
    if curl -fsS "$url" >/dev/null 2>&1; then
      echo "$name is ready"
      return 0
    fi
    sleep 1
  done

  echo "timeout waiting for $name"
  exit 1
}

json_field() {
  local url="$1"
  local field="$2"

  curl -fsS "$url" | python -c "import sys, json; print(json.load(sys.stdin)['$field'])"
}

echo "starting containers"
docker compose up -d --build

trap 'docker compose down -v' EXIT

wait_for "http://localhost:8080/health" "node1"
wait_for "http://localhost:8081/health" "node2"

node1_len_before="$(json_field "http://localhost:8080/chain" "length")"
node2_len_before="$(json_field "http://localhost:8081/chain" "length")"

echo "node1 length: $node1_len_before"
echo "node2 length: $node2_len_before"

echo "creating block on node1"
curl -fsS -X POST "http://localhost:8080/blocks" \
  -H "Content-Type: application/json" \
  -d '{
    "transactions": [
      {
        "id": "tx-1",
        "from": "alice",
        "to": "bob",
        "amount": 10
      }
    ]
  }' >/dev/null

node1_len_after="$(json_field "http://localhost:8080/chain" "length")"
echo "node1 length after block: $node1_len_after"

if [ "$node1_len_after" -le "$node1_len_before" ]; then
  echo "node1 did not grow after block creation"
  exit 1
fi

echo "waiting for propagation to node2"
for _ in $(seq 1 60); do
  node2_len_after="$(json_field "http://localhost:8081/chain" "length")"
  if [ "$node2_len_after" -eq "$node1_len_after" ]; then
    break
  fi
  sleep 1
done

node2_len_after="$(json_field "http://localhost:8081/chain" "length")"
echo "node2 length after propagation: $node2_len_after"

if [ "$node2_len_after" -ne "$node1_len_after" ]; then
  echo "node2 did not receive the updated chain"
  exit 1
fi

echo "e2e test passed"