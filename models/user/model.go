package user

import (
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/m-row/finder"
	"github.com/m-row/finder-example/model"
)

const (
	ScopeOwn   = "own"
	ScopeAdmin = "admin"

	SuperAdminID = "322f3e97-4e7e-4c2e-a765-1c0ce517f2f8"
)

type Model struct {
	ID           uuid.UUID `db:"id"            json:"id"`
	Name         *string   `db:"name"          json:"name"`
	Phone        *string   `db:"phone"         json:"phone"`
	Email        *string   `db:"email"         json:"email"`
	Password     password  `db:"-"             json:"-"`
	PasswordHash *[]byte   `db:"password_hash" json:"-"`
	Img          *string   `db:"img"           json:"img"`
	Thumb        *string   `db:"thumb"         json:"thumb"`
	IsDisabled   bool      `db:"is_disabled"   json:"is_disabled"`
	CreatedAt    time.Time `db:"created_at"    json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"    json:"updated_at"`

	Roles       *[]int `db:"roles"       json:"roles,omitempty"`
	Permissions *[]int `db:"permissions" json:"permissions,omitempty"`
}

type MinimalModel struct {
	ID    *uuid.UUID `db:"id"    json:"id"`
	Name  *string    `db:"name"  json:"name"`
	Phone *string    `db:"phone" json:"phone"`
	Email *string    `db:"email" json:"email"`
}

type UserRole struct {
	UserID uuid.UUID `db:"user_id"`
	RoleID int       `db:"role_id"`
}

// Model ----------------------------------------------------------------------

func (m *Model) GetID() string {
	return m.ID.String()
}

func (m *Model) ModelName() string {
	return "user"
}

func (m *Model) TableName() string {
	return "users"
}

func (m *Model) DefaultSearch() string {
	return "name"
}

func (m *Model) SearchFields() *[]string {
	return &[]string{
		m.DefaultSearch(),
		"email",
		"phone",
	}
}

func (m *Model) Columns(pgInfo map[string][]string) *[]string {
	return finder.GetColumns(m, pgInfo)
}

func (m *Model) Relations() *[]finder.RelationField {
	return &[]finder.RelationField{
		{
			Table: "roles",
			Join: &finder.Join{
				From: "users.id",
				To:   "roles.id",
			},
			Through: &finder.Through{
				Table: "user_roles",
				Join: &finder.Join{
					From: "user_roles.user_id",
					To:   "user_roles.role_id",
				},
			},
		},
	}
}

// Initialize generates a uuid and crc32 checksum hash for the user based on
// that uuid, must be called for user model
func (m *Model) Initialize(v url.Values, conn finder.Connection) bool {
	isInsert := m.CreatedAt.Equal(time.Time{})
	if isInsert || m.ID == uuid.Nil {
		model.InputOrNewUUID(&m.ID, v)
	}
	return isInsert
}

// Has Image ------------------------------------------------------------------

func (m *Model) GetImg() *string {
	return m.Img
}

func (m *Model) SetImg(name *string) {
	m.Img = name
}

func (m *Model) GetThumb() *string {
	return m.Thumb
}

func (m *Model) SetThumb(name *string) {
	m.Thumb = name
}
