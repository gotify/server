package auth

import (
	"github.com/gin-contrib/cors"
	"github.com/gotify/server/config"
	"github.com/gotify/server/mode"
	"reflect"
	"testing"
	"time"
)

func TestCorsConfig(t *testing.T) {
	mode.Set(mode.Prod)
	type args struct {
		conf *config.Configuration
	}
	serverConf := config.Configuration{}
	serverConf.Server.Cors.AllowOrigins = []string{"http://test.com"}
	serverConf.Server.Cors.AllowHeaders = []string{"content-type"}
	serverConf.Server.Cors.AllowMethods = []string{"GET"}
	tests := []struct {
		name string
		args args
		want cors.Config
	}{
		{
			name: "",
			args: args{
				conf: &serverConf,
			},
			want: cors.Config{
				AllowHeaders: []string{"content-type"},
				AllowMethods: []string{"GET"},
				MaxAge:       12 * time.Hour,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CorsConfig(tt.args.conf)
			if !reflect.DeepEqual(got.AllowHeaders, tt.want.AllowHeaders) {
				t.Errorf("CorsConfig() = %v, want %v", got.AllowHeaders, tt.want.AllowHeaders)
			}
			if !reflect.DeepEqual(got.AllowMethods, tt.want.AllowMethods) {
				t.Errorf("CorsConfig() = %v, want %v", got.AllowMethods, tt.want.AllowMethods)
			}
			if !reflect.DeepEqual(got.MaxAge, tt.want.MaxAge) {
				t.Errorf("CorsConfig() = %v, want %v", got.MaxAge, tt.want.MaxAge)
			}
			if !reflect.DeepEqual(got.AllowAllOrigins, tt.want.AllowAllOrigins) {
				t.Errorf("CorsConfig() = %v, want %v", got.AllowAllOrigins, tt.want.AllowAllOrigins)
			}
			if !reflect.DeepEqual(got.AllowCredentials, tt.want.AllowCredentials) {
				t.Errorf("CorsConfig() = %v, want %v", got.AllowCredentials, tt.want.AllowCredentials)
			}
			if !reflect.DeepEqual(got.ExposeHeaders, tt.want.ExposeHeaders) {
				t.Errorf("CorsConfig() = %v, want %v", got.ExposeHeaders, tt.want.ExposeHeaders)
			}
			if !reflect.DeepEqual(got.AllowWildcard, tt.want.AllowWildcard) {
				t.Errorf("CorsConfig() = %v, want %v", got.AllowWildcard, tt.want.AllowWildcard)
			}
			if got.AllowOriginFunc("http://test.com") != true {
				t.Errorf("CorsConfig() = AllowOriginFunc is false, want true")
			}
		})
	}
}

func TestDevCorsConfig(t *testing.T) {
	mode.Set(mode.Dev)
	type args struct {
		conf *config.Configuration
	}
	serverConf := config.Configuration{}
	tests := []struct {
		name string
		args args
		want cors.Config
	}{
		{
			name: "Dev config",
			args: args{
				conf: &serverConf,
			},
			want: cors.Config{
				AllowHeaders: []string{"X-Gotify-Key", "Authorization", "Content-Type", "Upgrade", "Origin",
					"Connection", "Accept-Encoding", "Accept-Language", "Host"},
				AllowMethods:    []string{"GET", "POST", "DELETE", "OPTIONS", "PUT"},
				MaxAge:          12 * time.Hour,
				AllowAllOrigins: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CorsConfig(tt.args.conf)
			if !reflect.DeepEqual(got.AllowHeaders, tt.want.AllowHeaders) {
				t.Errorf("CorsConfig() = %v, want %v", got.AllowHeaders, tt.want.AllowHeaders)
			}
			if !reflect.DeepEqual(got.AllowMethods, tt.want.AllowMethods) {
				t.Errorf("CorsConfig() = %v, want %v", got.AllowMethods, tt.want.AllowMethods)
			}
			if !reflect.DeepEqual(got.MaxAge, tt.want.MaxAge) {
				t.Errorf("CorsConfig() = %v, want %v", got.MaxAge, tt.want.MaxAge)
			}
			if !reflect.DeepEqual(got.AllowAllOrigins, tt.want.AllowAllOrigins) {
				t.Errorf("CorsConfig() = %v, want %v", got.AllowAllOrigins, tt.want.AllowAllOrigins)
			}
			if !reflect.DeepEqual(got.AllowCredentials, tt.want.AllowCredentials) {
				t.Errorf("CorsConfig() = %v, want %v", got.AllowCredentials, tt.want.AllowCredentials)
			}
			if !reflect.DeepEqual(got.ExposeHeaders, tt.want.ExposeHeaders) {
				t.Errorf("CorsConfig() = %v, want %v", got.ExposeHeaders, tt.want.ExposeHeaders)
			}
			if !reflect.DeepEqual(got.AllowWildcard, tt.want.AllowWildcard) {
				t.Errorf("CorsConfig() = %v, want %v", got.AllowWildcard, tt.want.AllowWildcard)
			}
			if got.AllowOriginFunc != nil {
				t.Errorf("CorsConfig() = AllowOriginFunc is not nil, want nil")
			}
		})
	}
}
