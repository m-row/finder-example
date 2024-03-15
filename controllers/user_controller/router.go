//nolint:lll
package user_controller

import (
	"github.com/m-row/finder-example/controller"
	"github.com/m-row/finder-example/models/user"
)

func (m *Controllers) SetAuthRoutes(
	d *controller.RouterDependencies,
) {
	d.E.POST("/login", m.Auth.Login).Name = "auth:login:public"
	// requires jwt
	d.E.GET("/logout", m.Auth.Logout).Name = "auth:logout:public"
}

func (m *Controllers) SetProfileRoutes(
	d *controller.RouterDependencies,
) {
	f := d.E.Group("/me")

	f.GET("", m.Profile.Me).Name = "user:me:public"

	r := d.Requires(
		user.ScopeOwn,
	)

	f.PUT("", m.Profile.Update, r).Name = "user:me-update:own"
	f.DELETE("", m.Profile.Clear, r).Name = "user:me-clear:own"
}

func (m *Controllers) SetAdminRoutes(
	d *controller.RouterDependencies,
) {
	f := d.E.Group("/users")
	r := d.Requires(
		user.ScopeAdmin,
	)

	f.GET("", m.Basic.Index, r).Name = "users:index:admin"
	f.POST("", m.Basic.Store, r).Name = "users:store:admin"
	f.GET("/:id", m.Basic.Show, r).Name = "users:show:admin"
	f.PUT("/:id", m.Basic.Update, r).Name = "users:update:admin"
	f.DELETE("/:id", m.Basic.Clear, r).Name = "users:clear:admin"

	f.GET("/:id/become", m.Admin.Become, r).Name = "users:become:admin"
	f.POST("/grant-role", m.Admin.GrantRole, r).Name = "users:grant-role:admin"
	f.POST("/revoke-role", m.Admin.RevokeRole, r).Name = "users:revoke-role:admin"
}
