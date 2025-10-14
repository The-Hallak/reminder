package utils

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Group struct {
	GroupID     string
	Type        string
	PagesTopic  *int
	AthkarTopic *int
}

type Repo struct {
	db *pgxpool.Pool
}

func NewRepo(ctx context.Context, dbURL string) (*Repo, error) {
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, err
	}

	// Create table if not exists
	_, err = pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS groups (
			group_id     TEXT PRIMARY KEY,
			type         TEXT NOT NULL,
			pages_topic  INTEGER,
			athkar_topic INTEGER
		);
	`)
	if err != nil {
		return nil, err
	}

	return &Repo{db: pool}, nil
}

func (r *Repo) Close() { r.db.Close() }

// Get all records
func (r *Repo) GetAll(ctx context.Context) ([]Group, error) {
	rows, err := r.db.Query(ctx, `SELECT group_id, type, pages_topic, athkar_topic FROM groups;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Group
	for rows.Next() {
		var g Group
		if err := rows.Scan(&g.GroupID, &g.Type, &g.PagesTopic, &g.AthkarTopic); err != nil {
			return nil, err
		}
		result = append(result, g)
	}
	return result, rows.Err()
}

// Upsert (create or update)
func (r *Repo) Upsert(ctx context.Context, g Group) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO groups (group_id, type, pages_topic, athkar_topic)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (group_id) DO UPDATE
		SET
			type = EXCLUDED.type,
			pages_topic  = COALESCE(EXCLUDED.pages_topic,  groups.pages_topic),
			athkar_topic = COALESCE(EXCLUDED.athkar_topic, groups.athkar_topic);
	`, g.GroupID, g.Type, g.PagesTopic, g.AthkarTopic)
	return err
}

// Delete by ID
func (r *Repo) Delete(ctx context.Context, groupID string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM groups WHERE group_id = $1;`, groupID)
	return err
}
