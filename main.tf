terraform {

  cloud {
    organization = "ringkasin"

    workspaces {
      name = "analytics-open-search"
    }
  }

  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
    archive = {
      source = "hashicorp/archive"
    }
  }

  required_version = ">= 1.3.7"
}

provider "aws" {
  region = "ap-southeast-1"

  default_tags {
    tags = {
      app = "lambda-open-search-analytics"
    }
  }
}