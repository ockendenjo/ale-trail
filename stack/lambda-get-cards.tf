module "lambda_get_cards" {
  source                   = "github.com/ockendenjo/tfmods//lambda"
  project_name             = "ale-trail"
  aws_env                  = var.env
  name                     = "get-cards"
  permissions_boundary_arn = var.permissions_boundary_arn
  s3_bucket                = var.binary_bucket
  s3_object_key            = local.manifest["get-cards"]

  environment = {
  }
}

module "apig_get_cards" {
  source = "github.com/ockendenjo/tfmods//apig-endpoint"

  rest_api    = aws_api_gateway_rest_api.main
  http_method = "POST"
  path        = "deleteAccount"
  lambda      = module.lambda_get_cards
  parent_id   = aws_api_gateway_resource.api.id
}
