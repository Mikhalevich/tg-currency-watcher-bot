package buttonrespository

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/domain/button"
)

var (
	_ button.ButtonRepository = (*ButtonRepository)(nil)
)

type ButtonRepository struct {
	client redis.UniversalClient
	ttl    time.Duration
}

func New(client redis.UniversalClient, ttl time.Duration) *ButtonRepository {
	return &ButtonRepository{
		client: client,
		ttl:    ttl,
	}
}

func (r *ButtonRepository) IsNotFoundError(err error) bool {
	return errors.Is(err, redis.Nil)
}

func encodeButton(b button.Button) ([]byte, error) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(b); err != nil {
		return nil, fmt.Errorf("gob encode: %w", err)
	}

	return buf.Bytes(), nil
}

func decodeButton(b []byte) (*button.Button, error) {
	var btn button.Button
	if err := gob.NewDecoder(bytes.NewReader(b)).Decode(&btn); err != nil {
		return nil, fmt.Errorf("gob decode: %w", err)
	}

	return &btn, nil
}

func parseButtonID(btnID string) (string, string) {
	var (
		keySplitterIdx = strings.Index(btnID, "_")
	)

	if keySplitterIdx == -1 {
		return btnID, ""
	}

	return btnID[:keySplitterIdx], btnID[keySplitterIdx+1:]
}
