#!/bin/bash
docker build . -t cryptoBot-container
docker rm -f cryptoBot-container
docker run -d --restart=always --log-opt max-size=50m —name=cryptoBot cryptoBot-container