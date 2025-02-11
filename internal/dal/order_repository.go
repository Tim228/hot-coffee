package dal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"hot-coffee/models"
)

type OrderRepositoryInterface interface {
	CreateOrder(order models.Order) error
	LoadOrders() ([]models.Order, error)
	GetOrderByID(id string) (*models.Order, error)
	DeleteOrder(id string) (*models.Order, error)
	SaveOrders(orders []models.Order) error
}

type OrderRepositoryJSON struct {
	filePath string
}

func NewOrderRepositoryJSON(filePath string) *OrderRepositoryJSON {
	return &OrderRepositoryJSON{filePath: filePath}
}

func (r *OrderRepositoryJSON) CreateOrder(order models.Order) error {
	orders, err := r.LoadOrders()
	if err != nil {
		return err
	}

	orders = append(orders, order)

	return r.SaveOrders(orders)
}

func (r *OrderRepositoryJSON) SaveOrders(orders []models.Order) error {
	filePath := filepath.Join(r.filePath, "orders.json")
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return fmt.Errorf("could not create orders file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(orders); err != nil {
		return fmt.Errorf("could not encode orders to file: %v", err)
	}

	return nil
}

func (r *OrderRepositoryJSON) LoadOrders() ([]models.Order, error) {
	filePath := filepath.Join(r.filePath, "orders.json")
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var orders []models.Order
	if err := json.NewDecoder(file).Decode(&orders); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *OrderRepositoryJSON) GetOrderByID(id string) (*models.Order, error) {
	orders, err := r.LoadOrders()
	if err != nil && orders != nil {
		return &models.Order{}, err
	}

	if len(orders) > 0 {
		for i := 0; i < len(orders); i++ {
			if orders[i].ID == id {
				return &orders[i], nil
			}
		}
	}
	return nil, fmt.Errorf("Order with ID %s not found", id)
}

func (r *OrderRepositoryJSON) DeleteOrder(id string) (*models.Order, error) {
	orders, err := r.LoadOrders()
	if err != nil {
		return &models.Order{}, err
	}

	for i := 0; i < len(orders); i++ {
		if orders[i].ID == id {
			deletedOrder := orders[i]

			orders = append(orders[:i], orders[i+1:]...)

			updatedData, err := json.MarshalIndent(orders, "", "  ")
			if err != nil {
				return nil, fmt.Errorf("Error marshaling updated orders: %v", err)
			}

			err = ioutil.WriteFile("data/orders.json", updatedData, os.ModePerm)
			if err != nil {
				return nil, fmt.Errorf("Error writing updated file: %v", err)
			}

			return &deletedOrder, nil
		}
	}
	return nil, fmt.Errorf("Order with ID %s not found", id)
}
