package apple

import (
	"context"
	"fmt"

	gopay "github.com/go-pay/gopay/apple"
)

// VerifyReceipt 验证收据
func (apple *AppleStore) VerifyReceipt(receipt string) (*gopay.VerifyResponse, error) {
	verifyUrl := gopay.UrlProd
	if !apple.options.IsProduct {
		verifyUrl = gopay.UrlSandbox
	}
	resp, err := gopay.VerifyReceipt(context.Background(), verifyUrl, apple.options.SharedSecret, receipt)
	if err != nil {
		return nil, fmt.Errorf("apple: 验证收据失败: %w", err)
	}
	return resp, nil
}

// GetTransaction 获取交易信息
func (apple *AppleStore) GetTransaction(txId string) (*gopay.TransactionInfoRsp, error) {
	txInfo, err := apple.client.GetTransactionInfo(context.Background(), txId)
	if err != nil {
		return nil, fmt.Errorf("apple: 获取交易信息失败: %w", err)
	}
	return txInfo, nil
}

// GetTransactionHistory 获取交易历史
func (apple *AppleStore) GetTransactionHistory(txId string, query *TransactionHistoryQuery) (*gopay.TransactionHistoryRsp, error) {
	var bodyMap = make(map[string]interface{})
	if query != nil {
		if query.Revision != nil {
			bodyMap["revision"] = &query.Revision
		}
		if query.StartDate != nil {
			bodyMap["startDate"] = &query.StartDate
		}
		if query.EndDate != nil {
			bodyMap["endDate"] = &query.EndDate
		}
		if query.ProductId != nil {
			bodyMap["productId"] = &query.ProductId
		}
		if query.ProductType != nil {
			bodyMap["productType"] = &query.ProductType
		}
		if query.InAppOwnershipType != nil {
			bodyMap["inAppOwnershipType"] = &query.InAppOwnershipType
		}
		if query.Sort != nil {
			bodyMap["sort"] = &query.Sort
		}
		if query.Revoked != nil {
			bodyMap["revoked"] = &query.Revoked
		}
		if query.SubscriptionGroupIdentifier != nil {
			bodyMap["subscriptionGroupIdentifier"] = &query.SubscriptionGroupIdentifier
		}

	}
	resp, err := apple.client.GetTransactionHistory(context.Background(), txId, bodyMap)
	if err != nil {
		return nil, fmt.Errorf("apple: 获取交易历史失败: %w", err)
	}
	return resp, nil
}
