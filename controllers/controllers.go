package controllers

import (
	"github.com/m-row/finder-example/controller"
	"github.com/m-row/finder-example/controllers/category_controller"
)

type Controllers struct {
	Category *category_controller.Controllers
}

func Setup(d *controller.Dependencies) *Controllers {
	return &Controllers{
		Category: category_controller.Get(d),
	}
}
