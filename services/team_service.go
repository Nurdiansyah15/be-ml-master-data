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

// Fungsi untuk menghapus Team beserta semua relasi terkait termasuk Match
func DeleteTeam(db *gorm.DB, team models.Team) error {
	// Mulai transaksi
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	matches := []models.Match{}
	if err := tx.Where("team_a_id = ? OR team_b_id = ?", team.TeamID, team.TeamID).Find(&matches).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal mendapatkan semua Match: %w", err)
	}

	// 2. Panggil DeleteMatch untuk setiap Match terkait
	for _, match := range matches {
		if err := DeleteMatch(db, match.MatchID); err != nil {
			tx.Rollback()
			return err
		}
	}

	players := []models.Player{}
	if err := tx.Where("team_id = ?", team.TeamID).Find(&players).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal mendapatkan semua Player: %w", err)
	}

	//panggil DeletePlayer untuk setiap Player terkait
	for _, player := range players {
		if err := DeletePlayer(db, player); err != nil {
			tx.Rollback()
			return err
		}
	}

	coaches := []models.Coach{}
	if err := tx.Where("team_id = ?", team.TeamID).Find(&coaches).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal mendapatkan semua Coach: %w", err)
	}

	//panggil DeleteCoach untuk setiap Coach terkait
	for _, coach := range coaches {
		if err := DeleteCoach(db, coach); err != nil {
			tx.Rollback()
			return err
		}
	}

	// 7. Hapus Team itu sendiri
	if err := tx.Delete(&models.Team{}, team.TeamID).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal menghapus Team: %w", err)
	}

	// Commit transaksi
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("gagal commit transaksi: %w", err)
	}

	if team.Image != "" && team.Image != "https://placehold.co/400x600" && strings.HasPrefix(team.Image, os.Getenv("BASE_URL")) {
		// Ambil Public ID dari URL gambar
		publicID := utils.ExtractPublicID(team.Image)

		// Inisialisasi Cloudinary
		cld, err := cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))
		if err != nil {
			return fmt.Errorf("gagal menginisialisasi Cloudinary: %w", err)
		}

		// Hapus gambar lama dari Cloudinary
		_, err = cld.Upload.Destroy(context.Background(), uploader.DestroyParams{PublicID: publicID})
		if err != nil {
			return fmt.Errorf("gagal menghapus gambar lama dari Cloudinary: %w", err)
		}

		log.Printf("Gambar lama berhasil dihapus dari Cloudinary: %s", team.Image)
	}

	log.Printf("Team dengan ID %d dan semua data terkait telah dihapus.", team.TeamID)
	return nil
}
