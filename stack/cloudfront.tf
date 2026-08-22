resource "aws_cloudfront_origin_access_control" "s3" {
  name                              = "ale-trail-${var.env}"
  signing_behavior                  = "always"
  origin_access_control_origin_type = "s3"
  signing_protocol                  = "sigv4"
}

resource "aws_cloudfront_response_headers_policy" "robots" {
  name = "ale-trail-robots-${var.env}"

  custom_headers_config {
    items {
      header   = "X-Robots-Tag"
      value    = "noindex, nofollow, noarchive, nosnippet"
      override = true
    }
  }
}

resource "aws_cloudfront_response_headers_policy" "index_html" {
  name = "ale-trail-index-html-${var.env}"

  custom_headers_config {
    items {
      header   = "X-Robots-Tag"
      value    = "noindex, nofollow, noarchive, nosnippet"
      override = true
    }
    items {
      header   = "Cache-Control"
      value    = "no-cache, max-age=0, must-revalidate"
      override = true
    }
  }
}

resource "aws_cloudfront_distribution" "ale_trail" {
  comment             = "ale trail (${var.env})"
  enabled             = true
  default_root_object = "index.html"
  price_class         = "PriceClass_100"
  http_version        = "http2and3"
  aliases             = [var.cloudfront.domain]

  origin {
    domain_name              = aws_s3_bucket.static_web.bucket_regional_domain_name
    origin_id                = "s3-static"
    origin_access_control_id = aws_cloudfront_origin_access_control.s3.id
  }

  default_cache_behavior {
    target_origin_id       = "s3-static"
    viewer_protocol_policy = "redirect-to-https"
    compress               = true

    cache_policy_id            = data.aws_cloudfront_cache_policy.caching_optimized.id
    response_headers_policy_id = aws_cloudfront_response_headers_policy.robots.id

    allowed_methods = ["GET", "HEAD"]
    cached_methods  = ["GET", "HEAD"]
  }

  ordered_cache_behavior {
    path_pattern               = "index.html"
    target_origin_id           = "s3-static"
    viewer_protocol_policy     = "redirect-to-https"
    compress                   = true
    cache_policy_id            = data.aws_cloudfront_cache_policy.caching_disabled.id
    response_headers_policy_id = aws_cloudfront_response_headers_policy.index_html.id
    allowed_methods            = ["GET", "HEAD"]
    cached_methods             = ["GET", "HEAD"]
  }

  custom_error_response {
    error_code         = 403
    response_code      = 200
    response_page_path = "/index.html"
  }

  custom_error_response {
    error_code         = 404
    response_code      = 200
    response_page_path = "/index.html"
  }

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  viewer_certificate {
    acm_certificate_arn      = var.cloudfront.certificate_arn
    ssl_support_method       = "sni-only"
    minimum_protocol_version = "TLSv1.2_2021"
  }
}

module "s3_policy_cloudfront" {
  source = "github.com/ockendenjo/tfmods//s3-policy-cloudfront"
  bucket = aws_s3_bucket.static_web
  cloudfront_arns = [
    aws_cloudfront_distribution.ale_trail.arn,
  ]
}

resource "aws_route53_record" "ale_trail" {
  zone_id = data.aws_route53_zone.primary.zone_id
  name    = var.cloudfront.domain
  type    = "A"

  alias {
    name                   = aws_cloudfront_distribution.ale_trail.domain_name
    zone_id                = aws_cloudfront_distribution.ale_trail.hosted_zone_id
    evaluate_target_health = false
  }
}

data "aws_cloudfront_cache_policy" "caching_optimized" {
  name = "Managed-CachingOptimized"
}

data "aws_cloudfront_cache_policy" "caching_disabled" {
  name = "Managed-CachingDisabled"
}

data "aws_route53_zone" "primary" {
  name         = var.zone_domain
  private_zone = false
}
