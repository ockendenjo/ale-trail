resource "aws_api_gateway_method" "main" {
  rest_api_id   = var.rest_api.id
  resource_id   = var.resource_id
  http_method   = var.http_method
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "main" {
  rest_api_id             = var.rest_api.id
  resource_id             = var.resource_id
  http_method             = var.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = "arn:aws:apigateway:${var.region}:lambda:path/2015-03-31/functions/${var.lambda.arn}/invocations"
}

resource "aws_lambda_permission" "main" {
  count         = var.create_permissions ? 1 : 0
  statement_id  = "AllowAPIGInvokeV3"
  action        = "lambda:InvokeFunction"
  function_name = var.lambda.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${var.rest_api.execution_arn}/*/*"
}
