package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	firebase "firebase.google.com/go/v4"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/yohannes/sint-new-backend/internal/application"
	"github.com/yohannes/sint-new-backend/internal/config"
	httpdelivery "github.com/yohannes/sint-new-backend/internal/delivery/http"
	firebaseauth "github.com/yohannes/sint-new-backend/internal/infrastructure/auth"
	"github.com/yohannes/sint-new-backend/internal/infrastructure/database"
	cloudinarystorage "github.com/yohannes/sint-new-backend/internal/infrastructure/storage"
	"google.golang.org/api/option"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	var firebaseOptions []option.ClientOption
	if len(cfg.FirebaseJSON) > 0 {
		firebaseOptions = append(firebaseOptions, option.WithCredentialsJSON(cfg.FirebaseJSON))
	}
	firebaseApp, err := firebase.NewApp(ctx, nil, firebaseOptions...)
	if err != nil {
		log.Fatalf("initialize Firebase Admin SDK: %v", err)
	}
	firebaseClient, err := firebaseApp.Auth(ctx)
	if err != nil {
		log.Fatalf("initialize Firebase Auth client: %v", err)
	}
	cloudinaryClient, err := cloudinary.NewFromURL(cfg.CloudinaryURL)
	if err != nil {
		log.Fatalf("initialize Cloudinary client: %v", err)
	}
	databaseClient, err := database.NewPostgres(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("initialize PostgreSQL: %v", err)
	}
	defer databaseClient.Close()
	if err := databaseClient.ApplyMigrations(ctx, "migrations"); err != nil {
		log.Fatalf("apply database migrations: %v", err)
	}

	authenticator := firebaseauth.NewFirebaseVerifier(firebaseClient)
	storage := cloudinarystorage.NewCloudinaryStorage(cloudinaryClient, cfg.CloudinaryFolder)
	handler := httpdelivery.NewHandler(
		databaseClient,
		authenticator,
		application.UploadPhoto{Storage: storage},
		application.SyncUser{Users: databaseClient},
		application.ListProducts{Products: databaseClient},
		application.ListStores{Stores: databaseClient},
		application.ListFavorites{Favorites: databaseClient},
		application.AddFavorite{Favorites: databaseClient},
		application.RemoveFavorite{Favorites: databaseClient},
		application.GetProductDetail{Products: databaseClient},
		application.SubmitPrice{Updates: databaseClient},
	)
	server := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-stop
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("shutdown HTTP server: %v", err)
		}
	}()

	log.Printf("Sint New API listening on %s", cfg.HTTPAddr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
