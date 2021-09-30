package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
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
			expectedRespCode: http.StatusBadRequest,
		},
		{
			requestBody: []byte(`{"id": "1eacb0fa-d4ae-4d5e-9b69-268c1359db19", "name": "Organization 1","creation_date": "2021-09-26T00:00:00Z",
								"employee_count": 10,"is_public": false}`),
			expectedOrganization: emptyOrg,
			expectedRespCode: http.StatusBadRequest,
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
			} else if test.expectedRespCode ==http.StatusBadRequest {
				var respObj ErrorResponse
				err := json.NewDecoder(res.Body).Decode(&respObj)
				assert.NoError(t, err)
				assert.True(t, strings.Contains(respObj.Error, "invalid request body"))
			}
		})
	}
}

func TestGetOrganizations(t *testing.T) {
	//req := httptest.NewRequest(http.MethodPost, "/organizations", bytes.NewBuffer(test.requestBody))
	//w := httptest.NewRecorder()
	//q := req.URL.Query()
	//if test.userId != nil {
	//	q.Add("user_id", (*test.userId).String())
	//}
	//if test.modelId != nil {
	//	q.Add("model_id", (*test.modelId).String())
	//}
	//if test.page != nil {
	//	q.Add("page", strconv.Itoa(*test.page))
	//}
	//if test.pageSize != nil {
	//	q.Add("page_size", strconv.Itoa(*test.pageSize))
	//}
	//req.URL.RawQuery = q.Encode()
}
