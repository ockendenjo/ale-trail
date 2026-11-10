variable "rest_api" {
  type = object({
    id            = string
    execution_arn = string
  })
}

variable "resource_id" {
  type = string
}

variable "http_method" {
  type = string
}

variable "lambda" {
  type = object({
    arn           = string
    function_name = string
  })
}

variable "region" {
  type = string
}

variable "create_permissions" {
  type    = bool
  default = true
}