package repositories

import (
	"errors"
	"log"
	"telegram-inventory/models"

	"gorm.io/gorm"
)

func GetServerInfo(keyword string) (models.Server, error) {
	if keyword == "" {
		return models.Server{}, errors.New("keyword tidak boleh kosong")
	}

	var server models.Server
	if err := DB.Where("vm_name = ? OR ipoam = ?", keyword, keyword).First(&server).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Server dengan keyword '%s' tidak ditemukan", keyword)
			return models.Server{}, nil // Server tidak ditemukan
		}
		return server, err // Kesalahan lain
	}
	return server, nil
}
