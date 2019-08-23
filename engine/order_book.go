package engine

import "sort"

func (a XDOrders) Len() int   { return len(a) }
func (a XDOrders) Swap(i, j int)  { a[i], a[j] = a[j], a[i] }
// Ranking from small to large
func (a XDOrders) Less(i, j int) bool { return a[i].Price < a[j].Price }

func (a DXOrders) Len() int   { return len(a) }
func (a DXOrders) Swap(i, j int)  { a[i], a[j] = a[j], a[i] }
// Ranking from large to small
func (a DXOrders) Less(i, j int) bool { return a[i].Price > a[j].Price }


type  XDOrders []Order
type  DXOrders []Order

// OrderBook type
type OrderBook struct {
	BuyOrders  []Order
	SellOrders []Order
}

// Add a buy order to the order book
func (book *OrderBook) addBuyOrder(order Order) {
	book.BuyOrders = append(book.BuyOrders, order)
	sort.Sort(XDOrders(book.BuyOrders))
}

// Add a sell order to the order book
func (book *OrderBook) addSellOrder(order Order) {
	book.SellOrders = append(book.SellOrders, order)
	sort.Sort(DXOrders(book.SellOrders))

}

// Remove a buy order from the order book at a given index
func (book *OrderBook) removeBuyOrder(index int) {
	book.BuyOrders = append(book.BuyOrders[:index], book.BuyOrders[index+1:]...)
}

// Remove a sell order from the order book at a given index
func (book *OrderBook) removeSellOrder(index int) {
	book.SellOrders = append(book.SellOrders[:index], book.SellOrders[index+1:]...)
}
