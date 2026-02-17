package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mechanical-lich/mlge/event"
	"github.com/mechanical-lich/mlge/resource"
	ui "github.com/mechanical-lich/mlge/ui/themed"
	containers "github.com/mechanical-lich/mlge/ui/themed/containers"
	elements "github.com/mechanical-lich/mlge/ui/themed/elements"
	theming "github.com/mechanical-lich/mlge/ui/themed/theming"
	validation "github.com/mechanical-lich/mlge/ui/themed/validation"
	views "github.com/mechanical-lich/mlge/ui/themed/views"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

type Game struct {
	gui          *ui.GUI
	eventManager *event.EventManager
}

// UIEventLogger implements EventListener to log UI events
type UIEventLogger struct{}

func (l *UIEventLogger) HandleEvent(data event.EventData) error {
	switch e := data.(type) {
	case ui.ClickEventData:
		log.Printf("Click event from: %s", e.SourceName)
	case ui.ValueChangeEventData:
		log.Printf("Value change in: %s (new value: %v)", e.SourceName, e.NewValue)
	case ui.FocusEventData:
		log.Printf("Focus event: %s", e.SourceName)
	case ui.BlurEventData:
		log.Printf("Blur event: %s", e.SourceName)
	case ui.ModalOpenEventData:
		log.Printf("Modal opened: %s", e.ModalName)
	case ui.ModalCloseEventData:
		log.Printf("Modal closed: %s", e.ModalName)
	}
	return nil
}

// ExampleView demonstrates all the new UI features
type ExampleView struct {
	views.GUIViewBase
	eventManager *event.EventManager

	// Form elements
	nameInput      *elements.InputField
	emailInput     *elements.InputField
	passwordInput  *elements.InputField
	ageSlider      *elements.Slider
	volumeSlider   *elements.Slider
	agreeCheckbox  *elements.Checkbox
	notifyCheckbox *elements.Checkbox
	themeDropdown  *elements.Dropdown
	submitButton   *elements.Button
	clearButton    *elements.Button
	toggleModal    *elements.Button

	// Demo modal
	demoModal *containers.Modal

	// Status label
	statusLabel *elements.Label
	eventLabel  *elements.Label
}

func NewExampleView(eventMgr *event.EventManager) *ExampleView {
	resource.LoadAssetsFromJSON("assets.json")
	view := &ExampleView{
		eventManager: eventMgr,
	}

	// Create name input with validation
	view.nameInput = elements.NewInputField("nameInput", 20, 50, 200, 50)
	view.nameInput.Placeholder = "Enter your name"
	view.nameInput.Validator = validation.Combine(
		validation.Required("Name"),
		validation.MinLength("Name", 2),
	)
	view.nameInput.OnChange = func(value string) {
		eventMgr.SendEvent(ui.ValueChangeEventData{
			SourceName: "nameInput",
			OldValue:   "",
			NewValue:   value,
		})
	}

	// Create email input with validation
	view.emailInput = elements.NewInputField("emailInput", 20, 100, 200, 100)
	view.emailInput.Placeholder = "your.email@example.com"
	view.emailInput.Validator = validation.Combine(
		validation.Required("Email"),
		validation.Email("Email"),
	)

	// Create password input
	view.passwordInput = elements.NewInputField("passwordInput", 20, 150, 200, 50)
	view.passwordInput.Placeholder = "Password"
	view.passwordInput.IsPassword = true
	view.passwordInput.Validator = validation.MinLength("Password", 6)

	// Create age slider
	view.ageSlider = elements.NewSlider("ageSlider", 20, 220, 200, "Age", 0, 100, 25, 1)
	view.ageSlider.OnChanged = func(value float64) {
		eventMgr.SendEvent(ui.ValueChangeEventData{
			SourceName: "ageSlider",
			OldValue:   0.0,
			NewValue:   value,
		})
	}

	// Create volume slider
	view.volumeSlider = elements.NewSlider("volumeSlider", 20, 270, 200, "Volume", 0, 100, 50, 5)

	// Create checkboxes
	view.agreeCheckbox = elements.NewCheckbox("agreeCheckbox", 20, 320, "I agree to terms", false)
	view.agreeCheckbox.OnChanged = func(checked bool) {
		eventMgr.SendEvent(ui.ValueChangeEventData{
			SourceName: "agreeCheckbox",
			OldValue:   !checked,
			NewValue:   checked,
		})
	}

	view.notifyCheckbox = elements.NewCheckbox("notifyCheckbox", 20, 350, "Send notifications", true)

	// Create dropdown
	themeOptions := []string{"Default Theme", "Dark Theme", "Light Theme"}
	view.themeDropdown = elements.NewDropdown("themeDropdown", 20, 390, 200, themeOptions, 0)
	view.themeDropdown.OnChanged = func(index int, value string) {
		eventMgr.SendEvent(ui.ValueChangeEventData{
			SourceName: "themeDropdown",
			OldValue:   "",
			NewValue:   value,
		})
	}

	// Create buttons
	view.submitButton = elements.NewButton("submitButton", 20, 440, "Submit Form", "Submit the form")
	view.submitButton.OnClicked = func() {
		eventMgr.SendEvent(ui.ClickEventData{
			SourceName: "submitButton",
			Data: map[string]interface{}{
				"action": "submit",
			},
		})
	}

	view.clearButton = elements.NewButton("clearButton", 140, 440, "Clear", "Clear all fields")
	view.clearButton.OnClicked = func() {
		view.nameInput.SetValue("")
		view.emailInput.SetValue("")
		view.passwordInput.SetValue("")
		view.ageSlider.SetValue(25)
		view.volumeSlider.SetValue(50)
		view.agreeCheckbox.SetChecked(false)
		view.notifyCheckbox.SetChecked(true)
		view.themeDropdown.SetSelectedIndex(0)
		eventMgr.SendEvent(ui.ClickEventData{
			SourceName: "clearButton",
			Data:       nil,
		})
	}

	view.toggleModal = elements.NewButton("toggleModal", 20, 490, "Open Modal", "Show demo modal")
	view.toggleModal.OnClicked = func() {
		if view.demoModal != nil {
			view.demoModal.OpenModal()
		}
	}

	// Create status labels
	view.statusLabel = elements.NewLabel("statusLabel", 250, 50, "UI v2 Feature Demo")
	view.eventLabel = elements.NewLabel("eventLabel", 250, 80, "Last event: None")

	// Create demo modal
	view.demoModal = containers.NewModal("demoModal", 200, 150, 400, 300)
	modalLabel := elements.NewLabel("modalLabel", 20, 40, "This is a modal dialog!")
	modalButton := elements.NewButton("modalCloseBtn", 150, 250, "Close Modal", "Close this modal")
	modalButton.OnClicked = func() {
		view.demoModal.CloseModal()
	}
	modalLabel.SetParent(view.demoModal)
	modalButton.SetParent(view.demoModal)
	view.demoModal.Children = []elements.ElementInterface{modalLabel, modalButton}
	view.demoModal.OnClose = func() {
		eventMgr.SendEvent(ui.ModalCloseEventData{
			ModalName: "demoModal",
		})
	}

	// Add all elements to the view
	view.AddElement(view.nameInput)
	view.AddElement(view.emailInput)
	view.AddElement(view.passwordInput)
	view.AddElement(view.ageSlider)
	view.AddElement(view.volumeSlider)
	view.AddElement(view.agreeCheckbox)
	view.AddElement(view.notifyCheckbox)
	view.AddElement(view.themeDropdown)
	view.AddElement(view.submitButton)
	view.AddElement(view.clearButton)
	view.AddElement(view.toggleModal)
	view.AddElement(view.statusLabel)
	view.AddElement(view.eventLabel)
	view.AddModal(view.demoModal)

	return view
}

func (v *ExampleView) Update() {
	// Custom view update logic can go here
}

func (v *ExampleView) Draw(screen *ebiten.Image, theme *theming.Theme) {
	// Custom view drawing can go here
	// Labels to describe sections
	// (We could draw section headers or backgrounds here)
}

func NewGame() *Game {
	// Initialize resource manager with UI sprites
	err := resource.LoadImageAsTexture("ui", "assets/ux.png")
	if err != nil {
		log.Printf("Warning: Could not load UI texture: %v", err)
	}

	// Create event manager
	eventMgr := &event.EventManager{}

	// Register event logger
	logger := &UIEventLogger{}
	eventMgr.RegisterListener(logger, ui.EventTypeUIClick)
	eventMgr.RegisterListener(logger, ui.EventTypeUIValueChange)
	eventMgr.RegisterListener(logger, ui.EventTypeUIFocus)
	eventMgr.RegisterListener(logger, ui.EventTypeUIBlur)
	eventMgr.RegisterListener(logger, ui.EventTypeUIModalOpen)
	eventMgr.RegisterListener(logger, ui.EventTypeUIModalClose)

	// Create the example view
	exampleView := NewExampleView(eventMgr)

	// Create GUI with default theme
	gui := ui.NewGUI(exampleView, &theming.DefaultTheme)

	return &Game{
		gui:          gui,
		eventManager: eventMgr,
	}
}

func (g *Game) Update() error {
	g.gui.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(theming.DefaultTheme.Colors.Background)
	g.gui.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("MLGE UI v2 - Feature Demo")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	game := NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
