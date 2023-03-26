
[![Build Status](https://github.com/suzuki-shunsuke/aws-iam-cred-sender/workflows/build/badge.svg)](https://github.com/suzuki-shunsuke/aws-iam-cred-sender/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/suzuki-shunsuke/aws-iam-cred-sender)](https://goreportcard.com/report/github.com/suzuki-shunsuke/aws-iam-cred-sender)
[![GitHub last commit](https://img.shields.io/github/last-commit/suzuki-shunsuke/aws-iam-cred-sender.svg)](https://github.com/suzuki-shunsuke/aws-iam-cred-sender)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/suzuki-shunsuke/aws-iam-cred-sender/main/LICENSE)

Função AWS Lambda para enviar uma senha inicial para um novo usuário via Slack DM

## Assumption

O nome de usuário do IAM e o nome de exibição do Slack devem ser iguais.

## Overview

Quando um usuário IAM é criado, a função Lambda é acionada por meio do CloudWatch Event.
A função pesquisa o usuário do Slack com o nome de usuário do IAM.
Caso o usuário IAM não seja encontrado, a função notifica que o usuário IAM foi criado para um canal Slack.

![image](https://user-images.githubusercontent.com/13323303/114290928-3ba40200-9abe-11eb-8f9b-72b3680d4a1e.png)

Se o usuário IAM for encontrado, a função cria o perfil de login do usuário IAM e envia a senha inicial para o usuário via Slack DM.

![image](https://user-images.githubusercontent.com/13323303/114290993-bbca6780-9abe-11eb-9efe-ff2376400a96.png)
Se o perfil de login já existir, a senha será alterada por padrão. Esse comportamento pode ser alterado.

## Architecture

![aws-iam-cred-sender](https://user-images.githubusercontent.com/13323303/115412477-4cc3e000-a22f-11eb-9d28-062f996b9697.png)

## Configuration

Aqui [here](docs/configuration.md)

### Lambda Execution Role

Aqui [here](docs/lambda-execution-role.md)

### Slack App Permission

* chat:write (chat.postMessage)
* users:read (users.list)

### DynamoDB

This function needs a DynamoDB table.
Please see [here](docs/dynamodb.md).

## LICENSE

[MIT](LICENSE)
# AWS-project-Lambda
