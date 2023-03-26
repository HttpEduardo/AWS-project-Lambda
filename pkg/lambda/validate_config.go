package lambda

import (
	"errors"

	"github.com/suzuki-shunsuke/aws-iam-cred-sender/pkg/controller"
)

var errSecretIDIsRequired = errors.New("secret ID is required")

func (handler *Handler) validateConfig(cfg controller.Config) error {
	if cfg.SecretID == "" {
		return errSecretIDIsRequired
	}
	return nil
}
