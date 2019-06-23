package lz4

import (
	"bytes"
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
			comp := CompressLZ4{}

			rb := bytes.NewReader([]byte("hello, world."))
			wb := new(bytes.Buffer)

			assert.Assert(t, GetConfig() != nil)
			assert.NilError(t, InitPipe(wb, rb))
			err := InitModule(&Config{Level: tt.args.level})
			if !tt.wantErr && err != nil {
				t.Errorf("Compress.CompressLZ4() result = %v, want result %v", err, tt.wantErr)
				return
			}

			err = Run()
			if err != nil {
				panic(err)
			}

			var buff2 = new(bytes.Buffer)
			gr := lz4.NewReader(wb)
			_, err = io.Copy(buff2, gr)
			assert.Equal(t, err, nil)
			assert.Equal(t, buff2.String(), "hello, world.")

			assert.NilError(t, Close())
		})
	}
}
