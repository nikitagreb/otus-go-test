package main

import (
	"io"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	cases := []struct {
		testName string
		from     string
		offset   int64
		limit    int64
	}{
		{"out offset0 limit0", "out_offset0_limit0.txt", 0, 0},
		{"out offset0 limit10", "out_offset0_limit10.txt", 0, 10},
		{"out offset0 limit1000", "out_offset0_limit1000.txt", 0, 1000},
		{"out offset0 limit10000", "out_offset0_limit10000.txt", 0, 10000},
	}

	for _, tcase := range cases {
		t.Run(tcase.testName, func(t *testing.T) {
			from := "testdata/" + tcase.from

			f, err := os.CreateTemp("", "test_copy_")
			if err != nil {
				log.Fatal(err)
			}

			defer os.Remove(f.Name())

			errCopy := Copy(from, f.Name(), tcase.offset, tcase.limit)
			require.NoError(t, errCopy)

			bCopy, err := os.ReadFile(f.Name())
			if err != nil {
				log.Fatal(err)
			}
			bOriginal, err := os.ReadFile(from)
			if err != nil {
				log.Fatal(err)
			}

			require.Equal(t, string(bOriginal), string(bCopy))
		})
	}

	t.Run("out offset100 limit10000", func(t *testing.T) {
		from := "testdata/" + "out_offset100_limit1000.txt"

		f, err := os.CreateTemp("", "test_copy_")
		if err != nil {
			log.Fatal(err)
		}

		defer os.Remove(f.Name())

		errCopy := Copy(from, f.Name(), 100, 1000)
		require.NoError(t, errCopy)

		bCopy, err := os.ReadFile(f.Name())
		if err != nil {
			log.Fatal(err)
		}

		fOriginal, err := os.Open(from)
		if err != nil {
			log.Fatal(err)
		}

		_, err = fOriginal.Seek(100, io.SeekStart)
		if err != nil {
			log.Fatal(err)
		}

		bOriginal := make([]byte, 900)
		_, err = fOriginal.Read(bOriginal)
		if err != nil {
			log.Fatal(err)
		}

		require.Equal(t, string(bOriginal), string(bCopy))
	})

	t.Run("out offset6000 limit1000", func(t *testing.T) {
		from := "testdata/" + "out_offset6000_limit1000.txt"

		f, err := os.CreateTemp("", "test_copy_")
		if err != nil {
			log.Fatal(err)
		}

		defer os.Remove(f.Name())

		errCopy := Copy(from, f.Name(), 6000, 1000)
		require.Errorf(t, errCopy, "offset exceeds file size")
	})
}
