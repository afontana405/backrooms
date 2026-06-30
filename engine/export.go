package engine

import (
	"encoding/base64"
	"fmt"
	"net/http"
)

// DownloadExcel calls the Python export module to build an xlsx file and
// streams it to the browser as an attachment download.
//
// rows is a slice of maps — keys become column headers, values become cells.
// filename is the suggested download filename (e.g. "users.xlsx").
func DownloadExcel(w http.ResponseWriter, filename string, rows []map[string]interface{}) error {
	// Convert []map to []interface{} for the Python bridge
	irows := make([]interface{}, len(rows))
	for i, r := range rows {
		irows[i] = r
	}

	result, err := CallPython("export", "to_excel", map[string]interface{}{
		"rows":     irows,
		"filename": filename,
	})
	if err != nil {
		return fmt.Errorf("export.to_excel: %w", err)
	}

	// Unwrap the nested data object
	outer, ok := result.(map[string]interface{})
	if !ok {
		return fmt.Errorf("export: unexpected response type")
	}

	content, ok := outer["content"].(string)
	if !ok || content == "" {
		return fmt.Errorf("export: missing content in response")
	}

	raw, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return fmt.Errorf("export: base64 decode: %w", err)
	}

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(raw)))
	_, err = w.Write(raw)
	return err
}
