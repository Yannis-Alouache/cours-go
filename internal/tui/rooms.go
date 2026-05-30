package tui

import (
	"context"
	"strconv"
	"strings"
	"time"

	"cours-go/internal/domain"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (ui *UI) showRooms() {
	dateField := tview.NewInputField().
		SetLabel("Date (YYYY-MM-DD): ").
		SetText(time.Now().Format("2006-01-02"))

	roomList := tview.NewList()
	roomList.ShowSecondaryText(false)
	roomList.SetBorder(true)
	roomList.SetTitle("Salles")

	slotsTable := tview.NewTable()
	slotsTable.SetBorders(false)
	slotsTable.SetSelectable(true, false)
	slotsTable.SetBorder(true)
	slotsTable.SetTitle("Disponibilités")

	buttons := tview.NewFlex()

	rooms := make([]domain.Room, 0)
	slots := make([]domain.Slot, 0)
	var selectedRoom domain.Room

	populateSlots := func() {
		slotsTable.Clear()
		headers := []string{"Début", "Fin", "Statut"}
		for column, title := range headers {
			slotsTable.SetCell(0, column, tview.NewTableCell(title).SetSelectable(false).SetTextColor(tcell.ColorYellow))
		}
		for index, slot := range slots {
			status := "Libre"
			color := tcell.ColorLightGreen
			if !slot.Available {
				status = "Pris"
				color = tcell.ColorIndianRed
			}
			slotsTable.SetCell(index+1, 0, tview.NewTableCell(slot.Start))
			slotsTable.SetCell(index+1, 1, tview.NewTableCell(slot.End))
			slotsTable.SetCell(index+1, 2, tview.NewTableCell(status).SetTextColor(color))
		}
		if len(slots) == 0 {
			ui.app.SetFocus(buttons)
		}
	}

	loadAvailability := func() {
		if selectedRoom.ID == 0 {
			ui.info("Choisissez une salle")
			return
		}

		day := strings.TrimSpace(dateField.GetText())
		ui.info("Chargement des créneaux...")
		go func() {
			response, err := ui.client.GetAvailability(context.Background(), selectedRoom.ID, day)
			if err != nil {
				ui.reportError(err)
				return
			}
			slots = response
			ui.update(func() {
				populateSlots()
				ui.status.SetText("Disponibilités chargées pour " + selectedRoom.Name)
			})
		}()
	}

	slotsTable.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyTab, tcell.KeyBacktab:
			ui.app.SetFocus(buttons)
		case tcell.KeyEscape:
			ui.app.SetFocus(roomList)
		}
	})

	roomList.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		if index < 0 || index >= len(rooms) {
			return
		}
		selectedRoom = rooms[index]
		loadAvailability()
	})

	loadRooms := func() {
		ui.info("Chargement des salles...")
		go func() {
			response, err := ui.client.GetRooms(context.Background())
			if err != nil {
				ui.reportError(err)
				return
			}
			rooms = response
			ui.update(func() {
				roomList.Clear()
				for _, room := range rooms {
					roomList.AddItem(room.Name, "", 0, nil)
				}
				if len(rooms) > 0 {
					selectedRoom = rooms[0]
					roomList.SetCurrentItem(0)
					loadAvailability()
				}
			})
		}()
	}

	btnActualiser := tview.NewButton("Actualiser").SetSelectedFunc(loadAvailability)
	btnReserver := tview.NewButton("Réserver").SetSelectedFunc(func() {
		row, _ := slotsTable.GetSelection()
		if row <= 0 || row-1 >= len(slots) {
			ui.info("Sélectionnez un créneau")
			return
		}

		slot := slots[row-1]
		if !slot.Available {
			ui.info("Ce créneau est déjà pris")
			return
		}

		startHour, err := strconv.Atoi(strings.Split(slot.Start, ":")[0])
		if err != nil {
			ui.reportError(err)
			return
		}

		ui.info("Réservation en cours...")
		go func() {
			_, err := ui.client.CreateReservation(context.Background(), domain.CreateReservationRequest{
				RoomID:    selectedRoom.ID,
				Day:       strings.TrimSpace(dateField.GetText()),
				StartHour: startHour,
			})
			if err != nil {
				ui.reportError(err)
				return
			}
			ui.update(func() {
				ui.state.Stats.TotalReservations++
				if ui.state.Stats.RoomCount == nil {
					ui.state.Stats.RoomCount = make(map[string]int)
				}
				ui.state.Stats.RoomCount[selectedRoom.Name]++
				_ = ui.persistState()
				ui.status.SetText("Réservation créée pour " + selectedRoom.Name)
				loadAvailability()
			})
		}()
	})
	btnRetour := tview.NewButton("Retour").SetSelectedFunc(func() {
		ui.showHome()
	})

	*buttons = *tview.NewFlex().
		AddItem(btnActualiser, 0, 1, false).
		AddItem(btnReserver, 0, 1, false).
		AddItem(btnRetour, 0, 1, false)

	leftPanel := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(dateField, 3, 0, false).
		AddItem(roomList, 0, 1, true)

	rightPanel := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(slotsTable, 0, 1, false).
		AddItem(buttons, 1, 0, false)

	layout := tview.NewFlex().
		AddItem(leftPanel, 0, 1, true).
		AddItem(rightPanel, 0, 2, false)

	ui.setPage("rooms", layout)
	loadRooms()
}
