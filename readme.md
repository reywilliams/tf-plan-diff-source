# 🚀 Terraform Plan Diff Action: `tf-plan-diff`

Welcome to the **tf-plan-diff** action! This lightweight GitHub Action allows you to create a neat and readable diff of your Terraform plan JSON files, making it easier to visualize what changes will happen in your infrastructure. 🎉

## ✨ Features

- **Fast and Lightweight**: Built in Go and uses precompiled binaries for speedy execution.
- **Readable Output**: See exactly what will be created, modified, or destroyed in a clear diff format.
- **Flexible Options**: Customize your output based on what you want to see!

## 🛠️ How It Works

Simply pass your Terraform plan JSON file to this action, and it will output the changes in a beautiful diff format like this:

# `webhook-lambda` Plan Diff :building_construction:

```diff
- module.dynamodb_table.aws_dynamodb_table.table will be deleted
- module.github_PAT_secret.aws_secretsmanager_secret_version.this will be deleted
! aws_iam_policy.secret_access will be updated
! aws_iam_role.lambda_execution will be updated
+ module.github_PAT_secret.aws_secretsmanager_secret.this will be created
+ module.github_webhook_secret.aws_secretsmanager_secret.this will be created
+ module.github_webhook_secret.aws_secretsmanager_secret_version.this will be created
! aws_iam_policy.lambda_dynamodb_write_policy will be recreated
! aws_lambda_function.webhook will be recreated
```

:warning: This plan will: **CREATE** 3, **UPDATE** 2, **DELETE** 2, **RECREATE** 2 :warning:

## 🎛️ Inputs

The action takes the following inputs:

| Input Name     | Description                                                 | Required |
| -------------- | ----------------------------------------------------------- | -------- |
| `file_path`    | Path to the JSON plan file (ex. `tfplan.json`)              | ✅       |
| `app_name`     | Name of your application (ex. `Test App`)                   | ❌       |
| `include_noop` | Flag to include no-op actions (any value will eval to true) | ❌       |
| `include_read` | Flag to include read actions (any value will eval to true)  | ❌       |

⚡ Sample Usage

```yaml
name: Terraform Plan w/ Diff

on:
  push:
    branches:
      - main

jobs:
  plan-diff:
    runs-on: ubuntu-22.04
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Terraform Plan Diff
        uses: reywilliams/tf-plan-diff@v1.0.0
        with:
          file_path: "path/to/tfplan.json"
          app_name: "My Awesome App"
```

## 🔒 Data Privacy

This action **does not** transmit any data from your Terraform plan. It is designed as a parser that leverages HashiCorp's own [structs](https://pkg.go.dev/github.com/hashicorp/terraform-json@v0.13.0) for parsing Terraform plans. Your plan data remains local to your GitHub Actions environment, ensuring that your infrastructure details are kept private and secure - source code repo can be found [here](https://github.com/reywilliams/tf-plan-diff-source/tree/main/src).

## 📦 Installation

Simply add this action to your GitHub workflow, and you're ready to roll! No complex installation steps are needed, just a simple YAML snippet in your workflow file

## 🤝 Feedback

If you have ideas, improvements, or bug fixes, feel free to open an issue!

## 🌟 License

This project is licensed under the Apache License 2.0 [License](./LICENSE) - see the LICENSE file for details.
