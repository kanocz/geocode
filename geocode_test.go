package geocode

import (
	"encoding/json"
	"testing"
)

// Make sure Bounds is not included in Geometry field, since it is an optional field
func TestMarshallingResponseNoBounds(t *testing.T) {
	googleResponse := []byte(`{"results": [ { "address_components": [ { "long_name": "1600", "short_name": "1600", "types": [ "street_number" ] }, { "long_name": "Amphitheatre Pkwy", "short_name": "Amphitheatre Pkwy", "types": [ "route" ] }, { "long_name": "Mountain View", "short_name": "Mountain View", "types": [ "locality", "political" ] }, { "long_name": "Santa Clara", "short_name": "Santa Clara", "types": [ "administrative_area_level_2", "political" ] }, { "long_name": "California", "short_name": "CA", "types": [ "administrative_area_level_1", "political" ] }, { "long_name": "United States", "short_name": "US", "types": [ "country", "political" ] }, { "long_name": "94043", "short_name": "94043", "types": [ "postal_code" ] } ], "formatted_address": "1600 Amphitheatre Pkwy, Mountain View, CA 94043, USA", "geometry": { "location": { "lat": 37.4229181, "lng": -122.0854212 }, "location_type": "ROOFTOP", "viewport": { "northeast": { "lat": 37.42426708029149, "lng": -122.0840722197085 }, "southwest": { "lat": 37.4215691197085, "lng": -122.0867701802915 } } }, "types": [ "street_address" ] } ], "status": "OK" }
`)

	resp := new(Response)
	err := json.Unmarshal(googleResponse, resp)
	if err != nil {
		t.Fatalf("Decoding error: %v", err)
	}

	// Test it did basic marshal
	if resp.Status != "OK" {
		t.Errorf("Status == %q, want %q", resp.Status, "OK")
	}

	// Make sure no bounds is found because is optional
	for _, addr := range resp.Results {
		if addr.Geometry.Bounds != nil {
			t.Errorf("Bounds found, do not want: %v", addr)
		}
	}

	b, err := json.Marshal(resp)
	// blind test, but make sure keys are all lower case
	expectedResponse := `{"status":"OK","error_message":"","results":[{"formatted_address":"1600 Amphitheatre Pkwy, Mountain View, CA 94043, USA","address_components":[{"long_name":"1600","short_name":"1600","types":["street_number"]},{"long_name":"Amphitheatre Pkwy","short_name":"Amphitheatre Pkwy","types":["route"]},{"long_name":"Mountain View","short_name":"Mountain View","types":["locality","political"]},{"long_name":"Santa Clara","short_name":"Santa Clara","types":["administrative_area_level_2","political"]},{"long_name":"California","short_name":"CA","types":["administrative_area_level_1","political"]},{"long_name":"United States","short_name":"US","types":["country","political"]},{"long_name":"94043","short_name":"94043","types":["postal_code"]}],"geometry":{"location":{"lat":37.4229181,"lng":-122.0854212},"location_type":"ROOFTOP","viewport":{"northeast":{"lat":37.42426708029149,"lng":-122.0840722197085},"southwest":{"lat":37.4215691197085,"lng":-122.0867701802915}}},"types":["street_address"]}]}`

	if string(b) != expectedResponse {
		t.Errorf("Expected JSON is not correct: %s", b)
	}
}
