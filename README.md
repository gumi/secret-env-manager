# Secret Env Manager (sem)

## Description
GCP や AWS に保存してあるシークレットを環境変数にロードするツールです。
- init コマンドは、環境変数にロードするシークレットを選択し、env ファイルを生成します。ただし、リストは１アカウントからしかとってこないので、簡単にenvファイルを作りたいときにだけ使えます。複数アカウントからシークレットを取得したい場合は、直接envファイルを編集してください。
- load コマンドは、env ファイルを読み込み、export文を組み立てます。exportはキャッシュされるので、毎回シークレットを取得する必要はありません。
- update コマンドは、env ファイルを読み込み、シークレットを更新します。


## Support
- [x] Multi Account Support
- [ ] Service
  - [x] GCP Secret Manager
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

$ sem load
export SAMPLE_SECRET='knCFtym5URfRY#W9oaGUYGmxs4p'
export TEST_TEST_SECRET='{"test_secret":"test_secretXXXXX"}'

$ eval $(sem load)
$ echo $SAMPLE_SECRET
knCFtym5URfRY#W9oaGUYGmxs4p
$ echo $TEST_TEST_SECRET
'{"test_secret":"test_secretXXXXX"}'


```
## env (plain) example

```txt
SAMPLE_SECRET=sem://gcp:secretmanager/xxx-project/sample_secret
TEST_TEST_SECRET=sem://aws:secretsmanager/xxx-profile/test/test_secret
```

### format
`EXPORT_NAME=sem://<Platform>:<Service>/<Account>/<SecretName>?version=<Version>`

## env.toml example

```toml
[[Environments]]
  Platform = "gcp"
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
