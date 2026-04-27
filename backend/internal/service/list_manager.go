package service

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"sekaitext/backend/internal/model"
)

// ListManager manages story metadata (events, cards, main story, etc.).
type ListManager struct {
	Events     []EventEntry
	Festivals  []FestivalEntry
	Cards      []CardEntry
	MainStory  []MainStoryEntry
	AreaTalks  []AreaTalkEntry
	Greets     []GreetEntry
	Specials   []SpecialEntry
	Setting    map[string]interface{}

	settingDir string
	DBurl      string

	// For area talk navigation
	AreaTalkByTime []AreaTalkTimeEntry
	ChapterScenario []ChapterScenarioEntry

	// CDN URLs
	baseUrls map[string]string
}

// EventEntry mirrors events.json structure.
type EventEntry struct {
	ID               int                    `json:"id"`
	KdyicrID         int                    `json:"kdyicr_id"`
	Title            string                 `json:"title"`
	Name             string                 `json:"name"`
	Chapters         []EventChapter         `json:"chapters"`
	Cards            []int                  `json:"cards"`
	InferredVoiceIDs map[string]interface{} `json:"inferredVoiceIDs,omitempty"`
}

// EventChapter represents a chapter in an event.
type EventChapter struct {
	Title     string `json:"title"`
	AssetName string `json:"assetName"`
}

// FestivalEntry mirrors festivals.json structure.
type FestivalEntry struct {
	ID            int    `json:"id"`
	IsBirthday    bool   `json:"isBirthday"`
	Cards         []int  `json:"cards"`
	Collaboration string `json:"collaboration,omitempty"`
	LevelUp       bool   `json:"levelup,omitempty"`
}

// CardEntry mirrors cards.json structure.
type CardEntry struct {
	ID          int  `json:"id"`
	CharacterID int  `json:"characterId"`
	CardNo      string `json:"cardNo"`
	Birthday    bool `json:"birthday"`
	LevelUp     bool `json:"levelup,omitempty"`
}

// MainStoryEntry mirrors mainStory.json structure.
type MainStoryEntry struct {
	Unit      string              `json:"unit"`
	AssetName string              `json:"assetName"`
	Chapters  []EventChapter      `json:"chapters"`
}

// AreaTalkEntry mirrors areatalks.json structure.
type AreaTalkEntry struct {
	ID             int    `json:"id"`
	TalkID         string `json:"talkid"`
	AreaID         int    `json:"areaId"`
	CharacterIDs   []int  `json:"characterIds"`
	ScenarioID     string `json:"scenarioId"`
	Type           string `json:"type"`
	AddEventID     int    `json:"addEventId"`
	ReleaseEventID int    `json:"releaseEventId"`
}

// GreetEntry mirrors greets.json structure.
type GreetEntry struct {
	Theme  GreetTheme    `json:"theme"`
	Year   int           `json:"year"`
	Greets []GreetItem   `json:"greets"`
}

// GreetTheme represents a greet theme.
type GreetTheme struct {
	Ch string `json:"ch"`
	En string `json:"en"`
}

// GreetItem represents a single greet entry.
type GreetItem struct {
	CharacterID int    `json:"characterId"`
	Text        string `json:"text"`
}

// SpecialEntry mirrors specials.json structure.
type SpecialEntry struct {
	Title    string `json:"title"`
	DirName  string `json:"dirName"`
	FileName string `json:"fileName"`
}

// AreaTalkTimeEntry is used for "按时间" sorting.
type AreaTalkTimeEntry struct {
	AddEventID     int  `json:"addEventId"`
	ReleaseEventID int  `json:"releaseEventId"`
	Limited        bool `json:"limited"`
	Monthly        bool `json:"monthly"`
}

// ChapterScenarioEntry stores resolved scenario info for a chapter.
type ChapterScenarioEntry struct {
	ID         int    `json:"id"`
	ScenarioID string `json:"scenarioId"`
	TalkID     string `json:"talkid"`
	IsSeparator bool  `json:"isSeparator,omitempty"`
}

// NewListManager creates and loads metadata from the setting directory.
func NewListManager(settingDir string) *ListManager {
	lm := &ListManager{
		settingDir: settingDir,
		Setting:    make(map[string]interface{}),
		baseUrls: map[string]string{
			"best":    "https://storage.sekai.best/sekai-jp-assets/",
			"uni":     "https://assets.unipjsk.com/",
			"haruki":  "https://bot-assets.haruki.seiunx.com/assets/",
			"harukiJP": "https://sekai-assets-bdf29c81.seiunx.net/jp-assets/",
		},
	}
	lm.loadSetting()
	lm.loadAll()
	return lm
}

func (lm *ListManager) loadSetting() {
	path := filepath.Join(lm.settingDir, "setting.json")
	data, err := os.ReadFile(path)
	if err == nil {
		json.Unmarshal(data, &lm.Setting)
	}
}

func (lm *ListManager) loadAll() {
	lm.Events = loadJSONFile[[]EventEntry](lm.settingDir, "events.json")
	lm.Festivals = loadJSONFile[[]FestivalEntry](lm.settingDir, "festivals.json")
	lm.Cards = loadJSONFile[[]CardEntry](lm.settingDir, "cards.json")
	lm.MainStory = loadJSONFile[[]MainStoryEntry](lm.settingDir, "mainStory.json")
	lm.AreaTalks = loadJSONFile[[]AreaTalkEntry](lm.settingDir, "areatalks.json")
	lm.Greets = loadJSONFile[[]GreetEntry](lm.settingDir, "greets.json")
	lm.Specials = loadJSONFile[[]SpecialEntry](lm.settingDir, "specials.json")
	log.Println("All metadata loaded")
}

func loadJSONFile[T any](dir, fileName string) T {
	var zero T
	path := filepath.Join(dir, fileName)
	data, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Warning: could not load %s: %v", fileName, err)
		return zero
	}
	if err := json.Unmarshal(data, &zero); err != nil {
		log.Printf("Warning: could not parse %s: %v", fileName, err)
		return zero
	}
	return zero
}

// --- Story Type constants (Chinese labels) ---
const (
	StoryLabelEvent       = "活动剧情"
	StoryLabelMainStory   = "主线剧情"
	StoryLabelCardEvent   = "活动卡面"
	StoryLabelCardSpecial = "特殊卡面"
	StoryLabelCardInit    = "初始卡面"
	StoryLabelCardUpgrade = "升级卡面"
	StoryLabelAreaTalkInit = "初始地图对话"
	StoryLabelAreaTalkUpgrade = "升级地图对话"
	StoryLabelAreaTalkExtra = "追加地图对话"
	StoryLabelGreet       = "主界面语音"
	StoryLabelSpecial     = "特殊剧情"
)

// GetStoryTypes returns available story type names (Chinese labels).
func (lm *ListManager) GetStoryTypes() []string {
	return []string{
		StoryLabelEvent,
		StoryLabelMainStory,
		StoryLabelCardEvent,
		StoryLabelCardSpecial,
		StoryLabelCardInit,
		StoryLabelCardUpgrade,
		StoryLabelAreaTalkInit,
		StoryLabelAreaTalkUpgrade,
		StoryLabelAreaTalkExtra,
		StoryLabelGreet,
		StoryLabelSpecial,
	}
}

// GetStorySorts returns sort options for a given story type.
func (lm *ListManager) GetStorySorts(storyType string) []model.StorySort {
	switch storyType {
	case StoryLabelAreaTalkInit, StoryLabelAreaTalkUpgrade:
		return []model.StorySort{
			{Label: "按人物", Value: "character"},
			{Label: "按地点", Value: "area"},
		}
	case StoryLabelGreet:
		return []model.StorySort{
			{Label: "按人物", Value: "character"},
			{Label: "按时间", Value: "time"},
		}
	case StoryLabelAreaTalkExtra:
		return []model.StorySort{
			{Label: "按人物", Value: "character"},
			{Label: "按时间", Value: "time"},
			{Label: "按地点", Value: "area"},
		}
	default:
		return nil
	}
}

// GetStoryIndexList returns index options for a story type and sort.
func (lm *ListManager) GetStoryIndexList(storyType, sort string) []model.StoryIndex {
	var indices []model.StoryIndex

	switch storyType {
	case StoryLabelMainStory:
		for _, unit := range lm.MainStory {
			name := model.UnitDict[unit.Unit]
			indices = append(indices, model.StoryIndex{
				Label: name,
				Value: name,
			})
		}

	case StoryLabelEvent, StoryLabelCardEvent:
		eventSum := len(lm.Events)
		for i := len(lm.Events) - 1; i >= 0; i-- {
			ev := lm.Events[i]
			label := strconv.Itoa(eventSum-i) + " " + ev.Title
			indices = append(indices, model.StoryIndex{
				Label: label,
				Value: strconv.Itoa(eventSum - 1 - i),
			})
		}

	case StoryLabelCardSpecial:
		for i := len(lm.Festivals) - 1; i >= 0; i-- {
			f := lm.Festivals[i]
			var label string
			if f.Collaboration != "" {
				label = f.Collaboration
			} else if f.IsBirthday {
				idx := f.ID
				year := 2021 + (idx+2)/4
				month := (idx+2)%4*3 + 1
				label = "Birthday " + strconv.Itoa(year) + " " + padZero(month) + "-" + padZero(month+2)
			} else {
				idx := f.ID
				year := 2021 + idx/4
				month := idx%4*3 + 1
				label = "Festival " + strconv.Itoa(year) + " " + padZero(month)
			}
			indices = append(indices, model.StoryIndex{Label: label, Value: strconv.Itoa(len(lm.Festivals) - 1 - i)})
		}

	case StoryLabelCardInit, StoryLabelCardUpgrade:
		for idx, char := range model.CharacterDict[:26] {
			indices = append(indices, model.StoryIndex{Label: char.NameJ, Value: strconv.Itoa(idx)})
			if idx%4 == 3 && idx < 20 {
				indices = append(indices, model.StoryIndex{Label: "-", Value: "-"})
			}
		}

	case StoryLabelAreaTalkInit, StoryLabelAreaTalkUpgrade:
		if sort == "character" {
			for idx, char := range model.CharacterDict[:26] {
				indices = append(indices, model.StoryIndex{Label: char.NameJ, Value: strconv.Itoa(idx)})
				if idx%4 == 3 && idx < 20 {
					indices = append(indices, model.StoryIndex{Label: "-", Value: "-"})
				}
			}
		} else if sort == "time" {
			lm.buildAreaTalkByTime()
			for i := len(lm.AreaTalkByTime) - 1; i >= 0; i-- {
				indices = append(indices, model.StoryIndex{Label: "time", Value: strconv.Itoa(i)})
			}
		} else if sort == "area" {
			for _, area := range model.AreaDict {
				if area != "" {
					indices = append(indices, model.StoryIndex{Label: area, Value: area})
				}
			}
		}

	case StoryLabelGreet:
		if sort == "character" {
			for idx, char := range model.CharacterDict {
				indices = append(indices, model.StoryIndex{Label: char.NameJ, Value: strconv.Itoa(idx)})
				if (idx%4 == 3 && idx < 20) || idx == 25 {
					indices = append(indices, model.StoryIndex{Label: "-", Value: "-"})
				}
			}
		} else if sort == "time" {
			for i := len(lm.Greets) - 1; i >= 0; i-- {
				g := lm.Greets[i]
				label := g.Theme.Ch + " " + strconv.Itoa(g.Year)
				indices = append(indices, model.StoryIndex{Label: label, Value: strconv.Itoa(i)})
			}
		}

	case StoryLabelSpecial:
		for i := len(lm.Specials) - 1; i >= 0; i-- {
			indices = append(indices, model.StoryIndex{
				Label: lm.Specials[i].Title,
				Value: strconv.Itoa(i),
			})
		}
	}

	return indices
}

// GetStoryChapterList returns chapters for a given story.
func (lm *ListManager) GetStoryChapterList(storyType, sort, index string) []model.StoryChapter {
	lm.ChapterScenario = nil
	idx := parseIndex(index)
	var chapters []model.StoryChapter

	switch storyType {
	case StoryLabelMainStory:
		unitIdx := idx
		if unitIdx >= 0 && unitIdx < len(lm.MainStory) {
			for ci, chapter := range lm.MainStory[unitIdx].Chapters {
				var epNo int
				if unitIdx == 0 {
					epNo = ci%4 + 1
				} else {
					epNo = ci
				}
				chapters = append(chapters, model.StoryChapter{
					Number: ci,
					Label:  strconv.Itoa(epNo) + " " + chapter.Title,
				})
			}
		}

	case StoryLabelEvent:
		eventIdx := len(lm.Events) - idx
		if eventIdx >= 1 && eventIdx <= len(lm.Events) {
			for ci, chapter := range lm.Events[eventIdx-1].Chapters {
				chapters = append(chapters, model.StoryChapter{
					Number: ci,
					Label:  strconv.Itoa(ci+1) + " " + chapter.Title,
				})
			}
		}

	case StoryLabelCardEvent:
		content := lm.Events
		contentIdx := len(content) - idx
		if contentIdx >= 1 && contentIdx <= len(content) {
			for _, cardID := range content[contentIdx-1].Cards {
				if cardID >= 1 && cardID <= len(lm.Cards) {
					char := model.CharacterDict[lm.Cards[cardID-1].CharacterID-1]
					chapters = append(chapters,
						model.StoryChapter{Number: len(chapters), Label: char.NameJ + " 前篇"},
						model.StoryChapter{Number: len(chapters), Label: char.NameJ + " 后篇"},
						model.StoryChapter{Number: len(chapters), Label: "-"},
					)
				}
			}
			if len(chapters) > 0 {
				chapters = chapters[:len(chapters)-1]
			}
		}

	case StoryLabelCardSpecial:
		content := lm.Festivals
		contentIdx := len(content) - idx
		if contentIdx >= 1 && contentIdx <= len(content) {
			for _, cardID := range content[contentIdx-1].Cards {
				if cardID >= 1 && cardID <= len(lm.Cards) {
					char := model.CharacterDict[lm.Cards[cardID-1].CharacterID-1]
					chapters = append(chapters,
						model.StoryChapter{Number: len(chapters), Label: char.NameJ + " 前篇"},
						model.StoryChapter{Number: len(chapters), Label: char.NameJ + " 后篇"},
						model.StoryChapter{Number: len(chapters), Label: "-"},
					)
				}
			}
			if len(chapters) > 0 {
				chapters = chapters[:len(chapters)-1]
			}
		}

	case StoryLabelCardInit:
		if idx >= 0 && idx < 26 {
			chapters = append(chapters,
				model.StoryChapter{Number: 0, Label: "1☆ 前篇"},
				model.StoryChapter{Number: 1, Label: "1☆ 后篇"},
				model.StoryChapter{Number: 2, Label: "1☆ 其他"},
				model.StoryChapter{Number: 3, Label: "2☆ 前篇"},
				model.StoryChapter{Number: 4, Label: "2☆ 后篇"},
				model.StoryChapter{Number: 5, Label: "2☆ 其他"},
				model.StoryChapter{Number: 6, Label: "3☆ 前篇"},
				model.StoryChapter{Number: 7, Label: "3☆ 后篇"},
				model.StoryChapter{Number: 8, Label: "3☆ 其他"},
				model.StoryChapter{Number: 9, Label: "4☆ 前篇"},
				model.StoryChapter{Number: 10, Label: "4☆ 后篇"},
				model.StoryChapter{Number: 11, Label: "4☆ 其他"},
			)
		}

	case StoryLabelCardUpgrade:
		if idx >= 0 && idx < 26 {
			chapters = append(chapters,
				model.StoryChapter{Number: 0, Label: "前篇"},
				model.StoryChapter{Number: 1, Label: "后篇"},
				model.StoryChapter{Number: 2, Label: "其他"},
			)
		}

	case StoryLabelAreaTalkInit:
		lm.buildAreaTalkChapterScenario("init", sort)
		for ci := range lm.ChapterScenario {
			if lm.ChapterScenario[ci].IsSeparator {
				chapters = append(chapters, model.StoryChapter{Number: ci, Label: "-"})
			} else {
				chapters = append(chapters, model.StoryChapter{Number: ci, Label: lm.ChapterScenario[ci].ScenarioID})
			}
		}

	case StoryLabelAreaTalkUpgrade:
		lm.buildAreaTalkChapterScenario("upgrade", sort)
		for ci := range lm.ChapterScenario {
			if lm.ChapterScenario[ci].IsSeparator {
				chapters = append(chapters, model.StoryChapter{Number: ci, Label: "-"})
			} else {
				chapters = append(chapters, model.StoryChapter{Number: ci, Label: lm.ChapterScenario[ci].ScenarioID})
			}
		}

	case StoryLabelAreaTalkExtra:
		lm.buildAreaTalkChapterScenario("extra", sort)
		for ci := range lm.ChapterScenario {
			if lm.ChapterScenario[ci].IsSeparator {
				chapters = append(chapters, model.StoryChapter{Number: ci, Label: "-"})
			} else {
				chapters = append(chapters, model.StoryChapter{Number: ci, Label: lm.ChapterScenario[ci].ScenarioID})
			}
		}

	case StoryLabelGreet:
		chapters = append(chapters, model.StoryChapter{Number: 0, Label: "默认"})

	case StoryLabelSpecial:
		chapters = append(chapters, model.StoryChapter{Number: 0, Label: "默认"})
	}

	return chapters
}

// GetJsonPath returns the CDN URL and filename for a story's JSON.
func (lm *ListManager) GetJsonPath(storyType, sort, index string, chapterIdx int, source string) model.JsonPathResult {
	extension := "asset"
	format := "uni"
	baseURL := lm.baseUrls["uni"]

	if source == "sekai.best" {
		format = "best"
		baseURL = lm.baseUrls["best"]
	} else if source == "haruki" {
		baseURL = lm.baseUrls["haruki"]
		extension = "json"
	} else if source == "unipjsk" {
		baseURL = lm.baseUrls["uni"]
		extension = "json"
	}

	idx := parseIndex(index)

	makeCardURL := func(charID int, cardNo, chapter string) (string, string) {
		char := model.CharacterDict[charID-1]
		cid := padZero3(charID)
		var url string
		if format == "best" {
			url = baseURL + "character/member/res" + cid + "_no" + cardNo + "/" + cid + cardNo + "_" + char.Name + chapter + "." + extension
		} else {
			url = baseURL + "startapp/character/member/res" + cid + "_no" + cardNo + "/" + cid + cardNo + "_" + char.Name + chapter + "." + extension
		}
		return url, char.Name
	}

	switch storyType {
	case StoryLabelMainStory:
		unitIdx := idx
		if unitIdx >= len(lm.MainStory) {
			return model.JsonPathResult{}
		}
		unit := lm.MainStory[unitIdx]
		ch := unit.Chapters
		if unitIdx == 0 {
			chapterIdx = (chapterIdx+1)*4/5
		}
		if chapterIdx >= len(ch) {
			return model.JsonPathResult{}
		}
		chapter := ch[chapterIdx].AssetName
		var url string
		if format == "best" {
			url = baseURL + "scenario/unitstory/" + unit.AssetName + "/" + chapter + "." + extension
		} else {
			url = baseURL + "startapp/scenario/unitstory/" + unit.AssetName + "/" + chapter + "." + extension
		}
		return model.JsonPathResult{
			URL:      url,
			FileName: "mainStory_" + chapter + ".json",
		}

	case StoryLabelEvent:
		eventIdx := len(lm.Events) - idx
		if eventIdx < 1 || eventIdx > len(lm.Events) {
			return model.JsonPathResult{}
		}
		ev := lm.Events[eventIdx-1]
		if chapterIdx >= len(ev.Chapters) {
			return model.JsonPathResult{}
		}
		chapter := ev.Chapters[chapterIdx].AssetName
		var url string
		if format == "best" {
			url = baseURL + "event_story/" + ev.Name + "/scenario/" + chapter + "." + extension
		} else {
			url = baseURL + "ondemand/event_story/" + ev.Name + "/scenario/" + chapter + "." + extension
		}
		return model.JsonPathResult{
			URL:      url,
			FileName: chapter + ".json",
		}

	case StoryLabelCardEvent:
		eventIdx := len(lm.Events) - idx
		if eventIdx < 1 || eventIdx > len(lm.Events) {
			return model.JsonPathResult{}
		}
		ev := lm.Events[eventIdx-1]
		cardSlot := chapterIdx / 3
		if cardSlot >= len(ev.Cards) {
			return model.JsonPathResult{}
		}
		cardID := ev.Cards[cardSlot]
		cardCharID := lm.Cards[cardID-1].CharacterID
		cardNo := lm.Cards[cardID-1].CardNo
		ch := padZero(chapterIdx%3 + 1)
		char := model.CharacterDict[cardCharID-1]
		charID := padZero3(cardCharID)
		var url string
		if format == "best" {
			url = baseURL + "character/member/res" + charID + "_no" + cardNo + "/" + charID + cardNo + "_" + char.Name + ch + "." + extension
		} else {
			url = baseURL + "startapp/character/member/res" + charID + "_no" + cardNo + "/" + charID + cardNo + "_" + char.Name + ch + "." + extension
		}
		return model.JsonPathResult{
			URL:      url,
			FileName: "event" + padZero3(ev.ID) + "_" + char.Name + "_" + ch + ".json",
		}

	case StoryLabelCardSpecial:
		fesIdx := len(lm.Festivals) - idx
		if fesIdx < 1 || fesIdx > len(lm.Festivals) {
			return model.JsonPathResult{}
		}
		f := lm.Festivals[fesIdx-1]
		cardSlot := chapterIdx / 3
		if cardSlot >= len(f.Cards) {
			return model.JsonPathResult{}
		}
		cardID := f.Cards[cardSlot]
		if cardID < 1 || cardID > len(lm.Cards) {
			return model.JsonPathResult{}
		}
		cardCharID := lm.Cards[cardID-1].CharacterID
		cardNo := lm.Cards[cardID-1].CardNo
		ch := padZero(chapterIdx%3 + 1)
		url, charName := makeCardURL(cardCharID, cardNo, ch)
		return model.JsonPathResult{
			URL:      url,
			FileName: "festival_" + padZero3(f.ID) + "_" + charName + "_" + ch + ".json",
		}

	case StoryLabelCardInit:
		if idx > 25 {
			// virtual singer special mapping
			idx = idx - 4
		}
		charId := (idx/5)*4 + (idx+1)%5
		charName := model.CharacterDict[charId-1].Name
		rarity := chapterIdx/3 + 1
		rarityStr := padZero3(rarity)
		ch := padZero(chapterIdx%3 + 1)
		var url string
		if format == "best" {
			url = baseURL + "character/member/res" + padZero3(charId) + "_no" + rarityStr + "/" + padZero3(charId) + rarityStr + "_" + charName + ch + "." + extension
		} else {
			url = baseURL + "startapp/character/member/res" + padZero3(charId) + "_no" + rarityStr + "/" + padZero3(charId) + rarityStr + "_" + charName + ch + "." + extension
		}
		return model.JsonPathResult{
			URL:      url,
			FileName: "release_" + charName + "_" + padZero(rarity) + "_" + ch + ".json",
		}

	case StoryLabelCardUpgrade:
		if idx > 25 {
			idx = idx - 4
		}
		charId := (idx/5)*4 + (idx+1)%5
		var levelupcards []int
		for _, f := range lm.Festivals {
			if f.LevelUp {
				levelupcards = f.Cards
				break
			}
		}
		if levelupcards == nil {
			return model.JsonPathResult{}
		}
		cardID := levelupcards[charId-1]
		// Virtual singer special VS chapters
		if charId >= 21 && chapterIdx > 2 {
			vsIdx := 30 - 6 + chapterIdx/3
			if charId == 22 { // Rin
				cardID = levelupcards[len(model.CharacterDict)-4]
			} else if charId == 23 { // Len
				cardID = levelupcards[len(model.CharacterDict)-3]
			} else if charId == 24 { // MEIKO
				cardID = levelupcards[len(model.CharacterDict)-5]
			} else if charId == 25 { // KAITO
				cardID = levelupcards[len(model.CharacterDict)-2]
			} else if charId == 26 { // Miku_band
				cardID = levelupcards[len(model.CharacterDict)-1]
			} else { // Miku (21)
				cardID = levelupcards[vsIdx]
			}
		}
		if cardID < 1 || cardID > len(lm.Cards) {
			return model.JsonPathResult{}
		}
		cardNo := lm.Cards[cardID-1].CardNo
		ch := padZero(chapterIdx%3 + 1)
		url, charName := makeCardURL(charId, cardNo, ch)
		return model.JsonPathResult{
			URL:      url,
			FileName: "levelup_" + charName + "_" + ch + ".json",
		}

	case StoryLabelAreaTalkInit, StoryLabelAreaTalkUpgrade, StoryLabelAreaTalkExtra:
		if chapterIdx >= len(lm.ChapterScenario) {
			return model.JsonPathResult{}
		}
		cs := lm.ChapterScenario[chapterIdx]
		group := cs.ID / 100
		var url string
		if format == "best" {
			url = baseURL + "scenario/actionset/group" + strconv.Itoa(group) + "/" + cs.ScenarioID + "." + extension
		} else {
			url = baseURL + "startapp/scenario/actionset/group" + strconv.Itoa(group) + "/" + cs.ScenarioID + "." + extension
		}
		fileName := "areatalk_" + cs.TalkID + "_" + cs.ScenarioID + ".json"
		return model.JsonPathResult{
			URL:      url,
			FileName: fileName,
		}

	case StoryLabelSpecial:
		specialIdx := len(lm.Specials) - 1 - idx
		if specialIdx < 0 || specialIdx >= len(lm.Specials) {
			return model.JsonPathResult{}
		}
		story := lm.Specials[specialIdx]
		var url string
		if format == "best" {
			url = baseURL + "scenario/special/" + story.DirName + "/" + story.FileName + "." + extension
		} else {
			url = baseURL + "startapp/scenario/special/" + story.DirName + "/" + story.FileName + "." + extension
		}
		return model.JsonPathResult{
			URL:      url,
			FileName: story.FileName + ".json",
		}
	}

	return model.JsonPathResult{}
}

// --- Helpers for ListManager ---

func (lm *ListManager) buildAreaTalkByTime() {
	// Simplified: just stores add/release event IDs from area talks
	for _, at := range lm.AreaTalks {
		if at.ScenarioID == "none" || at.AddEventID < 0 {
			continue
		}
		isLimited := at.Type == "limited"
		isMonthly := strings.Contains(at.ScenarioID, "monthly")
		lm.AreaTalkByTime = append(lm.AreaTalkByTime, AreaTalkTimeEntry{
			AddEventID:     at.AddEventID,
			ReleaseEventID: at.ReleaseEventID,
			Limited:        isLimited,
			Monthly:        isMonthly,
		})
	}
}

// buildAreaTalkChapterScenario populates ChapterScenario for area talk types.
func (lm *ListManager) buildAreaTalkChapterScenario(talkType, sort string) {
	lm.ChapterScenario = nil

	for _, at := range lm.AreaTalks {
		if at.ScenarioID == "none" || at.ScenarioID == "" {
			lm.ChapterScenario = append(lm.ChapterScenario, ChapterScenarioEntry{IsSeparator: true})
			continue
		}

		// Filter by type
		isInit := talkType == "init" && at.AddEventID <= 1 && at.Type != "limited"
		isUpgrade := talkType == "upgrade" && at.AddEventID > 0 && at.Type != "limited"
		isExtra := talkType == "extra" && at.Type == "limited"

		if !isInit && !isUpgrade && !isExtra {
			continue
		}

		lm.ChapterScenario = append(lm.ChapterScenario, ChapterScenarioEntry{
			ID:         at.ID,
			ScenarioID: at.ScenarioID,
			TalkID:     at.TalkID,
		})
	}
}

func parseIndex(index string) int {
	i, err := strconv.Atoi(index)
	if err != nil {
		return 0
	}
	return i
}

func padZero(n int) string {
	if n < 10 {
		return "0" + strconv.Itoa(n)
	}
	return strconv.Itoa(n)
}

func padZero3(n int) string {
	return strconv.Itoa(1000 + n)[1:]
}

// --- Voice Clue Inference ---

// BuildVoiceIDClues builds a map of voiceID prefix -> event info.
func (lm *ListManager) BuildVoiceIDClues() map[string]EventEntry {
	clues := make(map[string]EventEntry)
	for _, ev := range lm.Events {
		if iv, ok := ev.InferredVoiceIDs["prefix"]; ok {
			if prefix, ok := iv.(string); ok {
				clues[prefix] = ev
			}
		}
	}
	return clues
}

// InferVoiceEventID infers voice event IDs from area talks and stores them in events.
func (lm *ListManager) InferVoiceEventID() {
	eventsByID := make(map[int]int)
	for ei, ev := range lm.Events {
		eventsByID[ev.ID] = ei
	}

	clues := make(map[string]int) // clue prefix -> event array index
	areatalkRe := regexp.MustCompile(`areatalk_(ev|wl)_(.+)_\d+$`)

	for _, at := range lm.AreaTalks {
		match := areatalkRe.FindStringSubmatch(at.ScenarioID)
		if match == nil || at.AddEventID <= 0 {
			continue
		}
		eventClue := match[2]
		if match[1] == "wl" {
			eventClue = "wl_" + eventClue
		}
		if ei, exists := eventsByID[at.AddEventID]; exists {
			if prevEi, exists := clues[eventClue]; !exists || prevEi > ei {
				clues[eventClue] = ei
			}
		}
	}

	// Hard-coded patterns
	if ei, ok := eventsByID[1]; ok {
		clues["band_01"] = ei
	}
	if ei, ok := eventsByID[53]; ok {
		clues["night__"] = ei
	}
	if ei, ok := eventsByID[9]; ok {
		clues["shuffle_03"] = ei
	}

	for clue, ei := range clues {
		chOffset := 0
		if lm.Events[ei].ID == 9 {
			chOffset = 1
		}
		lm.Events[ei].InferredVoiceIDs = map[string]interface{}{
			"prefix":   clue,
			"choffset": chOffset,
		}
	}
}
