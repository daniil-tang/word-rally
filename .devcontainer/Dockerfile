# Use an official Go image as the base image
FROM mcr.microsoft.com/vscode/devcontainers/go:latest

# Install Node.js (for React) and other dependencies
RUN apt-get update && apt-get install -y \
  nodejs \
  npm \
  curl \
  git \
  && rm -rf /var/lib/apt/lists/*

# Install AWS CLI if needed
RUN curl "https://d1vvhvl2y92vvt.cloudfront.net/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip" \
  && unzip awscliv2.zip \
  && sudo ./aws/install \
  && rm -rf awscliv2.zip

# Set up the Go workspace environment variables
ENV GOPATH=/go
ENV GOROOT=/usr/local/go
ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin

# Set the default work directory
WORKDIR /workspace