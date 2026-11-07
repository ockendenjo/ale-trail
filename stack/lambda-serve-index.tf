module "lambda_serve_index" {
  source                   = "github.com/ockendenjo/tfmods//lambda"
  project_name             = "ale-trail"
  aws_env                  = var.env
  name                     = "serve-index"
  permissions_boundary_arn = var.permissions_boundary_arn
  s3_bucket                = var.binary_bucket
  s3_object_key            = local.manifest["serve-index"]

  environment = {
    WEB_BUCKET = aws_s3_bucket.static_web.id
  }
}

module "iam_s3_serve_index" {
  source     = "github.com/ockendenjo/tfmods//iam-s3"
  role_id    = module.lambda_serve_index.role_id
  bucket_arn = aws_s3_bucket.static_web.arn
}

resource "aws_api_gateway_method" "serve_index" {
  rest_api_id   = aws_api_gateway_rest_api.main.id
  resource_id   = aws_api_gateway_resource.index.id
  http_method   = "GET"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "serve_index" {
  rest_api_id             = aws_api_gateway_rest_api.main.id
  resource_id             = aws_api_gateway_method.serve_index.resource_id
  http_method             = aws_api_gateway_method.serve_index.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = "arn:aws:apigateway:${data.aws_region.current.region}:lambda:path/2015-03-31/functions/${module.lambda_serve_index.arn}/invocations"
}

resource "aws_lambda_permission" "serve_index" {
  statement_id  = "AllowAPIGInvoke"
  action        = "lambda:InvokeFunction"
  function_name = module.lambda_serve_index.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.main.execution_arn}/*/*"
}