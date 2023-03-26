package controller

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/sirupsen/logrus"
)

type iamMock struct {
	cnt int
}

func (mock *iamMock) UpdateLoginProfileWithContext(ctx aws.Context, input *iam.UpdateLoginProfileInput, opts ...request.Option) (*iam.UpdateLoginProfileOutput, error) {
	if mock.cnt < 2 {
		mock.cnt++
		return nil, awserr.New(iam.ErrCodeEntityTemporarilyUnmodifiableException, "", errors.New(""))
	}
	return &iam.UpdateLoginProfileOutput{}, nil
}

func Test_updateLoginProfile(t *testing.T) {
	t.Parallel()
	data := []struct {
		title string
		iam   IAM
	}{
		{
			iam: &iamMock{},
		},
	}
	ctx := context.Background()
	logE := logrus.WithFields(logrus.Fields{})
	for _, d := range data {
		d := d
		t.Run(d.title, func(t *testing.T) {
			t.Parallel()
			if err := updateLoginProfile(ctx, logE, d.iam, &iam.UpdateLoginProfileInput{}); err != nil {
				t.Fatal(err)
			}
		})
	}
}
