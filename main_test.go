package main

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
		// пока сравнивать не будем, а просто выведем ответы
		// удалите потом этот вывод
		fmt.Println(response.Body.String())
	}
}

func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeCount(t *testing.T) {
	cities := []string{"moscow", "tula"} // передаваемые значения параметра city
	counts := []int{0, 1, 2, 3, 100}     // передаваемые значения параметра count
	handler := http.HandlerFunc(mainHandle)

	for _, city := range cities {
		for _, count := range counts {
			response := httptest.NewRecorder()
			url := fmt.Sprintf("/cafe?city=%s&count=%d", city, count)
			req := httptest.NewRequest("GET", url, nil)
			handler.ServeHTTP(response, req)

			switch {
			case count == 0:
				assert.Empty(t, response.Body.String())
			case count > 0 && count <= len(cafeList[city]):
				assert.Equal(t, count, len(strings.Split(response.Body.String(), ",")))
			default:
				assert.Equal(t, len(cafeList[city]), len(strings.Split(response.Body.String(), ",")))
			}
			require.Equal(t, http.StatusOK, response.Code)
		}
	}
}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)
	requests := []struct {
		search    string // передаваемое значение search
		wantCount int    // ожидаемое количество кафе в ответе
	}{
		{"фасоль", 0},
		{"кофе", 2},
		{"вилка", 1},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		url := fmt.Sprintf("/cafe?city=moscow&search=%s", v.search)
		req := httptest.NewRequest("GET", url, nil)
		handler.ServeHTTP(response, req)
		if v.wantCount == 0 {
			assert.Empty(t, response.Body.String())
		} else {
			assert.Equal(t, v.wantCount, len(strings.Split(response.Body.String(), ",")))
			assert.Contains(t, strings.ToLower(response.Body.String()), strings.ToLower(v.search))
		}
		require.Equal(t, http.StatusOK, response.Code)
	}
}
