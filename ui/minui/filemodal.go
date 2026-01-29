package minui

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// FileModal is a reusable file browser modal for Load/Save dialogs.
type FileModal struct {
	*Modal

	Mode        string // "load" or "save"
	DefaultPath string

	currentDir   string
	foldersList  *ListBox
	filesList    *ListBox
	pathInput    *TextInput
	okButton     *Button
	cancelButton *Button

	OnSelect func(path string)
	OnCancel func()
}

// NewFileModal creates a new FileModal. Mode should be "load" or "save".
func NewFileModal(id, title string, width, height int, mode string) *FileModal {
	fm := &FileModal{
		Modal: NewModal(id, title, width, height),
		Mode:  strings.ToLower(mode),
	}

	// Build UI elements
	folders := NewListBox(id+"_folders", []string{})
	folders.SetBounds(Rect{X: 0, Y: 0, Width: 270, Height: height - 120})

	files := NewListBox(id+"_files", []string{})
	files.SetBounds(Rect{X: 280, Y: 0, Width: 280, Height: height - 120})

	pathIn := NewTextInput(id+"_path", "")
	pathIn.SetBounds(Rect{X: 10, Y: height - 110, Width: width - 20, Height: 28})

	ok := NewButton(id+"_ok", "Load")
	ok.SetBounds(Rect{X: width - 220, Y: height - 70, Width: 100, Height: 32})

	cancel := NewButton(id+"_cancel", "Cancel")
	cancel.SetBounds(Rect{X: width - 110, Y: height - 70, Width: 100, Height: 32})

	// Mode-specific label
	if fm.Mode == "save" {
		ok.Text = "Save"
	}

	// Wire callbacks
	folders.OnSelect = func(i int, item string) {
		// navigate into selected folder
		if item == ".." {
			parent := filepath.Dir(fm.currentDir)
			fm.changeDirectory(parent)
			return
		}
		newDir := filepath.Join(fm.currentDir, item)
		fm.changeDirectory(newDir)
	}

	files.OnSelect = func(i int, item string) {
		full := filepath.Join(fm.currentDir, item)
		fm.pathInput.SetText(full)
	}

	pathIn.OnSubmit = func(text string) {
		// Treat Enter as OK
		full := text
		if fm.OnSelect != nil {
			fm.OnSelect(full)
		}
		fm.visible = false
	}

	ok.OnClick = func() {
		selected := fm.pathInput.GetText()
		if fm.Mode == "load" {
			// If path is directory, navigate into it
			if stat, err := os.Stat(selected); err == nil && stat.IsDir() {
				fm.changeDirectory(selected)
				return
			}
			// otherwise return file path
			if fm.OnSelect != nil {
				fm.OnSelect(selected)
			}
			fm.visible = false
			return
		}
		// save mode: return the path (no existence check)
		if fm.OnSelect != nil {
			fm.OnSelect(selected)
		}
		fm.visible = false
	}

	cancel.OnClick = func() {
		fm.visible = false
		if fm.OnCancel != nil {
			fm.OnCancel()
		}
	}

	// Attach elements to modal
	fm.AddChild(folders)
	fm.AddChild(files)
	fm.AddChild(pathIn)
	fm.AddChild(ok)
	fm.AddChild(cancel)

	// Store references
	fm.foldersList = folders
	fm.filesList = files
	fm.pathInput = pathIn
	fm.okButton = ok
	fm.cancelButton = cancel

	// Layout children to calculate dimensions
	folders.Layout()
	files.Layout()

	// Default start dir
	fm.currentDir = "."

	return fm
}

// SetDefaultPath sets the starting directory/path when the modal is first shown.
func (fm *FileModal) SetDefaultPath(path string) {
	if path == "" {
		return
	}
	fm.DefaultPath = path
	// if path is file, start in its dir
	if stat, err := os.Stat(path); err == nil && !stat.IsDir() {
		fm.currentDir = filepath.Dir(path)
		fm.pathInput.SetText(path)
	} else {
		fm.currentDir = path
		fm.pathInput.SetText("")
	}
	fm.changeDirectory(fm.currentDir)
}

// SetVisible overrides to refresh contents before displaying
func (fm *FileModal) SetVisible(v bool) {
	fm.visible = v
	if v {
		// Initialize currentDir from DefaultPath if provided
		if fm.DefaultPath != "" && fm.currentDir == "." {
			fm.SetDefaultPath(fm.DefaultPath)
		} else {
			fm.changeDirectory(fm.currentDir)
		}
	}
}

// changeDirectory updates the folder and file lists
func (fm *FileModal) changeDirectory(dir string) {
	if dir == "" {
		dir = "."
	}
	// resolve absolute path
	abs, err := filepath.Abs(dir)
	if err == nil {
		dir = abs
	}
	fm.currentDir = dir

	entries, err := os.ReadDir(dir)
	if err != nil {
		// put nothing
		fm.foldersList.SetItems([]string{})
		fm.filesList.SetItems([]string{})
		return
	}

	var dirs []string
	var files []string
	for _, e := range entries {
		if e.IsDir() {
			dirs = append(dirs, e.Name())
		} else {
			files = append(files, e.Name())
		}
	}
	sort.Strings(dirs)
	sort.Strings(files)

	// Prepend parent
	parent := filepath.Dir(dir)
	if parent != dir {
		dirs = append([]string{".."}, dirs...)
	}

	fm.foldersList.SetItems(dirs)
	fm.filesList.SetItems(files)

	// Update path input to current dir if empty
	if fm.pathInput.GetText() == "" {
		fm.pathInput.SetText(dir)
	}
}
