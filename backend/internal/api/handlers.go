package api

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"sekaitext/backend/internal/config"
	"sekaitext/backend/internal/model"
	"sekaitext/backend/internal/service"
)

// Handler holds shared dependencies for all HTTP handlers.
type Handler struct {
	cfg          *config.AppConfig
	lm           *service.ListManager
	editor       *service.EditorService
	jsonLoader   *service.JsonLoaderService
	fb           *service.FlashbackAnalyzer
	dl           *service.Downloader
	progress     *service.ProgressTracker
	logBuf       *service.LogBuffer
	downloadTasks sync.Map // map[string]*model.DownloadTaskProgress
}

// NewHandler creates a new Handler with all services initialized.
func NewHandler(cfg *config.AppConfig, logBuf *service.LogBuffer) *Handler {
	lm := service.NewListManager(cfg.CatalogDir)
	fb := service.NewFlashbackAnalyzer(lm)
	return &Handler{
		cfg:        cfg,
		lm:         lm,
		editor:     service.NewEditorService(),
		jsonLoader: service.NewJsonLoaderService(fb),
		fb:         fb,
		dl:         service.NewDownloader(cfg.DataDir),
		progress:   service.NewProgressTracker(),
		logBuf:     logBuf,
	}
}

// --- Story ---

func (h *Handler) StoryTypes(w http.ResponseWriter, r *http.Request) {
	types := h.lm.GetStoryTypes()
	writeJSON(w, http.StatusOK, types)
}

func (h *Handler) StorySorts(w http.ResponseWriter, r *http.Request) {
	storyType := r.URL.Query().Get("type")
	sorts := h.lm.GetStorySorts(storyType)
	if sorts == nil {
		sorts = []model.StorySort{}
	}
	writeJSON(w, http.StatusOK, sorts)
}

func (h *Handler) StoryIndex(w http.ResponseWriter, r *http.Request) {
	storyType := r.URL.Query().Get("type")
	sort := r.URL.Query().Get("sort")
	indices := h.lm.GetStoryIndexList(storyType, sort)
	if indices == nil {
		indices = []model.StoryIndex{}
	}
	writeJSON(w, http.StatusOK, indices)
}

func (h *Handler) StoryChapter(w http.ResponseWriter, r *http.Request) {
	storyType := r.URL.Query().Get("type")
	sort := r.URL.Query().Get("sort")
	index := r.URL.Query().Get("index")
	chapters := h.lm.GetStoryChapterList(storyType, sort, index)
	if chapters == nil {
		chapters = []model.StoryChapter{}
	}
	writeJSON(w, http.StatusOK, chapters)
}

func (h *Handler) JsonPath(w http.ResponseWriter, r *http.Request) {
	storyType := r.URL.Query().Get("type")
	sort := r.URL.Query().Get("sort")
	index := r.URL.Query().Get("index")
	chapter, _ := strconv.Atoi(r.URL.Query().Get("chapter"))
	source := r.URL.Query().Get("source")

	result := h.lm.GetJsonPath(storyType, sort, index, chapter, source)
	writeJSON(w, http.StatusOK, result)
}

func (h *Handler) StoryLoad(w http.ResponseWriter, r *http.Request) {
	var req model.LoadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Get JSON path from CDN
	path := h.lm.GetJsonPath(req.StoryType, req.Sort, req.Index, req.Chapter, req.Source)
	if path.URL == "" {
		writeError(w, http.StatusNotFound, "story not found")
		return
	}

	// Download and parse
	filePath, err := h.dl.DownloadJSON(path.URL, path.FileName)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to download story: "+err.Error())
		return
	}
	resp, err := h.jsonLoader.ParseFile(filePath)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to parse story: "+err.Error())
		return
	}

	resp.SaveTitle = path.SaveTitle
	resp.ChapterTitle = path.ChapterTitle
	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) StoryParseLocal(w http.ResponseWriter, r *http.Request) {
	var req struct {
		FilePath string `json:"filePath"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.jsonLoader.ParseFile(req.FilePath)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to parse file: "+err.Error())
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) StoryLoadLocal(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.jsonLoader.ParseBytes([]byte(req.Content))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to parse story: "+err.Error())
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

// --- Translation ---

func (h *Handler) TranslationCreate(w http.ResponseWriter, r *http.Request) {
	var req model.TranslationCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	talks := h.editor.CreateFile(req.SourceTalks, req.JP)
	writeJSON(w, http.StatusOK, talks)
}

func (h *Handler) TranslationLoad(w http.ResponseWriter, r *http.Request) {
	var req model.TranslationLoadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	talks, meta, err := h.editor.LoadFile(req.FilePath)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"talks": talks,
		"meta":  meta,
	})
}

func (h *Handler) TranslationLoadContent(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	talks, meta, err := h.editor.LoadContent(req.Content)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"talks": talks,
		"meta":  meta,
	})
}

func (h *Handler) TranslationSave(w http.ResponseWriter, r *http.Request) {
	var req model.TranslationSaveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	content := h.editor.SerializeWithMeta(req.Talks, req.SaveN, req.Meta)
	if err := os.WriteFile(req.FilePath, []byte(content), 0644); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "saved"})
}

func (h *Handler) TranslationSerialize(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Talks []model.DstTalk    `json:"talks"`
		SaveN bool               `json:"saveN"`
		Meta  *model.SaveMetadata `json:"meta,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	content := h.editor.SerializeWithMeta(req.Talks, req.SaveN, req.Meta)
	writeJSON(w, http.StatusOK, map[string]string{"content": content})
}

func (h *Handler) CheckLines(w http.ResponseWriter, r *http.Request) {
	var req model.CheckLinesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	talks := h.editor.CheckLines(req.SourceTalks, req.LoadedTalks)
	writeJSON(w, http.StatusOK, talks)
}

// --- Editor ---

func (h *Handler) ChangeText(w http.ResponseWriter, r *http.Request) {
	var req model.EditorChangeTextRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	talks, dstTalks := h.editor.ChangeText(req.Row, req.Text, req.EditorMode,
		req.Talks, req.DstTalks, req.ReferTalks, req.SourceTalks)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"talks":    talks,
		"dstTalks": dstTalks,
	})
}

func (h *Handler) AddLine(w http.ResponseWriter, r *http.Request) {
	var req model.EditorAddLineRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	talks, dstTalks := h.editor.AddLine(req.Row, req.Talks, req.DstTalks, req.IsProofread)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"talks":    talks,
		"dstTalks": dstTalks,
	})
}

func (h *Handler) RemoveLine(w http.ResponseWriter, r *http.Request) {
	var req model.EditorRemoveLineRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	talks, dstTalks := h.editor.RemoveLine(req.Row, req.Talks, req.DstTalks)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"talks":    talks,
		"dstTalks": dstTalks,
	})
}

func (h *Handler) Compare(w http.ResponseWriter, r *http.Request) {
	var req model.CompareRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	talks := h.editor.CompareText(req.ReferTalks, req.CheckTalks, req.EditorMode)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"talks":    talks,
		"dstTalks": talks,
	})
}

func (h *Handler) ReplaceBrackets(w http.ResponseWriter, r *http.Request) {
	var req model.EditorReplaceBracketsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	talks := h.editor.ReplaceBrackets(req.Talks, req.Row, req.Brackets)
	writeJSON(w, http.StatusOK, talks)
}

// --- Check Text ---

func (h *Handler) CheckText(w http.ResponseWriter, r *http.Request) {
	var req model.CheckTextRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp := h.editor.GetTextCheck(req)
	writeJSON(w, http.StatusOK, resp)
}

// --- Flashback ---

func (h *Handler) FlashbackAnalyze(w http.ResponseWriter, r *http.Request) {
	var req model.FlashbackAnalyzeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Analyze clues for each talk
	for i := range req.SourceTalks {
		if len(req.SourceTalks[i].Voices) == 0 {
			continue
		}
		clueSet := make(map[string]struct{})
		for _, voiceID := range req.SourceTalks[i].Voices {
			clue, ignore := h.fb.GetClueFromVoiceID(voiceID)
			if !ignore && clue != "" {
				clueSet[clue] = struct{}{}
			}
		}
		for clue := range clueSet {
			req.SourceTalks[i].Clues = append(req.SourceTalks[i].Clues, clue)
		}
	}

	writeJSON(w, http.StatusOK, model.FlashbackAnalyzeResponse{
		SourceTalks: req.SourceTalks,
	})
}

func (h *Handler) ClueHints(w http.ResponseWriter, r *http.Request) {
	clue := r.URL.Query().Get("clue")
	lang := r.URL.Query().Get("lang")

	hints := h.fb.GetClueHints(clue, lang)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"clue":  clue,
		"hints": hints,
	})
}

func (h *Handler) VoiceClues(w http.ResponseWriter, r *http.Request) {
	clues := h.lm.BuildVoiceIDClues()
	writeJSON(w, http.StatusOK, clues)
}

// --- Voice ---

func (h *Handler) VoiceURL(w http.ResponseWriter, r *http.Request) {
	scenarioID := r.URL.Query().Get("scenarioId")
	voiceID := r.URL.Query().Get("voiceId")
	source := r.URL.Query().Get("source")

	if source == "" {
		source = "sekai.best"
	}

	baseURL := "https://storage.sekai.best/sekai-jp-assets/"
	if source == "unipjsk" {
		baseURL = "https://assets.unipjsk.com/"
	}

	url := baseURL + "voice/" + scenarioID + "/" + voiceID + ".mp3"
	writeJSON(w, http.StatusOK, model.VoiceURLResponse{URL: url})
}

// --- Speaker ---

func (h *Handler) SpeakerCount(w http.ResponseWriter, r *http.Request) {
	var req model.SpeakerCountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Count speakers
	speakerMap := make(map[string]struct {
		japanese string
		count    int
	})

	for _, talk := range req.Talks {
		if talk.Speaker == "" || talk.Speaker == "场景" || talk.Speaker == "左上场景" || talk.Speaker == "选项" {
			continue
		}
		var srcSpeaker string
		if talk.Idx > 0 && talk.Idx-1 < len(req.SourceTalks) {
			srcSpeaker = req.SourceTalks[talk.Idx-1].Speaker
		}
		if srcSpeaker == "" {
			srcSpeaker = talk.Speaker
		}
		entry := speakerMap[srcSpeaker]
		entry.japanese = srcSpeaker
		entry.count++
		speakerMap[srcSpeaker] = entry
	}

	var speakers []model.SpeakerEntry
	for _, entry := range speakerMap {
		speakers = append(speakers, model.SpeakerEntry{
			Japanese: entry.japanese,
			Chinese:  "",
			Count:    entry.count,
		})
	}

	writeJSON(w, http.StatusOK, model.SpeakerCountResponse{Speakers: speakers})
}

// --- Settings ---

func (h *Handler) settingsPath() string {
	return h.cfg.CatalogDir + "/settings.json"
}

func (h *Handler) loadSettings() (model.Settings, error) {
	data, err := os.ReadFile(h.settingsPath())
	if err != nil {
		return model.Settings{}, err
	}
	var s model.Settings
	if err := json.Unmarshal(data, &s); err != nil {
		return model.Settings{}, err
	}
	return s, nil
}

func (h *Handler) saveSettings(s model.Settings) error {
	os.MkdirAll(h.cfg.CatalogDir, 0755)
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(h.settingsPath(), data, 0644)
}

func (h *Handler) GetSettings(w http.ResponseWriter, r *http.Request) {
	s, err := h.loadSettings()
	if err != nil {
		s = model.DefaultSettings()
	}
	writeJSON(w, http.StatusOK, s)
}

func (h *Handler) PutSettings(w http.ResponseWriter, r *http.Request) {
	var s model.Settings
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.saveSettings(s); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to save settings")
		return
	}
	writeJSON(w, http.StatusOK, s)
}

// --- Update ---

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	go h.lm.UpdateAll(h.cfg.CatalogDir, h.progress)
	writeJSON(w, http.StatusOK, map[string]string{"status": "started"})
}

func (h *Handler) UpdateProgress(w http.ResponseWriter, r *http.Request) {
	current, total, message, done := h.progress.Status()
	writeJSON(w, http.StatusOK, model.UpdateProgress{
		Current: current,
		Total:   total,
		Message: message,
		Done:    done,
	})
}

// --- JSON Download ---

func (h *Handler) DownloadStoryJSON(w http.ResponseWriter, r *http.Request) {
	var req model.JsonDownloadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	path := h.lm.GetJsonPath(req.StoryType, req.Sort, req.Index, req.Chapter, req.Source)
	if path.URL == "" {
		writeError(w, http.StatusNotFound, "story not found")
		return
	}

	outputDir := req.OutputDir
	if outputDir == "" {
		outputDir = h.cfg.DataDir + "/json"
	}

	taskID := strconv.FormatInt(time.Now().UnixNano(), 36)
	task := &model.DownloadTaskProgress{TaskID: taskID, Status: "downloading"}
	h.downloadTasks.Store(taskID, task)

	go func() {
		dlPath, err := h.dl.DownloadJSONToDir(path.URL, outputDir, path.FileName, func(read, total int64) {
			task.Read = read
			task.Total = total
		})
		if err != nil {
			task.Status = "error"
			task.Error = err.Error()
		} else {
			task.Status = "done"
			task.FilePath = dlPath
		}
	}()

	writeJSON(w, http.StatusOK, map[string]string{"taskId": taskID})
}

func (h *Handler) DownloadProgress(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Query().Get("task")
	if taskID == "" {
		writeError(w, http.StatusBadRequest, "missing task id")
		return
	}

	val, ok := h.downloadTasks.Load(taskID)
	if !ok {
		writeError(w, http.StatusNotFound, "task not found")
		return
	}

	task := val.(*model.DownloadTaskProgress)
	writeJSON(w, http.StatusOK, task)

	// Clean up completed tasks after serving
	if task.Status == "done" || task.Status == "error" {
		h.downloadTasks.Delete(taskID)
	}
}

// --- Assets ---

func (h *Handler) Characters(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.CharacterDict)
}

func (h *Handler) CharacterIcon(w http.ResponseWriter, r *http.Request) {
	indexStr := chi.URLParam(r, "index")
	index, err := strconv.Atoi(indexStr)
	if err != nil || index < 1 || index > 31 {
		writeError(w, http.StatusBadRequest, "invalid character index")
		return
	}
	iconPath := h.cfg.ImagesChrDir + "/chr_" + indexStr + ".png"
	w.Header().Set("Content-Type", "image/png")
	http.ServeFile(w, r, iconPath)
}

func (h *Handler) Units(w http.ResponseWriter, r *http.Request) {
	units := make([]model.UnitInfo, 0, len(model.UnitDict))
	for k, v := range model.UnitDict {
		units = append(units, model.UnitInfo{Key: k, Name: v})
	}
	writeJSON(w, http.StatusOK, units)
}

func (h *Handler) Areas(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.AreaDict)
}

// --- Debug ---

func (h *Handler) DebugLogs(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.logBuf.Lines())
}

func (h *Handler) DebugSaveLogs(w http.ResponseWriter, r *http.Request) {
	entries := h.logBuf.Lines()
	f, err := os.Create("debug.log")
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create log file: "+err.Error())
		return
	}
	defer f.Close()

	f.WriteString("=== SekaiText Debug Log === " + time.Now().Format("2006-01-02 15:04:05") + " ===\n\n")
	for _, e := range entries {
		f.WriteString("[" + e.Timestamp + "] " + e.Message + "\n")
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status": "saved",
		"lines":  len(entries),
	})
}

// --- Recovery (autosave) ---

func (h *Handler) recoveryPath() string {
	return h.cfg.DataDir + "/autosave.json"
}

func (h *Handler) RecoverySave(w http.ResponseWriter, r *http.Request) {
	var req model.RecoverySaveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	content := h.editor.SerializeContent(req.Talks, req.SaveN)
	data := model.RecoveryData{
		Content:    content,
		FilePath:   req.FilePath,
		EditorMode: req.EditorMode,
		SavedAt:    time.Now().Format("2006-01-02 15:04:05"),
	}

	f, err := os.Create(h.recoveryPath())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to save autosave: "+err.Error())
		return
	}
	defer f.Close()
	json.NewEncoder(f).Encode(data)
	writeJSON(w, http.StatusOK, map[string]string{"status": "saved"})
}

func (h *Handler) RecoveryLoad(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open(h.recoveryPath())
	if err != nil {
		writeJSON(w, http.StatusOK, map[string]interface{}{"exists": false})
		return
	}
	defer f.Close()

	var data model.RecoveryData
	if err := json.NewDecoder(f).Decode(&data); err != nil {
		writeJSON(w, http.StatusOK, map[string]interface{}{"exists": false})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"exists":     true,
		"content":    data.Content,
		"filePath":   data.FilePath,
		"editorMode": data.EditorMode,
		"savedAt":    data.SavedAt,
	})
}

func (h *Handler) RecoveryClear(w http.ResponseWriter, r *http.Request) {
	os.Remove(h.recoveryPath())
	writeJSON(w, http.StatusOK, map[string]string{"status": "cleared"})
}

// --- Helpers ---

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
