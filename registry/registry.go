package registry

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

type Registry struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	Type         string         `json:"type"`
	Endpoint     string         `json:"endpoint"`
	Protocol     string         `json:"protocol"`
	AuthRequired bool           `json:"auth_required"`
	PolicyIDs    []string       `json:"policy_ids"`
	Metadata     map[string]any `json:"metadata"`
	Status       string         `json:"status"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

var ErrNotFound = errors.New("registry not found")

func insert(ctx context.Context, r *Registry) error {
	pool, err := db()
	if err != nil {
		return err
	}
	metaBytes, err := json.Marshal(r.Metadata)
	if err != nil {
		return err
	}
	row := pool.QueryRow(ctx, `
		INSERT INTO registries (name, type, endpoint, protocol, auth_required, policy_ids, metadata, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`, r.Name, r.Type, r.Endpoint, r.Protocol, r.AuthRequired, r.PolicyIDs, metaBytes, r.Status)
	return row.Scan(&r.ID, &r.CreatedAt, &r.UpdatedAt)
}

func update(ctx context.Context, id string, r *Registry) error {
	pool, err := db()
	if err != nil {
		return err
	}
	metaBytes, err := json.Marshal(r.Metadata)
	if err != nil {
		return err
	}
	row := pool.QueryRow(ctx, `
		UPDATE registries
		SET name=$1, type=$2, endpoint=$3, protocol=$4, auth_required=$5,
		    policy_ids=$6, metadata=$7, status=$8, updated_at=now()
		WHERE id=$9
		RETURNING id, created_at, updated_at
	`, r.Name, r.Type, r.Endpoint, r.Protocol, r.AuthRequired, r.PolicyIDs, metaBytes, r.Status, id)
	err = row.Scan(&r.ID, &r.CreatedAt, &r.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	return err
}

func deleteByID(ctx context.Context, id string) error {
	pool, err := db()
	if err != nil {
		return err
	}
	tag, err := pool.Exec(ctx, `DELETE FROM registries WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func getByID(ctx context.Context, id string) (*Registry, error) {
	pool, err := db()
	if err != nil {
		return nil, err
	}
	row := pool.QueryRow(ctx, `
		SELECT id, name, type, endpoint, protocol, auth_required, policy_ids,
		       metadata::text, status, created_at, updated_at
		FROM registries WHERE id=$1
	`, id)
	r, err := scanRegistry(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return r, err
}

func getByName(ctx context.Context, name string) (*Registry, error) {
	pool, err := db()
	if err != nil {
		return nil, err
	}
	row := pool.QueryRow(ctx, `
		SELECT id, name, type, endpoint, protocol, auth_required, policy_ids,
		       metadata::text, status, created_at, updated_at
		FROM registries
		WHERE name=$1
		ORDER BY created_at DESC
		LIMIT 1
	`, name)
	r, err := scanRegistry(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return r, err
}

func getAll(ctx context.Context) ([]*Registry, error) {
	pool, err := db()
	if err != nil {
		return nil, err
	}
	rows, err := pool.Query(ctx, `
		SELECT id, name, type, endpoint, protocol, auth_required, policy_ids,
		       metadata::text, status, created_at, updated_at
		FROM registries ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*Registry
	for rows.Next() {
		r, err := scanRegistry(rows)
		if err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	return results, rows.Err()
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanRegistry(s rowScanner) (*Registry, error) {
	var r Registry
	var metaStr string
	err := s.Scan(
		&r.ID, &r.Name, &r.Type, &r.Endpoint, &r.Protocol,
		&r.AuthRequired, &r.PolicyIDs, &metaStr,
		&r.Status, &r.CreatedAt, &r.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(metaStr), &r.Metadata); err != nil {
		return nil, err
	}
	return &r, nil
}
