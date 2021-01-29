package model

import (
	"context"
	"errors"
	em "github.com/Etpmls/Etpmls-Micro/v2"
)

type Model struct {

}

func (this *Model) ReadTokenFromCtx(ctx context.Context) (string, error) {
	if ctx == nil {
		em.LogError.Output(em.MessageWithLineNum("Failed to obtain request!"))
		return "", errors.New("Failed to obtain request!")
	}

	token := ctx.Value("token");
	if token == nil {
		em.LogError.Output(em.MessageWithLineNum("Failed to obtain token!"))
		return "", errors.New("Failed to obtain token!")
	}

	return token.(string), nil
}