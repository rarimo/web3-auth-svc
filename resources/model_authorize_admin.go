/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type AuthorizeAdmin struct {
	Key
	Attributes AuthorizeAdminAttributes `json:"attributes"`
}
type AuthorizeAdminRequest struct {
	Data     AuthorizeAdmin `json:"data"`
	Included Included       `json:"included"`
}

type AuthorizeAdminListRequest struct {
	Data     []AuthorizeAdmin `json:"data"`
	Included Included         `json:"included"`
	Links    *Links           `json:"links"`
	Meta     json.RawMessage  `json:"meta,omitempty"`
}

func (r *AuthorizeAdminListRequest) PutMeta(v interface{}) (err error) {
	r.Meta, err = json.Marshal(v)
	return err
}

func (r *AuthorizeAdminListRequest) GetMeta(out interface{}) error {
	return json.Unmarshal(r.Meta, out)
}

// MustAuthorizeAdmin - returns AuthorizeAdmin from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustAuthorizeAdmin(key Key) *AuthorizeAdmin {
	var authorizeAdmin AuthorizeAdmin
	if c.tryFindEntry(key, &authorizeAdmin) {
		return &authorizeAdmin
	}
	return nil
}
