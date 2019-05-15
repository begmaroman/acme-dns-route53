variable "account_id" {
  description = "The ID of the AWS account in which resources are created"
}

variable "region" {
  description = "The region in which services are deployed"
}

variable "s3_source_bucket" {
  description = "S3 bucket containing code that doesn't get deployed through ECR, e.g. lambdas"
}

variable "acme_dns_route53_tag" {
  description = "Tag for the acme-dns-route53 lambda. If not set, the currently deployed version will be used. If that's also missing, a failing dummy lambda is used."
  default     = ""
}

variable "domains" {
  description = "Comma-separated domains list for which certificates will be issued."
  default     = ""
}

variable "letsencrypt_email" {
  description = "The Email which uses for getting an account in Let's Encrypt."
  default     = ""
}

// The names given for the items in the locals block must be unique throughout a module,
// e.g. they must be different from acme_dns_route53.tf
locals {
  # The lambda function name
  acme_dns_route53_function_name = "acme-dns-route53"

  # The S3 object key containing the lambda source code
  acme_dns_route53_src_s3_key = "acme-dns-route53/acme-dns-route53_${var.acme_dns_route53_tag}.zip"

  # The name of topin in SNS
  acme_dns_route53_sns_topic = "acme-dns-route53"
}

#--------------------------------------------------------------
# Certificate Issuer triggering event
#--------------------------------------------------------------

# Cloudwatch event rule that runs acme-dns-route53 lambda every 12 hours
resource "aws_cloudwatch_event_rule" "acme_dns_route53_sheduler" {
  name                = "acme-dns-route53-scheduler"
  schedule_expression = "cron(0 */12 * * ? *)"
}

# Specify the lambda function to run
resource "aws_cloudwatch_event_target" "acme_dns_route53r_sheduler_target" {
  rule = "${aws_cloudwatch_event_rule.acme_dns_route53_sheduler.name}"
  arn  = "${aws_lambda_function.acme_dns_route53.arn}"
}

# Give CloudWatch permission to invoke the function
resource "aws_lambda_permission" "permission" {
  action        = "lambda:InvokeFunction"
  function_name = "${aws_lambda_function.acme_dns_route53.function_name}"
  principal     = "events.amazonaws.com"
  source_arn    = "${aws_cloudwatch_event_rule.acme_dns_route53_sheduler.arn}"
}

#--------------------------------------------------------------
# Main Lambda Function
#--------------------------------------------------------------

# The source code hash of the file to upload.
# This helps Terraform determine whether the lambda function needs re-uploading.
data "aws_s3_bucket_object" "lambda_acme_dns_route53_hash" {
  bucket = "${var.s3_source_bucket}"
  key    = "${local.acme_dns_route53_src_s3_key}.base64sha256"
}

# Main lambda function which will runs by CloudWatch and send notifications to SNS
resource "aws_lambda_function" "acme_dns_route53" {
  s3_bucket        = "${var.s3_source_bucket}"
  s3_key           = "${local.acme_dns_route53_src_s3_key}"
  source_code_hash = "${data.aws_s3_bucket_object.lambda_acme_dns_route53_hash.body}"
  publish          = "true"

  function_name = "${local.acme_dns_route53_function_name}"
  role          = "${aws_iam_role.lambda_acme_dns_route53_executor.arn}"
  handler       = "acme-dns-route53"
  runtime       = "go1.x"
  memory_size   = 1024
  timeout       = 900

  environment {
    variables = {
      DOMAINS             = "${var.domains}"
      LETSENCRYPT_EMAIL   = "${var.letsencrypt_email}"
      STAGING             = 0
      NOTIFICATION_TOPIC  = "arn:aws:sns:${var.region}:${var.account_id}:${local.acme_dns_route53_sns_topic}"
    }
  }

  tags {
    Name           = "${local.acme_dns_route53_function_name}"
    SourceS3Key    = "${local.acme_dns_route53_src_s3_key}"
  }

  depends_on = [
    "aws_iam_role_policy.lambda_acme_dns_route53_executor"
  ]
}

# This allows someone to verify which version of the lambda is currently deployed
resource "aws_lambda_alias" "acme_dns_route53_latest_alias" {
  name             = "latest"
  description      = "Lambda alias pointing to the currently deployed acme-dns-route53 lambda version."
  function_name    = "${aws_lambda_function.acme_dns_route53.arn}"
  function_version = "${aws_lambda_function.acme_dns_route53.version}"

  depends_on = [
    "aws_lambda_function.acme_dns_route53"
  ]
}

#--------------------------------------------------------------
# IAM Role and Policy
#--------------------------------------------------------------

resource "aws_iam_role" "lambda_acme_dns_route53_executor" {
  name        = "lambda-acme-dns-route53-executor"
  description = "Allows Lambda Function to call AWS services on your behalf."

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

# Provides permissions to the lambda
# - Write permissions to CloudWatch logs
# - Permissions to read and import certificates to ACM
# - Permissions to create and delete records in Route53
# - Permissions to publish messages to SNS
resource "aws_iam_role_policy" "lambda_acme_dns_route53_executor" {
  name = "acme-dns-route53"
  role = "${aws_iam_role.lambda_acme_dns_route53_executor.id}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
          "logs:CreateLogGroup"
      ],
      "Resource": "arn:aws:logs:${var.region}:${var.account_id}:*"
    },
    {
      "Effect": "Allow",
      "Action": [
        "logs:PutLogEvents",
        "logs:CreateLogStream"
      ],
      "Resource": "arn:aws:logs:${var.region}:${var.account_id}:log-group:/aws/lambda/${local.acme_dns_route53_function_name}:*"
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
        "route53:GetChange",
        "route53:ChangeResourceRecordSets",
        "acm:ImportCertificate",
        "acm:DescribeCertificate"
      ],
      "Resource": [
        "arn:aws:sns:${var.region}:${var.account_id}:*",
        "arn:aws:route53:::hostedzone/*",
        "arn:aws:route53:::change/*",
        "arn:aws:acm:${var.region}:${var.account_id}:certificate/*"
      ]
    }
  ]
}
EOF
}
