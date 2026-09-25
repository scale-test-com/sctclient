# Scale-Test CLI Usage

## Overview

`scale-test` is a CLI client for the Scale-Test API. It supports creating, retrieving, deleting and analysing (results, time series) load test runs.

## Authentication

Provide your API key with:

- `--api-key <API_KEY>`
- or `SCALE_TEST_API_KEY`

The CLI uses `--api-key` first. If no flag is provided, it reads `SCALE_TEST_API_KEY`.

## API base URL

The API base URL can be overridden with:

- `--base-url <URL>`
- or `SCALE_TEST_BASE_URL`

Default: `https://scale-test.com/api/v1`

## Commands

### `run create`

Create a new run.

Examples:

```bash
./scale-test --api-key <API_KEY> run create --scenario-id 123
```

```bash
./scale-test run create --file scenario.yaml
```

```bash
./scale-test run create --scenario-id 123 --wait --poll-interval 2s
```

### `run get`

Retrieve run details.

```bash
./scale-test run get <RUN_UUID>
```

### `run delete`

Delete a run.

```bash
./scale-test run delete <RUN_UUID>
```

### `run results`

Retrieve the aggregated results of a finished run: total requests, success rate, average duration,
requests per second (average and peak) and HTTP status code breakdown.
While the run is still in progress, the API answers `409`.

```bash
./scale-test run results <RUN_UUID>
```

### `run timeseries`

Retrieve the per-second metrics of a finished run (UTC timestamps). Use `--metric` to restrict the
output to one series: `status_codes`, `success_rate`, `average_duration` or `requests_per_second`.

```bash
./scale-test run timeseries <RUN_UUID>
./scale-test run timeseries <RUN_UUID> --metric requests_per_second
```

## Example YAML

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

## Notes

- JSON output is printed to `stdout`
- Progress messages are printed to `stderr` when using `--wait`
