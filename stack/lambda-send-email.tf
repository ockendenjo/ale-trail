module "lambda_send_email" {
  source                   = "github.com/ockendenjo/tfmods//lambda"
  project_name             = "ale-trail"
  aws_env                  = var.env
  name                     = "send-email"
  permissions_boundary_arn = var.permissions_boundary_arn
  s3_bucket                = var.binary_bucket
  s3_object_key            = local.manifest["send-email"]
  alarm_topic_arn          = aws_sns_topic.alarms_topic.arn

  environment = {
    DATA_BUCKET     = aws_s3_bucket.data.id
    ALARM_TOPIC_ARN = aws_sns_topic.alarms_topic.arn
    DOMAIN          = var.cloudfront.domain
    REPLY_TO        = var.reply_to_addr
  }
}

module "iam_s3_send_email" {
  source     = "github.com/ockendenjo/tfmods//iam-s3"
  bucket_arn = aws_s3_bucket.data.arn
  role_id    = module.lambda_send_email.role_id
}

module "iam_sns_send_email" {
  source = "github.com/ockendenjo/tfmods//iam-sns"

  role_id = module.lambda_send_email.role_id
  sns_arns = [
    aws_sns_topic.alarms_topic.arn,
  ]
}

data "aws_iam_policy_document" "iam_ses_send_email" {
  statement {
    sid       = "AllowSES"
    effect    = "Allow"
    actions   = ["ses:SendEmail"]
    resources = ["*"]
  }
}

resource "aws_iam_role_policy" "iam_ses_send_email" {
  name_prefix = "ses-"
  role        = module.lambda_send_email.role_id
  policy      = data.aws_iam_policy_document.iam_ses_send_email.json
}

module "schedule_d_minus_2" {
  source              = "./schedule"
  lambda              = module.lambda_send_email
  name                = "ale-trail-D-minus-2-${var.env}"
  schedule_expression = "cron(30 19 * * ? *)"
  template_id         = "DAY-2"
}

module "schedule_d_minus_1" {
  source              = "./schedule"
  lambda              = module.lambda_send_email
  name                = "ale-trail-D-minus-1-${var.env}"
  schedule_expression = "cron(0 8 * * ? *)"
  template_id         = "DAY-1"
}

module "schedule_overdue" {
  source              = "./schedule"
  lambda              = module.lambda_send_email
  name                = "ale-trail-overdue-${var.env}"
  schedule_expression = "cron(0 21 * * ? *)"
  template_id         = "OVERDUE"
}

module "schedule_image_check" {
  source              = "./schedule"
  lambda              = module.lambda_send_email
  name                = "ale-trail-image-check-${var.env}"
  schedule_expression = "cron(0 12 * * ? *)"
  template_id         = "IMAGE-CHECK"
}
