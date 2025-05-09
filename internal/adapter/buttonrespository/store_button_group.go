package buttonrespository

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/domain/button"
)

func (r *ButtonRepository) StoreButtonGroup(
	ctx context.Context,
	groupID string,
	buttons []button.Button,
) error {
	hMap, err := processButtonRows(buttons)
	if err != nil {
		return fmt.Errorf("process button rows: %w", err)
	}

	if err := r.client.HSet(ctx, groupID, hMap).Err(); err != nil {
		return fmt.Errorf("hset: %w", err)
	}

	return nil
}

func processButtonRows(
	buttons []button.Button,
) (map[string]any, error) {
	hMap := make(map[string]any, len(buttons))

	for _, btn := range buttons {
		encodedButton, err := encodeButton(btn)
		if err != nil {
			return nil, fmt.Errorf("encode button: %w", err)
		}

		hMap[btn.ID] = encodedButton
	}

	return hMap, nil
}
