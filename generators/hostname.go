package generators

import (
	"fmt"
	"os"
	"path/filepath"

	lxdShared "github.com/canonical/lxd/shared"
	"github.com/canonical/lxd/shared/api"

	"github.com/canonical/lxd-imagebuilder/image"
	"github.com/canonical/lxd-imagebuilder/shared"
)

type hostname struct {
	common
}

// RunLXC creates a hostname template.
func (g *hostname) RunLXC(img *image.LXCImage, target shared.DefinitionTargetLXC) error {
	// Skip if the file doesn't exist
	if !lxdShared.PathExists(filepath.Join(g.sourceDir, g.defFile.Path)) {
		return nil
	}

	// Create new hostname file
	file, err := os.Create(filepath.Join(g.sourceDir, g.defFile.Path))
	if err != nil {
		return fmt.Errorf("Failed to create file %q: %w", filepath.Join(g.sourceDir, g.defFile.Path), err)
	}

	defer file.Close()

	// Write LXC specific string to the hostname file
	_, err = file.WriteString("LXC_NAME\n")
	if err != nil {
		return fmt.Errorf("Failed to write to hostname file: %w", err)
	}

	// Add hostname path to LXC's templates file
	err = img.AddTemplate(g.defFile.Path)
	if err != nil {
		return fmt.Errorf("Failed to add template: %w", err)
	}

	return nil
}

// RunLXD creates a hostname template.
func (g *hostname) RunLXD(img *image.LXDImage, target shared.DefinitionTargetLXD) error {
	// Skip if the file doesn't exist
	if !lxdShared.PathExists(filepath.Join(g.sourceDir, g.defFile.Path)) {
		return nil
	}

	templateDir := filepath.Join(g.cacheDir, "templates")

	err := os.MkdirAll(templateDir, 0755)
	if err != nil {
		return fmt.Errorf("Failed to create directory %q: %w", templateDir, err)
	}

	file, err := os.Create(filepath.Join(templateDir, "hostname.tpl"))
	if err != nil {
		return fmt.Errorf("Failed to create file %q: %w", filepath.Join(templateDir, "hostname.tpl"), err)
	}

	defer file.Close()

	_, err = file.WriteString("{{ container.name }}\n")
	if err != nil {
		return fmt.Errorf("Failed to write to hostname file: %w", err)
	}

	// Add to LXD templates
	img.Metadata.Templates[g.defFile.Path] = &api.ImageMetadataTemplate{
		Template:   "hostname.tpl",
		Properties: g.defFile.Template.Properties,
		When:       g.defFile.Template.When,
	}

	if len(g.defFile.Template.When) == 0 {
		img.Metadata.Templates[g.defFile.Path].When = []string{
			"create",
			"copy",
		}
	}

	return nil
}

// Run does nothing.
func (g *hostname) Run() error {
	return nil
}
