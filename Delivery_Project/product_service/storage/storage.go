package storage

import (
	"context"
	pb "product_service/genproto"
)

type StorageI interface {
	Category() CategoryI
	Product() ProductI
}

type CategoryI interface {
	Create(context.Context, *pb.CreateCategoryRequest) (*pb.Response, error)
	Get(context.Context, *pb.IdRequest) (*pb.Category, error)
	GetList(context.Context, *pb.ListCategoryRequest) (*pb.ListCategoryResponse, error)
	Update(context.Context, *pb.UpdateCategoryRequest) (string, error)
	Delete(context.Context, *pb.IdRequest) (string, error)
}

type ProductI interface {
	Create(context.Context, *pb.CreateProductRequest) (*pb.Response, error)
	Get(context.Context, *pb.IdRequest) (*pb.Product, error)
	GetList(context.Context, *pb.ListProductRequest) (*pb.ListProductResponse, error)
	Update(context.Context, *pb.UpdateProductRequest) (string, error)
	Delete(context.Context, *pb.IdRequest) (string, error)
}
