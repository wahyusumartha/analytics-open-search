resource "aws_api_gateway_rest_api" "analytics_api" {
  name        = "analytics-service-api-gateway"
  description = "An API Gateway For Analytics Service"

  endpoint_configuration {
    types = ["REGIONAL"]
  }
}

resource "aws_api_gateway_resource" "receive_event" {
  parent_id   = aws_api_gateway_rest_api.analytics_api.root_resource_id
  path_part   = "event"
  rest_api_id = aws_api_gateway_rest_api.analytics_api.id
}

resource "aws_api_gateway_method" "proxy" {
  authorization = "NONE"
  http_method   = "POST"
  resource_id   = aws_api_gateway_resource.receive_event.id
  rest_api_id   = aws_api_gateway_rest_api.analytics_api.id
}

resource "aws_api_gateway_integration" "receive_event_lambda_integration" {
  http_method = aws_api_gateway_method.proxy.http_method
  resource_id = aws_api_gateway_resource.receive_event.id
  rest_api_id = aws_api_gateway_rest_api.analytics_api.id
  integration_http_method = "POST"
  type        = "AWS_PROXY"
  uri = aws_lambda_function.receive_event_function.invoke_arn
  request_templates = {
    "application/json" = jsonencode(
      {
        statusCode = 200
      }
    )
  }
}

resource "aws_api_gateway_method_response" "proxy" {
  http_method = aws_api_gateway_method.proxy.http_method
  resource_id = aws_api_gateway_resource.receive_event.id
  rest_api_id = aws_api_gateway_rest_api.analytics_api.id
  status_code = 200
}

resource "aws_api_gateway_integration_response" "proxy" {
  http_method = aws_api_gateway_method_response.proxy.http_method
  resource_id = aws_api_gateway_resource.receive_event.id
  rest_api_id = aws_api_gateway_rest_api.analytics_api.id
  status_code = aws_api_gateway_method_response.proxy.status_code
  depends_on = [
    aws_api_gateway_method.proxy,
    aws_api_gateway_integration.receive_event_lambda_integration
  ]
}

resource "aws_api_gateway_deployment" "deployment" {
  depends_on = [
    aws_api_gateway_integration.receive_event_lambda_integration
  ]

  rest_api_id = aws_api_gateway_rest_api.analytics_api.id
  description = "the api gateway deployment for dev environment - ${timestamp()}"

  variables = {
    deployed_at = timestamp()
  }
  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_api_gateway_account" "default_account" {
  cloudwatch_role_arn = aws_iam_role.api_gw_assume_role.arn
}

variable "stage_name" {
  default = "dev"
  type = string
}

resource "aws_cloudwatch_log_group" "api_gateway_execution_log" {
  name = "API-Gateway-Execution-Logs_${aws_api_gateway_deployment.deployment.rest_api_id}/${var.stage_name}"
  retention_in_days = 3
}

resource "aws_api_gateway_stage" "global" {
  depends_on = [
    aws_cloudwatch_log_group.api_gateway_execution_log
  ]

  deployment_id = aws_api_gateway_deployment.deployment.id
  rest_api_id   = aws_api_gateway_rest_api.analytics_api.id
  stage_name    = var.stage_name
  description = "Deployed at ${timestamp()}"
  access_log_settings {
    destination_arn = aws_cloudwatch_log_group.api_gateway_execution_log.arn
    format          = jsonencode({
      requestId        = "$context.requestId"
      requestTime      = "$context.requestTime"
      requestTimeEpoch = "$context.requestTimeEpoch"
      path             = "$context.path"
      method           = "$context.httpMethod"
      status           = "$context.status"
      responseLength   = "$context.responseLength"
      message          = "$context.error.message"
    })
  }
}