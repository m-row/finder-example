package role

import (
	"fmt"
	"net/url"

	"github.com/m-row/finder"
	"github.com/m-row/finder-example/model"
)

const (
	ScopeAdmin  = "admin"
	ScopePublic = "public"
)

type Model struct {
	ID          int    `db:"id"          json:"id"`
	Name        string `db:"name"        json:"name"`
	Permissions []int  `db:"permissions" json:"permissions,omitempty"`
}

// Model ----------------------------------------------------------------------

func (m *Model) GetID() string {
	return fmt.Sprintf("%d", m.ID)
}

func (m *Model) ModelName() string {
	return "role"
}

func (m *Model) TableName() string {
	return "roles"
}

func (m *Model) DefaultSearch() string {
	return "name"
}

func (m *Model) SearchFields() *[]string {
	return &[]string{
		m.DefaultSearch(),
	}
}

func (m *Model) Columns(pgInfo map[string][]string) *[]string {
	return finder.GetColumns(m, pgInfo)
}

func (m *Model) Relations() *[]finder.RelationField {
	return &[]finder.RelationField{}
}

func (m *Model) Initialize(v url.Values, conn finder.Connection) bool {
	isInsert := m.ID == 0
	if isInsert {
		model.SelectSeqID(&m.ID, m.TableName(), conn)
	}
	return isInsert
}
