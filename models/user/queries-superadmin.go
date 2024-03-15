package user

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func (m *Queries) CreateSuperAdmin(roleID int) error {
	hash, err := bcrypt.GenerateFromPassword([]byte("password"), 12)
	if err != nil {
		return fmt.Errorf("error hashing at seeding users: %w", err)
	}
	c := context.Background()
	query := `
        INSERT INTO users (
            id,
            name,
            email,
            ref,
            password_hash
        ) 
        VALUES (
            $1,
            'superadmin',
            'superadmin@sadeem-tech.com',
            'ACFE1828',
            $2
        ) 
        ON CONFLICT (id) 
        DO NOTHING;
    `
	if _, err := m.DB.ExecContext(c, query, SuperAdminID, hash); err != nil {
		return err
	}

	query2 := `
        INSERT INTO user_roles (
            user_id,
            role_id
        ) 
        VALUES (
            $1,
            $2
        ) 
        ON CONFLICT (user_id, role_id) 
        DO NOTHING;
    `
	if _, err := m.DB.ExecContext(c, query2, SuperAdminID, roleID); err != nil {
		return err
	}
	return nil
}
