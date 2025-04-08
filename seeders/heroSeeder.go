package seeders

import (
	"log"
	"ml-master-data/config"
	"ml-master-data/models"
	"strings"
)

// Fungsi untuk seeding heroes
func seedHeroes() []models.Hero {
	urls := []string{
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256370/Gatotkaca_ktvidy.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256370/Harley_lxkti2.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256369/Diggie_lbcbgz.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256368/Irithel_w9xgso.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256368/Grock_gblr3i.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256367/Hylos_zlbns6.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256367/Pharsa_ymufhl.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256366/Zhask_tapj74.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256365/Lesley_o5jxu9.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256364/Jawhead_c84l6w.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256363/Gusion_o5ax7w.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256363/Uranus_rxi3oh.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256363/Angela_fjtjxz.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256362/Martis_ffaaqc.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256361/Hanabi_fjwvs4.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256361/Aldous_skkvju.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256359/Chang_e_clww2k.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256359/Kaja_cdtzzs.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256358/Selena_vpzpwo.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256358/Hanzo_gzrkis.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256357/Claude_h9yzlg.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256357/Kimmy_oq3aam.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256356/Vale_gyw0yb.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256354/Lunox_rglxfn.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256354/Thamuz_dx6pqg.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256353/Harith_lvpni7.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256353/Minshitar_y3fhk1.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256352/Kadita_evfhva.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256352/Badang_ajxqsn.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256351/Khufra_ay7lb2.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256350/Carmilla_kq6ph6.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256349/Masha_zq9jyo.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256349/Cecilion_vorhup.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256348/Atlas_qks4sk.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256348/Luo_Yi_vf6to1.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256346/Popol_scvhfo.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256346/Belerick_ixg0dd.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256345/YuZhong_xijdqh.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256345/Baxia_dw6uyp.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256344/Wanwan_bglaub.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256344/Khaleed_ofcimg.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256343/Saber_nko8ru.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256343/Zilong_dlcghs.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256341/Barats_ciyudp.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256340/Brody_wg7mzk.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256340/Yve_oqyn8k.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256340/YiShunShin_dureop.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256339/Paquito_awmevr.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256338/Mathilda_iqtync.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256337/Yin_ukvw9d.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256337/Hayabusa_ydbozv.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256336/Karina_hsyhxd.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256336/Fanny_jiosny.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256335/Bane_scuqgb.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256333/Floryn_fc2yuw.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256333/Lancelot_vnmqsc.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256332/Odette_r5btb4.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256332/Valir_j1spvp.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256332/Silvanna_x2j0i3.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256330/Aamon_swxtll.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256330/Akai_fsq62k.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256329/Alpha_biuzdb.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256328/Beatrix_c9usct.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256328/Argus_gijifw.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256327/Eudora_gwbfts.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256326/Edith_xkzakw.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256326/Freya_xsrw59.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256325/Gloo_futgvw.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256324/Layla_bdwog7.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256323/Lapu_Lapu_pilhiy.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256323/Ling_js8kn9.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256322/Melissa_rwbq6f.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256321/Xavier_ikvua8.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256321/Natan_awli7h.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256320/Ruby_xkd1ly.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256319/Valentina_rq65ct.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256318/Kagura_nu9wsc.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256318/Tigreal_zvg2ha.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256317/Clint_y417g8.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256317/Aulus_gloast.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256315/Minotaur_r1ejzs.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256315/Benedetta_lboqjv.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256314/Phoveus_uttxku.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256313/Julian_ixgmvb.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256313/novaria_rdpc0i.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256313/Fredrinn_dktjnl.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256311/Lolita_twj8mo.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256310/Arlott_bicrhj.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256309/Nolan_oswyqq.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256309/Ixia_ntkx3q.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256308/Aurora1_mbjeee.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256307/Vexana_tgtqik.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256307/Sun_cd14to.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256307/Faramis_wslwu6.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256305/Chip_xxsbaa.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256305/Cici_eve9g8.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256304/Leomord_udfc7r.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256304/Helcurt_hwoqxy.png",
		"https://res.cloudinary.com/dnbreym94/image/upload/v1729256303/Zhuxin_dhkutl.png",
	}

	var heroes []models.Hero
	for _, url := range urls {
		name := extractHeroName(url)
		hero := models.Hero{Name: name, Image: url}
		heroes = append(heroes, hero)
		config.DB.Create(&hero)
	}

	log.Println("Heroes seeded successfully!")
	return heroes
}

// Fungsi untuk mengekstrak nama hero dari URL
func extractHeroName(url string) string {
	parts := strings.Split(url, "/")
	filename := parts[len(parts)-1]
	name := strings.Split(filename, "_")[0]
	return strings.ReplaceAll(name, "-", " ")
}
