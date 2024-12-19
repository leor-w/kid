package apple

import (
	"context"
	"fmt"

	gopay "github.com/go-pay/gopay/apple"
)

// GetNotificationHistory 获取通知历史
func (apple *AppleStore) GetNotificationHistory(paginationToken string, body *NotificationHistoryQuery) (*gopay.NotificationHistoryRsp, error) {
	bodyMap := make(map[string]interface{})
	if body != nil {
		if body.StartDate != nil {
			bodyMap["startDate"] = &body.StartDate
		}
		if body.EndDate != nil {
			bodyMap["endDate"] = &body.EndDate
		}
		if body.NotificationType != nil {
			bodyMap["notificationType"] = &body.NotificationType
		}
		if body.NotificationSubtype != nil {
			bodyMap["notificationSubtype"] = &body.NotificationSubtype
		}
		if body.OnlyFailures != nil {
			bodyMap["onlyFailures"] = &body.OnlyFailures
		}
		if body.TransactionId != nil {
			bodyMap["transactionId"] = &body.TransactionId
		}
	}
	resp, err := apple.client.GetNotificationHistory(context.Background(), paginationToken, bodyMap)
	if err != nil {
		return nil, fmt.Errorf("apple: 获取通知历史失败: %w", err)
	}
	return resp, nil
}
