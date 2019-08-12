package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func removeFilesFromList(list []os.FileInfo) []os.FileInfo {
	filtered := make([]os.FileInfo, 0, len(list))

	for _, fileInfo := range list {
		if !fileInfo.IsDir() {
			continue
		}
		filtered = append(filtered, fileInfo)
	}

	return filtered
}

func processPath(out io.Writer, root string, showFiles bool, globalPrefix string) {
	const (
		voidIdent = "	"
		ident     = "│	"
		leaf      = "├───"
		deadend   = "└───"
	)

	var elementPrefix, elementIdent, leafName, fileSize string

	list, _ := ioutil.ReadDir(root)

	if !showFiles {
		list = removeFilesFromList(list)
	}

	for idx, fileInfo := range list {
		if idx < len(list)-1 {
			elementPrefix = leaf
			elementIdent = ident
		} else {
			elementPrefix = deadend
			elementIdent = voidIdent
		}

		if !fileInfo.IsDir() {
			if fileInfo.Size() == 0 {
				fileSize = "empty"
			} else {
				fileSize = fmt.Sprintf("%db", fileInfo.Size())
			}
			leafName = fmt.Sprintf("%v (%v)", fileInfo.Name(), fileSize)
		} else {
			leafName = fileInfo.Name()
		}

		fmt.Fprintf(out, "%v%v\n", globalPrefix+elementPrefix, leafName)

		if fileInfo.IsDir() {
			processPath(
				out,
				root+string(os.PathSeparator)+fileInfo.Name(),
				showFiles,
				globalPrefix+elementIdent)
		}
	}

	return
}

func dirTree(out io.Writer, root string, showFiles bool) error {
	processPath(out, root, showFiles, "")

	return nil
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
