package services

import (
	"fmt"
	"log"
	"os"
	"strings"

	"ml-master-data/models"

	"gorm.io/gorm"
)

// Fungsi untuk menghapus Hero dan semua relasi terkait
func DeleteHero(db *gorm.DB, hero models.Hero) error {
	// Mulai transaksi
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// 1. Hapus HeroPickGame terkait
	if err := tx.Where("hero_id = ?", hero.HeroID).Delete(&models.HeroPickGame{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal menghapus HeroPickGame: %w", err)
	}

	// 2. Hapus HeroBanGame terkait
	if err := tx.Where("hero_id = ?", hero.HeroID).Delete(&models.HeroBanGame{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal menghapus HeroBanGame: %w", err)
	}

	// 3. Hapus TrioMidHero terkait
	if err := tx.Where("hero_id = ?", hero.HeroID).Delete(&models.TrioMidHero{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal menghapus TrioMidHero: %w", err)
	}

	// 4. Hapus HeroPick terkait
	if err := tx.Where("hero_id = ?", hero.HeroID).Delete(&models.HeroPick{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal menghapus HeroPick: %w", err)
	}

	// 5. Hapus HeroBan terkait
	if err := tx.Where("hero_id = ?", hero.HeroID).Delete(&models.HeroBan{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal menghapus HeroBan: %w", err)
	}

	// 6. Hapus Hero itu sendiri
	if err := tx.Delete(&models.Hero{}, hero.HeroID).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal menghapus Hero: %w", err)
	}

	// Commit transaksi
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("gagal commit transaksi: %w", err)
	}

	if hero.Image != "" && hero.Image != "https://placehold.co/400x600" && strings.HasPrefix(hero.Image, os.Getenv("BASE_URL")) {
		hero.Image = strings.Replace(hero.Image, os.Getenv("BASE_URL")+"/", "", 1)
		// Cek apakah file Image lama ada di sistem
		if _, err := os.Stat(hero.Image); err == nil {
			// Jika file ada, hapus file Image lama dari folder images
			if err := os.Remove(hero.Image); err != nil {
				return fmt.Errorf("gagal menghapus gambar lama: %w", err)
			}
		} else if os.IsNotExist(err) {
			// Jika file tidak ada, lanjutkan ke tahap selanjutnya
			log.Printf("File gambar lama tidak ditemukan: %s", hero.Image)
		}
	}

	log.Printf("Hero dengan ID %d dan semua data terkait telah dihapus.", hero.HeroID)
	return nil
}
