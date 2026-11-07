package nsq

import "errors"

var (
	// ErrInvalidConfig indicates configuration validation failed.
	ErrInvalidConfig = errors.New("nsq: invalid config")
)
