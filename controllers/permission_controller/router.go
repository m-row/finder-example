package permission_controller

import "github.com/m-row/finder-example/controller"

func (m *Controllers) SetBasicRoutes(
	d *controller.RouterDependencies,
) {
	f := d.E.Group("/permissions")
	f.GET("", m.Basic.Index).Name = "permissions:index:public"
	f.GET("/:id", m.Basic.Show).Name = "permissions:show:public"
}
