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

## GitHub Action

This repository ships a reusable GitHub Actions workflow that compiles `.tjs`
files to native binaries. It builds the compiler from source at the pinned ref,
so no pre-published release is required.

### Usage

In your repository, create a workflow file that calls this one:

```yaml
# .github/workflows/build.yml
name: Build

on:
  push:
    branches: [main]
  pull_request:

jobs:
  build:
    uses: your-org/tinyjs-participle/.github/workflows/compile.yml@main
    with:
      glob: "src/**/*.tjs"
```

Replace `your-org` with the GitHub organization or user that owns this
repository.

### Inputs

| Input | Required | Default | Description |
|---|---|---|---|
| `glob` | yes | — | Glob pattern for `.tjs` files relative to the repo root, e.g. `"src/**/*.tjs"` |
| `platforms` | no | `["ubuntu-latest","macos-latest"]` | JSON array of GitHub-hosted runner labels to build for |
| `emit_ir` | no | `false` | When `true`, also upload LLVM IR (`.ll`) files as a workflow artifact |

### Artifacts

After the workflow completes, GitHub makes the following artifacts available
on the run:

- **`tinyjs-bin-<runner>`** — one artifact per requested platform (e.g.
  `tinyjs-bin-ubuntu-latest`, `tinyjs-bin-macos-latest`), each containing
  one binary per source file. A `foo.tjs` input produces a binary named `foo`.
- **`tinyjs-ir`** — present only when `emit_ir: true`. Contains one `.ll` file
  per source file. IR is platform-independent so it is uploaded once regardless
  of how many platforms are targeted.

### Examples

Build for Linux only:

```yaml
jobs:
  build:
    uses: your-org/tinyjs-participle/.github/workflows/compile.yml@main
    with:
      glob: "*.tjs"
      platforms: '["ubuntu-latest"]'
```

Build for both platforms and emit IR:

```yaml
jobs:
  build:
    uses: jeffscottbrown/tinyjs/.github/workflows/compile.yml@main
    with:
      glob: "src/**/*.tjs"
      platforms: '["ubuntu-latest","macos-latest"]'
      emit_ir: true
```

Pin to a specific release instead of `main` to get reproducible builds:

```yaml
    uses: jeffscottbrown/tinyjs/.github/workflows/compile.yml@v1.0.0
```
