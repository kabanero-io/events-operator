package model

import (
	"github.com/kabanero-io/events-operator/pkg/semverimage"
)

type StringToVersion struct {
	Str     string
	Version *semverimage.Version
}