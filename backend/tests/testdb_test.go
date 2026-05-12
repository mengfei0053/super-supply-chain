package tests

import (
	"testing"

	"super-supply-chain/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB 创建一个内存 SQLite 数据库，临时替换全局 models.DB，测试结束后自动恢复。
func setupTestDB(t *testing.T, modelsToMigrate ...any) {
	t.Helper()

	previousDB := models.DB
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}

	if len(modelsToMigrate) > 0 {
		if err := db.AutoMigrate(modelsToMigrate...); err != nil {
			t.Fatalf("migrate test database: %v", err)
		}
	}

	models.DB = db
	t.Cleanup(func() {
		sqlDB, err := db.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
		models.DB = previousDB
	})
}
