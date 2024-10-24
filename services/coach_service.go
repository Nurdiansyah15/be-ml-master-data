package services

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"ml-master-data/models"
	"ml-master-data/utils"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
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

		// Ambil Public ID dari URL gambar
		publicID := utils.ExtractPublicID(coach.Image)

		// Inisialisasi Cloudinary
		cld, err := cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))
		if err != nil {
			return fmt.Errorf("gagal menginisialisasi Cloudinary: %w", err)
		}

		// Hapus gambar dari Cloudinary
		_, err = cld.Upload.Destroy(context.Background(), uploader.DestroyParams{PublicID: publicID})
		if err != nil {
			return fmt.Errorf("gagal menghapus gambar dari Cloudinary: %w", err)
		}

		log.Printf("Gambar coach berhasil dihapus dari Cloudinary: %s", coach.Image)
	}

	log.Printf("Coach dengan ID %d dan semua data terkait telah dihapus.", coach.CoachID)
	return nil
}
