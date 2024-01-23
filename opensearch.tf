resource "aws_cloudwatch_log_group" "analytics_opensearch_log_group" {
  name = "/aws/opensearch/analytics"
  retention_in_days = 3
}

data "aws_iam_policy_document" "allow_analytics_opensearch_logging" {
  statement {
    effect = "Allow"

    principals {
      identifiers = ["es.amazonaws.com"]
      type        = "Service"
    }

    actions = [
      "logs:PutLogEvents",
      "logs:PutLogEventsBatch",
      "logs:CreateLogStream",
    ]

    resources = ["arn:aws:logs:*"]
  }
}

data "aws_region" "current" {}

resource "aws_cloudwatch_log_resource_policy" "allow_analytics_opensearch_cloudwatch_policy" {
  policy_document = data.aws_iam_policy_document.allow_analytics_opensearch_logging.json
  policy_name     = "AllowAnalyticsOpenSearchCloudWatchLogsPolicy"
}


variable "domain_name" {
  default = "analytics"
}

data "aws_iam_policy_document" "allow_analytics_domain_access" {
  statement {
    effect = "Allow"

    principals {
      identifiers = ["*"]
      type        = "AWS"
    }

    actions = ["es:*"]
    resources = ["arn:aws:es:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:domain/${var.domain_name}/*"]
  }
}

resource "aws_opensearch_domain" "analytics" {
  domain_name = var.domain_name
  engine_version = "OpenSearch_2.11"

  cluster_config {
    instance_type = "t3.small.search"
    dedicated_master_enabled = false
    instance_count = 1
    warm_enabled = false
    zone_awareness_enabled = false
  }

  advanced_security_options {
    enabled = true
    anonymous_auth_enabled = false
    internal_user_database_enabled = true

    master_user_options {
      master_user_name = local.env_variables["opensearch"]["username"]
      master_user_password = local.env_variables["opensearch"]["password"]
    }
  }

  domain_endpoint_options {
    enforce_https = true
    tls_security_policy = "Policy-Min-TLS-1-2-2019-07"
  }

  log_publishing_options {
    cloudwatch_log_group_arn = aws_cloudwatch_log_group.analytics_opensearch_log_group.arn
    log_type                 = "INDEX_SLOW_LOGS"
  }

  encrypt_at_rest {
    enabled = true
  }

  node_to_node_encryption {
    enabled = true
  }

  ebs_options {
    ebs_enabled = true
    volume_size = 10
    volume_type = "gp3"
    throughput = 125
  }

  access_policies = data.aws_iam_policy_document.allow_analytics_domain_access.json

  tags = {
    env = "dev"
  }
}