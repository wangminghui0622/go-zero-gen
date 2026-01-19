package models

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ OrderSrcModel = (*customOrderSrcModel)(nil)

type (
	// OrderSrcModel is an interface to be customized, add more methods here,
	// and implement the added methods in customOrderSrcModel.
	OrderSrcModel interface {
		orderSrcModel
	}

	customOrderSrcModel struct {
		*defaultOrderSrcModel
	}
)

// NewOrderSrcModel returns a model for the database table.
func NewOrderSrcModel(conn sqlx.SqlConn, c cache.CacheConf) OrderSrcModel {
	return &customOrderSrcModel{
		defaultOrderSrcModel: newOrderSrcModel(conn, c),
	}
}
