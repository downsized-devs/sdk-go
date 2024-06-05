package email

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/downsized-devs/sdk-go/log"
	"github.com/stretchr/testify/assert"
)

func Test_emailtemplate_FromHTML(t *testing.T) {
	pwd, _ := os.Getwd()
	fileDir := fmt.Sprintf("%s/test_files", pwd)
	email := Init(Config{Template: TemplateConfig{FileDirectory: fileDir}}, log.Init(log.Config{Level: "debug"}))
	tests := []struct {
		name    string
		params  BodyFromHTMLParams
		want    string
		wantErr bool
	}{
		{
			name: "success",
			params: BodyFromHTMLParams{
				Filename: "test.html",
				Data: map[string]string{
					"Text": "hello world",
				},
			},
			want: "<html>\n    <strong>hello world</strong>\n</html>\n",
		},
		{
			name: "success from cache",
			params: BodyFromHTMLParams{
				Filename: "test.html",
				Data: map[string]string{
					"Text": "hello world",
				},
			},
			want: "<html>\n    <strong>hello world</strong>\n</html>\n",
		},
		{
			name: "success with funcmap",
			params: BodyFromHTMLParams{
				Filename: "test2.html",
				Data: map[string]string{
					"Text": "hello world",
				},
				FuncMap: map[string]any{
					"add": func(a, b int) int {
						return a + b
					},
				},
			},
			want: "<html>\n    <strong>3</strong>\n</html>\n",
		},
		{
			name: "file not found",
			params: BodyFromHTMLParams{
				Filename: "notfound.html",
				Data: map[string]string{
					"Text": "hello world",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := email.GenerateBody().FromHTML(tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("emailtemplate.FromHTML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_emailtemplate_FromMJML(t *testing.T) {
	pwd, _ := os.Getwd()
	fileDir := fmt.Sprintf("%s/test_files", pwd)
	email := Init(Config{Template: TemplateConfig{FileDirectory: fileDir}}, log.Init(log.Config{Level: "debug"}))
	ctxTimeout, cancel := context.WithTimeout(context.Background(), 1)
	defer cancel()
	type args struct {
		ctx    context.Context
		params BodyFromMJMLParams
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				params: BodyFromMJMLParams{
					Filename: "test.mjml",
					Data: map[string]string{
						"Text": "hello world",
					},
				},
			},
			want: "<!doctype html><html xmlns=\"http://www.w3.org/1999/xhtml\" xmlns:v=\"urn:schemas-microsoft-com:vml\" xmlns:o=\"urn:schemas-microsoft-com:office:office\"><head><title></title><!--[if !mso]><!--><meta http-equiv=\"X-UA-Compatible\" content=\"IE=edge\"><!--<![endif]--><meta http-equiv=\"Content-Type\" content=\"text/html; charset=UTF-8\"><meta name=\"viewport\" content=\"width=device-width,initial-scale=1\"><style type=\"text/css\">#outlook a { padding:0; }\n      body { margin:0;padding:0;-webkit-text-size-adjust:100%;-ms-text-size-adjust:100%; }\n      table, td { border-collapse:collapse;mso-table-lspace:0pt;mso-table-rspace:0pt; }\n      img { border:0;height:auto;line-height:100%; outline:none;text-decoration:none;-ms-interpolation-mode:bicubic; }\n      p { display:block;margin:13px 0; }</style><!--[if mso]>\n    <noscript>\n    <xml>\n    <o:OfficeDocumentSettings>\n      <o:AllowPNG/>\n      <o:PixelsPerInch>96</o:PixelsPerInch>\n    </o:OfficeDocumentSettings>\n    </xml>\n    </noscript>\n    <![endif]--><!--[if lte mso 11]>\n    <style type=\"text/css\">\n      .mj-outlook-group-fix { width:100% !important; }\n    </style>\n    <![endif]--><!--[if !mso]><!--><link href=\"https://fonts.googleapis.com/css?family=Ubuntu:300,400,500,700\" rel=\"stylesheet\" type=\"text/css\"><style type=\"text/css\">@import url(https://fonts.googleapis.com/css?family=Ubuntu:300,400,500,700);</style><!--<![endif]--><style type=\"text/css\">@media only screen and (min-width:480px) {\n        .mj-column-per-100 { width:100% !important; max-width: 100%; }\n      }</style><style media=\"screen and (min-width:480px)\">.moz-text-html .mj-column-per-100 { width:100% !important; max-width: 100%; }</style><style type=\"text/css\"></style><style type=\"text/css\"></style></head><body style=\"word-spacing:normal;\"><div><!--[if mso | IE]><table align=\"center\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\" class=\"\" role=\"presentation\" style=\"width:600px;\" width=\"600\" ><tr><td style=\"line-height:0px;font-size:0px;mso-line-height-rule:exactly;\"><![endif]--><div style=\"margin:0px auto;max-width:600px;\"><table align=\"center\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\" role=\"presentation\" style=\"width:100%;\"><tbody><tr><td style=\"direction:ltr;font-size:0px;padding:20px 0;text-align:center;\"><!--[if mso | IE]><table role=\"presentation\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\"><tr><td class=\"\" style=\"vertical-align:top;width:600px;\" ><![endif]--><div class=\"mj-column-per-100 mj-outlook-group-fix\" style=\"font-size:0px;text-align:left;direction:ltr;display:inline-block;vertical-align:top;width:100%;\"><table border=\"0\" cellpadding=\"0\" cellspacing=\"0\" role=\"presentation\" style=\"vertical-align:top;\" width=\"100%\"><tbody><tr><td align=\"center\" style=\"font-size:0px;padding:10px 25px;word-break:break-word;\"><p style=\"border-top:solid 4px #F45E43;font-size:1px;margin:0px auto;width:100%;\"></p><!--[if mso | IE]><table align=\"center\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\" style=\"border-top:solid 4px #F45E43;font-size:1px;margin:0px auto;width:550px;\" role=\"presentation\" width=\"550px\" ><tr><td style=\"height:0;line-height:0;\"> &nbsp;\n</td></tr></table><![endif]--></td></tr><tr><td align=\"left\" style=\"font-size:0px;padding:10px 25px;word-break:break-word;\"><div style=\"font-family:Ubuntu, Helvetica, Arial, sans-serif;font-size:20px;line-height:1;text-align:left;color:#F45E43;\">hello world</div></td></tr></tbody></table></div><!--[if mso | IE]></td></tr></table><![endif]--></td></tr></tbody></table></div><!--[if mso | IE]></td></tr></table><![endif]--></div></body></html>",
		},
		{
			name: "success from cache",
			args: args{
				ctx: context.Background(),
				params: BodyFromMJMLParams{
					Filename: "test.mjml",
					Data: map[string]string{
						"Text": "hello world",
					},
				},
			},
			want: "<!doctype html><html xmlns=\"http://www.w3.org/1999/xhtml\" xmlns:v=\"urn:schemas-microsoft-com:vml\" xmlns:o=\"urn:schemas-microsoft-com:office:office\"><head><title></title><!--[if !mso]><!--><meta http-equiv=\"X-UA-Compatible\" content=\"IE=edge\"><!--<![endif]--><meta http-equiv=\"Content-Type\" content=\"text/html; charset=UTF-8\"><meta name=\"viewport\" content=\"width=device-width,initial-scale=1\"><style type=\"text/css\">#outlook a { padding:0; }\n      body { margin:0;padding:0;-webkit-text-size-adjust:100%;-ms-text-size-adjust:100%; }\n      table, td { border-collapse:collapse;mso-table-lspace:0pt;mso-table-rspace:0pt; }\n      img { border:0;height:auto;line-height:100%; outline:none;text-decoration:none;-ms-interpolation-mode:bicubic; }\n      p { display:block;margin:13px 0; }</style><!--[if mso]>\n    <noscript>\n    <xml>\n    <o:OfficeDocumentSettings>\n      <o:AllowPNG/>\n      <o:PixelsPerInch>96</o:PixelsPerInch>\n    </o:OfficeDocumentSettings>\n    </xml>\n    </noscript>\n    <![endif]--><!--[if lte mso 11]>\n    <style type=\"text/css\">\n      .mj-outlook-group-fix { width:100% !important; }\n    </style>\n    <![endif]--><!--[if !mso]><!--><link href=\"https://fonts.googleapis.com/css?family=Ubuntu:300,400,500,700\" rel=\"stylesheet\" type=\"text/css\"><style type=\"text/css\">@import url(https://fonts.googleapis.com/css?family=Ubuntu:300,400,500,700);</style><!--<![endif]--><style type=\"text/css\">@media only screen and (min-width:480px) {\n        .mj-column-per-100 { width:100% !important; max-width: 100%; }\n      }</style><style media=\"screen and (min-width:480px)\">.moz-text-html .mj-column-per-100 { width:100% !important; max-width: 100%; }</style><style type=\"text/css\"></style><style type=\"text/css\"></style></head><body style=\"word-spacing:normal;\"><div><!--[if mso | IE]><table align=\"center\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\" class=\"\" role=\"presentation\" style=\"width:600px;\" width=\"600\" ><tr><td style=\"line-height:0px;font-size:0px;mso-line-height-rule:exactly;\"><![endif]--><div style=\"margin:0px auto;max-width:600px;\"><table align=\"center\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\" role=\"presentation\" style=\"width:100%;\"><tbody><tr><td style=\"direction:ltr;font-size:0px;padding:20px 0;text-align:center;\"><!--[if mso | IE]><table role=\"presentation\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\"><tr><td class=\"\" style=\"vertical-align:top;width:600px;\" ><![endif]--><div class=\"mj-column-per-100 mj-outlook-group-fix\" style=\"font-size:0px;text-align:left;direction:ltr;display:inline-block;vertical-align:top;width:100%;\"><table border=\"0\" cellpadding=\"0\" cellspacing=\"0\" role=\"presentation\" style=\"vertical-align:top;\" width=\"100%\"><tbody><tr><td align=\"center\" style=\"font-size:0px;padding:10px 25px;word-break:break-word;\"><p style=\"border-top:solid 4px #F45E43;font-size:1px;margin:0px auto;width:100%;\"></p><!--[if mso | IE]><table align=\"center\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\" style=\"border-top:solid 4px #F45E43;font-size:1px;margin:0px auto;width:550px;\" role=\"presentation\" width=\"550px\" ><tr><td style=\"height:0;line-height:0;\"> &nbsp;\n</td></tr></table><![endif]--></td></tr><tr><td align=\"left\" style=\"font-size:0px;padding:10px 25px;word-break:break-word;\"><div style=\"font-family:Ubuntu, Helvetica, Arial, sans-serif;font-size:20px;line-height:1;text-align:left;color:#F45E43;\">hello world</div></td></tr></tbody></table></div><!--[if mso | IE]></td></tr></table><![endif]--></td></tr></tbody></table></div><!--[if mso | IE]></td></tr></table><![endif]--></div></body></html>",
		},
		{
			name: "success with funcmap",
			args: args{
				ctx: context.Background(),
				params: BodyFromMJMLParams{
					Filename: "test2.mjml",
					Data: map[string]string{
						"Text": "hello world",
					},
					FuncMap: map[string]any{
						"add": func(a, b int) int {
							return a + b
						},
					},
				},
			},
			want: "<!doctype html><html xmlns=\"http://www.w3.org/1999/xhtml\" xmlns:v=\"urn:schemas-microsoft-com:vml\" xmlns:o=\"urn:schemas-microsoft-com:office:office\"><head><title></title><!--[if !mso]><!--><meta http-equiv=\"X-UA-Compatible\" content=\"IE=edge\"><!--<![endif]--><meta http-equiv=\"Content-Type\" content=\"text/html; charset=UTF-8\"><meta name=\"viewport\" content=\"width=device-width,initial-scale=1\"><style type=\"text/css\">#outlook a { padding:0; }\n      body { margin:0;padding:0;-webkit-text-size-adjust:100%;-ms-text-size-adjust:100%; }\n      table, td { border-collapse:collapse;mso-table-lspace:0pt;mso-table-rspace:0pt; }\n      img { border:0;height:auto;line-height:100%; outline:none;text-decoration:none;-ms-interpolation-mode:bicubic; }\n      p { display:block;margin:13px 0; }</style><!--[if mso]>\n    <noscript>\n    <xml>\n    <o:OfficeDocumentSettings>\n      <o:AllowPNG/>\n      <o:PixelsPerInch>96</o:PixelsPerInch>\n    </o:OfficeDocumentSettings>\n    </xml>\n    </noscript>\n    <![endif]--><!--[if lte mso 11]>\n    <style type=\"text/css\">\n      .mj-outlook-group-fix { width:100% !important; }\n    </style>\n    <![endif]--><!--[if !mso]><!--><link href=\"https://fonts.googleapis.com/css?family=Ubuntu:300,400,500,700\" rel=\"stylesheet\" type=\"text/css\"><style type=\"text/css\">@import url(https://fonts.googleapis.com/css?family=Ubuntu:300,400,500,700);</style><!--<![endif]--><style type=\"text/css\">@media only screen and (min-width:480px) {\n        .mj-column-per-100 { width:100% !important; max-width: 100%; }\n      }</style><style media=\"screen and (min-width:480px)\">.moz-text-html .mj-column-per-100 { width:100% !important; max-width: 100%; }</style><style type=\"text/css\"></style><style type=\"text/css\"></style></head><body style=\"word-spacing:normal;\"><div><!--[if mso | IE]><table align=\"center\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\" class=\"\" role=\"presentation\" style=\"width:600px;\" width=\"600\" ><tr><td style=\"line-height:0px;font-size:0px;mso-line-height-rule:exactly;\"><![endif]--><div style=\"margin:0px auto;max-width:600px;\"><table align=\"center\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\" role=\"presentation\" style=\"width:100%;\"><tbody><tr><td style=\"direction:ltr;font-size:0px;padding:20px 0;text-align:center;\"><!--[if mso | IE]><table role=\"presentation\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\"><tr><td class=\"\" style=\"vertical-align:top;width:600px;\" ><![endif]--><div class=\"mj-column-per-100 mj-outlook-group-fix\" style=\"font-size:0px;text-align:left;direction:ltr;display:inline-block;vertical-align:top;width:100%;\"><table border=\"0\" cellpadding=\"0\" cellspacing=\"0\" role=\"presentation\" style=\"vertical-align:top;\" width=\"100%\"><tbody><tr><td align=\"center\" style=\"font-size:0px;padding:10px 25px;word-break:break-word;\"><p style=\"border-top:solid 4px #F45E43;font-size:1px;margin:0px auto;width:100%;\"></p><!--[if mso | IE]><table align=\"center\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\" style=\"border-top:solid 4px #F45E43;font-size:1px;margin:0px auto;width:550px;\" role=\"presentation\" width=\"550px\" ><tr><td style=\"height:0;line-height:0;\"> &nbsp;\n</td></tr></table><![endif]--></td></tr><tr><td align=\"left\" style=\"font-size:0px;padding:10px 25px;word-break:break-word;\"><div style=\"font-family:Ubuntu, Helvetica, Arial, sans-serif;font-size:20px;line-height:1;text-align:left;color:#F45E43;\">3</div></td></tr></tbody></table></div><!--[if mso | IE]></td></tr></table><![endif]--></td></tr></tbody></table></div><!--[if mso | IE]></td></tr></table><![endif]--></div></body></html>",
		},
		{
			name: "file not found",
			args: args{
				ctx: context.Background(),
				params: BodyFromMJMLParams{
					Filename: "notfound.mjml",
					Data: map[string]string{
						"Text": "hello world",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "convert timeout",
			args: args{
				ctx: ctxTimeout,
				params: BodyFromMJMLParams{
					Filename: "test.mjml",
					Data: map[string]string{
						"Text": "hello world",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := email.GenerateBody().FromMJML(tt.args.ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("emailtemplate.FromMJML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
