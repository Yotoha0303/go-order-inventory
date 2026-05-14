package service_test

import (
	"errors"
	"go-order-inventory/global"
	"go-order-inventory/internal/model"
	"go-order-inventory/internal/request"
	"go-order-inventory/internal/service"
	"testing"
)

func TestCreateProduct_InvalidPrice(t *testing.T) {
	setupTestDB(t)

	req := request.CreateProductRequest{
		Name:        "test product",
		Description: "desc",
		PriceFen:    0,
	}

	product, err := service.CreateProduct(req)
	if !errors.Is(err, service.ErrInvalidProductPrice) {
		t.Fatalf("expected ErrInvalidProductPrice, got err=%v", err)
	}
	if product != nil {
		t.Fatalf("expected nil product, got %+v", product)
	}
}

func TestCreateProduct_SuccessTrimmed(t *testing.T) {
	setupTestDB(t)

	req := request.CreateProductRequest{
		Name:        "  apple  ",
		Description: "  good  ",
		PriceFen:    199,
	}

	p, err := service.CreateProduct(req)
	if err != nil {
		t.Fatalf("create product failed: %v", err)
	}
	if p.Name != "apple" || p.Description != "good" {
		t.Fatalf("trim failed, got name=%q description=%q", p.Name, p.Description)
	}
	if p.Status != model.ProductStatusOffSale {
		t.Fatalf("unexpected status: %d", p.Status)
	}
}

func TestListProducts_OnlyOffSale(t *testing.T) {
	setupTestDB(t)

	seedProduct(t, "off-sale", 100, model.ProductStatusOffSale)
	seedProduct(t, "on-sale", 100, model.ProductStatusOnSale)

	products, err := service.ListProducts()
	if err != nil {
		t.Fatalf("list products failed: %v", err)
	}

	if len(products) != 1 {
		t.Fatalf("expected 1 off-sale product, got %d", len(products))
	}
	if products[0].Status != model.ProductStatusOffSale {
		t.Fatalf("unexpected status: %d", products[0].Status)
	}
}

func TestGetProductByID_NotFound(t *testing.T) {
	setupTestDB(t)

	_, err := service.GetProductByID(99999)
	if !errors.Is(err, service.ErrProductNotFound) {
		t.Fatalf("expected ErrProductNotFound, got %v", err)
	}
}

func TestOnSaleProduct_Success(t *testing.T) {
	setupTestDB(t)

	p := seedProduct(t, "p1", 100, model.ProductStatusOffSale)
	if err := service.OnSaleProduct(p.ID); err != nil {
		t.Fatalf("on sale failed: %v", err)
	}

	var got model.Product
	if err := global.DB.First(&got, p.ID).Error; err != nil {
		t.Fatalf("query product failed: %v", err)
	}
	if got.Status != model.ProductStatusOnSale {
		t.Fatalf("expected on-sale status, got %d", got.Status)
	}
}

func TestOffSaleProduct_Success(t *testing.T) {
	setupTestDB(t)

	p := seedProduct(t, "p1", 100, model.ProductStatusOnSale)
	if err := service.OffSaleProduct(p.ID); err != nil {
		t.Fatalf("off sale failed: %v", err)
	}

	var got model.Product
	if err := global.DB.First(&got, p.ID).Error; err != nil {
		t.Fatalf("query product failed: %v", err)
	}
	if got.Status != model.ProductStatusOffSale {
		t.Fatalf("expected off-sale status, got %d", got.Status)
	}
}
