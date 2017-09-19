package main

import (
	"fmt"
	"os"

	"github.com/grokify/gotilla/fmt/fmtutil"
	"github.com/grokify/oauth2-util-go/services/ringcentral"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

func main() {
	err := godotenv.Load(os.Getenv("ENV_PATH"))

	if err != nil {
		panic(err)
	}

	cfg := oauth2.Config{
		ClientID:     os.Getenv("RC_CLIENT_ID"),
		ClientSecret: os.Getenv("RC_CLIENT_SECRET"),
		Endpoint:     ringcentral.NewEndpoint(os.Getenv("RC_SERVER_HOSTNAME"))}

	token, err := cfg.PasswordCredentialsToken(
		oauth2.NoContext,
		os.Getenv("RC_USER_USERNAME"),
		os.Getenv("RC_USER_PASSWORD"))

	if err != nil {
		panic(err)
	}

	fmtutil.PrintJSON(token)

	fmt.Println("DONE")
}
