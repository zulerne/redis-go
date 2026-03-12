# Contributing

## Quick Start

```bash
task check    # run lint + tests
```

## Prerequisites

- Go 1.26+
- [Task](https://taskfile.dev/)
- [golangci-lint](https://golangci-lint.run/)

## Commands

| Command         | Description              |
|-----------------|--------------------------|
| `task run`      | Run the server           |
| `task build`    | Build the binary         |
| `task test`     | Run tests with `-race`   |
| `task lint`     | Run golangci-lint        |
| `task check`    | Lint + test              |
| `task coverage` | Test with coverage       |
| `task fmt`      | Format code              |

## Guidelines

- Run `task check` before pushing
- Use [Conventional Commits](https://www.conventionalcommits.org/): `feat:`, `fix:`, `refactor:`, `test:`, `docs:`, `chore:`
- Keep PRs focused — one concern per PR
- Add tests for new functionality

## Reporting Bugs

Open an issue with steps to reproduce.
