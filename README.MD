# Bot GPT 

BotGPT is a multifunctional bot that integrates various communication software such as Line and Telegram. It features translation, conversation, creating Azure Work Items, AI chat, yes/no question asking, chat without tagging bot, image generation, and more.

## Features

- `/help`: Shows current mode and commands.
- `/entw@mazeaibot`: Translates to English (en-US).
- `/chat {message}`: Engages in conversation.
- `/cw@mazeaibot`: Creates a new Azure Work Item with the title `{title}` and content `{content}`.
- `/ct@mazeaibot {message}`: Engages in AI chat.
- `/asker@mazeaibot`: Asks a yes/no question.
- `/chatnotag`: Chats without tagging bot.
- `/jptw@mazeaibot`: Translates to Japanese (ja-JP).
- `/image`: Generates an image using DALL-E.



## Installation

1. Clone this repository:
```
git clone https://github.com/maze0417/botgpt.git
```
2. Config telegram , line ,open ai token
```
cd internal/config/botgpt

#change local.yml

linebot:
  access_token: your-line-bot-token
  secret: your-access-secret
openai:
  access_token: your-access-token
tg:
  access_token: your-access-token

```


## Run the bot

1. use go command
```
# go version is 1.20.4
go mod vendor  

go run main.go
```

2. use docker compose
```
docker compose up -d
```