#!/bin/bash
set -e

DB_FILE="/db/shared-local-instance.db"
AWS_ENDPOINT=http://localhost:8000
DYNAMO_USER="dynamodblocal"
AWS_REGION=us-west-2
TABLE_NAME=client-entity-dev

echo "Fixing ownership of /db to $DYNAMO_USER user"
chown -R $DYNAMO_USER:$DYNAMO_USER /db

if [ ! -f "$DB_FILE" ]; then
  echo "No DB file found. Bootstrapping schema..."
  echo "Starting DynamoDBLocal as $DYNAMO_USER..."

  su $DYNAMO_USER -c "java $*" &
  DYNAMO_PID=$!

  echo "Waiting for DynamoDBLocal to become available..."
  until su $DYNAMO_USER -c "curl -s '$AWS_ENDPOINT' > /dev/null"; do
    sleep 0.5
  done

  echo "Creating table from /schema/table.json..."
  su $DYNAMO_USER -c "aws dynamodb create-table --cli-input-json file:///schema/table.json \
    --endpoint-url $AWS_ENDPOINT --region $AWS_REGION"

  echo "Enabling TTL..."
  su $DYNAMO_USER -c "aws dynamodb update-time-to-live --table-name $TABLE_NAME \
    --time-to-live-specification 'Enabled=true, AttributeName=ExpirationTime' \
    --endpoint-url $AWS_ENDPOINT --region $AWS_REGION"

  kill $DYNAMO_PID
  wait $DYNAMO_PID 2>/dev/null || true

  echo "Schema initialization complete. "
else
  echo "DB file already exists. Skipping initialization."
fi

echo "Starting DynamoDBLocal as $DYNAMO_USER..."
exec su $DYNAMO_USER -c "java $*"
