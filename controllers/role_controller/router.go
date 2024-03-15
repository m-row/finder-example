//nolint:lll
package role_controller

import (
	"github.com/m-row/finder-example/controller"
	"github.com/m-row/finder-example/models/role"
)

func (m *Controllers) SetBasicRoutes(
	d *controller.RouterDependencies,
) {
	f := d.E.Group("/roles")

	f.GET("", m.Basic.Index).Name = "roles:index:public"
	f.GET("/:id", m.Basic.Show).Name = "roles:show:public"

	r := d.Requires(role.ScopeAdmin)

	f.POST("", m.Basic.Store, r).Name = "roles:store:admin"
	f.PUT("/:id", m.Basic.Update, r).Name = "roles:update:admin"
	f.DELETE("/:id", m.Basic.Destroy, r).Name = "roles:destroy:admin"

	f.POST("/:id/grant-all", m.Basic.GrantAllPermissions, r).Name = "roles:grant-all:admin"
	f.POST("/:id/grant-by-scope", m.Basic.GrantByScope, r).Name = "roles:grant-scope:admin"
	f.POST("/:id/revoke-all", m.Basic.RevokeAllPermissions, r).Name = "roles:revoke-all:admin"
}
