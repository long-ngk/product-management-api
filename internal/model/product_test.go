package model

import (
	"encoding/json"
	"regexp"
	"strings"
	"testing"
	"time"

	"pgregory.net/rapid"
)

// TestProperty7_ProductJSONSerializationPreservesFormatInvariants verifies that
// for any Product with any valid price and any timestamp, serializing to JSON produces:
// - price with exactly 2 decimal places
// - created_at and updated_at in RFC 3339 format with "Z" suffix
// - description present as JSON null when the value is nil (not omitted)
func TestProperty7_ProductJSONSerializationPreservesFormatInvariants(t *testing.T) {
	// Regex for a number with exactly 2 decimal places (e.g., 120.50, 0.00, 99999.99)
	priceRegex := regexp.MustCompile(`^-?\d+\.\d{2}$`)

	// Regex for RFC 3339 timestamp with "Z" suffix
	rfc3339ZRegex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`)

	rapid.Check(t, func(t *rapid.T) {
		// Generate random price (positive float with various magnitudes)
		price := rapid.Float64Range(0.01, 999999.99).Draw(t, "price")

		// Generate random timestamps across a wide range
		minTime := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
		maxTime := time.Date(2030, 12, 31, 23, 59, 59, 0, time.UTC)
		createdAtUnix := rapid.Int64Range(minTime.Unix(), maxTime.Unix()).Draw(t, "createdAtUnix")
		updatedAtUnix := rapid.Int64Range(minTime.Unix(), maxTime.Unix()).Draw(t, "updatedAtUnix")
		createdAt := time.Unix(createdAtUnix, 0).UTC()
		updatedAt := time.Unix(updatedAtUnix, 0).UTC()

		// Generate nil or non-nil description
		isNilDescription := rapid.Bool().Draw(t, "isNilDescription")
		var description *string
		if !isNilDescription {
			desc := rapid.StringMatching(`[a-zA-Z0-9 ]{1,100}`).Draw(t, "description")
			description = &desc
		}

		product := Product{
			ID:          rapid.IntRange(1, 100000).Draw(t, "id"),
			Name:        rapid.StringMatching(`[a-zA-Z]{3,50}`).Draw(t, "name"),
			Description: description,
			Price:       price,
			Quantity:    rapid.IntRange(0, 10000).Draw(t, "quantity"),
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		}

		// Marshal to JSON
		data, err := json.Marshal(product)
		if err != nil {
			t.Fatalf("failed to marshal product: %v", err)
		}

		// Parse the JSON into a map to inspect individual fields
		var result map[string]json.RawMessage
		if err := json.Unmarshal(data, &result); err != nil {
			t.Fatalf("failed to unmarshal JSON: %v", err)
		}

		// Assert: price has exactly 2 decimal places
		priceRaw := strings.TrimSpace(string(result["price"]))
		if !priceRegex.MatchString(priceRaw) {
			t.Fatalf("price does not have exactly 2 decimal places: got %s", priceRaw)
		}

		// Assert: created_at is RFC 3339 with "Z" suffix
		var createdAtStr string
		if err := json.Unmarshal(result["created_at"], &createdAtStr); err != nil {
			t.Fatalf("failed to unmarshal created_at: %v", err)
		}
		if !rfc3339ZRegex.MatchString(createdAtStr) {
			t.Fatalf("created_at is not in RFC 3339 Z format: got %s", createdAtStr)
		}

		// Assert: updated_at is RFC 3339 with "Z" suffix
		var updatedAtStr string
		if err := json.Unmarshal(result["updated_at"], &updatedAtStr); err != nil {
			t.Fatalf("failed to unmarshal updated_at: %v", err)
		}
		if !rfc3339ZRegex.MatchString(updatedAtStr) {
			t.Fatalf("updated_at is not in RFC 3339 Z format: got %s", updatedAtStr)
		}

		// Assert: when description is nil, JSON has "description":null (not omitted)
		descRaw, exists := result["description"]
		if !exists {
			t.Fatal("description field is missing from JSON (should be present as null when nil)")
		}
		if isNilDescription {
			if strings.TrimSpace(string(descRaw)) != "null" {
				t.Fatalf("expected description to be null in JSON, got: %s", string(descRaw))
			}
		}

		// Also verify timestamps can be parsed back as valid time
		_, err = time.Parse(time.RFC3339, createdAtStr)
		if err != nil {
			t.Fatalf("created_at is not a valid RFC 3339 timestamp: %v", err)
		}
		_, err = time.Parse(time.RFC3339, updatedAtStr)
		if err != nil {
			t.Fatalf("updated_at is not a valid RFC 3339 timestamp: %v", err)
		}
	})
}
