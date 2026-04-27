package service

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"sekaitext/backend/internal/model"
)

// UnityStoryData represents the Unity JSON structure for a story scenario.
type UnityStoryData struct {
	ScenarioID        string                   `json:"ScenarioId"`
	Snippets          []SnippetData            `json:"Snippets"`
	TalkData          []TalkData               `json:"TalkData"`
	SpecialEffectData []SpecialEffectData      `json:"SpecialEffectData"`
}

// SnippetData represents an action snippet in the story.
type SnippetData struct {
	Action          int `json:"Action"`
	ReferenceIndex  int `json:"ReferenceIndex"`
}

// TalkData represents a dialogue entry.
type TalkData struct {
	WindowDisplayName     string    `json:"WindowDisplayName"`
	Body                  string    `json:"Body"`
	Voices                []VoiceData `json:"Voices"`
	WhenFinishCloseWindow int       `json:"WhenFinishCloseWindow"`
}

// VoiceData represents a voice clip reference.
type VoiceData struct {
	VoiceID string  `json:"VoiceId"`
	Volume  float64 `json:"Volume"`
}

// SpecialEffectData represents a special effect (scene/option text).
type SpecialEffectData struct {
	EffectType int    `json:"EffectType"`
	StringVal  string `json:"StringVal"`
}

// JsonLoaderService parses Unity story JSON into SourceTalk entries.
type JsonLoaderService struct {
	fb *FlashbackAnalyzer
}

// NewJsonLoaderService creates a new JsonLoaderService.
func NewJsonLoaderService(fb *FlashbackAnalyzer) *JsonLoaderService {
	return &JsonLoaderService{fb: fb}
}

// ParseFile loads and parses a Unity story JSON file.
func (j *JsonLoaderService) ParseFile(path string) (*model.LoadResponse, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read story file: %w", err)
	}

	var story UnityStoryData
	if err := json.Unmarshal(data, &story); err != nil {
		return nil, fmt.Errorf("failed to parse story JSON: %w", err)
	}

	return j.parse(&story), nil
}

// ParseBytes parses Unity story JSON from raw bytes.
func (j *JsonLoaderService) ParseBytes(data []byte) (*model.LoadResponse, error) {
	var story UnityStoryData
	if err := json.Unmarshal(data, &story); err != nil {
		return nil, fmt.Errorf("failed to parse story JSON: %w", err)
	}
	return j.parse(&story), nil
}

func (j *JsonLoaderService) parse(story *UnityStoryData) *model.LoadResponse {
	var talks []model.SourceTalk
	talkIndex := make(map[int]*TalkData)
	for i := range story.TalkData {
		talkIndex[i] = &story.TalkData[i]
	}

	for _, snippet := range story.Snippets {
		switch snippet.Action {
		case 1: // TalkData
			if snippet.ReferenceIndex >= len(story.TalkData) {
				continue
			}
			talkdata := story.TalkData[snippet.ReferenceIndex]
			speaker := splitSpeaker(talkdata.WindowDisplayName)
			text := talkdata.Body

			var voices []string
			var volume []int
			for _, v := range talkdata.Voices {
				voices = append(voices, v.VoiceID)
				volume = append(volume, int(v.Volume))
			}

			charIdx := -1
			for idx, c := range model.CharacterDict {
				if c.NameJ == talkdata.WindowDisplayName {
					charIdx = idx
					break
				}
			}

			talk := model.SourceTalk{
				Speaker:  speaker,
				Text:     text,
				Voices:   voices,
				Volume:   volume,
				CharIdx:  charIdx,
			}

			talks = append(talks, talk)

			// Window close separator
			if talkdata.WhenFinishCloseWindow != 0 {
				talks = append(talks, model.SourceTalk{
					Speaker: "",
					Text:    "",
				})
			}

		case 6: // SpecialEffectData
			if snippet.ReferenceIndex >= len(story.SpecialEffectData) {
				continue
			}
			effect := story.SpecialEffectData[snippet.ReferenceIndex]

			// EffectType 8=location, 18=upper-left scene, 23=choice
			if effect.EffectType == 8 || effect.EffectType == 18 || effect.EffectType == 23 {
				speaker := "场景"
				if effect.EffectType == 18 {
					speaker = "左上场景"
				} else if effect.EffectType == 23 {
					speaker = "选项"
				}

				talks = append(talks, model.SourceTalk{
					Speaker: speaker,
					Text:    effect.StringVal,
				})

				// Separator after effect
				talks = append(talks, model.SourceTalk{
					Speaker: "",
					Text:    "",
				})
			}
		}
	}

	// Remove trailing empty
	if len(talks) > 0 && talks[len(talks)-1].Speaker == "" {
		talks = talks[:len(talks)-1]
	}

	// Flashback analysis
	talks = j.checkFlashback(talks)

	return &model.LoadResponse{
		ScenarioID:  story.ScenarioID,
		SourceTalks: talks,
	}
}

// checkFlashback analyzes voice IDs for flashback clues.
func (j *JsonLoaderService) checkFlashback(talks []model.SourceTalk) []model.SourceTalk {
	if j.fb == nil {
		return talks
	}

	for i, talk := range talks {
		if len(talk.Voices) == 0 {
			continue
		}
		clueSet := make(map[string]struct{})
		for _, voiceID := range talk.Voices {
			clue, ignore := j.fb.GetClueFromVoiceID(voiceID)
			if ignore {
				continue
			}
			if clue != "" {
				clueSet[clue] = struct{}{}
			}
		}
		for clue := range clueSet {
			talks[i].Clues = append(talks[i].Clues, clue)
		}
	}

	return talks
}

// splitSpeaker extracts the speaker name from WindowDisplayName (strip _ suffix).
func splitSpeaker(displayName string) string {
	parts := strings.SplitN(displayName, "_", 2)
	return parts[0]
}
