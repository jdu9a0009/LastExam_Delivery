package service

import (
	"context"
	"product_service/config"
	product_service "product_service/genproto"
	"product_service/pkg/logger"
	"product_service/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CategoryService struct {
	cfg     config.Config
	log     logger.LoggerI
	storage storage.StorageI
	product_service.UnimplementedCategoryServiceServer
}

func NewCategoryService(cfg config.Config, log logger.LoggerI, strg storage.StorageI) *CategoryService {
	return &CategoryService{
		cfg:     cfg,
		log:     log,
		storage: strg,
	}
}

func (b *CategoryService) Create(ctx context.Context, req *product_service.CreateCategoryRequest) (*product_service.Response, error) {
	resp, err := b.storage.Category().Create(context.Background(), req)
	if err != nil {
		b.log.Error("error while creating product", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return resp, nil
}

func (b *CategoryService) Get(ctx context.Context, req *product_service.IdRequest) (*product_service.Category, error) {
	resp, err := b.storage.Category().Get(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (b *CategoryService) List(ctx context.Context, req *product_service.ListCategoryRequest) (*product_service.ListCategoryResponse, error) {
	Categorys, err := b.storage.Category().GetList(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return &product_service.ListCategoryResponse{Categories: Categorys.Categories,
		Count: Categorys.Count}, nil
}

func (s *CategoryService) Update(ctx context.Context, req *product_service.UpdateCategoryRequest) (*product_service.Response, error) {
	resp, err := s.storage.Category().Update(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return &product_service.Response{Message: resp}, nil
}

func (s *CategoryService) Delete(ctx context.Context, req *product_service.IdRequest) (*product_service.Response, error) {
	resp, err := s.storage.Category().Delete(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return &product_service.Response{Message: resp}, nil
}
