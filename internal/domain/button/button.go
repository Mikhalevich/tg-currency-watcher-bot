package button

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/google/uuid"
)

type ButtonType string

const (
	CurrencyPair ButtonType = "CurrencyPair"
)

type Button struct {
	ID      string
	Caption string
	Type    ButtonType
	Payload []byte
}

func encodePayload[T any](payload T) ([]byte, error) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(payload); err != nil {
		return nil, fmt.Errorf("gob encode: %w", err)
	}

	return buf.Bytes(), nil
}

func decodePaylaod[T any](b []byte) (T, error) {
	var payload T
	if err := gob.NewDecoder(bytes.NewReader(b)).Decode(&payload); err != nil {
		return payload, fmt.Errorf("gob decode: %w", err)
	}

	return payload, nil
}

func GetPayload[T any](b Button) (T, error) {
	payload, err := decodePaylaod[T](b.Payload)
	if err != nil {
		return payload, fmt.Errorf("decode paylod: %w", err)
	}

	return payload, nil
}

type CurrencyPairPayload struct {
	CurrencyID int
	IsInverted bool
}

func CurrencyPairButton(caption string, payload CurrencyPairPayload) (Button, error) {
	payloadGob, err := encodePayload(payload)
	if err != nil {
		return Button{}, fmt.Errorf("encode payload: %w", err)
	}

	return Button{
		ID:      uuid.NewString(),
		Caption: caption,
		Type:    CurrencyPair,
		Payload: payloadGob,
	}, nil
}
