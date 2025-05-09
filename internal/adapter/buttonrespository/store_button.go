package buttonrespository

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/domain/button"
)

func (r *ButtonRepository) StoreButton(ctx context.Context, btn button.Button) error {
	encodedButton, err := encodeButton(btn)
	if err != nil {
		return fmt.Errorf("encode button: %w", err)
	}

	if err := r.client.Set(ctx, btn.ID, encodedButton, r.ttl).Err(); err != nil {
		return fmt.Errorf("redis set: %w", err)
	}

	return nil
}
