package level

import (
	"testing"

	"github.com/rs/zerolog"
)

// Test ParseLogLevel
func TestParseLogLevel(t *testing.T) {
	tests := []struct {
		name    string
		v       string
		want    LogLevel
		wantErr bool
	}{
		{
			name:    "debug",
			v:       "debug",
			want:    LogLevel(zerolog.DebugLevel),
			wantErr: false,
		},
		{
			name:    "info",
			v:       "info",
			want:    LogLevel(zerolog.InfoLevel),
			wantErr: false,
		},
		{
			name:    "warn",
			v:       "warn",
			want:    LogLevel(zerolog.WarnLevel),
			wantErr: false,
		},
		{
			name:    "error",
			v:       "error",
			want:    LogLevel(zerolog.ErrorLevel),
			wantErr: false,
		},
		{
			name:    "fatal",
			v:       "fatal",
			want:    LogLevel(zerolog.FatalLevel),
			wantErr: false,
		},
		{
			name:    "panic",
			v:       "panic",
			want:    LogLevel(zerolog.PanicLevel),
			wantErr: false,
		},
		{
			name:    "unknown",
			v:       "unknown",
			want:    LogLevel(0),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseLogLevel(tt.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseLogLevel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseLogLevel() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test SetLevel
func TestLevel_SetLevel(t *testing.T) {
	type fields struct {
		LogLevel LogLevel
	}
	type args struct {
		level LogLevel
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "debug",
			fields: fields{
				LogLevel: LogLevel(zerolog.DebugLevel),
			},
			args: args{
				level: LogLevel(zerolog.DebugLevel),
			},
			wantErr: false,
		},
		{
			name: "info",
			fields: fields{
				LogLevel: LogLevel(zerolog.InfoLevel),
			},
			args: args{
				level: LogLevel(zerolog.InfoLevel),
			},
			wantErr: false,
		},
		{
			name: "warn",
			fields: fields{
				LogLevel: LogLevel(zerolog.WarnLevel),
			},
			args: args{
				level: LogLevel(zerolog.WarnLevel),
			},
			wantErr: false,
		},
		{
			name: "error",
			fields: fields{
				LogLevel: LogLevel(zerolog.ErrorLevel),
			},
			args: args{
				level: LogLevel(zerolog.ErrorLevel),
			},
			wantErr: false,
		},
		{
			name: "fatal",
			fields: fields{
				LogLevel: LogLevel(zerolog.FatalLevel),
			},
			args: args{
				level: LogLevel(zerolog.FatalLevel),
			},
			wantErr: false,
		},
		{
			name: "panic",
			fields: fields{
				LogLevel: LogLevel(zerolog.PanicLevel),
			},
			args: args{
				level: LogLevel(zerolog.PanicLevel),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Level{
				LogLevel: tt.fields.LogLevel,
			}
			if err := l.SetLevel(tt.args.level); (err != nil) != tt.wantErr {
				t.Errorf("Level.SetLevel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Test GetLevel
func TestLevel_GetLevel(t *testing.T) {
	type fields struct {
		LogLevel LogLevel
	}
	tests := []struct {
		name   string
		fields fields
		want   LogLevel
	}{
		{
			name: "debug",
			fields: fields{
				LogLevel: LogLevel(zerolog.DebugLevel),
			},
			want: LogLevel(zerolog.DebugLevel),
		},
		{
			name: "info",
			fields: fields{
				LogLevel: LogLevel(zerolog.InfoLevel),
			},
			want: LogLevel(zerolog.InfoLevel),
		},
		{
			name: "warn",
			fields: fields{
				LogLevel: LogLevel(zerolog.WarnLevel),
			},
			want: LogLevel(zerolog.WarnLevel),
		},
		{
			name: "error",
			fields: fields{
				LogLevel: LogLevel(zerolog.ErrorLevel),
			},
			want: LogLevel(zerolog.ErrorLevel),
		},
		{
			name: "fatal",
			fields: fields{
				LogLevel: LogLevel(zerolog.FatalLevel),
			},
			want: LogLevel(zerolog.FatalLevel),
		},
		{
			name: "panic",
			fields: fields{
				LogLevel: LogLevel(zerolog.PanicLevel),
			},
			want: LogLevel(zerolog.PanicLevel),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Level{
				LogLevel: tt.fields.LogLevel,
			}
			if got := l.GetLevel(); got != tt.want {
				t.Errorf("Level.GetLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test Setup
func TestSetup(t *testing.T) {

	l := Setup()

	if l == nil {
		t.Errorf("Setup() = %v, want %v", l, nil)
	}

}

// Test String
func TestLogLevel_String(t *testing.T) {
	type fields struct {
		LogLevel LogLevel
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "debug",
			fields: fields{
				LogLevel: LogLevel(zerolog.DebugLevel),
			},
			want: "debug",
		},
		{
			name: "info",
			fields: fields{
				LogLevel: LogLevel(zerolog.InfoLevel),
			},
			want: "info",
		},
		{
			name: "warn",
			fields: fields{
				LogLevel: LogLevel(zerolog.WarnLevel),
			},
			want: "warn",
		},
		{
			name: "error",
			fields: fields{
				LogLevel: LogLevel(zerolog.ErrorLevel),
			},
			want: "error",
		},
		{
			name: "fatal",
			fields: fields{
				LogLevel: LogLevel(zerolog.FatalLevel),
			},
			want: "fatal",
		},
		{
			name: "panic",
			fields: fields{
				LogLevel: LogLevel(zerolog.PanicLevel),
			},
			want: "panic",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := Setup()
			l.SetLevel(tt.fields.LogLevel)
			if got := l.GetLevel().String(); got != tt.want {
				t.Errorf("LogLevel.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
