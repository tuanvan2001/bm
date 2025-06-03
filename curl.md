# cURL Commands for Order API

## Create Order (POST /order)
```bash
curl -X POST \
  http://localhost:3000/order \
  -H 'Content-Type: application/json' \
  -d '{
    "order_code": "ORDER-001",
    "productId": "PRODUCT-A",
    "amount": 100,
    "status": "pending"
  }'
```

## Get All Orders (GET /order)
```bash
curl -X GET \
  http://localhost:3000/order
```

## Get Order by ID (GET /order/:id)
```bash
curl -X GET \
  http://localhost:3000/order/{id}
```

## Update Order (PATCH /order/:id)
```bash
curl -X PATCH \
  http://localhost:3000/order/{id} \
  -H 'Content-Type: application/json' \
  -d '{
    "status": "completed"
  }'
```

## Delete Order (DELETE /order/:id)
```bash
curl -X DELETE \
  http://localhost:3000/order/{id}
```