package parser

import (
	"context"
	"fmt"
	"strings"

	"github.com/downsized-devs/sdk-go/codes"
	"github.com/downsized-devs/sdk-go/errors"
	"github.com/downsized-devs/sdk-go/log"
	jsoniter "github.com/json-iterator/go"
	"github.com/xeipuuv/gojsonschema"
)

type jsonConfig string

const (
	vanillaCompatible jsonConfig = "standard"
	defaultConfig     jsonConfig = "default"
	fastestConfig     jsonConfig = "fastest"
	customConfig      jsonConfig = "custom"
	errJSON                      = `JSON parser error %s`
)

type JSONOptions struct {
	Config                        jsonConfig
	IndentionStep                 int
	MarshalFloatWith6Digits       bool
	EscapeHTML                    bool
	SortMapKeys                   bool
	UseNumber                     bool
	DisallowUnknownFields         bool
	TagKey                        string
	OnlyTaggedField               bool
	ValidateJSONRawMessage        bool
	ObjectFieldMustBeSimpleString bool
	CaseSensitive                 bool
	// Schema contains schema definitions with key as schema name and value as source path
	// schema sources can be file or URL. Schema definition will be iniatialized during
	// JSON parser object initialization.
	Schema map[string]string
}

type JSONInterface interface {
	// Marshal go structs into bytes
	Marshal(orig interface{}) ([]byte, error)
	// Marshal go structs into bytes and validates returned bytes
	MarshalWithSchemaValidation(sch string, orig interface{}) ([]byte, error)
	// Unmarshal bytes into go structs
	Unmarshal(blob []byte, dest interface{}) error
	// Validates input bytes and Unmarshal
	UnmarshalWithSchemaValidation(sch string, blob []byte, dest interface{}) error
}

type jsonParser struct {
	API    jsoniter.API
	schema map[string]*gojsonschema.Schema
	logger log.Interface
}

func initJSON(opt JSONOptions, log log.Interface) JSONInterface {
	var jsonAPI jsoniter.API
	switch opt.Config {
	case defaultConfig:
		jsonAPI = jsoniter.ConfigDefault
	case fastestConfig:
		jsonAPI = jsoniter.ConfigFastest
	case customConfig:
		jsonAPI = jsoniter.Config{
			IndentionStep:                 opt.IndentionStep,
			MarshalFloatWith6Digits:       opt.MarshalFloatWith6Digits,
			EscapeHTML:                    opt.EscapeHTML,
			SortMapKeys:                   opt.SortMapKeys,
			UseNumber:                     opt.UseNumber,
			DisallowUnknownFields:         opt.DisallowUnknownFields,
			TagKey:                        opt.TagKey,
			OnlyTaggedField:               opt.OnlyTaggedField,
			ValidateJsonRawMessage:        opt.ValidateJSONRawMessage,
			ObjectFieldMustBeSimpleString: opt.ObjectFieldMustBeSimpleString,
			CaseSensitive:                 opt.CaseSensitive,
		}.Froze()
	default:
		jsonAPI = jsoniter.ConfigCompatibleWithStandardLibrary
	}

	p := &jsonParser{
		API:    jsonAPI,
		schema: make(map[string]*gojsonschema.Schema),
		logger: log,
	}

	// init defined schema
	p.initSchema(opt.Schema)

	return p
}

// initSchema initialize all defined schema from sources
func (p *jsonParser) initSchema(sources map[string]string) {
	for sch, src := range sources {
		schema, err := gojsonschema.NewSchema(gojsonschema.NewReferenceLoader(src))
		if err != nil {
			p.logger.Fatal(context.Background(), errors.NewWithCode(codes.CodeJSONSchemaInvalid, fmt.Sprintf("error on load : %s", sch), err))
			return
		}
		p.schema[sch] = schema
	}
}

// Marshal marshal input blobs into go structs.
func (p *jsonParser) Marshal(orig interface{}) ([]byte, error) {
	stream := p.API.BorrowStream(nil)
	defer p.API.ReturnStream(stream)
	stream.WriteVal(orig)
	result := make([]byte, stream.Buffered())
	if stream.Error != nil {
		return nil, errors.NewWithCode(codes.CodeJSONMarshalError, stream.Error.Error())
	}
	copy(result, stream.Buffer())
	return result, nil
}

// MarshalWithSchemaValidation marshals go structs (interface {}) and validate returned bytes using defined schema reference
// It marshal the original struct then validate. It should marshall first then validate to ensure produced bytes
// matches the schema definition.
func (p *jsonParser) MarshalWithSchemaValidation(sch string, orig interface{}) ([]byte, error) {
	blob, err := p.Marshal(orig)
	if err != nil {
		return nil, err
	}
	s, ok := p.schema[sch]
	if !ok {
		return nil, errors.NewWithCode(codes.CodeJSONSchemaNotFound, fmt.Sprintf("schema not found : %s", sch))
	}

	blobLoader := gojsonschema.NewBytesLoader(blob)
	res, err := s.Validate(blobLoader)
	if err != nil {
		return nil, errors.NewWithCode(codes.CodeJSONValidationError, err.Error())
	}

	if !res.Valid() {
		var errString []string
		for _, desc := range res.Errors() {
			errString = append(errString, fmt.Sprintf("- %s", desc))
		}
		return nil, errors.NewWithCode(codes.CodeJSONValidationError, errJSON, strings.Join(errString, "\n"))
	}

	return blob, nil
}

// Unmarshal unmarshal input blobs into go structs.
func (p *jsonParser) Unmarshal(blob []byte, dest interface{}) error {
	iter := p.API.BorrowIterator(blob)
	defer p.API.ReturnIterator(iter)
	iter.ReadVal(dest)
	if iter.Error != nil {
		return errors.NewWithCode(codes.CodeJSONUnmarshalError, iter.Error.Error())
	}
	return nil
}

// UnmarshalWithSchemaValidation validates input bytes (blob) with respective schema definition
// and unmarshal them into go structs.
func (p *jsonParser) UnmarshalWithSchemaValidation(sch string, blob []byte, dest interface{}) error {
	s, ok := p.schema[sch]
	if !ok {
		return errors.NewWithCode(codes.CodeJSONSchemaNotFound, errJSON, fmt.Sprintf("schema not found : %s", sch))
	}
	blobLoader := gojsonschema.NewBytesLoader(blob)
	res, err := s.Validate(blobLoader)
	if err != nil {
		return errors.NewWithCode(codes.CodeJSONRawInvalid, errJSON, err)
	}
	if !res.Valid() {
		var errString []string
		for _, desc := range res.Errors() {
			errString = append(errString, fmt.Sprintf("- %s", desc))
		}
		return errors.NewWithCode(codes.CodeJSONValidationError, errJSON, strings.Join(errString, "\n"))
	}
	return p.Unmarshal(blob, dest)
}
