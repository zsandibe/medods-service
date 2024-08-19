package service

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/zsandibe/medods-service/config"
	"github.com/zsandibe/medods-service/internal/domain"
	"github.com/zsandibe/medods-service/internal/entity"
	"github.com/zsandibe/medods-service/internal/service/mocks"
)

// Из-за того что я возвращаю в repo error,а в бизнес логике domain.TokenPair не получается написать правильно юнит тесты.
func Test_service_Create(t *testing.T) {
	type fields struct {
		conf *config.Config
		repo *mocks.MockRepository
	}
	type args struct {
		ctx  context.Context
		guid uuid.UUID
		ip   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    domain.TokenPair
		wantErr bool
		prepare func(args args, fields fields)
	}{
		{
			name: "successful creation",
			args: args{
				ctx:  context.Background(),
				guid: uuid.New(),
				ip:   "127.0.0.1",
			},
			want: domain.TokenPair{
				AccessToken:  "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDc0ODE5MDIsInNlc3Npb25faWQiOiIyZjU1YjQ2ZS1lZTEyLTQ1MDUtYTkxNC1lZjRlOWYyMjlkYWQiLCJndWlkIjoiMzg4ZDZiNDgtMWE3ZS00ZGY3LWIxZWItNzljYTcyNmZiODE0In0.mP2Vbzgtqv_0hnFt5V050Xnyqg9kB_zSL2KWkmY7NsMLSxqGhEd7B1P0K_jq3IG64rmYHQR1rFygoAmkl8fQgQ",
				RefreshToken: []byte("iSgH6v0USyOqIui1sLexZA=="),
			},
			wantErr: false,
			prepare: func(args args, fields fields) {
				session := entity.Session{
					Id:                 uuid.New(),
					Guid:               args.guid,
					Ip:                 args.ip,
					CreatedAt:          time.Now(),
					UpdatedAt:          time.Now(),
					HashedRefreshToken: []byte("JDJhJDEwJGZSUHdIWEN4UUI3L3pGWHVaZTVTLi5LOFVmNlYxTWp5Z1Q2alk2QXlGVWozZzNkZEN1QnRD"),
				}

				fields.repo.EXPECT().
					Create(args.ctx, gomock.AssignableToTypeOf(session)).
					Return(nil)
			},
		},
		{
			name: "creation fails due to repo error",
			args: args{
				ctx:  context.Background(),
				guid: uuid.New(),
				ip:   "127.0.0.1",
			},
			want:    domain.TokenPair{},
			wantErr: true,
			prepare: func(args args, fields fields) {
				fields.repo.EXPECT().
					Create(args.ctx, gomock.Any()).
					Return(errors.New("repository error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockRepository(ctrl)
			f := fields{
				conf: &config.Config{
					Token: config.TokenConfig{
						AccessTokenAge:  time.Hour,
						AccessKey:       "secret",
						RefreshTokenAge: time.Hour * 720,
					},
				},
				repo: mockRepo,
			}
			if tt.prepare != nil {
				tt.prepare(tt.args, f)
			}

			s := NewService(f.repo, f.conf)
			got, err := s.Create(tt.args.ctx, tt.args.guid, tt.args.ip)
			// fmt.Println("-----------------------")
			// fmt.Println(got.AccessToken)
			// fmt.Println(tt.want.AccessToken)
			// fmt.Println("++++++++++++++")
			// fmt.Println(string(got.RefreshToken))
			// fmt.Println("++++++++++++++")
			// fmt.Println(string(tt.want.RefreshToken))
			// fmt.Println("-----------------------")
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.AccessToken != tt.want.AccessToken {
				t.Errorf("Create() got = %v, want %v", got.AccessToken, tt.want.AccessToken)
			}
			if !bytes.Equal(got.RefreshToken, tt.want.RefreshToken) {
				t.Errorf("Create() got = %v, want %v", got.RefreshToken, tt.want.RefreshToken)
			}
		})
	}
}
