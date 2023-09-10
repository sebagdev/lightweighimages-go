package main

import (
	"fmt"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"testing"
)

var images = []struct {
	image, tag string
}{
	{"docker.io/library/golang_raw_build_image", "1"},
	{"docker.io/library/multistage_build_image", "2"},
	{"docker.io/library/alpine_build_image", "3"},
	{"docker.io/library/smaller_out_alpine_image", "4"},
	{"docker.io/library/scratch_image", "5"},
	{"docker.io/library/distroless_image", "6"},
	{"docker.io/library/distroless_image", "7"},
}

func TestForTimeZone(t *testing.T) {
	for _, tt := range images {
		testname := fmt.Sprintf("TestForTimeZone %s:%s", tt.image, tt.tag)

		t.Run(testname, func(t *testing.T) {
			pool, err, resource := initContainer(t, tt.image, tt.tag)
			var resp *http.Response

			err = pool.Retry(func() error {
				resp, err = http.Get(fmt.Sprint("http://localhost:", resource.GetPort("8080/tcp"), "/currentTime/Europe%2FWarsaw"))
				if err != nil {
					t.Log("container not ready, waiting...")
					return err
				}
				return nil
			})
			require.NoError(t, err, "HTTP error")
			defer resp.Body.Close()

			require.Equal(t, http.StatusOK, resp.StatusCode, "HTTP status code")

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err, "failed to read HTTP body")

			require.Contains(t, string(body), "current_time", "should respond with current time")
		})

	}
}

func TestForOddNumber(t *testing.T) {
	for _, tt := range images {
		testname := fmt.Sprintf("TestForOddNumber %s:%s", tt.image, tt.tag)

		t.Run(testname, func(t *testing.T) {
			pool, err, resource := initContainer(t, tt.image, tt.tag)

			var resp *http.Response

			err = pool.Retry(func() error {
				resp, err = http.Get(fmt.Sprint("http://localhost:", resource.GetPort("8080/tcp"), "/isEven/1"))
				if err != nil {
					t.Log("container not ready, waiting...")
					return err
				}
				return nil
			})
			require.NoError(t, err, "HTTP error")
			defer resp.Body.Close()

			require.Equal(t, http.StatusOK, resp.StatusCode, "HTTP status code")

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err, "failed to read HTTP body")

			require.Contains(t, string(body), "{\"even\":false}", "does not respond with love?")
		})
	}
}

func TestForEvenNumber(t *testing.T) {
	for _, tt := range images {
		testname := fmt.Sprintf("TestForEvenNumber %s:%s", tt.image, tt.tag)

		t.Run(testname, func(t *testing.T) {
			pool, err, resource := initContainer(t, tt.image, tt.tag)

			var resp *http.Response

			err = pool.Retry(func() error {
				resp, err = http.Get(fmt.Sprint("http://localhost:", resource.GetPort("8080/tcp"), "/isEven/2"))
				if err != nil {
					t.Log("container not ready, waiting...")
					return err
				}
				return nil
			})
			require.NoError(t, err, "HTTP error")
			defer resp.Body.Close()

			require.Equal(t, http.StatusOK, resp.StatusCode, "HTTP status code")

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err, "failed to read HTTP body")

			require.Contains(t, string(body), "{\"even\":true}", "does not respond with love?")
		})
	}
}

func initContainer(t *testing.T, image, tag string) (*dockertest.Pool, error, *dockertest.Resource) {
	pool, err := dockertest.NewPool("")
	require.NoError(t, err, "could not connect to Docker")

	resource, err := pool.Run(image, tag, []string{})
	require.NoError(t, err, "could not start container")

	t.Cleanup(func() {
		require.NoError(t, pool.Purge(resource), "failed to remove container")
	})
	return pool, err, resource
}
