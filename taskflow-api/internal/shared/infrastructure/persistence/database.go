package persistence

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewDatabase ouvre la connexion PostgreSQL via GORM.
func NewDatabase(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("impossible de se connecter à PostgreSQL: %v", err)
	}

	log.Println("connexion PostgreSQL établie")
	return db
}

// Migrate lance les migrations automatiques pour les modèles GORM donnés.
func Migrate(db *gorm.DB, models ...interface{}) {
	if err := db.AutoMigrate(models...); err != nil {
		log.Fatalf("erreur migration: %v", err)
	}
	log.Println("migrations exécutées")
}
