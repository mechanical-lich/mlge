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
	progressBar     *minui.ProgressBar
}

func NewGame() *Game {
	g := &Game{
		gui: minui.NewGUI(),
	}

	g.setupToolbar()
	g.setupFileBrowserDemo()
	g.setupDropdownDemo()
	g.setupProgressDemo()

	// Start with modal hidden
	g.fileBrowser.SetVisible(false)

	return g
}

func (g *Game) setupProgressDemo() {
	// Floating progress bar near the bottom-left of the window
	pb := minui.NewProgressBar("mainProgress")
	pb.SetBounds(minui.Rect{X: 20, Y: 700, Width: 400, Height: 20})
	pb.ShowPercent = true

	incBtn := minui.NewButton("incProgress", "Inc")
	incBtn.SetBounds(minui.Rect{X: 430, Y: 696, Width: 60, Height: 28})
	incBtn.OnClick = func() {
		v := pb.GetValue()
		pb.SetValue(v + 0.1)
	}

	resetBtn := minui.NewButton("resetProgress", "Reset")
	resetBtn.SetBounds(minui.Rect{X: 500, Y: 696, Width: 60, Height: 28})
	resetBtn.OnClick = func() {
		pb.SetValue(0)
	}

	g.gui.AddElement(pb)
	g.gui.AddElement(incBtn)
	g.gui.AddElement(resetBtn)
	g.progressBar = pb
}

func (g *Game) setupToolbar() {
	// Create toolbar at the top (use Panel, not HBox, since we're manually positioning)
	toolbar := minui.NewPanel("toolbar")
	toolbar.SetBounds(minui.Rect{X: 20, Y: 20, Width: 984, Height: 60})

	// Style the toolbar
	bgColor := color.Color(color.RGBA{240, 240, 245, 255})
	borderColor := color.Color(color.RGBA{120, 120, 130, 255})
	borderWidth := 1
	borderRadius := 1

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

	// Button to show size constraint demo
	sizeConstraintBtn := minui.NewButton("sizeConstraint", "Size Demo")
	sizeConstraintBtn.SetBounds(minui.Rect{X: 870, Y: 12, Width: 90, Height: 36})
	sizeConstraintBtn.OnClick = func() {
		g.createSizeConstraintDemo()
	}

	// Button to show radio button demo
	radioBtn := minui.NewButton("radioDemo", "Radio")
	radioBtn.SetBounds(minui.Rect{X: 970, Y: 12, Width: 70, Height: 36})
	radioBtn.OnClick = func() {
		g.createRadioButtonDemo()
	}

	toolbar.AddChild(titleLabel)
	toolbar.AddChild(openBrowserBtn)
	toolbar.AddChild(toggleDemoBtn)
	toolbar.AddChild(newModalBtn)
	toolbar.AddChild(sizeConstraintBtn)
	toolbar.AddChild(radioBtn)

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

func (g *Game) createSizeConstraintDemo() {
	demoModal := minui.NewModal("sizeDemo", "Size Constraint Demo", 600, 500)
	demoModal.SetPosition(212, 134)

	// Title
	title := minui.NewLabel("demoTitle", "Min/Max Size Constraints")
	title.SetBounds(minui.Rect{X: 20, Y: 20, Width: 560, Height: 30})

	// Example 1: Button with MinWidth constraint
	example1Label := minui.NewLabel("ex1Label", "Button with MinWidth=150 (text is short)")
	example1Label.SetBounds(minui.Rect{X: 20, Y: 60, Width: 350, Height: 20})

	shortTextBtn := minui.NewButton("shortBtn", "Hi")
	shortTextBtn.SetBounds(minui.Rect{X: 20, Y: 85, Width: 0, Height: 0}) // Will be calculated
	minWidth := 150
	shortTextBtn.GetStyle().MinWidth = &minWidth

	// Example 2: Button with MaxWidth constraint
	example2Label := minui.NewLabel("ex2Label", "Button with MaxWidth=200 (text is very long)")
	example2Label.SetBounds(minui.Rect{X: 20, Y: 135, Width: 400, Height: 20})

	longTextBtn := minui.NewButton("longBtn", "This is a very long button text that would normally be wide")
	longTextBtn.SetBounds(minui.Rect{X: 20, Y: 160, Width: 0, Height: 0}) // Will be calculated
	maxWidth := 200
	longTextBtn.GetStyle().MaxWidth = &maxWidth

	// Example 3: Label with MinHeight
	example3Label := minui.NewLabel("ex3Label", "Label with MinHeight=80")
	example3Label.SetBounds(minui.Rect{X: 20, Y: 210, Width: 300, Height: 20})

	tallLabel := minui.NewLabel("tallLabel", "Short text")
	tallLabel.SetBounds(minui.Rect{X: 20, Y: 235, Width: 0, Height: 0}) // Will be calculated
	minHeight := 80
	bgColor := color.Color(color.RGBA{220, 220, 230, 255})
	borderColor := color.Color(color.RGBA{100, 100, 120, 255})
	borderWidth := 1
	tallLabel.GetStyle().MinHeight = &minHeight
	tallLabel.GetStyle().BackgroundColor = &bgColor
	tallLabel.GetStyle().BorderColor = &borderColor
	tallLabel.GetStyle().BorderWidth = &borderWidth

	// Example 4: TextInput with both constraints
	example4Label := minui.NewLabel("ex4Label", "TextInput with MinWidth=250, MaxWidth=400")
	example4Label.SetBounds(minui.Rect{X: 20, Y: 330, Width: 400, Height: 20})

	constrainedInput := minui.NewTextInput("constrainedInput", "Type here...")
	constrainedInput.SetBounds(minui.Rect{X: 20, Y: 355, Width: 300, Height: 28})
	inputMinWidth := 250
	inputMaxWidth := 400
	constrainedInput.GetStyle().MinWidth = &inputMinWidth
	constrainedInput.GetStyle().MaxWidth = &inputMaxWidth

	// Info text
	infoText := minui.NewLabel("infoText", "Try resizing these elements - they'll respect\ntheir min/max constraints during layout!")
	infoText.SetBounds(minui.Rect{X: 20, Y: 395, Width: 560, Height: 40})

	// Close button
	closeBtn := minui.NewButton("closeSizeDemo", "Close")
	closeBtn.SetBounds(minui.Rect{X: 250, Y: 445, Width: 100, Height: 32})
	closeBtn.OnClick = func() {
		g.gui.RemoveModal(demoModal)
	}

	demoModal.AddChild(title)
	demoModal.AddChild(example1Label)
	demoModal.AddChild(shortTextBtn)
	demoModal.AddChild(example2Label)
	demoModal.AddChild(longTextBtn)
	demoModal.AddChild(example3Label)
	demoModal.AddChild(tallLabel)
	demoModal.AddChild(example4Label)
	demoModal.AddChild(constrainedInput)
	demoModal.AddChild(infoText)
	demoModal.AddChild(closeBtn)

	g.gui.AddModal(demoModal)
}

func (g *Game) createRadioButtonDemo() {
	radioModal := minui.NewModal("radioDemo", "Radio Button Demo", 500, 450)
	radioModal.SetPosition(262, 159)

	// Title
	title := minui.NewLabel("radioTitle", "Radio Buttons & Radio Groups")
	title.SetBounds(minui.Rect{X: 20, Y: 20, Width: 460, Height: 30})

	// Section 1: Radio Group
	groupLabel := minui.NewLabel("groupLabel", "Radio Group (select one):")
	groupLabel.SetBounds(minui.Rect{X: 20, Y: 60, Width: 300, Height: 20})

	// Create radio group
	colorGroup := minui.NewRadioGroup("colorGroup")

	// Result label
	selectedLabel := minui.NewLabel("selectedLabel", "Selected: None")
	selectedLabel.SetBounds(minui.Rect{X: 20, Y: 90, Width: 460, Height: 20})

	// Radio buttons for the group
	redRadio := minui.NewRadioButton("redRadio", "Red")
	redRadio.SetBounds(minui.Rect{X: 20, Y: 115, Width: 20, Height: 20})
	redLabel := minui.NewLabel("redLabel", "Red")
	redLabel.SetBounds(minui.Rect{X: 45, Y: 115, Width: 100, Height: 20})

	greenRadio := minui.NewRadioButton("greenRadio", "Green")
	greenRadio.SetBounds(minui.Rect{X: 20, Y: 145, Width: 20, Height: 20})
	greenLabel := minui.NewLabel("greenLabel", "Green")
	greenLabel.SetBounds(minui.Rect{X: 45, Y: 145, Width: 100, Height: 20})

	blueRadio := minui.NewRadioButton("blueRadio", "Blue")
	blueRadio.SetBounds(minui.Rect{X: 20, Y: 175, Width: 20, Height: 20})
	blueLabel := minui.NewLabel("blueLabel", "Blue")
	blueLabel.SetBounds(minui.Rect{X: 45, Y: 175, Width: 100, Height: 20})

	// Add buttons to group
	colorGroup.AddButton(redRadio)
	colorGroup.AddButton(greenRadio)
	colorGroup.AddButton(blueRadio)

	// Set callback
	colorGroup.OnSelectionChange = func(selectedID string, selectedButton *minui.RadioButton) {
		selectedLabel.Text = fmt.Sprintf("Selected: %s", selectedButton.Label)
	}

	// Section 2: Standalone radio buttons (like round checkboxes)
	standaloneLabel := minui.NewLabel("standaloneLabel", "Standalone Radio Buttons (toggle independently):")
	standaloneLabel.SetBounds(minui.Rect{X: 20, Y: 220, Width: 460, Height: 20})

	option1Radio := minui.NewRadioButton("option1", "Option 1")
	option1Radio.SetBounds(minui.Rect{X: 20, Y: 250, Width: 20, Height: 20})
	option1Label := minui.NewLabel("option1Label", "Option 1")
	option1Label.SetBounds(minui.Rect{X: 45, Y: 250, Width: 100, Height: 20})

	option2Radio := minui.NewRadioButton("option2", "Option 2")
	option2Radio.SetBounds(minui.Rect{X: 20, Y: 280, Width: 20, Height: 20})
	option2Label := minui.NewLabel("option2Label", "Option 2")
	option2Label.SetBounds(minui.Rect{X: 45, Y: 280, Width: 100, Height: 20})

	option3Radio := minui.NewRadioButton("option3", "Option 3")
	option3Radio.SetBounds(minui.Rect{X: 20, Y: 310, Width: 20, Height: 20})
	option3Label := minui.NewLabel("option3Label", "Option 3")
	option3Label.SetBounds(minui.Rect{X: 45, Y: 310, Width: 100, Height: 20})

	// Info text
	infoText := minui.NewLabel("radioInfo", "Radio buttons in groups enforce single selection.\nStandalone ones toggle like round checkboxes.")
	infoText.SetBounds(minui.Rect{X: 20, Y: 350, Width: 460, Height: 40})

	// Close button
	closeBtn := minui.NewButton("closeRadioDemo", "Close")
	closeBtn.SetBounds(minui.Rect{X: 200, Y: 400, Width: 100, Height: 32})
	closeBtn.OnClick = func() {
		g.gui.RemoveModal(radioModal)
	}

	// Add all elements to modal
	radioModal.AddChild(title)
	radioModal.AddChild(groupLabel)
	radioModal.AddChild(selectedLabel)
	radioModal.AddChild(redRadio)
	radioModal.AddChild(redLabel)
	radioModal.AddChild(greenRadio)
	radioModal.AddChild(greenLabel)
	radioModal.AddChild(blueRadio)
	radioModal.AddChild(blueLabel)
	radioModal.AddChild(standaloneLabel)
	radioModal.AddChild(option1Radio)
	radioModal.AddChild(option1Label)
	radioModal.AddChild(option2Radio)
	radioModal.AddChild(option2Label)
	radioModal.AddChild(option3Radio)
	radioModal.AddChild(option3Label)
	radioModal.AddChild(infoText)
	radioModal.AddChild(closeBtn)

	g.gui.AddModal(radioModal)
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
