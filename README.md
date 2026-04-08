# tinyjs

This is a deliberately tiny teaching compiler frontend written in Go.

At this stage, the language supports exactly one statement form:

```js
x = 1;
y = 42;
```

That means:
- an identifier on the left
- a decimal integer literal on the right
- a trailing semicolon

Everything else is rejected on purpose.

## Why it is structured this way

The project is intentionally split into a few small packages:

- `internal/ast`: syntax tree types
- `internal/parser`: Participle lexer and parser setup
- `internal/compiler`: tiny pseudo-backend
- `cmd/tinyjs`: CLI entry point

That keeps the first version trivial while leaving room to evolve toward a more realistic compiler.

## Why there is no testcontainers usage

There is nothing external to containerize yet. The project has:
- no database
- no message broker
- no network dependency
- no external compiler service

So testcontainers would add ceremony without adding value.

Once you later introduce something like an external toolchain, object-file linker, embedded runtime service, or maybe a golden-test environment that depends on a real executable tool, then testcontainers could become useful. For this milestone, ordinary unit tests are the right tool.

## Run

```bash
go run ./cmd/tinyjs
```

Then type or pipe source such as:

```js
x = 1;
y = 2;
```

## Test

```bash
go test ./...
```
