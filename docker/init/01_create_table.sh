#!/bin/bash
echo "Creating DynamoDB table..."
awslocal dynamodb create-table \
    --table-name Jobs \
    --attribute-definitions AttributeName=jobId,AttributeType=S \
    --key-schema AttributeName=jobId,KeyType=HASH \
    --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 \
    --region ap-northeast-1

echo "Inserting test data..."
awslocal dynamodb put-item \
    --table-name Jobs \
    --item '{"jobId": {"S": "job_a"}, "jobTitle": {"S": "Software Engineer"}, "jobContent": {"S": "Develop awesome software."}}' \
    --region ap-northeast-1

awslocal dynamodb put-item \
    --table-name Jobs \
    --item '{"jobId": {"S": "job_b"}, "jobTitle": {"S": "Product Manager"}, "jobContent": {"S": "Manage product roadmap."}}' \
    --region ap-northeast-1

echo "Done."
