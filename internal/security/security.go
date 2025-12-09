// Package security provides cryptographic primitives and secure memory handling
// for the glyphic password generator. All random operations use crypto/rand only.
package security

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"runtime"

	"golang.org/x/sys/unix"
)

var (
	// ErrInvalidRange indicates an invalid range for random number generation
	ErrInvalidRange = errors.New("invalid range: max must be positive")

	// ErrCryptoRandFailed indicates crypto/rand failed to generate random bytes
	ErrCryptoRandFailed = errors.New("crypto/rand failed")

	// ErrMemoryLockFailed indicates memory locking failed
	ErrMemoryLockFailed = errors.New("failed to lock memory")

	// ErrMemoryUnlockFailed indicates memory unlocking failed
	ErrMemoryUnlockFailed = errors.New("failed to unlock memory")
)

// SecureRandomIndex returns a cryptographically secure random index in the range [0, max).
// It uses crypto/rand and avoids modulo bias.
func SecureRandomIndex(max int) (int, error) {
	if max <= 0 {
		return 0, ErrInvalidRange
	}

	// For small ranges, use simple approach
	if max <= 256 {
		var buf [1]byte
		if _, err := rand.Read(buf[:]); err != nil {
			return 0, fmt.Errorf("%w: %v", ErrCryptoRandFailed, err)
		}
		return int(buf[0]) % max, nil
	}

	// For larger ranges, use 64-bit random value
	var buf [8]byte
	if _, err := rand.Read(buf[:]); err != nil {
		return 0, fmt.Errorf("%w: %v", ErrCryptoRandFailed, err)
	}

	n := binary.BigEndian.Uint64(buf[:])
	return int(n % uint64(max)), nil
}

// SecureRandomBytes fills the provided byte slice with cryptographically secure random bytes.
func SecureRandomBytes(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	if _, err := rand.Read(data); err != nil {
		return fmt.Errorf("%w: %v", ErrCryptoRandFailed, err)
	}
	return nil
}

// SecureRandomInt64 returns a cryptographically secure random int64.
func SecureRandomInt64() (int64, error) {
	var buf [8]byte
	if _, err := rand.Read(buf[:]); err != nil {
		return 0, fmt.Errorf("%w: %v", ErrCryptoRandFailed, err)
	}
	return int64(binary.BigEndian.Uint64(buf[:])), nil
}

// SecureRandomFloat returns a cryptographically secure random float64 in the range [0.0, 1.0).
func SecureRandomFloat() (float64, error) {
	n, err := SecureRandomInt64()
	if err != nil {
		return 0, err
	}
	// Convert to [0.0, 1.0) by treating as uint64 and dividing by max uint64
	return float64(uint64(n)) / float64(^uint64(0)), nil
}

// SecureZero zeroes the provided byte slice in a way that should not be optimized away.
// After zeroing, it uses runtime.KeepAlive to prevent the compiler from optimizing it out.
func SecureZero(data []byte) {
	for i := range data {
		data[i] = 0
	}
	runtime.KeepAlive(data)
}

// SecureBuffer represents a buffer with secure memory handling
type SecureBuffer struct {
	data   []byte
	locked bool
}

// NewSecureBuffer creates a new secure buffer of the specified size
func NewSecureBuffer(size int) (*SecureBuffer, error) {
	if size <= 0 {
		return nil, errors.New("buffer size must be positive")
	}

	buf := &SecureBuffer{
		data: make([]byte, size),
	}

	return buf, nil
}

// Lock locks the buffer's memory to prevent it from being swapped to disk
func (sb *SecureBuffer) Lock() error {
	if sb.locked {
		return nil // Already locked
	}

	if err := unix.Mlock(sb.data); err != nil {
		return fmt.Errorf("%w: %v", ErrMemoryLockFailed, err)
	}

	sb.locked = true
	return nil
}

// Unlock unlocks the buffer's memory
func (sb *SecureBuffer) Unlock() error {
	if !sb.locked {
		return nil // Already unlocked
	}

	if err := unix.Munlock(sb.data); err != nil {
		return fmt.Errorf("%w: %v", ErrMemoryUnlockFailed, err)
	}

	sb.locked = false
	return nil
}

// Bytes returns the underlying byte slice
func (sb *SecureBuffer) Bytes() []byte {
	return sb.data
}

// Zero securely zeros the buffer's contents
func (sb *SecureBuffer) Zero() {
	SecureZero(sb.data)
}

// Destroy zeros the buffer and unlocks its memory
func (sb *SecureBuffer) Destroy() error {
	sb.Zero()
	if sb.locked {
		return sb.Unlock()
	}
	return nil
}

// LockMemory locks the provided byte slice in memory to prevent swapping
func LockMemory(data []byte) error {
	if err := unix.Mlock(data); err != nil {
		return fmt.Errorf("%w: %v", ErrMemoryLockFailed, err)
	}
	return nil
}

// UnlockMemory unlocks the provided byte slice from memory
func UnlockMemory(data []byte) error {
	if err := unix.Munlock(data); err != nil {
		return fmt.Errorf("%w: %v", ErrMemoryUnlockFailed, err)
	}
	return nil
}

// ValidatePRNG performs a basic validation of the crypto/rand PRNG
func ValidatePRNG() error {
	// Generate two random values and ensure they're different
	var buf1, buf2 [32]byte

	if err := SecureRandomBytes(buf1[:]); err != nil {
		return fmt.Errorf("PRNG validation failed (first read): %w", err)
	}

	if err := SecureRandomBytes(buf2[:]); err != nil {
		return fmt.Errorf("PRNG validation failed (second read): %w", err)
	}

	// Check that the two values are different
	if string(buf1[:]) == string(buf2[:]) {
		return errors.New("PRNG validation failed: generated identical random values")
	}

	return nil
}

// ConstantTimeCompare performs constant-time comparison of two byte slices
// Returns true if they are equal, false otherwise
func ConstantTimeCompare(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	var diff byte
	for i := range a {
		diff |= a[i] ^ b[i]
	}

	return diff == 0
}
