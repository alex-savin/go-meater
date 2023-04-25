package meater

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func TestNewMeaterClient(t *testing.T) {
	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		assert.Equal(t, req.Method, POST)
		assert.Equal(t, req.Header.Get("Accept"), "application/json")
		assert.Equal(t, req.Header.Get("Content-Type"), "application/json")
		assert.Equal(t, req.URL.String(), "/v1/login")

		rw.WriteHeader(200)
		rw.Write([]byte(`{"status": "OK","statusCode": 200,"data": {"token": "TOKENISHERE","userId": "USERIDISHERE"},"meta": {}}`))
	}))
	defer server.Close()

	meater, err := New(
		Username("username"),
		Password("password"),
		BaseURL(server.URL),
		LogLevel("debug"),
	)

	fmt.Printf("MEATER CLIENT: %+v\n", meater)
	ok(t, err)

}

func TestGetProbes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.Method, GET)
		assert.Equal(t, req.Header.Get("Accept"), "application/json")
		// assert.Equal(t, req.Header.Get("Content-Type"), "application/json")
		assert.Equal(t, req.URL.String(), "/v1/devices")

		rw.WriteHeader(200)
		// rw.Header().Set("Content-Type", "application/json")
		rw.Write([]byte(`{"status":"OK","statusCode":200,"data":{"devices":[]},"meta":{}}`))
		// rw.Write([]byte(`{"status":"OK","statusCode":200,"data":{"devices":[{"id":"0fd5c70240a2c5c161f30964568b9117a46fcf6463fd3400b52e9c8c376196eb","temperature":{"internal":23.4,"ambient":23.4},"cook":{"id":"ded9c7fdcff311da52380cee3d3c419fd798a098a85d31581e12dd7c6e14ec51","name":null,"state":"Not Started","temperature":{"target":null,"peak":23.4},"time":{"elapsed":243,"remaining":-1}},"updated_at":1681750060},{"id":"9e14f62597268b6ec0ffbb1064d33ce953d5c817106c9128a7cccf21d4fab2e8","temperature":{"internal":23.2,"ambient":23.2},"cook":{"id":"764a386b1f539db3c80090112af2e0ef","name":null,"state":"Not Started","temperature":{"target":null,"peak":null},"time":{"elapsed":0,"remaining":-1}},"updated_at":1681749760},{"id":"b9e04921d5e1bb65a3b0a5025c5bac018b59b5d217aeb3c249fa5acc8426b032","temperature":{"internal":23.4,"ambient":23.4},"cook":{"id":"5cfa8d3d089a2b933f67b39841263c63fd078914bb5b7c8a0428dffbe2441a83","name":null,"state":"Not Started","temperature":{"target":null,"peak":23.4},"time":{"elapsed":1673,"remaining":-1}},"updated_at":1681750038},{"id":"81f7172c1ac47f2e74c990e48b35dbaa4f26b3f46413f506791b4dd3e229aa6b","temperature":{"internal":23.6,"ambient":23.6},"cook":{"id":"79a623047f9d01403b7351ae2852e5895b73eadb5ee61743601a74f6d67a3fbd","name":null,"state":"Not Started","temperature":{"target":null,"peak":23.6},"time":{"elapsed":146,"remaining":-1}},"updated_at":1681750062}]},"meta":{}}`))
	}))
	defer server.Close()

	httpClient := resty.New()
	httpClient.
		SetBaseURL(server.URL).
		SetHeaders(map[string]string{
			"Accept":          "application/json",
			"User-Agent":      "Mozilla/5.0 (Linux; Android 10; Android SDK built for x86 Build/QSR1.191030.002; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/74.0.3729.185 Mobile Safari/537.36",
			"Accept-Language": "en-US,en;q=0.9",
			"Accept-Encoding": "gzip, deflate"})

	meater := Client{
		baseURL: server.URL,
		// credentials:     credentials{},
		httpClient:      httpClient,
		isAuthenticated: true,
		logLevel:        "info",
	}

	meater.GetProbes()

	fmt.Printf("MEATER CLIENT: %+v\n", meater)

}

func TestGetProbeByID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.Method, GET)
		assert.Equal(t, req.Header.Get("Accept"), "application/json")
		assert.Equal(t, req.URL.String(), "/v1/devices/b9e04921d5e1bb65a3b0a5025c5bac018b59b5d217aeb3c249fa5acc8426b032")

		rw.WriteHeader(200)
		rw.Write([]byte(`{"status":"OK","statusCode":200,"data":{"id":"b9e04921d5e1bb65a3b0a5025c5bac018b59b5d217aeb3c249fa5acc8426b032","temperature":{"internal":23.3,"ambient":23.3},"cook":{"id":"5cfa8d3d089a2b933f67b39841263c63fd078914bb5b7c8a0428dffbe2441a83","name":"Sirloin Steak","state":"Configured","temperature":{"target":55,"peak":23.3},"time":{"elapsed":245,"remaining":-1}},"updated_at":1681748609},"meta":{}}`))
	}))
	defer server.Close()

	httpClient := resty.New()
	httpClient.
		SetBaseURL(server.URL).
		SetHeaders(map[string]string{
			"Accept":          "application/json",
			"User-Agent":      "Mozilla/5.0 (Linux; Android 10; Android SDK built for x86 Build/QSR1.191030.002; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/74.0.3729.185 Mobile Safari/537.36",
			"Accept-Language": "en-US,en;q=0.9",
			"Accept-Encoding": "gzip, deflate"})

	meater := Client{
		baseURL: server.URL,
		// credentials:     credentials{},
		httpClient:      httpClient,
		isAuthenticated: true,
		logLevel:        "info",
	}

	meater.GetProbeByID("b9e04921d5e1bb65a3b0a5025c5bac018b59b5d217aeb3c249fa5acc8426b032")

	fmt.Printf("MEATER CLIENT: %+v\n", meater)

}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
// func equals(tb testing.TB, exp, act interface{}) {
// 	if !reflect.DeepEqual(exp, act) {
// 		_, file, line, _ := runtime.Caller(1)
// 		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
// 		tb.FailNow()
// 	}
// }

// type RoundTripFunc func(req *http.Request) *http.Response

// func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
// 	return f(req), nil
// }

// // Test4xxError tests that a 4xx error is returned as an APIError.
// func Test4xxError(t *testing.T) {
// 	hclient := &http.Client{Transport: RoundTripFunc(func(req *http.Request) *http.Response {
// 		return &http.Response{
// 			StatusCode: http.StatusBadRequest,
// 			Body:       ioutil.NopCloser(strings.NewReader(`{"error":"Character not found"}`)),
// 		}
// 	})}

// 	credentials := credentials{
// 		username: "username",
// 		password: "password",
// 	}

// 	meater, err := New(credentials)

// 	_, err := client.GetCharacter(-1) // -1 is not a valid character ID so this should return an error
// 	if err == nil {
// 		t.Error("expected error, got nil")
// 	}

// 	apiError, ok := err.(*APIError)
// 	if !ok {
// 		t.Error("expected error to be of type APIError")
// 	}
// 	if apiError.Err != "Character not found" {
// 		t.Errorf("expected error message to be Character not found, got %s", apiError.Err)
// 	}
// }

// // Test5xxError tests that a 5xx error is returned as an APIError.
// func Test5xxError(t *testing.T) {
// 	hclient := &http.Client{Transport: RoundTripFunc(func(req *http.Request) *http.Response {
// 		return &http.Response{
// 			StatusCode: http.StatusInternalServerError,
// 			Body:       ioutil.NopCloser(strings.NewReader(`Internal server error`)),
// 		}
// 	})}

// 	client := Client{client: hclient, baseURL: "https://public-api.cloud.meater.com/"}

// 	_, err := client.GetCharacter(1)
// 	if err == nil {
// 		t.Error("expected error, got nil")
// 	}

// 	if _, ok := err.(*APIError); ok {
// 		t.Error("expected error to not be of type APIError")
// 	}
// }

// func TestGetCharacterByID(t *testing.T) {
// 	characterJSON := `{"id":2,"name":"Morty Smith","status":"Alive","species":"Human","type":"","gender":"Male","image":"https://rickandmortyapi.com/api/character/avatar/2.jpeg","url":"https://rickandmortyapi.com/api/character/2","created":"2017-11-04T18:50:21.651Z"}`

// 	hclient := &http.Client{Transport: RoundTripFunc(func(req *http.Request) *http.Response {
// 		if req.URL.String() != "https://rickandmortyapi.com/api/character/2" {
// 			t.Error("expected request to https://rickandmortyapi.com/api/character/2, got", req.URL.String())
// 		}
// 		if req.Method != "GET" {
// 			t.Error("expected request method to be GET, got", req.Method)
// 		}
// 		if req.Header.Get("Accept") != "application/json" {
// 			t.Error("expected request Accept header to be application/json, got", req.Header.Get("Accept"))
// 		}
// 		if req.Header.Get("Content-Type") != "application/json" {
// 			t.Error("expected request Content-Type header to be application/json, got", req.Header.Get("Content-Type"))
// 		}

// 		return &http.Response{
// 			StatusCode: http.StatusOK,
// 			Body:       ioutil.NopCloser(strings.NewReader(characterJSON)),
// 		}
// 	})}

// 	client := Client{client: hclient, baseURL: "https://public-api.cloud.meater.com/"}

// 	character, err := client.GetCharacter(2)
// 	if err != nil {
// 		t.Errorf("unexpected error: %v", err)
// 	}

// 	if character.ID != 2 {
// 		t.Errorf("expected character ID to be 2, got %d", character.ID)
// 	}
// 	if character.Name != "Morty Smith" {
// 		t.Errorf("expected character name to be Morty Smith, got %s", character.Name)
// 	}
// }
