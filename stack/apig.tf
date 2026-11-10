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

      module.apig_set_stasher_v2,
      module.apig_delete_stasher,

      aws_api_gateway_resource.location,
      aws_api_gateway_method.get_location,
      aws_api_gateway_method.put_location,
      aws_api_gateway_integration.get_location,
      aws_api_gateway_integration.put_location,
      aws_api_gateway_method.get_stash,
      aws_api_gateway_integration.get_stash,
      aws_api_gateway_method.list_stashes,
      aws_api_gateway_integration.list_stashes,

      // index.html
      aws_api_gateway_resource.index,
      aws_api_gateway_method.serve_index,
      aws_api_gateway_integration.serve_index,
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

resource "aws_api_gateway_resource" "stashes" {
  parent_id   = aws_api_gateway_resource.api.id
  path_part   = "stashes"
  rest_api_id = aws_api_gateway_rest_api.main.id
}

resource "aws_api_gateway_resource" "index" {
  parent_id   = aws_api_gateway_resource.api.id
  path_part   = "index.html"
  rest_api_id = aws_api_gateway_rest_api.main.id
}

resource "aws_api_gateway_resource" "date" {
  parent_id   = aws_api_gateway_resource.stashes.id
  path_part   = "{date}"
  rest_api_id = aws_api_gateway_rest_api.main.id
}

resource "aws_api_gateway_resource" "location" {
  parent_id   = aws_api_gateway_resource.date.id
  path_part   = "location"
  rest_api_id = aws_api_gateway_rest_api.main.id
}

resource "aws_api_gateway_resource" "stasher_v3" {
  parent_id   = aws_api_gateway_resource.date.id
  path_part   = "stasherV3"
  rest_api_id = aws_api_gateway_rest_api.main.id
}
