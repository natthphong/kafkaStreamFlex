package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func UnzipFile(src string, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fPath := filepath.Join(dest, f.Name)

		// Prevent ZipSlip vulnerability
		if !strings.HasPrefix(fPath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", fPath)
		}

		if f.FileInfo().IsDir() {
			// Make directory and continue
			if err := os.MkdirAll(fPath, os.ModePerm); err != nil {
				return err
			}
			continue
		}

		// Make parent directories
		if err := os.MkdirAll(filepath.Dir(fPath), os.ModePerm); err != nil {
			return err
		}

		// Extract file
		inFile, err := f.Open()
		if err != nil {
			return err
		}
		outFile, err := os.Create(fPath)
		if err != nil {
			inFile.Close()
			return err
		}

		_, err = io.Copy(outFile, inFile)

		// Always close files
		outFile.Close()
		inFile.Close()

		if err != nil {
			return err
		}
	}
	return nil
}
