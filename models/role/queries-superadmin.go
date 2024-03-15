package role

import "context"

func (m *Queries) CreateSuperAdmin() (int, error) {
	var id int
	query := `
        INSERT INTO roles (name) 
        VALUES ('superadmin') 
        ON CONFLICT (name) 
        DO UPDATE SET name='superadmin'
        RETURNING id;
    `
	if err := m.DB.GetContext(
		context.Background(),
		&id,
		query,
	); err != nil {
		return 0, err
	}

	if _, err := m.DB.ExecContext(
		context.Background(),
		`SELECT setval('roles_id_seq', (SELECT MAX(id) FROM roles));`,
	); err != nil {
		return 0, err
	}

	return id, nil
}
