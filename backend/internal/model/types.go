package model

// SourceTalk represents a source text entry parsed from Unity JSON.
type SourceTalk struct {
	Speaker  string   `json:"speaker"`
	Text     string   `json:"text"`
	Voices   []string `json:"voices,omitempty"`
	Volume   []int    `json:"volume,omitempty"`
	CharIdx  int      `json:"charIndex"`
	Clues    []string `json:"clues,omitempty"`
}

// DstTalk represents a translation entry in the editor.
type DstTalk struct {
	Idx       int    `json:"idx"`
	Speaker   string `json:"speaker"`
	Text      string `json:"text"`
	Start     bool   `json:"start"`
	End       bool   `json:"end"`
	Checked   bool   `json:"checked"`
	Save      bool   `json:"save"`
	DstIdx    int    `json:"dstidx"`
	ReferID   int    `json:"referid,omitempty"`
	Proofread *bool  `json:"proofread,omitempty"`
	CheckMode bool   `json:"checkmode,omitempty"`
}

// EditorMode constants
const (
	ModeTranslate = 0
	ModeProofread = 1
	ModeCheck     = 2
)

// TalkColor represents row background color in the editor.
type TalkColor int

const (
	ColorWhite  TalkColor = iota
	ColorRed
	ColorYellow
	ColorGreen
	ColorBlue
)

// EventInfo represents an event from events.json.
type EventInfo struct {
	ID               int                    `json:"id"`
	Title            string                 `json:"title"`
	EventType        string                 `json:"eventType,omitempty"`
	Chapters         []ChapterInfo          `json:"chapters"`
	InferredVoiceIDs map[string]interface{} `json:"inferredVoiceIDs,omitempty"`
}

// ChapterInfo represents a single chapter within an event or story.
type ChapterInfo struct {
	Title string `json:"title"`
}

// MainStoryInfo represents a main story entry from mainStory.json.
type MainStoryInfo struct {
	Unit     string        `json:"unit"`
	Chapters []ChapterInfo `json:"chapters"`
}

// CardInfo represents a card entry from cards.json.
type CardInfo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

// FestivalInfo represents a festival entry.
type FestivalInfo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

// AreaTalkInfo represents an area talk entry.
type AreaTalkInfo struct {
	ID int `json:"id"`
}

// GreetInfo represents a greet entry.
type GreetInfo struct {
	ID int `json:"id"`
}

// SpecialInfo represents a special story entry.
type SpecialInfo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

// StoryType defines the type of story.
type StoryType string

const (
	StoryTypeEvent     StoryType = "event"
	StoryTypeMainStory StoryType = "mainstory"
	StoryTypeCard      StoryType = "card"
	StoryTypeFestival  StoryType = "festival"
	StoryTypeAreaTalk  StoryType = "areatalk"
	StoryTypeGreet     StoryType = "greet"
	StoryTypeSpecial   StoryType = "special"
)

// StorySort represents a sorting/filter option for stories.
type StorySort struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// StoryIndex represents a story index entry (episode/book).
type StoryIndex struct {
	Label    string `json:"label"`
	Value    string `json:"value"`
	Chapters []int  `json:"chapters,omitempty"`
}

// StoryChapter represents a chapter within a story.
type StoryChapter struct {
	Number int    `json:"number"`
	Label  string `json:"label"`
}

// JsonPathResult contains CDN URL and filename for loading a story.
type JsonPathResult struct {
	URL      string `json:"url"`
	FileName string `json:"fileName"`
}

// Settings represents application settings.
type Settings struct {
	FontSize        int    `json:"fontSize"`
	DownloadSource  string `json:"downloadSource"`
	SaveLineBreakN  bool   `json:"saveN"`
	SaveVoice       bool   `json:"saveVoice"`
	DisableSSL      bool   `json:"disableSSL"`
	VoiceOutputDir  string `json:"voiceOutputDir,omitempty"`
	LastStoryType   string `json:"lastStoryType,omitempty"`
	LastStorySort   string `json:"lastStorySort,omitempty"`
	LastStoryIndex  string `json:"lastStoryIndex,omitempty"`
	LastChapter     int    `json:"lastChapter,omitempty"`
	LastDataSource  string `json:"lastDataSource,omitempty"`
}

// DefaultSettings returns sensible defaults.
func DefaultSettings() Settings {
	return Settings{
		FontSize:       18,
		DownloadSource: "best",
		SaveLineBreakN: true,
		SaveVoice:      false,
		DisableSSL:     false,
	}
}

// UpdateProgress tracks metadata refresh progress.
type UpdateProgress struct {
	Current int    `json:"current"`
	Total   int    `json:"total"`
	Message string `json:"message,omitempty"`
	Done    bool   `json:"done"`
}

// LoadRequest is the request body for loading a story JSON.
type LoadRequest struct {
	StoryType string `json:"storyType" validate:"required"`
	Sort      string `json:"sort"`
	Index     string `json:"index"`
	Chapter   int    `json:"chapter"`
	Source    string `json:"source"`
}

// LoadResponse contains source talks after loading.
type LoadResponse struct {
	ScenarioID string       `json:"scenarioId"`
	SourceTalks []SourceTalk `json:"sourceTalks"`
}

// TranslationCreateRequest creates a new translation from source talks.
type TranslationCreateRequest struct {
	SourceTalks []SourceTalk `json:"sourceTalks" validate:"required"`
	JP          bool         `json:"jp"`
}

// TranslationLoadRequest loads a translation file.
type TranslationLoadRequest struct {
	FilePath string `json:"filePath" validate:"required"`
}

// TranslationSaveRequest saves a translation file.
type TranslationSaveRequest struct {
	FilePath string `json:"filePath" validate:"required"`
	Talks    []DstTalk `json:"talks" validate:"required"`
	SaveN    bool   `json:"saveN"`
}

// EditorChangeTextRequest edits text in a talk entry.
type EditorChangeTextRequest struct {
	Row        int      `json:"row"`
	Text       string   `json:"text"`
	EditorMode int      `json:"editorMode"`
	Talks      []DstTalk `json:"talks"`
	DstTalks   []DstTalk `json:"dstTalks"`
	ReferTalks []DstTalk `json:"referTalks"`
	SourceTalks []SourceTalk `json:"sourceTalks"`
}

// EditorAddLineRequest adds a sub-line.
type EditorAddLineRequest struct {
	Row         int        `json:"row"`
	Talks       []DstTalk  `json:"talks"`
	DstTalks    []DstTalk  `json:"dstTalks"`
	IsProofread bool       `json:"isProofreading"`
	SourceTalks []SourceTalk `json:"sourceTalks"`
}

// EditorRemoveLineRequest removes a sub-line.
type EditorRemoveLineRequest struct {
	Row      int       `json:"row"`
	Talks    []DstTalk `json:"talks"`
	DstTalks []DstTalk `json:"dstTalks"`
}

// EditorReplaceBracketsRequest replaces brackets in a talk.
type EditorReplaceBracketsRequest struct {
	Row      int       `json:"row"`
	Brackets string    `json:"brackets"`
	Talks    []DstTalk `json:"talks"`
}

// CheckLinesRequest aligns and validates loaded talks.
type CheckLinesRequest struct {
	SourceTalks []SourceTalk `json:"sourceTalks" validate:"required"`
	LoadedTalks []DstTalk    `json:"loadedTalks" validate:"required"`
}

// CheckTextRequest validates text content.
type CheckTextRequest struct {
	Speaker string `json:"speaker"`
	Text    string `json:"text"`
}

// CheckTextResponse returns validation results.
type CheckTextResponse struct {
	Text    string `json:"text"`
	Checked bool   `json:"checked"`
	Message string `json:"message,omitempty"`
}

// CompareRequest compares referTalks and checkTalks.
type CompareRequest struct {
	ReferTalks []DstTalk `json:"referTalks"`
	CheckTalks []DstTalk `json:"checkTalks"`
	EditorMode int       `json:"editorMode"`
}

// SpeakerCountRequest counts lines per speaker.
type SpeakerCountRequest struct {
	Talks []DstTalk    `json:"talks"`
	SourceTalks []SourceTalk `json:"sourceTalks"`
}

// SpeakerCountResponse contains per-speaker counts.
type SpeakerCountResponse struct {
	Speakers []SpeakerEntry `json:"speakers"`
}

// SpeakerEntry represents a single speaker's count.
type SpeakerEntry struct {
	Japanese string `json:"japanese"`
	Chinese  string `json:"chinese"`
	Count    int    `json:"count"`
}

// FlashbackAnalyzeRequest analyzes flashbacks for source talks.
type FlashbackAnalyzeRequest struct {
	SourceTalks []SourceTalk `json:"sourceTalks"`
}

// FlashbackAnalyzeResponse contains flashback analysis results.
type FlashbackAnalyzeResponse struct {
	MajorClue string       `json:"majorClue"`
	SourceTalks []SourceTalk `json:"sourceTalks"`
}

// VoiceURLRequest gets the URL for a voice file.
type VoiceURLRequest struct {
	ScenarioID string `json:"scenarioId"`
	VoiceID    string `json:"voiceId"`
	Source     string `json:"source"`
}

// VoiceURLResponse contains the voice playback URL.
type VoiceURLResponse struct {
	URL string `json:"url"`
}

// CharacterInfo from characterDict.
type CharacterInfo struct {
	Index int    `json:"index"`
	Name  string `json:"name"`
	NameJ string `json:"nameJ"`
	NameC string `json:"nameC"`
}

// UnitInfo from unitDict.
type UnitInfo struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}
