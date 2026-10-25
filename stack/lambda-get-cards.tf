module "lambda_get_cards" {
  source                   = "github.com/ockendenjo/tfmods//lambda"
  project_name             = "ale-trail"
  aws_env                  = var.env
  name                     = "get-cards"
  permissions_boundary_arn = var.permissions_boundary_arn
  s3_bucket                = var.binary_bucket
  s3_object_key            = local.manifest["get-cards"]

  environment = {
    DATA_BUCKET = aws_s3_bucket.data.id
  }
}

module "apig_get_cards" {
  source = "github.com/ockendenjo/tfmods//apig-endpoint"

  rest_api    = aws_api_gateway_rest_api.main
  http_method = "GET"
  path        = "cards"
  lambda      = module.lambda_get_cards
  parent_id   = aws_api_gateway_resource.api.id
}

module "iam_s3_get_cards" {
  source     = "github.com/ockendenjo/tfmods//iam-s3"
  bucket_arn = aws_s3_bucket.data.arn
  role_id    = module.lambda_get_cards.role_id
}
