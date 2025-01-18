package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// URL struct represents the database table
type URL struct {
	ID       uint   `gorm:"primaryKey"`
	ShortURL string `gorm:"uniqueIndex"`
	LongURL  string `gorm:"uniqueIndex"` // Add unique constraint to LongURL
}

var db *gorm.DB

func main() {
	// Load environment variables with default values
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "shortly")
	baseURL := getEnv("BASE_URL", "http://localhost:8080")

	// Database connection
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Auto-migrate the schema
	if err := db.AutoMigrate(&URL{}); err != nil {
		log.Fatalf("Failed to migrate the database schema: %v", err)
	}

	// Set up the Gin router
	r := gin.Default()

	// Enable CORS with specific configuration
	corsConfig := cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
	}
	if gin.Mode() == gin.ReleaseMode {
		allowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
		if allowedOrigins != "" {
			corsConfig.AllowAllOrigins = false
			corsConfig.AllowOrigins = strings.Split(allowedOrigins, ",")
		}
	}
	r.Use(cors.New(corsConfig))

	// Health check endpoints
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Server is running!"})
	})

	r.GET("/health/db", func(c *gin.Context) {
		sqlDB, err := db.DB()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "Database connection failed"})
			return
		}

		if err := sqlDB.Ping(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "Database not responding"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "Database is healthy"})
	})

	// Shorten URL endpoint
	r.POST("/shorten", func(c *gin.Context) {
		var request struct {
			URL string `json:"url" binding:"required"`
		}
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		// Check if the long URL already exists in the database
		var existingURL URL
		result := db.Where("long_url = ?", request.URL).First(&existingURL)
		if result.Error == nil {
			// If the long URL exists, return the existing short URL
			c.JSON(http.StatusOK, gin.H{"short_url": baseURL + "/" + existingURL.ShortURL})
			return
		}

		// If the long URL doesn't exist, create a new short URL
		short := HashURL(request.URL)

		// Save to the database
		url := URL{ShortURL: short, LongURL: request.URL}
		if err := db.Create(&url).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save to database"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"short_url": baseURL + "/" + short})
	})

	// Redirect to original URL
	r.GET("/:short", func(c *gin.Context) {
		short := c.Param("short")

		var url URL
		if err := db.First(&url, "short_url = ?", short).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"})
			return
		}

		c.Redirect(http.StatusMovedPermanently, url.LongURL)
	})

	// Graceful server shutdown
	srv := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %s\n", err)
		}
	}()

	log.Println("Server is running on http://localhost:8080")

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down server...")

	// Context timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %s", err)
	}

	sqlDB, _ := db.DB()
	sqlDB.Close()
}

// HashURL generates a short hash for the URL
func HashURL(url string) string {
	hash := sha256.Sum256([]byte(url))
	return base64.URLEncoding.EncodeToString(hash[:])[:8]
}

// getEnv retrieves environment variables with a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
