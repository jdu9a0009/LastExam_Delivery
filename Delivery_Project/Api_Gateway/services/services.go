package services

import (
	"api-gateway-service/config"
	order_service "api-gateway-service/genproto/order_service"
	product_service "api-gateway-service/genproto/product_service"
	user_service "api-gateway-service/genproto/user_service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ServiceManagerI interface {
	// Product Service
	ProductService() product_service.ProductServiceClient
	CategoryService() product_service.CategoryServiceClient

	// UserService
	BranchService() user_service.BranchServiceClient
	UserService() user_service.UserServiceClient
	CourierService() user_service.CourierServiceClient
	ClientService() user_service.ClientServiceClient

	//  Order Service
	DeliveryTariffService() order_service.DeliveryTariffServiceClient
	OrderService() order_service.OrderServiceClient
}

type grpcClients struct {
	// // Product Service
	productService  product_service.ProductServiceClient
	categoryService product_service.CategoryServiceClient

	// // UserService
	branchService  user_service.BranchServiceClient
	userService    user_service.UserServiceClient
	courierService user_service.CourierServiceClient
	clientService  user_service.ClientServiceClient

	// // Order Service
	deliveryTariffService order_service.DeliveryTariffServiceClient
	orderService          order_service.OrderServiceClient
}

func NewGrpcClients(cfg config.Config) (ServiceManagerI, error) {
	// // Product Microservice
	connProductService, err := grpc.Dial(
		cfg.ProductServiceHost+cfg.ProductGRPCPort,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	// User Microservice
	connUserService, err := grpc.Dial(
		cfg.UserServiceHost+cfg.UserGRPCPort,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	// // Order Microservice
	connOrderService, err := grpc.Dial(
		cfg.OrderServiceHost+cfg.OrderGRPCPort,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &grpcClients{
		// // Product Service
		productService:  product_service.NewProductServiceClient(connProductService),
		categoryService: product_service.NewCategoryServiceClient(connProductService),
		// // User Service
		branchService:  user_service.NewBranchServiceClient(connUserService),
		userService:    user_service.NewUserServiceClient(connUserService),
		clientService:  user_service.NewClientServiceClient(connUserService),
		courierService: user_service.NewCourierServiceClient(connUserService),

		// // Order Service
		deliveryTariffService: order_service.NewDeliveryTariffServiceClient(connOrderService),
		orderService:          order_service.NewOrderServiceClient(connOrderService),
	}, nil
}

// // Product Service
func (g *grpcClients) ProductService() product_service.ProductServiceClient {
	return g.productService
}

func (g *grpcClients) CategoryService() product_service.CategoryServiceClient {
	return g.categoryService
}

// User Service
func (g *grpcClients) BranchService() user_service.BranchServiceClient {
	return g.branchService
}

func (g *grpcClients) UserService() user_service.UserServiceClient {
	return g.userService
}
func (g *grpcClients) CourierService() user_service.CourierServiceClient {
	return g.courierService
}

func (g *grpcClients) ClientService() user_service.ClientServiceClient {
	return g.clientService
}

// // Order Service
func (g *grpcClients) DeliveryTariffService() order_service.DeliveryTariffServiceClient {
	return g.deliveryTariffService
}

func (g *grpcClients) OrderService() order_service.OrderServiceClient {
	return g.orderService
}
