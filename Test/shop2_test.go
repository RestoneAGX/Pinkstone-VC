package Test

import (
	"BronzeHermes/Database"
	"testing"
)

func TestProcessingItemNormal(t *testing.T) {
	s := Database.Sale{Price: 5, Cost: 5, Quantity: 1}
	err := Database.ProcessNewItemData("", "", "", &s)
	switch err {
	case 0:
		t.Error("Invalid Input passed into ProcessNewItemData")
	case 1:
		t.Error("No input sent to the piece || total")
	case -1:
		// Check that teh proper transformation
		if s.Quantity != 1 {
			t.Errorf("Illegal Quantity Modification, want: 1, have: %f", s.Quantity)
		}

		if s.Price != 5 {
			t.Errorf("Illegal Price Modificaiton, want: 5, have: %f", s.Price)
		}

		if s.Cost != 5 {
			t.Errorf("Illegal Cost Modificaiton, want: 5, have: %f", s.Cost)
		}
	}
}

func TestProcessingItemWithBargin(t *testing.T) {
	s := Database.Sale{Price: 5, Cost: 5, Quantity: 1}
	err := Database.ProcessNewItemData("10", "", "", &s)
	switch err {
	case 0:
		t.Error("Invalid Input passed into ProcessNewItemData")
	case 1:
		t.Error("No input sent to the piece || total")
	case -1:
		// Check that teh proper transformation
		if s.Quantity != 1 {
			t.Errorf("Illegal Quantity Modification, want: 1, have: %f", s.Quantity)
		}

		if s.Price != 10 {
			t.Errorf("Illegal Price Modificaiton, want: 10, have: %f", s.Price)
		}

		if s.Cost != 5 {
			t.Errorf("Illegal Cost Modificaiton, want: 5, have: %f", s.Cost)
		}
	}
}

func TestProcessingItemInPieces(t *testing.T) {
	s := Database.Sale{Price: 5, Cost: 5, Quantity: 1}
	err := Database.ProcessNewItemData("", "1", "12", &s)
	switch err {
	case 0:
		t.Error("Invalid Input passed into ProcessNewItemData")
	case 1:
		t.Error("No input sent to the piece || total")
	case -1:
		// Check that teh proper transformation
		if s.Quantity != 1.0/12.0 {
			t.Errorf("Illegal Quantity Modification, want: 1/12, have: %f", s.Quantity)
		}

		if s.Price != 5 {
			t.Errorf("Illegal Price Modificaiton, want: 5, have: %f", s.Price)
		}

		if s.Cost != 5 {
			t.Errorf("Illegal Cost Modificaiton, want: 5, have: %f", s.Cost)
		}
	}
}

func TestProcessingItemWithQuantityInPieces(t *testing.T) {
	s := Database.Sale{Price: 5, Cost: 5, Quantity: 1}
	err := Database.ProcessNewItemData("25", "1", "12", &s)
	switch err {
	case 0:
		t.Error("Invalid Input passed into ProcessNewItemData")
	case 1:
		t.Error("No input sent to the piece || total")
	case -1:
		// Check that teh proper transformation
		if s.Quantity != 1.0/12.0 {
			t.Errorf("Illegal Quantity Modification, want: 1/12, have: %f", s.Quantity)
		}

		if s.Price != 25.0*12.0 {
			t.Errorf("Illegal Price Modificaiton, want: 5, have: %f", s.Price)
		}

		if s.Cost != 5 {
			t.Errorf("Illegal Cost Modificaiton, want: 5, have: %f", s.Cost)
		}
	}
}

func TestInvalidBargin(t *testing.T) {
	s := Database.Sale{Price: 5, Cost: 5, Quantity: 1}
	err := Database.ProcessNewItemData("-122-123-2", "", "", &s)
	switch err {
	case 0:
	case 1:
		t.Error("No input sent to the piece || total")
	case -1:
		t.Log(s)
		t.Error("This Data is invalid and should not pass")
	}
}

func TestMissingPiece(t *testing.T) {
	s := Database.Sale{Price: 5, Cost: 5, Quantity: 1}
	err := Database.ProcessNewItemData("", "", "12", &s)
	switch err {
	case 0:
		t.Error("Invalid Input passed into ProcessNewItemData")
	case 1:
	case -1:
		t.Log(s)
		t.Error("This Data is invalid and should not pass")
	}
}

func TestMissingTotal(t *testing.T) {
	s := Database.Sale{Price: 5, Cost: 5, Quantity: 1}
	err := Database.ProcessNewItemData("", "1", "", &s)
	switch err {
	case 0:
		t.Error("Invalid Input passed into ProcessNewItemData")
	case 1:
	case -1:
		t.Log(s)
		t.Error("This Data is invalid and should not pass")
	}
}

func TestAddingIndividual(t *testing.T) {
	s := Database.Sale{Price: 5, Cost: 5, Quantity: 1}
	cart := []Database.Sale{}
	cart = Database.AddToCart(s, cart)
	if len(cart) != 1 {
		t.Errorf("Error adding item to cart | len: %d, cart: %v", len(cart), cart)
	}

	for _, v := range cart {
		if v != s {
			t.Errorf("Cart and Item are != | item: %v, cart: %v", s, v)
		}
	}
}

func TestAddingToItemInCart(t *testing.T) {
	s := Database.Sale{Price: 5, Cost: 5, Quantity: 3}
	cart := []Database.Sale{
		{Price: 5, Cost: 5, Quantity: 1},
	}
	cart = Database.AddToCart(s, cart)
	s.Quantity = 4 //set this so that we can just compare them directly without having to check each individual stat

	if len(cart) != 1 {
		t.Errorf("Error adding item to cart | len: %d, cart: %v", len(cart), cart)
	}

	for _, v := range cart {
		if v != s {
			t.Errorf("Cart and Item are != | item: %v, cart: %v", s, v)
		}
	}
}

func TestAddingMoreToCart(t *testing.T) {
	answer := []Database.Sale{
		{Price: 5, Cost: 5, Quantity: 1},
		{Price: 6, Cost: 5, Quantity: 1},
	}
	cart := []Database.Sale{
		{Price: 5, Cost: 5, Quantity: 1},
	}
	cart = Database.AddToCart(answer[1], cart)

	if len(cart) != 2 {
		t.Errorf("Error adding item to cart | len: %d, cart: %v", len(cart), cart)
	}

	for i := range cart {
		if cart[i] != answer[i] {
			t.Errorf("Cart and Item are != | item: %v, cart: %v", answer[i], cart[i])
		}
	}
}

func TestDeductItemFromCart(t *testing.T) {
	cart := []Database.Sale{
		{Price: 5, Cost: 5, Quantity: 3},
	}

	cart = Database.DecreaseFromCart(0, cart)

	if cart[0].Quantity != 2 {
		t.Errorf("Error reducing quantity from item in cart | have: %f, want: 2.0", cart[0].Quantity)
	}

	if cart[0].Price != 5 || cart[0].Cost != 5 {
		t.Errorf("Error occured with price or cost of cart item | have: %v", cart[0])
	}
}

func TestItemRemovalFromCart(t *testing.T) {
	cart := []Database.Sale{
		{Price: 5, Cost: 5, Quantity: 1},
	}

	cart = Database.DecreaseFromCart(0, cart)

	if len(cart) != 0 {
		t.Errorf("An Erorr occured deleting an item from cart | have: %v", len(cart))
	}
}

func TestCartTotal(t *testing.T) {
	total := Database.GetCartTotal(testSales)
	if total != 30 {
		t.Log(total)
		t.Error("Some Sort of addition error in GetCartTotal")
	}
}
