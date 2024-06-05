package services

import (
	"errors"
	"os"

	"github.com/downsized-devs/sdk-go/generator/helper"
	"github.com/downsized-devs/sdk-go/generator/modifiers"
)

type services struct {
	EntityName string
	Location   string
	Api        []string
}

type Interface interface {
	Execute() error
}

type Param struct {
	EntityName string
	Location   string
	Api        []string
}

func Init(params Param) Interface {
	return &services{
		EntityName: params.EntityName,
		Location:   params.Location,
		Api:        params.Api,
	}
}

func (s *services) Execute() error {
	modifiers, err := s.initialize()
	if err != nil {
		return err
	}
	err = s.runModifiers(modifiers)
	if err != nil {
		return err
	}
	return nil
}

func (s *services) runModifiers(modifiers *modifiers.Modifiers) error {
	if !s.IsFileExist() {
		err := modifiers.DomainQuery.Replace()
		if err != nil {
			return errors.New("domain query " + err.Error())
		}
		err = modifiers.DomainRedis.Replace()
		if err != nil {
			return errors.New("domain redis " + err.Error())
		}

		err = modifiers.DomainItem.Replace()
		if err != nil {
			return errors.New("domain item " + err.Error())
		}

		err = modifiers.UsecaseItem.Replace()
		if err != nil {
			return errors.New("usecase item " + err.Error())
		}
		err = modifiers.HandlerRest.Replace()
		if err != nil {
			return errors.New("handler rest " + err.Error())
		}
		err = modifiers.Entity.Replace()
		if err != nil {
			return errors.New("entity " + err.Error())
		}
		err = modifiers.Sql.Replace()
		if err != nil {
			return errors.New("sql " + err.Error())
		}
		err = modifiers.DomainInitialize.Initialize()
		if err != nil {
			return errors.New("domain initialize " + err.Error())
		}
		err = modifiers.UseCaseInitialize.Initialize()
		if err != nil {
			return errors.New("usecase initialize " + err.Error())
		}
		err = modifiers.DomainTest.Replace()
		if err != nil {
			return errors.New("domain test " + err.Error())
		}
		err = modifiers.UsecaseTest.Replace()
		if err != nil {
			return errors.New("usecase test " + err.Error())
		}
		err = modifiers.MockFile.Replace()
		if err != nil {
			return errors.New("mock file " + err.Error())
		}
		err = modifiers.MockFile.Initialize()
		if err != nil {
			return errors.New("mock file initialize " + err.Error())
		}
	}
	err := modifiers.UsecaseItem.AppendInterfaceAndFunction()
	if err != nil {
		return errors.New("usecase append interface and function " + err.Error())
	}
	err = modifiers.HandlerRest.AppendInterfaceAndFunction()
	if err != nil {
		return errors.New("handler append interface and function " + err.Error())
	}
	err = modifiers.RestInitialize.Initialize()
	if err != nil {
		return errors.New("rest initialize " + err.Error())
	}
	err = modifiers.UsecaseTest.AppendInterfaceAndFunction()
	if err != nil {
		return errors.New("usecase test append interface and function " + err.Error())
	}
	return nil
}

func (s *services) initialize() (*modifiers.Modifiers, error) {
	modifiers := modifiers.Init(
		modifiers.Param{
			EntityName:           s.EntityName,
			EntityNameUpper:      helper.ConvertToUpperCase(s.EntityName),
			EntityNameSnakeCase:  helper.ConvertToSnakeCase(s.EntityName),
			EntityNameLowerSpace: helper.ConvertToLowerSpace(s.EntityName),
			EntityNameUpperSpace: helper.ConvertToUpperSpace(s.EntityName),
			EntityNameLowerDash:  helper.ConvertToLowerDash(s.EntityName),
			Location:             s.Location,
			Api:                  s.Api,
		},
	)
	return modifiers, nil
}

func (s *services) IsFileExist() bool {
	_, err := os.OpenFile(s.Location+"/src/business/usecase/"+helper.ConvertToSnakeCase(s.EntityName)+"/"+helper.ConvertToSnakeCase(s.EntityName)+".go", os.O_RDWR, 0644)
	return err == nil
}
