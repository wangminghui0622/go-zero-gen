package models

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ OrderRecordModel = (*customOrderRecordModel)(nil)

type (
	// OrderRecordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customOrderRecordModel.
	OrderRecordModel interface {
		orderRecordModel
	}

	customOrderRecordModel struct {
		*defaultOrderRecordModel
	}
)

// NewOrderRecordModel returns a model for the database table.
func NewOrderRecordModel(conn sqlx.SqlConn, c cache.CacheConf) OrderRecordModel {
	return &customOrderRecordModel{
		defaultOrderRecordModel: newOrderRecordModel(conn, c),
	}
}
