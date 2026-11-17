package services

import (
	"gmkdb/services/ninex"
	"gmkdb/services/system"
)

type serviceGroup struct {
	// ExampleServiceGroup example.ServiceGroup

	SystemSvcGroup system.ServiceGroup
	NineXSvcGroup  ninex.ServiceGroup
}

var ServiceGroupApp = new(serviceGroup)
