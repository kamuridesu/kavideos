package bot

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/go-telegram/bot"
)

type progress struct {
	MessageId     int
	LastTimestamp time.Time
	chatId        int64
	b             *bot.Bot
	ctx           context.Context
}

func progressStart(ctx context.Context, b *bot.Bot, chatId int64) *progress {

	m, _ := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatId,
		Text:   "Downloading... 0%",
	})
	return &progress{
		MessageId:     m.ID,
		LastTimestamp: time.Now(),
		chatId:        chatId,
		b:             b,
		ctx:           ctx,
	}

}

func (p *progress) update(total, current int) {
	if time.Since(p.LastTimestamp).Seconds() < 5 {
		return
	}
	p.LastTimestamp = time.Now()
	percent := int(math.Round(100 * float64(current) / float64(total)))
	p.b.EditMessageText(p.ctx, &bot.EditMessageTextParams{
		ChatID:    p.chatId,
		MessageID: p.MessageId,
		Text:      fmt.Sprintf("Downloading... %d%%", percent),
	})
}

func (p *progress) finishing() {
	p.b.EditMessageText(p.ctx, &bot.EditMessageTextParams{
		ChatID:    p.chatId,
		MessageID: p.MessageId,
		Text:      "Sending media...",
	})
}

func (p *progress) end() {
	p.b.DeleteMessage(p.ctx, &bot.DeleteMessageParams{
		ChatID:    p.chatId,
		MessageID: p.MessageId,
	})
}
