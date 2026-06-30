package engine

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ── Binary Schema Definition ────────────────────────────────────────────────

type BinarySchema struct {
	Collection string        `json:"collection"`
	BinaryCol  string        `json:"binaryCollection"`
	Fields     []BinaryField `json:"fields"`
	recordSize  int
	fieldIndex  map[string]int
	indexFields []string
	lookupMu    sync.RWMutex
	lookupTables map[string][]string
	reverseLookup map[string]map[string]uint32
}

type BinaryField struct {
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	Offset   int      `json:"offset"`
	Size     int      `json:"size"`
	Values   []string `json:"values"`
	Dynamic  bool     `json:"dynamic"`
	Index    bool     `json:"index"`
}

// ── Schema Registry ─────────────────────────────────────────────────────────

var (
	binarySchemas   = map[string]*BinarySchema{}
	binarySchemaMu  sync.RWMutex
)

// LoadBinarySchemas loads all binary schema definitions from a directory.
func LoadBinarySchemas(dir string) error {
	dirPath := strings.ReplaceAll(dir, "\\", "/")
	entries, err := ReadEmbedDir(dirPath)
	if err != nil {
		return nil
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		if entry.Name()[0] == '_' {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		path = strings.ReplaceAll(path, "\\", "/")
		raw, err := ReadEmbedFile(path)
		if err != nil {
			return fmt.Errorf("reading %s: %w", path, err)
		}

		var schema BinarySchema
		if err := json.Unmarshal(raw, &schema); err != nil {
			return fmt.Errorf("parsing %s: %w", path, err)
		}

		schema.compile()

		if SqlDB != nil {
			schema.ensureTable()
			schema.LoadLookupTables()
		}

		binarySchemaMu.Lock()
		binarySchemas[schema.Collection] = &schema
		binarySchemaMu.Unlock()

		indexNote := ""
		if len(schema.indexFields) > 0 {
			indexNote = fmt.Sprintf(" indexed:[%s]", strings.Join(schema.indexFields, ","))
		}
		fmt.Printf("[binary] Loaded schema: %s → %s (%d fields, %d bytes/record%s)\n",
			schema.Collection, schema.BinaryCol, len(schema.Fields), schema.recordSize, indexNote)
	}

	return nil
}

func GetBinarySchema(collection string) *BinarySchema {
	binarySchemaMu.RLock()
	defer binarySchemaMu.RUnlock()
	return binarySchemas[collection]
}

func (s *BinarySchema) compile() {
	s.fieldIndex = make(map[string]int)
	s.indexFields = nil
	s.lookupTables = make(map[string][]string)
	s.reverseLookup = make(map[string]map[string]uint32)

	offset := 0
	for i := range s.Fields {
		f := &s.Fields[i]
		s.fieldIndex[f.Name] = i
		if f.Index {
			s.indexFields = append(s.indexFields, f.Name)
		}
		f.Offset = offset

		switch f.Type {
		case "uint8", "lookup":
			f.Size = 1
		case "uint16", "lookup16":
			f.Size = 2
		case "uint32":
			f.Size = 4
		case "uint64", "timestamp":
			f.Size = 8
		default:
			f.Size = 4
		}

		if (f.Type == "lookup" || f.Type == "lookup16") && len(f.Values) > 0 && !f.Dynamic {
			s.lookupTables[f.Name] = f.Values
			rev := make(map[string]uint32)
			for j, v := range f.Values {
				rev[v] = uint32(j)
			}
			s.reverseLookup[f.Name] = rev
		}

		offset += f.Size
	}
	s.recordSize = offset
}

// ── SQLite Table Management ─────────────────────────────────────────────────

// ensureTable creates the SQLite table for this binary collection.
func (s *BinarySchema) ensureTable() {
	cols := "id INTEGER PRIMARY KEY AUTOINCREMENT, d BLOB"
	for _, fieldName := range s.indexFields {
		idx := s.fieldIndex[fieldName]
		f := s.Fields[idx]
		sqlType := "TEXT"
		switch f.Type {
		case "uint8", "uint16", "uint32", "uint64", "timestamp":
			sqlType = "INTEGER"
		}
		cols += fmt.Sprintf(", %s %s", fieldName, sqlType)
	}
	SqlDB.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", s.BinaryCol, cols))

	for _, fieldName := range s.indexFields {
		SqlDB.Exec(fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_%s_%s ON %s (%s)",
			s.BinaryCol, fieldName, s.BinaryCol, fieldName))
	}
}

// ── Encode ──────────────────────────────────────────────────────────────────

func (s *BinarySchema) Encode(doc map[string]interface{}) ([]byte, error) {
	buf := make([]byte, s.recordSize)

	for _, f := range s.Fields {
		val, exists := doc[f.Name]
		if !exists {
			continue
		}

		switch f.Type {
		case "uint8":
			buf[f.Offset] = toByte(val)
		case "uint16":
			binary.LittleEndian.PutUint16(buf[f.Offset:f.Offset+2], toUint16(val))
		case "uint32":
			binary.LittleEndian.PutUint32(buf[f.Offset:f.Offset+4], toUint32(val))
		case "uint64":
			binary.LittleEndian.PutUint64(buf[f.Offset:f.Offset+8], toUint64(val))
		case "timestamp":
			t := toTimestamp(val)
			binary.LittleEndian.PutUint64(buf[f.Offset:f.Offset+8], uint64(t))
		case "lookup":
			str := fmt.Sprintf("%v", val)
			id := s.lookupID(f.Name, str, f.Dynamic)
			buf[f.Offset] = byte(id)
		case "lookup16":
			str := fmt.Sprintf("%v", val)
			id := s.lookupID(f.Name, str, f.Dynamic)
			binary.LittleEndian.PutUint16(buf[f.Offset:f.Offset+2], uint16(id))
		}
	}

	return buf, nil
}

// ── Decode ──────────────────────────────────────────────────────────────────

func (s *BinarySchema) Decode(buf []byte) map[string]interface{} {
	doc := make(map[string]interface{}, len(s.Fields))

	for _, f := range s.Fields {
		if f.Offset+f.Size > len(buf) {
			continue
		}

		switch f.Type {
		case "uint8":
			doc[f.Name] = int(buf[f.Offset])
		case "uint16":
			doc[f.Name] = int(binary.LittleEndian.Uint16(buf[f.Offset : f.Offset+2]))
		case "uint32":
			doc[f.Name] = int(binary.LittleEndian.Uint32(buf[f.Offset : f.Offset+4]))
		case "uint64":
			doc[f.Name] = binary.LittleEndian.Uint64(buf[f.Offset : f.Offset+8])
		case "timestamp":
			ts := binary.LittleEndian.Uint64(buf[f.Offset : f.Offset+8])
			doc[f.Name] = time.Unix(int64(ts), 0).Format(time.RFC3339)
		case "lookup":
			id := int(buf[f.Offset])
			doc[f.Name] = s.lookupValue(f.Name, id)
		case "lookup16":
			id := int(binary.LittleEndian.Uint16(buf[f.Offset : f.Offset+2]))
			doc[f.Name] = s.lookupValue(f.Name, id)
		}
	}

	return doc
}

// ── Lookup Table Management ─────────────────────────────────────────────────

func (s *BinarySchema) lookupID(field, value string, dynamic bool) uint32 {
	s.lookupMu.RLock()
	if rev, ok := s.reverseLookup[field]; ok {
		if id, found := rev[value]; found {
			s.lookupMu.RUnlock()
			return id
		}
	}
	s.lookupMu.RUnlock()

	if !dynamic {
		return 0
	}

	s.lookupMu.Lock()
	defer s.lookupMu.Unlock()

	if rev, ok := s.reverseLookup[field]; ok {
		if id, found := rev[value]; found {
			return id
		}
	}

	if s.lookupTables[field] == nil {
		s.lookupTables[field] = []string{}
		s.reverseLookup[field] = make(map[string]uint32)
	}

	id := uint32(len(s.lookupTables[field]))
	s.lookupTables[field] = append(s.lookupTables[field], value)
	s.reverseLookup[field][value] = id
	return id
}

func (s *BinarySchema) lookupValue(field string, id int) string {
	s.lookupMu.RLock()
	defer s.lookupMu.RUnlock()

	table := s.lookupTables[field]
	if id >= 0 && id < len(table) {
		return table[id]
	}
	return fmt.Sprintf("?%d", id)
}

// ── Collection Operations ───────────────────────────────────────────────────

func (s *BinarySchema) BinaryInsert(doc map[string]interface{}) error {
	buf, err := s.Encode(doc)
	if err != nil {
		return err
	}

	cols := "d"
	placeholders := "?"
	args := []interface{}{buf}
	for _, fieldName := range s.indexFields {
		if val, ok := doc[fieldName]; ok {
			cols += ", " + fieldName
			placeholders += ", ?"
			args = append(args, s.coerceIndexValue(fieldName, val))
		}
	}

	_, err = SqlDB.Exec(
		fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", s.BinaryCol, cols, placeholders),
		args...,
	)
	return err
}

func (s *BinarySchema) BinaryInsertMany(docs []map[string]interface{}) error {
	for _, doc := range docs {
		if err := s.BinaryInsert(doc); err != nil {
			return err
		}
	}
	return nil
}

func (s *BinarySchema) BinaryFindAll() ([]map[string]interface{}, error) {
	rows, err := SqlDB.Query(fmt.Sprintf("SELECT id, d FROM %s", s.BinaryCol))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var id int64
		var d []byte
		if err := rows.Scan(&id, &d); err != nil {
			continue
		}
		doc := s.Decode(d)
		doc["_id"] = id
		results = append(results, doc)
	}
	return results, nil
}

func (s *BinarySchema) BinaryFind(filter map[string]interface{}) ([]map[string]interface{}, error) {
	where, args := s.buildWhere(filter)
	query := fmt.Sprintf("SELECT id, d FROM %s", s.BinaryCol)
	if where != "" {
		query += " WHERE " + where
	}

	rows, err := SqlDB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var id int64
		var d []byte
		if err := rows.Scan(&id, &d); err != nil {
			continue
		}
		doc := s.Decode(d)
		doc["_id"] = id
		results = append(results, doc)
	}
	return results, nil
}

// QueryOpts configures pagination and sorting for BinaryFindPage.
type QueryOpts struct {
	Page     int
	PageSize int
	SortBy   string
	SortDir  int // 1 = ASC, -1 = DESC
}

func (s *BinarySchema) BinaryFindPage(filter map[string]interface{}, opts QueryOpts) ([]map[string]interface{}, int64, error) {
	where, args := s.buildWhere(filter)

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s", s.BinaryCol)
	if where != "" {
		countQuery += " WHERE " + where
	}
	var total int64
	SqlDB.QueryRow(countQuery, args...).Scan(&total)

	pageSize := opts.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	page := opts.Page
	if page < 0 {
		page = 0
	}
	sortField := opts.SortBy
	if sortField == "" {
		sortField = "id"
	}
	sortDir := "ASC"
	if opts.SortDir == -1 {
		sortDir = "DESC"
	}

	query := fmt.Sprintf("SELECT id, d FROM %s", s.BinaryCol)
	if where != "" {
		query += " WHERE " + where
	}
	query += fmt.Sprintf(" ORDER BY %s %s LIMIT ? OFFSET ?", sortField, sortDir)
	args = append(args, pageSize, page*pageSize)

	rows, err := SqlDB.Query(query, args...)
	if err != nil {
		return nil, total, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var id int64
		var d []byte
		if err := rows.Scan(&id, &d); err != nil {
			continue
		}
		doc := s.Decode(d)
		doc["_id"] = id
		results = append(results, doc)
	}
	return results, total, nil
}

func (s *BinarySchema) BinaryUpdate(id interface{}, doc map[string]interface{}) error {
	buf, err := s.Encode(doc)
	if err != nil {
		return err
	}
	setClauses := "d = ?"
	args := []interface{}{buf}
	for _, fieldName := range s.indexFields {
		if val, ok := doc[fieldName]; ok {
			setClauses += fmt.Sprintf(", %s = ?", fieldName)
			args = append(args, s.coerceIndexValue(fieldName, val))
		}
	}
	args = append(args, id)
	_, err = SqlDB.Exec(
		fmt.Sprintf("UPDATE %s SET %s WHERE id = ?", s.BinaryCol, setClauses),
		args...,
	)
	return err
}

func (s *BinarySchema) BinaryDelete(id interface{}) error {
	_, err := SqlDB.Exec(fmt.Sprintf("DELETE FROM %s WHERE id = ?", s.BinaryCol), id)
	return err
}

// SaveLookupTables persists dynamic lookup tables to SQLite.
func (s *BinarySchema) SaveLookupTables() error {
	s.lookupMu.RLock()
	defer s.lookupMu.RUnlock()

	SqlDB.Exec(`CREATE TABLE IF NOT EXISTS _binary_lookups (
		collection TEXT PRIMARY KEY,
		tables TEXT
	)`)

	tablesJSON, _ := json.Marshal(s.lookupTables)
	SqlDB.Exec(`INSERT OR REPLACE INTO _binary_lookups (collection, tables) VALUES (?, ?)`,
		s.Collection, string(tablesJSON))
	return nil
}

// LoadLookupTables restores dynamic lookup tables from SQLite.
func (s *BinarySchema) LoadLookupTables() error {
	var tablesJSON string
	err := SqlDB.QueryRow(`SELECT tables FROM _binary_lookups WHERE collection = ?`, s.Collection).Scan(&tablesJSON)
	if err != nil {
		return nil
	}

	var tables map[string][]string
	if err := json.Unmarshal([]byte(tablesJSON), &tables); err != nil {
		return nil
	}

	s.lookupMu.Lock()
	defer s.lookupMu.Unlock()

	for field, values := range tables {
		s.lookupTables[field] = values
		rev := make(map[string]uint32)
		for i, v := range values {
			rev[v] = uint32(i)
		}
		s.reverseLookup[field] = rev
	}

	return nil
}

// ── Query Builder ───────────────────────────────────────────────────────────

// buildWhere converts a filter map to a SQL WHERE clause with args.
// Supports exact match and operator maps ($gt, $gte, $lt, $lte, $ne, $in).
func (s *BinarySchema) buildWhere(filter map[string]interface{}) (string, []interface{}) {
	if len(filter) == 0 {
		return "", nil
	}

	var clauses []string
	var args []interface{}

	for k, v := range filter {
		if opMap, ok := v.(map[string]interface{}); ok {
			for op, opVal := range opMap {
				sqlOp := ""
				switch op {
				case "$gt":
					sqlOp = ">"
				case "$gte":
					sqlOp = ">="
				case "$lt":
					sqlOp = "<"
				case "$lte":
					sqlOp = "<="
				case "$ne":
					sqlOp = "!="
				case "$in":
					if arr, ok := opVal.([]interface{}); ok {
						ph := make([]string, len(arr))
						for i, item := range arr {
							ph[i] = "?"
							args = append(args, s.coerceIndexValue(k, item))
						}
						clauses = append(clauses, fmt.Sprintf("%s IN (%s)", k, strings.Join(ph, ",")))
					}
					continue
				}
				if sqlOp != "" {
					clauses = append(clauses, fmt.Sprintf("%s %s ?", k, sqlOp))
					args = append(args, s.coerceIndexValue(k, opVal))
				}
			}
		} else {
			clauses = append(clauses, k+" = ?")
			args = append(args, s.coerceIndexValue(k, v))
		}
	}

	return strings.Join(clauses, " AND "), args
}

func (s *BinarySchema) coerceIndexValue(fieldName string, val interface{}) interface{} {
	idx, ok := s.fieldIndex[fieldName]
	if !ok {
		return fmt.Sprintf("%v", val)
	}
	f := s.Fields[idx]
	switch f.Type {
	case "uint8", "uint16", "uint32", "uint64":
		return toInt(val)
	case "timestamp":
		return toTimestamp(val)
	case "lookup", "lookup16":
		return fmt.Sprintf("%v", val)
	default:
		return fmt.Sprintf("%v", val)
	}
}

// ── Helpers ─────────────────────────────────────────────────────────────────

func toByte(v interface{}) byte {
	switch val := v.(type) {
	case float64:
		return byte(val)
	case int:
		return byte(val)
	case int64:
		return byte(val)
	}
	return 0
}

func toUint16(v interface{}) uint16 {
	switch val := v.(type) {
	case float64:
		return uint16(val)
	case int:
		return uint16(val)
	case int64:
		return uint16(val)
	}
	return 0
}

func toUint32(v interface{}) uint32 {
	switch val := v.(type) {
	case float64:
		return uint32(val)
	case int:
		return uint32(val)
	case int64:
		return uint32(val)
	case string:
		var h uint32
		for _, c := range val {
			h = h*31 + uint32(c)
		}
		return h
	}
	return 0
}

func toUint64(v interface{}) uint64 {
	switch val := v.(type) {
	case float64:
		return uint64(val)
	case int:
		return uint64(val)
	case int64:
		return uint64(val)
	}
	return 0
}

func toInt(v interface{}) int64 {
	switch val := v.(type) {
	case float64:
		return int64(val)
	case int:
		return int64(val)
	case int64:
		return val
	case uint64:
		return int64(val)
	case string:
		n, _ := strconv.ParseInt(val, 10, 64)
		return n
	}
	return 0
}

func toTimestamp(v interface{}) int64 {
	switch val := v.(type) {
	case string:
		t, err := time.Parse(time.RFC3339, val)
		if err == nil {
			return t.Unix()
		}
		return 0
	case float64:
		return int64(val)
	case int64:
		return val
	case time.Time:
		return val.Unix()
	}
	return 0
}
