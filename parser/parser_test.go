package parser

import (
	"testing"

	mock_log "github.com/downsized-devs/sdk-go/tests/mock/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestInit(t *testing.T) {
	ctrl := gomock.NewController(t)
	log := mock_log.NewMockInterface(ctrl)

	p := Init(log, Options{})
	assert.NotNil(t, p)
	assert.NotNil(t, p.JsonParser())
	assert.NotNil(t, p.CsvParser())
}

func TestInit_ReusesSameInstances(t *testing.T) {
	ctrl := gomock.NewController(t)
	log := mock_log.NewMockInterface(ctrl)

	p := Init(log, Options{})
	assert.Same(t, p.JsonParser(), p.JsonParser(), "JsonParser should be a stable handle")
	assert.Same(t, p.CsvParser(), p.CsvParser(), "CsvParser should be a stable handle")
}
