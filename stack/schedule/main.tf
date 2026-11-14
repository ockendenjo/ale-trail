variable "name" {
  type        = string
  description = "Schedule name"
}

variable "schedule_expression" {
  type = string
}

variable "template_id" {
  type = string
}

variable "lambda" {
  type = object({
    function_name = string
    arn           = string
  })
}

resource "aws_cloudwatch_event_rule" "main" {
  name                = var.name
  schedule_expression = var.schedule_expression
}

resource "aws_cloudwatch_event_target" "main" {
  rule      = aws_cloudwatch_event_rule.main.name
  target_id = "invoke-my-lambda"
  arn       = var.lambda.arn

  # Optional payload passed to the Lambda event.
  input = jsonencode({
    template = var.template_id,
  })
}

resource "aws_lambda_permission" "main" {
  statement_id  = "schedule-${var.name}"
  action        = "lambda:InvokeFunction"
  function_name = var.lambda.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.main.arn
}
