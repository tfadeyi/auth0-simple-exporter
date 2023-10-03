#!/bin/bash
ACCOUNT_ID=$1
if [[ "${ACCOUNT_ID}" == "" ]]; then
    echo "You must specify the account id as first argument."
    exit 1
fi

aws ecr get-login-password --region us-east-2 | docker login --username AWS --password-stdin ${ACCOUNT_ID}.dkr.ecr.us-east-2.amazonaws.com
docker build -t houserx-auth0-simple-exporter .
docker tag houserx-auth0-simple-exporter:latest ${ACCOUNT_ID}.dkr.ecr.us-east-2.amazonaws.com/houserx-auth0-simple-exporter:latest
docker push ${ACCOUNT_ID}.dkr.ecr.us-east-2.amazonaws.com/houserx-auth0-simple-exporter:latest