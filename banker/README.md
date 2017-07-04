# Banker

Banker is stupidly go-based tool to save and track personal bank account data. Banker will provide brief information about your financial information by using Telegram Bot API to send it right away to your favorite messager app.

# Todo

- ~~Create CSV parser~~
- ~~Create REST API to accept csv data~~
- Init Telegram BOT to track my data directly to phone

# Caveats

- Only supports BCA Klikpay mutation csv format
- Should export csv file manually

# Disclaimer

These projects under daily_tool are experimental. It's only intended to make my personal life a little bit easier and also to improve my understanding about certain programming language

# Endpoints
```
  POST  /banker/upload
```
param :
  - rahasianegara (required) : file to be uploaded, must be formatted as CSV

```
  GET   /banker/report/daily
```
param :
  - month (required) : month of transaction. formatted in MM
  - year (required) : year of transaction. formatted in YY

```
  GET /banker/report/monthly
```
param :
  - year (required) : year of transaction. formatted in YY
