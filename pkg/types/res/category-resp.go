package response

import "github.com/wafi04/backend/pkg/types"

type ListCategoriesResponse struct {
	Categories []*types.Category `json:"categories"`
	Total      int32             `json:"total"`
}

type CategoryHierarchyResponse struct {
	RootCategory     types.Category `json:"root_category"`
	TotalDescendants int32          `json:"total_descendants"`
	MaxDepth         int32          `json:"max_depth"`
}

type DeleteCategoryResponse struct {
	Success      bool  `json:"success"`
	DeletedCount int64 `json:"deleted_count"`
}
