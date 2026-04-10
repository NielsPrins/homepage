package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type Shortcut struct {
	ID       string `json:"id"`
	URL      string `json:"url"`
	Name     string `json:"name"`
	Section  string `json:"section,omitempty"`
	ImageURL string `json:"imageURL"`
}

func (s Shortcut) SafeImageURL() template.URL {
	return template.URL(s.ImageURL)
}

type Shortcuts []Shortcut

var DataFilePath = func() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "data/shortcuts.json"
	}
	return filepath.Join(configDir, "Homepage", "shortcuts.json")
}()

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

	defer func() {
		_ = file.Close()
	}()

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
	err := ensureParentDir(filePath)
	if err != nil {
		return err
	}

	emptyData := []byte("[]")
	err = os.WriteFile(filePath, emptyData, 0644)
	if err != nil {
		return err
	}
	return nil
}

func ensureParentDir(filePath string) error {
	dir := filepath.Dir(filePath)
	if dir == "." || dir == "" {
		return nil
	}

	return os.MkdirAll(dir, 0755)
}

func saveShortcuts(shortcuts Shortcuts) error {
	data, err := json.MarshalIndent(shortcuts, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}
	return os.WriteFile(DataFilePath, data, 0644)
}

func ensureValidUrl(urlString string) (string, error) {
	if urlString == "" {
		return "", errors.New("url is required")
	}

	//goland:noinspection HttpUrlsUsage
	if !strings.HasPrefix(urlString, "http://") && !strings.HasPrefix(urlString, "https://") {
		urlString = "https://" + urlString
	}

	_, err := url.Parse(urlString)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}

	return urlString, nil
}

func AddShortcut(url string, name string, section string, customImage string) (string, error) {
	shortcuts, err := loadShortcuts()
	if err != nil {
		return "", fmt.Errorf("failed to load shortcuts: %w", err)
	}

	url, err = ensureValidUrl(url)
	if err != nil {
		return "", err
	}

	imageUrl := customImage
	if imageUrl == "" {
		imageUrl, _ = GetImageForShortcut(url)
	}
	if imageUrl == "" {
		imageUrl = "/public/no-favicon.ico"
	}

	id := uuid.New().String()

	newShortcut := Shortcut{
		ID:       id,
		URL:      url,
		Name:     name,
		Section:  section,
		ImageURL: imageUrl,
	}
	shortcuts = append(shortcuts, newShortcut)

	return id, saveShortcuts(shortcuts)
}

func EditShortcut(id, url, name, section, customImage string) error {
	if name == "" {
		return fmt.Errorf("name is required")
	}

	shortcuts, err := loadShortcuts()
	if err != nil {
		return fmt.Errorf("failed to load shortcuts: %w", err)
	}

	url, err = ensureValidUrl(url)
	if err != nil {
		return err
	}

	foundShortcut := false
	for i, shortcut := range shortcuts {
		if shortcut.ID == id {
			shortcuts[i].URL = url
			shortcuts[i].Name = name
			shortcuts[i].Section = section

			if customImage != "" {
				shortcuts[i].ImageURL = customImage
			} else if !strings.HasPrefix(shortcut.ImageURL, "data:image/") {
				imageUrl, _ := GetImageForShortcut(url)
				if imageUrl != "" {
					shortcuts[i].ImageURL = imageUrl
				} else {
					shortcuts[i].ImageURL = "/public/no-favicon.ico"
				}
			}

			foundShortcut = true
			break
		}
	}
	if !foundShortcut {
		return ErrShortcutNotFound
	}

	return saveShortcuts(shortcuts)
}

func RemoveShortcut(id string) error {
	shortcuts, err := loadShortcuts()
	if err != nil {
		return fmt.Errorf("failed to load shortcuts: %w", err)
	}

	for i, shortcut := range shortcuts {
		if shortcut.ID == id {
			shortcuts = append(shortcuts[:i], shortcuts[i+1:]...)
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

type OrderShortcutDto struct {
	ID      string `json:"id"`
	Section string `json:"section"`
}

func ReorderShortcuts(items []OrderShortcutDto) error {
	shortcuts, err := loadShortcuts()
	if err != nil {
		return fmt.Errorf("failed to load shortcuts: %w", err)
	}

	if len(items) != len(shortcuts) {
		return fmt.Errorf("not enough shortcut items supplied to reorder")
	}

	shortcutByID := make(map[string]Shortcut, len(shortcuts))
	for _, shortcut := range shortcuts {
		shortcutByID[shortcut.ID] = shortcut
	}

	reordered := make(Shortcuts, 0, len(shortcuts))
	seen := make(map[string]bool, len(shortcuts))

	for _, item := range items {
		shortcut, exists := shortcutByID[item.ID]
		if !exists {
			return fmt.Errorf("no shortcut with ID: %s", item.ID)
		}
		if seen[item.ID] {
			return fmt.Errorf("duplicate id: %s", item.ID)
		}

		shortcut.Section = item.Section
		reordered = append(reordered, shortcut)
		seen[item.ID] = true
	}

	return saveShortcuts(reordered)
}
