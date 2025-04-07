#!/bin/bash

echo "Running pre-commit hooks..."
terraform fmt -recursive
terraform validate

