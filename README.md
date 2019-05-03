## ACME-DNS-ROUTE53

Tool for managing SSL server certificates using [Let's Encrypt CA](https://letsencrypt.org/) and [Amazon Web Services](https://aws.amazon.com/).

### Features:

- Register with CA
- Creating the initial server certificate
- Renewing already existing certificates
- Support DNS-01 challenge using [Route53](https://aws.amazon.com/route53/) by AWS
- Store certificates into [ACM](https://aws.amazon.com/certificate-manager/) by AWS
- Managing certificates of multiple domains within one request

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
    "Version": "2012-10-17",
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
    acme-dns-route53 obtain --domains=testserver.com,testserver1.com,testserver2.com --email=<email>
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
    acme-dns-route53 obtain --staging --domains=<domains> --email=<email>
    ```
    
- Configuration directory - defaults the configuration data storing in the current directory (where the CLI runs).
If you'd like to change config directory, set the desired path using **`--config-dir`** flag:
    ```sh
    acme-dns-route53 obtain --config-path=<config-dir-path> --domains=<domains> --email=<email>
    ```
    
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