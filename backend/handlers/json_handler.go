package handlers

import (
	"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"
)

type LogEntry struct {
	Level   string `json:"level"`
	Message string `json:"message"`
}

func PreviewCSVFile(ctx context.Context, filePath string) ([]LogEntry, error) {
	const maxPreviewLines = 20000
	var logs []LogEntry

	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	br := bufio.NewReader(f)
	delim := detectDelimiter(br)

	r := csv.NewReader(br)
	r.Comma = delim
	r.FieldsPerRecord = -1
	r.TrimLeadingSpace = true
	r.ReuseRecord = false

	// Read and process lines
	header, err := r.Read()
	if err == io.EOF {
		logs = append(logs, LogEntry{Level: "error", Message: "CSV file is empty"})
		return logs, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	if len(header) == 0 {
		logs = append(logs, LogEntry{Level: "error", Message: `Invalid CSV file: must be at least two column, :"Key","Name of langue`})
		return logs, nil
	}

	h0 := strings.TrimSpace(stripBOM(header[0]))
	h1 := strings.TrimSpace(header[1])

	if !strings.EqualFold((h0), "Key") || h1 == "" {
		logs = append(logs, LogEntry{Level: "error", Message: `Invalid Header, must be :"Key","Name of language - example: "Key","en-US"`})
	} else {
		logs = append(logs, LogEntry{Level: "ok", Message: `Valid CSV header detected`})
	}

	if len(header) < 2 {
		logs = append(logs, LogEntry{Level: "warn", Message: `CSV file has less than two columns, no translations to process`})
	}

	// Read preview lines
	keys := make(map[string]int)
	errCount, warnCount := 0, 0
	lineNo := 1 // Including header

	for ; lineNo <= maxPreviewLines; lineNo++ {
		select {
		case <-ctx.Done():
			logs = append(logs, LogEntry{Level: "warn", Message: "Preview generation cancelled"})
			goto SUM
		default:
		}

		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			errCount++
			logs = append(logs, LogEntry{Level: "error", Message: fmt.Sprintf("Error parse line %d: %v", lineNo+1, err)})
			continue
		}

		if len(rec) < 2 {
			errCount++
			logs = append(logs, LogEntry{Level: "error", Message: fmt.Sprintf("Line %d: less than two columns.", lineNo+1)})
			continue
		}

		key := strings.TrimSpace(rec[0])
		val := rec[1] // Keep {0}, \n v.v

		if key == "" {
			errCount++
			logs = append(logs, LogEntry{Level: "error", Message: fmt.Sprintf("Line %d: Key empty.", lineNo+1)})
			continue
		}
		if !utf8.ValidString(key) || !utf8.ValidString(val) {
			errCount++
			logs = append(logs, LogEntry{Level: "error", Message: fmt.Sprintf("Line %d: invalid UTF-8.", lineNo+1)})
		}
		if strings.ContainsAny(key, " \t") {
			warnCount++
			logs = append(logs, LogEntry{Level: "warn", Message: fmt.Sprintf("Line %d: Key has white space (should be Some.Key).", lineNo+1)})
		}
		if prevLine, dup := keys[key]; dup {
			warnCount++
			logs = append(logs, LogEntry{Level: "warn", Message: fmt.Sprintf("Line %d: Key same with line %d. Value gonna be override after convert.", lineNo+1, prevLine)})
		}
		keys[key] = lineNo + 1

		if strings.TrimSpace(val) == "" {
			warnCount++
			logs = append(logs, LogEntry{Level: "info", Message: fmt.Sprintf("Line %d: Translation Emptry.", lineNo+1)})
		}
		if unmatchedBraces(val) {
			warnCount++
			logs = append(logs, LogEntry{Level: "warn", Message: fmt.Sprintf("Line %d: '{}' in Translation not invalid.", lineNo+1)})
		}

	}

SUM:
	summary := fmt.Sprintf(
		"Done checking file ~%d Line (delimiter '%c'). Error: %d, Warning: %d",
		len(keys), delim, errCount, warnCount,
	)
	logs = append(logs, LogEntry{Level: "summary", Message: summary})
	return logs, nil
}

func stripBOM(s string) string {
	return strings.TrimPrefix(s, "\uFEFF")
}
func detectDelimiter(br *bufio.Reader) rune {
	peek, _ := br.Peek(4096)
	counts := map[rune]int{',': 0, '\t': 0, ';': 0}
	for _, b := range peek {
		if b == '\n' || b == '\r' {
			break
		}
		if b == ',' || b == ';' || b == '\t' {
			counts[rune(b)]++
		}
	}
	delim := ','
	max := 0
	for r, c := range counts {
		if c > max {
			max = c
			delim = r
		}
	}
	return delim
}

func unmatchedBraces(s string) bool {
	balance := 0
	for _, r := range s {
		if r == '{' {
			balance++
		} else if r == '}' {
			balance--
			if balance < 0 {
				return true
			}
		}
	}
	return balance != 0
}
