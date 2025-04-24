package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/JustBrowsing/query-service/api"
	"github.com/JustBrowsing/query-service/api/handlers"
	"github.com/JustBrowsing/query-service/config"
	"github.com/JustBrowsing/query-service/internal/repositories/mongo"
	"github.com/JustBrowsing/query-service/internal/services"
	"github.com/JustBrowsing/query-service/pkg/db"
	"github.com/JustBrowsing/query-service/pkg/logger"
	"github.com/go-redis/redis/v8"
	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"
)

func main() {
	// Parse flags
	configFile := flag.String("config", "config.yaml", "Path to config file")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load(*configFile)
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	log := logger.New(cfg.Logging.Level, cfg.Logging.Format)
	defer log.Sync()

	// Create context with cancelation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Connect to MongoDB
	log.Info("Connecting to MongoDB...", zap.String("uri", cfg.MongoDB.URI))
	mongoClient, err := db.NewMongoDBClient(ctx, cfg.MongoDB)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB", zap.Error(err))
	}
	defer mongoClient.Disconnect(ctx)

	// Create MongoDB database
	mongodb := db.NewMongoDBDatabase(mongoClient, cfg.MongoDB)

	// Connect to Redis
	log.Info("Connecting to Redis...", zap.String("address", cfg.Redis.Address))
	redisOptions := &redis.Options{
		Addr:     cfg.Redis.Address,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
		PoolSize: cfg.Redis.PoolSize,
	}
	redisClient := redis.NewClient(redisOptions)
	defer redisClient.Close()

	// Test Redis connection
	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		log.Fatal("Failed to connect to Redis", zap.Error(err))
	}

	// Connect to Elasticsearch
	log.Info("Connecting to Elasticsearch...", zap.Strings("addresses", cfg.Elasticsearch.Addresses))
	elasticClient, err := elastic.NewClient(
		elastic.SetURL(cfg.Elasticsearch.Addresses...),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(true),
		elastic.SetHealthcheckTimeout(5 * time.Second),
	)
	if err != nil {
		log.Fatal("Failed to connect to Elasticsearch", zap.Error(err))
	}

	// Create repositories
	productRepo := mongo.NewProductRepository(mongodb)
	orderRepo := mongo.NewOrderRepository(mongodb)

	// Create product index in Elasticsearch
	productIndex := cfg.Elasticsearch.IndexPrefix + "products"
	exists, err := elasticClient.IndexExists(productIndex).Do(ctx)
	if err != nil {
		log.Fatal("Failed to check if Elasticsearch index exists", zap.Error(err))
	}

	if !exists {
		log.Info("Creating Elasticsearch index", zap.String("index", productIndex))
		_, err = elasticClient.CreateIndex(productIndex).Do(ctx)
		if err != nil {
			log.Fatal("Failed to create Elasticsearch index", zap.Error(err))
		}
	}

	// Create services
	productService := services.NewProductService(
		log,
		productRepo,
		redisClient,
		elasticClient,
		cfg.Redis.TTL,
		productIndex,
	)

	orderService := services.NewOrderService(
		log,
		orderRepo,
	)

	eventService := services.NewEventService(
		log,
		productService,
		orderService,
		cfg.Kafka.Topics,
	)

	// Create handlers
	productHandler := handlers.NewProductHandler(log, productService)
	orderHandler := handlers.NewOrderHandler(log, orderService)
	healthHandler := handlers.NewHealthHandler()

	// Start Kafka consumer in a separate goroutine
	go func() {
		log.Info("Starting Kafka consumer...")
		if err := eventService.StartConsumer(ctx, cfg.Kafka); err != nil {
			log.Error("Failed to start Kafka consumer", zap.Error(err))
			cancel() // Cancel the context to shutdown the application
		}
	}()

	// Create router
	router := api.Router(
		log,
		cfg.Server.BasePath,
		productHandler,
		orderHandler,
		healthHandler,
	)

	// Create HTTP server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: router,
	}

	// Start server in a separate goroutine
	go func() {
		log.Info("Starting server...", zap.Int("port", cfg.Server.Port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for termination signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	// Create a deadline to wait for
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	// Shutdown server gracefully
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatal("Server forced to shutdown", zap.Error(err))
	}

	log.Info("Server exiting")
}