# Secret Env Manager (sem)

![GitHub Release](https://img.shields.io/github/v/release/gumi/secret-env-manager)

[English](./README.md) / [日本語](./README_ja.md)

A CLI tool that securely retrieves secrets from Google Cloud or AWS using IAM permissions and stores them as environment variables in env files. Simplifies secret management while maintaining security through cloud provider access controls.

---

## Features

- **Multi-Account Support**
- **Providers:**
  - Google Cloud Secret Manager
  - AWS Secrets Manager
- **Commands:**
  - `init`: Interactive secret selection from AWS or Google Cloud (choose provider at runtime)
  - `load`: Read env file and output environment variables (with optional `export` prefix)
  - `update`: Refresh secrets in the env file from cloud providers

---

## Installation

### Using Make
```bash
make install
```

### Uninstall
```bash
make uninstall
```

### Using Homebrew (macOS)
```bash
brew install gumi/tap/secret-env-manager
```

---

## Usage

### CLI Help
```bash
sem -h
```

### Commands
| Command | Description |
|---------|-------------|
| `init`  | Interactive secret selection from both AWS and Google Cloud providers |
| `load`  | Output environment variables from cached secrets |
| `update`| Update cached secrets by fetching latest values |

#### Environment Variables Required for Providers

For AWS Secrets Manager:
- `AWS_PROFILE`: AWS profile to use
- `AWS_REGION`: AWS region to query

For Google Cloud Secret Manager:
- `GOOGLE_CLOUD_PROJECT`: Google Cloud project ID

### Using with direnv

Secret Env Manager works seamlessly with direnv to automatically load environment variables when entering your project directory. Here's how to set it up:

1. First, create an environment file (e.g., `.env`) containing your secret URIs:
   ```
   DB_PASSWORD=sem://aws:secretsmanager/dev-profile/database/credentials?key=password
   API_KEY=sem://googlecloud:secretmanager/my-project/api-key
   ```

2. Run the `update` command to retrieve secrets and generate a cache file:
   ```bash
   sem update --input .env
   ```
   This will create a cache file named `.cache.env` containing the actual secret values.

3. Add the following to your `.envrc` file:
   ```bash
   # Update secrets cache (optional but ensures freshness)
   sem update --input .env
   
   # Load environment variables from the cache file
   dotenv .cache.env
   ```

4. Allow the direnv configuration:
   ```bash
   direnv allow
   ```

Now whenever you enter your project directory, direnv will automatically load the environment variables from the cache file, making your secrets available to your application.

---

## Env File Format Examples

The env file supports various ways to specify secrets. Each line in the file follows the SecretURI format.

### Basic Format

```
# Comment line starts with #
KEY=VALUE  # Regular environment variable (direct value assignment)

# Secret URIs - used to fetch secrets from providers
sem://aws:secretsmanager/profile/path/secret-name
ENV_VAR=sem://aws:secretsmanager/profile/path/secret-name
```

You can mix both direct value assignments and Secret URIs in the same env file. Direct value assignments are preserved as-is, while Secret URIs are processed to fetch values from cloud providers.

### AWS Secrets Examples

#### 1. Retrieving all key-value pairs from a JSON secret

**In your env file:**
```
sem://aws:secretsmanager/xxx-profile/test/test_secret
```

**After processing, becomes:**
```
key1=value1
key2=value2
...
keyN=valueN
```

#### 2. Adding a prefix to all keys in a JSON secret

**In your env file:**
```
KEY=sem://aws:secretsmanager/xxx-profile/test/test_secret
```

**After processing, becomes:**
```
KEY_key1=value1
KEY_key2=value2
...
KEY_keyN=valueN
```

#### 3. Retrieving a specific key from a JSON secret

**In your env file:**
```
sem://aws:secretsmanager/xxx-profile/test/test_secret?key=username
```

**After processing, becomes:**
```
username=value
```

#### 4. Assigning a specific key to an environment variable

**In your env file:**
```
DB_USER=sem://aws:secretsmanager/xxx-profile/test/test_secret?key=username
```

**After processing, becomes:**
```
DB_USER=value
```

### Google Cloud Secrets Examples

#### 1. Basic secret retrieval

**In your env file:**
```
sem://googlecloud:secretmanager/xxx-project/sample_secret
```

**After processing, becomes:**
```
sample_secret=value
```

#### 2. Custom environment variable name

**In your env file:**
```
API_KEY=sem://googlecloud:secretmanager/xxx-project/sample_secret
```

**After processing, becomes:**
```
API_KEY=value
```

### Complete Example of an env file

```
# Database credentials
DB_HOST=localhost
DB_PORT=5432
DB_NAME=myapp
DB_USER=sem://aws:secretsmanager/dev-profile/database/credentials?key=username
DB_PASSWORD=sem://aws:secretsmanager/dev-profile/database/credentials?key=password

# API Keys
GOOGLE_API_KEY=sem://googlecloud:secretmanager/my-project/google-api-key
PAYMENT_API_SECRET=sem://aws:secretsmanager/payments-profile/payment/api-keys?key=secret

# Import all email service settings with prefix
EMAIL_SERVICE=sem://aws:secretsmanager/marketing-profile/email-service/settings

# Version specific secret (retrieve a specific version)
CONFIG_V1=sem://aws:secretsmanager/dev-profile/versioned-config?version=v1
```

---

## SecretURI Format

```
EXPORT_NAME=sem://<Platform>:<Service>/<Account>/<SecretName>?version=<Version>&key=<Key>
```

| Field        | Description |
|--------------|-------------|
| Platform     | Cloud platform (`aws` or `gcp`) |
| Service      | Secret service (`secretsmanager` for AWS, `secretmanager` for GCP) |
| Account      | Account/profile/project name |
| SecretName   | Name of the secret |
| ExportName   | Environment variable name |
| Version      | Secret version (`AWSCURRENT` for AWS, `latest` for GCP) |
| Key          | (AWS only, for JSON secrets) Key to extract |

> **Note:** For GoogleCloud, key can only be specified when the value is in JSON format.

---

## License

Apache License, Version 2.0


