package button

import (
	"context"
	"fmt"
)

type ButtonRepository interface {
	GetButton(ctx context.Context, id string) (*Button, error)
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

func (bp *ButtonProvider) GetButton(ctx context.Context, btnID string) (*Button, error) {
	btn, err := bp.repository.GetButton(ctx, btnID)
	if err != nil {
		return nil, fmt.Errorf("load button: %w", err)
	}

	return btn, nil
}

func (bp *ButtonProvider) SetButton(ctx context.Context, btn Button) error {
	if err := bp.repository.StoreButton(ctx, btn); err != nil {
		return fmt.Errorf("store button: %w", err)
	}

	return nil
}

func (bp *ButtonProvider) SetButtonGroup(ctx context.Context, groupID string, buttons []Button) error {
	if err := bp.repository.StoreButtonGroup(ctx, groupID, buttons); err != nil {
		return fmt.Errorf("store button group: %w", err)
	}

	return nil
}
