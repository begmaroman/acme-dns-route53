### Usage by AWS Lambda:

#### Setting up the AWS CLI:

1. Throughout this instruction, we'll use the AWS CLI (command line interface) to configure our lambda functions and other AWS services. 
   Installation and basic usage instructions can be found [here](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-welcome.html), but if you’re using a Debian-based system like Ubuntu you can install the CLI with apt and run it using the aws command:

   ```bash
   $ sudo apt install awscli
   $ aws --version
   aws-cli/1.16.47 Python/3.6.3 Linux/4.15.0-47-generic botocore/1.12.37
   ```
   
2. Next we need to set up an AWS IAM user with *programmatic access permission* for the CLI to use. 
   A guide on how to do this can be found [here](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_users_create.html). 
   For testing purposes you can attach the all-powerful `AdministratorAccess` managed policy to this user, but in practice I would recommend using a more restrictive policy. 
   At the end of setting up the user you'll be given a *access key ID* and *secret access key*. 
   Make a note of these — you’ll need them in the next step.
   
3. Configure the CLI to use the credentials of the IAM user you've just created using the `aws configure` command. 
   You’ll also need to specify the [default region](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Concepts.RegionsAndAvailabilityZones.html) and [output format](https://docs.aws.amazon.com/cli/latest/userguide/controlling-output.html) you want the CLI to use.

   ```bash
   $ aws configure
   AWS Access Key ID [None]: access-key-ID
   AWS Secret Access Key [None]: secret-access-key
   Default region name [None]: us-east-1
   Default output format [None]: json
   ```

   (Throughout this instruction, I'll assume you're using the `us-east-1` region — you'll need to change the code snippets accordingly if you're using a different region.)

#### Creating and deploying an Lambda function:

---

The first step is to build an executable from remote repo using `go install`:
```bash
$ env GOOS=linux GOARCH=amd64 go install github.com/begmaroman/acme-dns-route53
```
The executable will be installed in `$GOPATH/bin` directory.
Important: as part of this command we're using env to temporarily set two environment variables for the duration for the command (GOOS=linux and GOARCH=amd64). 
These instruct the Go compiler to create an executable suitable for use with a linux OS and amd64 architecture — which is what it will be running on when we deploy it to AWS.

---

AWS requires us to upload our lambda functions in a zip file, so let's make a `acme-dns-route53.zip` zip file containing the executable we just made:
```bash
$ zip -j ~/acme-dns-route53.zip $GOPATH/bin/acme-dns-route53
```
***Note** that the executable must be in the root of the zip file — not in a folder within the zip file. 
To ensure this I've used the `-j` flag in the snippet above to junk directory names.*

---

The next step is a bit awkward, but critical to getting our lambda function working properly. 
We need to set up an IAM role which defines the permission that our lambda function will have when it is running. 

For now let's set up a `lambda-acme-dns-route53-executor` role and attach the `AWSLambdaBasicExecutionRole` managed policy to it. 
This will give our lambda function the basic permissions it need to run and log to the [AWS CloudWatch](https://aws.amazon.com/cloudwatch/) service.

First we have to create a trust policy JSON file. 
This will essentially instruct AWS to allow lambda services to assume the `lambda-acme-dns-route53-executor` role:

```
Filepath: ~/lambda-acme-dns-route53-executor-policy.json
```
```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "logs:CreateLogGroup"
            ],
            "Resource": "arn:aws:logs:<AWS_REGION>:<AWS_ACCOUNT_ID>:*"
        },
        {
            "Effect": "Allow",
            "Action": [
                "logs:PutLogEvents",
                "logs:CreateLogStream"
            ],
            "Resource": "arn:aws:logs:<AWS_REGION>:<AWS_ACCOUNT_ID>:log-group:/aws/lambda/acme-dns-route53:*"
        },
        {
            "Sid": "",
            "Effect": "Allow",
            "Action": [
                "route53:ListHostedZones",
                "cloudwatch:PutMetricData",
                "acm:ImportCertificate",
                "acm:ListCertificates"
            ],
            "Resource": "*"
        },
        {
            "Sid": "",
            "Effect": "Allow",
            "Action": [
                "sns:Publish",
                "route53:GetChange",
                "route53:ChangeResourceRecordSets",
                "acm:ImportCertificate",
                "acm:DescribeCertificate"
            ],
            "Resource": [
                "arn:aws:sns:${var.region}:<AWS_ACCOUNT_ID>:alarm_topic",
                "arn:aws:route53:::hostedzone/*",
                "arn:aws:route53:::change/*",
                "arn:aws:acm:<AWS_REGION>:<AWS_ACCOUNT_ID>:certificate/*"
            ]
        }
    ]
}
```

Then use the `aws iam create-role` command to create the role with this trust policy:

```bash
$ aws iam create-role --role-name lambda-acme-dns-route53-executor \
 --assume-role-policy-document ~/lambda-acme-dns-route53-executor-policy.json
```

***Note:** returned ARN (Amazon Resource Name) — you'll need this in the next step.*

Now the `lambda-acme-dns-route53-executor` role has been created we need to specify the permissions that the role has. 
The easiest way to do this it to use the `aws iam attach-role-policy` command, passing in the ARN of `AWSLambdaBasicExecutionRole` permission policy like so:

```bash
$ aws iam attach-role-policy --role-name lambda-acme-dns-route53-executor \
--policy-arn arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole
```

***Note:** you can find a list of other permission policies that might be useful [here](https://docs.aws.amazon.com/lambda/latest/dg/lambda-permissions.html#lambda-intro-execution-role).*

---

Now we're ready to actually deploy the lambda function to AWS, which we can do using the `aws lambda create-function` command.
Also, The lambda function needs to be configured with the following options:

   - `AWS_LAMBDA` environment variable with value `1` which adjusts the tool for using inside Lambda function.
   - `1024` MB as memory limit (can be changed if needed).
   - `900` secs (15 min) is the maximum timeout.
   - `acme-dns-route53` is the handler name of the lambda function
   - `fileb://~/acme-dns-route53.zip` is the created `.zip` file above.
   
Go ahead and try deploying it:

```
 $ aws lambda create-function \
 --function-name acme-dns-route53 \
 --runtime go1.x \
 --role arn:aws:iam::<AWS_ACCOUNT_ID>:role/lambda-acme-dns-route53-executor \
 --environment Variables="{AWS_LAMBDA=1}" \
 --memory-size 1024 \
 --timeout 900 \
 --handler acme-dns-route53 \
 --zip-file fileb://~/acme-dns-route53.zip

 {
     "FunctionName": "acme-dns-route53", 
     "LastModified": "2019-05-03T19:07:09.325+0000", 
     "RevisionId": "e3fadec9-2180-4bff-bb9a-999b1b71a558", 
     "MemorySize": 1024, 
     "Environment": {
         "Variables": {
             "AWS_LAMBDA": "1"
         }
     }, 
     "Version": "$LATEST", 
     "Role": "arn:aws:iam::<AWS_ACCOUNT_ID>:role/lambda-acme-dns-route53-executor", 
     "Timeout": 900, 
     "Runtime": "go1.x", 
     "TracingConfig": {
         "Mode": "PassThrough"
     }, 
     "CodeSha256": "+2KgE5mh5LGaOsni36pdmPP9O35wgZ6TbddspyaIXXw=", 
     "Description": "",
     "CodeSize": 8456317, 
     "FunctionArn": "arn:aws:lambda:us-east-1:<AWS_ACCOUNT_ID>:function:acme-dns-route53", 
     "Handler": "acme-dns-route53"
 }
```

---

So there it is. Our lambda function has been deployed and is now ready to use. 
First, needs to create JSON string with a configuration. Configuration structure:

| Field            | Type     | Description  |
|------------------|----------|--------------|
| `domains`        | []string | Domains list |
| `email`          | string   | [Let's Encrypt expiration Email](https://letsencrypt.org/docs/expiration-emails/) |
| `staging`        | string   | `1` for Let's Encrypt staging environment, and `0` for production one |
| `topic`          | string   | SNS Notification Topic ARN (optional) |
| `renew_before`   | int      | The number of days defining the period before expiration within which a certificate must be renewed |

Example of JSON configuration:

```json
{ 
  "domains":["example1.com","example2.com"],
  "email":"your@email.com",
  "staging":"1",
  "topic":"arn:aws:sns:<AWS_REGION>:<AWS_ACCOUNT_ID>:<SNS_TOPIC_NAME>",
  "renew_before":7
}
```

You can try it out by using the `aws lambda invoke` command (which requires you to specify an output file for the response — I've used `/tmp/output.json` in the snippet below).

```bash
$ aws lambda invoke \
 --function-name acme-dns-route53 \
 --payload "{\"domains\":[\"yourdomain.com\"],\"email\":\"your@email.com\",\"staging\":\"1\",\"topic\":\"arn:aws:sns:<AWS_REGION>:<AWS_ACCOUNT_ID>:<SNS_TOPIC_NAME>\",\"renew_before\":7}"
 /tmp/output.json
```

Then check logs on AWS CloudWatch, and obtained certificates on Amazon Certificate Manager.

And one **important** thing is that you can pass the parameters above (`domains`,`email`,`staging` etc.) via environment variables of the lambda function.
Environment variables has priority than payload.
Use the following environment variables to pass these parameters:
 
 - `DOMAINS` is the environment variable which contains comma-separated domains list. Equivalent to `domains` field in the payload object.
 - `LETSENCRYPT_EMAIL` is the environment variable which contains [Let's Encrypt expiration Email](https://letsencrypt.org/docs/expiration-emails/). Equivalent to `email` field in the payload object.
 - `STAGING` is the environment variable which must contain 1 value for using staging Let’s Encrypt environment or 0 for production environment. Equivalent to `staging` field in the payload object.
 - `NOTIFICATION_TOPIC` is the environment variable which contains SNS Notification Topic ARN.
 - `RENEW_BEFORE` is the number of days defining the period before expiration within which a certificate must be renewed.
