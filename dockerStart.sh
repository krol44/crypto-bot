#!/bin/bash
docker build . -t crypto-bot-container
docker rm -f crypto-bot
docker run -d --restart=always --log-opt max-size=50m --name=crypto-bot crypto-bot-container
docker prune -f