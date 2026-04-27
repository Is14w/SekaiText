package model

// UnitDict maps unit keys to display names.
var UnitDict = map[string]string{
	"piapro":         "VIRTUAL SINGER",
	"light_sound":    "Leo/need",
	"idol":           "MORE MORE JUMP！",
	"street":         "Vivid BAD SQUAD",
	"theme_park":     "ワンダーランズ×ショウタイム",
	"school_refusal": "25時、ナイトコードで。",
}

// SekaiDict lists unit keys in order.
var SekaiDict = []string{"leo", "mmj", "street", "wonder", "nightcode"}

// CharacterDict contains all 31 characters with their names.
var CharacterDict = []CharacterInfo{
	{Index: 0, Name: "ichika", NameJ: "一歌", NameC: "一歌"},
	{Index: 1, Name: "saki", NameJ: "咲希", NameC: "咲希"},
	{Index: 2, Name: "honami", NameJ: "穂波", NameC: "穗波"},
	{Index: 3, Name: "shiho", NameJ: "志歩", NameC: "志步"},
	{Index: 4, Name: "minori", NameJ: "みのり", NameC: "实乃里"},
	{Index: 5, Name: "haruka", NameJ: "遥", NameC: "遥"},
	{Index: 6, Name: "airi", NameJ: "愛莉", NameC: "爱莉"},
	{Index: 7, Name: "shizuku", NameJ: "雫", NameC: "雫"},
	{Index: 8, Name: "kohane", NameJ: "こはね", NameC: "心羽"},
	{Index: 9, Name: "an", NameJ: "杏", NameC: "杏"},
	{Index: 10, Name: "akito", NameJ: "彰人", NameC: "彰人"},
	{Index: 11, Name: "touya", NameJ: "冬弥", NameC: "冬弥"},
	{Index: 12, Name: "tsukasa", NameJ: "司", NameC: "司"},
	{Index: 13, Name: "emu", NameJ: "えむ", NameC: "笑梦"},
	{Index: 14, Name: "nene", NameJ: "寧々", NameC: "宁宁"},
	{Index: 15, Name: "rui", NameJ: "類", NameC: "类"},
	{Index: 16, Name: "kanade", NameJ: "奏", NameC: "奏"},
	{Index: 17, Name: "mafuyu", NameJ: "まふゆ", NameC: "真冬"},
	{Index: 18, Name: "ena", NameJ: "絵名", NameC: "绘名"},
	{Index: 19, Name: "mizuki", NameJ: "瑞希", NameC: "瑞希"},
	{Index: 20, Name: "miku", NameJ: "ミク", NameC: "MIKU"},
	{Index: 21, Name: "rin", NameJ: "リン", NameC: "RIN"},
	{Index: 22, Name: "len", NameJ: "レン", NameC: "LEN"},
	{Index: 23, Name: "luka", NameJ: "ルカ", NameC: "LUKA"},
	{Index: 24, Name: "meiko", NameJ: "MEIKO", NameC: "MEIKO"},
	{Index: 25, Name: "kaito", NameJ: "KAITO", NameC: "KAITO"},
	{Index: 26, Name: "miku_band", NameJ: "ミク_LeoN", NameC: "MIKU"},
	{Index: 27, Name: "miku_idol", NameJ: "ミク_MMJ", NameC: "MIKU"},
	{Index: 28, Name: "miku_street", NameJ: "ミク_VBS", NameC: "MIKU"},
	{Index: 29, Name: "miku_park", NameJ: "ミク_WS", NameC: "MIKU"},
	{Index: 30, Name: "miku_nothing", NameJ: "ミク_25", NameC: "MIKU"},
}

// AreaDict maps area index to name.
var AreaDict = []string{
	"",
	"十字路口",
	"商业街",
	"购物中心",
	"音乐商店",
	"教室的SEKAI",
	"",
	"舞台的SEKAI",
	"街道的SEKAI",
	"奇幻仙境的SEKAI",
	"空无一人的SEKAI",
	"神山高中",
	"凤凰仙境乐园",
	"宫益坂女子学园",
}

// GreetDictSeason contains season names.
var GreetDictSeason = []struct {
	Ch string `json:"ch"`
	En string `json:"en"`
}{
	{Ch: "春", En: "spring"},
	{Ch: "夏", En: "summer"},
	{Ch: "秋", En: "autumn"},
	{Ch: "冬", En: "winter"},
}

// GreetDictCelebrate lists character indices in celebration order.
var GreetDictCelebrate = []int{3, 17, 23, 16, 25, 8, 6, 4, 18, 1, 12, 11, 15, 14, 9, 0, 19, 20, 13, 5, 2, 24, 10, 7, 21, 22}

// GreetHolidayInfo represents a holiday entry.
type GreetHolidayInfo struct {
	Ch string `json:"ch"`
	En string `json:"en"`
}

// GreetDictHoliday contains holiday names.
var GreetDictHoliday = []GreetHolidayInfo{
	{Ch: "新年", En: "newyear"},
	{Ch: "节分", En: "setsubun"},
	{Ch: "情人节", En: "valentine"},
	{Ch: "女儿节", En: "hinamatsuri"},
	{Ch: "白情", En: "whiteday"},
	{Ch: "愚人节", En: "aprilfool"},
	{Ch: "七夕", En: "tanabata"},
	{Ch: "万圣节", En: "halloween"},
	{Ch: "圣诞节", En: "christmas"},
	{Ch: "年末", En: "endofyear"},
}

// FindCharacterByJapaneseName searches characterDict by Japanese name.
func FindCharacterByJapaneseName(name string) (CharacterInfo, bool) {
	for _, c := range CharacterDict {
		if c.NameJ == name {
			return c, true
		}
	}
	return CharacterInfo{}, false
}

// VoiceMsToMainStoryID maps voice prefixes to unit keys.
var VoiceMsToMainStoryID = map[string]string{
	"band":   "light_sound",
	"idol":   "idol",
	"street": "street",
	"wonder": "theme_park",
	"night":  "school_refusal",
	"piapro": "piapro",
}
