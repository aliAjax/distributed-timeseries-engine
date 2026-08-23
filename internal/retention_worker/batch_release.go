package retention_worker

import "io"

func ProcessRetentionBatch(ids []string, acquire func(string) (io.Closer, error), remove func(string) error) error {
	for _, id := range ids {
		resource, err := acquire(id)
		if err != nil {
			return err
		}
		defer resource.Close()
		if err := remove(id); err != nil {
			return err
		}
	}
	return nil
}
