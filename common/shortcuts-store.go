package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"os"
	"sort"
)

type Shortcut struct {
	ID       string `json:"id"`
	URL      string `json:"url"`
	Name     string `json:"name"`
	ImageURL string `json:"imageURL"`
	OrderNum int    `json:"order"`
}

type Shortcuts []Shortcut

const DataFilePath = "data/shortcuts.json"

var ErrShortcutNotFound = errors.New("shortcut not found")

func loadShortcuts() (Shortcuts, error) {
	_, err := os.Stat(DataFilePath)
	if os.IsNotExist(err) {
		err := createEmptyJSONFile(DataFilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to create file: %w", err)
		}
	}

	file, err := os.Open(DataFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var shortcuts Shortcuts
	data, err := os.ReadFile(DataFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	err = json.Unmarshal(data, &shortcuts)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return shortcuts, nil
}

func createEmptyJSONFile(filePath string) error {
	emptyData := []byte("[]")
	err := os.WriteFile(filePath, emptyData, 0644)
	if err != nil {
		return err
	}
	return nil
}

func saveShortcuts(shortcuts Shortcuts) error {
	data, err := json.MarshalIndent(shortcuts, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}
	return os.WriteFile(DataFilePath, data, 0644)
}

func AddShortcut(url, name, imageURL string) (string, error) {
	shortcuts, err := loadShortcuts()
	if err != nil {
		return "", fmt.Errorf("failed to load shortcuts: %w", err)
	}

	id := uuid.New().String()

	newShortcut := Shortcut{
		ID:       id,
		URL:      url,
		Name:     name,
		ImageURL: imageURL,
		OrderNum: len(shortcuts) + 1,
	}
	shortcuts = append(shortcuts, newShortcut)

	return id, saveShortcuts(shortcuts)
}

func EditShortcut(id, url, name, imageURL string, order int) error {
	shortcuts, err := loadShortcuts()
	if err != nil {
		return fmt.Errorf("failed to load shortcuts: %w", err)
	}

	for i, shortcut := range shortcuts {
		if shortcut.ID == id {
			shortcuts[i].URL = url
			shortcuts[i].Name = name
			shortcuts[i].ImageURL = imageURL
			shortcuts[i].OrderNum = order
			sortShortcutsByOrder(shortcuts)
			return saveShortcuts(shortcuts)
		}
	}
	return ErrShortcutNotFound
}

func RemoveShortcut(id string) error {
	shortcuts, err := loadShortcuts()
	if err != nil {
		return fmt.Errorf("failed to load shortcuts: %w", err)
	}

	for i, shortcut := range shortcuts {
		if shortcut.ID == id {
			shortcuts = append(shortcuts[:i], shortcuts[i+1:]...)
			reorderShortcutsAfterRemoval(shortcuts)
			return saveShortcuts(shortcuts)
		}
	}
	return ErrShortcutNotFound
}

func GetShortcut(id string) (*Shortcut, error) {
	shortcuts, err := loadShortcuts()
	if err != nil {
		return nil, fmt.Errorf("failed to load shortcuts: %w", err)
	}

	for _, shortcut := range shortcuts {
		if shortcut.ID == id {
			return &shortcut, nil
		}
	}
	return nil, ErrShortcutNotFound
}

func GetAllShortcuts() (Shortcuts, error) {
	shortcuts, err := loadShortcuts()
	if err != nil {
		return nil, fmt.Errorf("failed to load shortcuts: %w", err)
	}

	return shortcuts, nil
}

func sortShortcutsByOrder(shortcuts Shortcuts) {
	sort.Slice(shortcuts, func(i, j int) bool {
		return shortcuts[i].OrderNum < shortcuts[j].OrderNum
	})
}

func reorderShortcutsAfterRemoval(shortcuts Shortcuts) {
	for i := range shortcuts {
		shortcuts[i].OrderNum = i + 1
	}
}
