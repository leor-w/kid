package apple

import (
	"context"
	"fmt"
)

func (apple *AppleStore) SendConsumptionInformation(txId string, params *ConsumptionParams) error {
	bodyMap := make(map[string]interface{})
	if params.AccountTenure != nil {
		bodyMap["accountTenure"] = params.AccountTenure
	} else {
		return fmt.Errorf("apple: 账户年限不能为空")
	}
	if params.AppAccountToken != nil {
		bodyMap["appAccountToken"] = params.AppAccountToken
	} else {
		return fmt.Errorf("apple: 应用内用户帐户令牌不能为空")
	}
	if params.ConsumptionStatus != nil {
		bodyMap["consumptionStatus"] = params.ConsumptionStatus
	} else {
		return fmt.Errorf("apple: 消费状态不能为空")
	}
	if params.CustomerConsented != nil {
		bodyMap["customerConsented"] = params.CustomerConsented
	} else {
		return fmt.Errorf("apple: 客户是否同意消费不能为空")
	}
	if params.DeliveryStatus != nil {
		bodyMap["deliveryStatus"] = params.DeliveryStatus
	} else {
		return fmt.Errorf("apple: 交付状态不能为空")
	}
	if params.LifetimeDollarsPurchased != nil {
		bodyMap["lifetimeDollarsPurchased"] = params.LifetimeDollarsPurchased
	} else {
		return fmt.Errorf("apple: 终身消费金额不能为空")
	}
	if params.LifetimeDollarsRefunded != nil {
		bodyMap["lifetimeDollarsRefunded"] = params.LifetimeDollarsRefunded
	} else {
		return fmt.Errorf("apple: 终身退款金额不能为空")
	}
	if params.Platform != nil {
		bodyMap["platform"] = params.Platform
	} else {
		return fmt.Errorf("apple: 平台不能为空")
	}
	if params.PlayTime != nil {
		bodyMap["playTime"] = params.PlayTime
	} else {
		return fmt.Errorf("apple: 使用该应用程序的时间量的值不能为空")
	}
	if params.RefundPreference != nil {
		bodyMap["refundPreference"] = params.RefundPreference
	}
	if params.SampleContentProvided != nil {
		bodyMap["sampleContentProvided"] = params.SampleContentProvided
	} else {
		return fmt.Errorf("apple: 是否提供了内容的免费样本或试用版不能为空")
	}
	if params.UserStatus != nil {
		bodyMap["userStatus"] = params.UserStatus
	} else {
		return fmt.Errorf("apple: 客户账户的状态不能为空")
	}
	if err := apple.client.SendConsumptionInformation(context.Background(), txId, bodyMap); err != nil {
		return fmt.Errorf("apple: 发送消费信息失败: %w", err)
	}
	return nil
}
