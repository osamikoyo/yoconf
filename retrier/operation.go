package retrier

import "time"

type Operation func() error

func Try(count int, opr Operation) error {
	var err error

	timeout := 1 * time.Second

	for range count {
		err = opr()
		if err == nil {
			return nil
		}

		time.Sleep(timeout)

		timeout *= 2
	}

	return err
}
