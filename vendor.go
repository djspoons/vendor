package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Package struct {
	Dir           string // directory containing package sources
	ImportPath    string // import path of package in dir
	ImportComment string // path in import comment on package statement
	Name          string // package name
	Doc           string // package documentation string
	Target        string // install path
	Shlib         string // the shared library that contains this package (only set when -linkshared)
	Goroot        bool   // is this package in the Go root?
	Standard      bool   // is this package part of the standard Go library?
	Stale         bool   // would 'go install' do anything for this package?
	Root          string // Go root or Go path dir containing this package

	// Source files
	GoFiles        []string // .go source files (excluding CgoFiles, TestGoFiles, XTestGoFiles)
	CgoFiles       []string // .go sources files that import "C"
	IgnoredGoFiles []string // .go sources ignored due to build constraints
	CFiles         []string // .c source files
	CXXFiles       []string // .cc, .cxx and .cpp source files
	MFiles         []string // .m source files
	HFiles         []string // .h, .hh, .hpp and .hxx source files
	SFiles         []string // .s source files
	SwigFiles      []string // .swig files
	SwigCXXFiles   []string // .swigcxx files
	SysoFiles      []string // .syso object files to add to archive

	// Cgo directives
	CgoCFLAGS    []string // cgo: flags for C compiler
	CgoCPPFLAGS  []string // cgo: flags for C preprocessor
	CgoCXXFLAGS  []string // cgo: flags for C++ compiler
	CgoLDFLAGS   []string // cgo: flags for linker
	CgoPkgConfig []string // cgo: pkg-config names

	// Dependency information
	Imports []string // import paths used by this package
	Deps    []string // all (recursively) imported dependencies

	// Error information
	Incomplete bool            // this package or a dependency has an error
	Error      *PackageError   // error loading package
	DepsErrors []*PackageError // errors loading dependencies

	TestGoFiles  []string // _test.go files in package
	TestImports  []string // imports from TestGoFiles
	XTestGoFiles []string // _test.go files outside package
	XTestImports []string // imports from XTestGoFiles
}

type PackageError struct {
	ImportStack []string // shortest path from package named on command line to this one
	Pos         string   // position of error (if present, file:line:col)
	Err         string   // the error itself
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("vendor: ")

	flag.Parse()

	ps, err := listPackages(flag.Args())
	if err != nil {
		log.Fatalf("error encountered listing packages: %v", err)
	}
	for _, p := range ps {
		if p.Error != nil {
			log.Printf("encountered package error: %v", p.Error.Err)
			continue
		}
		if !isStandardOrLocal(p) {
			if err := copyPackage(p); err != nil {
				log.Printf("error copying package %s: %v", p.ImportPath, err)
				continue
			}
		}
		deps, err := listPackages(p.Deps)
		if err != nil {
			log.Printf("error encountered listing packages: %v", err)
		}
		for _, d := range deps {
			if d.Error != nil {
				log.Printf("encountered package error: %v", d.Error.Err)
				continue
			}
			if isStandardOrLocal(d) {
				continue
			}
			if strings.Contains(d.ImportPath, "/vendor/") {
				log.Printf("cowardly refusing to vendor a vendor: %s", d.ImportPath)
				continue
			}
			if err := copyPackage(d); err != nil {
				log.Printf("error copying package %s: %v", d.ImportPath, err)
				continue
			}
		}
	}
}

func isStandardOrLocal(p *Package) bool {
	return p.Standard || isLocal(p)
}

var cwd string

func getwd() string {
	if cwd == "" {
		var err error
		cwd, err = os.Getwd()
		if err != nil {
			log.Fatalf("error getting current directory: %v", err)
		}
	}
	return cwd
}

func isLocal(d *Package) bool {
	return strings.HasPrefix(d.Dir, getwd())
}

// listPackages returns all packages in name
func listPackages(names []string) ([]*Package, error) {
	args := append([]string{"list", "-json"}, names...)
	cmd := exec.Command("go", args...)
	cmd.Stderr = os.Stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	d := json.NewDecoder(stdout)
	var ps []*Package
	for {
		p := new(Package)
		err := d.Decode(p)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		ps = append(ps, p)
	}
	if err := cmd.Wait(); err != nil {
		return nil, err
	}
	return ps, nil
}

func copyPackage(p *Package) error {
	vdir := filepath.Join("vendor", p.ImportPath)
	if err := os.MkdirAll(vdir, 0755); err != nil {
		return err
	}

	files := flatten(
		p.GoFiles,
		p.CgoFiles,
		p.IgnoredGoFiles,
		p.CFiles,
		p.CXXFiles,
		p.MFiles,
		p.HFiles,
		p.SFiles,
		p.SwigFiles,
		p.SwigCXXFiles,
		p.SysoFiles,
	)

	for _, fname := range files {
		if err := copyFile(
			filepath.Join(vdir, fname),
			filepath.Join(p.Dir, fname),
		); err != nil {
			return err
		}
	}
	return nil
}

func copyFile(dstpath, srcpath string) error {
	dst, err := os.Create(dstpath)
	if err != nil {
		return err
	}
	defer dst.Close()
	src, err := os.Open(srcpath)
	if err != nil {
		return err
	}
	defer src.Close()
	_, err = io.Copy(dst, src)
	return err
}

func flatten(sss ...[]string) (ss []string) {
	for _, v := range sss {
		ss = append(ss, v...)
	}
	return
}
