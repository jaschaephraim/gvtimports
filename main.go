package main

import (
	"fmt"
	"go/build"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	pwd, err := filepath.Abs(".")
	if err != nil {
		log.Fatal(err)
	}

	pkgs := []*build.Package{}
	pkgs, err = addPackages(pkgs, ".", pwd)
	if err != nil {
		log.Fatal(err)
	}

	s := newSet()
	for _, pkg := range pkgs {
		addExternalImports(s, pkg.Imports, pkgs)
		addExternalImports(s, pkg.TestImports, pkgs)
	}

	for _, imp := range s.elements() {
		cmd := exec.Command("gvt", "fetch", imp)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}
}

func addPackages(pkgs []*build.Package, path string, srcDir string) ([]*build.Package, error) {
	pkg, err := build.Import(path, srcDir, 0)
	if err == nil {
		pkgs = append(pkgs, pkg)
	}

	filesInfo, err := ioutil.ReadDir(pkg.Dir)
	if err != nil {
		return pkgs, err
	}

	for _, info := range filesInfo {
		fileName := info.Name()
		isHidden := fileName[0] == 46
		isVendor := fileName == "vendor"
		if !isHidden && !isVendor && info.IsDir() {
			pkgPath := fmt.Sprintf("%s/%s", path, fileName)
			pkgs, err = addPackages(pkgs, pkgPath, srcDir)
			if err != nil {
				return pkgs, err
			}
		}
	}
	return pkgs, nil
}

func addExternalImports(s *set, imports []string, pkgs []*build.Package) {
	for _, name := range imports {
		imp, err := build.Import(name, ".", 0)

		isInternal := false
		for _, pkg := range pkgs {
			if strings.HasPrefix(imp.ImportPath, pkg.ImportPath) {
				isInternal = true
			}
		}

		if err != nil || !imp.Goroot && !isInternal {
			s.add(imp.ImportPath)
		}
	}
}
