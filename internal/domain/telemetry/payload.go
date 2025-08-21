package telemetry

import "errors"

type Payload map[string]any

func NewPayload(raw any) (Payload, error) {
	if raw == nil {
		return Payload{}, errors.New("telemetry payload is required")
	}

	p, ok := raw.(map[string]any)
	if !ok {
		return Payload{}, errors.New("telemetry payload must be a valid JSON object")
	}

	if len(p) == 0 {
		return Payload{}, errors.New("telemetry payload cannot be empty")
	}

	return Payload(p), nil
}
