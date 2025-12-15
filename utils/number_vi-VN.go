package utils

import "strings"

var (
	units     = []string{"không", "một", "hai", "ba", "bốn", "năm", "sáu", "bảy", "tám", "chín"}
	chunks    = []string{"", "nghìn", "triệu", "tỷ"}
	unitNames = []string{"", "mươi", "trăm"}
)

func NumberToVietnamese(n int64) string {
	if n == 0 {
		return units[0]
	}

	var parts []string
	chunkCount := 0

	for n > 0 {
		chunk := n % 1000
		n = n / 1000

		if chunk > 0 {
			chunkText := convertThreeDigits(chunk, chunkCount == 0)
			if chunkCount > 0 {
				chunkText += " " + chunks[chunkCount]
			}
			parts = append([]string{chunkText}, parts...)
		} else if len(parts) > 0 && strings.HasPrefix(parts[0], chunks[len(chunks)-1]) {
			// Handle cases like 1,000,000,000 (1 tỷ)
			parts = append([]string{chunks[len(chunks)-1]}, parts...)
		}

		chunkCount++
		if chunkCount >= len(chunks) {
			chunkCount = 1 // Reset to triệu, tỷ, etc. after tỷ
		}
	}

	return strings.Join(parts, " ")
}

func convertThreeDigits(n int64, isLastChunk bool) string {
	var parts []string
	hundreds := n / 100
	tens := (n % 100) / 10
	ones := n % 10

	// Hundreds place
	if hundreds > 0 {
		parts = append(parts, units[hundreds]+" "+unitNames[2])
	} else if (tens > 0 || ones > 0) && !isLastChunk {
		parts = append(parts, units[0]+" "+unitNames[2])
	}

	// Tens and ones places
	if tens > 0 {
		if tens == 1 {
			parts = append(parts, "mười")
		} else {
			parts = append(parts, units[tens]+" "+unitNames[1])
		}

		if ones > 0 {
			if tens == 1 && ones == 5 {
				parts = append(parts, "lăm")
			} else if (tens > 1 || isLastChunk) && ones == 1 {
				parts = append(parts, "mốt")
			} else if ones == 4 && tens != 1 {
				parts = append(parts, "tư")
			} else {
				parts = append(parts, units[ones])
			}
		}
	} else if ones > 0 {
		if hundreds > 0 && ones == 5 {
			parts = append(parts, "lăm")
		} else if hundreds > 0 && ones == 1 {
			parts = append(parts, "mốt")
		} else {
			parts = append(parts, units[ones])
		}
	}

	return strings.Join(parts, " ")
}
