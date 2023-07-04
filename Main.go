package main

import (
	"BronzeHermes/Database"
	"BronzeHermes/UI"
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.NewWithID("PINKSTONE")
	// go Graph.StartServer()

	Database.DataInit(false)

	// TODO: Access USB
	// TODO: Connect to Printer

	CreateWindow(a)
}

var w fyne.Window

func CreateWindow(a fyne.App) {
	w = a.NewWindow("Pinkstone")
	w.SetOnClosed(
		func() {
			// Graph.StopSever()
			Database.CleanUpDeadItems()
			Database.SaveData()
			Database.SaveBackUp()
		},
	)

	if UI.HandleErrorWindow(Database.LoadData(), w) {
		dialog.ShowInformation("Back Up", "Loading BackUp", w)
		UI.HandleErrorWindow(Database.LoadBackUp(), w)
	}

	w.SetContent(container.NewVBox(container.NewAppTabs(
		container.NewTabItem("Main", makeMainMenu(a)),
		container.NewTabItem("Shop", makeShoppingMenu()),
		container.NewTabItem("Inventory", Database.MakeInfoMenu(w)),
		container.NewTabItem("Statistics", makeStatsMenu()),
	)))

	// Start Sign In Menu
	w.Content().(*fyne.Container).Objects[0].(*container.AppTabs).Items[0].Content.(*fyne.Container).Objects[1].(*widget.Button).OnTapped()
	w.Content().(*fyne.Container).Objects[0].(*container.AppTabs).OnSelected = func(ti *container.TabItem) {
		updateReport()
	}

	w.ShowAndRun()
}

func makeMainMenu(a fyne.App) fyne.CanvasObject {
	var SignInStartUp dialog.Dialog
	var CreateUser dialog.Dialog

	usrData := binding.NewStringList()
	usrData.Set(Database.FilterUsers())
	titleText := widget.NewLabelWithStyle("Welcome", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	return container.NewVBox(
		titleText,
		widget.NewButton("Sign-In", func() {

			nameEntry := widget.NewEntry()
			usrList := widget.NewListWithData(usrData, func() fyne.CanvasObject {
				return container.NewBorder(nil, nil, nil, widget.NewButton("x", nil), widget.NewLabel(""))
			}, func(di binding.DataItem, co fyne.CanvasObject) {
				text := co.(*fyne.Container).Objects[0].(*widget.Label)
				s, _ := di.(binding.String).Get()
				text.SetText(s)

				co.(*fyne.Container).Objects[1].(*widget.Button).OnTapped = func() {
					for i := 0; i < len(Database.Users); i++ {
						if Database.Users[i] == s {
							Database.Users[i] = string([]byte{216}) + Database.Users[i]
							break
						}
					}
				}
			})

			usrList.OnSelected = func(id widget.ListItemID) {
				Database.Current_User = uint8(id)
			}

			CreateUser = dialog.NewForm("New", "Create", "Back", []*widget.FormItem{
				widget.NewFormItem("Username", nameEntry),
			}, func(b bool) {
				if !b || nameEntry.Text == "" {
					SignInStartUp.Show()
					return
				}

				Database.Users = append(Database.Users, nameEntry.Text)
				Database.Current_User = uint8(len(Database.Users) - 1)
				titleText.SetText("Welcome " + nameEntry.Text) // change the title Text
				usrData.Set(Database.Users)
				Database.SaveData()
			}, w)

			SignInStartUp = dialog.NewCustomConfirm("Sign In", "Login", "Create New", container.NewMax(usrList), func(b bool) {
				if b && len(Database.Users) > 0 {
					titleText.SetText("Welcome " + Database.Users[Database.Current_User])
				} else {
					CreateUser.Show()
				}
			}, w)

			SignInStartUp.Show()
		}),

		widget.NewButton("Save Backup Data", func() {
			go UI.HandleErrorWindow(Database.SaveBackUp(), w)
		}),
		widget.NewButton("Load Backup Data", func() {
			dialog.ShowInformation("Loading Back up Data", "Wait until back up is done loading...", w)
			UI.HandleErrorWindow(Database.LoadBackUp(), w)
			dialog.ShowInformation("Loaded", "Back Up Loaded", w)
		}),
		widget.NewButton("Delete Database", func() {
			dialog.ShowConfirm("Are you sure?", "DELETE EVERYTHING",
				func(confirmed bool) {
					if !confirmed {
						return
					}
					dialog.ShowConfirm("You sure you sure?", "You sure you sure?", func(b bool) {
						if !confirmed {
							return
						}

						Database.Items = map[uint16]*Database.Entry{}
						Database.Sales = []Database.Sale{}
						Database.Current_User = 0
						Database.Users = []string{}
						Database.Customers = []string{}
						CreateUser.Show()

						usrData.Set(nil)
						Database.InventoryData.Set(nil)
						titleText.SetText("Welcome")
						Database.SaveData()
					}, w)
				}, w)
		}),

		//Add inventory features here
	)
}

var shoppingCart []Database.Sale

func makeShoppingMenu() fyne.CanvasObject {

	title := widget.NewLabelWithStyle("Cart Total: 0.0", fyne.TextAlignCenter, fyne.TextStyle{})

	cartData := binding.NewUntypedList()

	shoppingList := widget.NewListWithData(cartData,
		func() fyne.CanvasObject {
			return container.NewBorder(nil, nil, nil, widget.NewButton("-", nil), widget.NewLabel(""))
		}, func(item binding.DataItem, obj fyne.CanvasObject) {})

	shoppingList.OnSelected = func(id widget.ListItemID) {
		shoppingCart[id].Quantity++
		cartData.Set(Database.ConvertCart(shoppingCart))
		title.SetText(fmt.Sprintf("Cart Total: %1.2f", Database.GetCartTotal(shoppingCart)))
		shoppingList.Unselect(id)
	}

	shoppingList.UpdateItem = func(id widget.ListItemID, obj fyne.CanvasObject) {
		text := obj.(*fyne.Container).Objects[0].(*widget.Label)
		btn := obj.(*fyne.Container).Objects[1].(*widget.Button)
		v, _ := cartData.GetValue(id)
		val := v.(Database.Sale)
		text.SetText(Database.Items[val.ID].Name + " x" + fmt.Sprint(val.Quantity))
		text.SetText(fmt.Sprintf("%s ₵%1.2f x%1.2f -> ₵%1.2f", Database.Items[val.ID].Name, val.Price, val.Quantity, val.Price*val.Quantity))
		btn.OnTapped = func() {
			shoppingCart = Database.DecreaseFromCart(val, shoppingCart)
			cartData.Set(Database.ConvertCart(shoppingCart))
			title.SetText(fmt.Sprintf("Cart Total: %1.2f", Database.GetCartTotal(shoppingCart)))
			text.SetText(Database.Items[val.ID].Name + " x" + fmt.Sprint(val.Quantity))
			shoppingList.Refresh()
		}
	}

	customerEntry := UI.NewSearchBar("Customer Name Here", Database.SearchCustomers)

	return container.New(layout.NewGridLayoutWithRows(3),
		title,
		container.NewMax(shoppingList),
		container.NewGridWithColumns(3,
			widget.NewButton("Buy Cart", func() {
				customerEntry.SetText("")
				dialog.ShowForm("Do you want to buy all items in the Cart?", "Yes", "No",
					[]*widget.FormItem{widget.NewFormItem("Customer", customerEntry)}, func(b bool) {
						if !b || len(shoppingCart) == 0 {
							return
						}

						i := 0

						for ; i < len(Database.Customers) && customerEntry.Text != Database.Customers[i]; i++ {
						}

						if i == len(Database.Customers) {
							Database.Customers = append(Database.Customers, customerEntry.Text)
						}

						receipt := Database.MakeReceipt(shoppingCart, customerEntry.Text)
						shoppingCart = Database.BuyCart(shoppingCart, i)
						cartData.Set(Database.ConvertCart(shoppingCart))
						title.SetText("Cart Total: 0.0")
						txtDisplay := widget.NewLabelWithStyle(receipt, fyne.TextAlignCenter, fyne.TextStyle{})

						dialog.ShowCustomConfirm("Complete", "Print", "Done", container.NewVBox(
							widget.NewLabelWithStyle("PINKSTONE TRADING", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
							txtDisplay,
						), func(printing bool) {
							if printing {
								// TODO: Send print msg to Printer & the receipt,
							}
						}, w)
					}, w)
			}),

			widget.NewButton("Clear Cart", func() {
				cartData.Set([]interface{}{})
				shoppingCart = shoppingCart[:0]
				title.SetText(fmt.Sprintf("Cart Total: %1.2f", Database.GetCartTotal(shoppingCart)))
			}),

			widget.NewButton("New Item", func() {
				searchBar := UI.NewSearchBar("Type Item Name", Database.SearchInventory)

				dialog.ShowCustomConfirm("Scan Item", "Confirm", "Cancel", searchBar, func(confirmed bool) {
					if !confirmed {
						return
					}

					id := searchBar.Result()

					if UI.HandleKnownError(0, id < 0, w) {
						return
					}

					val := Database.Items[uint16(id)]

					barginEntry := UI.NewNumEntry("An Adjusted price based on the customer")
					pieceEntry := UI.NewNumEntry("How many pieces are you buying?")
					totalEntry := UI.NewNumEntry("How many in total?")

					options := widget.NewAccordion(
						widget.NewAccordionItem("Bargin Price", barginEntry),
						widget.NewAccordionItem("Pieces", container.NewVBox(pieceEntry, totalEntry)))

					var menu dialog.Dialog
					menu = dialog.NewCustomConfirm("Just Checking...", "Yes", "No", container.NewVBox(widget.NewLabel(val.Name), options),
						func(b bool) {
							if !b {
								return
							}
							s := Database.ConvertItem(uint16(id))

							if pieceEntry.Text != "" || totalEntry.Text != "" {
								if UI.HandleKnownError(1, pieceEntry.Text == "" || totalEntry.Text == "", w) {
									menu.Show()
								} else {
									piece, err := strconv.ParseFloat(pieceEntry.Text, 32)

									if UI.HandleKnownError(0, err != nil || piece < 0, w) {
										menu.Show()
									}

									total, err := strconv.ParseFloat(totalEntry.Text, 32)
									if UI.HandleKnownError(0, err != nil || total < 0, w) {
										menu.Show()
									}
									s.Quantity = float32(piece / total)
								}
							}

							if barginEntry.Text != "" {
								f, err := strconv.ParseFloat(barginEntry.Text, 32)
								UI.HandleKnownError(0, err != nil, w)
								s.Price = float32(f) / s.Quantity
							}

							shoppingCart = Database.AddToCart(s, shoppingCart)
							cartData.Set(Database.ConvertCart(shoppingCart))
							title.SetText(fmt.Sprintf("Cart Total: %1.2f", Database.GetCartTotal(shoppingCart)))
							shoppingList.Refresh()
							shoppingList.ScrollToBottom()
						}, w)
					menu.Show()
				}, w)
			}),
		),
	)
}

var updateReport func()

func makeStatsMenu() fyne.CanvasObject {
	/*
		u, _ := url.Parse("http://localhost:8081/line")
		r, _ := url.Parse("http://localhost:8081/pie")

		link := widget.NewHyperlink("Go To Graph", u)

		selectionEntry := UI.NewNumEntry("YYYY-MM")

		var buttonType int
	*/

	reportDisplay := widget.NewLabel("")
	financeEntry := UI.NewNumEntry("YYYY-MM-DD")
	financeEntry.Hidden = true

	var variant uint8
	date := []uint8{}

	updateReport = func() {
		reportDisplay.SetText(Database.CompileReport(variant, date))
	}

	/*
		updateGraph := func() {
			switch buttonType {
			case 0:
				Graph.Labels, Graph.LineInputs = Database.GetLine(selectionEntry.Text, 0, 0)
			case 1:
				Graph.Labels, Graph.Inputs = Database.GetPie(selectionEntry.Text, 1)
			case 2:
				Graph.Labels, Graph.LineInputs = Database.GetLine(selectionEntry.Text, 1, 0)
			}
		}
	*/

	customerSearch := UI.NewSearchBar("Customer Name here...", Database.SearchCustomers)

	reportData := binding.NewUntypedList()
	reportData.Set(Database.ConvertCart(Database.Sales))

	reportList := widget.NewListWithData(reportData, func() fyne.CanvasObject {
		return container.NewBorder(nil, nil, nil, nil, widget.NewLabel(""))
	}, func(di binding.DataItem, co fyne.CanvasObject) {
		v, _ := di.(binding.Untyped).Get()
		val := v.(Database.Sale)
		display := co.(*fyne.Container).Objects[0].(*widget.Label)
		display.SetText(fmt.Sprintf("%s x%1.2f for ₵%1.2f [%2d-%2d-20%2d] Customer: %s, Cashier: %s",
			Database.Items[val.ID].Name, val.Quantity, val.Price*val.Quantity, val.Day, val.Month, val.Year, Database.Customers[val.Customer], Database.Users[val.Usr]))
	})

	reportList.OnSelected = func(id widget.ListItemID) {
		v, err := reportData.GetValue(id)
		UI.HandleError(err)
		val := v.(Database.Sale)

		infoText := fmt.Sprintf("Name: %s\nPrice: %1.2f\nCost: %1.2f\nQuantity: %1.2f\nTotal Revenue: %1.2f\nTotal Profit: %1.2f\nCustomer: %s\nCashier:%s",
			Database.Items[val.ID].Name, val.Price, val.Cost, val.Quantity, val.Price*val.Quantity, (val.Price-val.Cost)*val.Quantity, Database.Customers[val.Customer], Database.Users[val.Usr])

		dialog.ShowCustomConfirm("Info", "Refund", "Close", widget.NewLabel(infoText), func(b bool) {
			if !b {
				return
			}

			Database.RemoveReportEntry(id)
			UI.HandleError(Database.SaveData())
			reportData.Set(Database.ConvertCart(Database.Sales))
			reportList.Refresh()
		}, w)
		reportList.UnselectAll()
	}

	content := container.New(layout.NewGridLayoutWithRows(3),
		widget.NewCard("Financial Reports", "", container.NewVBox(
			financeEntry,
			widget.NewSelect([]string{"Day", "Month", "Year", "Date"}, func(time string) {
				financeEntry.Hidden = true

				switch time {
				case "Day":
					variant = Database.ONCE
				case "Month":
					variant = Database.MONTHLY
				case "Year":
					variant = Database.YEARLY
				case "Date": //The user will have to double tap when using Dates
					financeEntry.Hidden = false
					if financeEntry.Text == "" {
						reportDisplay.SetText("Type a date and select the date option again to get a report")
						return
					}

					raw := strings.SplitN(financeEntry.Text, "-", 3)

					year, err := strconv.Atoi(raw[0][1:])
					if err != nil {
						return
					}

					variant = Database.YEARLY

					var month int
					var day int

					if len(raw) > 1 {
						month, _ = strconv.Atoi(raw[1]) // NOTE: Error handling may be needed here, unknown for now
						variant = Database.MONTHLY
					}

					if len(raw) > 2 {
						day, _ = strconv.Atoi(raw[2])
						variant = Database.ONCE
					}

					date = []uint8{uint8(day), uint8(month), uint8(year)}
				}

				updateReport()
			}),
			reportDisplay,
		)),

		widget.NewCard("Sales", "", container.NewVBox(
			customerSearch,
			widget.NewButton("Search", func() {
				customerIdx := customerSearch.Result()
				found := []Database.Sale{}

				fmt.Println(customerIdx)

				if customerIdx == -1 {
					found = Database.Sales

				} else {
					for _, v := range Database.Sales {
						if v.Customer == uint8(customerIdx) {
							found = append(found, v)
						}
					}
				}

				reportData.Set(Database.ConvertCart(found))
				reportList.Refresh()
			}),
		)),
		container.NewMax(reportList),

		/*
			widget.NewCard("Data Graphs", "", container.NewVBox(
				selectionEntry,
				widget.NewSelect([]string{"Items Graph", "Price Changes", "Item Popularity", "Item Sales"}, func(graph string) {
					switch graph {
					case "Items Graph":
						buttonType = 0
						link.URL = u
					case "Price Changes":
						buttonType = 1
						link.URL = u
					case "Item Popularity":
						buttonType = 2
						link.URL = r
					case "Item Sales":
						buttonType = 3
						link.URL = r
					case "Sales Over Time":
						buttonType = 4
						link.URL = u
					}
				}),
				link,
			)),
		*/
	)
	return content
}
