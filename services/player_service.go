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

		// Ambil Public ID dari URL gambar
		publicID := utils.ExtractPublicID(player.Image)

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

		log.Printf("Gambar player berhasil dihapus dari Cloudinary: %s", player.Image)
	}

	log.Printf("Player dengan ID %d dan semua data terkait telah dihapus.", player.PlayerID)
	return nil
}
