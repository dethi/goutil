# goutil [![GoReportCard](https://goreportcard.com/badge/github.com/dethi/goutil)](https://goreportcard.com/report/github.com/dethi/goutil)

Yet another monolithic repository for Go packages and binaries.

The packages are small. It is a good idea to just copy/paste what you need instead of adding this repository as a dependency. This repository will probably never follow semantic versioning.

## Packages

- `envflag`: `flag` with support for environment variable
- `errgroup`: `errgroup` with bounded concurrently
- `fs`: filesystem related functions
- `safepprof`: `pprof` configured to be used in production
- `sqlstore`: `sql.DB` with prepared statement cache

## Binaries

- `cmd/statico`: Yet another static server
- `cmd/jsony`: Stdin JSON formatter
