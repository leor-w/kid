package apple

import (
	"context"
	"fmt"

	gopay "github.com/go-pay/gopay/apple"
)

func (apple *AppleStore) LookUpOrderId(orderId string) (*gopay.LookUpOrderIdRsp, error) {
	resp, err := apple.client.LookUpOrderId(context.Background(), orderId)
	if err != nil {
		return nil, fmt.Errorf("apple: 获取订单号失败: %w", err)
	}
	return resp, nil
}
