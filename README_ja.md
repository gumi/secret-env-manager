# Secret Env Manager (sem)

![GitHub Release](https://img.shields.io/github/v/release/gumi/secret-env-manager)

[English](./README.md) / [日本語](./README_ja.md)

Google CloudやAWSのIAM権限を使用してシークレットを安全に取得し、環境変数としてenvファイルに保存するCLIツールです。クラウドプロバイダのアクセス制御を通じてセキュリティを維持しながら、シークレット管理を簡素化します。

---

## 特徴

- **マルチアカウントサポート**
- **プロバイダ:**
  - Google Cloud Secret Manager
  - AWS Secrets Manager
- **コマンド:**
  - `init`: AWSまたはGoogle Cloudからのインタラクティブなシークレット選択（実行時にプロバイダを選択）
  - `load`: envファイルから環境変数を読み込み出力（オプションで`export`プレフィックスあり）
  - `update`: クラウドプロバイダからenvファイル内のシークレットを更新

---

## インストール

### Makeを使う場合
```bash
make install
```

### アンインストール
```bash
make uninstall
```

### Homebrew (macOS)
```bash
brew install gumi/tap/secret-env-manager
```

---

## 使い方

### CLIヘルプ
```bash
sem -h
```

### コマンド
| コマンド | 説明 |
|---------|-------------|
| `init`  | AWSとGoogle Cloudの両方のプロバイダからのインタラクティブなシークレット選択 |
| `load`  | キャッシュされたシークレットから環境変数を出力 |
| `update`| 最新の値を取得してキャッシュされたシークレットを更新 |

#### プロバイダに必要な環境変数

AWS Secrets Managerの場合:
- `AWS_PROFILE`: 使用するAWSプロファイル
- `AWS_REGION`: クエリするAWSリージョン

Google Cloud Secret Managerの場合:
- `GOOGLE_CLOUD_PROJECT`: Google CloudプロジェクトのプロジェクトID

### direnvとの併用

Secret Env Managerはdirenvと連携することで、プロジェクトディレクトリに入った際に自動的に環境変数をロードすることができます。設定方法は次の通りです：

1. まず、シークレットURIを含む環境変数ファイル（例：`.env`）を作成します：
   ```
   DB_PASSWORD=sem://aws:secretsmanager/dev-profile/database/credentials?key=password
   API_KEY=sem://googlecloud:secretmanager/my-project/api-key
   ```

2. `update`コマンドを実行して、シークレットを取得しキャッシュファイルを生成します：
   ```bash
   sem update --input .env
   ```
   これにより、実際のシークレット値を含む`.cache.env`という名前のキャッシュファイルが作成されます。

3. `.envrc`ファイルに以下の内容を追加します：
   ```bash
   # シークレットキャッシュを更新（任意ですが、最新の値を保証します）
   sem update --input .env
   
   # キャッシュファイルから環境変数を読み込む
   dotenv .cache.env
   ```

4. direnv設定を許可します：
   ```bash
   direnv allow
   ```

これで、プロジェクトディレクトリに入るたびに、direnvが自動的にキャッシュファイルから環境変数を読み込み、アプリケーションでシークレットが利用できるようになります。

---

## Envファイルの書き方

Envファイルではさまざまな方法でシークレットを指定できます。各行はSecretURI形式に従います。

### 基本形式

```
# コメント行は # で始まります
KEY=VALUE  # 通常の環境変数（直接値を代入）

# シークレットURI - プロバイダーからシークレットを取得するために使用
sem://aws:secretsmanager/profile/path/secret-name
ENV_VAR=sem://aws:secretsmanager/profile/path/secret-name
```

同じenvファイル内で直接値の代入とシークレットURIを混在させることができます。直接値の代入はそのまま保持され、シークレットURIはクラウドプロバイダーから値を取得して処理されます。

### AWS Secretsの例

#### 1. JSONシークレットからすべてのキーと値を取得

**envファイルの記述：**
```
sem://aws:secretsmanager/xxx-profile/test/test_secret
```

**処理後の結果：**
```
key1=value1
key2=value2
...
keyN=valueN
```

#### 2. JSONシークレットのすべてのキーにプレフィックスを追加

**envファイルの記述：**
```
KEY=sem://aws:secretsmanager/xxx-profile/test/test_secret
```

**処理後の結果：**
```
KEY_key1=value1
KEY_key2=value2
...
KEY_keyN=valueN
```

#### 3. JSONシークレットから特定のキーを取得

**envファイルの記述：**
```
sem://aws:secretsmanager/xxx-profile/test/test_secret?key=username
```

**処理後の結果：**
```
username=value
```

#### 4. 特定のキーを環境変数に割り当て

**envファイルの記述：**
```
DB_USER=sem://aws:secretsmanager/xxx-profile/test/test_secret?key=username
```

**処理後の結果：**
```
DB_USER=value
```

### Google Cloud Secretsの例

#### 1. 基本的なシークレット取得

**envファイルの記述：**
```
sem://googlecloud:secretmanager/xxx-project/sample_secret
```

**処理後の結果：**
```
sample_secret=value
```

#### 2. カスタム環境変数名の指定

**envファイルの記述：**
```
API_KEY=sem://googlecloud:secretmanager/xxx-project/sample_secret
```

**処理後の結果：**
```
API_KEY=value
```

### Envファイルの完全な例

```
# データベース認証情報
DB_HOST=localhost
DB_PORT=5432
DB_NAME=myapp
DB_USER=sem://aws:secretsmanager/dev-profile/database/credentials?key=username
DB_PASSWORD=sem://aws:secretsmanager/dev-profile/database/credentials?key=password

# APIキー
GOOGLE_API_KEY=sem://googlecloud:secretmanager/my-project/google-api-key
PAYMENT_API_SECRET=sem://aws:secretsmanager/payments-profile/payment/api-keys?key=secret

# メールサービスの設定をプレフィックス付きですべてインポート
EMAIL_SERVICE=sem://aws:secretsmanager/marketing-profile/email-service/settings

# 特定バージョンのシークレットを取得
CONFIG_V1=sem://aws:secretsmanager/dev-profile/versioned-config?version=v1
```

---

## SecretURI仕様

```
EXPORT_NAME=sem://<Platform>:<Service>/<Account>/<SecretName>?version=<Version>&key=<Key>
```

| 項目        | 説明 |
|-------------|------|
| Platform    | クラウド種別（`aws` or `gcp`） |
| Service     | シークレットサービス名（AWSは`secretsmanager`、GCPは`secretmanager`） |
| Account     | アカウント/プロファイル/プロジェクト名 |
| SecretName  | シークレット名 |
| ExportName  | 環境変数名 |
| Version     | シークレットバージョン（AWSは`AWSCURRENT`、GCPは`latest`） |
| Key         | （AWSのみ、JSONシークレット用）抽出するキー名 |

> **注意:** Google Cloudの場合、値がJSON形式の場合のみkeyを指定できます。

---

## License

Apache License, Version 2.0
