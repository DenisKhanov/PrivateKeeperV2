package lib

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const (
	perm = 0o755
)

func SaveToFile(path string, data string) error {
	path = filepath.FromSlash(path)

	dir := filepath.Dir(path)
	fmt.Println(dir)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, perm)
		if err != nil {
			return fmt.Errorf("unable to create directory: %s %w", dir, err)
		}
	}

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, perm)
	if err != nil {
		return fmt.Errorf("unable to create file: %s %w", path, err)
	}
	defer file.Close()

	_, err = file.WriteString(data)
	if err != nil {
		return fmt.Errorf("unable to write to file: %s %w", path, err)
	}

	return nil
}

func SaveBinaryToFile(path string, data []byte) error {
	path = filepath.FromSlash(path)

	dir := filepath.Dir(path)
	fmt.Println(dir)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, perm)
		if err != nil {
			return fmt.Errorf("unable to create directory: %s %w", dir, err)
		}
	}

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, perm)
	if err != nil {
		return fmt.Errorf("unable to create file: %s %w", path, err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("unable to write to file: %s %w", path, err)
	}

	return nil
}

func LoadFromFile(path string) ([]byte, error) {
	path = filepath.FromSlash(path)

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open file: %s %w", path, err)
	}

	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("read file: %s %w", path, err)
	}

	return data, nil
}
