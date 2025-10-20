package handlers

import (
	"bufio"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode"
	"unicode/utf8"
)

type CompressResult struct {
	Logs       []LogEntry `json:"logs"`
	OutputPath string     `json:"outputPath"`
	KeyCount   int        `json:"keyCount"`
}

// CompressCSVToJSON converts CSV to JSON file
func CompressCSVToJSON(
	ctx context.Context,
	csvPath string,
	outputPath string, // đường dẫn JSON output (nếu rỗng: cùng thư mục với CSV)
	includeEmpty bool, // có ghi cả bản dịch rỗng không
	pretty bool, // JSON thụt dòng
	perKeyLog bool, // log [OK] cho từng key
) (CompressResult, error) {
	var result CompressResult
	var logs []LogEntry
	push := func(level, msg string) {
		logs = append(logs, LogEntry{Level: level, Message: msg})
	}

	// Mở CSV & detect delimiter
	f, err := os.Open(csvPath)
	if err != nil {
		return result, fmt.Errorf("open csv: %w", err)
	}
	defer f.Close()

	br := bufio.NewReader(f)
	delim := detectDelimiter(br)

	r := csv.NewReader(br)
	r.Comma = delim
	r.FieldsPerRecord = -1
	r.TrimLeadingSpace = true
	r.ReuseRecord = false

	// Đọc header
	header, err := r.Read()
	if err == io.EOF {
		push("error", "CSV file is empty")
		result.Logs = logs
		return result, nil
	}
	if err != nil {
		return result, fmt.Errorf("read header: %w", err)
	}
	if len(header) < 2 {
		push("error", `Invalid header: must have at least 2 columns "Key","Translation"`)
		result.Logs = logs
		return result, nil
	}

	h0 := strings.TrimSpace(stripBOM(header[0]))
	h1 := strings.TrimSpace(header[1])
	if !strings.EqualFold(h0, "Key") || h1 == "" {
		push("error", fmt.Sprintf(`Header must be "Key","Translation" (received: "%s","%s")`, h0, h1))
	} else {
		push("ok", `Valid header detected: "Key","Translation"`)
	}

	// Duyệt dòng & gom key/val
	out := make(map[string]interface{}) // Changed from map[string]string to support nesting
	keys := make(map[string]int)        // để cảnh báo trùng
	row := 1                            // header = dòng 1
	errCount, warnCount := 0, 0

	for {
		select {
		case <-ctx.Done():
			push("warn", "Compression cancelled by context")
			goto CONVERT
		default:
		}

		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		row++
		if err != nil {
			errCount++
			push("error", fmt.Sprintf("Line %d: parse error: %v", row, err))
			continue
		}
		if len(rec) < 2 {
			errCount++
			push("error", fmt.Sprintf("Line %d: less than 2 columns", row))
			continue
		}

		key := strings.TrimSpace(replaceNBSP(rec[0]))
		val := rec[1] // giữ nguyên để bảo toàn placeholder và xuống dòng

		if key == "" {
			errCount++
			push("error", fmt.Sprintf("Line %d: Key is empty", row))
			continue
		}
		if !utf8.ValidString(key) || !utf8.ValidString(val) {
			errCount++
			push("error", fmt.Sprintf("Line %d: invalid UTF-8", row))
			continue
		}
		if prev, dup := keys[key]; dup {
			warnCount++
			push("warn", fmt.Sprintf("Line %d: Duplicate key with line %d. Value will be overwritten", row, prev))
		}
		keys[key] = row

		if !includeEmpty && strings.TrimSpace(val) == "" {
			push("info", fmt.Sprintf("Line %d: Empty translation → skipped (includeEmpty=false)", row))
			continue
		}

		// Split only at FIRST dot to create category
		// Example: "CyborgNames.A.N.D.Y." -> ["CyborgNames", "A.N.D.Y."]
		idx := strings.Index(key, ".")
		if idx > 0 && idx < len(key)-1 {
			// Has a dot and it's not at the end
			category := key[:idx]
			subKey := key[idx+1:]

			// Create nested structure
			if _, exists := out[category]; !exists {
				out[category] = make(map[string]interface{})
			}
			if categoryMap, ok := out[category].(map[string]interface{}); ok {
				categoryMap[subKey] = val
			}
		} else {
			// No dot or dot at end - keep flat
			out[key] = val
		}

		if perKeyLog {
			push("ok", fmt.Sprintf("%s", key))
		}
	}

CONVERT:
	// Xác định output path - use temp file if empty
	if outputPath == "" {
		dir := filepath.Dir(csvPath)
		base := filepath.Base(csvPath)
		name := strings.TrimSuffix(base, filepath.Ext(base))
		outputPath = filepath.Join(dir, name+"_temp.json")
	}

	// Ghi JSON
	var data []byte
	if pretty {
		data, err = json.MarshalIndent(out, "", "  ")
	} else {
		data, err = json.Marshal(out)
	}
	if err != nil {
		return result, fmt.Errorf("marshal json: %w", err)
	}

	if err := os.WriteFile(outputPath, data, 0o644); err != nil {
		return result, fmt.Errorf("write json: %w", err)
	}

	push("ok", fmt.Sprintf("JSON prepared: %d keys processed", len(out)))
	if errCount > 0 || warnCount > 0 {
		push("summary", fmt.Sprintf("Processing complete: %d keys, %d errors, %d warnings", len(out), errCount, warnCount))
	} else {
		push("summary", fmt.Sprintf("All keys processed successfully: %d total", len(out)))
	}

	result.Logs = logs
	result.OutputPath = outputPath
	result.KeyCount = len(out)
	return result, nil
}

// Helper functions
func replaceNBSP(s string) string {
	return strings.ReplaceAll(s, "\u00A0", " ")
}

func containsUnicodeSpace(s string) bool {
	for _, r := range s {
		if unicode.IsSpace(r) && r != '_' {
			return true
		}
	}
	return false
}
