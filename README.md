# Wassail Ale Trail

## tasks

### apply

directory: stack
environment: AWS_PROFILE=aletrail

```shell
terraform apply -var-file=tfvars/dev.auto.tfvars -auto-approve
```

### build-cmd

```shell
go run ./scripts/build-cmd --zip
```

### env

```shell
set -o allexport; source .env; set +o allexport
```

### format

requires: format-tf, format-ng

### format-ng

directory: src

```shell
npx prettier --write .
```

### format-tf

directory: stack

```shell
terraform fmt --recursive
```

### init

directory: stack
environment: AWS_PROFILE=aletrail

```shell
terraform init -backend-config=tfvars/dev.backend.tfvars
```

### push-config

environment: AWS_PROFILE=aletrail
directory: stack

```shell
aws s3 cp tfvars/dev.auto.tfvars s3://ale-trail-dev-state-20260822135440864900000002
aws s3 cp tfvars/dev.backend.tfvars s3://ale-trail-dev-state-20260822135440864900000002
aws s3 cp tfvars/pro.auto.tfvars s3://ale-trail-pro-state-20261107134939439300000002
aws s3 cp tfvars/pro.backend.tfvars s3://ale-trail-pro-state-20261107134939439300000002
```

### pull-config

environment: AWS_PROFILE=aletrail
directory: stack

```shell
mkdir -p tfvars
aws s3 cp s3://ale-trail-dev-state-20260822135440864900000002/dev.auto.tfvars tfvars/
aws s3 cp s3://ale-trail-dev-state-20260822135440864900000002/dev.backend.tfvars tfvars/
```

### push-dev-data

environment: AWS_PROFILE=aletrail
inputs: DATA_BUCKET

```shell
aws s3 cp data.json s3://${DATA_BUCKET}/data.json
```

### push-pro-data

environment: AWS_PROFILE=aletrail
inputs: PRO_DATA_BUCKET

```shell
aws s3 cp pro.data.json s3://${PRO_DATA_BUCKET}/data.json
```

### plan

directory: stack
environment: AWS_PROFILE=aletrail

```shell
terraform plan -var-file=tfvars/dev.auto.tfvars
```

### sast

```shell
go fix --diff $(go list ./... | grep -v /src)
wget -O .golangci.json https://raw.githubusercontent.com/ockendenjo/actions/refs/heads/main/.golangci.json
golangci-lint run
govulncheck ./...
```

### serve

directory: src

```shell
npm run start
```

### test

```shell
go test -json $(go list ./... | grep -v /e2e) | tparse
```

### upload-cmd

requires: build-cmd
environment: AWS_PROFILE=aletrail
inputs: BINARY_BUCKET

```shell
go run ./scripts/upload-binaries
```

### test-e2e

environment: AWS_PROFILE=aletrail
environment: BASE_URL=http://localhost:4200
environment: SHOW_BROWSER=1
inputs: DATA_BUCKET

```shell
go test -count=1 ./e2e/... -v
```
