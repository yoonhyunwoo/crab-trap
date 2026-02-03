#!/bin/bash
set -e

# Configuration
REGION="${AWS_REGION:-ap-northeast-2}"
REPO_NAME="crab-trap"
IMAGE_TAG="${IMAGE_TAG:-latest}"
AWS_PROFILE="${AWS_PROFILE:-yhw}"

# Get account ID
ACCOUNT_ID=$(aws sts get-caller-identity --profile ${AWS_PROFILE} --query Account --output text)

# Construct ECR repository URL
ECR_URL="${ACCOUNT_ID}.dkr.ecr.${REGION}.amazonaws.com/${REPO_NAME}"

echo "Building and pushing Docker image to ECR..."
echo "Repository: ${ECR_URL}"
echo "Tag: ${IMAGE_TAG}"
echo "Profile: ${AWS_PROFILE}"
echo ""

# Login to ECR
echo "Logging in to ECR..."
aws ecr get-login-password --region ${REGION} --profile ${AWS_PROFILE} | docker login --username AWS --password-stdin ${ACCOUNT_ID}.dkr.ecr.${REGION}.amazonaws.com

# Build Docker image
echo "Building Docker image..."
docker build -t ${REPO_NAME}:${IMAGE_TAG} .

# Tag image for ECR
echo "Tagging image..."
docker tag ${REPO_NAME}:${IMAGE_TAG} ${ECR_URL}:${IMAGE_TAG}

# Push to ECR
echo "Pushing to ECR..."
docker push ${ECR_URL}:${IMAGE_TAG}

echo ""
echo "Successfully pushed ${ECR_URL}:${IMAGE_TAG}"
echo ""
echo "To deploy to EC2, run:"
echo "  cd terraform"
echo "  terraform apply -var='image_tag=${IMAGE_TAG}'"
