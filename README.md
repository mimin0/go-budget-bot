[![Build Status][travis-badge]][travis]

# Go-Budget-Bot
this is the tool that can help you to simplify managing and/or handling your expences within Telegram bot

## Keys feature
- works with google spread sheet: you can generate your own report in any time
- cross platform: you are able to add work with your expeces via telegram client on your phone/PC/etc.


## Prepare the workspace
 - Set the GOPATH environment variable to your working directory.
 - Get the Google Sheets API Go client library and OAuth2 package using the following commands:

    ```go get -u google.golang.org/api/sheets/v4```</br>
    ```go get -u golang.org/x/oauth2/...```

 - Get the Telegram bot lib:

    ```go get gopkg.in/telegram-bot-api.v4```

## Work notes
- GO + Google Sheets https://developers.google.com/sheets/api/quickstart/go https://developers.google.com/sheets/api/guides/authorizing
