package id

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// New returns an application-generated ID suitable for SQLite TEXT primary keys.
// It is timestamp-prefixed for rough ordering and random-suffixed for uniqueness.
func New() (string, error) {
	var random [10]byte
	if _, err := rand.Read(random[:]); err != nil {
		return "", fmt.Errorf("read random bytes: %w", err)
	}

	return fmt.Sprintf("%d%s", time.Now().UTC().UnixMilli(), hex.EncodeToString(random[:])), nil
}
