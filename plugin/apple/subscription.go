package apple

import (
	"context"
	"fmt"

	gopay "github.com/go-pay/gopay/apple"
)

// GetAllSubscriptionStatuses 获取所有订阅状态
func (apple *AppleStore) GetAllSubscriptionStatuses(txId string) (*gopay.AllSubscriptionStatusesRsp, error) {
	resp, err := apple.client.GetAllSubscriptionStatuses(context.Background(), txId)
	if err != nil {
		return nil, fmt.Errorf("apple: 获取订阅状态失败: %w", err)
	}
	return resp, nil
}
