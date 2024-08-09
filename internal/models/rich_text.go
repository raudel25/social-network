package models

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

// RichText representa un campo JSONB en tu modelo con dos campos: Text y HTML.
type RichText struct {
	Text string `json:"text"`
	HTML string `json:"html"`
}

// Scan implementa la interfaz Scanner para convertir un valor de la base de datos a RichText.
func (r *RichText) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		// Intenta convertir el slice de bytes a RichText.
		var rt RichText
		if err := json.Unmarshal(v, &rt); err != nil {
			return err
		}
		*r = rt
		return nil
	default:
		return sql.ErrConnDone
	}
}

// Value implementa la interfaz Valuer para convertir RichText a un valor que la base de datos pueda entender.
func (r RichText) Value() (driver.Value, error) {
	// Convierte RichText a []byte, que es el tipo que la base de datos espera para JSONB.
	return json.Marshal(r)
}
