// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/cart/add": {
            "post": {
                "description": "Adds a product variant to the cart with specified quantity",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "cart"
                ],
                "summary": "Add item to cart",
                "parameters": [
                    {
                        "description": "Add item request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.AddItemRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.CartResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/api/v1/cart/apply-discount": {
            "post": {
                "description": "Applies a discount code to the cart",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "cart"
                ],
                "summary": "Apply discount to cart",
                "parameters": [
                    {
                        "description": "Apply discount request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.ApplyDiscountRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.CartResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/api/v1/cart/remove": {
            "post": {
                "description": "Removes a product variant from the cart",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "cart"
                ],
                "summary": "Remove item from cart",
                "parameters": [
                    {
                        "description": "Remove item request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.RemoveItemRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.CartResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/api/v1/cart/update-quantity": {
            "post": {
                "description": "Updates the quantity of a product variant in the cart",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "cart"
                ],
                "summary": "Update cart item quantity",
                "parameters": [
                    {
                        "description": "Update quantity request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.UpdateItemRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.CartResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/api/v1/cart/{cart_id}": {
            "get": {
                "description": "Retrieves cart details including items and totals",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "cart"
                ],
                "summary": "Get cart details",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Cart ID",
                        "name": "cart_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.CartResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/api/v1/orders": {
            "post": {
                "description": "Create a new order from the items in the cart",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "orders"
                ],
                "summary": "Create a new order from cart",
                "parameters": [
                    {
                        "description": "Order details",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.CreateOrderRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.OrderResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/api/v1/orders/{order_id}": {
            "get": {
                "description": "Get details of an order by ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "orders"
                ],
                "summary": "Get order details",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Order ID",
                        "name": "order_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.OrderResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/api/v1/payments/notification": {
            "post": {
                "description": "Handle payment notification from payment gateway",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "payments"
                ],
                "summary": "Handle payment notification",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/api/v1/payments/status/{payment_id}": {
            "get": {
                "description": "Get the status of a payment by ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "payments"
                ],
                "summary": "Get payment status",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Payment ID",
                        "name": "payment_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.PaymentStatusResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/api/v1/payments/{order_id}": {
            "post": {
                "description": "Create a payment transaction for an order",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "payments"
                ],
                "summary": "Create payment for an order",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Order ID",
                        "name": "order_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.PaymentResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/api/v1/products": {
            "get": {
                "description": "Retrieves all active products with their variants",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "products"
                ],
                "summary": "Get all active products",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/entity.Product"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/api/v1/products/{id}": {
            "get": {
                "description": "Retrieves a single product by its ID with variants",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "products"
                ],
                "summary": "Get a product by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Product ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entity.Product"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "entity.Dimensions": {
            "type": "object",
            "properties": {
                "height": {
                    "type": "number"
                },
                "length": {
                    "type": "number"
                },
                "unit": {
                    "type": "string"
                },
                "width": {
                    "type": "number"
                }
            }
        },
        "entity.JSONMap": {
            "type": "object",
            "additionalProperties": true
        },
        "entity.PaymentStatus": {
            "type": "string",
            "enum": [
                "pending",
                "success",
                "failed",
                "expired",
                "cancelled",
                "refunded"
            ],
            "x-enum-varnames": [
                "PaymentStatusPending",
                "PaymentStatusSuccess",
                "PaymentStatusFailed",
                "PaymentStatusExpired",
                "PaymentStatusCancelled",
                "PaymentStatusRefunded"
            ]
        },
        "entity.Product": {
            "type": "object",
            "properties": {
                "brand": {
                    "type": "string"
                },
                "category": {
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "deleted_at": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "features": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "id": {
                    "type": "string"
                },
                "image_urls": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "in_box_items": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "is_active": {
                    "type": "boolean"
                },
                "name": {
                    "type": "string"
                },
                "shopee_url": {
                    "description": "Added",
                    "type": "string"
                },
                "tokopedia_url": {
                    "description": "Added",
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                },
                "variants": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/entity.ProductVariant"
                    }
                }
            }
        },
        "entity.ProductVariant": {
            "type": "object",
            "properties": {
                "attribute_values": {
                    "$ref": "#/definitions/entity.JSONMap"
                },
                "created_at": {
                    "type": "string"
                },
                "deleted_at": {
                    "type": "string"
                },
                "dimensions": {
                    "$ref": "#/definitions/entity.Dimensions"
                },
                "id": {
                    "type": "string"
                },
                "image_url": {
                    "type": "string"
                },
                "is_active": {
                    "type": "boolean"
                },
                "name": {
                    "type": "string"
                },
                "price": {
                    "type": "number"
                },
                "product_id": {
                    "type": "string"
                },
                "sku": {
                    "type": "string"
                },
                "specifications": {
                    "description": "e.g., {\"Display\": \"0.49 Inch, OLED\", \"Material\": \"Plastic\", \"Battery\": \"45mAh\"}",
                    "allOf": [
                        {
                            "$ref": "#/definitions/entity.JSONMap"
                        }
                    ]
                },
                "stock_quantity": {
                    "type": "integer"
                },
                "updated_at": {
                    "type": "string"
                },
                "weight": {
                    "type": "number"
                }
            }
        },
        "request.AddItemRequest": {
            "type": "object",
            "required": [
                "quantity",
                "variant_id"
            ],
            "properties": {
                "cart_id": {
                    "type": "string"
                },
                "quantity": {
                    "type": "integer",
                    "minimum": 1
                },
                "variant_id": {
                    "type": "string"
                }
            }
        },
        "request.ApplyDiscountRequest": {
            "type": "object",
            "required": [
                "cart_id",
                "discount_code"
            ],
            "properties": {
                "cart_id": {
                    "type": "string"
                },
                "discount_code": {
                    "type": "string"
                }
            }
        },
        "request.CreateOrderRequest": {
            "type": "object",
            "required": [
                "cart_id",
                "customer_email",
                "customer_name",
                "customer_phone",
                "shipping_address"
            ],
            "properties": {
                "cart_id": {
                    "type": "string"
                },
                "customer_email": {
                    "type": "string"
                },
                "customer_name": {
                    "type": "string"
                },
                "customer_phone": {
                    "type": "string"
                },
                "notes": {
                    "type": "string"
                },
                "shipping_address": {
                    "type": "string"
                }
            }
        },
        "request.RemoveItemRequest": {
            "type": "object",
            "required": [
                "cart_id",
                "variant_id"
            ],
            "properties": {
                "cart_id": {
                    "type": "string"
                },
                "variant_id": {
                    "type": "string"
                }
            }
        },
        "request.UpdateItemRequest": {
            "type": "object",
            "required": [
                "cart_id",
                "quantity",
                "variant_id"
            ],
            "properties": {
                "cart_id": {
                    "type": "string"
                },
                "quantity": {
                    "type": "integer",
                    "minimum": 0
                },
                "variant_id": {
                    "type": "string"
                }
            }
        },
        "response.CartItemResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "image_url": {
                    "type": "string"
                },
                "product_attributes": {
                    "type": "object",
                    "additionalProperties": true
                },
                "quantity": {
                    "type": "integer"
                },
                "variant_id": {
                    "type": "string"
                },
                "variant_name": {
                    "type": "string"
                },
                "variant_price": {
                    "type": "number"
                }
            }
        },
        "response.CartResponse": {
            "type": "object",
            "properties": {
                "cart_id": {
                    "type": "string"
                },
                "discount_amount": {
                    "type": "number"
                },
                "discount_code_applied": {
                    "type": "string"
                },
                "items": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/response.CartItemResponse"
                    }
                },
                "subtotal_amount": {
                    "type": "number"
                },
                "total_items": {
                    "type": "integer"
                }
            }
        },
        "response.OrderItemResponse": {
            "type": "object",
            "properties": {
                "attributes": {
                    "type": "object",
                    "additionalProperties": true
                },
                "id": {
                    "type": "string"
                },
                "image_url": {
                    "type": "string"
                },
                "product_variant_id": {
                    "type": "string"
                },
                "quantity": {
                    "type": "integer"
                },
                "subtotal": {
                    "type": "number"
                },
                "unit_price": {
                    "type": "number"
                },
                "variant_name": {
                    "type": "string"
                }
            }
        },
        "response.OrderResponse": {
            "type": "object",
            "properties": {
                "cart_id": {
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "customer_email": {
                    "type": "string"
                },
                "customer_name": {
                    "type": "string"
                },
                "customer_phone": {
                    "type": "string"
                },
                "discount_amount": {
                    "type": "number"
                },
                "discount_code": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "items": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/response.OrderItemResponse"
                    }
                },
                "notes": {
                    "type": "string"
                },
                "order_status": {
                    "type": "string"
                },
                "shipping_address": {
                    "type": "string"
                },
                "shipping_cost": {
                    "type": "number"
                },
                "subtotal": {
                    "type": "number"
                },
                "total_amount": {
                    "type": "number"
                }
            }
        },
        "response.PaymentResponse": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "number"
                },
                "created_at": {
                    "type": "string"
                },
                "expiry_time": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "order_id": {
                    "type": "string"
                },
                "payment_method": {
                    "type": "string"
                },
                "payment_token": {
                    "type": "string"
                },
                "payment_url": {
                    "type": "string"
                },
                "status": {
                    "$ref": "#/definitions/entity.PaymentStatus"
                },
                "transaction_id": {
                    "type": "string"
                }
            }
        },
        "response.PaymentStatusResponse": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "number"
                },
                "id": {
                    "type": "string"
                },
                "order_id": {
                    "type": "string"
                },
                "payment_method": {
                    "type": "string"
                },
                "status": {
                    "$ref": "#/definitions/entity.PaymentStatus"
                },
                "transaction_id": {
                    "type": "string"
                },
                "transaction_time": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
