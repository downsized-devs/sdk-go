package featureflag

import (
	"context"
	"time"

	"github.com/downsized-devs/sdk-go/logger"
	ffclient "github.com/thomaspoignant/go-feature-flag"
	"github.com/thomaspoignant/go-feature-flag/ffuser"
	"github.com/thomaspoignant/go-feature-flag/retriever"
)

const (
	userAnonymousType = "user-anonymous"
	userLoginType     = "user-login"
	userCustomType    = "user-custom"
)

type Interface interface {
	CheckUserFlags(flagKey, userKey, userType string, value map[string]interface{}) (bool, error)
	GetAllUserFlags(userKey, userType string, value map[string]interface{}) ([]byte, error)
	Refresh()
}

type Config struct {
	Enabled         bool
	PollingInterval time.Duration
}

type featureFlag struct {
	cfg       Config
	log       logger.Interface
	client    *ffclient.GoFeatureFlag
	retriever retriever.Retriever
}

func Init(cfg Config, log logger.Interface, ret retriever.Retriever) Interface {
	ffLag := &featureFlag{
		log:       log,
		cfg:       cfg,
		retriever: ret,
	}

	if ffLag.cfg.Enabled {
		ffLag.initClient()
	}

	return ffLag
}

func (ff *featureFlag) initClient() {
	ctx := context.Background()
	ffClient, err := ffclient.New(ffclient.Config{
		Context:                 ctx,
		PollingInterval:         ff.cfg.PollingInterval,
		Retriever:               ff.retriever,
		StartWithRetrieverError: true,
	})
	if err != nil {
		ff.log.Fatal(ctx, err)
	}

	ff.client = ffClient
}

func (ff *featureFlag) Refresh() {
	if ff.cfg.Enabled {
		// sleep for 500ms to avoid replication lag when refreshing the feature flag
		// on a database with read replica configuration
		time.Sleep(500 * time.Millisecond)

		ff.initClient()
	}
}

func (ff *featureFlag) CheckUserFlags(flagKey, userKey, userType string, value map[string]interface{}) (bool, error) {
	if !ff.cfg.Enabled {
		return false, nil
	}

	return ff.client.BoolVariation(flagKey, ff.getUserType(userKey, userType, value), false)
}

func (ff *featureFlag) GetAllUserFlags(userKey, userType string, value map[string]interface{}) ([]byte, error) {
	if !ff.cfg.Enabled {
		return []byte{}, nil
	}

	return ff.client.AllFlagsState(ff.getUserType(userKey, userType, value)).MarshalJSON()
}

func (ff *featureFlag) getUserType(userKey, userType string, value map[string]interface{}) ffuser.User {
	user := ffuser.User{}
	switch userType {
	case userAnonymousType:
		user = ff.registerAnonymousUser(userKey)
	case userLoginType:
		user = ff.registerNewUser(userKey)
	case userCustomType:
		user = ff.registerCustomUser(userKey, value)
	}
	return user
}

func (ff *featureFlag) registerAnonymousUser(userKey string) ffuser.User {
	return ffuser.NewAnonymousUser(userKey)
}

func (ff *featureFlag) registerNewUser(userKey string) ffuser.User {
	return ffuser.NewUser(userKey)
}

func (ff *featureFlag) registerCustomUser(userKey string, value map[string]interface{}) ffuser.User {
	user := ffuser.NewUserBuilder(userKey)
	for k, val := range value {
		user.AddCustom(k, val)
	}

	return user.Build()
}
