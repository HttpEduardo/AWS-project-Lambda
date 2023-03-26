package event

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

type Parser struct{}

type detail struct {
	RequestParameters RequestParameters `json:"requestParameters"`
}

type RequestParameters struct {
	UserName string `json:"userName"`
}

func (detail *detail) UserName() string {
	return detail.RequestParameters.UserName
}

type User struct {
	Name string
}

func (*Parser) Parse(ctx context.Context, ev events.CloudWatchEvent) (User, error) {
	ed := detail{}
	user := User{}
	if err := json.Unmarshal(ev.Detail, &ed); err != nil {
		return user, fmt.Errorf("parse a request body detail: %w", err)
	}
	user.Name = ed.UserName()
	if user.Name == "" {
		return user, errors.New(`request body is invalid. created user name is missing`)
	}
	return user, nil
}
