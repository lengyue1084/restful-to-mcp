package biz

import (
	"context"
	"github.com/google/wire"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(
	NewOpenapiUserCase,
	NewMcpFileUserCase,
	NewMcpServerUseCase,
	NewMcpToolsUserCase,
	NewMcpGatewayUseCase,
)

type Transaction interface {
	ExecTx(ctx context.Context, fn func(ctx context.Context) error) error
}
