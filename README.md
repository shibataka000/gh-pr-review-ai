# gh-pr-review-ai

Review a pull request using [Amazon Bedrock](https://aws.amazon.com/bedrock/).

## Requirements

- [GitHub CLI.](https://cli.github.com/)
- [Configuration and credentials for AWS SDK should be configured.](https://docs.aws.amazon.com/sdkref/latest/guide/overview.html)
- [Access to Amazon Bedrock foundation models should be added.](https://docs.aws.amazon.com/bedrock/latest/userguide/model-access-modify.html)

## Usage

```
$ gh-pr-review-ai --help
Review a pull request using a generative AI.

Usage:
  gh-pr-review-ai [<url>] [flags]

Flags:
  -h, --help                 help for gh-pr-review-ai
      --instruction string   The prompt that provides instructions to the model about the task it should perform. (default "You are an IT expert. Review the pull request below and tell me points which should be fixed or improved if found. Write your response in Japanese.")
      --model-id string      The model with which to run inference. (default "anthropic.claude-3-5-sonnet-20240620-v1:0")
```

## Install

```
go install github.com/shibataka000/gh-pr-review-ai@main
```
