package bizcache

import (
	"context"
	"encoding/json"
	"fmt"
	"go-order-inventory/global"
	"go-order-inventory/internal/model"
	"log"
	"time"
)

const ProductDetailCacheTTL = 10 * time.Minute

func ProductDetailCacheKey(productID int64) string {
	return fmt.Sprintf("product:detail:%d", productID)
}

func GetProductDetail(ctx context.Context, productID int64) (*model.Product, bool) {
	if global.Redis == nil {
		return nil, false
	}

	ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	val, err := global.Redis.Get(ctx, ProductDetailCacheKey(productID)).Result()
	if err != nil {
		log.Printf("get product cache failed: product_id=%d err=%v", productID, err)
		return nil, false
	}

	var product model.Product
	if err := json.Unmarshal([]byte(val), &product); err != nil {
		return nil, false
	}

	return &product, true
}

func SetProductDetail(ctx context.Context, product *model.Product) {
	if global.Redis == nil || product == nil {
		return
	}

	data, err := json.Marshal(product)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	_ = global.Redis.Set(ctx, ProductDetailCacheKey(product.ID), data, ProductDetailCacheTTL).Err()
}

func DeleteProductDetailCache(ctx context.Context, productID int64) {
	if global.Redis == nil {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	_ = global.Redis.Del(ctx, ProductDetailCacheKey(productID)).Err()
}
