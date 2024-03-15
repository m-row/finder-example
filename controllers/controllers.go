package controllers

import (
	"github.com/m-row/finder-example/controller"
	"github.com/m-row/finder-example/controllers/permission_controller"
	"github.com/m-row/finder-example/controllers/role_controller"
	"github.com/m-row/finder-example/controllers/user_controller"
)

type Controllers struct {
	// API ---------------------------------------------------------------------

	User       *user_controller.Controllers
	Role       *role_controller.Controllers
	Permission *permission_controller.Controllers
}

func Setup(d *controller.Dependencies) *Controllers {
	return &Controllers{
		User:       user_controller.Get(d),
		Role:       role_controller.Get(d),
		Permission: permission_controller.Get(d),
	}
}
