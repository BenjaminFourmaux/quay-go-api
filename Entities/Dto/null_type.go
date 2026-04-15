package Dto

import (
	"database/sql"
	"encoding/json"
)

// NullString is a custom type for sql.NullString
type NullString sql.NullString

// MarshalJSON for NullString to handle null values
func (ns NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

// UnmarshalJSON for NullString to handle deserialization
func (ns *NullString) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		ns.Valid = false
		ns.String = ""
		return nil
	}
	ns.Valid = true
	return json.Unmarshal(data, &ns.String)
}
