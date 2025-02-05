package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/wafi04/backend/pkg/types"
	request "github.com/wafi04/backend/pkg/types/req"
	response "github.com/wafi04/backend/pkg/types/res"
	category "github.com/wafi04/backend/services/category/repository"
)

type CategoryService struct {
	categoryRepo category.CategoryRepository
}

func NewCategoryService(categoryRepo category.CategoryRepository) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
	}
}

func (s *CategoryService) CreateCategory(ctx context.Context, req *request.CreateCategoryRequest) (*types.Category, error) {
	categoryID := uuid.New().String()

	var depth int32 = 0
	if req.ParentID != nil {
		parentDepth, err := s.categoryRepo.GetParentDepth(ctx, *req.ParentID)
		if err != nil {
			return nil, err
		}
		depth = parentDepth + 1
	}

	category := &types.Category{
		ID:          categoryID,
		Name:        req.Name,
		Description: req.Description,
		Image:       req.Image,
		ParentID:    req.ParentID,
	}

	return s.categoryRepo.Create(ctx, category, depth)
}

func (s *CategoryService) GetCategories(ctx context.Context, req *request.ListCategoriesRequest) (*response.ListCategoriesResponse, error) {
	categoryMap, rootCategories, err := s.categoryRepo.GetCategoryTree(ctx)
	if err != nil {
		return nil, err
	}

	return &response.ListCategoriesResponse{
		Categories: rootCategories,
		Total:      int32(len(categoryMap)),
	}, nil
}

func (s *CategoryService) UppdateCategory(ctx context.Context, req *request.UpdateCategoryRequest) (*types.Category, error) {
	return s.categoryRepo.UpdateCategory(ctx, req)
}

func (s *CategoryService) DeleteCategory(ctx context.Context, req *request.DeleteCategoryRequest) (*response.DeleteCategoryResponse, error) {
	return s.categoryRepo.DeleteCategory(ctx, req)
}
