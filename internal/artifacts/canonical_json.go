package artifacts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
)

const maxCanonicalSafeInteger int64 = 9007199254740991

func isJSONContentType(contentType string) bool {
	trimmed := strings.TrimSpace(strings.ToLower(contentType))
	if trimmed == "application/json" {
		return true
	}
	return strings.HasPrefix(trimmed, "application/json;")
}

func canonicalizeJSONBytes(payload []byte) ([]byte, error) {
	decoder := json.NewDecoder(bytes.NewReader(payload))
	decoder.UseNumber()
	var value any
	if err := decoder.Decode(&value); err != nil {
		return nil, err
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return nil, fmt.Errorf("extra trailing JSON data")
		}
		return nil, err
	}
	canonical, err := canonicalizeJSONValue(value)
	if err != nil {
		return nil, err
	}
	return []byte(canonical), nil
}

func canonicalizeJSONValue(value any) (string, error) {
	switch typed := value.(type) {
	case nil:
		return "null", nil
	case bool:
		if typed {
			return "true", nil
		}
		return "false", nil
	case string:
		encoded, err := json.Marshal(typed)
		if err != nil {
			return "", err
		}
		return string(encoded), nil
	case json.Number:
		return canonicalizeNumber(typed)
	case []any:
		return canonicalizeArray(typed)
	case map[string]any:
		return canonicalizeObject(typed)
	default:
		return "", fmt.Errorf("unsupported JSON type %T", value)
	}
}

func canonicalizeNumber(value json.Number) (string, error) {
	parsed, err := canonicalIntegerFromText(value.String(), "number")
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(parsed, 10), nil
}

func canonicalizeArray(values []any) (string, error) {
	parts := make([]string, 0, len(values))
	for _, item := range values {
		canonical, err := canonicalizeJSONValue(item)
		if err != nil {
			return "", err
		}
		parts = append(parts, canonical)
	}
	return "[" + strings.Join(parts, ",") + "]", nil
}

func canonicalizeObject(object map[string]any) (string, error) {
	keys := make([]string, 0, len(object))
	for key := range object {
		if !isASCIIString(key) {
			return "", fmt.Errorf("object key %q is outside the MVP ASCII-only canonicalization profile", key)
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		keyJSON, err := canonicalizeJSONValue(key)
		if err != nil {
			return "", err
		}
		valueJSON, err := canonicalizeJSONValue(object[key])
		if err != nil {
			return "", err
		}
		parts = append(parts, keyJSON+":"+valueJSON)
	}
	return "{" + strings.Join(parts, ",") + "}", nil
}

func canonicalIntegerFromText(text string, location string) (int64, error) {
	parsed, err := strconv.ParseInt(text, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%s = %q is not a supported integer: %w", location, text, err)
	}
	if parsed < -maxCanonicalSafeInteger || parsed > maxCanonicalSafeInteger {
		return 0, fmt.Errorf("%s = %q is outside the shared Go/TS safe integer range", location, text)
	}
	return parsed, nil
}

func isASCIIString(text string) bool {
	for _, r := range text {
		if r > 0x7f {
			return false
		}
	}
	return true
}
