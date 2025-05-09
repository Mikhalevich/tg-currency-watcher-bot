package button

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type ButtonRepository interface {
	StoreButton(ctx context.Context, btn Button) error
	StoreButtonGroup(ctx context.Context, groupID string, btns []Button) error
}

type ButtonProvider struct {
	repository ButtonRepository
}

func NewButtonProvider(repository ButtonRepository) *ButtonProvider {
	return &ButtonProvider{
		repository: repository,
	}
}

func (bp *ButtonProvider) SetButton(ctx context.Context, btn Button) error {
	if err := bp.repository.StoreButton(ctx, btn); err != nil {
		return fmt.Errorf("store button: %w", err)
	}

	return nil
}

func (bp *ButtonProvider) ButtonGroup() (*ButtonGroup, string) {
	groupID := uuid.NewString()

	return NewButtonGroup(groupID, bp.repository), groupID
}
