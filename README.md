## Demo
https://t.me/Crypto4cBot

## Setup
1. create and change docker-start.sh in crypto-bot folder
```
#!/bin/bash
docker build . -t crypto-bot-cont
docker rm -f crypto-bot
docker run -d --restart=always \
-e TZ="Europe/Moscow" \
-e is_dev="false" \
-e bot_token="00000:AAAAAAA" \
--mount type=bind,source=/var/crypto-bot,target=/db \
--log-opt max-size=5m \
--name=crypto-bot crypto-bot-cont

docker image prune -f
```

2. change bot_token and folder for db - /var/crypto-bot
3. chmod +x docker-start.sh
4. ./docker-start.sh
5. add commands in your bot in the @BotFather
```
notify_binance - Example: /notify_binance dogeusdt up 0.33
notify_binance_my_rates_list - show your rates
notify_binance_my_rates_reset - reset all your rates
```