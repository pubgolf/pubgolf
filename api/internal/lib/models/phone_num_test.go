package models

import (
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPhoneNum(t *testing.T) {
	t.Parallel()

	t.Run("parses valid number", func(t *testing.T) {
		t.Parallel()

		testNum := faker.E164PhoneNumber()
		num, err := NewPhoneNum(testNum)

		require.NoError(t, err)
		assert.Equal(t, testNum, num.String())
	})

	t.Run("rejects invalid number", func(t *testing.T) {
		t.Parallel()

		t.Run("empty string", func(t *testing.T) {
			t.Parallel()

			_, err := NewPhoneNum("")
			require.ErrorIs(t, err, ErrInvalidPhoneNumFormat)
		})

		t.Run("missing plus sign", func(t *testing.T) {
			t.Parallel()

			_, err := NewPhoneNum("5551231234")
			require.ErrorIs(t, err, ErrInvalidPhoneNumFormat)
		})

		t.Run("too many digits", func(t *testing.T) {
			t.Parallel()

			_, err := NewPhoneNum("+12345678901234567890")
			require.ErrorIs(t, err, ErrInvalidPhoneNumFormat)
		})

		t.Run("dashes", func(t *testing.T) {
			t.Parallel()

			_, err := NewPhoneNum("555-555-5555")
			require.ErrorIs(t, err, ErrInvalidPhoneNumFormat)
		})

		t.Run("dots", func(t *testing.T) {
			t.Parallel()

			_, err := NewPhoneNum("555.555.5555")
			require.ErrorIs(t, err, ErrInvalidPhoneNumFormat)
		})

		t.Run("spaces", func(t *testing.T) {
			t.Parallel()

			_, err := NewPhoneNum("555 555 5555")
			require.ErrorIs(t, err, ErrInvalidPhoneNumFormat)
		})

		t.Run("american without country code", func(t *testing.T) {
			t.Parallel()

			_, err := NewPhoneNum("(555) 555-5555")
			require.ErrorIs(t, err, ErrInvalidPhoneNumFormat)
		})

		t.Run("american with country code", func(t *testing.T) {
			t.Parallel()

			_, err := NewPhoneNum("1 (555) 555-5555")
			require.ErrorIs(t, err, ErrInvalidPhoneNumFormat)
		})
	})
}
