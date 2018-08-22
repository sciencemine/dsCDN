package server

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetCes(t *testing.T) {
	path := "/ce"
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		log.Fatalln(err)
		return
	}

	rr := httptest.NewRecorder()

	router := NewRouter()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf(testFailed(path, rr.Code, http.StatusOK))
	}
}

func TestGetCe(t *testing.T) {
	tt := []struct {
		name          string
		routeVariable string
		shouldPass    bool
	}{
		{"valid id", "5b5f730b95dfef70164be84b", true},
		{"nonsense id", "puppiesandkittensandfish", false},
		{"empty id", "", false},
	}

	router := NewRouter()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			path := fmt.Sprintf("/ce/%s", tc.routeVariable)
			req, err := http.NewRequest("GET", path, nil)
			if err != nil {
				t.Fatal(err)
				return
			}

			rr := httptest.NewRecorder()

			//start the server and send a request
			router.ServeHTTP(rr, req)

			if rr.Code == http.StatusOK && !tc.shouldPass {
				t.Errorf(testFailed(path, rr.Code, http.StatusBadRequest))
			}

			if rr.Code != http.StatusOK && tc.shouldPass {
				t.Errorf(testFailed(path, rr.Code, http.StatusOK))
			}
		})
	}
}

func TestGetDsms(t *testing.T) {
	path := "/dsm"
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		log.Fatalln(err)
		return
	}

	rr := httptest.NewRecorder()

	router := NewRouter()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf(testFailed(path, rr.Code, http.StatusOK))
	}
}

func TestGetDsm(t *testing.T) {
	tt := []struct {
		name          string
		routeVariable string
		shouldPass    bool
	}{
		{"valid id", "5b5f730b95dfef70164be85b", true},
		{"nonsense id", "puppiesandkittensandfish", false},
		{"empty id", "", false},
	}

	router := NewRouter()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			path := fmt.Sprintf("/dsm/%s", tc.routeVariable)
			req, err := http.NewRequest("GET", path, nil)
			if err != nil {
				log.Fatalln("Could not make request: ", err)
				return
			}

			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			//It passed when it shouldn't have
			if rr.Code == http.StatusOK && !tc.shouldPass {
				t.Errorf(testFailed(path, rr.Code, http.StatusBadRequest))
			}

			//It didn't pass and it should've
			if rr.Code != http.StatusOK && tc.shouldPass {
				t.Errorf(testFailed(path, rr.Code, http.StatusOK))
			}
		})
	}
}

func TestGetAssets(t *testing.T) {
	path := "/asset"
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		log.Fatalln(err)
		return
	}

	rr := httptest.NewRecorder()

	router := NewRouter()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf(testFailed(path, rr.Code, http.StatusOK))
	}
}

func TestGetPaths(t *testing.T) {
	path := "/path"
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		log.Fatalln(err)
		return
	}

	rr := httptest.NewRecorder()

	router := NewRouter()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf(testFailed(path, rr.Code, http.StatusOK))
	}
}

func TestGetPath(t *testing.T) {
	tt := []struct {
		name          string
		routeVariable string
		shouldPass    bool
	}{
		{"valid id", "5b4e27ae6f10865ae67fdd25", true},
		{"nonsense id", "puppiesandkittensandfish", false},
		{"empty id", "", false},
	}

	router := NewRouter()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			path := fmt.Sprintf("/path/%s", tc.routeVariable)
			req, err := http.NewRequest("GET", path, nil)
			if err != nil {
				t.Errorf("Could not make request: %v", err)
				return
			}

			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			//It passed when it shouldn't have
			if rr.Code == http.StatusOK && !tc.shouldPass {
				t.Errorf(testFailed(path, rr.Code, http.StatusBadRequest))
			}

			//It didn't pass and it should've
			if rr.Code != http.StatusOK && tc.shouldPass {
				t.Errorf(testFailed(path, rr.Code, http.StatusOK))
			}
		})
	}
}

type testData struct {
	name         string
	requestBody  string
	expectedCode int
}

func TestPostCe(t *testing.T) {
	tt := []testData{
		{"empty body", "", http.StatusBadRequest},
		{"valid object", `{"version": "0.0.0","title": "Test","description": "Silverbow test", "playlist": { "teaser": "5b5f730b95dfef70164be82f", "queue" : [{"primary": "5b5f730b95dfef70164be830","backgrounds": [{"asset": "5b5f730b95dfef70164be830","x" : 0.0,"y" : 0.0,"start" : 0,"duration" : 0,"width" : 0,"height" : 0}],"tracks": [],"overlays": []}]}}`,
			http.StatusCreated},
		{"invalid json", `{whoopsie}`, 422},
	}

	path := "/ce"

	router := NewRouter()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var reader io.Reader

			reader = strings.NewReader(tc.requestBody)

			req, err := http.NewRequest("POST", path, reader)
			if err != nil {
				log.Println("Could not make request")
				return
			}
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			if tc.expectedCode == rr.Code {
				t.Logf("Passed Case: %s \n expectedStatus : %d \n  observedStatusCode : %d \n",
					tc.name, tc.expectedCode, rr.Code)
			} else {
				t.Errorf("Failed Case: %s \n expectedStatus : %d  \n observedStatusCode : %d \n",
					tc.name, tc.expectedCode, rr.Code)
			}
		})
	}
}

func TestPostAsset(t *testing.T) {
	tt := []testData{
		{"empty body", "", http.StatusBadRequest},
		{"valid object",
			`{
			"id": "5b15cbfcc0d0633fd211e7ec",
			"version": "0.0.0",
			"url": "http://66.62.91.16:60080/smds/SILVERBOW-02.webm",
			"type": "video/webm",
			"title": "Silverbow 2",
			"options": {}
		}`, http.StatusCreated},
		{"invalid json", `{whoopsie}`, 422},
	}

	path := "/asset"

	router := NewRouter()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var reader io.Reader

			reader = strings.NewReader(tc.requestBody)

			req, err := http.NewRequest("POST", path, reader)
			if err != nil {
				log.Println("Could not make request")
				return
			}
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			if tc.expectedCode == rr.Code {
				t.Logf("Passed Case: %s \n expectedStatus : %d \n  observedStatusCode : %d \n",
					tc.name, tc.expectedCode, rr.Code)
			} else {
				t.Errorf("Failed Case: %s \n expectedStatus : %d  \n observedStatusCode : %d \n",
					tc.name, tc.expectedCode, rr.Code)
			}
		})
	}
}

func TestPostPath(t *testing.T) {
	tt := []testData{
		{"empty body", "", http.StatusBadRequest},
		{"valid object",
			`{
			"id" : "",
			"model": {
				"id_model": "5b47c6b8132bb67388ee2534",
				"version_model": "0.0.0",
				"descriptionn_model": "What makes a stream healthy?",
				"author": "Michael Fryer"
			},
			"relations": [
				{
				"ce_list": [
					{
					"id_ce": "5b47c6b8132bb67388ee2530",
					"version_ce": "0.0.0",
					"title_ce": "What are \"the Five Cs\"?",
					"description_ce": "The 5Cs",
					"events": []
					}
				]
				},
				{
				"title_attr": "The 5Cs",
				"description_attr": "Cold, Clear, Clean, Complex, and Connected",
				"weight": 1,
				"ce_list": [
					{
					"id_ce": "5b47c6b8132bb67388ee2530",
					"version_ce": "0.0.0",
					"title_ce": "What are \"the Five Cs\"?",
					"description_ce": "The 5Cs",
					"events": []
					},
					{
					"id_ce": "5b47c6b8132bb67388ee2529",
					"version_ce": "0.0.0",
					"title_ce": "Characteristics of a healthy stream?",
					"description_ce": "Complexity",
					"events": []
					}
				]
				}
			]
			}`,
			http.StatusCreated},
		{"invalid json", `{whoopsie}`, 422},
	}

	path := "/path"

	router := NewRouter()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var reader io.Reader

			reader = strings.NewReader(tc.requestBody)

			req, err := http.NewRequest("POST", path, reader)
			if err != nil {
				log.Println("Could not make request")
				return
			}
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			if tc.expectedCode == rr.Code {
				t.Logf("Passed Case: %s \n expectedStatus : %d \n  observedStatusCode : %d \n",
					tc.name, tc.expectedCode, rr.Code)
			} else {
				t.Errorf("Failed Case: %s \n expectedStatus : %d  \n observedStatusCode : %d \n",
					tc.name, tc.expectedCode, rr.Code)
			}
		})
	}
}

func TestPostDsm(t *testing.T) {
	tt := []testData{
		{"empty body", "", http.StatusBadRequest},
		{"invalid json", `{whoopsie}`, 422},
		{"valid object", `{
    "title": "Test Streams",
    "description": "What makes a test healthy?",
    "version": "0.0.0",
    "author": "Michael Test",
    "config": {
        "menuLocale": "top",
        "idle": 15,
        "menuDwell": 3,
        "popoverDwell": 5,
        "popoverShowDelay": 0.2
    },
    "tactile": {
        "select": [ "w" ],
        "previous": [ "a" ],
        "cancel": [ "s" ],
        "next": [ "d" ]
    },
    "stylesheet": "",
    "style": {
        "idle": "",
        "menu": "",
        "select": "",
        "playing": ""
    },
    "contributors": [""],
    "idle_backgrounds": [
        "5b572dc5040aa32a9f61e6ba",
        "5b572dc5040aa32a9f61e6b0",
        "5b572dc5040aa32a9f61e6bc",
        "5b572dc5040aa32a9f61e6b1"
    ],
    "video_select_backgrounds": [],
    "ce_set": {
        "5b572dc5040aa32a9f61e6bb": {
            "attributes": [
                2
            ],
            "relationships": []
        },
        "5b572dc5040aa32a9f61e6b5": {
            "attributes": [
                4,
                0
            ],
            "relationships": []
        },
        "5b572dc5040aa32a9f61e6b8": {
            "attributes": [
                1
            ],
            "relationships": []
        },
        "5b572dc5040aa32a9f61e6bd": {
            "attributes": [
                2
            ],
            "relationships": []
        },
        "5b572dc5040aa32a9f61e6b3": {
            "attributes": [
                2
            ],
            "relationships": []
        },
        "5b572dc5040aa32a9f61e6b4": {
            "attributes": [
                2,
                0
            ],
            "relationships": []
        },
        "5b572dc5040aa32a9f61e6b9": {
            "attributes": [
                1
            ],
            "relationships": []
        },
        "5b572dc5040aa32a9f61e6b7": {
            "attributes": [
                1
            ],
            "relationships": []
        },
        "5b572dc5040aa32a9f61e6b2": {
            "attributes": [
                1
            ],
            "relationships": []
        },
        "5b572dc5040aa32a9f61e6af": {
            "attributes": [
                2
            ],
            "relationships": []
        },
        "5b572dc5040aa32a9f61e6b6": {
            "attributes": [
                3,
                0
            ],
            "relationships": []
        },
        "5b572dc5040aa32a9f61e6be": {
            "attributes": [
                3,
                0
            ],
            "relationships": []
        }
    },
    "attributes": [
        {
            "title": "Test Habitats",
            "description": "Raparian Tests"
        },
        {
            "title": "Tests",
            "description": "Tests"
        },
        {
            "title": "The 5 Tests",
            "description": "Cold, Clear, Clean, Complex, and Connected"
        },
        {
            "title": "Test",
            "description": "Plants"
        },
        {
            "title": "Test Beavers",
            "description": "Beavers"}]}`, http.StatusCreated},
	}

	path := "/dsm"

	router := NewRouter()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var reader io.Reader

			reader = strings.NewReader(tc.requestBody)

			req, err := http.NewRequest("POST", path, reader)
			if err != nil {
				log.Println("Could not make request")
				return
			}
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			if tc.expectedCode == rr.Code {
				t.Logf("Passed Case: %s \n expectedStatus : %d \n  observedStatusCode : %d \n",
					tc.name, tc.expectedCode, rr.Code)
			} else {
				t.Errorf("Failed Case: %s \n expectedStatus : %d  \n observedStatusCode : %d \n",
					tc.name, tc.expectedCode, rr.Code)
			}
		})
	}
}

func testFailed(path string, got, wanted int) string {
	return fmt.Sprintf("test %s failed: Got %d Wanted %d", path, got, wanted)
}

func Test_writeAssetToPart(t *testing.T) {
	type args struct {
		assetObj Asset
		writer   *multipart.Writer
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := writeAssetToPart(tt.args.assetObj, tt.args.writer); (err != nil) != tt.wantErr {
				t.Errorf("writeAssetToPart() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_writeDescToParts(t *testing.T) {
	type args struct {
		arr    []desc
		writer *multipart.Writer
		seen   map[string]bool
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writeDescToParts(tt.args.arr, tt.args.writer, tt.args.seen)
		})
	}
}
