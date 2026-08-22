resource "aws_s3_bucket" "static_web" {
  bucket_prefix = "ale-trail-${var.env}-web-"
  force_destroy = true
}
