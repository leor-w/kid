package apple

import (
	"context"
	"fmt"

	gopay "github.com/go-pay/gopay/apple"
)

// GetRefundHistory 获取退款历史
func (apple *AppleStore) GetRefundHistory(txId string, revision string) (*gopay.RefundHistoryRsp, error) {
	resp, err := apple.client.GetRefundHistory(context.Background(), txId, revision)
	if err != nil {
		return nil, fmt.Errorf("apple: 获取退款历史失败: %w", err)
	}
	return resp, nil
}
