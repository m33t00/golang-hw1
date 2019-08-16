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

func getNodeName(fileInfo os.FileInfo) string {
	fileSize := "empty"

	if fileInfo.IsDir() {
		return fileInfo.Name()
	}

	if fileInfo.Size() != 0 {
		fileSize = fmt.Sprintf("%db", fileInfo.Size())
	}

	return fmt.Sprintf("%v (%v)", fileInfo.Name(), fileSize)
}

func processPath(out io.Writer, root string, showFiles bool, linePrefix string) {
	const (
		voidIdent = "\t"
		ident     = "│\t"
		leaf      = "├───"
		deadend   = "└───"
	)

	var elementPrefix, nextLineIdentation string

	list, _ := ioutil.ReadDir(root)

	if !showFiles {
		list = removeFilesFromList(list)
	}

	for idx, fileInfo := range list {
		if idx < len(list)-1 {
			elementPrefix, nextLineIdentation = leaf, ident
		} else {
			elementPrefix, nextLineIdentation = deadend, voidIdent
		}

		fmt.Fprintf(out, "%v%v\n", linePrefix+elementPrefix, getNodeName(fileInfo))

		if fileInfo.IsDir() {
			processPath(
				out,
				root+string(os.PathSeparator)+fileInfo.Name(),
				showFiles,
				linePrefix+nextLineIdentation)
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
