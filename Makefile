TOKEN?=not_set
CHAT_ID?=-1002083536699
SECRET_TOKEN?=secret

check_values:
	@if [ $(TOKEN) == "not_set" ] ; then \
		echo "TOKEN not set!"; \
		exit 1;\
	fi

docker_build: 
	docker build -t s1kai/memes_bot:develop .

docker_run: check_values
	docker run --name bot_develop \
		-e TOKEN=$(TOKEN) \
		-e PORT=5050 \
		-e CHAT_ID=$(CHAT_ID) \
		-e SECRET_TOKEN=$(SECRET_TOKEN) \
		-e USER_WHITE_LIST=user_id \
		-p 5050:5050 s1kai/memes_bot:develop

req:
	curl localhost:5050/webhook \
		-H "X-Telegram-Bot-Api-Secret-Token:secret" \
		-d '{
  "update_id": 569506765,
  "message": {
    "message_id": 160,
    "from": {
      "username": "user_id",
    },
    "chat": {
      "username": "user_id",
    },
    "text": "https://www.instagram.com/reel/C0KJ4qtAStm/?igsh=aWY3cDMxaHRqdGZ2",
    "entities": [
      {
        "offset": 0,
        "length": 65,
        "type": "url"
      }
    ],
    "link_preview_options": {
      "is_disabled": true
    }
  }
}'
	

