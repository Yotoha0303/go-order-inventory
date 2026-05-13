package service_test

import (
	"errors"
	"go-order-inventory/internal/request"
	"go-order-inventory/internal/service"
	"testing"
)

func TestCreateProduct_InvalidPrice(t *testing.T) {
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

func TestListProduct(t *testing.T) {
	//0、先插入Product数据

	//1、再调用service层的ListProduct方法

	//2、判断数据是否相等，错误打印结果
}

func TestGetProductByID(t *testing.T) {

}

func TestOnSaleProduct(t *testing.T) {

}

func TestOffSaleProduct(t *testing.T) {

}
