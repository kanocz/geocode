// Package geocode is an interface to The Google Geocoding API.
//
// See http://code.google.com/apis/maps/documentation/geocoding/ for details.
package geocode

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

const api = "https://maps.googleapis.com/maps/api/geocode/json"

type Request struct {
	// One (and only one) of these must be set.
	Address  string
	Location *Point

	// Optional fields.
	Bounds     *Bounds // Lookup within this viewport.
	Region     string
	Language   string
	Components string
	Channel    string

	Sensor bool

	// Google credential
	Googleclient    string
	Googlesignature string
	Googleapikey    string

	values url.Values
}

func (r *Request) Values() url.Values {
	if r.values == nil {
		r.values = make(url.Values)
	}
	var v = r.values
	if r.Address != "" {
		v.Set("address", r.Address)
	} else if r.Location != nil {
		v.Set("latlng", r.Location.String())
	} else {
		// well, the request will probably fail
		// let's return an empty string?
		return v
	}
	if r.Channel != "" {
		v.Set("channel", r.Channel)
	}
	if r.Bounds != nil {
		v.Set("bounds", r.Bounds.String())
	}
	if r.Region != "" {
		v.Set("region", r.Region)
	}
	if r.Language != "" {
		v.Set("language", r.Language)
	}
	if r.Components != "" {
		v.Set("components", r.Components)
	}
	if r.Googleclient != "" {
		v.Set("client", r.Googleclient)
	}
	if r.Googlesignature != "" && r.Googleclient != "" {
		v.Set("signature", r.Googlesignature)
	}
	if r.Googleapikey != "" && (r.Googleclient == "") {
		v.Set("key", r.Googleapikey)
	}

	v.Set("sensor", strconv.FormatBool(r.Sensor))

	return v
}

// Lookup makes the Request to the Google Geocoding API servers using
// the provided transport (or http.DefaultTransport if nil).
func (r *Request) Lookup(transport http.RoundTripper) (*Response, error) {
	c := http.Client{Transport: transport}
	params := r.Values().Encode()
	if len(params) == 0 {
		return nil, fmt.Errorf("Missing address or latlng argument")
	}
	u := fmt.Sprintf("%s?%s", api, params)
	getResp, err := c.Get(u)
	if err != nil {
		return nil, err
	}
	defer getResp.Body.Close()

	if getResp.StatusCode < 200 || getResp.StatusCode >= 300 {
		body, _ := ioutil.ReadAll(getResp.Body)
		return nil, fmt.Errorf("Failed to lookup address (code %d): %s", getResp.StatusCode, body)
	}

	resp := new(Response)
	err = json.NewDecoder(getResp.Body).Decode(resp)
	if err != nil {
		return nil, err
	}

	switch resp.Status {
	case "OVER_QUERY_LIMIT", "REQUEST_DENIED", "INVALID_REQUEST", "UNKNOWN_ERROR":
		return nil, fmt.Errorf("Lookup failed (%s): %s", resp.Status, resp.ErrorMessage)
	default:
		return resp, nil
	}
}

type Response struct {
	Status       string    `json:"status"`
	ErrorMessage string    `json:"error_message"`
	Results      []*Result `json:"results"`
}

type Result struct {
	Address      string         `json:"formatted_address"`
	AddressParts []*AddressPart `json:"address_components"`
	Geometry     *Geometry      `json:"geometry"`
	Types        []string       `json:"types"`
	PartialMatch bool           `json:"partial_match,omitempty"`
	PlaceId      string         `json:"place_id,omitempty"`
}

type AddressPart struct {
	Name      string   `json:"long_name"`
	ShortName string   `json:"short_name"`
	Types     []string `json:"types"`
}

type Geometry struct {
	Bounds   *Bounds `json:"bounds,omitempty"`
	Location Point   `json:"location"`
	Type     string  `json:"location_type"`
	Viewport Bounds  `json:"viewport"`
}

type Bounds struct {
	NorthEast Point `json:"northeast"`
	SouthWest Point `json:"southwest"`
}

func (b Bounds) String() string {
	return fmt.Sprintf("%v|%v", b.NorthEast, b.SouthWest)
}

type Point struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

func (p Point) String() string {
	return fmt.Sprintf("%g,%g", p.Lat, p.Lng)
}

type Address struct {
	PlaceID      string  `json:"placeId"`
	PartialMatch bool    `json:"partial_match,omitempty"`
	AddrStr      string  `json:"address"`
	Number       string  `json:"number"`
	Premise      string  `json:"premise"`
	Street       string  `json:"street"`
	City         string  `json:"city"`
	Country      string  `json:"country"`
	Postcode     string  `json:"postcode"`
	Lat          float64 `json:"lat"`
	Lng          float64 `json:"lng"`
}

// Parse is simple function that returns geocode addresses in parsed format
func (r *Response) Parse(includePartialMatch bool, requireNumber bool) []Address {
	if r.Status != "OK" {
		return nil
	}

	result := make([]Address, 0, len(r.Results))

	for _, re := range r.Results {
		if re.PartialMatch && !includePartialMatch {
			continue
		}

		addr := Address{}

		addr.AddrStr = re.Address
		addr.PlaceID = re.PlaceId
		addr.PartialMatch = re.PartialMatch
		addr.Lat = re.Geometry.Location.Lat
		addr.Lng = re.Geometry.Location.Lng

		for _, part := range re.AddressParts {
			for _, partType := range part.Types {
				switch partType {
				case "premise":
					addr.Premise = part.Name
				case "street_number":
					addr.Number = part.Name
				case "route":
					addr.Street = part.ShortName
				case "locality":
					addr.City = part.Name
				case "postal_town":
					addr.City = part.Name
				case "country":
					addr.Country = part.ShortName
				case "postal_code":
					addr.Postcode = part.Name
				}
			}
		}

		if requireNumber && "" == addr.Number && "" == addr.Premise {
			continue
		}

		result = append(result, addr)
	}

	return result
}
