resource "aws_ses_domain_identity" "main" {
  domain = var.cloudfront.domain
}

resource "aws_route53_record" "ses_verification" {
  zone_id = data.aws_route53_zone.primary.zone_id
  name    = "_amazonses.${aws_ses_domain_identity.main.domain}"
  type    = "TXT"
  ttl     = "600"
  records = [aws_ses_domain_identity.main.verification_token]
}

resource "aws_ses_domain_dkim" "main" {
  domain = aws_ses_domain_identity.main.domain
}

resource "aws_route53_record" "ses_dkim" {
  count   = 3
  zone_id = data.aws_route53_zone.primary.zone_id
  name    = "${aws_ses_domain_dkim.main.dkim_tokens[count.index]}._domainkey"
  type    = "CNAME"
  ttl     = "600"
  records = ["${aws_ses_domain_dkim.main.dkim_tokens[count.index]}.dkim.amazonses.com"]
}

locals {
  sestopic = ["Bounce", "Complaint", "Delivery"]
}

resource "aws_sns_topic" "ses" {
  for_each = toset(local.sestopic)

  name = "ale-trail-ses-${lower(each.key)}-${var.env}"
}

resource "aws_ses_identity_notification_topic" "ses" {
  for_each = toset(local.sestopic)

  topic_arn                = aws_sns_topic.ses[each.key].arn
  notification_type        = each.key
  identity                 = aws_ses_domain_identity.main.domain
  include_original_headers = true
}

data "aws_iam_policy_document" "ses" {
  for_each = toset(local.sestopic)
  statement {
    sid     = "AllowSESNotifications"
    effect  = "Allow"
    actions = ["sns:Publish"]

    principals {
      type        = "Service"
      identifiers = ["ses.amazonaws.com"]
    }

    resources = [aws_sns_topic.ses[each.key].arn]

    condition {
      test     = "StringEquals"
      variable = "AWS:SourceAccount"
      values   = [var.aws_account_id]
    }
  }
}

resource "aws_sns_topic_policy" "ses" {
  for_each = toset(local.sestopic)
  arn      = aws_sns_topic.ses[each.key].arn
  policy   = data.aws_iam_policy_document.ses[each.value].json
}

resource "aws_sesv2_email_identity_mail_from_attributes" "main" {
  email_identity   = aws_ses_domain_identity.main.domain
  mail_from_domain = aws_ses_domain_identity.main.domain

  # Safer default: if the bounce MX record is missing, SES falls back to
  # its own MAIL FROM domain rather than rejecting the message.
  behavior_on_mx_failure = "USE_DEFAULT_VALUE"
}

resource "aws_route53_record" "mail_from_mx" {
  zone_id = data.aws_route53_zone.primary.zone_id
  name    = aws_ses_domain_identity.main.domain
  type    = "MX"
  ttl     = 300

  records = [
    "10 feedback-smtp.${var.aws_region}.amazonses.com",
  ]
}

resource "aws_route53_record" "mail_from_spf" {
  zone_id = data.aws_route53_zone.primary.zone_id
  name    = aws_ses_domain_identity.main.domain
  type    = "TXT"
  ttl     = 300

  records = [
    "\"v=spf1 include:amazonses.com ~all\"",
  ]
}
