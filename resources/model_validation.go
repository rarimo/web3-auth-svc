/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type Validation struct {
	Key
	Attributes ValidationAttributes `json:"attributes"`
}
type ValidationResponse struct {
	Data     Validation `json:"data"`
	Included Included   `json:"included"`
}

type ValidationListResponse struct {
	Data     []Validation    `json:"data"`
	Included Included        `json:"included"`
	Links    *Links          `json:"links"`
	Meta     json.RawMessage `json:"meta,omitempty"`
}

func (r *ValidationListResponse) PutMeta(v interface{}) (err error) {
	r.Meta, err = json.Marshal(v)
	return err
}

func (r *ValidationListResponse) GetMeta(out interface{}) error {
	return json.Unmarshal(r.Meta, out)
}

// MustValidation - returns Validation from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustValidation(key Key) *Validation {
	var validation Validation
	if c.tryFindEntry(key, &validation) {
		return &validation
	}
	return nil
}
