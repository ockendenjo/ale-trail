# Wassail Ale Trail

## tasks

### apply

directory: stack
environment: AWS_PROFILE=aletrail

```shell
terraform apply -var-file=tfvars/dev.auto.tfvars -auto-approve
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

### plan

directory: stack
environment: AWS_PROFILE=aletrail

```shell
terraform plan -var-file=tfvars/dev.auto.tfvars
```

### serve

directory: src

```shell
npm run start
```
