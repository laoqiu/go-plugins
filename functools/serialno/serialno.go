package serialno

import (
	"fmt"
	"time"

	"github.com/laoqiu/go-plugins/functools/random"
	"github.com/rs/xid"
)

// 流水号
func NewSerialNo() string {
	return fmt.Sprintf("%s%d", time.Now().Format("060102150405"), xid.New().Counter())
}

// 订单号
func NewOrderNo(prefix string) string {
	return prefix + time.Now().Format("060102150405") + random.RandDigitsStr(5)
}
