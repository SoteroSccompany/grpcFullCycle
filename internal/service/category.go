package service

import (
	"context"
	"fmt"
	"io"

	"github.com/Soter-Tec/grpc/internal/database"
	"github.com/Soter-Tec/grpc/internal/pb"
)

type CategoryService struct {
	pb.UnimplementedCategoryServiceServer
	CategoryDB database.Category
}

func NewCategoryService(categoryDB database.Category) *CategoryService {
	return &CategoryService{CategoryDB: categoryDB}
}

func (c *CategoryService) CreateCategory(ctx context.Context, in *pb.CreateCategoryRequest) (*pb.Category, error) {
	category, err := c.CategoryDB.Create(in.Name, in.Description)
	if err != nil {
		return nil, err
	}
	categoryResponse := &pb.Category{
		Id:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}
	return categoryResponse, nil
}

func (c *CategoryService) ListCategories(ctx context.Context, in *pb.Blank) (*pb.CategoryList, error) {
	categories, err := c.CategoryDB.FindAll()
	if err != nil {
		return nil, err
	}
	var categoryResponses []*pb.Category
	for _, category := range categories {
		categoryResponses = append(categoryResponses, &pb.Category{
			Id:          category.ID,
			Name:        category.Name,
			Description: category.Description,
		})
	}
	return &pb.CategoryList{Categories: categoryResponses}, nil
}

func (c *CategoryService) ListCategoryById(ctx context.Context, in *pb.CategoryListById) (*pb.Category, error) {
	categorie, err := c.CategoryDB.FindByID(in.Id)
	if err != nil {
		return nil, err
	}
	categoryResponse := &pb.Category{
		Id:          categorie.ID,
		Name:        categorie.Name,
		Description: categorie.Description,
	}
	return categoryResponse, nil
}

func (c *CategoryService) UpdateCategory(ctx context.Context, in *pb.UpdateCategoryRequest) (*pb.Category, error) {
	categorie, err := c.CategoryDB.FindByID(in.Id)
	if err != nil {
		return nil, err
	}
	categorie.Name = in.Name
	categorie.Description = in.Description
	_, err = c.CategoryDB.UpdateCategory(categorie.ID, categorie.Name, categorie.Description)
	if err != nil {
		return nil, err
	}
	categoryResponse := &pb.Category{
		Id:          categorie.ID,
		Name:        categorie.Name,
		Description: categorie.Description,
	}
	return categoryResponse, nil
}

func (c *CategoryService) CreateCategoryStream(stream pb.CategoryService_CreateCategoryStreamServer) error {
	categories := &pb.CategoryList{}

	for {
		category, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("End of stream")
			return stream.SendAndClose(categories)
		}
		if err != nil {
			fmt.Printf("Error creating category: %v\n", err)
			return err
		}
		categoryResult, err := c.CategoryDB.Create(category.Name, category.Description)
		if err != nil {
			fmt.Printf("Error creating category: %v\n", err)
			return err
		}
		categories.Categories = append(categories.Categories, &pb.Category{
			Id:          categoryResult.ID,
			Name:        categoryResult.Name,
			Description: categoryResult.Description,
		})
	}
}

func (c *CategoryService) CreateCategoryStreamBidirectional(stream pb.CategoryService_CreateCategoryStreamBidirectionalServer) error {
	for {
		category, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		categoryResult, err := c.CategoryDB.Create(category.Name, category.Description)
		if err != nil {
			return err
		}
		err = stream.Send(
			&pb.Category{
				Id:          categoryResult.ID,
				Name:        categoryResult.Name,
				Description: categoryResult.Description,
			},
		)
		if err != nil {
			return err
		}
	}
}
