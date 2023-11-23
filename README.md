## Overview
This is the library for Telegram API. With this library you can set response on user's commands.

## How to run

`go get` it using command: 

    go get git.foxminded.ua/foxstudent106264/tgapi@latest

Import to your project:

    import "git.foxminded.ua/foxstudent106264/tgapi"

## How to use

1. Load your Telegram Token from .env file
2. Create Api Struct using your Telegram Token and &http.Client{}
3. Add callback that is bounded to the command
4. Create a serever with handler set to our webhook
5. Set the server to listen and serve

## Example 

    package example

    import (
        "os"

        "git.foxminded.ua/foxstudent106264/tgapi"
        "github.com/Valgard/godotenv"
    )

    func main (){
        // Load your Telegram Token from .env file
        godotenv.Load()

        // Create Api Struct using your Telegram Token and &http.Client{}
        api := tgapi.Api{
            TelegramAPI: "https://api.telegram.org/bot" + os.Getenv("TELEGRAM_TOKEN") + "/sendMessage",
            HTTPClient:  &http.Client{},
        }

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



