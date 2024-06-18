package integrationtests

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestHelloPage(t *testing.T) {
	status, resp := sendRequest("http://localhost:8080/hello", http.MethodGet, nil, nil)

	require.Equal(t, http.StatusOK, status)
	require.Equal(t, "hello", string(resp))
}

func TestCreateEventSuccess(t *testing.T) {
	uuid := uuid.NewString()

	formData := url.Values{}
	formData.Add("id", uuid)
	formData.Add("title", "Test")
	formData.Add("start_dt", "2024-06-13")
	formData.Add("end_dt", "2024-07-13")
	formData.Add("notify_before", "24h")

	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	status, resp := sendRequest("http://localhost:8080/event/create", http.MethodPost, formData, headers)
	require.Equal(t, http.StatusCreated, status, string(resp))

	formData = url.Values{}
	formData.Add("id", uuid)
	sendRequest("http://localhost:8080/event/delete", http.MethodPost, formData, headers)
}

func TestCreateEventValidationFail(t *testing.T) {
	status, resp := sendRequest("http://localhost:8080/event/create", http.MethodPost, nil, nil)
	require.Equal(t, http.StatusBadRequest, status, string(resp))

	formData := url.Values{}
	formData.Add("id", "invalid uuid")
	formData.Add("title", "Test")
	formData.Add("start_dt", "2024-06-13")
	formData.Add("end_dt", "2024-07-13")
	formData.Add("notify_before", "24h")

	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	status, resp = sendRequest("http://localhost:8080/event/create", http.MethodPost, formData, headers)
	require.Equal(t, http.StatusBadRequest, status, string(resp))

	formData.Set("id", uuid.NewString())
	formData.Set("start_dt", "invalid_date")

	status, resp = sendRequest("http://localhost:8080/event/create", http.MethodPost, formData, headers)
	require.Equal(t, http.StatusBadRequest, status, string(resp))

	formData.Set("start_dt", "2024-06-13")
	formData.Set("end_dt", "invalid_date")

	status, resp = sendRequest("http://localhost:8080/event/create", http.MethodPost, formData, headers)
	require.Equal(t, http.StatusBadRequest, status, string(resp))

	formData.Set("end_dt", "2024-07-13")
	formData.Set("notify_before", "invalid_duration")

	status, resp = sendRequest("http://localhost:8080/event/create", http.MethodPost, formData, headers)
	require.Equal(t, http.StatusBadRequest, status, string(resp))
}

func TestUpdateEventSuccess(t *testing.T) {
	uuid := uuid.NewString()

	formData := url.Values{}
	formData.Add("id", uuid)
	formData.Add("title", "Test")
	formData.Add("start_dt", "2024-06-13")
	formData.Add("end_dt", "2024-07-13")
	formData.Add("notify_before", "24h")

	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	sendRequest("http://localhost:8080/event/create", http.MethodPost, formData, headers)

	formData.Set("title", "Test Test")
	status, resp := sendRequest("http://localhost:8080/event/update", http.MethodPost, formData, headers)
	require.Equal(t, http.StatusOK, status, string(resp))

	formData = url.Values{}
	formData.Add("id", uuid)
	sendRequest("http://localhost:8080/event/delete", http.MethodPost, formData, headers)
}

func TestUpdateEventValidationFail(t *testing.T) {
	status, resp := sendRequest("http://localhost:8080/event/update", http.MethodPost, nil, nil)
	require.Equal(t, http.StatusBadRequest, status, string(resp))

	uuid := uuid.NewString()

	formData := url.Values{}
	formData.Add("id", uuid)
	formData.Add("title", "Test")
	formData.Add("start_dt", "2024-06-13")
	formData.Add("end_dt", "2024-07-13")
	formData.Add("notify_before", "24h")

	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	status, resp = sendRequest("http://localhost:8080/event/create", http.MethodPost, formData, headers)
	if status != http.StatusCreated {
		t.Error("Create event failed", string(resp))
	}

	formData.Set("start_dt", "invalid_date")
	status, resp = sendRequest("http://localhost:8080/event/update", http.MethodPost, formData, headers)
	require.Equal(t, http.StatusBadRequest, status, string(resp))

	formData.Set("start_dt", "2024-06-13")
	formData.Set("end_dt", "invalid_date")
	status, resp = sendRequest("http://localhost:8080/event/update", http.MethodPost, formData, headers)
	require.Equal(t, http.StatusBadRequest, status, string(resp))

	formData.Set("end_dt", "2024-07-13")
	formData.Set("notify_before", "invalid_duration")
	status, resp = sendRequest("http://localhost:8080/event/update", http.MethodPost, formData, headers)
	require.Equal(t, http.StatusBadRequest, status, string(resp))

	formData = url.Values{}
	formData.Add("id", uuid)
	sendRequest("http://localhost:8080/event/delete", http.MethodPost, formData, headers)
}

func TestDeleteEvent(t *testing.T) {
	uuid := uuid.NewString()

	formData := url.Values{}
	formData.Add("id", uuid)
	formData.Add("title", "Test")
	formData.Add("start_dt", "2024-06-13")
	formData.Add("end_dt", "2024-07-13")
	formData.Add("notify_before", "24h")

	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	sendRequest("http://localhost:8080/event/create", http.MethodPost, formData, headers)

	formData = url.Values{}
	formData.Add("id", uuid)
	status, resp := sendRequest("http://localhost:8080/event/delete", http.MethodPost, formData, headers)
	require.Equal(t, http.StatusOK, status, string(resp))
}

func TestGetEventFail(t *testing.T) {
	uuid := uuid.NewString()
	status, _ := sendRequest("http://localhost:8080/event/get?id="+uuid, http.MethodGet, nil, nil)
	require.Equal(t, http.StatusInternalServerError, status)
}

func TestGetEventSuccess(t *testing.T) {
	uuid := uuid.NewString()

	formData := url.Values{}
	formData.Add("id", uuid)
	formData.Add("title", "Test")
	formData.Add("start_dt", "2024-06-13")
	formData.Add("end_dt", "2024-07-13")
	formData.Add("notify_before", "24h")

	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	sendRequest("http://localhost:8080/event/create", http.MethodPost, formData, headers)

	status, resp := sendRequest("http://localhost:8080/event/get?id="+uuid, http.MethodGet, nil, nil)
	require.Equal(t, http.StatusOK, status, string(resp))
	require.Equal(
		t,
		//nolint: all
		`{"id":"`+uuid+`","title":"Test","description":"","start_dt":"2024-06-13T00:00:00Z","end_dt":"2024-07-13T00:00:00Z","creator_id":0,"notify_before":86400000000000}`,
		string(resp),
	)

	formData = url.Values{}
	formData.Add("id", uuid)
	sendRequest("http://localhost:8080/event/delete", http.MethodPost, formData, headers)
}

func TestListByDates(t *testing.T) {
	uuid1 := uuid.NewString()

	formData := url.Values{}
	formData.Add("id", uuid1)
	formData.Add("title", "Test")
	formData.Add("start_dt", "2024-06-01")
	formData.Add("end_dt", "2024-06-10")
	formData.Add("notify_before", "24h")

	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	sendRequest("http://localhost:8080/event/create", http.MethodPost, formData, headers)

	uuid2 := uuid.NewString()

	formData = url.Values{}
	formData.Add("id", uuid2)
	formData.Add("title", "Test2")
	formData.Add("start_dt", "2024-06-05")
	formData.Add("end_dt", "2024-06-10")
	formData.Add("notify_before", "24h")

	sendRequest("http://localhost:8080/event/create", http.MethodPost, formData, headers)

	uuid3 := uuid.NewString()

	formData = url.Values{}
	formData.Add("id", uuid3)
	formData.Add("title", "Test3")
	formData.Add("start_dt", "2024-06-06")
	formData.Add("end_dt", "2024-06-13")
	formData.Add("notify_before", "24h")

	sendRequest("http://localhost:8080/event/create", http.MethodPost, formData, headers)

	status, resp := sendRequest(
		"http://localhost:8080/event/listByDates?from=2024-06-01&to=2024-06-10",
		http.MethodGet,
		nil,
		nil,
	)

	require.Equal(t, http.StatusOK, status, string(resp))
	require.Equal(
		t,
		//nolint: all
		`[{"id":"`+uuid1+`","title":"Test","description":"","start_dt":"2024-06-01T00:00:00Z","end_dt":"2024-06-10T00:00:00Z","creator_id":0,"notify_before":86400000000000}{"id":"`+uuid2+`","title":"Test2","description":"","start_dt":"2024-06-05T00:00:00Z","end_dt":"2024-06-10T00:00:00Z","creator_id":0,"notify_before":86400000000000}]`,
		string(resp),
	)

	formData = url.Values{}
	formData.Add("id", uuid1)
	sendRequest("http://localhost:8080/event/delete", http.MethodPost, formData, headers)

	formData = url.Values{}
	formData.Add("id", uuid2)
	sendRequest("http://localhost:8080/event/delete", http.MethodPost, formData, headers)

	formData = url.Values{}
	formData.Add("id", uuid3)
	sendRequest("http://localhost:8080/event/delete", http.MethodPost, formData, headers)
}

func TestListOnDate(t *testing.T) {
	uuid1 := uuid.NewString()

	formData := url.Values{}

	formData.Add("id", uuid1)
	formData.Add("title", "Test")
	formData.Add("start_dt", "2024-06-01")
	formData.Add("end_dt", "2024-06-10")
	formData.Add("notify_before", "24h")

	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	sendRequest("http://localhost:8080/event/create", http.MethodPost, formData, headers)

	uuid2 := uuid.NewString()

	formData = url.Values{}
	formData.Add("id", uuid2)
	formData.Add("title", "Test2")
	formData.Add("start_dt", "2024-06-05")
	formData.Add("end_dt", "2024-06-10")
	formData.Add("notify_before", "24h")

	sendRequest("http://localhost:8080/event/create", http.MethodPost, formData, headers)

	status, resp := sendRequest("http://localhost:8080/event/listOnDate?date=2024-06-01", http.MethodGet, nil, nil)
	require.Equal(t, http.StatusOK, status, string(resp))
	require.Equal(
		t,
		//nolint: all
		`[{"id":"`+uuid1+`","title":"Test","description":"","start_dt":"2024-06-01T00:00:00Z","end_dt":"2024-06-10T00:00:00Z","creator_id":0,"notify_before":86400000000000}]`,
		string(resp),
	)

	formData = url.Values{}
	formData.Add("id", uuid1)
	sendRequest("http://localhost:8080/event/delete", http.MethodPost, formData, headers)

	formData = url.Values{}
	formData.Add("id", uuid2)
	sendRequest("http://localhost:8080/event/delete", http.MethodPost, formData, headers)
}

func TestListOnWeek(t *testing.T) {
	uuid1 := uuid.NewString()

	formData := url.Values{}
	formData.Add("id", uuid1)
	formData.Add("title", "Test")
	formData.Add("start_dt", "2024-06-03")
	formData.Add("end_dt", "2024-06-09")
	formData.Add("notify_before", "24h")

	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	sendRequest("http://localhost:8080/event/create", http.MethodPost, formData, headers)

	uuid2 := uuid.NewString()

	formData = url.Values{}
	formData.Add("id", uuid2)
	formData.Add("title", "Test2")
	formData.Add("start_dt", "2024-06-05")
	formData.Add("end_dt", "2024-06-10")
	formData.Add("notify_before", "24h")

	sendRequest("http://localhost:8080/event/create", http.MethodPost, formData, headers)

	status, resp := sendRequest(
		"http://localhost:8080/event/listOnWeek?weekStartDate=2024-06-03",
		http.MethodGet,
		nil,
		nil,
	)

	require.Equal(t, http.StatusOK, status, string(resp))
	require.Equal(
		t,
		//nolint: all
		`[{"id":"`+uuid1+`","title":"Test","description":"","start_dt":"2024-06-03T00:00:00Z","end_dt":"2024-06-09T00:00:00Z","creator_id":0,"notify_before":86400000000000}]`,
		string(resp),
	)

	formData = url.Values{}
	formData.Add("id", uuid1)
	sendRequest("http://localhost:8080/event/delete", http.MethodPost, formData, headers)

	formData = url.Values{}
	formData.Add("id", uuid2)
	sendRequest("http://localhost:8080/event/delete", http.MethodPost, formData, headers)
}

func TestListOnMonth(t *testing.T) {
	uuid1 := uuid.NewString()

	formData := url.Values{}
	formData.Add("id", uuid1)
	formData.Add("title", "Test")
	formData.Add("start_dt", "2024-06-03")
	formData.Add("end_dt", "2024-06-09")
	formData.Add("notify_before", "24h")

	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	sendRequest("http://localhost:8080/event/create", http.MethodPost, formData, headers)

	uuid2 := uuid.NewString()

	formData = url.Values{}
	formData.Add("id", uuid2)
	formData.Add("title", "Test2")
	formData.Add("start_dt", "2024-06-05")
	formData.Add("end_dt", "2024-06-10")
	formData.Add("notify_before", "24h")

	sendRequest("http://localhost:8080/event/create", http.MethodPost, formData, headers)

	status, resp := sendRequest(
		"http://localhost:8080/event/listOnMonth?monthStartDate=2024-06-01",
		http.MethodGet,
		nil,
		nil,
	)

	require.Equal(t, http.StatusOK, status, string(resp))
	require.Equal(
		t,
		//nolint: all
		`[{"id":"`+uuid1+`","title":"Test","description":"","start_dt":"2024-06-03T00:00:00Z","end_dt":"2024-06-09T00:00:00Z","creator_id":0,"notify_before":86400000000000}{"id":"`+uuid2+`","title":"Test2","description":"","start_dt":"2024-06-05T00:00:00Z","end_dt":"2024-06-10T00:00:00Z","creator_id":0,"notify_before":86400000000000}]`,
		string(resp),
	)

	formData = url.Values{}
	formData.Add("id", uuid1)
	sendRequest("http://localhost:8080/event/delete", http.MethodPost, formData, headers)

	formData = url.Values{}
	formData.Add("id", uuid2)
	sendRequest("http://localhost:8080/event/delete", http.MethodPost, formData, headers)
}

func sendRequest(url string, method string, formData url.Values, headers map[string]string) (status int, body []byte) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, method, url, strings.NewReader(formData.Encode()))

	for key := range headers {
		req.Header.Add(key, headers[key])
	}

	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("%v", err)
	}

	resData, _ := io.ReadAll(res.Body)

	defer res.Body.Close()

	return res.StatusCode, resData
}
