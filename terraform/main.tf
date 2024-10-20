provider "aws" {
  region = "us-east-1" # Adjust this to your preferred region
} 

# ECR for frontend
resource "aws_ecr_repository" "frontend" {
  name                 = "frontend"
  image_tag_mutability = "MUTABLE"
}

# ECR for backend
resource "aws_ecr_repository" "backend" {
  name                 = "backend"
  image_tag_mutability = "MUTABLE"
}

# Build and push frontend Docker image
resource "null_resource" "frontend_image" {
  provisioner "local-exec" {
    command = <<EOT
      aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin ${aws_ecr_repository.frontend.repository_url}
      docker build -t ${aws_ecr_repository.frontend.repository_url}:latest ../frontend
      docker push ${aws_ecr_repository.frontend.repository_url}:latest
    EOT
  }
}

# Build and push backend Docker image
resource "null_resource" "backend_image" {
  provisioner "local-exec" {
    command = <<EOT
      aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin ${aws_ecr_repository.backend.repository_url}
      docker build -t ${aws_ecr_repository.backend.repository_url}:latest ../backend
      docker push ${aws_ecr_repository.backend.repository_url}:latest
    EOT
  }
}
