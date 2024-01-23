data "aws_ssm_parameter" "analytics-service-env" {
  name = "analytics-service"
  with_decryption = false
}

output "analytics_service_env" {
  value = data.aws_ssm_parameter.analytics-service-env.value
  sensitive = true
}

locals {
  env_variables = jsondecode(data.aws_ssm_parameter.analytics-service-env.value)
}
