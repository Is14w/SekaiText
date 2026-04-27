package service

import (
	"regexp"
	"strconv"
	"strings"

	"sekaitext/backend/internal/model"
)

// FlashbackAnalyzer analyzes voice IDs to identify flashback scenes.
type FlashbackAnalyzer struct {
	lm         *ListManager
	flashbackRe *regexp.Regexp
	areatalkRe  *regexp.Regexp
	mainstoryEpRe *regexp.Regexp
	cardrarityEpRe *regexp.Regexp
	clueDict    map[string]EventEntry
	mainstory   map[string]MainStoryEntry
	voiceMsToMainStory map[string]string
}

// NewFlashbackAnalyzer creates a new FlashbackAnalyzer.
func NewFlashbackAnalyzer(lm *ListManager) *FlashbackAnalyzer {
	fb := &FlashbackAnalyzer{
		lm:          lm,
		flashbackRe: regexp.MustCompile(`voice_(.+)_\d+[a-z]?_\d+(?:_?.*)?$`),
		areatalkRe:  regexp.MustCompile(`areatalk_(ev|wl)_(.+)_\d+$`),
		mainstoryEpRe: regexp.MustCompile(`(.*?)(\d+)$`),
		cardrarityEpRe: regexp.MustCompile(`(\d+)(.*?)$`),
		mainstory:   make(map[string]MainStoryEntry),
		voiceMsToMainStory: map[string]string{
			"band":   "light_sound",
			"idol":   "idol",
			"street": "street",
			"wonder": "theme_park",
			"night":  "school_refusal",
			"piapro": "piapro",
		},
	}
	fb.updateClues()
	return fb
}

func (fb *FlashbackAnalyzer) updateClues() {
	for _, ms := range fb.lm.MainStory {
		fb.mainstory[ms.Unit] = ms
	}
	fb.clueDict = fb.lm.BuildVoiceIDClues()
}

// GetClueFromVoiceID analyzes a voice ID string and returns a clue.
// Returns: ("", false) = no idea, ("", true) = ignore, (clue, false) = scenario clue.
func (fb *FlashbackAnalyzer) GetClueFromVoiceID(voiceID string) (string, bool) {
	if strings.Contains(voiceID, "partvoice") {
		return "", true
	}
	match := fb.flashbackRe.FindStringSubmatch(voiceID)
	if match != nil {
		return match[1], false
	}
	return "", false
}

// GetClueHints returns human-readable hints for a clue.
func (fb *FlashbackAnalyzer) GetClueHints(clue, lang string) []string {
	if lang == "" {
		lang = "zh-cn"
	}
	words := strings.Split(clue, "_")
	var hints []string

	// Skip 'sc'
	firstIdx := 0
	if len(words) > 0 && words[0] == "sc" {
		firstIdx = 1
	}

	if firstIdx >= len(words) {
		return hints
	}

	first := words[firstIdx]

	// Check if it's a card label
	if first == "ev" && len(words) > 0 {
		lastWord := words[len(words)-1]
		if strings.HasSuffix(lastWord, "a") || strings.HasSuffix(lastWord, "b") {
			first = "card"
		}
	}

	switch first {
	case "ev":
		return fb.getEventHints(words[firstIdx+1:])

	case "ms", "op", "unit":
		return fb.getMainStoryHints(words[firstIdx:], first)

	case "card":
		return fb.getCardHints(words[firstIdx+1:])

	default:
		hints = append(hints, "闪回：未知来源")
	}

	return hints
}

func (fb *FlashbackAnalyzer) getEventHints(words []string) []string {
	var hints []string
	ep := -1
	if len(words) > 0 {
		if v, err := strconv.Atoi(words[len(words)-1]); err == nil {
			ep = v
			words = words[:len(words)-1]
		}
	}

	clueKey := strings.Join(words, "_")
	eventInfo, ok := fb.clueDict[clueKey]
	if !ok {
		return append(hints, "未知活动")
	}

	chOffset := 0
	if chOff, ok := eventInfo.InferredVoiceIDs["choffset"]; ok {
		if v, ok := chOff.(float64); ok {
			chOffset = int(v)
		}
	}

	if ep >= 0 {
		ep += chOffset
	}

	hints = append(hints, strconv.Itoa(eventInfo.ID)+"-"+padZero(ep))
	hints = append(hints, eventInfo.Title)
	if ep > 0 && ep <= len(eventInfo.Chapters) {
		hints = append(hints, eventInfo.Chapters[ep-1].Title)
	} else {
		hints = append(hints, "未知章节")
	}
	return hints
}

func (fb *FlashbackAnalyzer) getMainStoryHints(words []string, first string) []string {
	var hints []string
	if len(words) == 0 {
		return append(hints, "未知主线剧情")
	}

	w := words[0]
	match := fb.mainstoryEpRe.FindStringSubmatch(w)
	if match == nil {
		return append(hints, "未知主线剧情")
	}

	team := match[1]
	ep, _ := strconv.Atoi(match[2])

	if unitKey, ok := fb.voiceMsToMainStory[team]; ok {
		unitName := model.UnitDict[unitKey]
		hints = append(hints, unitName+" 主线剧情 - "+padZero(ep)+"话")

		if ms, ok := fb.mainstory[unitKey]; ok {
			if first == "unit" {
				epHints := []string{"ln", "mmj", "vbs", "ws", "25时"}
				epIdx := (ep - 1) / 4
				epSub := ((ep - 1) % 4) + 1
				if epIdx >= 0 && epIdx < len(epHints) && ep-1 >= 0 && ep-1 < len(ms.Chapters) {
					hints = append(hints, "("+epHints[epIdx]+"-"+strconv.Itoa(epSub)+") "+ms.Chapters[ep-1].Title)
				} else {
					hints = append(hints, "未知章节")
				}
			} else {
				if ep >= 0 && ep < len(ms.Chapters) {
					hints = append(hints, ms.Chapters[ep].Title)
				} else {
					hints = append(hints, "未知章节")
				}
			}
		}
	} else {
		hints = append(hints, "未知主线剧情")
	}
	return hints
}

func (fb *FlashbackAnalyzer) getCardHints(words []string) []string {
	var hints []string
	if len(words) < 2 {
		return append(hints, "未知卡面")
	}

	// Last word has star rating and episode
	starsep := words[len(words)-1]
	charIDStr := words[len(words)-2]

	match := fb.cardrarityEpRe.FindStringSubmatch(starsep)
	stars := "?"
	ep := "未知章节"
	if match != nil {
		stars = match[1]
		epPart := match[2]
		if epPart == "a" {
			ep = "前篇"
		} else if epPart == "b" {
			ep = "后篇"
		}
	}

	charID, _ := strconv.Atoi(charIDStr)
	charName := charIDStr
	if charID >= 1 && charID <= len(model.CharacterDict) {
		charName = model.CharacterDict[charID-1].NameJ
	}

	eventID := strings.Join(words[:len(words)-2], "_")

	eventHints := "卡面来自 " + eventID
	if eventID == "" {
		eventHints = "初期卡面"
	} else if ev, ok := fb.clueDict[eventID]; ok {
		eventHints = strconv.Itoa(ev.ID) + " - " + ev.Title
	}

	hints = append(hints, eventHints)
	hints = append(hints, charName+" ☆"+stars+" "+ep)
	return hints
}

// NoClue returns a placeholder for unknown event info.
func (fb *FlashbackAnalyzer) NoClue() EventEntry {
	return EventEntry{
		ID:    -1,
		Title: "未知剧情",
	}
}
