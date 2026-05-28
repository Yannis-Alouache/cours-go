package tui

import (
	"context"
	"strings"

	"github.com/rivo/tview"
)

func (ui *UI) showLogin() {
	emailField := tview.NewInputField().
		SetLabel("Email: ").
		SetText("alice@example.com")
	passwordField := tview.NewInputField().
		SetLabel("Mot de passe: ").
		SetMaskCharacter('*').
		SetText("password")

	form := tview.NewForm().
		AddFormItem(emailField).
		AddFormItem(passwordField).
		AddButton("Se connecter", func() {
			email := strings.TrimSpace(emailField.GetText())
			password := passwordField.GetText()
			if email == "" || password == "" {
				ui.info("Email et mot de passe requis")
				return
			}

			ui.info("Connexion en cours...")
			go func() {
				response, err := ui.client.Login(context.Background(), email, password)
				if err != nil {
					ui.reportError(err)
					return
				}

				ui.state.Token = response.Token
				ui.state.UserEmail = response.User.Email
				if err := ui.persistState(); err != nil {
					ui.reportError(err)
					return
				}

				ui.update(func() {
					ui.showHome()
					ui.status.SetText("Connecté avec succès")
				})
			}()
		}).
		AddButton("Retour", func() {
			ui.showHome()
		})

	form.SetBorder(true)
	form.SetTitle("Authentification")
	form.SetFocus(0)
	ui.setPage("login", form)
}
