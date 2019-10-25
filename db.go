package redix

import (
	"errors"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

var (
	db      redis.Conn
	timeout = 1 * time.Second
)

var (
	errUnhealthyConn    = errors.New("Redis connection does not seem to be healthy")
	errUnexpectedOutput = errors.New("Unexpected output")
)

func connect(redisURL string) error {
	var err error

	db, err = redis.DialTimeout("tcp", "0.0.0.0:6379", timeout, timeout, timeout)
	if err != nil {
		return fmt.Errorf("Could not dial redis server with error: %w", err)
	}

	msg := "This is Ripley, last survivor of the Nostromo, signing off"

	reply, err := db.Do("ping", msg)
	if err != nil {
		return fmt.Errorf("Could not ping redis server with error: %w", err)
	}

	if v, ok := reply.(string); ok {
		if v != msg {
			return errUnhealthyConn
		}
	} else {
		return errUnhealthyConn
	}

	return nil
}

func command(command string, args ...string) (interface{}, error) {
	reply, err := db.Do(command, args)
	if err != nil {
		return nil, fmt.Errorf("Could not send command with error: %w", err)
	}

	return reply, nil
}

func close() error {
	return db.Close()
}
