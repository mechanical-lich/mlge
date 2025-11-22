package resource

import (
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-fonts/liberation/liberationsansregular"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mechanical-lich/mlge/audio"
	"github.com/mechanical-lich/mlge/config"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var Textures map[string]*ebiten.Image
var Fonts map[string]font.Face
var Sounds map[string]audio.AudioResource

// LoadAssetsFromJSON - loads asset entries described in the JSON file.
// Support both map entries (name->path) and a "folders" key containing a list of directories to walk.
func LoadAssetsFromJSON(jsonPath string) error {
	// Parse into a generic map to detect both named assets and folder lists.
	raw := make(map[string]interface{})
	data, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("failed to unmarshal assets JSON: %w", err)
	}

	for key, val := range raw {
		// special-case the "folders" key as an array of folder paths
		if key == "folders" {
			arr, ok := val.([]interface{})
			if !ok {
				return fmt.Errorf("folders must be an array of strings")
			}
			for _, item := range arr {
				folderPath, ok := item.(string)
				if !ok {
					return fmt.Errorf("folder entries must be strings")
				}
				if err := processFolder(folderPath); err != nil {
					return fmt.Errorf("failed to process folder %s: %w", folderPath, err)
				}
			}
			continue
		}

		// otherwise treat entry as name -> path
		path, ok := val.(string)
		if !ok {
			// ignore non-string entries
			continue
		}
		if strings.Contains(strings.ToLower(path), ".ttf") {
			if err := LoadFont(key, path); err != nil {
				return fmt.Errorf("failed to load font %s: %w", key, err)
			}
		} else {
			if err := LoadImageAsTexture(key, path); err != nil {
				return fmt.Errorf("failed to load image %s: %w", key, err)
			}
		}
	}

	return nil
}

// processFolder - recursively walks the provided folder and loads assets.
// Names are generated as relative path with directory separators replaced by underscores
// and without the file extension, e.g. bullfrog/foot.png -> bullfrog_foot
func processFolder(folderPath string) error {
	slog.Debug("Loading assets from folder", slog.String("folder", folderPath))
	// normalize folder path
	folderPath = filepath.Clean(folderPath)

	info, err := os.Stat(folderPath)
	if err != nil {
		return fmt.Errorf("folder not found: %s", folderPath)
	}
	if !info.IsDir() {
		return fmt.Errorf("specified path is not a folder: %s", folderPath)
	}

	// Walk the folder recursively
	err = filepath.Walk(folderPath, func(path string, fi os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if fi.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))

		// Only consider font or image files; ignore other extensions
		switch ext {
		case ".ttf":
			name := assetNameFromFile(folderPath, path)
			if err := LoadFont(name, path); err != nil {
				return fmt.Errorf("failed loading font %s from %s: %w", name, path, err)
			}
		case ".png", ".jpg", ".jpeg", ".gif", ".webp":
			name := assetNameFromFile(folderPath, path)
			if err := LoadImageAsTexture(name, path); err != nil {
				return fmt.Errorf("failed loading image %s from %s: %w", name, path, err)
			}
			slog.Debug("Loaded texture", slog.String("name", name), slog.String("path", path))
		default:
			// not a recognized asset type; ignore
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

func assetNameFromFile(baseFolder, fullPath string) string {
	rel, err := filepath.Rel(baseFolder, fullPath)
	if err != nil {
		rel = fullPath
	}
	// Remove extension
	rel = strings.TrimSuffix(rel, filepath.Ext(rel))
	// Replace separators with underscore
	name := strings.ReplaceAll(rel, string(filepath.Separator), "_")
	// normalize to lowercase
	name = strings.ToLower(name)
	return name
}

// LoadImageAsTexture - Loads an image in the texture map with the given name and path.
func LoadImageAsTexture(name string, path string) error {
	if Textures == nil {
		log.Print("Initialize resource manager")
		Textures = make(map[string]*ebiten.Image)
	}
	img, err := LoadImage(path)
	if err != nil {
		return err
	}
	// Warn about duplicate names and overwrite to keep behavior simple (existing behavior uses a map)
	if _, exists := Textures[name]; exists {
		log.Printf("warning: texture name already exists, overwriting: %s", name)
	}
	Textures[name] = img
	return nil
}

// LoadImage - Loads an image from the specified path and returns it as an ebiten.Image
func LoadImage(path string) (*ebiten.Image, error) {
	imgFile, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening tileset " + path)
		return nil, errors.New("error opening tileset " + path)
	}

	img, _, err := image.Decode(imgFile)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return ebiten.NewImageFromImage(img), nil
}

// LoadFont - Loads a font into the font map with the given name and path.
func LoadFont(name string, path string) error {
	if Fonts == nil {
		log.Print("Initialize fonts")
		Fonts = make(map[string]font.Face)
	}
	raw := liberationsansregular.TTF
	tt, err := opentype.Parse(raw)
	if err != nil {
		return err
	}

	fontFace, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    16,
		DPI:     config.DPI,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return err
	}

	Fonts[name] = fontFace
	return nil
}
