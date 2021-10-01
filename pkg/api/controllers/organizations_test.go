package controllers

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"organization_manager/pkg/api/services"
	"organization_manager/pkg/database"
	"organization_manager/pkg/database/models"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestSaveNewOrganization(t *testing.T) {
	emptyOrg := models.Organization{}

	var testCases = []struct {
		requestBody          []byte
		expectedOrganization models.Organization
		expectedRespCode     int
	}{
		{
			requestBody: []byte(`{"name": "Organization 1","creation_date": "2021-09-26T00:00:00Z",
								"employee_count": 10,"is_public": false}`),
			expectedOrganization: models.Organization{
				Name:          "Organization 1",
				CreationDate:  time.Date(2021, 9, 26, 0, 0, 0, 0, time.UTC),
				EmployeeCount: 10,
				IsPublic:      false,
			},
			expectedRespCode: http.StatusCreated,
		},
		{
			requestBody: []byte(`{"invalid":"invalid",creation_date": "2021-09-26T00:00:00Z",
								"employee_count": 10,"is_public": false}`),
			expectedOrganization: emptyOrg,
			expectedRespCode:     http.StatusBadRequest,
		},
		{
			requestBody: []byte(`{"id": "1eacb0fa-d4ae-4d5e-9b69-268c1359db19", "name": "Organization 1","creation_date": "2021-09-26T00:00:00Z",
								"employee_count": 10,"is_public": false}`),
			expectedOrganization: emptyOrg,
			expectedRespCode:     http.StatusBadRequest,
		},
	}

	_, mock, err := database.InitializeTest()
	assert.NoError(t, err)

	for i, test := range testCases {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/organizations", bytes.NewBuffer(test.requestBody))
			w := httptest.NewRecorder()

			if test.expectedOrganization != emptyOrg {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "organizations" ("id","name","creation_date","employee_count","is_public") VALUES ($1,$2,$3,$4,$5)`)).
					WithArgs(sqlmock.AnyArg(), test.expectedOrganization.Name,
						test.expectedOrganization.CreationDate, test.expectedOrganization.EmployeeCount,
						test.expectedOrganization.IsPublic).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			}

			CreateOrganization(w, req)
			res := w.Result()
			assert.Equal(t, test.expectedRespCode, res.StatusCode)

			if test.expectedRespCode == http.StatusCreated {
				// ensures a new id was assigned to the created org and returned
				var respObj models.Organization
				err := json.NewDecoder(res.Body).Decode(&respObj)
				assert.NoError(t, err)
				assert.NotEqual(t, uuid.Nil, respObj.ID)
			} else if test.expectedRespCode == http.StatusBadRequest {
				var respObj ErrorResponse
				err := json.NewDecoder(res.Body).Decode(&respObj)
				assert.NoError(t, err)
				assert.True(t, strings.Contains(respObj.Error, "invalid request body"))
			}
		})
	}
}

func TestGetOrganizations(t *testing.T) {
	var tests = []struct {
		queryParams              map[string][]string
		expectedQueryConditional string
		expectedQueryLimit       string
		expectedArgs             []driver.Value
		expectedResponseCode     int
	}{
		{
			queryParams:              map[string][]string{"filter": {"name:CLEAR"}},
			expectedQueryConditional: `WHERE name = $1`,
			expectedQueryLimit:       `LIMIT 20`,
			expectedArgs:             []driver.Value{"CLEAR"},
			expectedResponseCode:     http.StatusOK,
		},
		{
			queryParams:              map[string][]string{"filter": {"name:CLEAR"}, "range_filter": {"creation_date:[2002-09-22T00:00:00ZTO*]"}, "page": {"2"}},
			expectedQueryConditional: `WHERE name = $1 AND creation_date >= $2`,
			expectedQueryLimit:       `LIMIT 20 OFFSET 20`,
			expectedArgs:             []driver.Value{"CLEAR", "2002-09-22T00:00:00Z"},
			expectedResponseCode:     http.StatusOK,
		},
		{
			queryParams:              map[string][]string{"filter": {"name:CLEAR"}, "range_filter": {"creation_date:[2002-09-22T00:00:00ZTO*]"}, "page": {"2"}},
			expectedQueryConditional: `WHERE name = $1 AND creation_date >= $2`,
			expectedQueryLimit:       `LIMIT 20 OFFSET 20`,
			expectedArgs:             []driver.Value{"CLEAR", "2002-09-22T00:00:00Z"},
			expectedResponseCode:     http.StatusNotFound,
		},
		{
			queryParams:              map[string][]string{"filter": {"name:CLEAR"}, "range_filter": {"creation_date:[2002-09-22T00:00:00ZTO*]"}, "page": {"r"}},
			expectedQueryConditional: `WHERE name = $1 AND creation_date >= $2`,
			expectedQueryLimit:       `LIMIT 20 OFFSET 20`,
			expectedArgs:             []driver.Value{"CLEAR", "2002-09-22T00:00:00Z"},
			expectedResponseCode:     http.StatusBadRequest,
		},
	}

	_, mock, err := database.InitializeTest()
	assert.NoError(t, err)

	for i, test := range tests {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/organizations", nil)
			w := httptest.NewRecorder()
			q := req.URL.Query()

			for queryKey, queryValue := range test.queryParams {
				for _, value := range queryValue {
					q.Add(queryKey, value)
				}
			}
			req.URL.RawQuery = q.Encode()

			mockSearchQueries(test.expectedResponseCode, mock, test.expectedArgs, test.expectedQueryLimit,
				test.expectedQueryConditional)

			GetOrganizations(w, req)
			res := w.Result()
			assert.Equal(t, test.expectedResponseCode, res.StatusCode)

			if test.expectedResponseCode == http.StatusOK {
				assert.NoError(t, mock.ExpectationsWereMet())
				var respObj services.PaginatedOrganizationResponse
				err := json.NewDecoder(res.Body).Decode(&respObj)
				assert.NoError(t, err)
				// mock db call will always return a single org
				assert.Equal(t, 1, len(respObj.Organizations))
			} else {
				var respObj ErrorResponse
				err := json.NewDecoder(res.Body).Decode(&respObj)
				assert.NoError(t, err)
				assert.NotEqual(t, "", respObj.Error)
			}
		})
	}
}

func mockSearchQueries(expectedResponseCode int, mock sqlmock.Sqlmock, expectedArgs []driver.Value, expectedQueryLimit,
	expectedQueryConditional string) {

	var organizationColumns = []string{"id", "name", "created_date", "employee_count", "is_public"}
	var countCol = []string{"count"}

	var returnedCount = 1
	if expectedResponseCode == http.StatusNotFound {
		returnedCount = 0
	}
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "organizations" ` + expectedQueryConditional)).
		WithArgs(expectedArgs...).WillReturnRows(sqlmock.NewRows(countCol).AddRow(returnedCount))

	searchQuery := mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "organizations" ` + expectedQueryConditional + ` ORDER BY id ` + expectedQueryLimit)).
		WithArgs(expectedArgs...)
	if expectedResponseCode == http.StatusNotFound {
		searchQuery.WillReturnRows(sqlmock.NewRows(organizationColumns))
	} else {
		searchQuery.WillReturnRows(sqlmock.NewRows(organizationColumns).
			AddRow(uuid.New(), "CLEAR", "2002-09-22T00:00:00Z", 10000, true))
	}
}
