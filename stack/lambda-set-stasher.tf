module "lambda_set_stasher" {
  source                   = "github.com/ockendenjo/tfmods//lambda"
  project_name             = "ale-trail"
  aws_env                  = var.env
  name                     = "set-stasher"
  permissions_boundary_arn = var.permissions_boundary_arn
  s3_bucket                = var.binary_bucket
  s3_object_key            = local.manifest["set-stasher"]

  environment = {
    DATA_BUCKET = aws_s3_bucket.data.id
  }
}

resource "aws_api_gateway_resource" "stashes" {
  parent_id   = aws_api_gateway_resource.api.id
  path_part   = "stashes"
  rest_api_id = aws_api_gateway_rest_api.main.id
}

resource "aws_api_gateway_resource" "date" {
  parent_id   = aws_api_gateway_resource.stashes.id
  path_part   = "{date}"
  rest_api_id = aws_api_gateway_rest_api.main.id
}

module "apig_set_stasher" {
  source = "github.com/ockendenjo/tfmods//apig-endpoint"

  rest_api    = aws_api_gateway_rest_api.main
  http_method = "POST"
  path        = "stasher"
  lambda      = module.lambda_set_stasher
  parent_id   = aws_api_gateway_resource.date.id
}

module "iam_s3_set_stasher" {
  source      = "github.com/ockendenjo/tfmods//iam-s3"
  bucket_arn  = aws_s3_bucket.data.arn
  role_id     = module.lambda_set_stasher.role_id
  allow_write = true
}
