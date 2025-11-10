package json

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"unicode"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var httpCodeToAppCode = map[int]string{
	http.StatusBadRequest:          "BAD_REQUEST",
	http.StatusUnauthorized:        "UNAUTHORIZED",
	http.StatusForbidden:           "FORBIDDEN",
	http.StatusNotFound:            "NOT_FOUND",
	http.StatusConflict:            "CONFLICT",
	422:                            "VALIDATION_ERROR",
	429:                            "TOO_MANY_REQUESTS",
	http.StatusInternalServerError: "INTERNAL_ERROR",
	http.StatusServiceUnavailable:  "SERVICE_UNAVAILABLE",
}

func ParseProtoJSON(r io.Reader, m proto.Message) error {
	bodyBytes, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("cannot read request body: %w", err)
	}

	unmarshalOpts := protojson.UnmarshalOptions{
		DiscardUnknown: false,
	}

	if err := unmarshalOpts.Unmarshal(bodyBytes, m); err != nil {
		return errors.New("invalid JSON format or data")
	}
	return nil
}

// ParseJSON decodes the request body into the provided model.
func ParseJSON(r *http.Request, model any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}
	return json.NewDecoder(r.Body).Decode(model)
}

type jsonNumberMarshaler struct {
	json.Number
}

func (j jsonNumberMarshaler) MarshalJSON() ([]byte, error) {
	if j.Number == "" {
		return []byte("0"), nil
	}
	i, err := j.Int64()
	if err != nil {
		return []byte(fmt.Sprintf("%v", j.Number)), nil
	}
	return []byte(fmt.Sprintf("%d", i)), nil
}

var _ = jsonNumberMarshaler{}

// WriteJSON writes the data v as JSON with the specified HTTP status.
// For protobuf messages, it marshals using protojson with EmitUnpopulated,
// then post-processes the JSON to (1) parse numeric fields using proto reflection,
// (2) fill missing fields without overwriting valid values, and
// (3) normalize keys to snake_case.
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	var data []byte
	var err error

	if protoMsg, ok := v.(proto.Message); ok {
		marshaler := protojson.MarshalOptions{
			EmitUnpopulated: true,
			UseProtoNames:   true, // forces keys defined in the proto (ideally snake_case)
			UseEnumNumbers:  true,
		}

		// Marshal the protobuf message to JSON.
		data, err = marshaler.Marshal(protoMsg)
		if err != nil {
			return err
		}

		// Decode the JSON into a map.
		var objMap map[string]interface{}
		decoder := json.NewDecoder(bytes.NewReader(data))
		decoder.UseNumber()
		if err := decoder.Decode(&objMap); err != nil {
			return err
		}

		// Parse numeric fields based on the protobuf definition.
		convertNumericFields(protoMsg.ProtoReflect(), objMap)

		// Recursively fill missing fields using proto reflection,
		// but do not override valid values.
		fillMissingFields(protoMsg.ProtoReflect(), objMap)

		// Normalize all keys to snake_case.
		objMap = normalizeKeysMap(objMap)

		// Marshal the modified map back to JSON.
		data, err = json.Marshal(objMap)
		if err != nil {
			return err
		}
	} else {
		data, err = json.Marshal(v)
		if err != nil {
			return err
		}
	}

	_, err = w.Write(data)
	return err
}

// fillMissingFields recursively processes the proto message m and ensures that
// for each field the normalized key (snake_case) exists in objMap. It will not
// overwrite a value that is already non-nil. If a duplicate camelCase key exists,
// its non-nil value is used if the normalized key is missing or nil.
func fillMissingFields(m protoreflect.Message, objMap map[string]interface{}) {
	fields := m.Descriptor().Fields()
	for i := 0; i < fields.Len(); i++ {
		fd := fields.Get(i)

		protoFieldName := string(fd.Name())    // proto name (usually snake_case)
		camelKey := fd.JSONName()              // camelCase
		normalizedKey := toSnakeCase(camelKey) // canonical snake_case

		// 1) Consolidate value to the canonical key, preserving null
		// Priority: existing normalizedKey (even if nil) -> protoFieldName -> camelKey -> nil
		val, has := objMap[normalizedKey]
		if !has {
			if v, ok := objMap[protoFieldName]; ok {
				val = v // may be nil (=> JSON null)
			} else if v, ok := objMap[camelKey]; ok {
				val = v // may be nil (=> JSON null)
			} else {
				val = nil // explicitly ensure presence as null
			}
			objMap[normalizedKey] = val
		}

		// Remove aliases to avoid duplicates
		if protoFieldName != normalizedKey {
			delete(objMap, protoFieldName)
		}
		if camelKey != normalizedKey {
			delete(objMap, camelKey)
		}

		// 2) Recurse into nested messages
		if fd.Kind() == protoreflect.MessageKind {
			// List of messages
			if fd.IsList() {
				if listIface, ok := objMap[normalizedKey].([]interface{}); ok {
					pfList := m.Get(fd).List()
					for idx := 0; idx < pfList.Len() && idx < len(listIface); idx++ {
						if subMap, ok := listIface[idx].(map[string]interface{}); ok {
							fillMissingFields(pfList.Get(idx).Message(), subMap)
						}
					}
				}
				continue
			}
			// Singular message
			if subMap, ok := objMap[normalizedKey].(map[string]interface{}); ok {
				fillMissingFields(m.Get(fd).Message(), subMap)
			}
		}
	}
}

// normalizeKeysMap creates a new map where every key is converted to snake_case.
// It recurses into nested maps and slices.
func normalizeKeysMap(m map[string]interface{}) map[string]interface{} {
	newMap := make(map[string]interface{})
	for k, v := range m {
		newKey := toSnakeCase(k)
		switch val := v.(type) {
		case map[string]interface{}:
			newMap[newKey] = normalizeKeysMap(val)
		case []interface{}:
			newSlice := make([]interface{}, len(val))
			for i, item := range val {
				if subMap, ok := item.(map[string]interface{}); ok {
					newSlice[i] = normalizeKeysMap(subMap)
				} else {
					newSlice[i] = item
				}
			}
			newMap[newKey] = newSlice
		default:
			newMap[newKey] = v
		}
	}
	return newMap
}

// toSnakeCase converts a given string from camelCase (or mixed) to snake_case.
func toSnakeCase(s string) string {
	var result []rune
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result = append(result, '_')
			}
			result = append(result, unicode.ToLower(r))
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}

// convertNumericFields walks the proto message m and converts any numeric field
// values in objMap from strings or json.Number to native Go numeric types. It
// recurses into nested messages and lists.
func convertNumericFields(m protoreflect.Message, objMap map[string]interface{}) {
	fields := m.Descriptor().Fields()
	for i := 0; i < fields.Len(); i++ {
		fd := fields.Get(i)
		camel := fd.JSONName()
		snake := toSnakeCase(camel)

		keys := []string{snake}
		if camel != snake {
			keys = append(keys, camel)
		}

		for _, key := range keys {
			val, exists := objMap[key]
			if !exists {
				continue
			}

			// Handle lists
			if fd.IsList() {
				list, ok := val.([]interface{})
				if !ok {
					continue
				}
				for idx := range list {
					elem := list[idx]
					if fd.Kind() == protoreflect.MessageKind {
						subMap, ok := elem.(map[string]interface{})
						if ok {
							pfList := m.Get(fd).List()
							if idx < pfList.Len() {
								convertNumericFields(pfList.Get(idx).Message(), subMap)
							}
						}
					} else if isNumericKind(fd.Kind()) {
						list[idx] = convertScalar(fd.Kind(), elem)
					}
				}
				objMap[key] = list
				continue
			}

			if fd.Kind() == protoreflect.MessageKind {
				if subMap, ok := val.(map[string]interface{}); ok {
					convertNumericFields(m.Get(fd).Message(), subMap)
					objMap[key] = subMap
				}
				continue
			}

			if isNumericKind(fd.Kind()) {
				objMap[key] = convertScalar(fd.Kind(), val)
			}
		}
	}
}

// convertScalar converts val to a Go numeric type according to kind if val is a
// string or json.Number. On conversion failure the original value is returned.
func convertScalar(kind protoreflect.Kind, val interface{}) interface{} {
	switch v := val.(type) {
	case string:
		if isFloatKind(kind) {
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				return f
			}
		} else if isUintKind(kind) {
			if ui, err := strconv.ParseUint(v, 10, 64); err == nil {
				return ui
			}
		} else if isSignedKind(kind) {
			if i, err := strconv.ParseInt(v, 10, 64); err == nil {
				return i
			}
		}
	case json.Number:
		if isFloatKind(kind) {
			if f, err := v.Float64(); err == nil {
				return f
			}
		} else if isUintKind(kind) {
			if ui, err := strconv.ParseUint(v.String(), 10, 64); err == nil {
				return ui
			}
		} else if isSignedKind(kind) {
			if i, err := v.Int64(); err == nil {
				return i
			}
		}
	}
	return val
}

// isNumericKind reports whether kind represents a numeric protobuf scalar.
func isNumericKind(kind protoreflect.Kind) bool {
	switch kind {
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind,
		protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind,
		protoreflect.Uint32Kind, protoreflect.Fixed32Kind,
		protoreflect.Uint64Kind, protoreflect.Fixed64Kind,
		protoreflect.FloatKind, protoreflect.DoubleKind:
		return true
	default:
		return false
	}
}

func isFloatKind(kind protoreflect.Kind) bool {
	return kind == protoreflect.FloatKind || kind == protoreflect.DoubleKind
}

func isSignedKind(kind protoreflect.Kind) bool {
	switch kind {
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind,
		protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return true
	default:
		return false
	}
}

func isUintKind(kind protoreflect.Kind) bool {
	switch kind {
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind,
		protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return true
	default:
		return false
	}
}

func WriteError(w http.ResponseWriter, status int, err error) {
	appCode, ok := httpCodeToAppCode[status]
	if !ok {
		appCode = "INTERNAL_ERROR"
		status = http.StatusInternalServerError
	}

	resp := HTTPError{
		Code:    appCode,
		Message: err.Error(),
	}

	if writeErr := WriteJSON(w, status, resp); writeErr != nil {
		http.Error(w, "failed to write error response", http.StatusInternalServerError)
	}
}

func WriteValidationError(w http.ResponseWriter, details interface{}) {
	resp := HTTPError{
		Code:    "VALIDATION_ERROR",
		Message: "validation failed",
		Details: details,
	}
	_ = WriteJSON(w, 422, resp)
}

type HTTPError struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}
