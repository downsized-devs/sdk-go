package auth

import (
	firebase_auth "firebase.google.com/go/auth"
	"github.com/downsized-devs/sdk-go/null"
)

const (
	revokedTokenMessage = "ID token has been revoked"
	expiredTokenMessage = "ID token has expired"
)

type Token struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type FirebaseUser struct {
	ID                 string    `json:"id"`
	Email              string    `json:"email"`
	IsEmailVerified    null.Bool `json:"is_email_verified"`
	PhoneNumber        string    `json:"phone_number"`
	Password           string    `json:"password"`
	DisplayName        string    `json:"display_name"`
	PhotoURL           string    `json:"photo_url"`
	IsDisabled         null.Bool `json:"is_disabled"`
	CreationTimestamp  int64     `json:"creation_timestamp"`
	LastLoginTimestamp int64     `json:"last_login_timestamp"`
}

type FirebaseUserParam struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

type UserAuthInfo struct {
	User           User                `json:"user"`
	FirebaseToken  firebase_auth.Token `json:"firebaseToken"`
	UserCredential UserCredential      `json:"userCredential"`
}

type UserAuthParam struct {
	User           User                 `json:"user"`
	FirebaseToken  *firebase_auth.Token `json:"firebaseToken"`
	UserCredential *UserCredential      `json:"userCredential"`
}

type UserCredential struct {
	ID           int64     `db:"id" json:"id"`
	UserID       int64     `db:"fk_user_id" json:"userId"`
	ServiceID    int64     `db:"fk_service_id" json:"serviceId"`
	AccessToken  string    `db:"access_token" json:"accessToken"`
	RefreshToken string    `db:"refresh_token" json:"refreshToken"`
	UserAgent    string    `db:"user_agent" json:"userAgent"`
	ExpiredAt    null.Time `db:"expired_at" json:"expiredAt"`
	IsRevoke     bool      `db:"is_revoke" json:"isRevoke"`
}

type User struct {
	ID          int64  `db:"id" json:"id"`
	CompanyID   int64  `db:"fk_company_id" json:"companyId"`
	Name        string `db:"name" json:"name"`
	Email       string `db:"email" json:"email"`
	UID         string `db:"uid" json:"uid"`
	RoleID      int64  `db:"fk_role_id" json:"roleId"`
	RoleRank    int64  `db:"rank" json:"roleRank"`
	PhoneNumber string `db:"phone_num" json:"phoneNumber"`
	IsQA        bool   `db:"is_qa" json:"isQa"`
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RefreshTokenRequest struct {
	GrantType    string `json:"grant_type"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenResponse struct {
	ExpiresIn    string `json:"expires_in"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
	UserID       string `json:"user_id"`
	ProjectID    string `json:"project_id"`
}

type UserRefreshTokenParam struct {
	RefreshToken string `form:"refreshToken"`
}

type UserLoginResponse struct {
	Kind           string `json:"kind"`
	LocalID        string `json:"localId"`
	Email          string `json:"email"`
	DisplayName    string `json:"displayName"`
	IDToken        string `json:"idToken"`
	Registered     bool   `json:"registered"`
	ProfilePicture string `json:"profilePicture"`
	RefreshToken   string `json:"refreshToken"`
	ExpiresIn      int64  `json:"expiresIn"`
}
