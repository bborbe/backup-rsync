package archiver

import "context"

type backupArchiver struct {
}

func New() *backupArchiver {
	return new(backupArchiver)
}

func (b *backupArchiver) Archiv(ctx context.Context) error {

	return nil
}
