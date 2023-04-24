package types

type LoginSession struct {
	Id uint64 `json:"id,omitempty"`

	RoleType int `json:"role_type,omitempty"`
}
