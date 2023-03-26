# Configuration

## Envrionment Variable

`CONFIG`

Default Config

```yaml
log_level: info
message_for_system_user: "AWS IAM User for system has been created: `{{.UserName}}`"
message: |
  Your AWS IAM User has been created.
  AWS Account ID: `{{.AWSAccountID}}`
  User Name: `{{.UserName}}`
  Initial password:

  {{.Password}}

  Please sign in AWS and change your password.

  https://{{AWSAccountID}}.signin.aws.amazon.com/console

  And please create your AWS Access Key if needed.
  https://console.aws.amazon.com/iam/home#/users/{{.UserName}}?section=security_credentials
initial_password_length: 32
slack:
  channel_id_for_system_user: ""
secret_version_id: ""
# change_password, ignore, error
when_login_profile_exist: change_password
dynamodb_table_name: aws-iam-cred-sender
dynamodb_ttl: 600
```

Required

```yaml
secret_id: xxx
region: xxx
```

#### Template

The template is rendered with Go's [text/template](https://golang.org/pkg/text/template/).
And [sprig Function](http://masterminds.github.io/sprig/) can be used.

Template variables:

* AWSAccountID
* UserName
* Password

### AWS Secrets Manager's Secret

* slack_bot_token
