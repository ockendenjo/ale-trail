resource "aws_api_gateway_rest_api" "main" {
  name        = "ale-trail-api-${var.env}"
  description = "Ale Trail (${var.env})"

  disable_execute_api_endpoint = false

  endpoint_configuration {
    types = ["REGIONAL"]
  }
}

resource "aws_api_gateway_deployment" "main" {
  rest_api_id = aws_api_gateway_rest_api.main.id

  triggers = {
    redeployment = sha1(jsonencode([
      aws_api_gateway_resource.api,
      aws_api_gateway_resource.stashes,
      aws_api_gateway_resource.date,

      module.apig_get_cards,
      module.apig_set_stasher,
    ]))
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_api_gateway_stage" "stg" {
  deployment_id = aws_api_gateway_deployment.main.id
  rest_api_id   = aws_api_gateway_rest_api.main.id
  stage_name    = "stg"
}

resource "aws_api_gateway_resource" "api" {
  parent_id   = aws_api_gateway_rest_api.main.root_resource_id
  path_part   = "api"
  rest_api_id = aws_api_gateway_rest_api.main.id
}
