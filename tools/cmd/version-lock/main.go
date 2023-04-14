// This package is to allow "stub" imports of CLI tools so their version can be tracked via go.mod
package main

import (
	// Version locks.
	_ "github.com/bufbuild/buf/cmd/buf"
	_ "github.com/dmarkham/enumer"
	_ "github.com/golang-migrate/migrate/v4/cmd/migrate"
	_ "github.com/kyleconroy/sqlc/cmd/sqlc"
	_ "github.com/vburenin/ifacemaker"
	_ "github.com/vektra/mockery/v2"
)

func main() {}
