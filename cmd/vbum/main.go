package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/coreos/go-semver/semver"
)

func main() {
	lauge := log.New(os.Stderr, "", 0)
	flag.CommandLine.SetOutput(os.Stderr)
	var verbose bool
	var vers string
	var tgt string
	var pre string
	var meta string
	flag.BoolVar(&verbose, "v", false, "verbose output")
	flag.StringVar(&vers, "ver", "", "semantic versions string to manage")
	flag.StringVar(&tgt, "tgt", "", "component to bump: major, minor, or patch, optional")
	flag.StringVar(&pre, "pre-release", "", "pre-release tag, if applicable")
	flag.StringVar(&meta, "meta", "", "meta tag, optional, gitsha is a good example of something to use here")
	flag.Parse()
	if vers == "" {
		flag.Usage()
		return
	}
	if vers[0] == 'v' {
		vers = vers[1:]
	}
	if verbose {
		lauge.Printf("checking %s\n", vers)
	}
	v, err := semver.NewVersion(vers)
	if err != nil {
		lauge.Printf("Error: %s\n\n", err)
		flag.Usage()
		return
	}
	if verbose {
		lauge.Printf("major version number is %d\n", v.Major)
		lauge.Printf("minor version number is %d\n", v.Minor)
		lauge.Printf("patch version number is %d\n", v.Patch)
		lauge.Printf("pre version number is %s\n", v.PreRelease)
		lauge.Printf("meta version number is %s\n", v.Metadata)
	}
	switch tgt {
	case "major":
		v.BumpMajor()
	case "minor":
		v.BumpMajor()
	case "patch":
		v.BumpMajor()
	}
	if pre != "" {
		v.PreRelease = semver.PreRelease(pre)
	}
	v.Metadata = meta
	fmt.Printf("v%s", v.String())
}
