#!/bin/bash
docker build . -t crypto-bot-container
docker rm -f crypto-bot-container
docker run -d --restart=always --log-opt max-size=50m —name=cryptoBot crypto-bot-container