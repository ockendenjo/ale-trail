module "lambda_put_location" {
  source                   = "github.com/ockendenjo/tfmods//lambda"
  project_name             = "ale-trail"
  aws_env                  = var.env
  name                     = "put-location"
  permissions_boundary_arn = var.permissions_boundary_arn
  s3_bucket                = var.binary_bucket
  s3_object_key            = local.manifest["put-location"]

  environment = {
    DATA_BUCKET = aws_s3_bucket.data.id
  }
}

module "apig_put_location" {
  source = "github.com/ockendenjo/tfmods//apig-endpoint"

  rest_api    = aws_api_gateway_rest_api.main
  http_method = "GET"
  path        = "location"
  lambda      = module.lambda_put_location
  resource_id = aws_api_gateway_resource.location.id
}

module "iam_s3_put_location" {
  source      = "github.com/ockendenjo/tfmods//iam-s3"
  bucket_arn  = aws_s3_bucket.data.arn
  role_id     = module.lambda_put_location.role_id
  allow_write = true
}
