package etcd

import (
	"testing"
	"time"

	"github.com/azrod/zr/pkg/db"
	"github.com/azrod/zr/pkg/mocks/tests"
	"github.com/stretchr/testify/assert"
)

func TestNewDBETCD(t *testing.T) {

	type args struct {
		cfg ConfigDBETCD
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				cfg: ConfigDBETCD{
					Endpoints:   []string{"127.0.0.1:2379"},
					DialTimeout: 5 * time.Second,
					TLS:         nil,
				},
			},
			wantErr: false,
		},
		{
			name: "test2",
			args: args{
				cfg: ConfigDBETCD{
					Endpoints:   nil,
					DialTimeout: 5 * time.Second,
					TLS:         nil,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewClient(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("newDBETCD() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestEtcdServerRWOLD(t *testing.T) {

	client := tests.SetupClientETCD()

	key := "myKey-1"
	value := db.DataHotReload{
		LogLevel:  "debug",
		LogFormat: "json",
	}

	err := client.Set(key, value, 10)
	assert.Nil(t, err)

	kv, err := client.Get(key, 10)
	assert.Nil(t, err)
	assert.NotNil(t, kv)
	assert.Equal(t, value, kv)
}
