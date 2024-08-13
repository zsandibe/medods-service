package repository

import (
	"context"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/zsandibe/medods-service/internal/entity"
)

func Test_repositoryPostgres_Create(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		ctx     context.Context
		session entity.Session
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &repositoryPostgres{
				db: tt.fields.db,
			}
			if err := r.Create(tt.args.ctx, tt.args.session); (err != nil) != tt.wantErr {
				t.Errorf("repositoryPostgres.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_repositoryPostgres_GetSessionById(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		ctx       context.Context
		sessionId uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.Session
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &repositoryPostgres{
				db: tt.fields.db,
			}
			got, err := r.GetSessionById(tt.args.ctx, tt.args.sessionId)
			if (err != nil) != tt.wantErr {
				t.Errorf("repositoryPostgres.GetSessionById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("repositoryPostgres.GetSessionById() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_repositoryPostgres_DeleteSessionById(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		ctx       context.Context
		sessionId uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &repositoryPostgres{
				db: tt.fields.db,
			}
			if err := r.DeleteSessionById(tt.args.ctx, tt.args.sessionId); (err != nil) != tt.wantErr {
				t.Errorf("repositoryPostgres.DeleteSessionById() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_repositoryPostgres_Update(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		ctx     context.Context
		session entity.Session
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &repositoryPostgres{
				db: tt.fields.db,
			}
			if err := r.Update(tt.args.ctx, tt.args.session); (err != nil) != tt.wantErr {
				t.Errorf("repositoryPostgres.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
