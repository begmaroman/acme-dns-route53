## ACME-DNS-ROUTE53

Tool for managing SSL server certificates using [Let's Encrypt CA](https://letsencrypt.org/) and [Amazon Web Services](https://aws.amazon.com/).

### Features:

- Register with CA
- Creating the initial server certificate
- Renewing already existing certificates
- Support DNS-01 challenge using [Route53](https://aws.amazon.com/route53/) by AWS
- Store certificates into [ACM](https://aws.amazon.com/certificate-manager/) by AWS

### Installation:

Make sure that [GoLang](https://golang.org/doc/install) already installed

    go install github.com/begmaroman/acme-dns-route53

### Usage:

- Domains - TBD.

- Let's Encrypt Email - TBD.

- Let’s Encrypt ACME server - defaults to communicating with the production Let’s Encrypt ACME server. 
If you'd like to test something without issuing real certificates, consider using  `--staging` flag: 
    ```sh
    acme-dns-route53 obtain --staging --domains=<domains> --email=<email>
    ```
    
- Configuration directory - defaults the configuration data storing in the current directory (where the CLI runs).
If you'd like to change config directory, set the desired path using `--config-dir` flag:
    ```sh
    acme-dns-route53 obtain --config-path=<config-dir-path> --domains=<domains> --email=<email>
    ```
    
### Links:

Let's Encrypt Website: [https://letsencrypt.org](https://letsencrypt.org)

Amazon Certificate Manager: [https://aws.amazon.com/certificate-manager](https://aws.amazon.com/certificate-manager)

Route53 by Amazon: [https://aws.amazon.com/route53](https://aws.amazon.com/route53)

### Dependencies:

- [github.com/go-acme/lego](https://github.com/go-acme/lego) - Let's Encrypt client
- [github.com/aws/aws-sdk-go](https://github.com/aws/aws-sdk-go) - AWS SDK to manage certificates