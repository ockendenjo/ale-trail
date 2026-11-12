resource "aws_sns_topic" "alarms_topic" {
  name         = "ale-trail-alarm-topic-${var.env}"
  display_name = "Ale Trail Alarms (${var.env})"
}

resource "aws_sns_topic_policy" "alarms_topic" {
  arn    = aws_sns_topic.alarms_topic.arn
  policy = data.aws_iam_policy_document.alarms_topic.json
}

data "aws_iam_policy_document" "alarms_topic" {
  statement {
    sid    = "AllowCloudWatch"
    effect = "Allow"
    principals {
      type        = "Service"
      identifiers = ["cloudwatch.amazonaws.com"]
    }
    resources = [
      aws_sns_topic.alarms_topic.arn,
    ]
    condition {
      test     = "ArnLike"
      values   = ["arn:aws:cloudwatch:${data.aws_region.current.region}:${var.aws_account_id}:alarm:*"]
      variable = "aws:SourceArn"
    }
  }
}
