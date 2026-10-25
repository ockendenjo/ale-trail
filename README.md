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
```

### pull-config

environment: AWS_PROFILE=aletrail
directory: stack

```shell
mkdir -p tfvars
aws s3 cp s3://ale-trail-dev-state-20260822135440864900000002/dev.auto.tfvars tfvars/
aws s3 cp s3://ale-trail-dev-state-20260822135440864900000002/dev.backend.tfvars tfvars/
```

### push-data

environment: AWS_PROFILE=aletrail
inputs: DATA_BUCKET

```shell
aws s3 cp data.json s3://${DATA_BUCKET}/data.json
```

### plan

directory: stack
environment: AWS_PROFILE=aletrail

```shell
terraform plan -var-file=tfvars/dev.auto.tfvars
```

### sast

```shell
go fix --diff ./...
wget -O .golangci.json https://raw.githubusercontent.com/ockendenjo/actions/refs/heads/main/.golangci.json
golangci-lint run
govulncheck ./...
```

### serve

directory: src

```shell
npm run start
```

### upload-cmd

requires: build-cmd
environment: AWS_PROFILE=aletrail
inputs: BINARY_BUCKET

```shell
go run ./scripts/upload-binaries
```
