package query

import (
	"context"
	"fmt"

	"github.com/joaodddev/order-fulfillment-platform/internal/domain/order"
)

// Query
type ListOrdersQuery struct {
	CustomerID string
	Status     order.OrderStatus
	Limit      int
	Offset     int
}

// Query Handler
type ListOrdersHandler struct {
	repo order.Repository
}

func NewListOrdersHandler(repo order.Repository) *ListOrdersHandler {
	return &ListOrdersHandler{repo: repo}
}

func (h *ListOrdersHandler) Handle(ctx context.Context, query ListOrdersQuery) ([]*OrderDTO, error) {
	// Validação básica
	if query.Limit == 0 {
		query.Limit = 10 // default
	}
	if query.Limit > 100 {
		query.Limit = 100 // max
	}

	var orders []*order.Order
	var err error

	// Buscar por diferentes critérios
	if query.CustomerID != "" {
		orders, err = h.repo.FindByCustomerID(ctx, query.CustomerID)
	} else if query.Status != "" {
		orders, err = h.repo.FindByStatus(ctx, query.Status)
	} else {
		// Se não houver filtros, buscar todos (com limites)
		// Nota: Idealmente teríamos um método FindAll com paginação
		// Por enquanto, vamos usar FindByStatus com um status vazio
		// ou implementar um método FindAll mais tarde
		orders, err = h.repo.FindByStatus(ctx, query.Status)
		if err != nil {
			// Se não encontrar com status vazio, tentar buscar todos
			// Isso é um fallback, idealmente teríamos um método GetAll
			return nil, fmt.Errorf("failed to list orders: %w", err)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}

	// Aplicar paginação manual (em memória)
	// Idealmente, o repositório deveria suportar paginação
	start := query.Offset
	end := query.Offset + query.Limit
	if start >= len(orders) {
		return []*OrderDTO{}, nil
	}
	if end > len(orders) {
		end = len(orders)
	}
	paginatedOrders := orders[start:end]

	// Mapear para DTOs
	result := make([]*OrderDTO, len(paginatedOrders))
	for i, orderEntity := range paginatedOrders {
		items := make([]OrderItemDTO, len(orderEntity.Items))
		for j, item := range orderEntity.Items {
			items[j] = OrderItemDTO{
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				Price:     item.Price,
			}
		}

		result[i] = &OrderDTO{
			ID:          orderEntity.ID,
			CustomerID:  orderEntity.CustomerID,
			Status:      string(orderEntity.Status),
			TotalAmount: orderEntity.TotalAmount,
			Items:       items,
			CreatedAt:   orderEntity.CreatedAt,
			UpdatedAt:   orderEntity.UpdatedAt,
		}
	}

	return result, nil
}
