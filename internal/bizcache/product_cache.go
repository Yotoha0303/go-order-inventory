package bizcache

import (
	"context"
	"encoding/json"
	"fmt"
	"go-order-inventory/global"
	"go-order-inventory/internal/model"
	"time"
)

const ProductDetailCacheTTL = 10 * time.Minute

func ProductDetailCacheKey(productID int64) string {
	return fmt.Sprintf("product:detail:%d", productID)
}

func GetProductDetail(productID int64) (*model.Product, bool) {
	if global.Redis == nil {
		return nil, false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	val, err := global.Redis.Get(ctx, ProductDetailCacheKey(productID)).Result()
	if err != nil {
		return nil, false
	}

	// Parse the cached value into a Product model
	var product model.Product
	if err := json.Unmarshal([]byte(val), &product); err != nil {
		return nil, false
	}

	return &product, true
}

func SetProductDetail(product *model.Product) {
	if global.Redis == nil || product == nil {
		return
	}

	data, err := json.Marshal(product)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Microsecond)
	defer cancel()

	_ = global.Redis.Set(ctx, ProductDetailCacheKey(product.ID), data, ProductDetailCacheTTL).Err()
}

func DeleteProductDetailCache(productID int64) {
	if global.Redis == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Microsecond)
	defer cancel()

	_ = global.Redis.Del(ctx, ProductDetailCacheKey(productID)).Err()
}
