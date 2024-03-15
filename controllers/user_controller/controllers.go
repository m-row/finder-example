package user_controller

import "github.com/m-row/finder-example/controller"

type Controllers struct {
	Auth    *ControllerAuth
	Admin   *ControllerAdmin
	Basic   *ControllerBasic
	Profile *ControllerProfile
}

func Get(deps *controller.Dependencies) *Controllers {
	return &Controllers{
		Auth:    &ControllerAuth{deps},
		Admin:   &ControllerAdmin{deps},
		Basic:   &ControllerBasic{deps},
		Profile: &ControllerProfile{deps},
	}
}
