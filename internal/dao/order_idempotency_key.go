package dao

import (
	"context"
	"go-order-inventory/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func TryCreateOrderIdempotencyKey(db *gorm.DB, ctx context.Context, idempotencyKey, requestHash string) (bool, error) {
	orderIdempotencyKey := &model.OrderIdempotencyKey{
		IdempotencyKey: idempotencyKey,
		RequestHash:    requestHash,
		Status:         model.OrderBeingCreated,
	}
	result := db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(orderIdempotencyKey)
	return result.RowsAffected == 1, result.Error
}

func GetOrderIdempotencyKey(db *gorm.DB, ctx context.Context, idempotencyKey string) (*model.OrderIdempotencyKey, error) {
	var orderIdempotencyKey model.OrderIdempotencyKey
	if err := db.WithContext(ctx).Where("idempotency_key = ?", idempotencyKey).First(&orderIdempotencyKey).Error; err != nil {
		return nil, err
	}
	return &orderIdempotencyKey, nil
}

func CompleteOrderIdempotencyKey(db *gorm.DB, ctx context.Context, idempotencyKey string, orderID int64) (int64, error) {
	result := db.WithContext(ctx).
		Model(&model.OrderIdempotencyKey{}).
		Where("idempotency_key = ? AND status = ?", idempotencyKey, model.OrderBeingCreated).
		Updates(map[string]interface{}{
			"order_id": orderID,
			"status":   model.OrderAlreadyCreated,
		})
	return result.RowsAffected, result.Error
}
