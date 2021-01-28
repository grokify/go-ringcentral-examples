package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/grokify/oauth2more/ringcentral"
	"github.com/grokify/simplego/fmt/fmtutil"
	"github.com/grokify/simplego/net/http/httpsimple"
	"github.com/jessevdk/go-flags"
)

type Options struct {
	Credentials string `short:"c" long:"credentials" description:"Credentials File"`
	Account     string `short:"a" long:"account" description:"Credentials Account"`
	To          string `short:"t" long:"to" description:"To"`
}

func getClient(filename, credsKey string) (ringcentral.Credentials, *http.Client, error) {
	cset, err := ringcentral.ReadFileCredentialsSet(filename)
	if err != nil {
		return ringcentral.Credentials{}, nil, err
	}
	creds, err := cset.Get(credsKey)
	if err != nil {
		return creds, nil, err
	}
	client, err := creds.NewClient()
	return creds, client, err
}

func main() {
	opts := Options{}
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}
	fmtutil.PrintJSON(opts)

	creds, client, err := getClient(opts.Credentials, opts.Account)
	if err != nil {
		log.Fatal(err)
	}

	sClient := httpsimple.SimpleClient{
		BaseURL:    creds.Application.ServerURL,
		HTTPClient: client}

	reqBody := HVSmsBatch{
		From: creds.PasswordCredentials.Username,
		Text: "High Volume Test from Go",
		Messages: []HVSmsMessage{
			{To: []string{opts.To}}}}

	simpleReq := httpsimple.SimpleRequest{
		Method: http.MethodPost,
		URL: creds.Application.InflateURL(
			"/restapi/v1.0/account/~/a2p-sms/batch"),
		Body:   reqBody,
		IsJSON: true}

	resp, err := sClient.Do(simpleReq)
	if err != nil {
		log.Fatal(err)
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(bytes))

	fmt.Println("DONE")
}

type HVSmsBatch struct {
	From     string         `json:"from,omitempty"`
	Text     string         `json:"text,omitempty"`
	Messages []HVSmsMessage `json:"messages,omitempty"`
}

type HVSmsMessage struct {
	To   []string `json:"to"`
	Text string   `json:"text"`
}
