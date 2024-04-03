package main

import (
	"embed"
	"fmt"
	"log"
	"os"

	"github.com/darkqiank/who-dat/api"
	"github.com/valyala/fasthttp"
)

var staticAssets embed.FS

func main() {

	// Custom request handler for fasthttp
	requestHandler := func(ctx *fasthttp.RequestCtx) {
		// Handle API requests
		api.MainHandler(ctx)
	}

	// Choose the port to start server on
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	serverAddress := fmt.Sprintf(":%s", port)

	asciiArt := `
__          ___             _____        _  ___  
\ \        / / |           |  __ \      | ||__ \ 
 \ \  /\  / /| |__   ___   | |  | | __ _| |_  ) |
  \ \/  \/ / | '_ \ / _ \  | |  | |/ _` + "`" + ` | __|/ / 
   \  /\  /  | | | | (_) | | |__| | (_| | |_|_|  
    \/  \/   |_| |_|\___/  |_____/ \__,_|\__(_)																							
`
	log.Println(asciiArt)
	log.Printf("\nWelcome to Who-Dat - WHOIS Lookup Service.\nApp up and running at %s", serverAddress)

	// Start fasthttp server
	if err := fasthttp.ListenAndServe(serverAddress, requestHandler); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}

}
