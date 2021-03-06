package serializers

import (
	"clean/app/utils/consts"
	"clean/app/utils/methodsutil"
	"time"
)

type ResolveUserResponse struct {
	CompanyName  string             `json:"company_name"`
	CompanyID    uint               `json:"company_id"`
	Subordinates []*ResolveUserResp `json:"subordinates"`
}

type UserReq struct {
	UserName   string  `json:"user_name,omitempty"`
	FirstName  string  `json:"first_name,omitempty"`
	LastName   string  `json:"last_name,omitempty"`
	Email      string  `json:"email,omitempty"`
	Password   *string `json:"password,omitempty"`
	ProfilePic *string `json:"profile_pic,omitempty"`
	Phone      string  `json:"phone,omitempty"`
}

type LoggedInUser struct {
	ID          int      `json:"user_id"`
	CompanyID   uint     `json:"company_id"`
	AccessUuid  string   `json:"access_uuid"`
	RefreshUuid string   `json:"refresh_uuid"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
}

type UserResp struct {
	ID             uint       `json:"id"`
	Name           string     `json:"name"`
	Email          string     `json:"email"`
	TotalPoint     uint       `json:"total_point"`
	Phone          string     `json:"phone"`
	Specialization string     `json:"specialization"`
	LastLoginAt    *time.Time `json:"last_login_at"`
	FirstLogin     bool       `json:"first_login"`
}

type ResolveUserResp struct {
	ID          int        `json:"id"`
	UserName    string     `json:"user_name"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	Email       string     `json:"email"`
	Phone       *string    `json:"phone"`
	ProfilePic  *string    `json:"profile_pic"`
	AppKey      string     `json:"app_key,omitempty"`
	RoleID      uint       `json:"role_id"`
	CreatedAt   time.Time  `json:"-"`
	UpdatedAt   time.Time  `json:"-"`
	LastLoginAt *time.Time `json:"last_login_at"`
	FirstLogin  bool       `json:"first_login"`
}

type UserWithParamsResp struct {
	UserResp
	RoleName    string   `json:"role_name"`
	Permissions []string `json:"permissions,omitempty"`
}

type VerifyTokenResp struct {
	ID           int      `json:"id"`
	FirstName    string   `json:"first_name"`
	LastName     string   `json:"last_name"`
	Email        string   `json:"email"`
	Phone        *string  `json:"phone"`
	ProfilePic   *string  `json:"profile_pic"`
	BusinessID   *int     `json:"business_id"`
	BusinessName string   `json:"business_name"`
	CompanyID    *int     `json:"company_id"`
	CompanyName  string   `json:"company_name"`
	Permissions  []string `json:"permissions"`
	Admin        bool     `json:"admin"`
}

type UserWithLocations struct {
	UserID   uint   `json:"user_id"`
	UserName string `json:"user_name"`
}

// func (lu LoggedInUser) IsSuperAdmin() bool {
// 	return consts.RoleSuperAdmin == lu.Role
// }

func (lu LoggedInUser) IsAdmin() bool {
	return consts.RoleAdmin == lu.Role
}

func (lu LoggedInUser) IsSales() bool {
	return consts.RoleSales == lu.Role
}

func (lu LoggedInUser) HasPermission(perm string) bool {
	return methodsutil.InArray(perm, lu.Permissions)
}
