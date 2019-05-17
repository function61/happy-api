
variable "zip_filename" { type = "string" }
variable "region" { type = "string" }

provider "aws" {
	region = "${var.region}"
}

resource "aws_lambda_function" "onni" {
	function_name = "Onni"
	# description = "Checks that important web properties are working."

	# won't get updated to Lambda unless the path changes every time - that's why we need
	# to embed version in filename and make this a variable
	filename = "${var.zip_filename}"

	handler = "onni"

	# https://docs.aws.amazon.com/lambda/latest/dg/current-supported-versions.html
	runtime = "go1.x"

	# FIXME
	role = "arn:aws:iam::329074924855:role/AlertManager"

	timeout = 60
}

module "apigateway" {
	lambda_name = "${aws_lambda_function.onni.function_name}"
	lambda_arn = "${aws_lambda_function.onni.arn}"
	lambda_invoke_arn = "${aws_lambda_function.onni.invoke_arn}"
  source = "./apigatewaybullshit"
}

output "base_url" {
	value = "${module.apigateway.invoke_url}"
}
