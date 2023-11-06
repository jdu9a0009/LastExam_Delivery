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

type ProductService struct {
	cfg     config.Config
	log     logger.LoggerI
	storage storage.StorageI
	product_service.UnimplementedProductServiceServer
}

func NewProductService(cfg config.Config, log logger.LoggerI, strg storage.StorageI) *ProductService {
	return &ProductService{
		cfg:     cfg,
		log:     log,
		storage: strg,
	}
}

func (b *ProductService) Create(ctx context.Context, req *product_service.CreateProductRequest) (*product_service.Response, error) {
	resp, err := b.storage.Product().Create(context.Background(), req)
	if err != nil {
		b.log.Error("error while creating product", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return resp, nil
}

func (b *ProductService) Get(ctx context.Context, req *product_service.IdRequest) (*product_service.Product, error) {
	resp, err := b.storage.Product().Get(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (b *ProductService) List(ctx context.Context, req *product_service.ListProductRequest) (*product_service.ListProductResponse, error) {
	Products, err := b.storage.Product().GetList(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return &product_service.ListProductResponse{Products: Products.Products,
		Count: Products.Count}, nil
}

func (s *ProductService) Update(ctx context.Context, req *product_service.UpdateProductRequest) (*product_service.Response, error) {
	resp, err := s.storage.Product().Update(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return &product_service.Response{Message: resp}, nil
}

func (s *ProductService) Delete(ctx context.Context, req *product_service.IdRequest) (*product_service.Response, error) {
	resp, err := s.storage.Product().Delete(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return &product_service.Response{Message: resp}, nil
}
