package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sirupsen/logrus"
	lmb "github.com/suzuki-shunsuke/aws-iam-cred-sender/pkg/lambda"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	handler := lmb.Handler{}
	ctx := context.Background()
	if err := handler.Init(ctx); err != nil {
		logrus.Fatal(err)
	}
	lambda.Start(handler.Start)
}
