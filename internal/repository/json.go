package repository

import "encoding/json"

func encodingJSONMarshal(v any) ([]byte, error)   { return json.Marshal(v) }
func encodingJSONUnmarshal(b []byte, v any) error { return json.Unmarshal(b, v) }
