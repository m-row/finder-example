package permission_controller

import "github.com/m-row/finder-example/controller"

type Controllers struct {
	Basic *ControllerBasic
}

func Get(deps *controller.Dependencies) *Controllers {
	return &Controllers{
		Basic: &ControllerBasic{deps},
	}
}
