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

// Fungsi untuk menghapus Hero dan semua relasi terkait
func DeleteHero(db *gorm.DB, hero models.Hero) error {
	// Mulai transaksi
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	//hapus priority pick
	if err := tx.Where("hero_id = ?", hero.HeroID).Delete(&models.PriorityPick{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal menghapus PriorityPick: %w", err)
	}

	// hapus priorityban
	if err := tx.Where("hero_id = ?", hero.HeroID).Delete(&models.PriorityBan{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal menghapus PriorityBan: %w", err)
	}

	// hapus flexpick
	if err := tx.Where("hero_id = ?", hero.HeroID).Delete(&models.FlexPick{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal menghapus FlexPick: %w", err)
	}

	// 1. Hapus HeroPickGame terkait
	if err := tx.Where("hero_pick_id IN (SELECT hero_pick_id FROM hero_picks WHERE hero_id = ?)", hero.HeroID).Delete(&models.HeroPickGame{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// hapus heropikc
	if err := tx.Where("hero_id = ?", hero.HeroID).Delete(&models.HeroPick{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal menghapus HeroPick: %w", err)
	}

	//  hapus herobangame
	if err := tx.Where("hero_ban_id IN (SELECT hero_ban_id FROM hero_bans WHERE hero_id = ?)", hero.HeroID).Delete(&models.HeroBanGame{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// hapus heroban
	if err := tx.Where("hero_id = ?", hero.HeroID).Delete(&models.HeroBan{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal menghapus HeroBan: %w", err)
	}

	// hapus explaner
	if err := tx.Where("hero_id = ?", hero.HeroID).Delete(&models.Explaner{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal menghapus Explaner: %w", err)
	}

	// hapus goldlaner
	if err := tx.Where("hero_id = ?", hero.HeroID).Delete(&models.Goldlaner{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal menghapus Goldlaner: %w", err)
	}

	// hapus triomidhero
	type trioMidID struct {
		TrioMidID uint
	}

	trioMidIDs := []trioMidID{}
	if err := tx.Raw("SELECT trio_mid_id FROM trio_mid_heros WHERE hero_id = ?", hero.HeroID).Scan(&trioMidIDs).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal mengambil trioMidID: %w", err)
	}

	// hapus trioMidHero
	if err := tx.Where("hero_id = ?", hero.HeroID).Delete(&models.TrioMidHero{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal menghapus TrioMidHero: %w", err)
	}

	for _, trioMidID := range trioMidIDs {
		var trioMidHeroCount int64
		if err := tx.Model(&models.TrioMidHero{}).Where("trio_mid_id = ?", trioMidID.TrioMidID).Count(&trioMidHeroCount).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("gagal menghitung trioMidHeroCount: %w", err)
		}
		if trioMidHeroCount == 0 {
			// hapus trioMid
			if err := tx.Where("trio_mid_id = ?", trioMidID.TrioMidID).Delete(&models.TrioMid{}).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("gagal menghapus TrioMid: %w", err)
			}
		}

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

		// Ambil Public ID dari URL gambar
		publicID := utils.ExtractPublicID(hero.Image)

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

		log.Printf("Gambar hero berhasil dihapus dari Cloudinary: %s", hero.Image)
	}

	log.Printf("Hero dengan ID %d dan semua data terkait telah dihapus.", hero.HeroID)
	return nil
}
