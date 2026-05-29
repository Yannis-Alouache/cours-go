package tui

import (
	"context"
	"strconv"

	"cours-go/internal/domain"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (ui *UI) showReservations() {
	table := tview.NewTable()
	table.SetSelectable(true, false)
	table.SetBorder(true)
	table.SetTitle("Mes réservations")

	reservations := make([]domain.Reservation, 0)

	buttons := tview.NewFlex()

	populate := func() {
		table.Clear()
		headers := []string{"ID", "Salle", "Jour", "Début", "Fin"}
		for column, title := range headers {
			table.SetCell(0, column, tview.NewTableCell(title).SetSelectable(false).SetTextColor(tcell.ColorYellow))
		}
		for index, reservation := range reservations {
			values := []string{
				strconv.FormatInt(reservation.ID, 10),
				reservation.RoomName,
				reservation.Day,
				reservation.StartTime,
				reservation.EndTime,
			}
			for column, value := range values {
				table.SetCell(index+1, column, tview.NewTableCell(value))
			}
		}
		if len(reservations) == 0 {
			ui.app.SetFocus(buttons)
		}
	}

	loadReservations := func() {
		ui.info("Chargement des réservations...")
		go func() {
			response, err := ui.client.GetMyReservations(context.Background())
			if err != nil {
				ui.reportError(err)
				return
			}
			reservations = response
			ui.update(func() {
				populate()
				ui.status.SetText("Réservations chargées")
			})
		}()
	}

	btnActualiser := tview.NewButton("Actualiser").SetSelectedFunc(loadReservations)
	btnAnnuler := tview.NewButton("Annuler").SetSelectedFunc(func() {
		row, _ := table.GetSelection()
		if row <= 0 || row-1 >= len(reservations) {
			ui.info("Sélectionnez une réservation")
			return
		}

		reservation := reservations[row-1]
		ui.info("Annulation en cours...")
		go func() {
			if err := ui.client.CancelReservation(context.Background(), reservation.ID); err != nil {
				ui.reportError(err)
				return
			}
			ui.update(func() {
				ui.status.SetText("Réservation annulée")
				loadReservations()
			})
		}()
	})
	btnRetour := tview.NewButton("Retour").SetSelectedFunc(func() {
		ui.showHome()
	})

	tabCapture := func(next, prev tview.Primitive) func(event *tcell.EventKey) *tcell.EventKey {
		return func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyTab:
				ui.app.SetFocus(next)
				return nil
			case tcell.KeyBacktab:
				ui.app.SetFocus(prev)
				return nil
			}
			return event
		}
	}

	table.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyTab:
			ui.app.SetFocus(btnActualiser)
		case tcell.KeyBacktab:
			ui.app.SetFocus(btnRetour)
		case tcell.KeyEscape:
			ui.showHome()
		}
	})

	btnActualiser.SetInputCapture(tabCapture(btnAnnuler, table))
	btnAnnuler.SetInputCapture(tabCapture(btnRetour, btnActualiser))
	btnRetour.SetInputCapture(tabCapture(table, btnAnnuler))

	*buttons = *tview.NewFlex().
		AddItem(btnActualiser, 0, 1, false).
		AddItem(btnAnnuler, 0, 1, false).
		AddItem(btnRetour, 0, 1, false)

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(table, 0, 1, true).
		AddItem(buttons, 1, 0, false)

	ui.setPage("reservations", layout)
	loadReservations()
}
