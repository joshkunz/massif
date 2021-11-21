# massif

[![API Reference](https://img.shields.io/badge/go.pkg.dev-reference-5272B4)](
https://pkg.go.dev/github.com/joshkunz/massif?tab=doc)
[![Build Status](https://github.com/joshkunz/massif/actions/workflows/test.yaml/badge.svg)](https://github.com/joshkunz/massif/actions/workflows/test.yaml)
[![LICENSE](
https://img.shields.io/github/license/joshkunz/massif?color=informational)](
LICENSE)

`massif` is a library for working with logs produced by [`massif`](
https://valgrind.org/docs/manual/ms-manual.html), a tool contained in the
Valgrind suite. This library can make it easy to do things like write an
automated test to verify peak heap usage in a binary. The linked godoc contains
usage examples.
