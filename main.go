package main

import (
	"os"

	"github.com/begmaroman/acme-dns-route53/cmd"
	"github.com/begmaroman/acme-dns-route53/lambda"
)

// IsLambda contains true value if this tool uses by AWS Lambda
var IsLambda bool

func init() {
	if os.Getenv("AWS_LAMBDA") == "1" {
		IsLambda = true
	}
}

func main() {
	if IsLambda {
		// Initialize Lambda handler
		lambda.Init()
	} else {
		// Initialize CLI
		cmd.Execute()
	}
}
