package oauth2

import (
	"context"
	"gitee.com/unitedrhino/share/conf"
	"golang.org/x/oauth2"
	"reflect"
	"testing"
)

func TestGithubClient_GetAuthCodeURL(t *testing.T) {
	type fields struct {
		config *oauth2.Config
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GithubClient{
				config: tt.fields.config,
			}
			if got := g.GetAuthCodeURL(); got != tt.want {
				t.Errorf("GetAuthCodeURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGithubClient_GetUserInfo(t *testing.T) {
	type fields struct {
		config *oauth2.Config
	}
	type args struct {
		ctx  context.Context
		code string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *GitHubUser
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GithubClient{
				config: tt.fields.config,
			}
			got, err := g.GetUserInfo(tt.args.ctx, tt.args.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUserInfo() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewGithub(t *testing.T) {
	ctx := context.Background()
	c := NewGithub(ctx, &conf.ThirdConf{
		AppID:     "Ov23liEdUgBqBwlW77Au",
		AppSecret: "4a89ff39187f83bdd68a75f34c89e04071567c9c",
	})
	url := c.GetAuthCodeURL()
	t.Log(url)
	var code string
	c.GetUserInfo(ctx, code)
}
