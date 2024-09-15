package category

import (
	"context"
	"net/url"
	"time"

	"github.com/m-row/finder-example/model"
	"github.com/m-row/finder-example/types"

	"github.com/google/uuid"
	"github.com/m-row/finder"
)

type Model struct {
	ID            uuid.UUID   `db:"id"              json:"id"`
	Name          types.JSONB `db:"name"            json:"name"`
	Img           *string     `db:"img"             json:"img"`
	Thumb         *string     `db:"thumb"           json:"thumb"`
	Depth         int         `db:"depth"           json:"depth"`
	Sort          int         `db:"sort"            json:"sort"`
	IsDisabled    bool        `db:"is_disabled"     json:"is_disabled"`
	IsFeatured    bool        `db:"is_featured"     json:"is_featured"`
	ParentID      *uuid.UUID  `db:"parent_id"       json:"-"`
	SuperParentID *uuid.UUID  `db:"super_parent_id" json:"-"`
	CreatedAt     time.Time   `db:"created_at"      json:"created_at"`
	UpdatedAt     time.Time   `db:"updated_at"      json:"updated_at"`

	Parent      MinimalModel `db:"parent"       json:"parent"`
	SuperParent MinimalModel `db:"super_parent" json:"super_parent"`
}

type MinimalModel struct {
	ID   *uuid.UUID   `db:"id"   json:"id"`
	Name *types.JSONB `db:"name" json:"name"`
}

// Model ----------------------------------------------------------------------

func (m *Model) GetID() string {
	return m.ID.String()
}

func (m *Model) ModelName() string {
	return "category"
}

func (m *Model) TableName() string {
	return "categories"
}

func (m *Model) DefaultSearch() string {
	return "name->>'ar'"
}

func (m *Model) SearchFields() *[]string {
	return &[]string{
		m.DefaultSearch(),
		"name->>'en'",
		"code",
	}
}

func (m *Model) Columns(pgInfo map[string][]string) *[]string {
	return finder.GetColumns(m, pgInfo)
}

func (m *Model) Relations() *[]finder.RelationField {
	return &[]finder.RelationField{}
}

func (m *Model) Initialize(v url.Values, conn finder.Connection) bool {
	isInsert := m.CreatedAt.Equal(time.Time{})
	if isInsert && m.ID == uuid.Nil {
		model.InputOrNewUUID(&m.ID, v)
	}
	return isInsert
}

// HasImage -------------------------------------------------------------------

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
} // HasImage interface end

// Utilities ------------------------------------------------------------------

// func (m *Model) MergeAndValidate(v *validator.Validator) bool {
// 	m.Initialize(v.Data.Values, v.DB)
// 	v.AssignString("code", &m.Code)
// 	v.UnmarshalInto("name", &m.Name)
// 	v.AssignBool("is_disabled", &m.IsDisabled)
// 	v.AssignBool("is_featured", &m.IsFeatured)
// 	v.AssignInt("sort", &m.Sort)
//
// 	if err := v.AssignImage("img", m, false); err != nil {
// 		v.Check(false, "img", err.Error())
// 	}
//
// 	m.BusinessID = v.AssignUUID(
// 		"business_id",
// 		"businesses",
// 		m.BusinessID,
// 		false,
// 		ScopeAdmin,
// 	)
// 	// parent,super_parent and depth assignment
// 	v.UnmarshalInto("parent", &m.Parent)
// 	if m.Parent.ID != nil {
// 		if *m.Parent.ID == m.ID {
// 			v.Check(false, "parent.id", v.T.ValidateCategoryParent())
// 		} else {
// 			if err := m.AssignSuperParent(v.DB); err != nil {
// 				v.Check(false, "parent.id", err.Error())
// 			}
// 		}
// 	}
//
// 	v.ValidateModelSchema(m, v.Schema)
// 	return v.Valid()
// }

// AssignSuperParent gets parent super_parent and depth assigned to body.
func (m *Model) AssignSuperParent(db finder.Connection) error {
	if m.Parent.ID != nil {
		var parent Model
		if err := db.GetContext(
			context.Background(),
			&parent,
			`
                SELECT 
                    id,
                    name,
                    parent_id,
                    super_parent_id,
                    depth
                FROM 
                    categories 
                WHERE 
                    id=$1
            `,
			m.Parent.ID,
		); err != nil {
			return err
		}
		if parent.Depth == 0 {
			m.SuperParent.ID = &parent.ID
		} else {
			m.SuperParent.ID = parent.SuperParentID
		}
		m.Depth = parent.Depth + 1
	}
	return nil
}
