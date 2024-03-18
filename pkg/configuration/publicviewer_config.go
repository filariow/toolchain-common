package configuration

import (
	toolchainv1alpha1 "github.com/codeready-toolchain/api/api/v1alpha1"
)

type PublicViewerConfig struct {
	Config *toolchainv1alpha1.PublicViewerConfiguration
}

func (c PublicViewerConfig) Enabled() bool {
	return c.Config != nil && c.Config.Enabled
}
