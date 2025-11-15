#!/usr/bin/env bash
set -euo pipefail

ELASTIC_URL="${ELASTIC_URL:-http://localhost:9200}"

curl -s -X POST "${ELASTIC_URL}/_aliases" -H 'Content-Type: application/json' -d '{
  "actions": [
    { "add": { "index": "articles_fa_v1", "alias": "articles_fa" } },
    { "add": { "index": "articles_en_v1", "alias": "articles_en" } },
    { "add": { "index": "forum_fa_v1",    "alias": "forum_fa"    } },
    { "add": { "index": "forum_en_v1",    "alias": "forum_en"    } }
  ]
}'
echo
echo "✓ Aliases applied."