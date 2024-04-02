package config

import (
	"fmt"
	"strings"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

func TestPhoneNumSet(t *testing.T) {
	t.Parallel()

	numTrials := 10_000

	t.Run("Universal set accepts all numbers", func(t *testing.T) {
		t.Parallel()

		pns := PhoneNumSet{}
		require.NoError(t, pns.Set("*"))

		for range numTrials {
			matchNum := models.PhoneNum(faker.E164PhoneNumber())
			assert.True(t, pns.Match(matchNum), fmt.Sprintf("Expected valid phone number %q to match", matchNum))
		}
	})

	t.Run("Empty set accepts no numbers", func(t *testing.T) {
		t.Parallel()

		pns := PhoneNumSet{}
		require.NoError(t, pns.Set(""))

		for range numTrials {
			matchNum := models.PhoneNum(faker.E164PhoneNumber())
			assert.False(t, pns.Match(matchNum), fmt.Sprintf("Expected valid phone number %q to fail to match", matchNum))
		}
	})

	t.Run("parser handles extra commas or whitespace", func(t *testing.T) {
		t.Parallel()

		for numPos := range 5 {
			// Empty inputs
			inputElems := make([]string, 5)

			// Add varying amounts of whitespace to each
			for i := range 5 {
				inputElems[i] = strings.Repeat(" ", i)
			}

			// Append an actual phone number to just one of the inputElems
			matchNum := faker.E164PhoneNumber()
			inputElems[numPos] += matchNum

			// Valid matches
			pns := PhoneNumSet{}
			require.NoError(t, pns.Set(strings.Join(inputElems, ",")))
			assert.True(t, pns.Match(models.PhoneNum(matchNum)))

			// Invalid matches
			for range numTrials {
				genNum := faker.E164PhoneNumber()
				if genNum == matchNum {
					continue
				}

				assert.False(t, pns.Match(models.PhoneNum(genNum)), fmt.Sprintf("Expected valid phone number %q to fail to match", genNum))
			}
		}
	})

	for _, expectedMatchCount := range []int{1, 3, 10, 100} {
		t.Run(fmt.Sprintf("Correctly handles %d matches", expectedMatchCount), func(t *testing.T) {
			t.Parallel()

			// Store as strings to allow constructing the parsable form and as a set of models.PhoneNum to validate the format and ensure easy lookup to avoid collisions when generating phone numbers that *shouldn't* match.
			matchedStrings := make([]string, 0, expectedMatchCount)
			matchedNums := make(map[models.PhoneNum]struct{}, expectedMatchCount)

			for range expectedMatchCount {
				numStr := faker.E164PhoneNumber()
				matchedStrings = append(matchedStrings, numStr)

				num, err := models.NewPhoneNum(numStr)
				require.NoError(t, err, "parse test phone num")

				matchedNums[num] = struct{}{}
			}

			pns := PhoneNumSet{}
			require.NoError(t, pns.Set(strings.Join(matchedStrings, ",")))

			// All expected numbers match.
			for num := range matchedNums {
				assert.True(t, pns.Match(num))
			}

			// Unexpected numbers do not match.
			for i := 0; i < numTrials; {
				num := models.PhoneNum(faker.E164PhoneNumber())
				if _, has := matchedNums[num]; has {
					// Skip this iteration if we accidentally generated one of the test numbers.
					continue
				}

				assert.False(t, pns.Match(num), fmt.Sprintf("Expected valid phone number %q to fail to match", num))

				i++
			}
		})
	}
}
