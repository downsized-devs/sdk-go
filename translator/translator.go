package translator

import (
	"context"
	"fmt"
	"os"

	"github.com/downsized-devs/sdk-go/appcontext"
	"github.com/downsized-devs/sdk-go/codes"
	"github.com/downsized-devs/sdk-go/errors"
	"github.com/downsized-devs/sdk-go/language"
	"github.com/downsized-devs/sdk-go/logger"
	"github.com/go-playground/locales"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
)

type Interface interface {
	Translate(ctx context.Context, key interface{}, params ...string) (string, error)
}

type Config struct {
	FallbackLanguageID   string
	SupportedLanguageIDs []string
	TranslationDir       string
}

type translator struct {
	translator *ut.UniversalTranslator
	log        logger.Interface
}

func Init(conf Config, log logger.Interface) Interface {
	fallback, supported, err := parseLanguageId(conf)
	if err != nil {
		log.Fatal(context.Background(), err)
	}

	t := ut.New(fallback, supported...)

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(context.Background(), err)
	}
	t.Import(ut.FormatJSON, fmt.Sprintf("%s/%s", pwd, conf.TranslationDir))
	if err := t.VerifyTranslations(); err != nil {
		log.Fatal(context.Background(), err)
	}

	return &translator{
		translator: t,
		log:        log,
	}
}

func parseLanguageId(conf Config) (locales.Translator, []locales.Translator, error) {
	var (
		locales []locales.Translator
	)

	localeIds := []string{conf.FallbackLanguageID}
	localeIds = append(localeIds, conf.SupportedLanguageIDs...)

	for _, v := range localeIds {
		switch v {
		case language.English:
			locales = append(locales, en.New())
		case language.Indonesian:
			locales = append(locales, id.New())
		default:
			return nil, nil, errors.NewWithCode(codes.CodeTranslatorError, fmt.Sprintf("unsupported languages ID %s", v))
		}
	}

	if len(locales) < 1 {
		return nil, nil, errors.NewWithCode(codes.CodeTranslatorError, "unsupported fallback language")
	}

	return locales[0], locales, nil
}

func (ut *translator) Translate(ctx context.Context, key interface{}, params ...string) (string, error) {
	if key == nil {
		return "", nil
	}

	strKey, ok := key.(string)
	if !ok {
		return "", errors.NewWithCode(codes.CodeTranslatorError, "key must be a string")
	}

	if strKey == "" {
		return "", nil
	}

	language := appcontext.GetAcceptLanguage(ctx)
	trans, found := ut.translator.GetTranslator(language)
	if !found {
		trans = ut.translator.GetFallback()
	}

	return trans.T(strKey, params...)

}
