data "archive_file" "received_event_function_archive" {
  type = "zip"

  source_file = local.receive_event_function.binary_path
  output_path = local.receive_event_function.archive_path
}

resource "aws_lambda_function" "receive_event_function" {
  function_name = local.receive_event_function.name
  description = "Function to receive an event via REST API and Push the event to SQS"
  role          = aws_iam_role.lambda.arn
  handler = local.receive_event_function.binary_name
  memory_size = 128

  filename = local.receive_event_function.archive_path
  source_code_hash = data.archive_file.received_event_function_archive.output_base64sha256

  runtime = "go1.x"

  environment {
    variables = {
      aws_key = local.env_variables["aws_key"]
      aws_secret = local.env_variables["aws_secret"]
      incoming_event_queue_url = local.env_variables["analytics_incoming_event_queue_url"]
      role_arn = local.env_variables["role_arn"]
    }
  }
}

resource "aws_lambda_function_url" "receive_event_function_url" {
  function_name = aws_lambda_function.receive_event_function.arn
  authorization_type = "NONE"
  cors {
    allow_methods = ["POST"]
    allow_origins = ["*"]
  }
}

output "receive_event_function_url" {
  value = aws_lambda_function_url.receive_event_function_url.function_url
}

resource "aws_cloudwatch_log_group" "received_event_log_group" {
  name = "/aws/lambda/${aws_lambda_function.receive_event_function.function_name}"
  retention_in_days = 3
}

data "archive_file" "ingest_event_function_archive" {
  type = "zip"

  source_file = local.ingest_event_function.binary_path
  output_path = local.ingest_event_function.archive_path
}

resource "aws_lambda_function" "ingest_event_function" {
  function_name = local.ingest_event_function.name
  description = "Function to consume sqs and ingest event to open search"
  role          = aws_iam_role.lambda.arn
  handler = local.ingest_event_function.binary_name
  memory_size =  128

  filename = local.ingest_event_function.archive_path
  source_code_hash = data.archive_file.ingest_event_function_archive.output_base64sha256

  runtime = "go1.x"

  environment {
    variables = {
      opensearch_endpoint = local.env_variables["opensearch"]["endpoint"]
      opensearch_username = local.env_variables["opensearch"]["username"]
      opensearch_password = local.env_variables["opensearch"]["password"]
    }
  }
}

resource "aws_lambda_event_source_mapping" "ingest_event_sqs_trigger" {
  event_source_arn = aws_sqs_queue.analytics_incoming_event_queue.arn
  function_name    = aws_lambda_function.ingest_event_function.arn
  depends_on = [
    aws_sqs_queue.analytics_incoming_event_queue
  ]
}

resource "aws_cloudwatch_log_group" "ingest_event_log_group" {
  name = "/aws/lambda/${aws_lambda_function.ingest_event_function.function_name}"
  retention_in_days = 3
}

