# env-ssm

`env-ssm` is a command-line interface (CLI) tool designed to synchronize your local `.env` file with AWS Systems Manager (SSM) Parameter Store.

## Download

The binaries are published to the [Releases](https://github.com/localopsco/env-ssm/releases) page. You can download the binary for your OS and architecture.

### Usage

```
envssm
```

This will open a TUI to enter the env file path and the AWS SSM path. Which will be used to sync the env file with AWS SSM.

While syncing, it will keep the old env in SSM if the `keep old env in SSM` option is not selected.
