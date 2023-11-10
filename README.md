# Secret Env Manager (sem)

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
$ sem init
project: xxx-project
Select secrets:
secret name (created at)
▶ [x] sample_secret(2023-01-01T00:00:00+00:00)
env.toml has been saved.
ExportName can be changed to any value you like.


$ eval $(sem load)
$ echo $SAMPLE_SECRET
knCFtym5URfRY#W9oaGUYGmxs4p

```

## env.toml example

```toml
[aws]
  Profile = ""

[gcp]
  Project = "xxx-project"

  [[gcp.env]]
    SecretName = "sample_secret"
    ExportName = "SAMPLE_SECRET"
    Version = "latest"
```