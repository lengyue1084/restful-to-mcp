package database

import (
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	//NewMysqlClient,
	NewPgClient,
	NewData,
	NewTransaction,
)
