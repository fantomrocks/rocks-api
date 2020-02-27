#!/bin/bash
#
# Try to send test query to the API
#
curl \
  -X POST \
  -H "Content-Type: application/json" \
  --data '{"query": "{accounts{id,name,balance}}"}' \
  https://fantom.rocks/api
