
# inputs
variable "lambda_name" {}

data "aws_lambda_function" "fn" {
  function_name = "${var.lambda_name}"
  qualifier     = "" # workaround to ensure the arn doesn't container a qualifier
}

resource "aws_api_gateway_rest_api" "restapi" {
	name        = "${var.lambda_name}"
	description = "Boilerplate for Lambda proxy"
}

resource "aws_api_gateway_resource" "proxyres" {
	rest_api_id = "${aws_api_gateway_rest_api.restapi.id}"
	parent_id   = "${aws_api_gateway_rest_api.restapi.root_resource_id}"
	path_part   = "{proxy+}"
}

resource "aws_api_gateway_method" "proxymeth" {
	rest_api_id   = "${aws_api_gateway_rest_api.restapi.id}"
	resource_id   = "${aws_api_gateway_resource.proxyres.id}"
	http_method   = "ANY"
	authorization = "NONE"
}

resource "aws_api_gateway_method" "proxymeth_root" {
	rest_api_id   = "${aws_api_gateway_rest_api.restapi.id}"
	resource_id   = "${aws_api_gateway_rest_api.restapi.root_resource_id}"
	http_method   = "ANY"
	authorization = "NONE"
}

resource "aws_api_gateway_integration" "lambda_root" {
	rest_api_id = "${aws_api_gateway_rest_api.restapi.id}"
	resource_id = "${aws_api_gateway_method.proxymeth_root.resource_id}"
	http_method = "${aws_api_gateway_method.proxymeth_root.http_method}"

	integration_http_method = "POST"
	type                    = "AWS_PROXY"
	uri                     = "${data.aws_lambda_function.fn.invoke_arn}"
}

resource "aws_api_gateway_integration" "lambda" {
	rest_api_id = "${aws_api_gateway_rest_api.restapi.id}"
	resource_id = "${aws_api_gateway_method.proxymeth.resource_id}"
	http_method = "${aws_api_gateway_method.proxymeth.http_method}"

	integration_http_method = "POST"
	type                    = "AWS_PROXY"
	uri                     = "${data.aws_lambda_function.fn.invoke_arn}"
}

resource "aws_api_gateway_deployment" "deployment" {
	depends_on = [
		"aws_api_gateway_integration.lambda",
		"aws_api_gateway_integration.lambda_root",
	]

	rest_api_id = "${aws_api_gateway_rest_api.restapi.id}"
	stage_name  = "prod"
}

resource "aws_lambda_permission" "apigw" {
	statement_id  = "AllowAPIGatewayInvoke"
	action        = "lambda:InvokeFunction"
	function_name = "${data.aws_lambda_function.fn.arn}"
	principal     = "apigateway.amazonaws.com"

	# The /*/* portion grants access from any method on any resource
	# within the API Gateway "REST API".
	source_arn = "${aws_api_gateway_deployment.deployment.execution_arn}/*/*"
}

output "invoke_url" {
	value = "${aws_api_gateway_deployment.deployment.invoke_url}"
}
