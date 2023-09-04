package spdy

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpgradeResponse(t *testing.T) {
	testCases := []struct {
		connectionHeader string
		upgradeHeader    string
		shouldError      bool
	}{
		{
			connectionHeader: "",
			upgradeHeader:    "",
			shouldError:      true,
		},
		{
			connectionHeader: "Upgrade",
			upgradeHeader:    "",
			shouldError:      true,
		},
		{
			connectionHeader: "",
			upgradeHeader:    "SPDY/3.1",
			shouldError:      true,
		},
		{
			connectionHeader: "Upgrade",
			upgradeHeader:    "SPDY/3.1",
			shouldError:      false,
		},
	}

	for i, testCase := range testCases {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			upgrader := NewResponseUpgrader()
			conn := upgrader.UpgradeResponse(w, req, nil)
			haveErr := conn == nil
			if e, a := testCase.shouldError, haveErr; e != a {
				t.Fatalf("%d: expected shouldErr=%t, got %t", i, testCase.shouldError, haveErr)
			}
			if haveErr {
				return
			}
			if conn == nil {
				t.Fatalf("%d: unexpected nil conn", i)
			}
			defer conn.Close()
		}))
		defer server.Close()

		req, err := http.NewRequest("GET", server.URL, nil)
		if err != nil {
			t.Fatalf("%d: error creating request: %s", i, err)
		}

		req.Header.Set("Connection", testCase.connectionHeader)
		req.Header.Set("Upgrade", testCase.upgradeHeader)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("%d: unexpected non-nil err from client.Do: %s", i, err)
		}

		if testCase.shouldError {
			continue
		}

		if resp.StatusCode != http.StatusSwitchingProtocols {
			t.Fatalf("%d: expected status 101 switching protocols, got %d", i, resp.StatusCode)
		}
	}
}
