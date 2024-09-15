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

func (m *Controllers) SetRoutes(
	d *controller.RouterDependencies,
) {
	r := d.E.Group("/categories")

	r.GET("", m.Index).Name = "categories:index"
	r.POST("", m.Store).Name = "categories:store"
	r.GET("/:id", m.Show).Name = "categories:show"
	r.PUT("/:id", m.Update).Name = "categories:update"
	r.DELETE("/:id", m.Destroy).Name = "categories:destroy"
}
