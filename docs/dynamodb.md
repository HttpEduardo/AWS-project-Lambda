# DynamoDB

[#12](https://github.com/suzuki-shunsuke/aws-iam-cred-sender/pull/12)
[#15](https://github.com/suzuki-shunsuke/aws-iam-cred-sender/pull/15)

Sometimes Lambda Function is called at multiple times by CloudWatch Event.

https://docs.aws.amazon.com/AmazonCloudWatch/latest/events/CWE_Troubleshooting.html#RuleTriggeredMoreThanOnce

So we use DynamoDB to handle multiple function call.
The function registers the User Name at DynamoDB, and if the User Name is already registered at DynamoDB table the function aborts the request.

```json
{
    "Table": {
        "AttributeDefinitions": [
            {
                "AttributeName": "UserName",
                "AttributeType": "S"
            }
        ],
        "TableName": "aws-iam-cred-sender",
        "KeySchema": [
            {
                "AttributeName": "UserName",
                "KeyType": "HASH"
            }
        ]
    }
}
```

```json
{
    "TimeToLiveDescription": {
        "TimeToLiveStatus": "ENABLED",
        "AttributeName": "ExpiredAt"
    }
}
```
