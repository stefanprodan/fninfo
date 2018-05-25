package function

import (
	"path/filepath"
	"os"
)

type Response struct {
	Hostname    string
	Secrets     []string
	Namespaces  []Namespace
	Environment []string
	Request     string
}

type Namespace struct {
	Name        string
	Pods        int
	Deployments int
	Services    int
}

func walkDir(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}
