package service

import (
	"context"
	"io"

	"github.com/lucasacoutinho/go-grpc/internal/database"
	"github.com/lucasacoutinho/go-grpc/internal/pb"
)

type CategoryService struct {
	pb.UnimplementedCategoryServiceServer
	DB *database.Category
}

func NewCategoryService(db *database.Category) *CategoryService {
	return &CategoryService{DB: db}
}

func (cs *CategoryService) ListCategories(ctx context.Context, in *pb.Blank) (*pb.CategoryListResponse, error) {
	categories, err := cs.DB.FindAll()
	if err != nil {
		return nil, err
	}

	var categoriesResponse []*pb.Category

	for _, category := range categories {
		categorieResponse := &pb.Category{
			Id:          category.ID,
			Name:        category.Name,
			Description: category.Description,
		}

		categoriesResponse = append(categoriesResponse, categorieResponse)
	}

	return &pb.CategoryListResponse{Categories: categoriesResponse}, nil
}

func (cs *CategoryService) CreateCategory(ctx context.Context, in *pb.CategoryCreateRequest) (*pb.CategoryResponse, error) {
	category, err := cs.DB.Create(in.Name, in.Description)
	if err != nil {
		return nil, err
	}

	result := &pb.Category{
		Id:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}

	return &pb.CategoryResponse{Category: result}, nil
}

func (cs *CategoryService) CreateCategoryStream(
	stream pb.CategoryService_CreateCategoryStreamServer,
) error {
	categories := &pb.CategoryListResponse{}

	for {
		category, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(categories)
		}
		if err != nil {
			return err
		}

		categoryResult, err := cs.DB.Create(category.Name, category.Description)
		if err != nil {
			return err
		}

		categories.Categories = append(categories.Categories, &pb.Category{
			Id:          categoryResult.ID,
			Name:        categoryResult.Name,
			Description: categoryResult.Description,
		})
	}
}

func (cs *CategoryService) CreateCategoryStreamBidirectional(
	stream pb.CategoryService_CreateCategoryStreamBidirectionalServer,
) error {
	for {
		category, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		categoryResult, err := cs.DB.Create(category.Name, category.Description)
		if err != nil {
			return err
		}

		err = stream.Send(&pb.Category{
			Id:          categoryResult.ID,
			Name:        categoryResult.Name,
			Description: categoryResult.Description,
		})
		if err != nil {
			return err
		}
	}
}

func (cs *CategoryService) GetCategory(ctx context.Context, in *pb.CategoryGetRequest) (*pb.CategoryResponse, error) {
	category, err := cs.DB.Find(in.Id)
	if err != nil {
		return nil, err
	}

	result := &pb.Category{
		Id:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}

	return &pb.CategoryResponse{Category: result}, nil
}
