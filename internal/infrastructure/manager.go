package infrastructure

import (
	"fmt"

	"github.com/vashafaris/nucleus/internal/infrastructure/persistence/postgres"
	"github.com/vashafaris/nucleus/internal/infrastructure/persistence/redis"
	"github.com/vashafaris/nucleus/pkg/config"
)

// Manager holds all infrastructure connections
type Manager struct {
	DB    *postgres.DB
	Redis *redis.Client
	cfg   *config.Config
}

// NewManager creates a new infrastructure manager
func NewManager(cfg *config.Config) (*Manager, error) {
	// Initialize PostgreSQL
	db, err := postgres.NewConnection(&cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	// Initialize Redis
	redisClient, err := redis.NewConnection(&cfg.Redis)
	if err != nil {
		db.Close() // Clean up if Redis fails
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &Manager{
		DB:    db,
		Redis: redisClient,
		cfg:   cfg,
	}, nil
}

// Health checks all infrastructure components
func (m *Manager) Health() map[string]error {
	health := make(map[string]error)

	// Check PostgreSQL
	if err := m.DB.Health(); err != nil {
		health["postgres"] = err
	}

	// Check Redis
	if err := m.Redis.Health(); err != nil {
		health["redis"] = err
	}

	return health
}

// Close closes all infrastructure connections
func (m *Manager) Close() error {
	var errs []error

	if err := m.DB.Close(); err != nil {
		errs = append(errs, fmt.Errorf("failed to close PostgreSQL: %w", err))
	}

	if err := m.Redis.Close(); err != nil {
		errs = append(errs, fmt.Errorf("failed to close Redis: %w", err))
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing connections: %v", errs)
	}

	return nil
}
