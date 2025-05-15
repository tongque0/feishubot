package callback

import (
	"context"
	"feishu/config"
	"feishu/util"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/schema"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

// 接收消息事件
func MessageReceive(ctx context.Context, event *larkim.P2MessageReceiveV1) error {

	client := lark.NewClient(config.GetConf().AppID, config.GetConf().AppSecret)

	go func() {
		cardID, _ := util.CreateCard(client)
		err := util.SendCard(client, *event.Event.Message.ChatId, cardID)
		if err != nil {
			fmt.Printf("发送卡片失败: %v\n", err)
		}
		model, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
			APIKey:  config.GetConf().AiKey,
			Timeout: 38 * time.Second,
			Model:   config.GetConf().AiModel,
			BaseURL: config.GetConf().AiBaseurl,
		})
		if err != nil {
			log.Fatalf("new chat model failed: %v", err)
		}
		messages := []*schema.Message{
			schema.UserMessage(*event.Event.Message.Content),
		}
		streamResult, err := model.Stream(ctx, messages)
		if err != nil {
			log.Fatalf("stream failed: %v", err)
		}
		defer streamResult.Close()
		i := 0
		step := 10
		nextUpdate := step
		var fullText strings.Builder
		for {
			message, err := streamResult.Recv()
			if err == io.EOF {
				step++
				util.UpdateElementText(client, cardID, "streaming_txt", fullText.String(), step)
				break
			}
			if err != nil {
				log.Fatalf("recv failed: %v", err)
			}

			fullText.WriteString(message.Content)
			i++

			if i >= nextUpdate {
				step++
				util.UpdateElementText(client, cardID, "streaming_txt", fullText.String(), step)
				nextUpdate += step
			}
		}
	}()

	return nil
}
