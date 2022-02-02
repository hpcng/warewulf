package warewulfd

import (
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/pkg/errors"
)

func sendFile(w http.ResponseWriter, filename string, sendto string) error {
	fd, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fd.Close()

	FileHeader := make([]byte, 512)
	_, err = fd.Read(FileHeader)
	if err != nil {
		return errors.Wrap(err, "failed to read header")
	}
	FileContentType := http.DetectContentType(FileHeader)
	FileStat, _ := fd.Stat()
	FileSize := strconv.FormatInt(FileStat.Size(), 10)

	w.Header().Set("Content-Disposition", "attachment; filename=kernel")
	w.Header().Set("Content-Type", FileContentType)
	w.Header().Set("Content-Length", FileSize)

	_, err = fd.Seek(0, 0)
	if err != nil {
		return errors.Wrap(err, "failed to seek")
	}

	_, err = io.Copy(w, fd)
	if err != nil {
		return errors.Wrap(err, "failed to copy")
	}

	daemonLogf("SEND:  %15s: %s\n", sendto, filename)

	return nil
}
