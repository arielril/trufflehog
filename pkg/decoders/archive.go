package decoders

import (
	"bytes"

	"github.com/mholt/archiver/v4"
	"github.com/sirupsen/logrus"
	"github.com/trufflesecurity/trufflehog/v3/pkg/sources"
)

// Ensure the Decoder satisfies the interface at compile time
var _ Decoder = (*Archive)(nil)

type Archive struct{}

func (d *Archive) FromChunk(chunk *sources.Chunk) *sources.Chunk {
	logrus.Debug("ARCHIVE")

	chunkReader := bytes.NewReader(chunk.Data)
	format, _, err := archiver.Identify("test.file", chunkReader)
	if err != nil {
		logrus.WithError(err).Errorf("error identifying file")
		return chunk
	}
	// switch archive := format.(type) {
	switch format.(type) {
	case archiver.Extractor:
		logrus.Debug("chunk is archived")
	case archiver.Decompressor:
		logrus.Debug("chunk is compressed")
	default:
		logrus.Debugf("%s", format.Name())
	}
	return chunk
}
