#!/bin/bash

set -e -u -o pipefail

if [[ $# != 2 ]] || [[ $1 =  "" ]] || [[ $1 =  "-h" ]] || [[ $1 =  "--help" ]] || [[ $2 =  "" ]]; then
  echo "Script to run things locally"
  echo "Usage:"
  echo "  $0 stack_suffix aws_region"
  exit 1
fi

STACK_SUFFIX="$1"
AWS_REGION="$2"

GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o output/syncer cmd/syncer.go && cd output && zip syncer.zip syncer && cd ..

PROJECT_DIRECTORY="$(dirname "$(cd "$(dirname "$0")" && pwd)")"

aws cloudformation deploy \
    --region "$AWS_REGION" \
    --template-file "$PROJECT_DIRECTORY"/scripts/aws/infra.yml \
    --stack-name notion-ical-syncer-infra

UPLOAD_BUCKET_NAME="$(aws cloudformation describe-stacks \
    --region "$AWS_REGION" \
    --stack-name notion-ical-syncer-infra \
    --query "Stacks[0].Outputs[0].OutputValue" | tr -d  '"')"

aws cloudformation package \
    --region "$AWS_REGION" \
    --template-file "$PROJECT_DIRECTORY"/scripts/aws/template.yml \
    --s3-bucket "$UPLOAD_BUCKET_NAME" \
    --output-template-file "$PROJECT_DIRECTORY"/output/output-template.yml

aws cloudformation deploy \
    --region "$AWS_REGION" \
    --template-file "$PROJECT_DIRECTORY"/output/output-template.yml \
    --stack-name notion-ical-syncer-stack-"$STACK_SUFFIX" \
    --capabilities CAPABILITY_NAMED_IAM