package v1

import (
	"gmkdb/api/v1/ninex"
	"gmkdb/api/v1/system"
)

type apiGroup struct {
	// ExampleApiGroup example.ApiGroup

	SystemApiGroup system.ApiGroup
	NineXApiGroup  ninex.ApiGroup
}

var ApiGroupApp = new(apiGroup)
