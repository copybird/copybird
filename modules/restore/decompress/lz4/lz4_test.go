package lz4

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/pierrec/lz4"
	"github.com/stretchr/testify/assert"
)

func TestDecompress(t *testing.T) {
	wr := new(bytes.Buffer)
	var decompressor RestoreDecompressLz4

	s := `I bomb atomically, Socrates' philosophies
And hypotheses can't define how I be droppin' these
Mockeries, lyrically perform armed robbery
Flee with the lottery, possibly they spotted me
Battle-scarred shogun, explosion when my pen hits
Tremendous, ultra-violet shine blind forensics
I inspect view through the future see millennium
Killa Beez sold fifty gold sixty platinum
Shackling the masses with drastic rap tactics
Graphic displays melt the steel like blacksmiths
Black Wu jackets Queen Beez ease the guns in
Rumblin' patrolmen tear gas laced the function
Heads by the score take flight incite a war
Chicks hit the floor, die hard fans demand more
Behold the bold soldier, control the globe slowly
Proceeds to blow swingin' swords like Shinobi
Stomp grounds I pound footprints in solid rock
Wu got it locked, performin' live on your hottest block`
	hw, err := helper(s)
	assert.NoError(t, err)
	rd := bytes.NewReader(hw)
	assert.Equal(t, hw, nil)

	assert.NoError(t, InitModule(nil))
	assert.NoError(t, InitPipe(wr, rd))
	assert.NoError(t, Run())
	spew.Dump(wr)

	assert.Equal(t, wr.String(), "hello world")
}

func helper(s string) ([]byte, error) {
	var wr bytes.Buffer
	lw := lz4.NewWriter(&wr)
	lw.Header = lz4.Header{CompressionLevel: 4}
	defer lw.Close()
	r := strings.NewReader(s)
	spew.Dump(r)
	if _, err := io.Copy(lw, r); err != nil {
		return nil, err
	}
	spew.Dump(wr.Bytes())
	return wr.Bytes(), nil
}
