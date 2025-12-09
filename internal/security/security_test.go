package security

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSecureRandomIndex(t *testing.T) {
	tests := []struct {
		name    string
		max     int
		wantErr bool
	}{
		{"valid range", 100, false},
		{"single element", 1, false},
		{"large range", 10000, false},
		{"zero max", 0, true},
		{"negative max", -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				_, err := SecureRandomIndex(tt.max)
				assert.Error(t, err)
				return
			}

			// Generate multiple random indices
			seen := make(map[int]bool)
			for range 100 {
				idx, err := SecureRandomIndex(tt.max)
				assert.NoError(t, err)
				assert.GreaterOrEqual(t, idx, 0)
				assert.Less(t, idx, tt.max)
				seen[idx] = true
			}

			// Should have some variety (not all the same)
			if tt.max > 1 {
				assert.Greater(t, len(seen), 1, "should generate varied indices")
			}
		})
	}
}

func TestSecureRandomBytes(t *testing.T) {
	tests := []struct {
		name string
		size int
	}{
		{"small buffer", 8},
		{"medium buffer", 64},
		{"large buffer", 1024},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := make([]byte, tt.size)
			err := SecureRandomBytes(data)
			assert.NoError(t, err)

			// Should not be all zeros
			allZero := true
			for _, b := range data {
				if b != 0 {
					allZero = false
					break
				}
			}
			assert.False(t, allZero, "random bytes should not be all zeros")
		})
	}
}

func TestSecureRandomInt64(t *testing.T) {
	// Generate multiple random int64s
	seen := make(map[int64]bool)
	for range 100 {
		n, err := SecureRandomInt64()
		assert.NoError(t, err)
		seen[n] = true
	}

	// Should have variety
	assert.Greater(t, len(seen), 90, "should generate varied int64 values")
}

func TestSecureRandomFloat(t *testing.T) {
	// Generate multiple random floats
	for range 100 {
		f, err := SecureRandomFloat()
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, f, 0.0)
		assert.Less(t, f, 1.0)
	}
}

func TestSecureZero(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{"small buffer", []byte{1, 2, 3, 4, 5}},
		{"large buffer", bytes.Repeat([]byte{0xFF}, 1024)},
		{"empty buffer", []byte{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := make([]byte, len(tt.data))
			copy(data, tt.data)

			SecureZero(data)

			// Should be all zeros
			for i, b := range data {
				assert.Equal(t, byte(0), b, "byte at index %d should be zero", i)
			}
		})
	}
}

func TestSecureBuffer(t *testing.T) {
	t.Run("create and destroy", func(t *testing.T) {
		buf, err := NewSecureBuffer(64)
		require.NoError(t, err)
		require.NotNil(t, buf)

		// Write some data
		copy(buf.data, []byte("test password"))

		// Destroy should zero the data
		buf.Destroy()

		// Verify all zeros
		for i, b := range buf.data {
			assert.Equal(t, byte(0), b, "byte at index %d should be zero after destroy", i)
		}
	})

	t.Run("lock and unlock", func(t *testing.T) {
		buf, err := NewSecureBuffer(64)
		require.NoError(t, err)
		defer buf.Destroy()

		// Lock memory
		err = buf.Lock()
		// Memory locking might fail in some environments (e.g., Docker, insufficient permissions)
		// So we don't assert NoError, but we check it doesn't panic
		if err != nil {
			t.Logf("Memory locking not available: %v", err)
		}

		// Unlock memory
		err = buf.Unlock()
		if err != nil {
			t.Logf("Memory unlocking not available: %v", err)
		}
	})

	t.Run("zero", func(t *testing.T) {
		buf, err := NewSecureBuffer(64)
		require.NoError(t, err)
		defer buf.Destroy()

		// Write some data
		copy(buf.data, []byte("sensitive data"))

		// Zero
		buf.Zero()

		// Verify all zeros
		for i, b := range buf.data {
			assert.Equal(t, byte(0), b, "byte at index %d should be zero", i)
		}
	})
}

func TestLockMemory(t *testing.T) {
	data := make([]byte, 64)
	err := LockMemory(data)
	if err != nil {
		t.Logf("Memory locking not available: %v", err)
		t.Skip("Skipping test - memory locking not supported in this environment")
	}

	// Cleanup
	_ = UnlockMemory(data)
}

func TestUnlockMemory(t *testing.T) {
	data := make([]byte, 64)
	err := LockMemory(data)
	if err != nil {
		t.Logf("Memory locking not available: %v", err)
		t.Skip("Skipping test - memory locking not supported in this environment")
	}

	err = UnlockMemory(data)
	assert.NoError(t, err)
}

func TestValidatePRNG(t *testing.T) {
	// crypto/rand should be working
	err := ValidatePRNG()
	assert.NoError(t, err)
}

func TestConstantTimeCompare(t *testing.T) {
	tests := []struct {
		name string
		a    []byte
		b    []byte
		want bool
	}{
		{"equal slices", []byte("hello"), []byte("hello"), true},
		{"different slices", []byte("hello"), []byte("world"), false},
		{"different lengths", []byte("hello"), []byte("hi"), false},
		{"empty slices", []byte{}, []byte{}, true},
		{"one empty", []byte("hello"), []byte{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConstantTimeCompare(tt.a, tt.b)
			assert.Equal(t, tt.want, got)
		})
	}
}

// Benchmarks
func BenchmarkSecureRandomIndex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = SecureRandomIndex(10000)
	}
}

func BenchmarkSecureRandomBytes(b *testing.B) {
	data := make([]byte, 64)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = SecureRandomBytes(data)
	}
}

func BenchmarkSecureRandomInt64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = SecureRandomInt64()
	}
}

func BenchmarkSecureRandomFloat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = SecureRandomFloat()
	}
}

func BenchmarkSecureZero(b *testing.B) {
	data := make([]byte, 1024)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SecureZero(data)
	}
}

func BenchmarkConstantTimeCompare(b *testing.B) {
	a := bytes.Repeat([]byte("test"), 16)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ConstantTimeCompare(a, a)
	}
}
