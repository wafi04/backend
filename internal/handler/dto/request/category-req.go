package request

type CreateCategoryRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Image       *string `json:"image,omitempty"`
	ParentID    *string `json:"parent_id,omitempty"`
}

type GetCategoryRequest struct {
	ID string `json:"id"`
}

type UpdateCategoryRequest struct {
	ID          string  `json:"id"`
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Image       *string `json:"image,omitempty"`
	Depth       *int32  `json:"depth,omitempty"`
	ParentID    *string `json:"parent_id,omitempty"`
}

type DeleteCategoryRequest struct {
	ID             string `json:"id"`
	DeleteChildren bool   `json:"delete_children"`
}

type ListCategoriesRequest struct {
	Page            int32   `json:"page"`
	Limit           int32   `json:"limit"`
	ParentID        *string `json:"parent_id,omitempty"`
	IncludeChildren bool    `json:"include_children"`
}
