package engine

import (
	"fmt"
	"strings"
)

// ResolveSchemaProps checks if a component has a "schema" prop.
// If so, loads data from binary and injects it directly as native props.
// Returns true if schema was resolved (data is in props["rows"]).
func ResolveSchemaProps(props map[string]interface{}, diag *DiagCollector) bool {
	schemaName, ok := props["schema"].(string)
	if !ok || schemaName == "" {
		return false
	}

	schema := GetBinarySchema(schemaName)
	if schema == nil {
		if diag != nil {
			diag.Error("schema", fmt.Sprintf("Binary schema '%s' not found", schemaName))
		}
		return false
	}

	// Load all records
	docs, err := schema.BinaryFindAll()
	if err != nil {
		if diag != nil {
			diag.Error("schema", fmt.Sprintf("BinaryFindAll failed for '%s': %v", schemaName, err))
		}
		return false
	}

	// Apply filter if present
	if filter, ok := props["filter"].(map[string]interface{}); ok && len(filter) > 0 {
		docs = filterDocs(docs, filter)
	}

	// Apply sort if present
	if sortField, ok := props["sort"].(string); ok && sortField != "" {
		sortDocs(docs, sortField)
	}

	// Apply limit if present
	if limit, ok := props["limit"].(float64); ok && int(limit) > 0 && int(limit) < len(docs) {
		docs = docs[:int(limit)]
	}

	// Inject as native props — no JSON round-trip
	props["rows"] = docs
	props["count"] = len(docs)

	if diag != nil {
		diag.InfoDetail("schema",
			fmt.Sprintf("Schema '%s' → %d records", schemaName, len(docs)),
			fmt.Sprintf("fields=%d filter=%v", len(schema.Fields), props["filter"]))
	}

	// Clean up schema-specific props so components don't see them as HTML attrs
	delete(props, "schema")
	delete(props, "filter")
	delete(props, "sort")
	delete(props, "limit")

	return true
}

// filterDocs returns only docs where all filter fields match.
func filterDocs(docs []map[string]interface{}, filter map[string]interface{}) []map[string]interface{} {
	var result []map[string]interface{}
	for _, doc := range docs {
		match := true
		for k, v := range filter {
			expected := fmt.Sprintf("%v", v)
			actual := fmt.Sprintf("%v", doc[k])
			if !strings.EqualFold(actual, expected) {
				match = false
				break
			}
		}
		if match {
			result = append(result, doc)
		}
	}
	return result
}

// sortDocs sorts docs by a field (ascending string compare).
func sortDocs(docs []map[string]interface{}, field string) {
	for i := 1; i < len(docs); i++ {
		for j := i; j > 0; j-- {
			a := fmt.Sprintf("%v", docs[j-1][field])
			b := fmt.Sprintf("%v", docs[j][field])
			if a > b {
				docs[j-1], docs[j] = docs[j], docs[j-1]
			} else {
				break
			}
		}
	}
}
