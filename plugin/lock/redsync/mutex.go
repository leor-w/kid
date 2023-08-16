package redsync

import (
	"context"
	"time"

	"github.com/hashicorp/go-multierror"

	"github.com/go-redsync/redsync/v4/redis"
)

type DelayFunc func(tries int) time.Duration

type RedsyncMutex struct {
	name          string
	genValueFunc  func() (string, error)
	delayFunc     DelayFunc
	tries         int
	expiry        time.Duration
	quorum        int
	timeoutFactor float64
	driftFactor   float64
	value         string
	until         time.Time
	pools         []redis.Pool
}

func (m *RedsyncMutex) Lock(ctx context.Context, key string, expiry time.Duration) error {

	if ctx == nil {
		ctx = context.Background()
	}

	value, err := m.genValueFunc()
	if err != nil {
		return err
	}

	for i := 0; i < m.tries; i++ {
		if i != 0 {
			select {
			case <-ctx.Done():
				return ErrFailed
			case <-time.After(m.delayFunc(i)):
				// Fall-through when the delay timer completes.
			}
		}

		start := time.Now()

		n, err := func() (int, error) {
			ctx, cancel := context.WithTimeout(ctx, time.Duration(int64(float64(expiry)*m.timeoutFactor)))
			defer cancel()
			return m.actOnPoolsAsync(func(pool redis.Pool) (bool, error) {
				return m.acquire(ctx, pool, value)
			})
		}()

		now := time.Now()
		until := now.Add(m.expiry - now.Sub(start) - time.Duration(int64(float64(m.expiry)*m.driftFactor)))
		if n >= m.quorum && now.Before(until) {
			m.value = value
			m.until = until
			return nil
		}
		func() (int, error) {
			ctx, cancel := context.WithTimeout(ctx, time.Duration(int64(float64(m.expiry)*m.timeoutFactor)))
			defer cancel()
			return m.actOnPoolsAsync(func(pool redis.Pool) (bool, error) {
				return m.release(ctx, pool, value)
			})
		}()
		if i == m.tries-1 && err != nil {
			return err
		}
	}

	return ErrFailed
}

func (m *RedsyncMutex) actOnPoolsAsync(actFn func(redis.Pool) (bool, error)) (int, error) {
	type result struct {
		Node   int
		Status bool
		Err    error
	}

	ch := make(chan result)
	for node, pool := range m.pools {
		go func(node int, pool redis.Pool) {
			r := result{Node: node}
			r.Status, r.Err = actFn(pool)
			ch <- r
		}(node, pool)
	}
	n := 0
	var taken []int
	var err error
	for range m.pools {
		r := <-ch
		if r.Status {
			n++
		} else if r.Err != nil {
			err = multierror.Append(err, &RedisError{Node: r.Node, Err: r.Err})
		} else {
			taken = append(taken, r.Node)
			err = multierror.Append(err, &ErrNodeTaken{Node: r.Node})
		}
	}

	if len(taken) >= m.quorum {
		return n, &ErrTaken{Nodes: taken}
	}
	return n, err
}

func (m *RedsyncMutex) acquire(ctx context.Context, pool redis.Pool, value string) (bool, error) {
	conn, err := pool.Get(ctx)
	if err != nil {
		return false, err
	}
	defer conn.Close()
	reply, err := conn.SetNX(m.name, value, m.expiry)
	if err != nil {
		return false, err
	}
	return reply, nil
}

var deleteScript = redis.NewScript(1, `
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		return redis.call("DEL", KEYS[1])
	else
		return 0
	end
`)

func (m *RedsyncMutex) release(ctx context.Context, pool redis.Pool, value string) (bool, error) {
	conn, err := pool.Get(ctx)
	if err != nil {
		return false, err
	}
	defer conn.Close()
	status, err := conn.Eval(deleteScript, m.name, value)
	if err != nil {
		return false, err
	}
	return status != int64(0), nil
}

func (m *RedsyncMutex) Unlock(ctx context.Context, key string) (bool, error) {
	n, err := m.actOnPoolsAsync(func(pool redis.Pool) (bool, error) {
		return m.release(ctx, pool, m.value)
	})
	if n < m.quorum {
		return false, err
	}
	return true, nil
}

func (m *RedsyncMutex) Extend() (bool, error) {
	return m.ExtendContext(nil)
}

// ExtendContext resets the mutex's expiry and returns the status of expiry extension.
func (m *RedsyncMutex) ExtendContext(ctx context.Context) (bool, error) {
	start := time.Now()
	n, err := m.actOnPoolsAsync(func(pool redis.Pool) (bool, error) {
		return m.touch(ctx, pool, m.value, int(m.expiry/time.Millisecond))
	})
	if n < m.quorum {
		return false, err
	}
	now := time.Now()
	until := now.Add(m.expiry - now.Sub(start) - time.Duration(int64(float64(m.expiry)*m.driftFactor)))
	if now.Before(until) {
		m.until = until
		return true, nil
	}
	return false, ErrExtendFailed
}

var touchScript = redis.NewScript(1, `
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		return redis.call("PEXPIRE", KEYS[1], ARGV[2])
	else
		return 0
	end
`)

func (m *RedsyncMutex) touch(ctx context.Context, pool redis.Pool, value string, expiry int) (bool, error) {
	conn, err := pool.Get(ctx)
	if err != nil {
		return false, err
	}
	defer conn.Close()
	status, err := conn.Eval(touchScript, m.name, value, expiry)
	if err != nil {
		return false, err
	}
	return status != int64(0), nil
}

func (m *RedsyncMutex) Valid() (bool, error) {
	return m.ValidContext(nil)
}

// ValidContext returns true if the lock acquired through m is still valid. It may
// also return true erroneously if quorum is achieved during the call and at
// least one node then takes long enough to respond for the lock to expire.
//
// Deprecated: Use Until instead. See https://github.com/go-redsync/redsync/issues/72.
func (m *RedsyncMutex) ValidContext(ctx context.Context) (bool, error) {
	n, err := m.actOnPoolsAsync(func(pool redis.Pool) (bool, error) {
		return m.valid(ctx, pool)
	})
	return n >= m.quorum, err
}

func (m *RedsyncMutex) valid(ctx context.Context, pool redis.Pool) (bool, error) {
	if m.value == "" {
		return false, nil
	}
	conn, err := pool.Get(ctx)
	if err != nil {
		return false, err
	}
	defer conn.Close()
	reply, err := conn.Get(m.name)
	if err != nil {
		return false, err
	}
	return m.value == reply, nil
}
