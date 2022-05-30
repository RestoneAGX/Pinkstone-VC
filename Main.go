package main

import (
	"BronzeHermes/Cam"
	"BronzeHermes/Database"
	"BronzeHermes/Graph"
	"BronzeHermes/UI"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.NewWithID("Bronze Hermes")
	// go Graph.StartServer()

	Database.DataInit()

	CreateWindow(a)
}

func CreateWindow(a fyne.App) {
	w := a.NewWindow("Bronze Hermes")
	w.SetOnClosed(Graph.StopSever)

	if UI.HandleErrorWindow(Database.LoadData(), w) {
		dialog.ShowInformation("Back Up", "Loading BackUp", w)
		UI.HandleErrorWindow(Database.LoadBackUp(), w)
	}

	w.SetContent(container.NewVBox(container.NewAppTabs(
		container.NewTabItem("Main", makeMainMenu(a, w)),
		container.NewTabItem("Shop", makeShoppingMenu(w)),
		container.NewTabItem("Inventory", makeInfoMenu(w)),
		container.NewTabItem("Statistics", makeStatsMenu(w)),
	)))

	w.ShowAndRun()
}

func makeMainMenu(a fyne.App, w fyne.Window) fyne.CanvasObject {
	return container.NewVBox(
		widget.NewLabelWithStyle("Welcome", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewButton("Save Backup Data", func() {
			go UI.HandleErrorWindow(Database.BackUpAllData(), w)
		}),
		widget.NewButton("Save Backup Data", func() {
			dialog.ShowInformation("Loading Back up Data", "Wait until back up is done loading...", w)
			go UI.HandleErrorWindow(Database.LoadBackUp(), w)
			dialog.ShowInformation("Loaded", "Back Up Loaded", w)
		}),
		//Add inventory features here
	)
}

func makeShoppingMenu(w fyne.Window) fyne.CanvasObject {
	var shoppingCart []Database.Sale

	title := widget.NewLabelWithStyle("Cart Total: 0.0", fyne.TextAlignCenter, fyne.TextStyle{})

	cartList := binding.BindUntypedList(&[]interface{}{})

	shoppingList := widget.NewListWithData(cartList,
		func() fyne.CanvasObject {
			return container.NewBorder(nil, nil, nil, widget.NewButton("X", nil), widget.NewLabel(""))
		}, func(item binding.DataItem, obj fyne.CanvasObject) {})

	shoppingList.OnSelected = func(id widget.ListItemID) {
		shoppingCart[id].Quantity++
		cartList.Reload()
		title.SetText(fmt.Sprintf("Cart Total: %f", Database.GetCartTotal(shoppingCart)))
		shoppingList.Unselect(id)
	}

	shoppingList.UpdateItem = func(id widget.ListItemID, obj fyne.CanvasObject) {
		text := obj.(*fyne.Container).Objects[0].(*widget.Label)
		btn := obj.(*fyne.Container).Objects[1].(*widget.Button)
		val := shoppingCart[id]

		text.SetText(Database.NameKeys[val.ID] + " x" + strconv.Itoa(int(val.Quantity)))
		btn.OnTapped = func() {
			cartList.Set(Database.ConvertCart(Database.DecreaseFromCart(val, shoppingCart)))
			title.SetText(fmt.Sprintf("Cart Total: %1.1f", Database.GetCartTotal(shoppingCart)))
			text.SetText(Database.NameKeys[val.ID] + " x" + strconv.Itoa(int(val.Quantity)))
			shoppingList.Refresh()
		}
	}

	screen := container.New(layout.NewGridLayoutWithRows(3),
		title,
		container.NewMax(shoppingList),
		container.NewGridWithColumns(3,
			widget.NewButton("Buy Cart", func() {
				dialog.ShowConfirm("Buying", "Do you want to buy all items in the Cart?", func(b bool) {
					if !b {
						return
					}
					cartList.Set(Database.ConvertCart(Database.BuyCart(shoppingCart)))
					title.SetText(fmt.Sprintf("Cart Total: %1.1f", Database.GetCartTotal(shoppingCart)))
					dialog.ShowInformation("Complete", "You're Purchase has been made.", w)
				}, w)
			}),
			widget.NewButton("Clear Cart", func() {
				cartList.Set([]interface{}{})
				shoppingCart = shoppingCart[:0]
				title.SetText(fmt.Sprintf("Cart Total: %1.1f", Database.GetCartTotal(shoppingCart)))
			}),
			widget.NewButton("New Item", func() {
				id := Cam.OpenCam(&w)
				if id == 0 {
					return
				}

				item := Database.FindItem(id)

				dialog.ShowCustomConfirm("Just Checking...", "Yes", "No", container.NewVBox(widget.NewLabel("Is this the right item: "+Database.NameKeys[item.ID])),
					func(b bool) {
						if !b {
							return
						}
						cartList.Set(Database.ConvertCart(Database.AddToCart(item, shoppingCart)))
						title.SetText(fmt.Sprintf("Cart Total: %1.1f", Database.GetCartTotal(shoppingCart)))
						shoppingList.Refresh()
					}, w)
			}),
		),
	)
	return screen
}

func makeInfoMenu(w fyne.Window) fyne.CanvasObject {
	idLabel := widget.NewLabel("ID")
	nameLabel := widget.NewLabel("Name")
	priceLabel := widget.NewLabel("Price")
	costLabel := widget.NewLabel("Cost")
	inventoryLabel := widget.NewLabel("Inventory")

	boundData := binding.BindUntypedList(&[]interface{}{})
	boundData.Set(Database.ConvertCart(Database.Databases[0]))

	inventoryList := widget.NewListWithData(boundData, func() fyne.CanvasObject {
		return container.NewBorder(nil, nil, nil, nil, widget.NewLabel("name"))
	}, func(item binding.DataItem, obj fyne.CanvasObject) {})

	inventoryList.UpdateItem = func(idx widget.ListItemID, obj fyne.CanvasObject) {
		obj.(*fyne.Container).Objects[0].(*widget.Label).SetText(Database.NameKeys[Database.Databases[0][idx].ID])
	}

	inventoryList.OnSelected = func(id widget.ListItemID) {
		item := Database.Databases[0][id]
		values := Database.ConvertSale(item)

		idLabel.SetText(strconv.Itoa(int(item.ID)))
		nameLabel.SetText(Database.NameKeys[item.ID])
		priceLabel.SetText(values[0])
		costLabel.SetText(values[1])
		inventoryLabel.SetText(values[2])
	}

	return container.New(layout.NewGridLayout(2),
		container.NewVBox(
			widget.NewLabelWithStyle("Inventory Info", fyne.TextAlign(1), fyne.TextStyle{Bold: true}),
			idLabel,
			nameLabel,
			priceLabel,
			costLabel,
			inventoryLabel,
			widget.NewButton("New", func() {
				id := Cam.OpenCam(&w)
				if id == 0 {
					return
				}

				result := Database.FindItem(id)
				labels := Database.ConvertSale(result)

				idLabel.SetText(strconv.Itoa(id))
				nameLabel.SetText(Database.NameKeys[result.ID])
				priceLabel.SetText(labels[0])
				costLabel.SetText(labels[1])
				inventoryLabel.SetText(labels[2])

			}),
			widget.NewButton("Modify", func() {
				conID, _ := strconv.Atoi(idLabel.Text)
				idLabel := widget.NewLabel(strconv.Itoa(conID))

				nameEntry := widget.NewEntry()
				nameEntry.SetPlaceHolder("Product Name with _ for spaces.")
				nameEntry.Validator = validation.NewRegexp(`^[A-Za-z0-9_-]+$`, "username can only contain letters, numbers, '_', and '-'")

				priceEntry := UI.NewNumEntry("Selling Price")
				costEntry := UI.NewNumEntry("How much you bought it for")
				inventoryEntry := UI.NewNumEntry("Current Inventory")

				dialog.ShowForm("Item", "Save", "Cancel",
					[]*widget.FormItem{
						widget.NewFormItem("ID", idLabel),
						widget.NewFormItem("Name", nameEntry),
						widget.NewFormItem("Price", priceEntry),
						widget.NewFormItem("Cost", costEntry),
						widget.NewFormItem("Inventory", inventoryEntry),
					},
					func(b bool) {
						if !b {
							return
						}

						price, cost, inventory := Database.ConvertString(priceEntry.Text, costEntry.Text, inventoryEntry.Text)
						newItem := Database.Sale{ID: uint64(conID), Price: price, Cost: cost, Quantity: inventory}

						Database.Databases[2] = append(Database.Databases[2], newItem)
						Database.NameKeys[uint64(conID)] = nameEntry.Text

						func(found bool) {
							for i, v := range Database.Databases[0] {
								if v.ID == newItem.ID {
									Database.Databases[0][i] = newItem
									found = true
									break
								}
							}

							if !found {
								Database.Databases[0] = append(Database.Databases[0], newItem)
							}
						}(false)

						boundData.Set(Database.ConvertCart(Database.Databases[0]))

						UI.HandleErrorWindow(Database.SaveData(), w)

						//Updating Entries
						nameLabel.Text = nameEntry.Text
						priceLabel.Text = priceEntry.Text
						costLabel.Text = costEntry.Text
						inventoryLabel.Text = inventoryEntry.Text

						dialog.NewInformation("Success!", "Your data has been saved successfully!", w)
					}, w)
			}),
		),
		container.NewMax(
			inventoryList,
		))
}

func makeStatsMenu(w fyne.Window) fyne.CanvasObject {
	u, _ := url.Parse("http://localhost:8081/line")
	r, _ := url.Parse("http://localhost:8081/pie")

	link := widget.NewHyperlink("Go To Graph", u)

	selectionEntry := UI.NewNumEntry("Year/Month")

	var profitDataSelect int
	var buttonType int

	dataSelectOptions := widget.NewSelect([]string{"Revenue", "Cost", "Profit"}, func(dataType string) {
		switch dataType {
		case "Revenue":
			profitDataSelect = 0
		case "Cost":
			profitDataSelect = 1
		case "Profit":
			profitDataSelect = 2
		}
	})

	financeEntry := UI.NewNumEntry("YYYY/MM/DD")
	reportDisplay := widget.NewLabel("")

	var expense_frequency uint8
	expense_entry := widget.NewEntry()
	expense_amount := UI.NewNumEntry("The amount gained or lost.")

	items := []*widget.FormItem{
		widget.NewFormItem("Name ", expense_entry),
		widget.NewFormItem("Amount ", expense_amount),
		widget.NewFormItem("Frequency ", widget.NewSelect([]string{"Once", "Weekly", "Monthly", "Yearly"}, func(s string) {
			switch s {
			case "Once":
				expense_frequency = 0
			case "Weekly":
				expense_frequency = 1
			case "Monthly":
				expense_frequency = 2
			case "Yearly":
				expense_frequency = 3
			}
		})),
	}

	return container.NewVScroll(container.NewMax(container.NewVBox(
		widget.NewButton("Expense/Gift", func() {
			dialog.ShowForm("Expense", "Create", "Cancel", items, func(b bool) {
				if !b {
					return
				}
				amount, err := strconv.ParseFloat(expense_amount.Text, 32)
				if err != nil {
					log.Println(err)
				}

				Database.Expenses = append(Database.Expenses, Database.Expense{
					Name:      expense_entry.Text,
					Amount:    float32(amount),
					Day:       uint8(time.Now().Day()),
					Month:     uint8(time.Now().Month()),
					Year:      uint8(time.Now().Year()),
					Frequency: expense_frequency,
				})
			}, w)
		}),
		widget.NewCard("Financial Reports", "", container.NewVBox(
			financeEntry,
			widget.NewSelect([]string{"Day", "Month", "Year", "Date"}, func(time string) { // TODO: Finish the times
				financeEntry.Hidden = true
				switch time {
				case "Day":
				case "Month":
				case "Year":
				case "Date":
					financeEntry.Hidden = false
				}
			}),
			reportDisplay,
		)),
		widget.NewCard("Data Graphs", "", container.NewVBox(
			selectionEntry,
			widget.NewSelect([]string{"Items Graph", "Price Changes", "Item Popularity", "Item Sales"}, func(graph string) {
				switch graph {
				case "Items Graph":
					buttonType = 0
					link.URL = u
					dataSelectOptions.Hidden = false
				case "Price Changes":
					buttonType = 1
					link.URL = u
					dataSelectOptions.Hidden = false
				case "Item Popularity":
					buttonType = 2
					link.URL = r
					dataSelectOptions.Hidden = false
				case "Item Sales":
					buttonType = 3
					link.URL = r
					dataSelectOptions.Hidden = true
				case "Sales Over Time":
					buttonType = 4
					link.URL = u
					dataSelectOptions.Hidden = true
				}
			}),
			dataSelectOptions,
			widget.NewButton("Graph", func() {
				switch buttonType {
				case 0:
					Graph.Labels, Graph.LineInputs = Database.GetLine(selectionEntry.Text, profitDataSelect, Database.Databases[1])
				case 1:
					Graph.Labels, Graph.LineInputs = Database.GetLine(selectionEntry.Text, profitDataSelect, Database.Databases[2])
				case 4:
					Graph.Labels, Graph.LineInputs = Database.GetLine(selectionEntry.Text, 3, Database.Databases[1])
				case 2:
					Graph.Labels, Graph.Inputs = Database.GetPie(selectionEntry.Text, profitDataSelect)
				case 3:
					Graph.Labels, Graph.Inputs = Database.GetPie(selectionEntry.Text, 3)
				}
			}),
			link,
		)),
	)))
}
