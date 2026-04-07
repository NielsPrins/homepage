package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/google/uuid"
)

type Shortcut struct {
	ID       string `json:"id"`
	URL      string `json:"url"`
	Name     string `json:"name"`
	ImageURL string `json:"imageURL"`
	OrderNum int    `json:"order"`
}

func (s Shortcut) SafeImageURL() template.URL {
	return template.URL(s.ImageURL)
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

func AddShortcut(url string, name string, customImage string) (string, error) {
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
		ImageURL: imageUrl,
		OrderNum: len(shortcuts) + 1,
	}
	shortcuts = append(shortcuts, newShortcut)

	return id, saveShortcuts(shortcuts)
}

func EditShortcut(id, url, name, customImage string) error {
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

	sortShortcutsByOrder(shortcuts)
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

func ReorderShortcuts(ids []string) error {
	shortcuts, err := loadShortcuts()
	if err != nil {
		return fmt.Errorf("failed to load shortcuts: %w", err)
	}

	if len(ids) != len(shortcuts) {
		return fmt.Errorf("not enough shortcut ID's supplied to reorder")
	}

	shortcutByID := make(map[string]Shortcut, len(shortcuts))
	for _, shortcut := range shortcuts {
		shortcutByID[shortcut.ID] = shortcut
	}

	reordered := make(Shortcuts, 0, len(shortcuts))
	seen := make(map[string]bool, len(shortcuts))

	for i, id := range ids {
		shortcut, exists := shortcutByID[id]
		if !exists {
			return fmt.Errorf("no shortcut with ID: %s", id)
		}
		if seen[id] {
			return fmt.Errorf("duplicate id: %s", id)
		}

		shortcut.OrderNum = i + 1
		reordered = append(reordered, shortcut)
		seen[id] = true
	}

	return saveShortcuts(reordered)
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
