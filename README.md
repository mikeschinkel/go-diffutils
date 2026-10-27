# go-diffutils

Canonical types and algorithms for computing and representing text diffs in Go.

The root package compares string slices (`DiffLines`) and produces the canonical interchange types `DiffContent`, `DiffChange` and `LineRange`, using a port of Python's difflib `SequenceMatcher`. Two adapter modules bridge major Go diff libraries.

## Modules

| Module | Import path | What it is |
|--------|-------------|------------|
| root   | `github.com/mikeschinkel/go-diffutils` | Canonical diff types, `DiffLines`, and the difflib engine (`SequenceMatcher`, unified and context diffs) |
| dmp    | `github.com/mikeschinkel/go-diffutils/dmp` | An `InlineHighlighter` using diff-match-patch (via [sergi/go-diff](https://github.com/sergi/go-diff)) for character-level highlighting |
| sgdiff | `github.com/mikeschinkel/go-diffutils/sgdiff` | Parses unified diff text into `DiffContent` (via [sourcegraph/go-diff](https://github.com/sourcegraph/go-diff)) |

Each is a separate Go module with its own version. Install the ones you use:

```sh
go get github.com/mikeschinkel/go-diffutils
go get github.com/mikeschinkel/go-diffutils/dmp
go get github.com/mikeschinkel/go-diffutils/sgdiff
```

Root is tagged `vX.Y.Z`; dmp and sgdiff are tagged `dmp/vX.Y.Z` and `sgdiff/vX.Y.Z`. All three are provisional as of v0.1.0: the API may change in minor releases.

## Local Development

dmp and sgdiff require a **published** version of the root module and carry no `replace` directive, because a consumer never sees a replace — only the require — and a replace would let a wrong require go unnoticed.

To build every module from the code on disk, create a Go workspace:

```sh
just work
```

`just work` writes a `go.work` listing every directory that holds a `go.mod`, so it needs no edit when a module is added. Re-run it after adding one.

`go.work` is not committed (it is in `.gitignore`). CI builds and tests with `GOWORK=off`, against the versions each `go.mod` requires, so a change in dmp or sgdiff that depends on an unpublished root change fails in CI until that root version is tagged and required. To see what CI sees locally, prefix a command with `GOWORK=off`:

```sh
cd dmp && GOWORK=off go test ./...
```

`just tidy` runs `go mod tidy` in every module with `GOWORK=off`, and CI fails when any `go.mod` is not tidy. Run `just` to list the other recipes (`test`, `lint`, `vet`, `build`, `ci`, …).

## Releasing

Releases are made by the **Release** GitHub Actions workflow: choose the module (root, dmp or sgdiff) and the version. It refuses a module whose `go.mod` declares the wrong path or carries a `replace`, or whose tag already exists, runs that module's tests, then tags and publishes it.

## License

MIT; see [LICENSE](LICENSE).

`difflib.go` is derived from [pmezard/go-difflib](https://github.com/pmezard/go-difflib), Copyright (c) 2013 Patrick Mezard, and is used under its BSD 3-Clause license, whose full text is included in [LICENSE](LICENSE). Thanks to Patrick Mezard for the port of Python's difflib it is built on.
