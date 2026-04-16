package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dapr/go-sdk/actor"
	"github.com/dapr/go-sdk/actor/config"
	"github.com/dapr/go-sdk/client"
	daprd "github.com/dapr/go-sdk/service/http"
)

const (
	cartActorType = "CartActor"
	stateKeyCart   = "cart"
	timerName      = "cart-expiry"
	expiryDuration = "PT30M" // 30 minutes
)

// Cart holds the shopping cart state for a user.
type Cart struct {
	UserID    string     `json:"user_id"`
	Items     []CartItem `json:"items"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// CartItem represents a single item in the cart.
type CartItem struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

// AddItemRequest is the input for AddItem.
type AddItemRequest struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

// RemoveItemRequest is the input for RemoveItem.
type RemoveItemRequest struct {
	ItemID string `json:"item_id"`
}

// CheckoutResponse is the output of Checkout.
type CheckoutResponse struct {
	OrderID string  `json:"order_id"`
	Items   []CartItem `json:"items"`
	Total   float64 `json:"total"`
}

// CartActor implements the Dapr actor for shopping cart management.
type CartActor struct {
	actor.ServerImplBaseCtx
}

func (a *CartActor) Type() string {
	return cartActorType
}

// getCart loads the cart from actor state, returning an empty cart if none exists.
func (a *CartActor) getCart(ctx context.Context) (*Cart, error) {
	sm := a.GetStateManager()
	var cart Cart
	if err := sm.Get(ctx, stateKeyCart, &cart); err != nil {
		// No state yet — return empty cart.
		cart = Cart{
			UserID:    a.ID(),
			Items:     []CartItem{},
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		}
	}
	return &cart, nil
}

// saveCart persists the cart to actor state and resets the expiry timer.
func (a *CartActor) saveCart(ctx context.Context, cart *Cart) error {
	sm := a.GetStateManager()
	cart.UpdatedAt = time.Now().UTC()
	if err := sm.Set(ctx, stateKeyCart, cart); err != nil {
		return fmt.Errorf("set state: %w", err)
	}
	if err := sm.Save(ctx); err != nil {
		return fmt.Errorf("save state: %w", err)
	}

	// Reset the expiry timer — if the cart is idle for 30 minutes, it auto-clears.
	a.resetExpiryTimer(ctx)
	return nil
}

// resetExpiryTimer registers (or re-registers) a timer that fires Clear after inactivity.
func (a *CartActor) resetExpiryTimer(ctx context.Context) {
	dc, err := client.NewClient()
	if err != nil {
		log.Printf("cartactor: timer client: %s", err)
		return
	}
	defer dc.Close()

	err = dc.RegisterActorTimer(ctx, &client.RegisterActorTimerRequest{
		ActorType: cartActorType,
		ActorID:   a.ID(),
		Name:      timerName,
		DueTime:   expiryDuration,
		CallBack:  "ExpireCart",
	})
	if err != nil {
		log.Printf("cartactor: register timer: %s", err)
	}
}

// GetCart returns the current cart contents.
func (a *CartActor) GetCart(ctx context.Context) (*Cart, error) {
	return a.getCart(ctx)
}

// AddItem adds an item to the cart or increments its quantity if it already exists.
func (a *CartActor) AddItem(ctx context.Context, req *AddItemRequest) (*Cart, error) {
	cart, err := a.getCart(ctx)
	if err != nil {
		return nil, err
	}

	found := false
	for i := range cart.Items {
		if cart.Items[i].ID == req.ID {
			cart.Items[i].Quantity += req.Quantity
			if req.Price > 0 {
				cart.Items[i].Price = req.Price
			}
			if req.Name != "" {
				cart.Items[i].Name = req.Name
			}
			found = true
			break
		}
	}
	if !found {
		cart.Items = append(cart.Items, CartItem{
			ID:       req.ID,
			Name:     req.Name,
			Price:    req.Price,
			Quantity: req.Quantity,
		})
	}

	if err := a.saveCart(ctx, cart); err != nil {
		return nil, err
	}
	log.Printf("cartactor [%s]: added item %s (qty %d)", a.ID(), req.ID, req.Quantity)
	return cart, nil
}

// RemoveItem removes an item from the cart by ID.
func (a *CartActor) RemoveItem(ctx context.Context, req *RemoveItemRequest) (*Cart, error) {
	cart, err := a.getCart(ctx)
	if err != nil {
		return nil, err
	}

	filtered := cart.Items[:0]
	for _, item := range cart.Items {
		if item.ID != req.ItemID {
			filtered = append(filtered, item)
		}
	}
	cart.Items = filtered

	if err := a.saveCart(ctx, cart); err != nil {
		return nil, err
	}
	log.Printf("cartactor [%s]: removed item %s", a.ID(), req.ItemID)
	return cart, nil
}

// Checkout finalizes the cart: computes the total, clears the cart, and returns an order summary.
func (a *CartActor) Checkout(ctx context.Context) (*CheckoutResponse, error) {
	cart, err := a.getCart(ctx)
	if err != nil {
		return nil, err
	}
	if len(cart.Items) == 0 {
		return nil, fmt.Errorf("cart is empty")
	}

	var total float64
	for _, item := range cart.Items {
		total += item.Price * float64(item.Quantity)
	}

	resp := &CheckoutResponse{
		OrderID: fmt.Sprintf("order-%s-%d", a.ID(), time.Now().UnixMilli()),
		Items:   cart.Items,
		Total:   total,
	}

	// Clear the cart after checkout.
	if err := a.clearCart(ctx); err != nil {
		return nil, fmt.Errorf("clear after checkout: %w", err)
	}

	log.Printf("cartactor [%s]: checkout complete — order %s, total %.2f", a.ID(), resp.OrderID, resp.Total)
	return resp, nil
}

// Clear empties the cart.
func (a *CartActor) Clear(ctx context.Context) error {
	return a.clearCart(ctx)
}

// ExpireCart is the timer callback — auto-clears abandoned carts.
func (a *CartActor) ExpireCart(ctx context.Context) error {
	log.Printf("cartactor [%s]: cart expired due to inactivity", a.ID())
	return a.clearCart(ctx)
}

func (a *CartActor) clearCart(ctx context.Context) error {
	sm := a.GetStateManager()
	if err := sm.Remove(ctx, stateKeyCart); err != nil {
		return fmt.Errorf("remove state: %w", err)
	}
	if err := sm.Save(ctx); err != nil {
		return fmt.Errorf("save state: %w", err)
	}
	return nil
}

func main() {
	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "5050"
	}

	svc := daprd.NewService(fmt.Sprintf(":%s", appPort))

	svc.RegisterActorImplFactoryContext(
		func() actor.ServerContext {
			return &CartActor{}
		},
		config.WithSerializerName("json"),
	)

	log.Printf("cartsvc: starting actor host: port %s", appPort)

	if err := svc.Start(); err != nil {
		log.Fatalf("cartsvc: %s", err)
	}
}

// marshalJSON is a helper used by tests and utilities.
func marshalJSON(v any) []byte {
	data, _ := json.Marshal(v)
	return data
}
