package types

import (
	"time"

	"github.com/rs/zerolog"
)

type ServiceInfo struct {
	Name        string
	UpTime      time.Time
	Environment string
	Version     string
}

func (s ServiceInfo) MarshalZerologObject(e *zerolog.Event) {
	e.Str("name", s.Name).
		Str("environment", s.Environment).
		Time("started", s.UpTime).
		Str("version", s.Version)
}
