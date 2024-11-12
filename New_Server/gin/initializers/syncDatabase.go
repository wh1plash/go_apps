package initializers

import "gin/models"

func SyncDatabase() {
	err := DB.AutoMigrate(&models.User{})
	if err != nil {
		return
	}
}
