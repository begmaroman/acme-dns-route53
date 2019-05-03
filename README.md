**acme-dns-route53** is the tool to obtain SSL certificates from [Let's Encrypt CA](https://letsencrypt.org/) using DNS-01 challenge with Route53 and Amazon Certificate Manager by [AWS](https://aws.amazon.com/).

### Features:

- Register with CA
- Creating the initial server certificate
- Renewing already existing certificates
- Support DNS-01 challenge using [Route53](https://aws.amazon.com/route53/) by AWS
- Store certificates into [ACM](https://aws.amazon.com/certificate-manager/) by AWS
- Managing certificates of multiple domains within one request
- Build-in [AWS Lambda](https://aws.amazon.com/lambda/) tolerance

### Installation:

Make sure that [GoLang](https://golang.org/doc/install) already installed

    go install github.com/begmaroman/acme-dns-route53
    
### Credentials:

Use of this tool requires a configuration file containing Amazon Web Sevices API credentials for an account with the following permissions:

- `route53:ListHostedZones`
- `route53:GetChange`
- `route53:ChangeResourceRecordSets`
- `acm:ImportCertificate`
- `acm:ListCertificates`
- `acm:DescribeCertificate`

These permissions can be captured in an AWS policy like the one below. 
Amazon provides [information about managing](https://docs.aws.amazon.com/Route53/latest/DeveloperGuide/access-control-overview.html) access and [information about the required permissions](https://docs.aws.amazon.com/Route53/latest/DeveloperGuide/r53-api-permissions-ref.html)

*Example AWS policy file:*
```json
{
    "Version": "2019-05-01",
    "Statement": [
        {
            "Sid": "",
            "Effect": "Allow",
            "Action": [
                "route53:ListHostedZones",
                "acm:ImportCertificate",
                "acm:ListCertificates"
            ],
            "Resource": "*"
        },
        {
            "Sid": "",
            "Effect": "Allow",
            "Action": [
                "route53:GetChange",
                "route53:ChangeResourceRecordSets",
                "acm:ImportCertificate",
                "acm:DescribeCertificate"
            ],
            "Resource": [
                "arn:aws:route53:::hostedzone/<HOSTED_ZONE_ID>",
                "arn:aws:route53:::change/*",
                "arn:aws:acm:us-east-1:<AWS_ACCOUNT_ID>:certificate/*"
            ]
        }
    ]
}
```

The [access keys](https://docs.aws.amazon.com/general/latest/gr/aws-sec-cred-types.html#access-keys-and-secret-access-keys) for an account with these permissions must be supplied in one of the following ways:

- Using the `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` environment variables.
- Using a credentials configuration file at the default location, `~/.aws/config`.
- Using a credentials configuration file at a path supplied using the `AWS_CONFIG_FILE` environment variable.

*Example credentials config file:*
```
[default]
aws_access_key_id=AKIAIOSFODNN7EXAMPLE
aws_secret_access_key=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
```

### Usage:

- Domains (required) - use **`--domains`** flag to determine comma-separated domains list, certificates of which should be obtained. Example:
    ```sh
    $ acme-dns-route53 obtain --domains=testserver.com,testserver1.com,testserver2.com --email=<email>
    ```

- Let's Encrypt Email (required) - use **`--email`** flag to determine Let's Encrypt account email. 
If account's private key is not provided, registers a new account. Private key expected by path `<config-dir>/<email>.pem`. Example:
    
    Path: `/tmp/letsencrypt/test@test.test.pem`
    
    Content:
    ```pem
    -----BEGIN RSA PRIVATE KEY-----
    somecontentoftheprivatekey
    -----END RSA PRIVATE KEY-----
    ```

- Let’s Encrypt ACME server - defaults to communicating with the production Let’s Encrypt ACME server. 
If you'd like to test something without issuing real certificates, consider using  **`--staging`** flag: 
    ```sh
    $ acme-dns-route53 obtain --staging --domains=<domains> --email=<email>
    ```
    
- Configuration directory - defaults the configuration data storing in the current directory (where the CLI runs).
If you'd like to change config directory, set the desired path using **`--config-dir`** flag:
    ```sh
    $ acme-dns-route53 obtain --config-path=<config-dir-path> --domains=<domains> --email=<email>
    ```
    
### Usage by AWS Lambda:

1. The first step is to build an executable from remote repo using `go install`:
    ```bash
    $ env GOOS=linux GOARCH=amd64 go install github.com/begmaroman/acme-dns-route53
    ```
    The executable will be installed in `$GOPATH/bin` directory.
    Important: as part of this command we're using env to temporarily set two environment variables for the duration for the command (GOOS=linux and GOARCH=amd64). 
    These instruct the Go compiler to create an executable suitable for use with a linux OS and amd64 architecture — which is what it will be running on when we deploy it to AWS.

2. AWS requires us to upload our lambda functions in a zip file, so let's make a `acme-dns-route53.zip` zip file containing the executable we just made:
    ```bash
    $ zip -j ~/acme-dns-route53.zip $GOPATH/bin/acme-dns-route53
    ```
    Note that the executable must be in the root of the zip file — not in a folder within the zip file. To ensure this I've used the `-j` flag in the snippet above to junk directory names.
    
3. The next step is a bit awkward, but critical to getting our lambda function working properly. 
   We need to set up an IAM role which defines the permission that our lambda function will have when it is running. 
   
   For now let's set up a `lambda-acme-dns-route53-executor` role and attach the `AWSLambdaVPCAccessExecutionRole` managed policy to it. 
   This will give our lambda function the basic permissions it need to run and log to the [AWS CloudWatch](https://aws.amazon.com/cloudwatch/) service.
   
   First we have to create a trust policy JSON file. 
   This will essentially instruct AWS to allow lambda services to assume the `lambda-acme-dns-route53-executor` role:

   ```
   Filepath: ~/lambda-acme-dns-route53-executor-policy.json
   ```
   ```json
   {
       "Version": "2019-05-01",
       "Statement": [
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
                   "route53:GetChange",
                   "route53:ChangeResourceRecordSets",
                   "acm:ImportCertificate",
                   "acm:DescribeCertificate"
               ],
               "Resource": [
                   "arn:aws:route53:::hostedzone/<HOSTED_ZONE_ID>",
                   "arn:aws:route53:::change/*",
                   "arn:aws:acm:us-east-1:<AWS_ACCOUNT_ID>:certificate/*"
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
   
   Make a note of the returned ARN (Amazon Resource Name) — you'll need this in the next step.
   
   Now the `lambda-acme-dns-route53-executor` role has been created we need to specify the permissions that the role has. 
   The easiest way to do this it to use the `aws iam attach-role-policy` command, passing in the ARN of `AWSLambdaVPCAccessExecutionRole` permission policy like so:
   
   ```bash
   $ aws iam attach-role-policy --role-name lambda-acme-dns-route53-executor \
   --policy-arn arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole
   ```
   
   Note: you can find a list of other permission policies that might be useful [here](https://docs.aws.amazon.com/lambda/latest/dg/lambda-permissions.html#lambda-intro-execution-role).
   
4. Now we're ready to actually deploy the lambda function to AWS, which we can do using the `aws lambda create-function` command.
   Go ahead and try deploying it:
   
   ```bash
    $ aws lambda create-function --function-name acme-dns-route53 --runtime go1.x \
    --role arn:aws:iam::<AWS_ACCOUNT_ID>:role/lambda-acme-dns-route53-executor \
    --handler acme-dns-route53 --zip-file ~/acme-dns-route53.zip
   ```
   
5. So there it is. Our lambda function has been deployed and is now ready to use. 
   You can try it out by using the `aws lambda invoke` command (which requires you to specify an output file for the response — I've used `/tmp/output.json` in the snippet below).
   
   ```bash
   $ aws lambda invoke --function-name acme-dns-route53 /tmp/output.json
   ```
   
   Then check logs on AWS CloudWatch, and obtained certificates on Amazon Certificate Manager.
   
### Links:

Let's Encrypt Website: [https://letsencrypt.org](https://letsencrypt.org)

Community: [https://community.letsencrypt.org](https://community.letsencrypt.org)

Amazon Certificate Manager: [https://aws.amazon.com/certificate-manager](https://aws.amazon.com/certificate-manager)

Route53 by Amazon: [https://aws.amazon.com/route53](https://aws.amazon.com/route53)

ACME spec: [http://ietf-wg-acme.github.io/acme/](http://ietf-wg-acme.github.io/acme/)

### Dependencies:

- [github.com/go-acme/lego](https://github.com/go-acme/lego) - Let's Encrypt client
- [github.com/aws/aws-sdk-go](https://github.com/aws/aws-sdk-go) - AWS SDK to manage certificates

### Inspired by:

- [https://arkadiyt.com/2018/01/26/deploying-effs-certbot-in-aws-lambda/](https://arkadiyt.com/2018/01/26/deploying-effs-certbot-in-aws-lambda/)
- Lego library - [https://github.com/go-acme/lego](https://github.com/go-acme/lego)