package model

type OrderLine struct {
	ID      int64
	OrderID string
	SKU     string
	Qty     int
}

func NewOrderLine(orderID string, SKU string, qty int) OrderLine {
	return OrderLine{OrderID: orderID, SKU: SKU, Qty: qty}
}

func (line OrderLine) EqualTo(orderLine OrderLine) bool {
	return line.OrderID == orderLine.OrderID && line.Qty == orderLine.Qty && line.SKU == orderLine.SKU
}

func (line OrderLine) Any(lines []OrderLine) bool {
	for _, orderLine := range lines {
		if line.EqualTo(orderLine) {
			return true
		}
	}
	return false
}

func OrderLinesEqual(first, second []OrderLine) bool {
	for _, line := range first {
		if !line.Any(second) {
			return false
		}
	}
	return true
}
