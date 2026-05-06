# 08 - Actors: Shopping Cart

A real-world shopping cart built with Dapr virtual actors. Each user gets their own `CartActor` instance — isolated, single-threaded, and automatically persisted.

## What This Demonstrates

| Dapr Actor Feature | How It's Used |
|--------------------|---------------|
| Actor registration | `cartsvc` registers `CartActor` with the Dapr runtime |
| Actor invocation   | `frontendsvc` calls actor methods via the Dapr client |
| Actor state        | Cart contents persisted per-actor in Redis |
| Actor timers       | Abandoned carts auto-expire after 30 minutes of inactivity |

## Architecture

```
                  ┌─────────────┐
  HTTP requests   │ frontendsvc │   (port 8080)
 ───────────────► │  (Go HTTP)  │
                  └──────┬──────┘
                         │ Dapr Actor proxy
                         ▼
                  ┌─────────────┐
                  │   cartsvc   │   (port 5050)
                  │ (Actor Host)│
                  └──────┬──────┘
                         │ Actor state store
                         ▼
                  ┌─────────────┐
                  │    Redis    │
                  └─────────────┘
```

## API

| Method   | Path                              | Description                |
|----------|-----------------------------------|----------------------------|
| `GET`    | `/cart/{userId}`                  | Get cart contents          |
| `POST`   | `/cart/{userId}/items`            | Add item to cart           |
| `DELETE`  | `/cart/{userId}/items/{itemId}`   | Remove item from cart      |
| `POST`   | `/cart/{userId}/checkout`         | Checkout and get order     |
| `DELETE`  | `/cart/{userId}`                  | Clear entire cart          |

### Example: Add an item

```bash
curl -X POST http://localhost:8080/cart/user-1/items \
  -H "Content-Type: application/json" \
  -d '{"id": "item-1", "name": "Wireless Mouse", "price": 29.99, "quantity": 2}'
```

### Example: Checkout

```bash
curl -X POST http://localhost:8080/cart/user-1/checkout
```

Returns:
```json
{
  "order_id": "order-user-1-1713200000000",
  "items": [{"id": "item-1", "name": "Wireless Mouse", "price": 29.99, "quantity": 2}],
  "total": 59.98
}
```

## Build & Deploy

### Install redis on k8s
```
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update
helm install redis bitnami/redis
```

### Build

```
ko build -B -L ./frontendsvc --platform=linux/arm64
ko build -B -L ./cartsvc --platform=linux/arm64
```

### Deploy

```bash
# Apply Dapr state store component (with actorStateStore: true)
kubectl apply -f manifest/state-store.yaml

# Build and deploy services
kubectl apply -f manifest/cart.yaml
kubectl apply -f manifest/frontend.yaml
```

### Test

```bash
# Add items
curl -s -X POST http://192.168.147.254/cart/alice/items \
  -d '{"id":"sku-100","name":"Keyboard","price":79.99,"quantity":1}' | jq

curl -s -X POST http://192.168.147.254/cart/alice/items \
  -d '{"id":"sku-200","name":"Monitor","price":349.00,"quantity":1}' | jq

# View cart
curl -s http://192.168.147.254/cart/alice | jq

# Checkout
curl -s -X POST http://192.168.107.254/cart/alice/checkout | jq

# Cart is now empty
curl -s http://192.168.107.254/cart/alice | jq
```

## How the Actor Timer Works

Every time the cart is modified, a 30-minute timer (`cart-expiry`) is registered. If no further activity happens, the timer fires `ExpireCart` which clears the cart — mimicking real-world abandoned cart behavior. Each modification resets the timer.
