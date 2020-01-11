package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	/*
		gazelle: /home/cmcfarland/Code/highlander-bazel-setup/highlander/go/src/paperless/vendor/github.com/3d0c/gmf/frame_go112.go: error reading go file: /home/cmcfarland/Code/highlander-bazel-setup/highlander/go/src/paperless/vendor/github.com/3d0c/gmf/frame_go112.go: pkg-config not supported: #cgo pkg-config: libavcodec libavutil
		gazelle: finding module path for import paperless/service/entitlements: exit status 1: can't load package: package paperless/service/entitlements: malformed module path "paperless/service/entitlements": missing dot in first path element
	*/
	scnr := bufio.NewScanner(os.Stdin)
	for scnr.Scan() {
		line := scnr.Text()
		found, parsed := resolveAnnotation(line)
		if found {
			fmt.Println(parsed)
		}
	}
}

func resolveAnnotation(line string) (bool, string) {
	if strings.HasPrefix(line, "gazelle: finding module path for import") {
		splitted := strings.Split(line, " ")
		if len(splitted) < 7 {
			return false, ""
		}
		// get lib path like "paperless/service/entitlements:" and remove the ":"
		lib := strings.Replace(splitted[6], ":", "", 1)
		// template the annotation
		lib = fmt.Sprintf("# gazelle:resolve go %s //go/src/%s:go_default_library", lib, lib)
		return true, lib
	}
	return false, ""
}
