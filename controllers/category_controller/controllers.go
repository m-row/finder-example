package category_controller

import (
	"github.com/m-row/finder-example/controller"
)

type Controllers struct {
	*controller.Dependencies
}

func Get(d *controller.Dependencies) *Controllers {
	return &Controllers{d}
}
