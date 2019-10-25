package redix

import (
	"fmt"
	"io/ioutil"

	toml "github.com/pelletier/go-toml"
)

type Record struct {
	Key   string
	Value interface{}
}

// Load - Connects to the Redis server using the given address,
// then tries to load that database up with fixtures from the fixtures file at the path provided
func Load(address, fixturePath string) error {
	err := connect(address)
	if err != nil {
		return fmt.Errorf("Error encountered while trying to connect to provided redis server: %s, %w", address, err)
	}

	return loadFixtures(fixturePath)
}

// Teardown - should be called after using the test fixtures, flushes the redis dbs
func Teardown() error {
	// FLUSHALL
	output, err := command("FLUSHALL")
	if err != nil {
		return fmt.Errorf("Encountered error while trying to tear down fixture: %w", err)
	}

	if v, ok := output.(string); ok {
		if v != "OK" {
			return errUnexpectedOutput
		}
	}

	return close()
}

func loadFixtures(fixturePath string) error {
	// parse toml fixture file at path and add the records to the db
	b, err := ioutil.ReadFile(fixturePath) // TODO: change this to a proper os.Open and Read
	if err != nil {
		return fmt.Errorf("Encountered error while trying to open file")
	}

	var rec Record
	return toml.Unmarshal(b, &rec)
	// also add records to redis from each of these. Might have to do a bit of surgery to read in unexpected TOML tuples
}
