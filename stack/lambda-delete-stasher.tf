module "lambda_delete_stasher" {
  source                   = "github.com/ockendenjo/tfmods//lambda"
  project_name             = "ale-trail"
  aws_env                  = var.env
  name                     = "delete-stasher"
  permissions_boundary_arn = var.permissions_boundary_arn
  s3_bucket                = var.binary_bucket
  s3_object_key            = local.manifest["delete-stasher"]

  environment = {
    DATA_BUCKET = aws_s3_bucket.data.id
  }
}

module "apig_delete_stasher" {
  source = "./apig-endpoint"

  http_method = "DELETE"
  lambda      = module.lambda_delete_stasher
  region      = data.aws_region.current.region
  resource_id = aws_api_gateway_resource.stasher.id
  rest_api    = aws_api_gateway_rest_api.main
}

module "iam_s3_delete_stasher" {
  source      = "github.com/ockendenjo/tfmods//iam-s3"
  bucket_arn  = aws_s3_bucket.data.arn
  role_id     = module.lambda_delete_stasher.role_id
  allow_write = true
}
