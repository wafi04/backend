package category

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/wafi04/backend/pkg/types"
	request "github.com/wafi04/backend/pkg/types/req"
	response "github.com/wafi04/backend/pkg/types/res"
)

type CategoryRepository interface {
	Create(ctx context.Context, category *types.Category, depth int32) (*types.Category, error)
	GetParentDepth(ctx context.Context, parentID string) (int32, error)
	GetCategoryTree(ctx context.Context) (map[string]*types.Category, []*types.Category, error)
	DeleteCategory(ctx context.Context, req *request.DeleteCategoryRequest) (*response.DeleteCategoryResponse, error)
	UpdateCategory(ctx context.Context, req *request.UpdateCategoryRequest) (*types.Category, error)
}

type categoryRepository struct {
	db *sqlx.DB
}

func NewCategoryRepository(db *sqlx.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(ctx context.Context, category *types.Category, depth int32) (*types.Category, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	query := `
        INSERT INTO categories (
            id,
            name,
            description,
            image,
            parent_id,
            depth,
            created_at
        ) VALUES (
            $1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP
        )
        RETURNING id, name, description, image, parent_id, depth, created_at`

	var createdAt sql.NullTime
	var parentID, image sql.NullString

	err = tx.QueryRowContext(ctx, query,
		category.ID,
		category.Name,
		category.Description,
		category.Image,
		category.ParentID,
		depth,
	).Scan(
		&category.ID,
		&category.Name,
		&category.Description,
		&image,
		&parentID,
		&depth,
		&createdAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to insert category: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	if parentID.Valid {
		category.ParentID = &parentID.String
	}
	if image.Valid {
		category.Image = &image.String
	}
	if createdAt.Valid {
		category.CreatedAt = createdAt.Time
	}

	return category, nil
}

func (r *categoryRepository) GetParentDepth(ctx context.Context, parentID string) (int32, error) {
	var depth int32
	err := r.db.QueryRowContext(ctx, "SELECT depth FROM categories WHERE id = $1", parentID).Scan(&depth)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("parent category not found")
		}
		return 0, fmt.Errorf("failed to get parent category: %v", err)
	}
	return depth, nil
}

func (r *categoryRepository) GetCategoryTree(ctx context.Context) (map[string]*types.Category, []*types.Category, error) {
	query := `
        WITH RECURSIVE category_tree AS (
            SELECT 
                c.id, c.name, c.description, c.image, 
                c.parent_id, c.depth, c.created_at,
                ARRAY[]::VARCHAR[] AS path,
                0 as level
            FROM categories c
            WHERE c.parent_id IS NULL
            UNION ALL
            SELECT 
                c.id, c.name, c.description, c.image,
                c.parent_id, c.depth, c.created_at,
                path || c.parent_id,
                ct.level + 1
            FROM categories c
            INNER JOIN category_tree ct ON ct.id = c.parent_id
        )
        SELECT 
            id, name, description, image,
            parent_id, depth, created_at,
            path
        FROM category_tree
        ORDER BY path, level`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to query categories: %v", err)
	}
	defer rows.Close()

	categoryMap := make(map[string]*types.Category)
	var rootCategories []*types.Category
	for rows.Next() {
		var cat types.Category
		var createdAt sql.NullTime
		var parentID, image sql.NullString
		cat.Path = pq.StringArray{}

		err := rows.Scan(
			&cat.ID,
			&cat.Name,
			&cat.Description,
			&image,
			&parentID,
			&cat.Depth,
			&createdAt,
			&cat.Path,
		)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to scan category: %v", err)
		}

		if parentID.Valid {
			cat.ParentID = &parentID.String
		}

		if image.Valid {
			cat.Image = &image.String
		}
		if createdAt.Valid {
			cat.CreatedAt = createdAt.Time
		}

		categoryMap[cat.ID] = &cat

		if !parentID.Valid {
			rootCategories = append(rootCategories, &cat)
		} else {
			parent := categoryMap[parentID.String]
			if parent != nil {
				parent.Children = append(parent.Children, &cat)
			}
		}
	}

	if err = rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("error iterating categories: %v", err)
	}

	return categoryMap, rootCategories, nil
}

func (r *categoryRepository) UpdateCategory(ctx context.Context, req *request.UpdateCategoryRequest) (*types.Category, error) {
	query := `UPDATE categories SET `

	updates := []string{}
	args := []interface{}{}
	argCount := 1

	if req.Name != nil {
		updates = append(updates, fmt.Sprintf("name = $%d", argCount))
		args = append(args, *req.Name)
		argCount++
	}

	if req.Description != nil {
		updates = append(updates, fmt.Sprintf("description = $%d", argCount))
		args = append(args, *req.Description)
		argCount++
	}

	if req.Image != nil {
		updates = append(updates, fmt.Sprintf("image = $%d", argCount))
		args = append(args, *req.Image)
		argCount++
	}

	if req.ParentID != nil {
		updates = append(updates, fmt.Sprintf("parent_id = $%d", argCount))
		args = append(args, *req.ParentID)
		argCount++
	}

	if len(updates) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	query += strings.Join(updates, ", ")
	query += fmt.Sprintf(" WHERE id = $%d RETURNING id, name, description, image, depth, parent_id, created_at", argCount)
	args = append(args, req.ID)

	var category types.Category
	var parentID, image sql.NullString
	var createdAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&category.ID,
		&category.Name,
		&category.Description,
		&image,
		&category.Depth,
		&parentID,
		&createdAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("category not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update category: %v", err)
	}

	if image.Valid {
		category.Image = &image.String
	}
	if parentID.Valid {
		category.ParentID = &parentID.String
	}
	if createdAt.Valid {
		category.CreatedAt = createdAt.Time
	}

	return &category, nil
}

func (s *categoryRepository) DeleteCategory(ctx context.Context, req *request.DeleteCategoryRequest) (*response.DeleteCategoryResponse, error) {
	// Start a transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	// First check if category exists
	var exists bool
	err = tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM categories WHERE id = $1)", req.ID).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("failed to check category existence: %v", err)
	}
	if !exists {
		return nil, fmt.Errorf("category not found")
	}

	var deletedCount int64
	if req.DeleteChildren {
		withDescendantsQuery := `
            WITH RECURSIVE category_tree AS (
                -- Base case: the category we want to delete
                SELECT id FROM categories WHERE id = $1
                UNION ALL
                -- Recursive case: get all children
                SELECT c.id 
                FROM categories c
                INNER JOIN category_tree ct ON c.parent_id = ct.id
            )
            DELETE FROM categories 
            WHERE id IN (SELECT id FROM category_tree)
            RETURNING id`

		result, err := tx.ExecContext(ctx, withDescendantsQuery, req.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to delete category and its children: %v", err)
		}
		deletedCount, err = result.RowsAffected()
		if err != nil {
			return nil, fmt.Errorf("failed to get affected rows: %v", err)
		}
	} else {
		var hasChildren bool
		err = tx.QueryRowContext(ctx,
			"SELECT EXISTS(SELECT 1 FROM categories WHERE parent_id = $1)",
			req.ID,
		).Scan(&hasChildren)
		if err != nil {
			return nil, fmt.Errorf("failed to check for children: %v", err)
		}
		if hasChildren {
			return nil, fmt.Errorf("cannot delete category with children, set DeleteChildren to true to delete all")
		}

		// Delete single category
		result, err := tx.ExecContext(ctx,
			"DELETE FROM categories WHERE id = $1",
			req.ID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to delete category: %v", err)
		}
		deletedCount, err = result.RowsAffected()
		if err != nil {
			return nil, fmt.Errorf("failed to get affected rows: %v", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return &response.DeleteCategoryResponse{
		Success:      true,
		DeletedCount: deletedCount,
	}, nil
}
