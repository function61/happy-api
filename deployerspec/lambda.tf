
variable "zip_filename" { type = "string" }
variable "region" { type = "string" }
variable "s3_bucket" { type = "string" }
variable "base_url" { type = "string" }

provider "aws" {
	region = "${var.region}"
}

resource "aws_lambda_function" "onni" {
	function_name = "Onni"
	description = "REST API for delivering happiness"

	# won't get updated to Lambda unless the path changes every time - that's why we need
	# to embed version in filename and make this a variable
	filename = "${var.zip_filename}"

	handler = "onni"

	# https://docs.aws.amazon.com/lambda/latest/dg/current-supported-versions.html
	runtime = "go1.x"

	# FIXME
	role = "arn:aws:iam::329074924855:role/AlertManager"

	timeout = 10

	environment {
		variables = {
			S3_BUCKET = "${var.s3_bucket}"
			BASE_URL = "${var.base_url}"
		}
	}
}

resource "aws_s3_bucket" "mediabucket" {
  bucket = "${var.s3_bucket}"
  acl    = "public-read"
  policy = "{\"Version\": \"2008-10-17\", \"Statement\": [{ \"Sid\": \"AllowPublicRead\", \"Effect\": \"Allow\", \"Principal\": { \"AWS\": \"*\" }, \"Action\": \"s3:GetObject\", \"Resource\": \"arn:aws:s3:::${var.s3_bucket}/*\" } ]}"
}

module "apigateway" {
	lambda_name = "${aws_lambda_function.onni.function_name}"
  source = "./apigatewaybullshit"
}

output "base_url" {
	value = "${module.apigateway.invoke_url}"
}
