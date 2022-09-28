package format

import (
	"io"
	"os"
	"testing"
)

// TestParseLogFormat tests the ParseLogFormat function
func TestParseLogFormat(t *testing.T) {
	type args struct {
		v string
	}
	tests := []struct {
		name    string
		args    args
		want    LogFormat
		wantErr bool
	}{
		{
			name: "json",
			args: args{
				v: "json",
			},
			want:    LogFormatJson,
			wantErr: false,
		},
		{
			name: "human",
			args: args{
				v: "human",
			},
			want:    LogFormatHuman,
			wantErr: false,
		},
		{
			name: "unknown",
			args: args{
				v: "unknown",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Format{}
			got, err := f.ParseLogFormat(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseLogFormat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseLogFormat() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test SetFormat
func TestFormat_SetFormat(t *testing.T) {
	type fields struct {
		LogFormat LogFormat
	}
	type args struct {
		format LogFormat
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "json",
			fields: fields{
				LogFormat: LogFormatJson,
			},
			args: args{
				format: LogFormatJson,
			},
			wantErr: false,
		},
		{
			name: "human",
			fields: fields{
				LogFormat: LogFormatHuman,
			},
			args: args{
				format: LogFormatHuman,
			},
			wantErr: false,
		},
		{
			name: "unknown",
			fields: fields{
				LogFormat: "",
			},
			args: args{
				format: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := Setup()
			if err := f.SetFormat(tt.args.format); (err != nil) != tt.wantErr {
				t.Errorf("SetFormat() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Test CustomOutput
func TestFormat_CustomOutput(t *testing.T) {
	type fields struct {
		LogFormat LogFormat
	}
	type args struct {
		output io.Writer
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "json",
			fields: fields{
				LogFormat: LogFormatJson,
			},
			args: args{
				output: os.Stdout,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := Setup()
			if err := f.SetOptions(CustomOutput(tt.args.output)); (err != nil) != tt.wantErr {
				t.Errorf("CustomOutput() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}

// Test GetFormat
func TestFormat_GetFormat(t *testing.T) {
	type fields struct {
		LogFormat LogFormat
	}
	tests := []struct {
		name   string
		fields fields
		want   LogFormat
	}{
		{
			name: "json",
			fields: fields{
				LogFormat: LogFormatJson,
			},
			want: LogFormatJson,
		},
		{
			name: "human",
			fields: fields{
				LogFormat: LogFormatHuman,
			},
			want: LogFormatHuman,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := Setup()
			err := f.SetFormat(tt.fields.LogFormat)
			if err != nil {
				t.Errorf("SetFormat() error = %v", err)
			}
			if got := f.GetFormat(); got != tt.want {
				t.Errorf("GetFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test String
func TestFormat_String(t *testing.T) {
	type fields struct {
		LogFormat LogFormat
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "json",
			fields: fields{
				LogFormat: LogFormatJson,
			},
			want: "json",
		},
		{
			name: "human",
			fields: fields{
				LogFormat: LogFormatHuman,
			},
			want: "human",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := Setup()
			err := f.SetFormat(tt.fields.LogFormat)
			if err != nil {
				t.Errorf("SetFormat() error = %v", err)
			}
			if got := f.format.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
