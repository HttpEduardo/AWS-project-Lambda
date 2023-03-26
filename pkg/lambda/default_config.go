package lambda

import "github.com/suzuki-shunsuke/aws-iam-cred-sender/pkg/controller"

const (
	systemUserCreatedMsg = "AWS IAM User for system has been created: `{{.UserName}}`"
	userCreatedMsg       = `Your AWS IAM User has been created.
AWS Account ID: ` + "`{{.AWSAccountID}}`" + `
User Name: ` + "`{{.UserName}}`" + `
Initial password:

` + "```" + `
{{.Password}}
` + "```" + `

Please sign in AWS and change your password.

https://{{.AWSAccountID}}.signin.aws.amazon.com/console

And please create your AWS Access Key if needed.
https://console.aws.amazon.com/iam/home#/users/{{.UserName}}?section=security_credentials
`
)

func (handler *Handler) setDefaultConfig(cfg *controller.Config) {
	if cfg.InitialPasswordLength == 0 {
		cfg.InitialPasswordLength = 32
	}
	if cfg.Message == "" {
		cfg.Message = userCreatedMsg
	}
	if cfg.MessageForSystemUser == "" {
		cfg.MessageForSystemUser = systemUserCreatedMsg
	}
	if cfg.WhenLoginProfileExist == "" {
		cfg.WhenLoginProfileExist = "change_password"
	}
	if cfg.DynamoDBTableName == "" {
		cfg.DynamoDBTableName = "aws-iam-cred-sender"
	}
	if cfg.DynamoDBTTL == 0 {
		cfg.DynamoDBTTL = 600 // default ttl: 600 seconds (10 minutes)
	}
}
