package services

import (
	"fmt"
	"log"
	"os"
	"strings"

	"ml-master-data/models"

	"gorm.io/gorm"
)

// Fungsi untuk menghapus Player dan semua relasi terkait
func DeletePlayer(db *gorm.DB, player models.Player) error {
	// Mulai transaksi
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// 1. Hapus semua PlayerMatch yang terkait dengan Player ini
	if err := tx.Where("player_id = ?", player.PlayerID).Delete(&models.PlayerMatch{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal menghapus PlayerMatch: %w", err)
	}

	// 2. Hapus Player itu sendiri
	if err := tx.Delete(&models.Player{}, player.PlayerID).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal menghapus Player: %w", err)
	}

	// Commit transaksi
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("gagal commit transaksi: %w", err)
	}

	if player.Image != "" && player.Image != "https://placehold.co/400x600" && strings.HasPrefix(player.Image, os.Getenv("BASE_URL")) {
		player.Image = strings.Replace(player.Image, os.Getenv("BASE_URL")+"/", "", 1)
		// Cek apakah file Image lama ada di sistem
		if _, err := os.Stat(player.Image); err == nil {
			// Jika file ada, hapus file Image lama dari folder images
			if err := os.Remove(player.Image); err != nil {
				return fmt.Errorf("gagal menghapus gambar lama: %w", err)
			}
		} else if os.IsNotExist(err) {
			// Jika file tidak ada, lanjutkan ke tahap selanjutnya
			log.Printf("File gambar lama tidak ditemukan: %s", player.Image)
		}
	}

	log.Printf("Player dengan ID %d dan semua data terkait telah dihapus.", player.PlayerID)
	return nil
}
