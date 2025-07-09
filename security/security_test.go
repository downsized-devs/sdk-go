package security

import (
	"context"
	"testing"

	"github.com/downsized-devs/sdk-go/logger"
	"github.com/stretchr/testify/assert"
)

func Test_security_Encrypt_Decrypt(t *testing.T) {
	type args struct {
		ctx        context.Context
		passphrase string
		timestamp  int64
		plaintext  string
	}
	tests := []struct {
		name      string
		args      args
		wantError bool
	}{
		{
			name: "success encrypted & decrypted text: Hello World",
			args: args{
				ctx:        context.Background(),
				passphrase: "testing123",
				timestamp:  12345,
				plaintext:  "Hello World",
			},
			wantError: false,
		},
		{
			name: "success encrypted & decrypted text: generic",
			args: args{
				ctx:        context.Background(),
				passphrase: "testing123",
				timestamp:  54321,
				plaintext:  "generic",
			},
			wantError: false,
		},
		{
			name: "failed encrypted & decrypted text: chipper is empty or invalid",
			args: args{
				ctx:        context.Background(),
				passphrase: "",
				timestamp:  0,
				plaintext:  "",
			},
			wantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &security{
				log: logger.Init(logger.Config{}),
			}

			enc := s.Encrypt(tt.args.ctx, tt.args.passphrase, tt.args.timestamp, tt.args.plaintext)
			dec, err := s.Decrypt(tt.args.ctx, tt.args.passphrase, tt.args.timestamp, enc)
			if !tt.wantError {
				assert.Nil(t, err)
				assert.NotEqual(t, enc, tt.args.plaintext)
				assert.Equal(t, dec, tt.args.plaintext)
			} else {
				// set chipper empty
				_, err := s.Decrypt(tt.args.ctx, tt.args.passphrase, tt.args.timestamp, "")
				assert.NotNil(t, err)
			}
		})
	}
}

func Test_security_HashPassword(t *testing.T) {
	type fields struct {
		log logger.Interface
	}
	type args struct {
		ctx       context.Context
		secretKey string
		password  string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "success",
			fields: fields{
				log: logger.Init(logger.Config{}),
			},
			args: args{
				ctx:       context.Background(),
				secretKey: "secret-key",
				password:  "password",
			},
			want: "93435c1a22a97eb31d56ee721b650d90cb17daf535b8f8d4727035e1cbf6821c",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &security{
				log: tt.fields.log,
			}
			if got := s.HashPassword(tt.args.ctx, tt.args.secretKey, tt.args.password); got != tt.want {
				t.Errorf("security.HashPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_security_CompareHashPassword(t *testing.T) {
	type fields struct {
		log logger.Interface
	}
	type args struct {
		ctx          context.Context
		secretKey    string
		hashPassword string
		password     string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "success",
			fields: fields{
				log: logger.Init(logger.Config{}),
			},
			args: args{
				ctx:          context.Background(),
				secretKey:    "secret-key",
				hashPassword: "93435c1a22a97eb31d56ee721b650d90cb17daf535b8f8d4727035e1cbf6821c",
				password:     "password",
			},
			want: true,
		},
		{
			name: "failed",
			fields: fields{
				log: logger.Init(logger.Config{}),
			},
			args: args{
				ctx:          context.Background(),
				secretKey:    "secret-key",
				hashPassword: "93435c1a22a97eb31d56ee721b650d90cb17daf535b8f8d4727035e1cbf6821c-failed",
				password:     "password",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &security{
				log: tt.fields.log,
			}
			if got := s.CompareHashPassword(tt.args.ctx, tt.args.secretKey, tt.args.hashPassword, tt.args.password); got != tt.want {
				t.Errorf("security.CompareHashPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_security_ScryptPassword(t *testing.T) {
	config := ScryptConfig{
		Base64SignerKey:     "MoaHZJjRSE9Ktj6HnIkoldV+BmXpD7YVboHgJOY4SDnUNiNMTUILxlsY4igO3Uzx/n/VwFju9IC4fQfgDy7LwQ==",
		Base64SaltSeperator: "Bw==",
		Rounds:              8,
		MemoryCost:          14,
	}
	type fields struct {
		log    logger.Interface
		scrypt ScryptConfig
	}
	type args struct {
		ctx      context.Context
		salt     string
		password string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "success scrypt my password",
			fields: fields{
				log:    logger.Init(logger.Config{}),
				scrypt: config,
			},
			args: args{
				ctx:      context.Background(),
				salt:     "ekr/rlgB6tovww==",
				password: "D0wn5izeDd3v5",
			},
			want: "8WTYvjqmK1naiAZnAFthrXzvdiW2SKFW6RWCsFe8bhhCr7PZ9EUr9WgZOQSYNMoSRujIlaNluzLl7u268P1FJQ==",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &security{
				log:    tt.fields.log,
				scrypt: tt.fields.scrypt,
			}
			if got := s.ScryptPassword(tt.args.ctx, tt.args.salt, tt.args.password); got != tt.want {
				t.Errorf("security.ScryptPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_security_CompareScryptPassword(t *testing.T) {
	config := ScryptConfig{
		Base64SignerKey:     "MoaHZJjRSE9Ktj6HnIkoldV+BmXpD7YVboHgJOY4SDnUNiNMTUILxlsY4igO3Uzx/n/VwFju9IC4fQfgDy7LwQ==",
		Base64SaltSeperator: "Bw==",
		Rounds:              8,
		MemoryCost:          14,
	}
	type fields struct {
		log    logger.Interface
		scrypt ScryptConfig
	}
	type args struct {
		ctx          context.Context
		passwordHash string
		salt         string
		password     string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "success compare scrypt password",
			fields: fields{
				log:    logger.Init(logger.Config{}),
				scrypt: config,
			},
			args: args{
				ctx:          context.Background(),
				passwordHash: "8WTYvjqmK1naiAZnAFthrXzvdiW2SKFW6RWCsFe8bhhCr7PZ9EUr9WgZOQSYNMoSRujIlaNluzLl7u268P1FJQ==",
				salt:         "ekr/rlgB6tovww==",
				password:     "D0wn5izeDd3v5",
			},
			want: true,
		},
		{
			name: "success compare scrypt password: test+3",
			fields: fields{
				log:    logger.Init(logger.Config{}),
				scrypt: config,
			},
			args: args{
				ctx:          context.Background(),
				passwordHash: "gpBOQ2f6nwfoip848UDPBxvkUNEn/laOQkH5CDnj6NnrHsG8QGoHpL+cIOANU4mJw+XqD0s1+l3sUEHMRjYrWQ==",
				salt:         "AlerfOEjlbSEbQ==",
				password:     "freedom09",
			},
			want: true,
		},
		{
			name: "failed compare scrypt password: wrong password",
			fields: fields{
				log:    logger.Init(logger.Config{}),
				scrypt: config,
			},
			args: args{
				ctx:          context.Background(),
				passwordHash: "8bVI07w7zmyTgnSU8fE2C6Nn/pzEiukZJEFkWSFch5zOxypmjRt67C0aikpfQ/3z5T2XF9vmJF5PkaskD9D1sw==",
				salt:         "ekr/rlgB6tovww==",
				password:     "D0wn5izeDd3v5",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &security{
				log:    tt.fields.log,
				scrypt: tt.fields.scrypt,
			}
			if got := s.CompareScryptPassword(tt.args.ctx, tt.args.passwordHash, tt.args.salt, tt.args.password); got != tt.want {
				t.Errorf("security.CompareScryptPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}
