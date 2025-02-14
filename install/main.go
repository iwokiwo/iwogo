package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// CopyFile copies a single file from src to dst.
func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	return os.Chmod(dst, srcInfo.Mode())
}

// CopyDir copies a directory from src to dst.
func CopyDir(src, dst string) error {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("source folder does not exist: %v", err)
	}

	if !srcInfo.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	err = os.MkdirAll(dst, srcInfo.Mode())
	if err != nil {
		return err
	}

	err = filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		return CopyFile(path, destPath)
	})

	return err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run install/main.go <destination_folder>")
		return
	}

	// Correct the source folder path
	src := "iwogo/install/package" // Change this based on your project structure

	// Verify source exists
	if _, err := os.Stat(src); os.IsNotExist(err) {
		fmt.Println("Error: Source folder does not exist:", src)
		return
	}

	// Get the destination folder from the command-line argument
	dst := os.Args[1]

	// Perform the copy
	err := CopyDir(src, dst)
	if err != nil {
		fmt.Println("Error copying folder:", err)
	} else {
		fmt.Println("Configuration successfully copied to", dst)
	}
}
