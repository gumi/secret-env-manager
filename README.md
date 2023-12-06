# Secret Env Manager (sem)

## Description
googlecloud や AWS に保存してあるシークレットを環境変数にロードするツールです。
- init コマンドは、環境変数にロードするシークレットを選択し、env ファイルを生成します。ただし、リストは１アカウントからしかとってこないので、簡単にenvファイルを作りたいときにだけ使えます。複数アカウントからシークレットを取得したい場合は、直接envファイルを編集してください。
- load コマンドは、env ファイルを読み込み、export文を組み立てます。exportはキャッシュされるので、毎回シークレットを取得する必要はありません。
- update コマンドは、env ファイルを読み込み、シークレットを更新します。


## Support
- [x] Multi Account Support
- [ ] Service
  - [x] googlecloud Secret Manager
  - [x] AWS Secrets Manager
  - [ ] openstack ?

## Install

```bash
make install
```

## uninstall
    
```bash
make uninstall
```

## Install (macOS, Homebrew)

```bash
$ brew install gumi/tap/secret-env-manager
```

## Usage

```bash
# other format
# $ sem init -f toml
$ sem init
File `env` already exists. Do you want to overwrite it? (yes/no): yes
------------------------------------------
project : xxx-project
Select secrets:
secret name (created at)
▶ [x] sample_secret (2023-01-01T00:00:00+00:00)
==========================================
profile : xxx-profile
Select secrets:
secret name (created at)
▶ [x] test/test_secret (2023-01-01 08:00:00.000 +0000 UTC)
------------------------------------------
env has been saved.

$ sem update
cache updated.

$ sem load
SAMPLE_SECRET='knCFtym5URfRY#W9oaGUYGmxs4p'
TEST_TEST_SECRET='{"test_secret":"test_secretXXXXX"}'

# if you use fish shell
$ env -i $(sem load) env
SAMPLE_SECRET=knCFtym5URfRY#W9oaGUYGmxs4p
TEST_TEST_SECRET='{"test_secret":"test_secretXXXXX"}'

# if you use bash or zsh shell
$ eval $(sem load -e) | env
…
SAMPLE_SECRET=knCFtym5URfRY#W9oaGUYGmxs4p
TEST_TEST_SECRET='{"test_secret":"test_secretXXXXX"}'
…

```

## params

- Platform
  - クラウドプラットフォームを指定します。現在は、gcp と aws が指定可能です。
- Service
  - シークレットを保存しているサービスを指定します。現在は、gcp の secretmanager と aws の secretsmanager が指定可能です。
- Account
  - シークレットを保存しているアカウントを指定します。
- SecretName
  - シークレット名を指定します。
- ExportName
  - 環境変数名を指定します。
- Version
  - シークレットのバージョンを指定します。gcp は latest と指定すると最新のバージョンを取得します。aws は AWSCURRENT と指定すると最新のバージョンを取得します。
- Key
  - シークレットのキーを指定します。gcp は指定しないでください。aws かつ 値が json の場合は指定したキーの値を取得します。


## env (plain) example

```txt
SAMPLE_SECRET=sem://googlecloud:secretmanager/xxx-project/sample_secret
TEST_TEST_SECRET=sem://aws:secretsmanager/xxx-profile/test/test_secret
```

### format
`EXPORT_NAME=sem://<Platform>:<Service>/<Account>/<SecretName>?version=<Version>&key=<Key>`

## env.toml example

```toml
[[Environments]]
  Platform = "googlecloud"
  Service = "secretmanager"
  Account = "xxx-project"
  SecretName = "sample_secret"
  ExportName = "SAMPLE_SECRET"
  Version = "latest"


[[Environments]]
  Platform = "aws"
  Service = "secretsmanager"
  Account = "xxx-profile"
  SecretName = "test/test_secret"
  ExportName = "TEST_TEST_SECRET"
  Version = "AWSCURRENT"

```


# Help
```bash
# sem -h

NAME:
   secret-env-manager - manage secret environment variables

USAGE:
   secret-env-manager [global options] command [command options] [arguments...]

COMMANDS:
   init     Save the credentials stored in googlecloud Secret Manager as file.
   load     Output a string to read credentials from SecretManager based on the file and export them as environment variables.
   update   Forcefully update the cached information for the load command.
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help

$ sem init -h

USAGE:
   secret-env-manager init [command options] [arguments...]

OPTIONS:
   --file value, -f value  Load configuration from the specified format file. Available formats are Plain and Toml, and the default is Plain.
   --help, -h              show help


$ sem load -h

USAGE:
   secret-env-manager load [command options] [arguments...]

OPTIONS:
   --file value, -f value  Load configuration from the specified format file. Available formats are Plain and Toml, and the default is Plain.
   --with-export, -e       When this option is enabled, 'export' is added when displaying to standard output. (default: false)
   --help, -h              show help


$ sem update -h

USAGE:
   secret-env-manager update [command options] [arguments...]

OPTIONS:
   --file value, -f value  Load configuration from the specified format file. Available formats are Plain and Toml, and the default is Plain.
   --with-quote, -q        When this option is enabled, single quotes are added to the cached environment variable values. It is generally recommended to use this option when the value contains spaces. (default: false)
   --help, -h              show help

```

