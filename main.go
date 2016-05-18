package main

import (
	"go/build"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	pwd, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	pkg, err := build.Import(".", pwd, 0)
	if err != nil {
		log.Fatal(err)
	}

	s := newSet()
	addExternalImports(s, pkg.Imports, pkg.ImportPath)
	addExternalImports(s, pkg.TestImports, pkg.ImportPath)

	for _, imp := range s.elements() {
		cmd := exec.Command("gvt", "fetch", imp)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}
}

func addExternalImports(s *set, imports []string, baseImport string) {
	for _, name := range imports {
		imp, err := build.Import(name, ".", 0)
		if err != nil || !imp.Goroot && !strings.HasPrefix(imp.ImportPath, baseImport) {
			s.add(imp.ImportPath)
		}
	}
}
