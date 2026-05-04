## GitHub Actions OIDC Authentication for ECR

To publish Docker images securely without long-lived static AWS access keys, this project uses GitHub OpenID Connect (OIDC).

A dedicated IAM Role `GitHubActionsShareBiteECRPublisherRole` is provisioned in AWS.
*(Note: In our environment, this role must be created with the `GolangBound` permissions boundary).*

### 1. Trust Policy (OIDC Condition)
This policy allows GitHub Actions to assume the role, restricted strictly to the `share-bite` repository and `main` branch or release tags (`v*`).
```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Federated": "arn:aws:iam::<AWS_ACCOUNT_ID>:oidc-provider/token.actions.githubusercontent.com"
      },
      "Action": "sts:AssumeRoleWithWebIdentity",
      "Condition": {
        "StringEquals": {
          "token.actions.githubusercontent.com:aud": "sts.amazonaws.com"
        },
        "StringLike": {
          "token.actions.githubusercontent.com:sub": [
            "repo:ua-academy-projects/share-bite:ref:refs/heads/main",
            "repo:ua-academy-projects/share-bite:ref:refs/tags/v*"
          ]
        }
      }
    }
  ]
}
```

### 2. Permission Policy (Least Privilege)
This policy grants the role permission to push and pull images ONLY to/from the `share-bite/*` ECR repositories.
```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "GetAuthorizationToken",
      "Effect": "Allow",
      "Action": "ecr:GetAuthorizationToken",
      "Resource": "*"
    },
    {
      "Sid": "AllowPushPullToShareBiteRepos",
      "Effect": "Allow",
      "Action": [
        "ecr:BatchCheckLayerAvailability",
        "ecr:GetDownloadUrlForLayer",
        "ecr:GetRepositoryPolicy",
        "ecr:DescribeRepositories",
        "ecr:ListImages",
        "ecr:DescribeImages",
        "ecr:BatchGetImage",
        "ecr:InitiateLayerUpload",
        "ecr:UploadLayerPart",
        "ecr:CompleteLayerUpload",
        "ecr:PutImage"
      ],
      "Resource": "arn:aws:ecr:<REGION>:<AWS_ACCOUNT_ID>:repository/share-bite/*"
    }
  ]
}
```