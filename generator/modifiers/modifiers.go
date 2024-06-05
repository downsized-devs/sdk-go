package modifiers

import (
	domainInitialize "github.com/downsized-devs/sdk-go/generator/modifiers/domain_initialize"
	domainItem "github.com/downsized-devs/sdk-go/generator/modifiers/domain_item"
	domainQuery "github.com/downsized-devs/sdk-go/generator/modifiers/domain_query"
	domainRedis "github.com/downsized-devs/sdk-go/generator/modifiers/domain_redis"
	domainTest "github.com/downsized-devs/sdk-go/generator/modifiers/domain_test"
	"github.com/downsized-devs/sdk-go/generator/modifiers/entity"
	handlerRest "github.com/downsized-devs/sdk-go/generator/modifiers/handler_rest"
	mockfile "github.com/downsized-devs/sdk-go/generator/modifiers/mock_file"
	restInitialize "github.com/downsized-devs/sdk-go/generator/modifiers/rest_initialize"
	"github.com/downsized-devs/sdk-go/generator/modifiers/sql"
	usecaseInitialize "github.com/downsized-devs/sdk-go/generator/modifiers/usecase_initialize"
	usecaseItem "github.com/downsized-devs/sdk-go/generator/modifiers/usecase_item"
	usecaseUnittest "github.com/downsized-devs/sdk-go/generator/modifiers/usecase_test"
)

type Modifiers struct {
	DomainQuery       domainQuery.Interface
	DomainRedis       domainRedis.Interface
	DomainItem        domainItem.Interface
	UsecaseItem       usecaseItem.Interface
	HandlerRest       handlerRest.Interface
	Entity            entity.Interface
	Sql               sql.Interface
	DomainInitialize  domainInitialize.Interface
	UseCaseInitialize usecaseInitialize.Interface
	RestInitialize    restInitialize.Interface
	DomainTest        domainTest.Interface
	UsecaseTest       usecaseUnittest.Interface
	MockFile          mockfile.Interface
}

type Param struct {
	EntityName           string
	EntityNameUpper      string
	EntityNameSnakeCase  string
	EntityNameLowerSpace string
	EntityNameUpperSpace string
	EntityNameLowerDash  string
	Location             string
	Api                  []string
}

func Init(params Param) *Modifiers {
	return &Modifiers{
		DomainQuery:       domainQuery.Init(domainQuery.Params{EntityName: params.EntityName, EntityNameUpper: params.EntityNameUpper, EntityNameSnakeCase: params.EntityNameSnakeCase, Location: params.Location}),
		DomainRedis:       domainRedis.Init(domainRedis.Params{EntityName: params.EntityName, EntityNameUpper: params.EntityNameUpper, EntityNameSnakeCase: params.EntityNameSnakeCase, Location: params.Location}),
		DomainItem:        domainItem.Init(domainItem.Params{EntityName: params.EntityName, EntityNameUpper: params.EntityNameUpper, EntityNameSnakeCase: params.EntityNameSnakeCase, EntityNameLowerSpace: params.EntityNameLowerSpace, Location: params.Location}),
		UsecaseItem:       usecaseItem.Init(usecaseItem.Params{EntityName: params.EntityName, EntityNameUpper: params.EntityNameUpper, EntityNameSnakeCase: params.EntityNameSnakeCase, Location: params.Location, Api: params.Api}),
		HandlerRest:       handlerRest.Init(handlerRest.Params{EntityName: params.EntityName, EntityNameUpper: params.EntityNameUpper, EntityNameSnakeCase: params.EntityNameSnakeCase, EntityNameLowerSpace: params.EntityNameLowerSpace, EntityNameUpperSpace: params.EntityNameUpperSpace, EntityNameLowerDash: params.EntityNameLowerDash, Location: params.Location, Api: params.Api}),
		Entity:            entity.Init(entity.Params{EntityNameUpper: params.EntityNameUpper, EntityNameSnakeCase: params.EntityNameSnakeCase, EntityName: params.EntityName, Location: params.Location}),
		Sql:               sql.Init(sql.Params{EntityNameUpperSpace: params.EntityNameUpperSpace, EntityNameSnakeCase: params.EntityNameSnakeCase, Location: params.Location}),
		DomainInitialize:  domainInitialize.Init(domainInitialize.Params{EntityNameUpper: params.EntityNameUpper, EntityNameSnakeCase: params.EntityNameSnakeCase, Location: params.Location}),
		UseCaseInitialize: usecaseInitialize.Init(usecaseInitialize.Params{EntityNameUpper: params.EntityNameUpper, EntityNameSnakeCase: params.EntityNameSnakeCase, Location: params.Location}),
		RestInitialize:    restInitialize.Init(restInitialize.Params{EntityNameUpper: params.EntityNameUpper, EntityNameSnakeCase: params.EntityNameSnakeCase, EntityNameLowerSpace: params.EntityNameLowerSpace, EntityNameLowerDash: params.EntityNameLowerDash, Location: params.Location, Api: params.Api}),
		DomainTest:        domainTest.Init(domainTest.Params{EntityName: params.EntityName, EntityNameUpper: params.EntityNameUpper, EntityNameSnakeCase: params.EntityNameSnakeCase, Location: params.Location}),
		UsecaseTest:       usecaseUnittest.Init(usecaseUnittest.Params{EntityName: params.EntityName, EntityNameUpper: params.EntityNameUpper, EntityNameSnakeCase: params.EntityNameSnakeCase, Location: params.Location, Api: params.Api}),
		MockFile:          mockfile.Init(mockfile.Params{EntityName: params.EntityName, EntityNameUpper: params.EntityNameUpper, EntityNameSnakeCase: params.EntityNameSnakeCase, Location: params.Location}),
	}
}
