package main

import (
	"fmt"
	"github.com/kaleocheng/goldmark"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	_ "time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: md <filename>")
		return
	}

	filename := os.Args[1]
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	absInputFile, err := filepath.Abs(filename)
	if err != nil {
		fmt.Printf("Error resolving file path: %v\n", err)
		return
	}
	baseDir := filepath.Dir(absInputFile)
	baseURL := "file://" + filepath.ToSlash(baseDir) + "/"

	// time := time.Now()
	tmpFile, err := os.CreateTemp("/tmp", "md-preview-*.html")
	if err != nil {
		fmt.Printf("Error creating temp file: %v\n", err)
		return
	}
	defer tmpFile.Close()

	fmt.Fprintf(tmpFile, `<html><head><base href="%s"><style>
        body { font-family: -apple-system, Segoe UI, Helvetica, Arial, sans-serif; 
        max-width: 850px; margin: 50px auto; padding: 0 30px; line-height: 1.6; color: #333; }
        pre { background: #f6f8fa; padding: 16px; border-radius: 6px; overflow: auto; }
        code { font-family: monospace; background: #afb8c133; padding: 0.2em 0.4em; border-radius: 6px; }
        img { max-width: 100%%; display: block; margin: 0 auto; } </style></head><body>`, baseURL)

	goldmark.Convert(content, tmpFile)

	fmt.Fprint(tmpFile, `</body></html>`)

	absPath, _ := filepath.Abs(tmpFile.Name())
	fileURI := "file://" + absPath

	openBrowser(fileURI)
}

func openBrowser(url string) {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	default:
		cmd = "xdg-open"
		args = []string{url}
	}
	exec.Command(cmd, args...).Start()
}
