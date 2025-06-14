package service

import (
	"github.com/hanifbg/landing_backend/internal/model/request"
	"github.com/hanifbg/landing_backend/internal/model/response"
)

type CartService interface {
	AddItem(req request.AddItemRequest) (*response.CartResponse, error)
	UpdateItemQuantity(req request.UpdateItemRequest) (*response.CartResponse, error)
	RemoveItem(req request.RemoveItemRequest) (*response.CartResponse, error)
	GetCart(cartID string) (*response.CartResponse, error)
	ApplyDiscount(req request.ApplyDiscountRequest) (*response.CartResponse, error)
}
