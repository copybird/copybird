package lz4

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/pierrec/lz4"
	"gotest.tools/assert"
)

func TestCompressLZ4(t *testing.T) {
	type args struct {
		level int
		input string
	}

	tests := []struct {
		name    string
		args    args
		want    int32
		wantErr bool
	}{
		// TODO: Add more test cases.
		{
			name:    "invalid-level",
			args:    args{level: -2, input: "hello, world."},
			wantErr: true,
		},
		{
			name:    "valid-level",
			args:    args{level: 0, input: "hello, world."},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			comp := BackupCompressLz4{}

			rb := bytes.NewReader([]byte("hello, world."))
			wb := new(bytes.Buffer)

			assert.Assert(t, comp.GetConfig() != nil)
			assert.NilError(t, comp.InitPipe(wb, rb))
			err := comp.InitModule(&Config{Level: tt.args.level})
			if !tt.wantErr && err != nil {
				t.Errorf("Compress.BackupCompressLz4() result = %v, want result %v", err, tt.wantErr)
				return
			}

			err = comp.Run(context.TODO())
			if err != nil {
				panic(err)
			}

			var buff2 = new(bytes.Buffer)
			gr := lz4.NewReader(wb)
			_, err = io.Copy(buff2, gr)
			assert.Equal(t, err, nil)
			assert.Equal(t, buff2.String(), "hello, world.")

			assert.NilError(t, comp.Close())
		})
	}
}
