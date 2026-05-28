package tui

import (
	"errors"
	"net/http"
	"sync/atomic"

	"cours-go/internal/client"
	"cours-go/internal/domain"
	"cours-go/internal/local"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type UI struct {
	app    *tview.Application
	pages  *tview.Pages
	status *tview.TextView

	client *client.Client
	config domain.ClientConfig
	state  domain.ClientState

	started atomic.Bool
}

func New(cfg domain.ClientConfig, state domain.ClientState, apiClient *client.Client) *UI {
	ui := &UI{
		app:    tview.NewApplication(),
		pages:  tview.NewPages(),
		status: tview.NewTextView(),
		client: apiClient,
		config: local.NormalizeConfig(cfg),
		state:  state,
	}

	ui.status.SetDynamicColors(true)
	ui.status.SetText("Prêt")
	ui.applyTheme()
	ui.showHome()

	return ui
}

func (ui *UI) Run() error {
	root := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(ui.pages, 0, 1, true).
		AddItem(ui.status, 1, 0, false)

	ui.started.Store(true)
	defer ui.started.Store(false)

	return ui.app.SetRoot(root, true).EnableMouse(true).Run()
}

func (ui *UI) applyTheme() {
	switch ui.config.Theme {
	case "sunset":
		tview.Styles.PrimitiveBackgroundColor = tcell.Color235
		tview.Styles.ContrastBackgroundColor = tcell.Color94
		tview.Styles.PrimaryTextColor = tcell.ColorOrange
		tview.Styles.SecondaryTextColor = tcell.ColorWhite
		tview.Styles.BorderColor = tcell.ColorOrange
		tview.Styles.TitleColor = tcell.ColorOrange
	case "default":
		tview.Styles = tview.Theme{}
	default:
		tview.Styles.PrimitiveBackgroundColor = tcell.Color17
		tview.Styles.ContrastBackgroundColor = tcell.Color24
		tview.Styles.PrimaryTextColor = tcell.ColorWhite
		tview.Styles.SecondaryTextColor = tcell.ColorLightCyan
		tview.Styles.BorderColor = tcell.ColorDeepSkyBlue
		tview.Styles.TitleColor = tcell.ColorDeepSkyBlue
	}
}

func (ui *UI) setPage(name string, primitive tview.Primitive) {
	ui.pages.RemovePage(name)
	ui.pages.AddAndSwitchToPage(name, primitive, true)
	ui.app.SetFocus(primitive)
}

func (ui *UI) update(fn func()) {
	if !ui.started.Load() {
		fn()
		return
	}

	go ui.app.QueueUpdateDraw(fn)
}

func (ui *UI) info(message string) {
	ui.update(func() {
		ui.status.SetTextColor(tcell.ColorLightGreen)
		ui.status.SetText(message)
	})
}

func (ui *UI) reportError(err error) {
	if err == nil {
		return
	}

	var httpErr *client.HTTPError
	if errors.As(err, &httpErr) && httpErr.StatusCode == http.StatusUnauthorized {
		ui.state = domain.ClientState{}
		ui.client.SetToken("")
		_ = local.SaveState(ui.state)
		ui.update(func() {
			ui.showHome()
			ui.status.SetTextColor(tcell.ColorIndianRed)
			ui.status.SetText("Session expirée — reconnectez-vous")
		})
		return
	}

	ui.update(func() {
		ui.status.SetTextColor(tcell.ColorIndianRed)
		ui.status.SetText(err.Error())
	})
}

func (ui *UI) persistConfig() error {
	ui.config = local.NormalizeConfig(ui.config)
	ui.client.SetBaseURL(ui.config.APIURL)
	ui.applyTheme()
	return local.SaveConfig(ui.config)
}

func (ui *UI) persistState() error {
	ui.client.SetToken(ui.state.Token)
	return local.SaveState(ui.state)
}
