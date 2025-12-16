package main

import (
	"context"
	"fmt"
	"net/url"

	"github.com/alanshaw/1up-service/pkg/capabilities/debug"
	"github.com/alanshaw/ucantone/client"
	"github.com/alanshaw/ucantone/did"
	"github.com/alanshaw/ucantone/execution"
	"github.com/alanshaw/ucantone/ipld"
	"github.com/alanshaw/ucantone/ipld/datamodel"
	"github.com/alanshaw/ucantone/principal/ed25519"
	"github.com/alanshaw/ucantone/result"
	"github.com/alanshaw/ucantone/ucan/invocation"
)

const (
	serviceID  = "did:key:z6MkiZfWmWbXpBj2bxF4w8ifBRi8PRSa83qUFTWq7rb73Hse"
	serviceURL = "http://localhost:3000"
)

func main() {
	alice, err := ed25519.Generate()
	if err != nil {
		panic(err)
	}

	service, err := did.Parse(serviceID)
	if err != nil {
		panic(err)
	}

	inv, err := debug.Echo.Invoke(
		alice,
		alice,
		&debug.EchoArguments{
			Message: "Hello, UCAN!",
		},
		invocation.WithAudience(service),
	)
	if err != nil {
		panic(err)
	}

	url, err := url.Parse(serviceURL)
	if err != nil {
		panic(err)
	}

	client, err := client.NewHTTP(url)
	if err != nil {
		panic(err)
	}

	res, err := client.Execute(execution.NewRequest(context.Background(), inv))
	if err != nil {
		panic(err)
	}

	result.MatchResultR0(
		res.Result(),
		func(o ipld.Any) {
			args := debug.EchoOK{}
			err := datamodel.Rebind(datamodel.NewAny(o), &args)
			if err != nil {
				panic(err)
			}
			fmt.Printf("Echo response: %+v\n", args)
		},
		func(x ipld.Any) {
			fmt.Printf("Invocation failed: %v\n", x)
		},
	)
}
