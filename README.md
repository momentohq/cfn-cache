# Momento::SimpleCache::Cache

Official cloudformation resource to manage momento serverless caches.

Please see our [Docs](docs/README.md) page for full end user docs and usage instructions for using this resource in your templates.

## Registering and Using this resource in your account

To use this resource in your account it currently needs to be registered as a private cloudformation extension 
per AWS account and region you want to use it in. You will need to install the AWS 
[CFN CLI](https://docs.aws.amazon.com/cloudformation-cli/latest/userguide/what-is-cloudformation-cli.html) locally
to submit the resource.


Once you have the CFN-CLI installed and have AWS credentials for the account you want to deploy you can submit this 
resource to be registered in your account in region of your choice.
```
make build
cfn submit  --set-default -v --region us-east-1
```

### Usage

Once you have the `Momento::SimpleCache::Cache` cfn extension installed you can start provisioning caches in your account

Bootstrap account with your momento auth token in secrets manager
```
export MOMENTO_AUTH_TOKEN=eyjbTestToken
aws secretsmanager create-secret \
    --name /momento/authToken \
    --secret-string $MOMENTO_AUTH_TOKEN
```
_replace $MOMENTO_AUTH_TOKEN value with token received [during signup](https://docs.momentohq.com/docs/overview)_

Create `test.yml` file with following contents
```yaml
AWSTemplateFormatVersion: 2010-09-09

Resources:
  MyCache:
    Type: Momento::SimpleCache::Cache
    Properties:
      Name: test-cache
      AuthToken: '{{resolve:secretsmanager:/momento/authToken}}'
```

Deploy test stack

```console
aws cloudformation create-stack \
    --region us-west-2 \
    --template-body "file://test.yml" \
    --stack-name "test-cache-stack"
```


### Development

If you want to contribute to this repo and develop on this resource please follow these instructions

**Pre-Reqs:**
1. Python version 3.6 or above.
2. [AWS CLI](https://aws.amazon.com/cli/)
3. [SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html)
4. [CFN CLI](https://docs.aws.amazon.com/cloudformation-cli/latest/userguide/what-is-cloudformation-cli.html)

Build resource
```
make
```

Start Local lambda
```
sam local start-lambda
```

Run Tests
```
cfn test
```
