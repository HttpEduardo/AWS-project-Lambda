# aws-iam-cred-sender

[![Build Status](https://github.com/suzuki-shunsuke/aws-iam-cred-sender/workflows/build/badge.svg)](https://github.com/suzuki-shunsuke/aws-iam-cred-sender/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/suzuki-shunsuke/aws-iam-cred-sender)](https://goreportcard.com/report/github.com/suzuki-shunsuke/aws-iam-cred-sender)
[![GitHub last commit](https://img.shields.io/github/last-commit/suzuki-shunsuke/aws-iam-cred-sender.svg)](https://github.com/suzuki-shunsuke/aws-iam-cred-sender)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/suzuki-shunsuke/aws-iam-cred-sender/main/LICENSE)

AWS Lambda Function to send an initial password to a new user via Slack DM

## Assumption

IAM User name and Slack's display name must be same.

## Overview

When an IAM User is created, the Lambda Function is triggered via CloudWatch Event.
The function searches the Slack User with IAM User name.
If the IAM User isn't found, the function notifies that the IAM User is created to a Slack channel.

![image](https://user-images.githubusercontent.com/13323303/114290928-3ba40200-9abe-11eb-8f9b-72b3680d4a1e.png)

If the IAM User is found, the function creates the IAM User's login profile and sends the initial password to the user via Slack DM.

![image](https://user-images.githubusercontent.com/13323303/114290993-bbca6780-9abe-11eb-9efe-ff2376400a96.png)

If the login profile already exists, the password is changed by default. This behavior can be changed.

## Architecture

![aws-iam-cred-sender](https://user-images.githubusercontent.com/13323303/115412477-4cc3e000-a22f-11eb-9d28-062f996b9697.png)

## Configuration

Please see [here](docs/configuration.md)

### Lambda Execution Role

Please see [here](docs/lambda-execution-role.md)

### Slack App Permission

* chat:write (chat.postMessage)
* users:read (users.list)

### DynamoDB

This function needs a DynamoDB table.
Please see [here](docs/dynamodb.md).

## LICENSE

[MIT](LICENSE)
# AWS-project-Lambda
