package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

const url = "http://localhost:8080/"

func TestRoute(t *testing.T) {

	t.Run("it returns 2 on request to '/'", func(t *testing.T) {
		queryData := "?firstNum=1&secondNum=1&sign=plus"
		request, _ := http.NewRequest(http.MethodGet, url+queryData, nil)
		response := httptest.NewRecorder()

		computationHandler(response, request)

		if response.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", response.Code)
		}

		expected := "2\n"
		if response.Body.String() != expected {
			t.Fatalf("expected %s got %s", expected, response.Body.String())
		}
	})

	t.Run("it returns 4 on request to '/'", func(t *testing.T) {
		queryData := "?firstNum=2&secondNum=2&sign=multiply"
		request, _ := http.NewRequest(http.MethodGet, url+queryData, nil)
		response := httptest.NewRecorder()

		computationHandler(response, request)

		if response.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", response.Code)
		}

		expected := "4\n"
		if response.Body.String() != expected {
			t.Fatalf("expected %s got %s", expected, response.Body.String())
		}
	})

	t.Run("it returns 3 on request to '/'", func(t *testing.T) {
		queryData := "?firstNum=5&secondNum=2&sign=minus"
		request, _ := http.NewRequest(http.MethodGet, url+queryData, nil)
		response := httptest.NewRecorder()

		computationHandler(response, request)

		if response.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", response.Code)
		}

		expected := "3\n"
		if response.Body.String() != expected {
			t.Fatalf("expected %s got %s", expected, response.Body.String())
		}
	})

	t.Run("it returns 4 on request to '/'", func(t *testing.T) {
		queryData := "?firstNum=16&secondNum=4&sign=divide"
		request, _ := http.NewRequest(http.MethodGet, url+queryData, nil)
		response := httptest.NewRecorder()

		computationHandler(response, request)

		if response.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", response.Code)
		}

		expected := "4\n"
		if response.Body.String() != expected {
			t.Fatalf("expected %s got %s", expected, response.Body.String())
		}
	})

	t.Run("it returns 'nothing' on request to '/'", func(t *testing.T) {
		queryData := "?firstNum=1&secondNum=1&sign=wrong"
		request, _ := http.NewRequest(http.MethodGet, url+queryData, nil)
		response := httptest.NewRecorder()

		computationHandler(response, request)

		if response.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", response.Code)
		}

		expected := "nothing\n"
		if response.Body.String() != expected {
			t.Fatalf("expected %s got %s", expected, response.Body.String())
		}
	})

	t.Run("it returns 'invalid firstNum' on request to '/'", func(t *testing.T) {
		queryData := "?firstNum=&secondNum=1&sign=plus"
		request, _ := http.NewRequest(http.MethodGet, url+queryData, nil)
		response := httptest.NewRecorder()

		computationHandler(response, request)

		if response.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400, got %d", response.Code)
		}

		expected := "invalid firstNum\n"
		if response.Body.String() != expected {
			t.Fatalf("expected %s got %s", expected, response.Body.String())
		}
	})
}
