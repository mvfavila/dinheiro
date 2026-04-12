# Contributing

Contributions are welcome! Please follow the guidelines below.

## Getting started

1. Fork the repository and clone your fork.
2. Ensure you have Go 1.21 or later installed.
3. Run the existing tests to confirm everything passes before making changes:

```bash
go test ./...
```

## Making changes

- Keep changes focused. One feature or fix per pull request.
- Follow standard Go conventions (`gofmt`, `go vet`).
- Add or update tests for any new behaviour.
- All tests must pass before submitting:

```bash
go test ./...
```

## Submitting a pull request

1. Push your branch to your fork.
2. Open a pull request against the `main` branch.
3. Describe what the change does and why.
4. Reference any related issues.

## Reporting issues

Open a GitHub issue with a clear description of the problem, including the input value and the output you expected vs. what you got.

## Code style

- Use `gofmt` to format your code.
- Exported functions must have doc comments.
- Internal helpers do not need doc comments unless the logic is non-obvious.
