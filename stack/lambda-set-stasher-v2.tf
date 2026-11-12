module "lambda_set_stasher_v2" {
  source                   = "github.com/ockendenjo/tfmods//lambda"
  project_name             = "ale-trail"
  aws_env                  = var.env
  name                     = "set-stasher-v2"
  permissions_boundary_arn = var.permissions_boundary_arn
  s3_bucket                = var.binary_bucket
  s3_object_key            = local.manifest["set-stasher"]
  alarm_topic_arn          = aws_sns_topic.alarms_topic.arn

  environment = {
    DATA_BUCKET = aws_s3_bucket.data.id
  }
}

resource "aws_api_gateway_method" "set_stasher_v2" {
  rest_api_id   = aws_api_gateway_rest_api.main.id
  resource_id   = aws_api_gateway_resource.stasher_v2.id
  http_method   = "POST"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "set_stasher_v2" {
  rest_api_id             = aws_api_gateway_rest_api.main.id
  resource_id             = aws_api_gateway_method.set_stasher_v2.resource_id
  http_method             = aws_api_gateway_method.set_stasher_v2.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = "arn:aws:apigateway:${data.aws_region.current.region}:lambda:path/2015-03-31/functions/${module.lambda_set_stasher_v2.arn}/invocations"
}

resource "aws_lambda_permission" "set_stasher_v2" {
  statement_id  = "AllowAPIGInvoke"
  action        = "lambda:InvokeFunction"
  function_name = module.lambda_set_stasher_v2.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.main.execution_arn}/*/*"
}

module "iam_s3_set_stasher_v2" {
  source      = "github.com/ockendenjo/tfmods//iam-s3"
  bucket_arn  = aws_s3_bucket.data.arn
  role_id     = module.lambda_set_stasher_v2.role_id
  allow_write = true
}
