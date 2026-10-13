package downloader

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func (d *Downloader) delete(name string) error {
	if len(name) == 0 {
		return fmt.Errorf("missing file name")
	}

	fullPath := d.resolvePath(name)
	if _, err := os.Stat(fullPath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("failed to remove file %s: %w", name, err)
		}
		return fmt.Errorf("failed to stat file %s: %w", name, err)
	}

	if err := os.Remove(fullPath); err != nil {
		return fmt.Errorf("failed to remove file %s: %w", name, err)
	}

	return nil
}

func (d *Downloader) saveStream(r io.Reader, name string) (string, error) {
	if r == nil {
		return "", fmt.Errorf("image reader is nil")
	}

	if len(name) == 0 {
		name = "."
	}

	name = d.resolvePath(name)

	if err := os.MkdirAll(filepath.Dir(name), 0755); err != nil {
		return "", fmt.Errorf("failed to create directories for %s: %w", name, err)
	}

	file, err := os.Create(name)
	if err != nil {
		return "", fmt.Errorf("failed to create file %s: %w", name, err)
	}
	defer file.Close()

	if _, err := io.Copy(file, r); err != nil {
		return "", fmt.Errorf("failed to save file %s: %w", name, err)
	}

	return name, nil
}
