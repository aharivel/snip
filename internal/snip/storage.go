package snip

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const DefaultDir = ".snip"
const CategoriesDir = "categories"

func DataDir() (string, error) {
	if env := strings.TrimSpace(os.Getenv("SNIP_DATA_DIR")); env != "" {
		return env, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, DefaultDir), nil
}

func CategoriesPath() (string, error) {
	dataDir, err := DataDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dataDir, CategoriesDir), nil
}

func CategoryFilePath(category string) (string, error) {
	categoriesPath, err := CategoriesPath()
	if err != nil {
		return "", err
	}
	return filepath.Join(categoriesPath, category+".md"), nil
}

func EnsureCategoriesDir() (string, error) {
	path, err := CategoriesPath()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(path, 0o755); err != nil {
		return "", err
	}
	return path, nil
}

func ListCategories() ([]string, error) {
	path, err := CategoriesPath()
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []string{}, nil
		}
		return nil, err
	}
	var categories []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, ".md") {
			categories = append(categories, strings.TrimSuffix(name, ".md"))
		}
	}
	sort.Strings(categories)
	return categories, nil
}

func ReadCategory(category string) (string, error) {
	path, err := CategoryFilePath(category)
	if err != nil {
		return "", err
	}
	bytes, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func CreateCategory(category string) (string, error) {
	if err := EnsureDataDir(); err != nil {
		return "", err
	}
	path, err := CategoryFilePath(category)
	if err != nil {
		return "", err
	}
	if _, err := os.Stat(path); err == nil {
		return "", ErrCategoryExists
	}
	if err := os.WriteFile(path, []byte(""), 0o644); err != nil {
		return "", err
	}
	return path, nil
}

func EnsureDataDir() error {
	dataDir, err := DataDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return err
	}
	_, err = EnsureCategoriesDir()
	return err
}
