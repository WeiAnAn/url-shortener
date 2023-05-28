package utils_test

import (
	"testing"

	"github.com/WeiAnAn/url-shortener/internal/utils"
)

func TestGenerateURLReturnValidString(t *testing.T) {
	var generator utils.RandomBase62StringGenerator

	for i := 1; i <= 7; i++ {
		str, err := generator.Generate(i)

		if err != nil {
			t.Error(err)
		}

		if len(str) != i {
			t.Errorf("received string length %d is not equal to given length %d", i, i)
		}

		for _, c := range str {
			if !(c >= 'A' && c <= 'Z') &&
				!(c >= 'a' && c <= 'z') &&
				!(c >= '0' && c <= '9') {
				t.Errorf("char %c is not base62 char", c)
			}
		}
	}
}
