## Overview
This is the library for Telegram API. With this library you can set responses on user's commands.

## How to run

`go get` it using command: 

    go get github.com/sosshik/tgapi@latest

Import to your project:

    import "github.com/sosshik/tgapi"

## How to use

1. Load your Telegram Token from .env file
2. Create Api Struct using your Telegram Token and &http.Client{}
2. Set instance that will handle user input
4. Add callback that is bounded to the command
5. Create a serever with handler set to our webhook
6. Set the server to listen and serve

## Example 

    package example

    import (
        "os"

        "github.com/sosshik/tgapi"
        "github.com/Valgard/godotenv"
    )

    func main (){
        // Load your Telegram Token from .env file
        godotenv.Load()

        // Create Api Struct using your Telegram Token and &http.Client{}
        api := tgapi.Api{
            SendMessageURL: "https://api.telegram.org/bot" + os.Getenv("TELEGRAM_TOKEN") + "/sendMessage",
			GetUpdatesURL:  "https://api.telegram.org/bot" + os.Getenv("TELEGRAM_TOKEN") + "/getUpdates",
            HTTPClient:  &http.Client{},
        }
        // Set instance that will handle user input 
        tgapi.UserInput = &YourStruct{}

        // Add callback that is bounded to the command
        api.AddCallback(YourWebhook, "/yourcommand")

        // Create a serever with handler set to our webhook
        server := &http.Server{
            Addr:              "Your address",
            ReadHeaderTimeout: 5 * time.Second,
            Handler:           http.HandlerFunc(api.HandleTelegramWebHook),
        }

        // Set the server to listen and serve 
        err := server.ListenAndServe()
        if err != nil {
            log.Fatal(err)
        }
    }

    // This function always should take *tgapi.Api and tgapi.Update as input, it cannot output anything
    func YourWebhook (a *tgapi.Api, update tgapi.Update){
        // Your function
    }

    type YourStruct struct{
        //Your fileds 
    }

    func (y *YourStruct) HandleUserInput(a *tgapi.Api, update Update){
        // Your function
    }



