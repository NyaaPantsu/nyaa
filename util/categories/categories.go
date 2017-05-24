package categories

import (
	"github.com/NyaaPantsu/nyaa/config"
)

var categories map[string]string

func GetCategories() map[string]string {
	if categories != nil {
		return categories
	}

	if config.IsSukebei() {
		categories = map[string]string{
			"1_":  "art",
			"1_1": "art_anime",
			"1_2": "art_doujinshi",
			"1_3": "art_games",
			"1_4": "art_manga",
			"1_5": "art_pictures",
			"2_":  "real_life",
			"2_1": "real_life_photobooks_and_pictures",
			"2_2": "real_life_videos",
		}
	} else {
		categories = map[string]string{
			"3_":   "anime",
			"3_12": "anime_amv",
			"3_5":  "anime_english_translated",
			"3_13": "anime_non_english_translated",
			"3_6":  "anime_raw",
			"2_":   "audio",
			"2_3":  "audio_lossless",
			"2_4":  "audio_lossy",
			"4_":   "literature",
			"4_7":  "literature_english_translated",
			"4_8":  "literature_raw",
			"4_14": "literature_non_english_translated",
			"5_":   "live_action",
			"5_9":  "live_action_english_translated",
			"5_10": "live_action_idol_pv",
			"5_18": "live_action_non_english_translated",
			"5_11": "live_action_raw",
			"6_":   "pictures",
			"6_15": "pictures_graphics",
			"6_16": "pictures_photos",
			"1_":   "software",
			"1_1":  "software_applications",
			"1_2":  "software_games",
		}
	}

	return categories
}

func CategoryExists(category string) bool {
	_, exists := GetCategories()[category]
	return exists
}

func GetCategoriesSelect(keepParent bool) map[string]string {
	categories := GetCategories()
	catSelect := make(map[string]string, len(categories))
	for k, v := range categories {
		if len(k) > 2 || keepParent {
			catSelect[v] = k
		}
	}
	return catSelect
}
