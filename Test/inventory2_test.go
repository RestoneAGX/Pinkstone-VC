package Test

import (
	"BronzeHermes/Database"
	"strconv"
	"testing"
	"time"
)

func TestRemoveSale(t *testing.T) {
	resetTestItemsAndSales()
	Database.Items = testItems
	Database.Sales = []Database.Sale{
		{ID: 6, Cost: 2, Quantity: 15},
	}

	Database.RemoveFromSales(0)

	// check the length of sales
	if len(Database.Sales) != 0 {
		t.Errorf("Item not removed | len: %d", len(Database.Sales))
	}
	// Check if item[6]'s quantity has increase for cost 2
	if Database.Items[6].Quantity[0] != 18 {
		t.Errorf("Error occured with Item's quantity[0] | have: %f", Database.Items[6].Quantity[0])
	}

	if Database.Items[6].Quantity[1] != 4 || Database.Items[6].Quantity[2] != 7 {
		t.Errorf("Error occured with other quantiites | want [1]: 4 & want [2]: 7, have: %v", Database.Items[6].Quantity)
	}
}

func TestRemoveSale2(t *testing.T) {
	resetTestItemsAndSales()
	Database.Items = testItems
	Database.Sales = []Database.Sale{
		{ID: 6, Cost: 3, Quantity: 15},
	}

	Database.RemoveFromSales(0)

	if len(Database.Sales) != 0 {
		t.Errorf("Item not removed | len: %d", len(Database.Sales))
	}

	if Database.Items[6].Quantity[1] != 19 {
		t.Errorf("Error occured with Item's quantity[0] | have: %f", Database.Items[6].Quantity[0])
	}

	if Database.Items[6].Quantity[0] != 3 || Database.Items[6].Quantity[2] != 7 {
		t.Errorf("Error occured with other quantiites | want [1]: 4 & want [2]: 7, have: %v", Database.Items[6].Quantity)
	}
}

func TestAddingDamages(t *testing.T) {
	Database.Sales = []Database.Sale{}
	Database.Items = testItems

	y, month, day := time.Now().Date()
	year, _ := strconv.Atoi(strconv.Itoa(y)[1:])

	answer := Database.Sale{
		ID:       6,
		Price:    0,
		Cost:     Database.Items[0].Cost[0],
		Quantity: 2,
		Year:     uint8(year),
		Month:    uint8(month),
		Day:      uint8(day),
		Usr:      255,
	}

	errID := Database.AddDamages(6, "2")
	if errID != -1 {
		t.Errorf("Some error has occured | have: %d, want: -1", errID)
	}

	if len(Database.Sales) != 1 {
		t.Log(Database.Sales)
		t.Errorf("Issue adding the damages to sales | have: %d, want: 1", len(Database.Sales))
	}

	if Database.Sales[0] != answer {
		t.Errorf("Sales do not match up with the expected | have: %v, want: %v", Database.Sales[0], answer)
	}
}
