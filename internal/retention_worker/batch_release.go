package retention_worker

import "io"

func ProcessRetentionBatch(ids []string, acquire func(string) (io.Closer, error), remove func(string) error) error {
	for _, id := range ids {
		if err := func() error {
			resource, err := acquire(id)
			if err != nil {
				return err
			}
			defer resource.Close()
			return remove(id)
		}(); err != nil {
			return err
		}
	}
	return nil
}
