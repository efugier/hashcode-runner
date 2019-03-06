package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func copyFile(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("Couldn't open source file: %s", err)
	}
	defer inputFile.Close()
	outputFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("Couldn't open destination file: %s", err)
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, inputFile)
	if err != nil {
		return fmt.Errorf("Writing to output file failed: %s", err)
	}

	return nil
}

func moveFile(sourcePath, destPath string) error {
	err := copyFile(sourcePath, destPath)
	if err != nil {
		return fmt.Errorf("Failed to copy source file: %s", err)
	}
	err = os.Remove(sourcePath)
	if err != nil {
		return fmt.Errorf("Failed to remove original file: %s", err)
	}
	return nil
}

func swapFiles(file1Path, file2Path string) error {
	tmpFilePath := file2Path + ".tmp"
	errtmp := moveFile(file2Path, tmpFilePath)
	if errtmp != nil {
		log.Println("Failed to write", file1Path, "in a temporary file:", errtmp)
	}
	err := moveFile(file1Path, file2Path)
	if err != nil {
		return fmt.Errorf("Failed to move %s to %s: %s", file1Path, file2Path, err)
	}
	if errtmp == nil {
		err := moveFile(tmpFilePath, file1Path)
		if err != nil {
			return fmt.Errorf("Failed to move file2: %s", err)
		}
	}
	return nil
}
