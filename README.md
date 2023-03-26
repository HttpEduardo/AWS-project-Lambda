## AWS

Função AWS Lambda para enviar uma senha inicial para um novo usuário via Slack DM

## Assumption

O nome de usuário do IAM e o nome de exibição do Slack devem ser iguais.

## Overview

Quando um usuário IAM é criado, a função Lambda é acionada por meio do CloudWatch Event.
A função pesquisa o usuário do Slack com o nome de usuário do IAM.
Caso o usuário IAM não seja encontrado, a função notifica que o usuário IAM foi criado para um canal Slack.

Se o usuário IAM for encontrado, a função cria o perfil de login do usuário IAM e envia a senha inicial para o usuário via Slack DM.

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
