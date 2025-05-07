#!/bin/bash

# Initialization script for LocalStack

# Basic JSON secret (for KEY=URL pattern)
awslocal secretsmanager create-secret \
  --name test \
  --description "My test secret" \
  --secret-string '{"username":"admin","password":"password"}'

# Single value secret (for URL?key=xxx pattern and URL only pattern)
awslocal secretsmanager create-secret \
  --name simple_secret \
  --description "Simple string secret" \
  --secret-string "simpleValue123"

# More complex JSON secret (for testing extraction of multiple values from JSON with URL only pattern)
awslocal secretsmanager create-secret \
  --name complex_json \
  --description "Complex JSON secret" \
  --secret-string '{"db":{"username":"dbuser","password":"dbpass","host":"db.example.com","port":5432},"api":{"key":"api-key-12345","endpoint":"https://api.example.com"}}'

# Secret with Japanese characters (for testing multibyte characters)
awslocal secretsmanager create-secret \
  --name japanese_secret \
  --description "Secret with Japanese characters" \
  --secret-string '{"message":"メッセージ","user":"ユーザー","company":"カンパニー"}'

# Secret with special characters (for testing special character handling)
awslocal secretsmanager create-secret \
  --name special_chars \
  --description "Secret with special characters" \
  --secret-string '{"value":"!@#$%^&*()_+","query":"SELECT * FROM users;","path":"/usr/local/bin"}'

# Secret with arrays (for testing array handling)
awslocal secretsmanager create-secret \
  --name array_secret \
  --description "Secret with arrays" \
  --secret-string '{"items":["item1","item2","item3"],"users":[{"name":"John","role":"admin"},{"name":"Alice","role":"user"}]}'

# Secret with different versions (for testing version specification)
awslocal secretsmanager create-secret \
  --name versioned_secret \
  --description "Versioned secret - initial version" \
  --secret-string '{"version":"v1","status":"active"}'
  
awslocal secretsmanager update-secret \
  --secret-id versioned_secret \
  --secret-string '{"version":"v2","status":"pending"}'

# Testing for KEY=URL with key pattern (accessing specific keys in JSON)
awslocal secretsmanager create-secret \
  --name nested_fields \
  --description "Secret with nested fields" \
  --secret-string '{"production":{"db_username":"prod_user","db_password":"prod_pass"},"staging":{"db_username":"staging_user","db_password":"staging_pass"}}'

# JSON secret with empty values (for testing special case handling)
awslocal secretsmanager create-secret \
  --name empty_values \
  --description "Secret with empty values" \
  --secret-string '{"empty_string":"","null_value":null,"defined_value":"exists"}'

# Large JSON secret (for testing large data processing)
awslocal secretsmanager create-secret \
  --name large_secret \
  --description "Large JSON secret" \
  --secret-string '{
    "app_name": "test-application",
    "environment": "production",
    "database": {
      "host": "db.example.com",
      "port": 5432,
      "username": "prod_db_user",
      "password": "very_secure_password_123",
      "ssl": true,
      "max_connections": 100,
      "timeout_seconds": 30
    },
    "api_keys": {
      "google": "google-api-key-12345",
      "aws": "aws-api-key-67890",
      "github": "github-api-key-abcdef"
    },
    "feature_flags": {
      "new_ui": true,
      "beta_features": false,
      "maintenance_mode": false
    },
    "cache": {
      "enabled": true,
      "ttl_seconds": 3600,
      "max_size_mb": 512
    },
    "log_levels": {
      "production": "ERROR",
      "staging": "INFO",
      "development": "DEBUG"
    },
    "contacts": [
      {
        "name": "Admin Team",
        "email": "admin@example.com",
        "phone": "+81-3-1234-5678"
      },
      {
        "name": "Support Team",
        "email": "support@example.com",
        "phone": "+81-3-8765-4321"
      }
    ]
  }'

# Secret with different newline codes (for testing handling of CRLF, LF, CR)
awslocal secretsmanager create-secret \
  --name newline_variants \
  --description "Secret with different newline characters" \
  --secret-string '{"unix_style":"line1\\nline2\\nline3","windows_style":"line1\\r\\nline2\\r\\nline3","old_mac_style":"line1\\rline2\\rline3","mixed_style":"line1\\nline2\\r\\nline3\\r"}'

# Secret with emojis (for testing Unicode emoji handling)
awslocal secretsmanager create-secret \
  --name emoji_secret \
  --description "Secret with emoji characters" \
  --secret-string '{"reaction":"👍","weather":"☀️🌧️❄️","faces":"😀😎🤔😱","flags":"🇯🇵🇺🇸🇪🇺","combined":"Family trip to 🏖️! 🎉"}'

# Secret with a very long value
long_str=$(printf '%0.s-' {1..10000})
awslocal secretsmanager create-secret \
  --name very_long_value \
  --description "Secret with a very long string value" \
  --secret-string "{\"long_text\":\"$long_str\"}"



# Add a third version to the versioned secret for testing version selection
awslocal secretsmanager update-secret \
  --secret-id versioned_secret \
  --secret-string '{"version":"v3","status":"deployed"}'

# Create a secret with numbered versions for explicit version ID testing
awslocal secretsmanager create-secret \
  --name explicit_version_secret \
  --description "Secret for testing explicit version ID retrieval" \
  --secret-string '{"message":"This is version 1","version_num":1}'

# Store the version ID for reference
VERSION_ID_1=$(awslocal secretsmanager describe-secret --secret-id explicit_version_secret --query 'VersionIdsToStages[0].VersionId' --output text)
echo "Created explicit_version_secret with initial version ID: $VERSION_ID_1" 

# Create second version
awslocal secretsmanager update-secret \
  --secret-id explicit_version_secret \
  --secret-string '{"message":"This is version 2","version_num":2}'

# Store the second version ID
VERSION_ID_2=$(awslocal secretsmanager describe-secret --secret-id explicit_version_secret --query 'VersionIdsToStages[0].VersionId' --output text)
echo "Updated explicit_version_secret with new version ID: $VERSION_ID_2"

# Create third version (which will become the AWSCURRENT/latest version)
awslocal secretsmanager update-secret \
  --secret-id explicit_version_secret \
  --secret-string '{"message":"This is version 3 (latest)","version_num":3}'