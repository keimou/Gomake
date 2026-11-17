package routers

import (
	"gmkdb/routers/ninex"
	"gmkdb/routers/system"
)

type routerGroup struct {
	// Example example.RouterGroup

	System system.RouterGroup
	Ninex  ninex.RouterGroup
}

var RouterGroupApp = new(routerGroup)
