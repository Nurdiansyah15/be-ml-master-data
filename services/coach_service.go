package services

import (
	"fmt"
	"log"
	"os"
	"strings"

	"ml-master-data/models"

	"gorm.io/gorm"
)

// Fungsi untuk menghapus Coach dan semua relasi terkait
func DeleteCoach(db *gorm.DB, coach models.Coach) error {
	// Mulai transaksi
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// 1. Hapus semua CoachMatch yang terkait dengan Coach ini
	if err := tx.Where("coach_id = ?", coach.CoachID).Delete(&models.CoachMatch{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal menghapus CoachMatch: %w", err)
	}

	// 2. Hapus Coach itu sendiri
	if err := tx.Delete(&models.Coach{}, coach.CoachID).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal menghapus Coach: %w", err)
	}

	// Commit transaksi
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("gagal commit transaksi: %w", err)
	}

	if coach.Image != "" && coach.Image != "https://placehold.co/400x600" && strings.HasPrefix(coach.Image, os.Getenv("BASE_URL")) {
		coach.Image = strings.Replace(coach.Image, os.Getenv("BASE_URL")+"/", "", 1)
		// Cek apakah file Image lama ada di sistem
		if _, err := os.Stat(coach.Image); err == nil {
			// Jika file ada, hapus file Image lama dari folder images
			if err := os.Remove(coach.Image); err != nil {
				return fmt.Errorf("gagal menghapus gambar lama: %w", err)
			}
		} else if os.IsNotExist(err) {
			// Jika file tidak ada, lanjutkan ke tahap selanjutnya
			log.Printf("File gambar lama tidak ditemukan: %s", coach.Image)
		}
	}

	log.Printf("Coach dengan ID %d dan semua data terkait telah dihapus.", coach.CoachID)
	return nil
}
