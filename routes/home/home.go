package home

import (
	_ "embed"
	"homepage/common"

	"github.com/gofiber/fiber/v2"
)

//go:embed home.gohtml
var htmlTemplate string

type SectionGroup struct {
	Name      string
	Shortcuts common.Shortcuts
}

type Data struct {
	SectionGroups []SectionGroup
}

func Handler(c *fiber.Ctx) error {
	shortcuts, _ := common.GetAllShortcuts()

	var sectionGroups []SectionGroup
	sectionIndex := map[string]int{}

	for _, shortcut := range shortcuts {
		sectionIndexValue, exists := sectionIndex[shortcut.Section]
		if !exists {
			sectionIndexValue = len(sectionGroups)
			sectionGroups = append(sectionGroups, SectionGroup{Name: shortcut.Section})
			sectionIndex[shortcut.Section] = sectionIndexValue
		}
		sectionGroups[sectionIndexValue].Shortcuts = append(sectionGroups[sectionIndexValue].Shortcuts, shortcut)
	}

	return common.RenderTemplate(c, htmlTemplate, Data{
		SectionGroups: sectionGroups,
	})
}
