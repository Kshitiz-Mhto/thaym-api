package cart

import (
	"ecom-api/internal/application/core/types/entity"
	"ecom-api/pkg/configs"
	"fmt"
	"math"
)

func getCartItemsIDs(items []entity.CartCheckoutItem) ([]string, error) {
	productIds := make([]string, len(items))
	for index, item := range items {
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("invalid quantity for product %s", item.ProductID)
		}
		productIds[index] = item.ProductID
	}
	return productIds, nil
}

func checkIfCartIsInStock(cartItems []entity.CartCheckoutItem, products map[string]entity.Product) error {
	if len(cartItems) == 0 {
		return fmt.Errorf("cart is empty")
	}
	for _, item := range cartItems {

		product, ok := products[item.ProductID]

		if !ok {
			return fmt.Errorf("product %s is not available in the store, please refresh your cart", item.ProductID)
		}
		if product.Quantity < item.Quantity {
			return fmt.Errorf("product %s is not available in the quantity requested", product.Name)
		}
	}
	return nil
}

func calculateTotalPrice(cartItems []entity.CartCheckoutItem, products map[string]entity.Product) (float64, float64) {
	var totalPriceAfterTaxAndDis float64
	var totalPriceBeforeTaxAndDis float64

	for _, item := range cartItems {
		product := products[item.ProductID]
		itemTotalBefore := product.Price * float64(item.Quantity)
		itemTotalAfter := itemTotalBefore - (item.Discount * product.Price) + (item.Tax * product.Price)

		totalPriceBeforeTaxAndDis += itemTotalBefore
		totalPriceAfterTaxAndDis += itemTotalAfter
	}

	return roundToTwoDecimals(totalPriceBeforeTaxAndDis), roundToTwoDecimals(totalPriceAfterTaxAndDis)
}

func calculateIndivisualProductPricing(product entity.Product, item entity.CartCheckoutItem) (float64, float64) {
	totalPriceBeforeTaxAndDis := product.Price * float64(item.Quantity)
	totalPriceAfterTaxAndDis := totalPriceBeforeTaxAndDis - (item.Discount * product.Price) + (item.Tax * product.Price)

	return roundToTwoDecimals(totalPriceBeforeTaxAndDis), roundToTwoDecimals(totalPriceAfterTaxAndDis)
}

func (handler *CartHandler) createOrder(products []entity.Product, cartItems []entity.CartCheckoutItem, userID string) (string, float64, float64, error) {
	productsMap := make(map[string]entity.Product)
	for _, product := range products {
		productsMap[product.ProductId] = product
	}

	if err := checkIfCartIsInStock(cartItems, productsMap); err != nil {
		return "", 0, 0, err
	}

	totalPriceBeforeTaxAndDis, totalPriceAfterTaxAndDis := calculateTotalPrice(cartItems, productsMap)

	for _, item := range cartItems {
		product := productsMap[item.ProductID]
		product.Quantity -= item.Quantity
		handler.store.UpdateProduct(product)
	}

	orderId, err := handler.orderStore.CreateOrder(entity.Order{
		UserID:        userID,
		Total:         totalPriceAfterTaxAndDis,
		Subtotal:      totalPriceBeforeTaxAndDis,
		Status:        configs.Envs.OrderStatusPending,
		PaymentStatus: configs.Envs.PaymentStatusPending,
		PaymentMethod: "Credit Card",
		Address:       "USA, New York, 123 wall street, Apt:10032, 143B",
		Currency:      configs.Envs.DEFAULT_CURRENCY,
	})

	if err != nil {
		return "", 0, 0, err
	}

	for _, item := range cartItems {
		indivisualSubTotalPrice, indivisualTotalPrice := calculateIndivisualProductPricing(productsMap[item.ProductID], item)
		handler.orderStore.CreateOrderItem(entity.OrderItem{
			OrderID:     orderId,
			ProductID:   item.ProductID,
			ProductName: productsMap[item.ProductID].Name,
			Quantity:    item.Quantity,
			Price:       productsMap[item.ProductID].Price,
			TotalPrice:  indivisualTotalPrice,
			Subtotal:    indivisualSubTotalPrice,
			Currency:    item.Currency,
			Discount:    item.Discount,
			Tax:         item.Tax,
		})
	}

	return orderId, totalPriceBeforeTaxAndDis, totalPriceAfterTaxAndDis, nil
}

func roundToTwoDecimals(value float64) float64 {
	return math.Round(value*100) / 100
}
