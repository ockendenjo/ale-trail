resource "aws_ses_domain_identity" "main" {
  count  = var.env == "pro" ? 1 : 0
  domain = var.zone_domain
}

resource "aws_route53_record" "ses_verification" {
  count   = var.env == "pro" ? 1 : 0
  zone_id = data.aws_route53_zone.primary.zone_id
  name    = "_amazonses.${var.zone_domain}"
  type    = "TXT"
  ttl     = "600"
  records = [aws_ses_domain_identity.main[0].verification_token]
}

resource "aws_ses_domain_dkim" "main" {
  count  = var.env == "pro" ? 1 : 0
  domain = aws_ses_domain_identity.main[0].domain
}

resource "aws_route53_record" "ses_dkim" {
  count   = var.env == "pro" ? 3 : 0
  zone_id = data.aws_route53_zone.primary.zone_id
  name    = "${aws_ses_domain_dkim.main[0].dkim_tokens[count.index]}._domainkey"
  type    = "CNAME"
  ttl     = "600"
  records = ["${aws_ses_domain_dkim.main[0].dkim_tokens[count.index]}.dkim.amazonses.com"]
}
