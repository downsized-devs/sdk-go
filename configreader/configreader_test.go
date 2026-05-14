package configreader

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type basicConfig struct {
	String   string
	Boolean  bool
	Duration time.Duration
}

type additionalConfig struct {
	Duration time.Duration
}

type mockConfig struct {
	String    string
	Boolean   bool
	Duration  time.Duration
	Object    basicConfig
	Reference basicConfig
	AddConfig additionalConfig
}

func TestInit(t *testing.T) {
	type args struct {
		options Options
	}
	tests := []struct {
		name      string
		args      args
		wantPanic bool
		panicMsg  string
	}{
		{
			name: "init success",
			args: args{
				options: Options{
					ConfigFile: "./files/conf.json",
					AdditionalConfig: []AdditionalConfigOptions{
						{
							ConfigKey:  "AddConfig",
							ConfigFile: "./files/conf-excel.json",
						},
					},
				},
			},
		},
		{
			name:      "init config panic",
			wantPanic: true,
			panicMsg:  "fatal error found during reading file. err: Config File \"config\" Not Found in \"[]\"",
		},
		{
			name: "init additional config panic",
			args: args{
				options: Options{
					ConfigFile: "./files/conf.json",
					AdditionalConfig: []AdditionalConfigOptions{
						{
							ConfigKey: "AddConfig",
						},
					},
				},
			},
			wantPanic: true,
			panicMsg:  "fatal error found during reading additional config file. err: Config File \"config\" Not Found in \"[]\"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				assert.PanicsWithError(t, tt.panicMsg, func() { Init(tt.args.options) })
			} else {
				Init(tt.args.options)
			}
		})
	}
}

func Test_config_ReadConfig(t *testing.T) {

	type args struct {
		options Options
	}
	tests := []struct {
		name      string
		args      args
		want      interface{}
		wantPanic bool
		panicMsg  string
	}{
		{
			name: "read config",
			args: args{
				options: Options{
					ConfigFile: "./files/conf.json",
					AdditionalConfig: []AdditionalConfigOptions{
						{
							ConfigKey:  "AddConfig",
							ConfigFile: "./files/conf-excel.json",
						},
					},
				},
			},
			want: mockConfig{
				String:   "Test String",
				Boolean:  true,
				Duration: 60000000000,
				Object: basicConfig{
					String:   "Test String",
					Boolean:  true,
					Duration: 60000000000,
				},
				Reference: basicConfig{
					String:   "Test String",
					Boolean:  true,
					Duration: 60000000000,
				},
				AddConfig: additionalConfig{
					Duration: 1000000000,
				},
			},
		},
		{
			name: "read empty config",
			args: args{
				options: Options{
					ConfigFile: "./files/empty-conf.json",
					AdditionalConfig: []AdditionalConfigOptions{
						{
							ConfigKey:  "AddConfig",
							ConfigFile: "./files/empty-conf.json",
						},
					},
				},
			},
			want: mockConfig{},
		},
		{
			name: "read config panic unmarshal duration",
			args: args{
				options: Options{
					ConfigFile: "./files/panic-conf.json",
					AdditionalConfig: []AdditionalConfigOptions{
						{
							ConfigKey:  "AddConfig",
							ConfigFile: "./files/conf-excel.json",
						},
					},
				},
			},
			wantPanic: true,
			panicMsg:  "fatal error found during unmarshaling config. err: 1 error(s) decoding:\n\n* error decoding 'Duration': time: invalid duration \"sadasfasfdas\"",
		},
		{
			name: "read additional config panic unmarshal duration",
			args: args{
				options: Options{
					ConfigFile: "./files/conf.json",
					AdditionalConfig: []AdditionalConfigOptions{
						{
							ConfigKey:  "AddConfig",
							ConfigFile: "./files/panic-conf.json",
						},
					},
				},
			},
			wantPanic: true,
			panicMsg:  "fatal error found during unmarshaling config. err: 1 error(s) decoding:\n\n* error decoding 'AddConfig.Duration': time: invalid duration \"sadasfasfdas\"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := Init(tt.args.options)

			cfg := mockConfig{}
			if tt.wantPanic {
				assert.PanicsWithError(t, tt.panicMsg, func() { reader.ReadConfig(&cfg) })
			} else {
				reader.ReadConfig(&cfg)
				assert.Equal(t, tt.want, cfg)
			}
		})
	}
}

func Test_configReader_AllSettings(t *testing.T) {
	type args struct {
		options Options
	}
	tests := []struct {
		name     string
		args     args
		mockFunc func()
		want     map[string]interface{}
	}{
		{
			name: "get all settings",
			args: args{
				options: Options{
					ConfigFile: "./files/conf.json",
					AdditionalConfig: []AdditionalConfigOptions{
						{
							ConfigKey:  "AddConfig",
							ConfigFile: "./files/conf-excel.json",
						},
					},
				},
			},
			want: map[string]interface{}{
				"string":   "Test String",
				"boolean":  "true",
				"duration": "1m",
				"object": map[string]interface{}{
					"string":   "Test String",
					"boolean":  "true",
					"duration": "1m",
				},
				"reference": map[string]interface{}{
					"string":   "Test String",
					"boolean":  "true",
					"duration": "1m",
				},
				"meta": map[string]interface{}{
					"version":     "dev",
					"environment": "dev",
				},
				"addconfig": map[string]interface{}{
					"duration": "1s",
				},
			},
		},
		{
			name: "get all settings with custom version",
			args: args{
				options: Options{
					ConfigFile: "./files/conf.json",
					AdditionalConfig: []AdditionalConfigOptions{
						{
							ConfigKey:  "AddConfig",
							ConfigFile: "./files/conf-excel.json",
						},
					},
				},
			},
			mockFunc: func() {
				os.Setenv("SERVICE_VERSION", "v1.0.0")
			},
			want: map[string]interface{}{
				"string":   "Test String",
				"boolean":  "true",
				"duration": "1m",
				"object": map[string]interface{}{
					"string":   "Test String",
					"boolean":  "true",
					"duration": "1m",
				},
				"reference": map[string]interface{}{
					"string":   "Test String",
					"boolean":  "true",
					"duration": "1m",
				},
				"meta": map[string]interface{}{
					"version":     "v1.0.0",
					"environment": "dev",
				},
				"addconfig": map[string]interface{}{
					"duration": "1s",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := Init(tt.args.options)

			if tt.mockFunc != nil {
				tt.mockFunc()
			}

			reader.ReadConfig(&mockConfig{})
			got := reader.AllSettings()
			assert.Equal(t, tt.want, got)
		})
	}
}
