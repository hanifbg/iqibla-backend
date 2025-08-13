# iQibla Backend API Documentation

Complete API documentation for all endpoints in the iQibla Backend application.

## Base URL
```
http://localhost:8080
```

## Table of Contents
1. [Product APIs](#product-apis)
2. [Category APIs](#category-apis)
3. [Cart APIs](#cart-apis)
4. [Shipping APIs](#shipping-apis)
5. [Payment APIs](#payment-apis)

---

## Product APIs

### Get All Products

Get a list of all active products.

- **URL**: `/api/v1/products`
- **Method**: `GET`
- **Success Response**:
  - **Code**: 200
  - **Content**:
    ```json
    {
      "message": "Products retrieved successfully",
      "data": [
        {
          "id": 1,
          "name": "Product Name",
          "description": "Product description",
          "category_id": 1,
          "category": {
            "id": 1,
            "name": "Category Name",
            "slug": "category-slug"
          },
          "variants": [
            {
              "id": 1,
              "product_id": 1,
              "name": "Variant Name",
              "price": 150000,
              "stock": 10,
              "weight": 500,
              "images": ["image1.jpg", "image2.jpg"]
            }
          ]
        }
      ]
    }
    ```

### Get Product by ID

Get a specific product by its ID.

- **URL**: `/api/v1/products/:id`
- **Method**: `GET`
- **URL Parameters**:
  - `id`: Product ID
- **Success Response**:
  - **Code**: 200
  - **Content**: Same as single product object above
- **Error Response**:
  - **Code**: 404
  - **Content**:
    ```json
    {
      "error": "Product not found"
    }
    ```

---

## Category APIs

### Get Category by Slug

Get a category by its slug.

- **URL**: `/api/v1/categories/:slug`
- **Method**: `GET`
- **URL Parameters**:
  - `slug`: Category slug
- **Success Response**:
  - **Code**: 200
  - **Content**:
    ```json
    {
      "message": "Category retrieved successfully",
      "data": {
        "id": 1,
        "name": "Category Name",
        "slug": "category-slug",
        "description": "Category description"
      }
    }
    ```
- **Error Response**:
  - **Code**: 404
  - **Content**:
    ```json
    {
      "error": "Category not found"
    }
    ```

---

## Cart APIs

### Add Item to Cart

Add a product variant to the cart with specified quantity.

- **URL**: `/api/v1/cart/add`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "cart_id": "optional-cart-id",
    "variant_id": "123",
    "quantity": 2
  }
  ```
- **Success Response**:
  - **Code**: 200
  - **Content**:
    ```json
    {
      "cart_id": "cart-uuid",
      "items": [
        {
          "variant_id": "123",
          "product_name": "Product Name",
          "variant_name": "Variant Name",
          "price": 150000,
          "quantity": 2,
          "subtotal": 300000,
          "weight": 1000,
          "image": "product-image.jpg"
        }
      ],
      "total_amount": 300000,
      "total_weight": 1000,
      "total_items": 1,
      "discount": {
        "code": "",
        "amount": 0,
        "percentage": 0
      }
    }
    ```
- **Error Response**:
  - **Code**: 400
  - **Content**:
    ```json
    {
      "error": "Invalid request format"
    }
    ```
  - **Code**: 404
  - **Content**:
    ```json
    {
      "error": "Product variant not found"
    }
    ```

### Update Item Quantity

Update the quantity of an item in the cart.

- **URL**: `/api/v1/cart/update-quantity`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "cart_id": "cart-uuid",
    "variant_id": "123",
    "quantity": 3
  }
  ```
- **Success Response**: Same as Add Item response
- **Error Response**: Same as Add Item response

### Remove Item from Cart

Remove an item from the cart.

- **URL**: `/api/v1/cart/remove`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "cart_id": "cart-uuid",
    "variant_id": "123"
  }
  ```
- **Success Response**: Same as Add Item response
- **Error Response**: Same as Add Item response

### Get Cart Details

Get cart details by cart ID.

- **URL**: `/api/v1/cart/:cart_id`
- **Method**: `GET`
- **URL Parameters**:
  - `cart_id`: Cart UUID
- **Success Response**: Same as Add Item response
- **Error Response**:
  - **Code**: 404
  - **Content**:
    ```json
    {
      "error": "Cart not found"
    }
    ```

### Apply Discount

Apply a discount code to the cart.

- **URL**: `/api/v1/cart/apply-discount`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "cart_id": "cart-uuid",
    "discount_code": "DISCOUNT10"
  }
  ```
- **Success Response**: Same as Add Item response with updated discount information
- **Error Response**:
  - **Code**: 400
  - **Content**:
    ```json
    {
      "error": "Invalid discount code"
    }
    ```

---

## Shipping APIs

### Get Provinces

Get a list of all provinces or a specific province by ID.

- **URL**: `/api/v1/shipping/provinces`
- **Method**: `GET`
- **Query Parameters**:
  - `id` (optional): Province ID to get a specific province
- **Success Response**:
  - **Code**: 200
  - **Content**:
    ```json
    {
      "message": "Provinces retrieved successfully",
      "data": [
        {
          "province_id": "1",
          "province": "NUSA TENGGARA BARAT (NTB)"
        }
      ]
    }
    ```
- **Error Response**:
  - **Code**: 500
  - **Content**:
    ```json
    {
      "error": "Failed to get provinces",
      "message": "Error details"
    }
    ```

### Get Cities

Get cities by province ID and optionally filter by city ID.

- **URL**: `/api/v1/shipping/cities/:province_id`
- **Method**: `GET`
- **URL Parameters**:
  - `province_id`: ID of the province
- **Query Parameters**:
  - `id` (optional): City ID to get a specific city
- **Success Response**:
  - **Code**: 200
  - **Content**:
    ```json
    {
      "message": "Cities retrieved successfully",
      "data": [
        {
          "city_id": "55",
          "province_id": "5",
          "province": "",
          "type": "",
          "city_name": "BANDUNG",
          "postal_code": ""
        }
      ]
    }
    ```
- **Error Response**:
  - **Code**: 500
  - **Content**:
    ```json
    {
      "error": "Failed to get cities",
      "message": "Error details"
    }
    ```

### Get Districts

Get districts by city ID.

- **URL**: `/api/v1/shipping/districts/:city_id`
- **Method**: `GET`
- **URL Parameters**:
  - `city_id`: ID of the city
- **Success Response**:
  - **Code**: 200
  - **Content**:
    ```json
    {
      "message": "Districts retrieved successfully",
      "data": [
        {
          "district_id": "423",
          "city_id": "55",
          "city": "",
          "district_name": "BANDUNG",
          "type": ""
        }
      ]
    }
    ```
- **Error Response**:
  - **Code**: 500
  - **Content**:
    ```json
    {
      "error": "Failed to get districts",
      "message": "Error details"
    }
    ```

### Calculate Shipping Cost

Calculate shipping cost based on origin, destination, weight, and courier.

- **URL**: `/api/v1/shipping/cost`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "origin": "501",
    "destination": "114",
    "weight": 1000,
    "courier": "jne"
  }
  ```
- **Success Response**:
  - **Code**: 200
  - **Content**:
    ```json
    {
      "message": "Shipping cost calculated successfully",
      "data": [
        {
          "service": "JTR",
          "description": "JNE Trucking",
          "cost": 220000,
          "etd": "10 day"
        }
      ]
    }
    ```
- **Error Response**:
  - **Code**: 400
  - **Content**:
    ```json
    {
      "error": "Invalid request",
      "message": "Error details"
    }
    ```
  - **Code**: 500
  - **Content**:
    ```json
    {
      "error": "Failed to calculate shipping cost",
      "message": "Error details"
    }
    ```

### Validate and Save AWB Number

Validate AWB (Air Way Bill) number with RajaOngkir API and save it to database for order tracking.

- **URL**: `/api/v1/shipping/awb/validate`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "invoice_number": "IQB-2025-00001",
    "awb_number": "JX3948927748",
    "courier": "jnt",
    "last_phone_number": "12345"
  }
  ```
- **Request Body Parameters**:
  - `invoice_number` (required): Invoice number to link AWB to an order
  - `awb_number` (required): AWB tracking number from courier
  - `courier` (required): Courier service name (jne, jnt, ninja, tiki, pos, anteraja, sicepat, sap, lion, wahana, first, ide)
  - `last_phone_number` (optional): Last 5 digits of recipient's phone number. **Only required for JNE courier** for additional verification (exactly 5 numeric characters)
- **Success Response** (Valid AWB):
  - **Code**: 200
  - **Content**:
    ```json
    {
      "message": "AWB validated and saved successfully",
      "data": {
        "id": "46d6774d-22b4-48c9-8705-d26e71f4cf84",
        "invoice_number": "IQB-2025-00001",
        "awb_number": "JX3948927748",
        "courier": "jnt",
        "is_validated": true,
        "message": "AWB number validated and saved successfully"
      }
    }
    ```
- **Success Response** (Invalid AWB):
  - **Code**: 400
  - **Content**:
    ```json
    {
      "error": "Invalid AWB number",
      "message": "Invalid AWB number",
      "data": {
        "awb_number": "JX3948927748",
        "courier": "jnt",
        "invoice_number": "IQB-2025-00001",
        "is_validated": false,
        "message": "Invalid AWB number"
      }
    }
    ```
- **Error Responses**:
  - **Code**: 400 (Validation Failed)
  - **Content**:
    ```json
    {
      "error": "Validation failed",
      "message": "Key: 'ValidateAWBRequest.LastPhoneNumber' Error:Field validation for 'LastPhoneNumber' failed on the 'len' tag"
    }
    ```
  - **Code**: 400 (Invalid Request)
  - **Content**:
    ```json
    {
      "error": "Invalid request",
      "message": "Error details"
    }
    ```
  - **Code**: 500 (Server Error)
  - **Content**:
    ```json
    {
      "error": "Failed to validate AWB",
      "message": "Error details"
    }
    ```

**Notes**:
- The endpoint validates the AWB number using RajaOngkir API before saving
- Each AWB number + courier combination must be unique
- The invoice number must exist in the system
- `last_phone_number` is **only required for JNE courier** and must contain exactly the last 5 digits of the recipient's phone number
- For other couriers (JNT, Ninja, Tiki, etc.), the `last_phone_number` parameter should be omitted
- The system will return appropriate error messages for duplicate AWB numbers, invalid invoice numbers, or API validation failures

---

## Payment APIs

### Create Order

Create a new order from cart.

- **URL**: `/api/v1/orders`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "cart_id": "cart-uuid",
    "customer_name": "John Doe",
    "customer_email": "john@example.com",
    "customer_phone": "+6281234567890",
    "shipping_address": "Jl. Example No. 123",
    "shipping_city_name": "BANDUNG",
    "shipping_province_name": "JAWA BARAT",
    "shipping_district_name": "BANDUNG KULON",
    "shipping_postal_code": "40123",
    "shipping_courier": "jne",
    "shipping_service": "REG",
    "shipping_cost": 65000,
    "total_weight": 1000,
    "notes": "Optional notes"
  }
  ```
- **Success Response**:
  - **Code**: 200
  - **Content**:
    ```json
    {
      "message": "Order created successfully",
      "data": {
        "order_id": "order-uuid",
        "customer_name": "John Doe",
        "customer_email": "john@example.com",
        "customer_phone": "+6281234567890",
        "shipping_address": "Jl. Example No. 123",
        "shipping_city_name": "BANDUNG",
        "shipping_province_name": "JAWA BARAT",
        "shipping_district_name": "BANDUNG KULON",
        "shipping_postal_code": "40123",
        "subtotal": 300000,
        "shipping_cost": 65000,
        "discount_amount": 0,
        "total_amount": 365000,
        "status": "pending_payment",
        "items": [
          {
            "product_name": "Product Name",
            "variant_name": "Variant Name",
            "quantity": 2,
            "price": 150000,
            "subtotal": 300000
          }
        ]
      }
    }
    ```
- **Error Response**:
  - **Code**: 400
  - **Content**:
    ```json
    {
      "error": "Validation failed",
      "message": "Error details"
    }
    ```

### Get Order Details

Get order details by order ID.

- **URL**: `/api/v1/orders/:order_id`
- **Method**: `GET`
- **URL Parameters**:
  - `order_id`: Order UUID
- **Success Response**: Same as Create Order response
- **Error Response**:
  - **Code**: 404
  - **Content**:
    ```json
    {
      "error": "Order not found"
    }
    ```

### Create Payment

Create payment for an order.

- **URL**: `/api/v1/payments/:order_id`
- **Method**: `POST`
- **URL Parameters**:
  - `order_id`: Order UUID
- **Success Response**:
  - **Code**: 200
  - **Content**:
    ```json
    {
      "message": "Payment created successfully",
      "data": {
        "payment_id": "payment-uuid",
        "order_id": "order-uuid",
        "amount": 365000,
        "status": "pending",
        "payment_url": "https://app.sandbox.midtrans.com/snap/v2/vtweb/...",
        "snap_token": "snap-token"
      }
    }
    ```
- **Error Response**:
  - **Code**: 404
  - **Content**:
    ```json
    {
      "error": "Order not found"
    }
    ```

### Get Payment Status

Get payment status by payment ID.

- **URL**: `/api/v1/payments/status/:payment_id`
- **Method**: `GET`
- **URL Parameters**:
  - `payment_id`: Payment UUID
- **Success Response**:
  - **Code**: 200
  - **Content**:
    ```json
    {
      "message": "Payment status retrieved successfully",
      "data": {
        "payment_id": "payment-uuid",
        "order_id": "order-uuid",
        "status": "paid",
        "amount": 365000,
        "paid_at": "2025-07-29T14:30:00Z"
      }
    }
    ```
- **Error Response**:
  - **Code**: 404
  - **Content**:
    ```json
    {
      "error": "Payment not found"
    }
    ```

### Handle Payment Notification

Handle payment notification from Midtrans.

- **URL**: `/api/v1/payments/notification`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "transaction_time": "2025-07-29 14:30:00",
    "transaction_status": "settlement",
    "transaction_id": "midtrans-transaction-id",
    "status_code": "200",
    "signature_key": "signature",
    "payment_type": "credit_card",
    "order_id": "order-uuid",
    "merchant_id": "merchant-id",
    "gross_amount": "365000.00",
    "fraud_status": "accept",
    "currency": "IDR"
  }
  ```
- **Success Response**:
  - **Code**: 200
  - **Content**:
    ```json
    {
      "message": "Notification processed successfully"
    }
    ```
- **Error Response**:
  - **Code**: 400
  - **Content**:
    ```json
    {
      "error": "Invalid notification"
    }
    ```

---

## Static Files

### Product Images

Access product images via static file serving.

- **URL**: `/uploads/:category/:filename`
- **Method**: `GET`
- **URL Parameters**:
  - `category`: Product category folder (e.g., `jood_pro`, `lite`, `noor`)
  - `filename`: Image filename (e.g., `Jood-Pro-1.png`)
- **Example**: `http://localhost:8080/uploads/jood_pro/Jood-Pro-1.png`

---

## Error Codes

- `200`: Success
- `400`: Bad Request - Invalid request format or validation failed
- `404`: Not Found - Resource not found
- `500`: Internal Server Error - Server error

## Authentication

Currently, the API does not require authentication. All endpoints are publicly accessible.

## Rate Limiting

No rate limiting is currently implemented.

## CORS

CORS is enabled for all origins with the following configuration:
- **Allowed Origins**: `*`
- **Allowed Methods**: `GET`, `POST`, `PUT`, `DELETE`, `OPTIONS`
- **Allowed Headers**: `Origin`, `Content-Type`, `Accept`, `Authorization`
- **Allow Credentials**: `true`
