output "resource_id" {
  value = var.resource_id
}

output "method_id" {
  value = aws_api_gateway_method.main.id
}

output "integration_id" {
  value = aws_api_gateway_integration.main.id
}

output "lambda" {
  value = var.lambda
}
