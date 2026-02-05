package core

// DependenciesRegistry maps feature keys to their corresponding Go package import paths.
var DependenciesRegistry = map[string][]string{
	// Web Frameworks
	"Gin": {
		"github.com/gin-gonic/gin",
	},
	"Fiber": {
		"github.com/gofiber/fiber/v2",
	},
	"Echo": {
		"github.com/labstack/echo/v4",
	},
	"Chi": {
		"github.com/go-chi/chi/v5",
	},

	// Database Drivers
	"SQLite": {
		"modernc.org/sqlite",
	},
	"PostgreSQL": {
		"github.com/lib/pq",
	},
	"MySQL": {
		"github.com/go-sql-driver/mysql",
	},
	"MongoDB": {
		"go.mongodb.org/mongo-driver/mongo",
	},

	// ORM
	"Gorm": {
		"gorm.io/gorm",
	},
	"GormPostgres": {
		"gorm.io/driver/postgres",
	},
	"GormMySQL": {
		"gorm.io/driver/mysql",
	},
	"GormSQLite": {
		"gorm.io/driver/sqlite",
	},

	// CLI Libraries
	"Cobra": {
		"github.com/spf13/cobra",
	},
	"Viper": {
		"github.com/spf13/viper",
	},

	// Bot Libraries
	"TelegramBotAPI": {
		"github.com/go-telegram-bot-api/telegram-bot-api/v5",
	},

	// Environment Variable Management
	"Godotenv": {
		"github.com/joho/godotenv",
	},

	// Testing Frameworks
	"Testify": {
		"github.com/stretchr/testify",
	},
	"Ginkgo": {
		"github.com/onsi/ginkgo/v2",
	},
	"Gomega": {
		"github.com/onsi/gomega",
	},

	// Logging Frameworks
	"Zap": {
		"go.uber.org/zap",
	},
	"Logrus": {
		"github.com/sirupsen/logrus",
	},
}

// GetPackages returns the list of packages for the given key from the DependenciesRegistry.
// If the key does not exist, it returns an empty slice instead of nil for safety.
func GetPackages(key string) []string {
	if packages, exists := DependenciesRegistry[key]; exists {
		return packages
	}

	return []string{}
}
