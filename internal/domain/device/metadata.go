package device

import "errors"

type Metadata map[string]any

func NewMetadata(raw any) (Metadata, error) {
	if raw == nil {
		return Metadata{}, errors.New("device metadata is required")
	}

	p, ok := raw.(map[string]any)
	if !ok {
		return Metadata{}, errors.New("device metadata must be a valid JSON object")
	}

	if len(p) == 0 {
		return Metadata{}, errors.New("device metadata cannot be empty")
	}

	return Metadata(p), nil
}
