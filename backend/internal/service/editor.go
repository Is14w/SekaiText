package service

import (
	"encoding/json"
	"fmt"
	"math"

	diff "github.com/sergi/go-diff/diffmatchpatch"
	"os"
	"strings"

	"sekaitext/backend/internal/model"
)

// EditorService implements the core editor logic (port of Editor.py).
type EditorService struct{}

// NewEditorService creates a new EditorService.
func NewEditorService() *EditorService {
	return &EditorService{}
}

// CreateFile creates dstTalks from source talks (translation template).
func (e *EditorService) CreateFile(srctalks []model.SourceTalk, jp bool) []model.DstTalk {
	var dsttalks []model.DstTalk

	for idx, srctalk := range srctalks {
		subsrctalks := strings.Split(srctalk.Text, "\n")
		for iidx, subsrctalk := range subsrctalks {
			text := ""
			if jp {
				text = subsrctalk
			} else {
				for _, char := range subsrctalk {
					if strings.ContainsRune("♪☆/『』", char) {
						text += string(char)
					}
				}
			}
			dsttalks = append(dsttalks, model.DstTalk{
				Idx:     idx + 1,
				Speaker: srctalk.Speaker,
				Text:    text,
				Start:   iidx == 0,
				End:     false,
				Checked: true,
				Save:    true,
			})
		}
		dsttalks[len(dsttalks)-1].End = true
	}

	// Replace Japanese speaker names with Chinese
	for i := range dsttalks {
		speaker := strings.ReplaceAll(dsttalks[i].Speaker, "の声", "")
		parts := strings.Split(speaker, "・")
		for pi, part := range parts {
			if char, ok := model.FindCharacterByJapaneseName(part); ok {
				parts[pi] = char.NameC
			}
		}
		newSpeaker := strings.Join(parts, "・")
		newSpeaker = strings.ReplaceAll(newSpeaker, "の声", "的声音")
		newSpeaker = strings.ReplaceAll(newSpeaker, "ネネロボ", "宁宁号")
		dsttalks[i].Speaker = newSpeaker
	}

	return dsttalks
}

// LoadFile parses a translation .txt file into DstTalk entries.
func (e *EditorService) LoadFile(filepath string) ([]model.DstTalk, *model.SaveMetadata, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read file: %w", err)
	}
	return e.LoadContent(string(data))
}

// SaveFile serializes dstTalks to translation .txt format.
func (e *EditorService) SaveFile(filepath string, dsttalks []model.DstTalk, saveN bool) error {
	var out strings.Builder
	for _, talk := range dsttalks {
		switch talk.Speaker {
		case "场景", "左上场景", "":
			if (talk.Speaker == "场景" || talk.Speaker == "左上场景") && talk.Text == "" {
				talk.Text = talk.Speaker
			} else if talk.Speaker == "选项" && !strings.Contains(talk.Text, "/") {
				talk.Text = talk.Text + "/"
			}
			if talk.Speaker == "" && talk.Text != "" {
				out.WriteString(talk.Text)
			} else {
				out.WriteString(talk.Text)
			}
			out.WriteString("\n")

		default:
			if talk.Start {
				out.WriteString(talk.Speaker + "：")
			}
			lines := strings.Split(talk.Text, "\n")
			out.WriteString(lines[0])
			if !talk.End {
				if saveN {
					out.WriteString("\\N")
				}
			} else {
				out.WriteString("\n")
			}
		}
	}

	result := strings.TrimRight(out.String(), "\n")
	return os.WriteFile(filepath, []byte(result), 0644)
}

// LoadContent parses translation content from a string instead of a file.
// Returns the parsed talks and any embedded SaveMetadata (nil if absent).
func (e *EditorService) LoadContent(content string) ([]model.DstTalk, *model.SaveMetadata, error) {
	var meta *model.SaveMetadata
	if strings.HasPrefix(content, "#SekaiText ") {
		idx := strings.Index(content, "\n")
		if idx > 0 {
			header := content[len("#SekaiText "):idx]
			var m model.SaveMetadata
			if json.Unmarshal([]byte(header), &m) == nil {
				meta = &m
			}
			content = content[idx+1:]
		}
	}

	lines := strings.Split(content, "\n")
	var talks []model.DstTalk
	preblank := false

	for idx, line := range lines {
		line = strings.Replace(line, ":", "：", 1)
		var speaker, fulltext string

		if strings.Contains(line, "：") {
			parts := strings.SplitN(line, "：", 2)
			speaker = parts[0]
			fulltext = parts[1]
		} else if strings.Contains(line, "/") {
			speaker = "选项"
			fulltext = line
		} else {
			trimmed := strings.TrimSpace(line)
			if trimmed == "" {
				speaker = ""
				fulltext = line
			} else {
				speaker = "场景"
				fulltext = line
			}
		}

		if speaker == "" {
			if preblank {
				continue
			}
			preblank = true
		} else {
			preblank = false
		}

		texts := strings.Split(fulltext, "\\N")
		for iidx, text := range texts {
			text, checked, msg := e.checkText(speaker, text)
			talk := model.DstTalk{
				Idx:     idx + 1,
				Speaker: speaker,
				Text:    text,
				Start:   iidx == 0,
				End:     false,
				Checked: checked,
				Save:    true,
				Message: msg,
			}
			talks = append(talks, talk)
		}
		talks[len(talks)-1].End = true
	}

	if preblank && len(talks) > 0 {
		talks = talks[:len(talks)-1]
	}

	return talks, meta, nil
}

// SerializeWithMeta prepends a metadata header to the serialized content.
func (e *EditorService) SerializeWithMeta(dsttalks []model.DstTalk, saveN bool, meta *model.SaveMetadata) string {
	content := e.SerializeContent(dsttalks, saveN)
	if meta == nil {
		return content
	}
	b, err := json.Marshal(meta)
	if err != nil {
		return content
	}
	return "#SekaiText " + string(b) + "\n" + content
}

// SerializeContent serializes dstTalks to translation text format and returns the string.
func (e *EditorService) SerializeContent(dsttalks []model.DstTalk, saveN bool) string {
	var out strings.Builder
	for _, talk := range dsttalks {
		switch talk.Speaker {
		case "场景", "左上场景", "":
			if (talk.Speaker == "场景" || talk.Speaker == "左上场景") && talk.Text == "" {
				talk.Text = talk.Speaker
			} else if talk.Speaker == "选项" && !strings.Contains(talk.Text, "/") {
				talk.Text = talk.Text + "/"
			}
			if talk.Speaker == "" && talk.Text != "" {
				out.WriteString(talk.Text)
			} else {
				out.WriteString(talk.Text)
			}
			out.WriteString("\n")

		default:
			if talk.Start {
				out.WriteString(talk.Speaker + "：")
			}
			lines := strings.Split(talk.Text, "\n")
			out.WriteString(lines[0])
			if !talk.End {
				if saveN {
					out.WriteString("\\N")
				}
			} else {
				out.WriteString("\n")
			}
		}
	}

	return strings.TrimRight(out.String(), "\n")
}

// CheckLines aligns loaded talks with source talks (handles line mismatches).
func (e *EditorService) CheckLines(srctalks []model.SourceTalk, loadtalks []model.DstTalk) []model.DstTalk {
	// Trim leading scene/option lines that exceed source
	srcCount := 0
	for _, st := range srctalks {
		if st.Speaker == "场景" || st.Speaker == "左上场景" || st.Speaker == "选项" {
			srcCount++
		} else if st.Speaker != "" {
			break
		}
	}
	count := 0
	for _, lt := range loadtalks {
		if lt.Speaker == "场景" || lt.Speaker == "左上场景" || lt.Speaker == "选项" {
			count++
		} else if lt.Speaker != "" {
			break
		}
	}
	if count > srcCount {
		loadtalks = loadtalks[count-srcCount:]
	}
	for len(loadtalks) > 0 && loadtalks[0].Text == "" {
		loadtalks = loadtalks[1:]
	}

	var newtalks []model.DstTalk
	idx := 0
	dstend := false

	for srcidx, srctalk := range srctalks {
		if idx >= len(loadtalks) {
			dstend = true
		}

		// Fix "左上场景" vs "场景"
		if !dstend && srctalk.Speaker == "左上场景" && loadtalks[idx].Speaker == "场景" {
			loadtalks[idx].Speaker = "左上场景"
		}

		// Scene/option lines
		if srctalk.Speaker == "场景" || srctalk.Speaker == "左上场景" || srctalk.Speaker == "选项" || srctalk.Speaker == "" {
			if dstend || srctalk.Speaker != loadtalks[idx].Speaker {
				newtalks = append(newtalks, model.DstTalk{
					Idx:     srcidx + 1,
					Speaker: srctalk.Speaker,
					Text:    srctalk.Text,
					Start:   true,
					End:     true,
					Checked: true,
					Save:    true,
				})
				continue
			}
		}

		subsrctalks := strings.Split(srctalk.Text, "\n")
		dstidx := -1
		if !dstend {
			dstidx = loadtalks[idx].Idx
		}

		for iidx := range subsrctalks {
			if idx >= len(loadtalks) {
				dstend = true
			}

			if !dstend && loadtalks[idx].Idx == dstidx {
				text, _, msg := e.checkText(srctalk.Speaker, loadtalks[idx].Text)
				talk := model.DstTalk{
					Idx:     srcidx + 1,
					Speaker: loadtalks[idx].Speaker,
					Text:    text,
					Start:   iidx == 0,
					End:     false,
					Checked: true,
					Save:    true,
					Message: msg,
				}
				idx++
				newtalks = append(newtalks, talk)
			} else if dstend {
				talk := model.DstTalk{
					Idx:     srcidx + 1,
					Speaker: srctalk.Speaker,
					Text:    " ",
					Start:   iidx == 0,
					End:     false,
					Checked: false,
					Save:    true,
				}
				newtalks = append(newtalks, talk)
			} else {
				// \N lost
				if len(newtalks) > 0 {
					newtalks[len(newtalks)-1].Text += "\n【分行不一致】"
					newtalks[len(newtalks)-1].End = true
					newtalks[len(newtalks)-1].Checked = false
				}
				continue
			}
		}

		// Too many \N
		for idx < len(loadtalks) && loadtalks[idx].Idx == dstidx {
			text, _, chkMsg := e.checkText(srctalk.Speaker, loadtalks[idx].Text)
			lineMsg := "分行不一致"
			if chkMsg != "" {
				lineMsg = chkMsg + "\n" + lineMsg
			}
			talk := model.DstTalk{
				Idx:     srcidx + 1,
				Speaker: loadtalks[idx].Speaker,
				Text:    text,
				Start:   false,
				End:     true,
				Checked: false,
				Save:    true,
				Message: lineMsg,
			}
			idx++
			newtalks = append(newtalks, talk)
		}

		if len(newtalks) > 0 {
			newtalks[len(newtalks)-1].End = true
		}
	}

	// Extra lines at the end
	if idx < len(loadtalks) {
		idxdiff := 0
		if len(newtalks) > 0 && idx < len(loadtalks) {
			idxdiff = newtalks[len(newtalks)-1].Idx - loadtalks[idx].Idx + 1
		}
		for _, talk := range loadtalks[idx:] {
			newtalk := talk
			newtalk.Idx = talk.Idx + idxdiff
			newtalk.Checked = false
			if newtalk.Message != "" {
				newtalk.Message += "\n多余行"
			} else {
				newtalk.Message = "多余行"
			}
			newtalks = append(newtalks, newtalk)
		}
	}

	return newtalks
}

// CompareText compares referTalks with checkTalks sub-line by sub-line,
// grouped by source idx. Each sub-line pair within an idx is compared individually.
func (e *EditorService) CompareText(refertalks, checktalks []model.DstTalk, editormode int) []model.DstTalk {
	refGrp := groupByIdx(refertalks)
	chkGrp := groupByIdx(checktalks)

	// Collect all unique idxs in order
	var keys []int
	seen := make(map[int]bool)
	for _, t := range refertalks {
		if !seen[t.Idx] {
			seen[t.Idx] = true
			keys = append(keys, t.Idx)
		}
	}
	for _, t := range checktalks {
		if !seen[t.Idx] {
			seen[t.Idx] = true
			keys = append(keys, t.Idx)
		}
	}
	for i := 0; i < len(keys); i++ {
		for j := i + 1; j < len(keys); j++ {
			if keys[j] < keys[i] {
				keys[i], keys[j] = keys[j], keys[i]
			}
		}
	}

	var newtalks []model.DstTalk
	dstBase := 0

	for _, idx := range keys {
		ref := refGrp[idx]
		chk := chkGrp[idx]
		maxN := len(ref)
		if len(chk) > maxN {
			maxN = len(chk)
		}

		if len(ref) == 0 {
			for k, t := range chk {
				t.Proofread = boolPtr(true)
				t.DstIdx = dstBase + k
				newtalks = append(newtalks, t)
			}
		} else if len(chk) == 0 {
			for k, t := range ref {
				t.Checked = false
				t.DstIdx = dstBase + k
				if editormode == 2 {
					t.CheckMode = true
				} else {
					t.Save = false
					t.Proofread = boolPtr(false)
				}
				newtalks = append(newtalks, t)
			}
		} else {
			// Compare sub-line by sub-line
			for k := 0; k < maxN; k++ {
				if k < len(ref) && k < len(chk) {
					r := ref[k]
					c := chk[k]
					if r.Text == c.Text {
						r.DstIdx = dstBase
						newtalks = append(newtalks, r)
						dstBase++
					} else {
						r.Checked = false
						r.DstIdx = dstBase
						if editormode == 2 {
							r.CheckMode = true
							referDiff, checkDiff := computeDiff(r.Text, c.Text)
							r.DiffParts = referDiff
							c.DiffParts = checkDiff
						} else {
							r.Save = false
							r.Proofread = boolPtr(false)
						}
						newtalks = append(newtalks, r)
						dstBase++

						c.Proofread = boolPtr(true)
						c.DstIdx = dstBase
						newtalks = append(newtalks, c)
						dstBase++
					}
				} else if k < len(ref) {
					r := ref[k]
					r.Checked = false
					r.DstIdx = dstBase
					if editormode == 2 {
						r.CheckMode = true
					} else {
						r.Save = false
						r.Proofread = boolPtr(false)
					}
					newtalks = append(newtalks, r)
					dstBase++
				} else {
					c := chk[k]
					c.Proofread = boolPtr(true)
					c.DstIdx = dstBase
					newtalks = append(newtalks, c)
					dstBase++
				}
			}
		}
	}

	return newtalks
}

// groupByIdx collects rows with the same source idx into slices.
func groupByIdx(talks []model.DstTalk) map[int][]model.DstTalk {
	groups := make(map[int][]model.DstTalk)
	for _, t := range talks {
		groups[t.Idx] = append(groups[t.Idx], t)
	}
	return groups
}

// ChangeText handles text editing with proofread mode logic.
func (e *EditorService) ChangeText(row int, text string, editormode int,
	talks, dsttalks, refertalks []model.DstTalk) ([]model.DstTalk, []model.DstTalk) {

	if row >= len(talks) {
		return talks, dsttalks
	}

	speaker := talks[row].Speaker
	text, checked, msg := e.checkText(speaker, text)

	if len(strings.Split(text, "\n")) > 1 {
		checked = false
		msg = "文本含换行符，请拆分至多行"
	}

	if speaker == "" {
		return talks, dsttalks
	}

	// Translate mode
	if editormode == 0 {
		talks[row].Text = text
		talks[row].Checked = checked
		talks[row].Message = msg
		dstidx := talks[row].DstIdx
		if dstidx < len(dsttalks) {
			dsttalks[dstidx].Text = text
			dsttalks[dstidx].Checked = checked
			dsttalks[dstidx].Message = msg
		}
		return talks, dsttalks
	}

	// Proofread mode — create new proofread line when editing original
	if editormode == 1 {
		if talks[row].Proofread == nil || !*talks[row].Proofread {
			newtalk := talks[row]
			newtalk.Text = text
			newtalk.Checked = true
			newtalk.Save = true
			newtalk.Proofread = boolPtr(true)
			newtalk.Message = msg

			dstidx := talks[row].DstIdx
			if dstidx < len(dsttalks) {
				dsttalks[dstidx].Text = text
				dsttalks[dstidx].Checked = checked
				dsttalks[dstidx].Message = msg
			}
			newtalk.DstIdx = dstidx

			talks = insertTalk(talks, row+1, newtalk)
			talks[row].Checked = false
			talks[row].Save = false
			talks[row].Proofread = boolPtr(false)
			talks[row].Message = ""

			return talks, dsttalks
		}

		// Update existing proofread line
		talks[row].Text = text
		talks[row].Checked = true
		talks[row].Message = msg
		dstidx := talks[row].DstIdx
		if dstidx < len(dsttalks) {
			dsttalks[dstidx].Text = text
			dsttalks[dstidx].Checked = checked
			dsttalks[dstidx].Message = msg
		}
	}

	// Check/合意 mode — always update in-place, recompute diffs
	if editormode == 2 {
		talks[row].Text = text
		talks[row].Checked = true
		talks[row].Message = msg
		dstidx := talks[row].DstIdx
		if dstidx < len(dsttalks) {
			dsttalks[dstidx].Text = text
			dsttalks[dstidx].Checked = checked
			dsttalks[dstidx].Message = msg
		}
		recomputeDiffs(talks, talks[row].Idx)
	}

	return talks, dsttalks
}

// recomputeDiffs recalculates diffs for sub-line pairs at a given idx after editing.
// Sub-lines are paired positionally: checkmode₀ with proofread₀, checkmode₁ with proofread₁, etc.
func recomputeDiffs(talks []model.DstTalk, rowIdx int) {
	// Collect rows at this idx
	var rows []int
	for i := range talks {
		if talks[i].Idx == rowIdx {
			rows = append(rows, i)
		}
	}
	// Pair consecutive checkmode+proofread rows by position
	for i := 0; i < len(rows)-1; {
		a := rows[i]
		b := rows[i+1]
		ra := talks[a]
		rb := talks[b]
		if ra.CheckMode && rb.Proofread != nil && *rb.Proofread {
			if ra.Text == rb.Text {
				ra.DiffParts = nil
				rb.DiffParts = nil
			} else {
				referDiff, checkDiff := computeDiff(ra.Text, rb.Text)
				talks[a].DiffParts = referDiff
				talks[b].DiffParts = checkDiff
			}
			i += 2
		} else {
			i++
		}
	}
}

// AddLine adds a sub-line after the given row.
func (e *EditorService) AddLine(row int, talks, dsttalks []model.DstTalk, isProofreading bool) ([]model.DstTalk, []model.DstTalk) {
	if row >= len(talks) {
		return talks, dsttalks
	}

	newtalk := talks[row]
	newtalk.Text = " "
	newtalk.End = true
	newtalk.Checked = true
	newtalk.Save = true
	newtalk.Start = false
	if isProofreading {
		newtalk.Proofread = boolPtr(true)
	}

	dstidx := talks[row].DstIdx
	dsttalks = insertDstTalk(dsttalks, dstidx+1, newtalk)

	talks = insertTalk(talks, row+1, newtalk)
	for i := row + 1; i < len(talks); i++ {
		talks[i].DstIdx++
	}

	// Update the original row
	talks[row].End = false
	talks[row].Checked = true
	talks[row].Save = true
	if dstidx < len(dsttalks) {
		dsttalks[dstidx].End = false
		dsttalks[dstidx].Checked = true
		dsttalks[dstidx].Save = true
	}

	return talks, dsttalks
}

// RemoveLine removes a sub-line at the given row.
func (e *EditorService) RemoveLine(row int, talks, dsttalks []model.DstTalk) ([]model.DstTalk, []model.DstTalk) {
	if row >= len(talks) || row < 0 {
		return talks, dsttalks
	}

	dstidx := talks[row].DstIdx

	// Mark previous row as end
	if row > 0 {
		talks[row-1].End = true
		if dstidx-1 >= 0 && dstidx-1 < len(dsttalks) {
			dsttalks[dstidx-1].End = true
		}
	}

	talks = append(talks[:row], talks[row+1:]...)
	if dstidx >= 0 && dstidx < len(dsttalks) {
		dsttalks = append(dsttalks[:dstidx], dsttalks[dstidx+1:]...)
	}

	for i := row; i < len(talks); i++ {
		talks[i].DstIdx--
	}

	return talks, dsttalks
}

// ReplaceBrackets replaces all bracket types with the specified pair.
func (e *EditorService) ReplaceBrackets(talks []model.DstTalk, row int, brackets string) []model.DstTalk {
	if row >= len(talks) || len(brackets) < 2 {
		return talks
	}
	openB := brackets[0]
	closeB := brackets[1]

	var newText strings.Builder
	for _, ch := range talks[row].Text {
		if strings.ContainsRune("「『（“‘【", ch) {
			newText.WriteRune(rune(openB))
		} else if strings.ContainsRune("」』（”’】", ch) {
			newText.WriteRune(rune(closeB))
		} else {
			newText.WriteRune(ch)
		}
	}
	talks[row].Text = newText.String()
	return talks
}

// ShowDiff updates proofread row visibility based on showDiff state.
func (e *EditorService) ShowDiff(talks []model.DstTalk) {
	// Frontend handles visibility; this is a no-op in the service
}

// CheckProofread toggles the checked state on a row.
func (e *EditorService) CheckProofread(talks []model.DstTalk, row int, checked bool) {
	if row < len(talks) {
		talks[row].Checked = checked
	}
}

// checkText validates text content and returns fixed text, pass/fail, and warning message.
func (e *EditorService) checkText(speaker, text string) (string, bool, string) {
	var msg string

	if speaker != "" && speaker != "场景" && speaker != "左上场景" && speaker != "选项" && text == "" {
		msg = "空行，若不需要改行请点右侧“-”删去本行"
		return text, false, msg
	}

	text = strings.TrimSpace(strings.Split(text, "\n")[0])
	if text == "" {
		return text, true, ""
	}

	if speaker == "场景" || speaker == "左上场景" || speaker == "" {
		return text, true, ""
	}

	if speaker == "选项" {
		if !strings.Contains(text, "/") {
			text += "/"
		}
		if strings.HasSuffix(text, "/") {
			msg = "选项必须用/分隔"
			return text, false, msg
		}
		return text, true, ""
	}

	// Standard text replacements
	text = strings.ReplaceAll(text, "…", "...")
	text = strings.ReplaceAll(text, "(", "（")
	text = strings.ReplaceAll(text, ")", "）")
	text = strings.ReplaceAll(text, ",", "，")
	text = strings.ReplaceAll(text, "?", "？")
	text = strings.ReplaceAll(text, "!", "！")
	text = strings.ReplaceAll(text, "~", "～")
	text = strings.ReplaceAll(text, "欸", "诶")

	check := true
	normalEnd := []rune{'、', '，', '。', '？', '！', '～', '♪', '☆', '.', '—'}
	unusualEnd := []rune{'）', '」', '』', '”'}

	runes := []rune(text)
	if len(runes) > 0 {
		last := runes[len(runes)-1]
		if containsRune(normalEnd, last) {
			if strings.Contains(text, ".，") || strings.Contains(text, ".。") {
				msg = "「……。」和「……，」只保留省略号"
				check = false
			}
		} else if containsRune(unusualEnd, last) {
			if len(runes) > 1 && !containsRune(normalEnd, runes[len(runes)-2]) {
				msg = "句尾缺少逗号句号"
				check = false
			}
		} else {
			msg = "句尾缺少逗号句号"
			check = false
		}
	}

	// Check dashes (only if no msg yet, to avoid overwriting)
	if msg == "" && strings.Contains(text, "—") {
		dashCount := strings.Count(text, "—")
		doubleDashCount := strings.Count(text, "——")
		if dashCount != doubleDashCount*2 {
			msg = "破折号用双破折——，或者视情况删掉"
			check = false
		}
	}

	// Check line length (only if no msg yet)
	if msg == "" {
		lineLen := lineLength(strings.Split(text, "\n")[0])
		if lineLen >= 30 {
			msg = "单行过长，请删减或换行"
			check = false
		}
	}

	return text, check, msg
}

// UpdateHiddenRowMap builds compression/decompression maps for hidden rows.
func (e *EditorService) UpdateHiddenRowMap(talks []model.DstTalk) (compressMap, decompressMap []int) {
	current := 0
	for idx, talk := range talks {
		compressMap = append(compressMap, current)
		if talk.Proofread != nil && !*talk.Proofread {
			continue
		}
		decompressMap = append(decompressMap, idx)
		current++
	}
	return
}

// lineLength calculates display width (half-width for ASCII).
func lineLength(s string) int {
	count := 0
	for _, ch := range s {
		if ch <= 127 {
			count++
		} else {
			count += 2
		}
	}
	return int(math.Ceil(float64(count) / 2.0))
}

// computeDiff uses diffmatchpatch for character-level diff between two strings.
func computeDiff(a, b string) ([]model.DiffPart, []model.DiffPart) {
	dmp := diff.New()
	diffs := dmp.DiffMain(a, b, false)
	diffs = dmp.DiffCleanupSemantic(diffs)

	var referParts, checkParts []model.DiffPart
	for _, d := range diffs {
		switch d.Type {
		case diff.DiffEqual:
			referParts = append(referParts, model.DiffPart{Text: d.Text, Type: "same"})
			checkParts = append(checkParts, model.DiffPart{Text: d.Text, Type: "same"})
		case diff.DiffDelete:
			referParts = append(referParts, model.DiffPart{Text: d.Text, Type: "remove"})
		case diff.DiffInsert:
			checkParts = append(checkParts, model.DiffPart{Text: d.Text, Type: "add"})
		}
	}
	return referParts, checkParts
}

func containsRune(runes []rune, r rune) bool {
	for _, v := range runes {
		if v == r {
			return true
		}
	}
	return false
}

func boolPtr(b bool) *bool {
	return &b
}

func insertTalk(slice []model.DstTalk, index int, talk model.DstTalk) []model.DstTalk {
	slice = append(slice, model.DstTalk{})
	copy(slice[index+1:], slice[index:])
	slice[index] = talk
	return slice
}

func insertDstTalk(slice []model.DstTalk, index int, talk model.DstTalk) []model.DstTalk {
	slice = append(slice, model.DstTalk{})
	copy(slice[index+1:], slice[index:])
	slice[index] = talk
	return slice
}

// GetTextCheck performs text validation and returns the result.
func (e *EditorService) GetTextCheck(req model.CheckTextRequest) model.CheckTextResponse {
	text, checked, msg := e.checkText(req.Speaker, req.Text)
	resp := model.CheckTextResponse{
		Text:    text,
		Checked: checked,
		Message: msg,
	}
	return resp
}
