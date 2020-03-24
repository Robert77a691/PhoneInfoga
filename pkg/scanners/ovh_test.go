package scanners

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sundowndev/phoneinfoga/pkg/utils"
	gock "gopkg.in/h2non/gock.v1"
)

func TestOVHScanner(t *testing.T) {
	assert := assert.New(t)

	t.Run("should find number on OVH", func(t *testing.T) {
		defer gock.Off() // Flush pending mocks after test execution

		gock.New("https://api.ovh.com").
			Get("/1.0/telephony/number/detailedZones").
			MatchParam("country", "fr").
			Reply(200).
			JSON([]ovhAPIResponseNumber{
				ovhAPIResponseNumber{
					ZneList:             []string{},
					MatchingCriteria:    "",
					Prefix:              33,
					InternationalNumber: "003336517xxxx",
					Country:             "fr",
					ZipCode:             "",
					Number:              "036517xxxx",
					City:                "Abbeville",
					AskedCity:           "",
				},
			})

		number, _ := LocalScan("+33 0365179268")

		result := ovhScanCLI(utils.LoggerService, number)

		assert.Equal(result, &OVHScannerResponse{
			Found:       true,
			NumberRange: "036517xxxx",
			City:        "Abbeville",
			ZipCode:     "",
		}, "they should be equal")

		assert.Equal(gock.IsDone(), true, "there should have no pending mocks")
	})

	t.Run("should not find number on OVH", func(t *testing.T) {
		defer gock.Off() // Flush pending mocks after test execution

		gock.New("https://api.ovh.com").
			Get("/1.0/telephony/number/detailedZones").
			MatchParam("country", "us").
			Reply(200).
			JSON([]ovhAPIResponseNumber{
				ovhAPIResponseNumber{
					ZneList:             []string{},
					MatchingCriteria:    "",
					Prefix:              33,
					InternationalNumber: "003336517xxxx",
					Country:             "fr",
					ZipCode:             "",
					Number:              "036517xxxx",
					City:                "Abbeville",
					AskedCity:           "",
				},
			})

		number, _ := LocalScan("+1 718-521-2994")

		result, err := OVHScan(number)

		assert.Equal(err, nil, "should not be errored")
		assert.Equal(result, &OVHScannerResponse{
			Found:       false,
			NumberRange: "",
			City:        "",
			ZipCode:     "",
		}, "they should be equal")

		assert.Equal(gock.IsDone(), true, "there should have no pending mocks")
	})
}
