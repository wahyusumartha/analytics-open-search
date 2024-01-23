resource "aws_sqs_queue" "analytics_incoming_event_queue" {
  name = "analytics-incoming-event"
  message_retention_seconds = 604800
  tags = {
    app = "lambda-open-search-analytics"
    environment = "development"
  }
}

data "aws_iam_policy_document" "allow_lambda_assume_role_to_analytics_queue" {
  statement {
    sid = "AnalyticsQueue_OnlyLambdaAccess"
    effect = "Allow"

    principals {
      identifiers = [aws_iam_role.lambda.arn]
      type        = "AWS"
    }

    actions = [
      "sqs:SendMessage",
      "sqs:ReceiveMessage",
      "sqs:DeleteMessage",
      "sqs:GetQueueUrl",
    ]
    resources = [aws_sqs_queue.analytics_incoming_event_queue.arn]
  }
}

resource "aws_sqs_queue_policy" "analytics_incoming_event_queue_policy" {
  policy    = data.aws_iam_policy_document.allow_lambda_assume_role_to_analytics_queue.json
  queue_url = aws_sqs_queue.analytics_incoming_event_queue.id
}
