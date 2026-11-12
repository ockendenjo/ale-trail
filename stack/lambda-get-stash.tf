module "lambda_get_stash" {
  source                   = "github.com/ockendenjo/tfmods//lambda"
  project_name             = "ale-trail"
  aws_env                  = var.env
  name                     = "get-stash"
  permissions_boundary_arn = var.permissions_boundary_arn
  s3_bucket                = var.binary_bucket
  s3_object_key            = local.manifest["get-stash"]

  environment = {
    DATA_BUCKET = aws_s3_bucket.data.id
  }
}

resource "aws_api_gateway_method" "get_stash" {
  rest_api_id   = aws_api_gateway_rest_api.main.id
  resource_id   = aws_api_gateway_resource.date.id
  http_method   = "GET"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "get_stash" {
  rest_api_id             = aws_api_gateway_rest_api.main.id
  resource_id             = aws_api_gateway_method.get_stash.resource_id
  http_method             = aws_api_gateway_method.get_stash.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = "arn:aws:apigateway:${data.aws_region.current.region}:lambda:path/2015-03-31/functions/${module.lambda_get_stash.arn}/invocations"
}

resource "aws_lambda_permission" "get_stash" {
  statement_id  = "AllowAPIGInvoke"
  action        = "lambda:InvokeFunction"
  function_name = module.lambda_get_stash.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.main.execution_arn}/*/*"
}

module "iam_s3_get_stash" {
  source     = "github.com/ockendenjo/tfmods//iam-s3"
  bucket_arn = aws_s3_bucket.data.arn
  role_id    = module.lambda_get_stash.role_id
}
