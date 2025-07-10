
# Payment API Documentation

This document provides details on the payment-related API endpoints, including payloads and responses.

## Endpoints

### Payments

#### `POST /api/v1/payments/:order_id`

Create a payment transaction for an order.

**URL Parameters:**

*   `order_id` (string, required): The ID of the order.

**Success Response:** `response.PaymentResponse`

```json
{
  "id": "string",
  "order_id": "string",
  "amount": "float64",
  "status": "string",
  "payment_method": "string",
  "transaction_id": "string",
  "payment_token": "string",
  "payment_url": "string",
  "expiry_time": "time.Time",
  "created_at": "time.Time"
}
```

#### `GET /api/v1/payments/status/:payment_id`

Get the status of a payment by ID.

**URL Parameters:**

*   `payment_id` (string, required): The ID of the payment.

**Success Response:** `response.PaymentStatusResponse`

```json
{
  "id": "string",
  "order_id": "string",
  "status": "string",
  "transaction_id": "string",
  "transaction_time": "time.Time",
  "payment_method": "string",
  "amount": "float64",
  "updated_at": "time.Time"
}
```

#### `POST /api/v1/payments/notification`

Handle payment notification from the payment gateway.

**Request Body:** `request.PaymentNotificationRequest`

```json
{
  "transaction_time": "string",
  "transaction_status": "string",
  "transaction_id": "string",
  "status_code": "string",
  "signature_key": "string",
  "payment_type": "string",
  "order_id": "string",
  "merchant_id": "string",
  "gross_amount": "string",
  "fraud_status": "string",
  "currency": "string"
}
```

**Success Response:**

```json
{
  "status": "ok"
}
```
