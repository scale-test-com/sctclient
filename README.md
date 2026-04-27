# Scale-Test CLI

A command-line client for the Scale-Test load testing API.

## Features

- Create a load test run
- Retrieve run details
- Delete a run
- Authenticate with an API key using a flag or environment variable
- Optional wait mode until a run reaches a terminal state

## Requirements

- Go 1.22+
- A valid Scale-Test API key

## Installation

### From source

```bash
go build -o scale-test .
```

### From GitHub Releases

Download the archive for your operating system from the Releases page.

## Configuration

Global settings can be provided with flags or environment variables.

- API key:
  - flag: `--api-key`
  - env: `SCALE_TEST_API_KEY`
- API base URL:
  - flag: `--base-url`
  - env: `SCALE_TEST_BASE_URL`
  - default: `https://scale-test.com/api/v1`

Priority: flag > environment variable > default.

> The CLI requires an API key. If `--api-key` is provided, it is used first. Otherwise the CLI reads `SCALE_TEST_API_KEY`.

## Commands

### Help

```bash
./scale-test --help
./scale-test run --help
```

### Create a run using an existing scenario

```bash
./scale-test --api-key <API_KEY> run create --scenario-id 123
```

### Create a run from a YAML file

```bash
./scale-test run create --file scenario.yaml
```

### Create a run and wait until completion

```bash
./scale-test run create --scenario-id 123 --wait --poll-interval 2s
```

### Get a run

```bash
./scale-test run get <RUN_UUID>
```

### Delete a run

```bash
./scale-test run delete <RUN_UUID>
```

## Example YAML scenario

```yaml
name: Test API
request_timeout: 5s
req_target_curve:
  - elapsed_time: 0s
    req_per_sec: 10
  - elapsed_time: 30s
    req_per_sec: 50
  - elapsed_time: 60s
    req_per_sec: 10
operations:
  - uri: https://api.example.com/
    method: GET
```

## Output

- JSON results are written to `stdout`
- Progress messages (in `--wait` mode) are written to `stderr`

## Localized documentation

For more usage details, see the `doc/` folder:

- `doc/en.md` — English usage guide
- `doc/fr.md` — French usage guide
