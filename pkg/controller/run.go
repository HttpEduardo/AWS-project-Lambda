package controller

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/sethvargo/go-password/password"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

type Secret struct {
	SlackBotToken string `json:"slack_bot_token"`
}

func (ctrl *Controller) Run(ctx context.Context, param Param) error { //nolint:funlen,cyclop
	sess := session.Must(session.NewSession())
	logE := logrus.WithFields(logrus.Fields{
		"user_name": param.UserName,
	})

	// check user is registered at DynamoDB
	dynamoDBSvc := dynamodb.New(sess)

	if f, err := ctrl.existUserNameAtDynamoDB(ctx, dynamoDBSvc, param.UserName); err != nil {
		logE.WithError(err).Error("check whether an IAM User is registered at DynamoDB")
	} else if f {
		logE.Info("this IAM User has already been registered at DynamoDB")
		return nil
	}

	// add UserName to DynamoDB
	if err := ctrl.addUserNameToDynamoDB(ctx, dynamoDBSvc, param.UserName, strconv.FormatInt(time.Now().Add(time.Duration(ctrl.Config.DynamoDBTTL)*time.Second).Unix(), 10)); err != nil {
		logE.WithError(err).Error("add IAM User Name to DynamoDB")
	}
	logE.Info("add IAM User Name to DynamoDB")

	// get a slack user id
	logE.Info("start getting a Slack User")
	user, err := ctrl.getSlackUser(ctx, param.UserName, logE)
	if err != nil {
		return ctrl.handleSystemUser(ctx, param)
	}
	logE = logE.WithFields(logrus.Fields{
		"slack_id": user.ID,
	})
	logE.Info("get a Slack User")
	// create an initial password
	passwd, err := password.Generate(ctrl.Config.InitialPasswordLength, 10, 10, false, false)
	if err != nil {
		return fmt.Errorf("generate an initial password: %w", err)
	}
	logE.Info("generate an initial password")
	// create a login profile
	iamSvc := iam.New(sess, aws.NewConfig().WithRegion(ctrl.Config.Region))
	if _, err := iamSvc.CreateLoginProfileWithContext(ctx, &iam.CreateLoginProfileInput{ //nolint:nestif
		Password:              aws.String(passwd),
		PasswordResetRequired: aws.Bool(true),
		UserName:              aws.String(param.UserName),
	}); err != nil {
		if ctrl.Config.WhenLoginProfileExist == "error" {
			return fmt.Errorf("create a login profile: %w", err)
		}
		var aerr awserr.Error
		if errors.As(err, &aerr) {
			if aerr.Code() != iam.ErrCodeEntityAlreadyExistsException {
				return fmt.Errorf("create a login profile: %w", err)
			}
			if ctrl.Config.WhenLoginProfileExist == "ignore" {
				return nil
			}
			if err := updateLoginProfile(ctx, logE, iamSvc, &iam.UpdateLoginProfileInput{
				Password:              aws.String(passwd),
				PasswordResetRequired: aws.Bool(true),
				UserName:              aws.String(param.UserName),
			}); err != nil {
				return fmt.Errorf("update a login profile: %w", err)
			}
			logE.Info("update a login profile")
		} else {
			return fmt.Errorf("create a login profile: %w", err)
		}
	} else {
		logE.Info("create a login profile")
	}
	// create an access key
	// if _, err := iamSvc.CreateAccessKeyWithContext(ctx, &iam.CreateAccessKeyInput{}); err != nil {
	// 	return err
	// }
	// create a message
	msg, err := ctrl.renderTemplate(ctrl.MessageTemplate, map[string]interface{}{
		"UserName":     param.UserName,
		"Password":     passwd,
		"AWSAccountID": ctrl.Config.AWSAccountID,
	})
	if err != nil {
		return fmt.Errorf("render a message: %w", err)
	}
	logE.Info("render a message")
	// send a message
	if _, _, _, err := ctrl.SlackBot.SendMessageContext(ctx, user.ID, slack.MsgOptionText(msg, false)); err != nil {
		return fmt.Errorf("send Slack DM to a created user(Slack User ID: %s): %w", user.ID, err)
	}
	logE.Info("send a message")
	return nil
}

func (ctrl *Controller) handleSystemUser(ctx context.Context, param Param) error {
	logE := logrus.WithFields(logrus.Fields{
		"user_name": param.UserName,
	})
	if ctrl.Config.Slack.ChannelIDForSystemUser == "" {
		logE.Debug("channel_id_for_system_user isn't set")
		return nil
	}
	// treat the user as a system account
	// send a notification to slack
	msg, err := ctrl.renderTemplate(ctrl.MessageTemplateForSystemUser, map[string]interface{}{
		"UserName":     param.UserName,
		"AWSAccountID": ctrl.Config.AWSAccountID,
	})
	if err != nil {
		return err
	}
	if _, _, _, err := ctrl.SlackBot.SendMessageContext(ctx, ctrl.Config.Slack.ChannelIDForSystemUser, slack.MsgOptionText(msg, false)); err != nil {
		return fmt.Errorf("send a notification that a system user has been created to Slack channel (channel id: %s): %w", ctrl.Config.Slack.ChannelIDForSystemUser, err)
	}
	logE.WithFields(logrus.Fields{
		"channel_id": ctrl.Config.Slack.ChannelIDForSystemUser,
	}).Info("send a notification that a system user has been created to Slack channel")
	return nil
}

type IAM interface {
	UpdateLoginProfileWithContext(ctx aws.Context, input *iam.UpdateLoginProfileInput, opts ...request.Option) (*iam.UpdateLoginProfileOutput, error)
}

func updateLoginProfile(ctx context.Context, logE *logrus.Entry, iamSvc IAM, input *iam.UpdateLoginProfileInput) error {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()
	for {
		if _, err := iamSvc.UpdateLoginProfileWithContext(ctx, input); err != nil {
			var aerr awserr.Error
			if errors.As(err, &aerr) {
				if aerr.Code() == iam.ErrCodeEntityTemporarilyUnmodifiableException {
					// retry
					select {
					case <-ctx.Done():
						return fmt.Errorf("context is cancelled: %w", ctx.Err())
					case <-time.After(1 * time.Second):
						logE.Info("retry to update a login profile")
						continue
					}
				}
			}
			return fmt.Errorf("update a login profile: %w", err)
		}
		return nil
	}
}
