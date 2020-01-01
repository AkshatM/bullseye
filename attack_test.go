package main

import (
	"strings"
	"testing"
)

// Helper function to generate a bitmap corresponding to header profiles
func generateAllBinaryStrings(n int) []string {

	if n == 1 {
		return []string{"0", "1"}
	}

	if n > 1 {
		stringsOfSmallerLength := generateAllBinaryStrings(n - 1)
		stringsOfLengthN := []string{}

		for _, binaryString := range stringsOfSmallerLength {
			stringsOfLengthN = append(stringsOfLengthN, string(binaryString)+"0")
			stringsOfLengthN = append(stringsOfLengthN, string(binaryString)+"1")
		}

		return stringsOfLengthN
	}

	return []string{}
}

func TestGenerateHeaders(test *testing.T) {

	all_header_profiles := generateAllBinaryStrings(len(possible_headers))

	for _, profile := range all_header_profiles {

		generated_headers := generateHeaders(profile)
		for index, bit := range profile {
			if string(bit) == "1" {
				header := possible_headers[index]
				ok := generated_headers.Get(header)

				if len(ok) == 0 {
					test.Errorf("Missing header. Profile string: %s. Expected header: %s. Generated header map: %s",
						profile, header, generated_headers)
				}
			}
		}

		// if no header is set, check if header map is empty
		if !strings.Contains(profile, "1") {
			if len(generated_headers) != 0 {
				test.Errorf("Header map not empty. Profile string: %s. Generated header map: %s", profile, generated_headers)
			}
		}
	}

}
