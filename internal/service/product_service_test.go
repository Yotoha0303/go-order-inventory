package service_test

import (
	"context"
	"errors"
	"go-order-inventory/global"
	"go-order-inventory/internal/model"
	"go-order-inventory/internal/request"
	"go-order-inventory/internal/service"
	"strings"
	"testing"
)

func TestCreateProduct_Success(t *testing.T) {
	setupTestDB(t)

	req := request.CreateProductRequest{
		Name:        "test product",
		Description: "desc",
		PriceFen:    199,
	}
	product, err := service.CreateProduct(req)
	if err != nil {
		t.Fatalf("create product failed: %v", err)
	}

	if product.ID <= 0 {
		t.Fatalf("expected product ID > 0,got %d", product.ID)
	}

	if product.Name != req.Name {
		t.Fatalf("expected name %q,got %q", req.Name, product.Name)
	}

	if product.Description != req.Description {
		t.Fatalf("expected description %q,got %q", req.Description, product.Description)
	}

	if product.PriceFen != req.PriceFen {
		t.Fatalf("expected price %d,got %d", req.PriceFen, product.PriceFen)
	}

	var saved model.Product

	if err := global.DB.First(&saved, product.ID).Error; err != nil {
		t.Fatalf("query product record failed: %v", err)
	}

	if saved.Name != req.Name || saved.Description != req.Description || saved.PriceFen != req.PriceFen || saved.Status != model.ProductStatusOffSale {
		t.Fatalf("saved record mismatch, got %+v", saved)
	}

}

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

func TestCreateProduct_EmptyName(t *testing.T) {
	setupTestDB(t)

	req := request.CreateProductRequest{
		Name:        "",
		Description: "name is empty",
		PriceFen:    199,
	}

	product, err := service.CreateProduct(req)
	if !errors.Is(err, service.ErrInvalidProductName) {
		t.Fatalf("expected ErrInvalidProductName, got err=%v", err)
	}

	if product != nil {
		t.Fatalf("expected product nil, got %+v", product)
	}

	var count int64
	if err := global.DB.Model(&model.Product{}).Where("description = ? AND name = ? ", req.Description, req.Name).Count(&count).Error; err != nil {
		t.Fatalf("count products failed: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected 0 products, got %d", count)
	}
}

func TestCreateProduct_DescriptionTooLong(t *testing.T) {
	setupTestDB(t)

	req := request.CreateProductRequest{
		Name:        "description-too-long-product",
		Description: strings.Repeat("a", 501),
		PriceFen:    199,
	}

	product, err := service.CreateProduct(req)
	if !errors.Is(err, service.ErrInvalidProductDescription) {
		t.Fatalf("expect desciption less 500 character:,got %v", err)
	}

	if product != nil {
		t.Fatalf("expect product nil,got %+v", product)
	}

	var count int64
	if err := global.DB.Model(&model.Product{}).Where("name = ? AND price_fen = ?", req.Name, 199).Count(&count).Error; err != nil {
		t.Fatalf("count products failed: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected 0 products, got %d", count)
	}
}

func TestCreateProduct_DescriptionExactly500(t *testing.T) {
	setupTestDB(t)

	req := request.CreateProductRequest{
		Name:        "description exactly 500",
		Description: strings.Repeat("a", 500),
		PriceFen:    199,
	}

	product, err := service.CreateProduct(req)
	if err != nil {
		t.Fatalf("create product failed: %v", err)
	}

	if product == nil {
		t.Fatal("expected product not nil")
	}

	if product.Name != req.Name {
		t.Fatalf("expected name %q, got %q", req.Name, product.Name)
	}

	if product.PriceFen != req.PriceFen {
		t.Fatalf("expected price_fen %d, got %d", req.PriceFen, product.PriceFen)
	}

	if product.Description != req.Description {
		t.Fatalf("expected description %q, got %q", req.Description, product.Description)
	}

	if len(product.Description) != 500 {
		t.Fatalf("expected description length 500, got %d", len(product.Description))
	}

	if product.Status != model.ProductStatusOffSale {
		t.Fatalf("expected status off-sale, got %d", product.Status)
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

	_, err := service.GetProductByID(context.Background(), 99999)
	if !errors.Is(err, service.ErrProductNotFound) {
		t.Fatalf("expected ErrProductNotFound, got %v", err)
	}
}

func TestOnSaleProduct_Success(t *testing.T) {
	setupTestDB(t)

	p := seedProduct(t, "p1", 100, model.ProductStatusOffSale)
	if err := service.OnSaleProduct(context.Background(), p.ID); err != nil {
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
	if err := service.OffSaleProduct(context.Background(), p.ID); err != nil {
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
