provider "aws" {
  region = "us-east-2"
}

data "aws_caller_identity" "current" {}

data "aws_ami" "ubuntu" {
  most_recent = true
  owners      = ["099720109477"]

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd-gp3/ubuntu-noble-24.04-amd64-server-*"]
  }
}

resource "aws_ecr_repository" "repo" {
  name                 = "share-bite/notifications-worker"
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }
}

resource "aws_sns_topic" "notifications" {
  name = "notifications"
}

resource "aws_sqs_queue" "dlq_sse" {
  name = "notifications-sse-dlq"
}

resource "aws_sqs_queue" "dlq_lambda" {
  name = "notifications-lambda-dlq"
}

resource "aws_sqs_queue" "notifications_sse" {
  name                       = "notifications-sse"
  visibility_timeout_seconds = 30
  message_retention_seconds  = 345600
  receive_wait_time_seconds  = 1
  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.dlq_sse.arn,
    maxReceiveCount     = 3
  })
}

resource "aws_sqs_queue" "notifications_lambda" {
  name                       = "notifications-lambda"
  visibility_timeout_seconds = 30
  message_retention_seconds  = 345600
  receive_wait_time_seconds  = 1
  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.dlq_lambda.arn,
    maxReceiveCount     = 3
  })
}

resource "aws_sqs_queue_policy" "sse_policy" {
  queue_url = aws_sqs_queue.notifications_sse.id
  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [{
      Sid       = "Allow-SNS-SendMessage",
      Effect    = "Allow",
      Principal = "*",
      Action    = "sqs:SendMessage",
      Resource  = aws_sqs_queue.notifications_sse.arn,
      Condition = { ArnEquals = { "aws:SourceArn" = aws_sns_topic.notifications.arn } }
    }]
  })
}

resource "aws_sqs_queue_policy" "lambda_policy" {
  queue_url = aws_sqs_queue.notifications_lambda.id
  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [{
      Sid       = "Allow-SNS-SendMessage",
      Effect    = "Allow",
      Principal = "*",
      Action    = "sqs:SendMessage",
      Resource  = aws_sqs_queue.notifications_lambda.arn,
      Condition = { ArnEquals = { "aws:SourceArn" = aws_sns_topic.notifications.arn } }
    }]
  })
}

resource "aws_sns_topic_subscription" "to_sse" {
  topic_arn            = aws_sns_topic.notifications.arn
  protocol             = "sqs"
  endpoint             = aws_sqs_queue.notifications_sse.arn
  raw_message_delivery = true
}

resource "aws_sns_topic_subscription" "to_lambda" {
  topic_arn            = aws_sns_topic.notifications.arn
  protocol             = "sqs"
  endpoint             = aws_sqs_queue.notifications_lambda.arn
  raw_message_delivery = true
}

resource "aws_iam_role" "ec2_ecr_role" {
  name                 = "Training-GolangShareBiteEC2Role"
  permissions_boundary = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:policy/GolangBound"

  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Action    = "sts:AssumeRole",
        Effect    = "Allow",
        Principal = { Service = "ec2.amazonaws.com" }
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "ecr_read" {
  role       = aws_iam_role.ec2_ecr_role.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
}

resource "aws_iam_role_policy_attachment" "ssm_core" {
  role       = aws_iam_role.ec2_ecr_role.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"
}

resource "aws_iam_role_policy_attachment" "ec2_sqs_full" {
  role       = aws_iam_role.ec2_ecr_role.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonSQSFullAccess"
}

resource "aws_iam_instance_profile" "ec2_profile" {
  name = "Training-GolangShareBiteEC2Profile"
  role = aws_iam_role.ec2_ecr_role.name
}

resource "aws_security_group" "share_bite_sg" {
  name        = "share-bite-sg"
  description = "Allow API traffic only (SSH is managed via SSM)"

  ingress {
    from_port   = 8080
    to_port     = 8082
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_instance" "app_server" {
  ami                    = data.aws_ami.ubuntu.id
  instance_type          = "t3.small"
  vpc_security_group_ids = [aws_security_group.share_bite_sg.id]
  iam_instance_profile   = aws_iam_instance_profile.ec2_profile.name

  root_block_device {
    encrypted   = true
    volume_type = "gp3"
  }

  metadata_options {
    http_tokens = "required"
  }

  user_data = <<-EOF
              #!/bin/bash
              apt-get update && apt-get upgrade -y
              apt-get install -y awscli curl golang-go

              curl -fsSL https://get.docker.com -o get-docker.sh
              sh get-docker.sh
              usermod -aG docker ubuntu

              echo "export PATH=\$PATH:\$HOME/go/bin" >> /home/ubuntu/.bashrc

              mkdir -p /home/ubuntu/share-bite
              chown -R ubuntu:ubuntu /home/ubuntu/share-bite
              EOF

  tags = {
    Name = "ShareBite-App-Server"
  }
}

resource "aws_iam_role" "training_lambda_role" {
  name                 = "Training-GolangShareBiteLambdaRole"
  permissions_boundary = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:policy/GolangBound"

  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Action    = "sts:AssumeRole",
        Effect    = "Allow",
        Principal = { Service = "lambda.amazonaws.com" }
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "lambda_basic_exec" {
  role       = aws_iam_role.training_lambda_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_role_policy_attachment" "lambda_sqs_full" {
  role       = aws_iam_role.training_lambda_role.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonSQSFullAccess"
}

resource "aws_iam_role_policy_attachment" "lambda_vpc_access" {
  count      = 0
  role       = aws_iam_role.training_lambda_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

resource "aws_iam_policy" "ecr_pull" {
  name        = "notifications-worker-ecr-pull"
  description = "Allow Lambda to pull images from ECR repository"

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Action = [
          "ecr:BatchGetImage",
          "ecr:GetDownloadUrlForLayer",
          "ecr:DescribeImages"
        ],
        Resource = aws_ecr_repository.repo.arn
      },
      {
        Effect   = "Allow",
        Action   = ["ecr:GetAuthorizationToken"],
        Resource = "*"
      }
    ]
  })
}

resource "aws_iam_policy_attachment" "attach_ecr_pull" {
  name       = "notifications-worker-attach-ecr"
  policy_arn = aws_iam_policy.ecr_pull.arn
  roles      = [aws_iam_role.training_lambda_role.name]
}

resource "aws_lambda_function" "worker" {
  function_name = "notifications-worker"
  package_type  = "Image"
  image_uri     = "${aws_ecr_repository.repo.repository_url}:20ab157"
  role          = aws_iam_role.training_lambda_role.arn
  memory_size   = 128
  timeout       = 3
  architectures = ["arm64"]
}

resource "aws_lambda_event_source_mapping" "sqs_to_lambda" {
  event_source_arn = aws_sqs_queue.notifications_lambda.arn
  function_name    = aws_lambda_function.worker.arn
  batch_size       = 10
  enabled          = true
}

resource "aws_instance" "notifications_sse" {
  ami                    = data.aws_ami.ubuntu.id
  instance_type          = "t3.small"
  vpc_security_group_ids = [aws_security_group.share_bite_sg.id]
  iam_instance_profile   = aws_iam_instance_profile.ec2_profile.name

  root_block_device {
    encrypted   = true
    volume_type = "gp3"
  }

  metadata_options {
    http_tokens = "required"
  }

  user_data = <<-EOF
              #!/bin/bash
              apt-get update && apt-get upgrade -y
              apt-get install -y awscli curl golang-go
              curl -fsSL https://get.docker.com -o get-docker.sh
              sh get-docker.sh
              usermod -aG docker ubuntu
              echo "export PATH=\$PATH:\$HOME/go/bin" >> /home/ubuntu/.bashrc
              mkdir -p /home/ubuntu/share-bite
              chown -R ubuntu:ubuntu /home/ubuntu/share-bite
              EOF

  tags = {
    Name = "ShareBite-Notifications-SSE"
  }
}

output "ecr_repository_url" { value = aws_ecr_repository.repo.repository_url }
output "sns_topic_arn" { value = aws_sns_topic.notifications.arn }
output "notifications_sse_queue_arn" { value = aws_sqs_queue.notifications_sse.arn }
output "notifications_lambda_queue_arn" { value = aws_sqs_queue.notifications_lambda.arn }
output "notifications_sse_dlq_arn" { value = aws_sqs_queue.dlq_sse.arn }
output "notifications_lambda_dlq_arn" { value = aws_sqs_queue.dlq_lambda.arn }
output "lambda_arn" { value = aws_lambda_function.worker.arn }
output "app_instance_public_ip" { value = aws_instance.app_server.public_ip }
output "sse_instance_public_ip" { value = aws_instance.notifications_sse.public_ip }
