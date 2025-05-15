package util

import (
	"context"
	"encoding/json"
	"fmt"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkcardkit "github.com/larksuite/oapi-sdk-go/v3/service/cardkit/v1"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

// CreateCard 创建互动卡片并返回卡片ID
func CreateCard(client *lark.Client) (string, error) {
	cardData := map[string]interface{}{
		"schema": "2.0",
		"config": map[string]interface{}{
			"streaming_mode": true,
			"summary":        map[string]interface{}{"content": "..."},
		},
		"body": map[string]interface{}{
			"direction":          "vertical",
			"horizontal_spacing": "8px",
			"vertical_spacing":   "8px",
			"horizontal_align":   "left",
			"vertical_align":     "top",
			"padding":            "12px 12px 12px 12px",
			"elements": []interface{}{
				map[string]interface{}{
					"tag":        "markdown",
					"content":    "",
					"text_align": "left",
					"text_size":  "normal",
					"margin":     "4px 4px 4px 4px",
					"element_id": "streaming_txt",
				},
			},
		},
	}

	dataBytes, err := json.Marshal(cardData)
	if err != nil {
		return "", fmt.Errorf("序列化卡片数据失败: %w", err)
	}

	req := larkcardkit.NewCreateCardReqBuilder().
		Body(larkcardkit.NewCreateCardReqBodyBuilder().
			Type("card_json").
			Data(string(dataBytes)).
			Build()).
		Build()

	resp, err := client.Cardkit.V1.Card.Create(context.Background(), req)
	if err != nil {
		return "", fmt.Errorf("创建卡片失败: %w", err)
	}
	if !resp.Success() {
		return "", fmt.Errorf("卡片创建响应失败: %s", larkcore.Prettify(resp.CodeError))
	}

	return *resp.Data.CardId, nil
}

// SendCard 发送卡片消息到指定会话
func SendCard(client *lark.Client, chatID, cardID string) error {
	content := fmt.Sprintf(`{"type":"card","data":{"card_id":"%s"}}`, cardID)

	req := larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(`chat_id`).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			ReceiveId(chatID).
			MsgType("interactive").
			Content(content).
			Build()).
		Build()

	resp, err := client.Im.V1.Message.Create(context.Background(), req)
	if err != nil {
		return fmt.Errorf("发送卡片请求失败: %w", err)
	}
	if !resp.Success() {
		return fmt.Errorf("发送卡片响应失败: %s", larkcore.Prettify(resp.CodeError))
	}

	return nil
}

// UpdateElementText 更新卡片中某个元素的内容
func UpdateElementText(client *lark.Client, cardID, elementID, newText string, sequence int) error {
	req := larkcardkit.NewContentCardElementReqBuilder().
		CardId(cardID).
		ElementId(elementID).
		Body(larkcardkit.NewContentCardElementReqBodyBuilder().
			Content(newText).
			Sequence(sequence).
			Build()).
		Build()

	resp, err := client.Cardkit.V1.CardElement.Content(context.Background(), req)
	if err != nil {
		return fmt.Errorf("更新元素失败: %w", err)
	}
	if !resp.Success() {
		return fmt.Errorf("更新元素响应失败: %s", larkcore.Prettify(resp.CodeError))
	}

	return nil
}
