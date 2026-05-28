package tui

import (
	"context"

	"cours-go/internal/local"

	"github.com/rivo/tview"
)

func (ui *UI) showSettings() {
	themes := []string{"ocean", "sunset", "default"}
	currentTheme := 0
	for index, value := range themes {
		if value == ui.config.Theme {
			currentTheme = index
			break
		}
	}

	apiURLField := tview.NewInputField().
		SetLabel("API URL: ").
		SetText(ui.config.APIURL)

	selectedTheme := ui.config.Theme
	themeField := tview.NewDropDown().
		SetLabel("Theme: ").
		SetOptions(themes, nil).
		SetCurrentOption(currentTheme)
	themeField.SetSelectedFunc(func(optionText string, _ int) {
		selectedTheme = optionText
	})

	form := tview.NewForm().
		AddFormItem(apiURLField).
		AddFormItem(themeField).
		AddButton("Sauvegarder", func() {
			ui.config.APIURL = apiURLField.GetText()
			ui.config.Theme = selectedTheme
			if err := ui.persistConfig(); err != nil {
				ui.reportError(err)
				return
			}
			ui.info("Configuration locale sauvegardée")
			ui.showSettings()
		}).
		AddButton("Importer config", func() {
			ui.info("Import config en cours...")
			go func() {
				cfg, err := ui.client.ImportConfig(context.Background())
				if err != nil {
					ui.reportError(err)
					return
				}
				ui.config = local.NormalizeConfig(cfg)
				if err := ui.persistConfig(); err != nil {
					ui.reportError(err)
					return
				}
				ui.update(func() {
					ui.showSettings()
					ui.status.SetText("Configuration importée depuis le serveur")
				})
			}()
		}).
		AddButton("Exporter config", func() {
			ui.info("Export config en cours...")
			go func() {
				if err := ui.client.ExportConfig(context.Background(), ui.config); err != nil {
					ui.reportError(err)
					return
				}
				ui.info("Configuration exportée vers le serveur")
			}()
		}).
		AddButton("Importer état", func() {
			ui.info("Import état en cours...")
			go func() {
				state, err := ui.client.ImportState(context.Background())
				if err != nil {
					ui.reportError(err)
					return
				}
				ui.state = state
				if err := ui.persistState(); err != nil {
					ui.reportError(err)
					return
				}
				ui.update(func() {
					ui.showHome()
					ui.status.SetText("État importé depuis le serveur")
				})
			}()
		}).
		AddButton("Exporter état", func() {
			ui.info("Export état en cours...")
			go func() {
				if err := ui.client.ExportState(context.Background(), ui.state); err != nil {
					ui.reportError(err)
					return
				}
				ui.info("État exporté vers le serveur")
			}()
		}).
		AddButton("Retour", func() {
			ui.showHome()
		})

	form.SetBorder(true)
	form.SetTitle("Paramètres")
	ui.setPage("settings", form)
}
