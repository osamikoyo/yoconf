package retrier

import "time"

func Connect[T any](retriers uint, try func() (T, error)) (T, error) {
	timeout := time.Second

	var (
		value T
		err   error
	)

	for range retriers {
		value, err = try()
		if err == nil {
			return value, nil
		}

		time.Sleep(timeout)

		timeout *= 2
	}

	return value, err
}
