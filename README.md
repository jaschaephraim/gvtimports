# gvtimports

Runs `gvt fetch` for all imports in a package, including imports in tests and subpackages.

```bash
go get github.com/jaschaephraim/gvtimports
cd path/to/package
gvtimports
```

To only list the imports to be fetched without actually fetching, pass the `--ls` flag.
