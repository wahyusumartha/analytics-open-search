data "aws_caller_identity" "current" {}

// allow lambda service to use assume role with such policy
data "aws_iam_policy_document" "assume_lambda_role" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      identifiers = ["lambda.amazonaws.com"]
      type        = "Service"
    }

    principals {
      identifiers = [data.aws_caller_identity.current.arn]
      type        = "AWS"
    }
  }
}

resource "aws_iam_role" "lambda" {
  name = "AssumeLambdaRole"
  description = "Role for lambda to assume lambda"
  assume_role_policy = data.aws_iam_policy_document.assume_lambda_role.json
}

data "aws_iam_policy_document" "allow_lambda_logging" {
  statement {
    effect = "Allow"
    actions = [
      "logs:CreateLogStream",
      "logs:PutLogEvents",
    ]

    resources = [
      "arn:aws:logs:*:*:*",
    ]
  }
}

resource "aws_iam_policy" "function_logging_policy" {
  name = "AllowLambdaLoggingPolicy"
  description = "Policy for lambda cloudwatch logging"
  policy = data.aws_iam_policy_document.allow_lambda_logging.json
}

resource "aws_iam_role_policy_attachment" "lambda_logging_policy_attachment" {
  policy_arn = aws_iam_policy.function_logging_policy.arn
  role       = aws_iam_role.lambda.id
}

data "aws_iam_policy_document" "allow_lambda_sqs" {
  statement {
    effect = "Allow"
    actions = [
      "sqs:SendMessage",
      "sqs:ReceiveMessage",
      "sqs:DeleteMessage",
      "sqs:GetQueueUrl",
      "sqs:GetQueueAttributes",
    ]

    resources = [
      aws_sqs_queue.analytics_incoming_event_queue.arn
    ]
  }
}

resource "aws_iam_policy" "function_sqs_policy" {
  name = "AllowAnalyticsLambdaSQSPolicy"
  description = "Policy for lambda sqs"
  policy = data.aws_iam_policy_document.allow_lambda_sqs.json
  tags = {
    app = "lambda-open-search-analytics"
  }
}

resource "aws_iam_role_policy_attachment" "lambda_sqs_policy_attachment" {
  policy_arn = aws_iam_policy.function_sqs_policy.arn
  role       = aws_iam_role.lambda.id
}

resource "aws_lambda_permission" "receive_event_api_gateway" {
  statement_id = "AllowInvokeReceiveEvent"
  action = "lambda:InvokeFunction"
  function_name = aws_lambda_function.receive_event_function.function_name
  principal     = "apigateway.amazonaws.com"

  # The /* part allows invocation from any stage, method and resource path
  # within API Gateway.
  source_arn = "${aws_api_gateway_rest_api.analytics_api.execution_arn}/*/*/*"
}

data "aws_iam_policy_document" "api_gw_assume_role" {
  statement {
    effect = "Allow"

    principals {
      identifiers = ["apigateway.amazonaws.com"]
      type        = "Service"
    }

    actions = ["sts:AssumeRole"]
  }
}

resource "aws_iam_role" "api_gw_assume_role" {
  name = "api_gateway_cloudwatch_global"
  assume_role_policy = data.aws_iam_policy_document.api_gw_assume_role.json
}

data "aws_iam_policy_document" "allow_api_gw_cloudwatch_logging" {
  statement {
    effect = "Allow"

    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:DescribeLogGroups",
      "logs:DescribeLogStreams",
      "logs:PutLogEvents",
      "logs:GetLogEvents",
      "logs:FilterLogEvents",
    ]

    resources = ["*"]
  }
}

resource "aws_iam_policy" "api_gw_cloudwatch_policy" {
  name = "AllowAPIGatewayCloudwatchLoggingPolicy"
  description = "Policy for API Gateway Cloudwatch Logging"
  policy = data.aws_iam_policy_document.allow_api_gw_cloudwatch_logging.json
}

resource "aws_iam_role_policy_attachment" "api_gw_cloudwatch_logging" {
  policy_arn = aws_iam_policy.api_gw_cloudwatch_policy.arn
  role       = aws_iam_role.api_gw_assume_role.id
}

