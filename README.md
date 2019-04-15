# Task Kitchen

## Prerequisite

### Tools

- go >= 1.12
- GNU Make >= 3.81
- yarn >= 1.15.2
  - npx >= 6.7.0
- awscli >= 1.16.130

### Resources

- Your domain name (e.g. `task-kitchen.example.com` )
- AWS CLI credentials for API (e.g. environment variable `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` )

## Development

### Setup

### API server

Run server with AWS credentials.

```bash
$ go run ./server/ <your-region> <your-dynamodb-name>
```

### Content server

```bash
$ npx webpack-dev-server --config ./webpack.config.js --hot
```

Then, open http://localhost:8080/

## Deployment

Create a configuration file as `config.json` like following.

```json
{
  "StackName": "your-task-kitchen-stack-name",
  "CodeS3Bucket": "your-bucket-name-for-artifact",
  "CodeS3Prefix": "path/to/artifact",
  "Region": "your-region",

  "ServiceDomainName": "task-kitchen.example.com",
  "ServiceCertArn": "arn:aws:acm:us-east-1:1234567890:certificate/f5e65492-77b9-40a9-abdd-78fa49493d46"
}
```

Run deploy task.

```bash
$ env STACK_CONFIG=config.json make deploy
```