package main

import (
	"fmt"
	_ "os"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"time"
)

func formatFileSize(fileSizeBytes int64) string {
	const (
		KB = 1 << (10 * (iota + 1))
		MB
	)
	if fileSizeBytes >= MB {
		return fmt.Sprintf("%.2f MB", float64(fileSizeBytes)/MB)
	}
	return fmt.Sprintf("%.2f KB", float64(fileSizeBytes)/KB)
}

func generateDirectoryStructure(rootDir string, depth int) string {
	entries, err := ioutil.ReadDir(rootDir)
	if err != nil {
		return ""
	}

	var treeStructure string
	for _, entry := range entries {
		fullPath := filepath.Join(rootDir, entry.Name())
		if entry.IsDir() {
			symbols := []string{"🗂", "📁", "📂", "📄"}
			symbol := symbols[depth%len(symbols)]
			treeStructure += fmt.Sprintf("<li>%s %s/</li>", symbol, entry.Name())
			treeStructure += generateDirectoryStructure(fullPath, depth+1)
		} else {
			creationTime := entry.ModTime().Format("2006-01-02 15:04:05")
			fileSizeStr := formatFileSize(entry.Size())

			fileURL := "file://" + fullPath
			treeStructure += fmt.Sprintf("<li><a href='%s' target='_blank' style='text-decoration: none; color: white;'>%s</a> - Date: %s - Size: %s</li>",
				fileURL, entry.Name(), creationTime, fileSizeStr)
		}
	}
	return "<ul>" + treeStructure + "</ul>"
}

func main() {
	targetDirectory := "D:\\onedrive\\3图书\\01_编程书"
	directoryStructure := generateDirectoryStructure(targetDirectory, 0)

	htmlContent := fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<head>
		<title>文件目录查看器</title>
		<style>
			body {
				background-color: #333;
				color: white;
				margin: 0;
				padding: 0;
			}
			.container {
				width: 70%;
				margin: 0 auto;
				background-color: #333;
				padding: 20px;
			}
			a {
				text-decoration: none;
			}
			li {
				display: flex;
				justify-content: space-between;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<h1>文件目录查看器</h1>
			<p>目录：%s</p>
			%s
		</div>
	</body>
	</html>
	`, targetDirectory, directoryStructure)

	err := ioutil.WriteFile("index.html", []byte(htmlContent), 0644)
	if err != nil {
		fmt.Println("Error writing HTML file:", err)
		return
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		content, err := ioutil.ReadFile("index.html")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(content)
	})

	port := 2000
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Printf("Serving at http://localhost:%d\n", port)
	err = server.ListenAndServe()
	if err != nil {
		fmt.Println("Error starting the server:", err)
	}
}
