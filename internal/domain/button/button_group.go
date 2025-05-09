package button

import (
	"context"
	"fmt"
	"strconv"
)

type ButtonGroupRepository interface {
	StoreButtonGroup(ctx context.Context, groupID string, btns []Button) error
}

type ButtonGroup struct {
	groupID string
	buttons []Button

	repository ButtonGroupRepository
}

func NewButtonGroup(groupID string, repository ButtonGroupRepository) *ButtonGroup {
	return &ButtonGroup{
		groupID:    groupID,
		repository: repository,
	}
}

func (bg *ButtonGroup) AddButton(btn Button) {
	bg.buttons = append(bg.buttons, btn)
}

func (bg *ButtonGroup) Store(ctx context.Context) error {
	for i := range bg.buttons {
		bg.buttons[i].ID = strconv.Itoa(i + 1)
	}

	if err := bg.repository.StoreButtonGroup(ctx, bg.groupID, bg.buttons); err != nil {
		return fmt.Errorf("store button group: %w", err)
	}

	return nil
}
