package main

import (
	"arisetyawan/adrena-auto/hatcher"

	"github.com/aws/aws-lambda-go/lambda"
)

type Response struct {
	Message string `json:"message"`
}

func Handler() (Response, error) {
	hatcherBridge := hatcher.NewBridge()
	hatcherBridge.DoCheck("CHECKOUT")

	return Response{
		Message: "Checking Runs",
	}, nil
}

func main() {
	lambda.Start(Handler)
}
