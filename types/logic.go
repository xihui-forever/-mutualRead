package types

type (
	LoginReq struct {
		RoleType int    `json:"role_type,omitempty"`
		Username string `json:"username,omitempty"`
		Password string `json:"password,omitempty"`
	}
	LoginRsp struct {
		Token  string      `json:"token,omitempty"`
		Expire uint32      `json:"expire,omitempty" help:"有效时长的秒数"`
		Info   interface{} `json:"info,omitempty"`
	}
)

type (
	ResetPasswordReq struct {
		RoleType int    `json:"role_type,omitempty"`
		Username string `json:"username,omitempty"`
		Password string `json:"password,omitempty"`
	}

	ChangePasswordReq struct {
		OldPassword string `json:"old_password,omitempty" validate:"required"`
		NewPassword string `json:"new_password,omitempty" validate:"required"`
	}
)
