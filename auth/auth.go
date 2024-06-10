package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	firebase "firebase.google.com/go"
	firebase_auth "firebase.google.com/go/auth"
	"github.com/downsized-devs/sdk-go/codes"
	"github.com/downsized-devs/sdk-go/errors"
	"github.com/downsized-devs/sdk-go/logger"
	"github.com/downsized-devs/sdk-go/null"
	"github.com/downsized-devs/sdk-go/parser"
	identitytoolkitv1 "google.golang.org/api/identitytoolkit/v1"
	identitytoolkitv3 "google.golang.org/api/identitytoolkit/v3"
	"google.golang.org/api/option"
)

type contextKey string

const (
	userAuthInfo contextKey = "UserAuthInfo"
)

type Interface interface {
	VerifyToken(ctx context.Context, bearertoken string) (*firebase_auth.Token, error)
	GetUser(ctx context.Context, userParam FirebaseUserParam) ([]FirebaseUser, error)
	RegisterUser(ctx context.Context, user FirebaseUser) (FirebaseUser, error)
	UpdateUser(ctx context.Context, user FirebaseUser) (FirebaseUser, error)
	DeleteUser(ctx context.Context, userID string) error
	SetUserAuthInfo(ctx context.Context, param UserAuthParam) context.Context
	GetUserAuthInfo(ctx context.Context) (UserAuthInfo, error)
	RevokeUserRefreshToken(ctx context.Context, uid string) error
	VerifyPassword(ctx context.Context, email, password string) (bool, error)
	GetUsers(ctx context.Context, userParams []FirebaseUserParam) ([]FirebaseUser, error)
	SignInWithPassword(ctx context.Context, param UserLogin) (UserLoginResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (RefreshTokenResponse, error)
}

type auth struct {
	log               logger.Interface
	json              parser.JsonInterface
	firebase          *firebase_auth.Client
	identitytoolkitv1 *identitytoolkitv1.Service
	identitytoolkitv3 *identitytoolkitv3.Service
	httpClient        *http.Client
	conf              Config
}

type Config struct {
	SkipFirebaseInit bool
	Firebase         FirebaseConf
}

type FirebaseConf struct {
	AccountKey FirebaseAccountKey
	ApiKey     string
}

type FirebaseAccountKey struct {
	Type                    string `json:"type"`
	ProjectID               string `json:"project_id"`
	PrivateKeyID            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientID                string `json:"client_id"`
	AuthURI                 string `json:"auth_uri"`
	TokenURI                string `json:"token_uri"`
	AuthProviderx509CertURL string `json:"auth_provider_x509_cert_url"`
	Clientx509CertURL       string `json:"client_x509_cert_url"`
}

func Init(cfg Config, log logger.Interface, json parser.JsonInterface, httpClient *http.Client) Interface {
	if cfg.SkipFirebaseInit {
		return &auth{
			log:  log,
			json: json,
			conf: cfg,
		}
	}
	ctx := context.Background()

	accountkey, err := json.Marshal(cfg.Firebase.AccountKey)
	if err != nil {
		log.Fatal(ctx, err)
	}

	app, err := firebase.NewApp(context.Background(), nil, option.WithCredentialsJSON(accountkey))
	if err != nil {
		log.Fatal(ctx, err)
	}

	firebaseAuth, err := app.Auth(ctx)
	if err != nil {
		log.Fatal(ctx, err)
	}

	identitytoolv3, err := identitytoolkitv3.NewService(ctx, option.WithAPIKey(cfg.Firebase.ApiKey))
	if err != nil {
		log.Fatal(ctx, err)
	}

	identitytoolv1, err := identitytoolkitv1.NewService(ctx, option.WithAPIKey(cfg.Firebase.ApiKey))
	if err != nil {
		log.Fatal(ctx, err)
	}

	return &auth{
		log:               log,
		json:              json,
		firebase:          firebaseAuth,
		identitytoolkitv1: identitytoolv1,
		identitytoolkitv3: identitytoolv3,
		httpClient:        httpClient,
		conf:              cfg,
	}
}

func (a *auth) SignInWithPassword(ctx context.Context, param UserLogin) (UserLoginResponse, error) {
	if a.conf.SkipFirebaseInit {
		return UserLoginResponse{}, errors.NewWithCode(http.StatusNotImplemented, "Auth Not Initiated")
	}

	toolkitCall := a.identitytoolkitv1.Accounts.SignInWithPassword(&identitytoolkitv1.GoogleCloudIdentitytoolkitV1SignInWithPasswordRequest{
		Email:             param.Email,
		Password:          param.Password,
		ReturnSecureToken: true,
	})

	res, err := toolkitCall.Do()
	if err != nil {
		return UserLoginResponse{}, err
	}

	result := UserLoginResponse{
		Kind:           res.Kind,
		LocalID:        res.LocalId,
		Email:          res.Email,
		DisplayName:    res.DisplayName,
		IDToken:        res.IdToken,
		Registered:     res.Registered,
		ProfilePicture: res.ProfilePicture,
		RefreshToken:   res.RefreshToken,
		ExpiresIn:      res.ExpiresIn,
	}

	return result, nil
}

func (a *auth) RefreshToken(ctx context.Context, refreshToken string) (RefreshTokenResponse, error) {
	if a.conf.SkipFirebaseInit {
		return RefreshTokenResponse{}, errors.NewWithCode(http.StatusNotImplemented, "Auth Not Initiated")
	}

	return a.exchangeRefreshToken(ctx, RefreshTokenRequest{
		GrantType:    "refresh_token",
		RefreshToken: refreshToken,
	})
}

func (a *auth) GetUser(ctx context.Context, userParam FirebaseUserParam) ([]FirebaseUser, error) {
	if a.conf.SkipFirebaseInit {
		return []FirebaseUser{}, errors.NewWithCode(http.StatusNotImplemented, "Auth Not Initiated")
	}

	var users []FirebaseUser
	var param []firebase_auth.UserIdentifier

	switch {
	case userParam.ID != "":
		param = append(param, firebase_auth.UIDIdentifier{UID: userParam.ID})
	case userParam.Email != "":
		param = append(param, firebase_auth.EmailIdentifier{Email: userParam.Email})
	case userParam.PhoneNumber != "":
		param = append(param, firebase_auth.PhoneIdentifier{PhoneNumber: userParam.PhoneNumber})
	}

	res, err := a.firebase.GetUsers(ctx, param)
	if err != nil {
		return users, err
	}

	for _, v := range res.Users {
		user := a.assignUser(v)
		users = append(users, user)
	}

	return users, nil
}

func (a *auth) RegisterUser(ctx context.Context, user FirebaseUser) (FirebaseUser, error) {
	if a.conf.SkipFirebaseInit {
		return FirebaseUser{}, errors.NewWithCode(http.StatusNotImplemented, "Auth Not Initiated")
	}

	params := (&firebase_auth.UserToCreate{}).
		Email(user.Email).
		Password(user.Password).
		DisplayName(user.DisplayName)

	if user.PhoneNumber != "" {
		params.PhoneNumber(user.PhoneNumber)
	}

	if user.PhotoURL != "" {
		params.PhotoURL(user.PhotoURL)
	}

	if user.IsEmailVerified.Valid {
		params.EmailVerified(user.IsEmailVerified.Bool)
	}

	if user.IsDisabled.Valid {
		params.Disabled(user.IsEmailVerified.Bool)
	}

	u, err := a.firebase.CreateUser(ctx, params)
	if err != nil {
		return user, errors.NewWithCode(codes.CodeAuthFailure, "cannot create new user with err: %v", err)
	}

	newUser := FirebaseUser{
		ID:              u.UID,
		Email:           u.Email,
		IsEmailVerified: null.BoolFrom(u.EmailVerified),
		PhoneNumber:     u.PhoneNumber,
		DisplayName:     u.DisplayName,
		PhotoURL:        u.PhotoURL,
		IsDisabled:      null.BoolFrom(u.Disabled),
	}

	return newUser, nil
}

func (a *auth) UpdateUser(ctx context.Context, user FirebaseUser) (FirebaseUser, error) {
	if a.conf.SkipFirebaseInit {
		return FirebaseUser{}, errors.NewWithCode(http.StatusNotImplemented, "Auth Not Initiated")
	}

	params := (&firebase_auth.UserToUpdate{})

	if user.Email != "" {
		params.Email(user.Email)
	}

	if user.PhoneNumber != "" {
		params.PhoneNumber(user.PhoneNumber)
	}

	if user.DisplayName != "" {
		params.DisplayName(user.DisplayName)
	}

	if user.PhotoURL != "" {
		params.PhotoURL(user.PhotoURL)
	}

	if user.IsEmailVerified.Valid {
		params.EmailVerified(user.IsEmailVerified.Bool)
	}

	if user.IsDisabled.Valid {
		params.Disabled(user.IsEmailVerified.Bool)
	}

	if user.Password != "" {
		params.Password(user.Password)
	}

	u, err := a.firebase.UpdateUser(ctx, user.ID, params)
	if err != nil {
		return user, errors.NewWithCode(codes.CodeAuthFailure, fmt.Sprintf("cannot update user %s", user.ID), err)
	}

	updatedUser := FirebaseUser{
		ID:              u.UID,
		Email:           u.Email,
		IsEmailVerified: null.BoolFrom(u.EmailVerified),
		PhoneNumber:     u.PhoneNumber,
		DisplayName:     u.DisplayName,
		PhotoURL:        u.PhotoURL,
		IsDisabled:      null.BoolFrom(u.Disabled),
	}

	return updatedUser, nil
}

func (a *auth) DeleteUser(ctx context.Context, userID string) error {
	if a.conf.SkipFirebaseInit {
		return errors.NewWithCode(http.StatusNotImplemented, "Auth Not Initiated")
	}

	err := a.firebase.DeleteUser(ctx, userID)
	if err != nil {
		return errors.NewWithCode(codes.CodeAuthFailure, "cannot update user %s with err: %v", userID, err)
	}

	return nil
}

func (a *auth) VerifyToken(ctx context.Context, bearertoken string) (*firebase_auth.Token, error) {
	if a.conf.SkipFirebaseInit {
		return nil, errors.NewWithCode(http.StatusNotImplemented, "Auth Not Initiated")
	}

	token, err := a.firebase.VerifyIDTokenAndCheckRevoked(ctx, bearertoken)
	if err != nil {
		a.log.Error(ctx, errors.NewWithCode(codes.CodeAuthFailure, "failed to get token info with err %v", err))
		if strings.Contains(err.Error(), expiredTokenMessage) {
			return nil, errors.NewWithCode(codes.CodeAuthAccessTokenExpired, "token is expired with err: %v", err)
		} else if strings.Contains(err.Error(), revokedTokenMessage) {
			return nil, errors.NewWithCode(codes.CodeAuthRefreshTokenExpired, "token is revoked with err: %v", err)
		} else {
			return nil, errors.NewWithCode(codes.CodeAuthInvalidToken, "invalid token with err: %v", err)
		}
	}
	return token, nil
}

func (a *auth) SetUserAuthInfo(ctx context.Context, param UserAuthParam) context.Context {
	var token firebase_auth.Token
	if param.FirebaseToken != nil {
		token = *param.FirebaseToken
	}

	var userCredential UserCredential
	if param.UserCredential != nil {
		userCredential = *param.UserCredential
	}

	userauth := UserAuthInfo{
		User:           param.User,
		FirebaseToken:  token,
		UserCredential: userCredential,
	}
	return context.WithValue(ctx, userAuthInfo, userauth)
}

func (a *auth) GetUserAuthInfo(ctx context.Context) (UserAuthInfo, error) {
	user, ok := ctx.Value(userAuthInfo).(UserAuthInfo)
	if !ok {
		return user, errors.NewWithCode(codes.CodeAuthFailure, "failed getting user auth info")
	}

	return user, nil
}

func (a *auth) RevokeUserRefreshToken(ctx context.Context, uid string) error {
	if a.conf.SkipFirebaseInit {
		return errors.NewWithCode(http.StatusNotImplemented, "Auth Not Initiated")
	}

	err := a.firebase.RevokeRefreshTokens(ctx, uid)
	if err != nil {
		return errors.NewWithCode(codes.CodeAuthRevokeRefreshTokenFailed, err.Error())
	}
	return nil
}

func (a *auth) VerifyPassword(ctx context.Context, email, password string) (bool, error) {
	if a.conf.SkipFirebaseInit {
		return false, errors.NewWithCode(http.StatusNotImplemented, "Auth Not Initiated")
	}

	toolkitCall := a.identitytoolkitv3.Relyingparty.VerifyPassword(&identitytoolkitv3.IdentitytoolkitRelyingpartyVerifyPasswordRequest{
		Email:    email,
		Password: password,
	})

	_, err := toolkitCall.Do()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (a *auth) GetUsers(ctx context.Context, userParams []FirebaseUserParam) ([]FirebaseUser, error) {
	if a.conf.SkipFirebaseInit {
		return []FirebaseUser{}, errors.NewWithCode(http.StatusNotImplemented, "Auth Not Initiated")
	}

	var users []FirebaseUser
	var param []firebase_auth.UserIdentifier

	for _, userParam := range userParams {
		switch {
		case userParam.ID != "":
			param = append(param, firebase_auth.UIDIdentifier{UID: userParam.ID})
		case userParam.Email != "":
			param = append(param, firebase_auth.EmailIdentifier{Email: userParam.Email})
		case userParam.PhoneNumber != "":
			param = append(param, firebase_auth.PhoneIdentifier{PhoneNumber: userParam.PhoneNumber})

		}
	}

	res, err := a.firebase.GetUsers(ctx, param)
	if err != nil {
		return users, err
	}

	for _, v := range res.Users {
		user := a.assignUser(v)
		users = append(users, user)
	}

	return users, nil
}

func (a *auth) assignUser(firebaseUser *firebase_auth.UserRecord) FirebaseUser {
	user := FirebaseUser{
		ID:                 firebaseUser.UID,
		Email:              firebaseUser.Email,
		IsEmailVerified:    null.BoolFrom(firebaseUser.EmailVerified),
		PhoneNumber:        firebaseUser.PhoneNumber,
		DisplayName:        firebaseUser.DisplayName,
		PhotoURL:           firebaseUser.PhotoURL,
		IsDisabled:         null.BoolFrom(firebaseUser.Disabled),
		CreationTimestamp:  firebaseUser.UserMetadata.CreationTimestamp / 1000,
		LastLoginTimestamp: firebaseUser.UserMetadata.LastLogInTimestamp / 1000,
	}

	return user
}
