package filesystem

import (
	"bytes"
	"context"
	"errors"
	"io"
	"path/filepath"
	"strings"

	"github.com/mholt/archiver/v4"
	log "github.com/sirupsen/logrus"
	"github.com/trufflesecurity/trufflehog/v3/pkg/sources"
)

func (s *Source) scanArchive(ctx context.Context, file io.Reader, fileName string, path string, chunksChan chan *sources.Chunk) error {
	buf := &bytes.Buffer{}
	tee := io.TeeReader(file, buf)
	format, reader, err := archiver.Identify(fileName, tee)
	if err != nil {
		if errors.Is(err, archiver.ErrNoMatch) {
			//junk, _ := ioutil.ReadAll(reader)
			//log.Debugf("BILL2: %s", string(junk))
			return s.scanFile(buf, path, chunksChan)
		}
		log.WithError(err).Errorf("error identifying file")
		return err
	}
	handler := func(ctx context.Context, f archiver.File) error {
		fReader, err := f.Open()
		if err != nil {
			return err
		}
		return s.scanArchive(ctx, fReader, f.Name(), fileName, chunksChan)
	}
	switch archive := format.(type) {
	case archiver.Extractor:
		log.Debug("chunk is archived")
		return archive.Extract(ctx, reader, nil, handler)
	case archiver.Decompressor:
		compReader, err := archive.OpenReader(reader)
		log.Debug("chunk is compressed")
		if err != nil {
			return err
		}
		//junk, _ := ioutil.ReadAll(compReader)
		//log.Debugf("BILL1: %s", string(junk))
		ext := filepath.Ext(path)
		return s.scanArchive(ctx, compReader, strings.TrimSuffix(fileName, ext), path, chunksChan)
	default:
		log.Errorf("Unknown archive type: %s", format.Name())
	}
	return nil
}
