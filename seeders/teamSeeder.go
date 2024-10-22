package seeders

import (
	"log"
	"ml-master-data/config"
	"ml-master-data/models"
)

// Fungsi untuk seeding heroes
func seedTeams() []models.Team {

	teams := []models.Team{
		{Name: "RRQ", Image: "https://res.cloudinary.com/dnbreym94/image/upload/v1729256279/600px-Rex_Regum_Qeon_allmode_dcgt4m.png"},
		{Name: "EVOS", Image: "https://res.cloudinary.com/dnbreym94/image/upload/v1729256279/600px-EVOS_Esports_allmode_knfb5x.png"},
		{Name: "ALTER EGO", Image: "https://res.cloudinary.com/dnbreym94/image/upload/v1729256279/ae-256_fdchvl.png"},
		{Name: "BIGETRON", Image: "https://res.cloudinary.com/dnbreym94/image/upload/v1729256278/btr-256_gtabwc.png"},
		{Name: "REBELION", Image: "https://res.cloudinary.com/dnbreym94/image/upload/v1729256278/rbl_new_logo_ttv5gw.png"},
		{Name: "GEEK FAM", Image: "https://res.cloudinary.com/dnbreym94/image/upload/v1729256278/geek-500_aaejld.png"},
		{Name: "DEWA ESPORT", Image: "https://res.cloudinary.com/dnbreym94/image/upload/v1729256278/dewa-united-500_mazqhb.png"},
		{Name: "FNATIC ONIC", Image: "https://res.cloudinary.com/dnbreym94/image/upload/v1729256278/fnoc_logo_500x500_vrahwy.png"},
		{Name: "TLID", Image: "https://res.cloudinary.com/dnbreym94/image/upload/v1729256277/TLID-Primary500x500_fhndo5.png"},
	}

	if err := config.DB.Create(&teams).Error; err != nil {
		log.Fatal("Failed to seed teams:", err)
	}

	return teams
}
