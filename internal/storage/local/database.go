package local

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"github.com/spf13/viper"
	"github.com/williamfotso/acc/internal/core/models/assignment"
	"github.com/williamfotso/acc/internal/core/models/course"
	"github.com/williamfotso/acc/internal/core/models/user"
	"github.com/williamfotso/acc/internal/core/models"
)

var (
	dbInstances = make(map[uint]*gorm.DB)
	dbLock      sync.Mutex
)

// getLocalDB returns a user-specific SQLite database instance
func GetLocalDB(userID uint) (*gorm.DB, error) {
	dbLock.Lock()
	defer dbLock.Unlock()

	// Return cached instance if available
	if db, exists := dbInstances[userID]; exists {
		return db, nil
	}

	// Determine database path
	dbPath, err := getDBPath(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get DB path: %w", err)
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create DB directory: %w", err)
	}

	// Open database connection
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		PrepareStmt: true, // Better performance
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open SQLite database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get SQL DB: %w", err)
	}
	sqlDB.SetMaxOpenConns(1) // SQLite works best with single connection

	// Initialize schema
	if err := initializeSchema(db); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	// Cache the instance
	dbInstances[userID] = db

	return db, nil
}

func getDBPath(userID uint) (string, error) {
	// Check for custom path in config
	if customPath := viper.GetString("localdb.path"); customPath != "" {
		return filepath.Join(customPath, fmt.Sprintf("user_%d.db", userID)), nil
	}

	// Use OS-specific default location
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get config directory: %w", err)
	}

	return filepath.Join(
		configDir,
		"acc-homework",
		"data",
		fmt.Sprintf("user_%d.db", userID),
	), nil
}

func initializeSchema(db *gorm.DB) error {
	// Enable foreign key support for SQLite
	if err := db.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
		return fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	// Run migrations
	return db.AutoMigrate(
		&user.LocalUser{},
		&course.LocalCourse{},
		&models.LocalAssignmentType{},
		&models.LocalAssignmentStatus{},
		&assignment.LocalAssignment{},
	)
}

// CloseAll closes all database connections
func CloseAll() error {
	dbLock.Lock()
	defer dbLock.Unlock()

	var lastErr error
	for userID, db := range dbInstances {
		if sqlDB, err := db.DB(); err == nil {
			if err := sqlDB.Close(); err != nil {
				lastErr = err
			}
		}
		delete(dbInstances, userID)
	}
	return lastErr
}
