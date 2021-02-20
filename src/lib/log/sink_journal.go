package log

import (
	"github.com/coreos/go-systemd/v22/journal"
)

type journalSink struct{}

func (j *journalSink) Write(p []byte) (int, error) {
	if err := journal.Print(journal.PriInfo, string(p)); err != nil {
		return 0, err
	}
	return 0, nil
}

func (j *journalSink) Sync() error {
	return nil
}

func (j *journalSink) Close() error {
	return nil
}
