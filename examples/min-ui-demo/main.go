package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	minui "github.com/mechanical-lich/mlge/ui/min-ui"
)

const (
	screenWidth  = 1024
	screenHeight = 768
)

type Game struct {
	gui             *minui.GUI
	fileBrowser     *minui.Modal
	dropdownDemo    *minui.Panel
	toolbar         *minui.Panel
	dropdownVisible bool
}

func NewGame() *Game {
	g := &Game{
		gui: minui.NewGUI(),
	}

	g.setupToolbar()
	g.setupFileBrowserDemo()
	g.setupDropdownDemo()

	// Start with modal hidden
	g.fileBrowser.SetVisible(false)

	return g
}

func (g *Game) setupToolbar() {
	// Create toolbar at the top (use Panel, not HBox, since we're manually positioning)
	toolbar := minui.NewPanel("toolbar")
	toolbar.SetBounds(minui.Rect{X: 20, Y: 20, Width: 984, Height: 60})

	// Style the toolbar
	bgColor := color.Color(color.RGBA{240, 240, 245, 255})
	borderColor := color.Color(color.RGBA{120, 120, 130, 255})
	borderWidth := 2
	borderRadius := 5

	toolbar.GetStyle().BackgroundColor = &bgColor
	toolbar.GetStyle().BorderColor = &borderColor
	toolbar.GetStyle().BorderWidth = &borderWidth
	toolbar.GetStyle().BorderRadius = &borderRadius

	// Add title label
	titleLabel := minui.NewLabel("toolbarTitle", "Min-UI Demo Application")
	titleLabel.SetBounds(minui.Rect{X: 20, Y: 18, Width: 300, Height: 24})

	// Button to open file browser modal
	openBrowserBtn := minui.NewButton("openBrowser", "Open File Browser")
	openBrowserBtn.SetBounds(minui.Rect{X: 340, Y: 12, Width: 180, Height: 36})
	openBrowserBtn.OnClick = func() {
		g.fileBrowser.SetVisible(true)
		fmt.Println("Click")
	}

	// Button to toggle dropdown demo
	toggleDemoBtn := minui.NewButton("toggleDemo", "Toggle Demo Panel")
	toggleDemoBtn.SetBounds(minui.Rect{X: 530, Y: 12, Width: 180, Height: 36})
	toggleDemoBtn.OnClick = func() {
		if g.dropdownVisible {
			g.gui.RemoveElement(g.dropdownDemo)
			g.dropdownVisible = false
		} else {
			g.gui.AddElement(g.dropdownDemo)
			g.dropdownVisible = true
		}
	}

	// Button to create new modal
	newModalBtn := minui.NewButton("newModal", "Create Modal")
	newModalBtn.SetBounds(minui.Rect{X: 720, Y: 12, Width: 140, Height: 36})
	newModalBtn.OnClick = func() {
		g.createInfoModal()
	}

	toolbar.AddChild(titleLabel)
	toolbar.AddChild(openBrowserBtn)
	toolbar.AddChild(toggleDemoBtn)
	toolbar.AddChild(newModalBtn)

	g.toolbar = toolbar
	g.gui.AddElement(toolbar)
}

func (g *Game) createInfoModal() {
	infoModal := minui.NewModal("infoModal", "Information", 500, 350)
	infoModal.SetPosition(262, 209)

	// Title
	titleLabel := minui.NewLabel("infoTitle", "Welcome to Min-UI!")
	titleLabel.SetBounds(minui.Rect{X: 20, Y: 20, Width: 460, Height: 30})

	// Description
	descLabel := minui.NewLabel("infoDesc", "This is a vector-based UI library for MLGE.\nIt supports various elements:")
	descLabel.SetBounds(minui.Rect{X: 20, Y: 60, Width: 460, Height: 50})

	// Feature list
	featureList := minui.NewListBox("features", []string{
		"Buttons and Labels",
		"Text Input Fields",
		"Checkboxes",
		"ListBoxes with scrolling",
		"Panels and Layout Containers",
		"Modals with dragging support",
		"CSS-like styling system",
	})
	featureList.SetBounds(minui.Rect{X: 20, Y: 120, Width: 460, Height: 140})

	// Text input example
	inputLabel := minui.NewLabel("inputLabel", "Try typing here:")
	inputLabel.SetBounds(minui.Rect{X: 20, Y: 270, Width: 150, Height: 20})

	input := minui.NewTextInput("exampleInput", "")
	input.SetBounds(minui.Rect{X: 180, Y: 268, Width: 200, Height: 28})

	// Close button
	closeBtn := minui.NewButton("closeInfo", "Close")
	closeBtn.SetBounds(minui.Rect{X: 200, Y: 305, Width: 100, Height: 32})
	closeBtn.OnClick = func() {
		g.gui.RemoveModal(infoModal)
	}

	infoModal.AddChild(titleLabel)
	infoModal.AddChild(descLabel)
	infoModal.AddChild(featureList)
	infoModal.AddChild(inputLabel)
	infoModal.AddChild(input)
	infoModal.AddChild(closeBtn)

	g.gui.AddModal(infoModal)
}

func (g *Game) setupFileBrowserDemo() {
	// Create file browser modal
	modal := minui.NewModal("fileBrowser", "kiss_sdl example 1", 580, 400)
	modal.SetPosition(30, 40)

	// Create main container (use Panel for manual positioning)
	mainContainer := minui.NewPanel("mainContainer")
	mainContainer.SetBounds(minui.Rect{X: 10, Y: 10, Width: 560, Height: 300})

	// Folders panel (use Panel for manual positioning)
	foldersPanel := minui.NewPanel("foldersPanel")
	foldersPanel.SetBounds(minui.Rect{X: 0, Y: 0, Width: 270, Height: 300})

	foldersLabel := minui.NewLabel("foldersLabel", "Folders")
	foldersLabel.SetBounds(minui.Rect{X: 0, Y: 0, Width: 270, Height: 20})

	foldersList := minui.NewListBox("foldersList", []string{
		"../",
		"./",
		".git/",
	})
	foldersList.SetBounds(minui.Rect{X: 0, Y: 22, Width: 270, Height: 278})

	foldersPanel.AddChild(foldersLabel)
	foldersPanel.AddChild(foldersList)

	// Files panel (use Panel for manual positioning)
	filesPanel := minui.NewPanel("filesPanel")
	filesPanel.SetBounds(minui.Rect{X: 280, Y: 0, Width: 280, Height: 300})

	filesLabel := minui.NewLabel("filesLabel", "Files")
	filesLabel.SetBounds(minui.Rect{X: 0, Y: 0, Width: 280, Height: 20})

	filesList := minui.NewListBox("filesList", []string{
		"README.md",
		"kiss_LICENSE",
		"kiss_active.png",
		"kiss_bar.png",
		"kiss_down.png",
		"README.md",
		"kiss_LICENSE",
		"kiss_active.png",
		"kiss_bar.png",
		"kiss_down.png",
		"README.md",
		"kiss_LICENSE",
		"kiss_active.png",
		"kiss_bar.png",
		"kiss_down.png",
	})
	filesList.SetBounds(minui.Rect{X: 0, Y: 22, Width: 280, Height: 278})

	filesPanel.AddChild(filesLabel)
	filesPanel.AddChild(filesList)

	mainContainer.AddChild(foldersPanel)
	mainContainer.AddChild(filesPanel)
	modal.AddChild(mainContainer)

	// Create bottom panel with path and buttons
	pathLabel := minui.NewLabel("pathLabel", "/usr/local/projects/kiss_sdl/")
	pathLabel.SetBounds(minui.Rect{X: 10, Y: 320, Width: 300, Height: 20})

	pathInput := minui.NewTextInput("pathInput", "kiss")
	pathInput.SetText("kiss")
	pathInput.SetBounds(minui.Rect{X: 10, Y: 342, Width: 560, Height: 28})

	// Buttons
	okButton := minui.NewButton("okButton", "OK")
	okButton.SetBounds(minui.Rect{X: 390, Y: 378, Width: 80, Height: 32})
	okButton.OnClick = func() {
		fmt.Println("OK clicked:", pathInput.GetText())
	}

	cancelButton := minui.NewButton("cancelButton", "Cancel")
	cancelButton.SetBounds(minui.Rect{X: 480, Y: 378, Width: 80, Height: 32})
	cancelButton.OnClick = func() {
		modal.SetVisible(false)
	}

	modal.AddChild(pathLabel)
	modal.AddChild(pathInput)
	modal.AddChild(okButton)
	modal.AddChild(cancelButton)

	// Add selection modal for folder click
	foldersList.OnSelect = func(index int, item string) {
		selectionModal := minui.NewModal("selectionModal", "Info", 450, 200)
		selectionModal.SetPosition(95, 140)

		message := minui.NewLabel("message", fmt.Sprintf("The following path was selected:\n%s%s", "/usr/local/projects/kiss_sdl/", item))
		message.SetBounds(minui.Rect{X: 20, Y: 20, Width: 410, Height: 80})

		messageInput := minui.NewTextInput("messageInput", "")
		messageInput.SetBounds(minui.Rect{X: 20, Y: 100, Width: 410, Height: 28})

		okBtn := minui.NewButton("modalOkBtn", "OK")
		okBtn.SetBounds(minui.Rect{X: 180, Y: 140, Width: 80, Height: 32})
		okBtn.OnClick = func() {
			g.gui.RemoveModal(selectionModal)
		}

		selectionModal.AddChild(message)
		selectionModal.AddChild(messageInput)
		selectionModal.AddChild(okBtn)

		g.gui.AddModal(selectionModal)
	}

	g.fileBrowser = modal
	g.gui.AddModal(modal)
}

func (g *Game) setupDropdownDemo() {
	// Create a simple panel with checkboxes for the dropdown demo
	panel := minui.NewPanel("dropdownPanel")
	panel.SetBounds(minui.Rect{X: 200, Y: 100, Width: 400, Height: 300})

	// Set panel style
	bgColor := color.Color(color.RGBA{240, 240, 245, 255})
	borderColor := color.Color(color.RGBA{100, 100, 110, 255})
	borderWidth := 2

	panel.GetStyle().BackgroundColor = &bgColor
	panel.GetStyle().BorderColor = &borderColor
	panel.GetStyle().BorderWidth = &borderWidth

	// Add title
	title := minui.NewLabel("title", "Population")
	title.SetBounds(minui.Rect{X: 80, Y: 20, Width: 100, Height: 20})

	// Add checkboxes
	popCheckbox := minui.NewCheckbox("popCheckbox", "")
	popCheckbox.SetBounds(minui.Rect{X: 320, Y: 20, Width: 18, Height: 18})

	areaLabel := minui.NewLabel("areaLabel", "Area")
	areaLabel.SetBounds(minui.Rect{X: 80, Y: 50, Width: 100, Height: 20})

	areaCheckbox := minui.NewCheckbox("areaCheckbox", "")
	areaCheckbox.SetBounds(minui.Rect{X: 320, Y: 50, Width: 18, Height: 18})

	// Add dropdown-style list
	cities := minui.NewListBox("cities", []string{
		"Kansas City",
		"New York",
		"Orlando",
		"Philadelphia",
	})
	cities.SetBounds(minui.Rect{X: 80, Y: 90, Width: 200, Height: 100})
	cities.SelectedIndex = 0 // Select first item

	// Add slider/scroll control
	scrollLabel := minui.NewLabel("scrollLabel", "Label test")
	scrollLabel.SetBounds(minui.Rect{X: 60, Y: 200, Width: 300, Height: 20})

	// Add OK button
	okButton := minui.NewButton("okButton2", "OK")
	okButton.SetBounds(minui.Rect{X: 150, Y: 240, Width: 80, Height: 32})

	panel.AddChild(title)
	panel.AddChild(popCheckbox)
	panel.AddChild(areaLabel)
	panel.AddChild(areaCheckbox)
	panel.AddChild(cities)
	panel.AddChild(scrollLabel)
	panel.AddChild(okButton)

	g.dropdownDemo = panel
	// Don't add it by default - only show file browser
	// g.gui.AddElement(panel)
}

func (g *Game) Update() error {
	g.gui.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{200, 200, 210, 255})

	g.gui.Layout()
	g.gui.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Min-UI Example - File Browser")

	game := NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
