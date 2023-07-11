#!/bin/bash

set -e -u -o pipefail

if [[ $# != 2 ]] || [[ $1 =  "" ]] || [[ $1 =  "-h" ]] || [[ $1 =  "--help" ]] || [[ $2 =  "" ]]; then
  echo "Script to get ical url"
  echo "Usage:"
  echo "  $0 stack_suffix aws_region"
  exit 1
fi

STACK_SUFFIX="$1"
AWS_REGION="$2"

aws cloudformation describe-stacks \
    --region "$AWS_REGION" \
    --stack-name notion-ical-syncer-stack-"$STACK_SUFFIX" \
    --query "Stacks[0].Outputs[0].OutputValue" | tr -d  '"'
