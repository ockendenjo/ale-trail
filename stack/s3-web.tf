resource "aws_s3_bucket" "static_web" {
  bucket_prefix = "ale-trail-${var.env}-web-"
  force_destroy = true
}

resource "aws_s3_bucket" "data" {
  bucket_prefix = "ale-trail-${var.env}-data-"
  force_destroy = true
}

resource "aws_s3_bucket_versioning" "data" {
  bucket = aws_s3_bucket.data.id

  versioning_configuration {
    status = "Enabled"
  }
}
