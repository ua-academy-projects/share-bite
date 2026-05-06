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

resource "aws_iam_role" "ec2_ecr_role" {
  name                 = "Training-GolangShareBiteEC2Role"
  permissions_boundary = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:policy/GolangBound"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = { Service = "ec2.amazonaws.com" }
    }]
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
  ami           = data.aws_ami.ubuntu.id
  instance_type = "t3.small"

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

output "instance_public_ip" {
  description = "Public IP of the EC2 instance"
  value       = aws_instance.app_server.public_ip
}
