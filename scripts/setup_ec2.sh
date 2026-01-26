#!/bin/bash
# scripts/setup_ec2.sh

# Update system
sudo apt-get update
sudo apt-get upgrade -y

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Install Docker Compose
sudo apt-get install -y docker-compose-plugin

# Setup permissions (so you don't always need sudo)
sudo usermod -aG docker ubuntu

echo "Docker installed successfully."
echo "Please log out and log back in for group permissions to take effect."
