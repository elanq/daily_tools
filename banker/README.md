# Banker

Banker is stupidly go-based tool to save and track personal bank account data. Banker will provide brief information about your financial information by using REST API to send it right away to your favorite app.

# Supported features

- Receive bank mutation data from klikbca in CSV format
- Parse csv data and save them to mongodb
- Provide endpoints to serve transactions data and return into different type (raw and summary)
- Real time google sheet backup

# Todo

- ~~Create CSV parser~~
- ~~Create REST API to accept csv data~~
- Provide support to return transaction data as chart
- Multiple bank support

# Caveats

- Only supports BCA Klikpay mutation csv format
- Should export csv file manually
- This intended to only support one bank account at a time. So multiple user support is not considered as the part of development

# Disclaimer

These projects under daily_tool are experimental. It's only intended to make my personal life a little bit easier and also to improve my understanding about certain programming language

# Endpoints
```
  POST  /banker/upload
```
param :
  - rahasianegara (required) : file to be uploaded, must be formatted as CSV

  - year (required) : year of transaction. formatted in YY

```
  GET   /banker/report/monthly
```
param :
  - month (required) : month of transaction. formatted in MM
  - year (required) : year of transaction. formatted in YY
  - type (optional) : specify return type of data, default is raw json per transaction
    - summary : will return summary of transaction

```
  GET /banker/report/yearly
```
param :
  - year (required) : year of transaction. formatted in YY
  - type (optional) : specify return type of data, default is raw json per transaction
    - summary : will return summary of transaction

