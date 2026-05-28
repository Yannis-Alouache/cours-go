package tui

import (
	"fmt"

	"github.com/rivo/tview"
)

func (ui *UI) showHome() {
	ui.setPage("home", ui.newHomePage())
	ui.info("Accueil prêt")
}

func (ui *UI) newHomePage() tview.Primitive {
	header := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)

	userLine := "Non connecté"
	if ui.state.UserEmail != "" {
		userLine = "Connecté : " + ui.state.UserEmail
	}

	header.SetText(fmt.Sprintf("[::b]Reservation de salles[-:-:-]\n%s\nAPI : %s\nCompte demo : alice@example.com / password", userLine, ui.config.APIURL))

	list := tview.NewList()
	list.ShowSecondaryText(false)
	list.SetBorder(true)
	list.SetTitle("Menu principal")

	if ui.state.Token == "" {
		list.AddItem("Login", "Se connecter au serveur", 'l', func() {
			ui.showLogin()
		})
	} else {
		list.AddItem("Salles", "Consulter les salles et réserver", 's', func() {
			ui.showRooms()
		})
		list.AddItem("Mes réservations", "Afficher et annuler mes réservations", 'r', func() {
			ui.showReservations()
		})
	}

	list.AddItem("Paramètres", "Configurer l'URL API, le thème et la synchro JSON", 'p', func() {
		ui.showSettings()
	})

	if ui.state.Token != "" {
		list.AddItem("Déconnexion", "Effacer la session locale", 'x', func() {
			ui.state.Token = ""
			ui.state.UserEmail = ""
			if err := ui.persistState(); err != nil {
				ui.reportError(err)
				return
			}
			ui.showHome()
			ui.info("Déconnecté")
		})
	}

	list.AddItem("Quitter", "Fermer l'application", 'q', func() {
		ui.app.Stop()
	})

	return tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(header, 5, 0, false).
		AddItem(list, 0, 1, true)
}
