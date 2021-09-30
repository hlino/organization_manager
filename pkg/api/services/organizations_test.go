package services

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func Test_checkFilter(t *testing.T) {
	var testCases = []struct {
		regex            string
		inputFilter      string
		expectedGroupLen int
		continuousFilter bool
		expectedResult   []string
		shouldFail       bool
	}{
		{
			// Testing creating a valid range filter with open upper bound
			regex: rangeFilterRegex,
			inputFilter: "creation_date:[2002-09-22T00:00:00ZTO*]",
			expectedGroupLen: 6,
			continuousFilter: true,
			expectedResult: []string{"creation_date:[2002-09-22T00:00:00ZTO*]",
				"creation_date", "[", "2002-09-22T00:00:00Z", "*", "]"},
			shouldFail: false,
		},
		{
			// Testing creating a valid range filter with open lower bound
			regex: rangeFilterRegex,
			inputFilter: "creation_date:[*TO2002-09-22T00:00:00Z]",
			expectedGroupLen: 6,
			continuousFilter: true,
			expectedResult: []string{"creation_date:[*TO2002-09-22T00:00:00Z]",
				"creation_date", "[", "*", "2002-09-22T00:00:00Z", "]"},
			shouldFail: false,
		},
		{
			// Testing creating a valid range filter with upper and lower bounds
			regex: rangeFilterRegex,
			inputFilter: "employee_count:[5TO10)",
			expectedGroupLen: 6,
			continuousFilter: true,
			expectedResult: []string{"employee_count:[5TO10)",
				"employee_count", "[", "5", "10", ")"},
			shouldFail: false,
		},
		{
			// Testing creating an invalid range filter
			regex: rangeFilterRegex,
			inputFilter: "employee_count: [5TO10)",
			expectedGroupLen: 6,
			continuousFilter: true,
			expectedResult: []string{},
			shouldFail: true,
		},
		{
			// Testing creating an invalid range filter
			regex: rangeFilterRegex,
			inputFilter: "employee_count:[5-10)",
			expectedGroupLen: 6,
			continuousFilter: true,
			expectedResult: []string{},
			shouldFail: true,
		},
		{
			// Testing creating an invalid range filter
			regex: rangeFilterRegex,
			inputFilter: "employee_count:[5TO10",
			expectedGroupLen: 6,
			continuousFilter: true,
			expectedResult: []string{},
			shouldFail: true,
		},
		{
			// Testing creating a valid categorical filter
			regex: categoryFilterRegex,
			inputFilter: "name:CLEAR",
			expectedGroupLen: 3,
			continuousFilter: false,
			expectedResult: []string{"name:CLEAR", "name", "CLEAR"},
			shouldFail: false,
		},
		{
			// Testing creating a valid categorical filter with wildcard character
			regex: categoryFilterRegex,
			inputFilter: "name:CLEAR*",
			expectedGroupLen: 3,
			continuousFilter: false,
			expectedResult: []string{"name:CLEAR*", "name", "CLEAR*"},
			shouldFail: false,
		},
		{
			// Testing creating an invalid categorical filter
			regex: categoryFilterRegex,
			inputFilter: "name=CLEAR*",
			expectedGroupLen: 3,
			continuousFilter: false,
			expectedResult: []string{},
			shouldFail: true,
		},
		{
			// Testing creating an invalid categorical filter
			regex: categoryFilterRegex,
			inputFilter: ":CLEAR*",
			expectedGroupLen: 3,
			continuousFilter: false,
			expectedResult: []string{},
			shouldFail: true,
		},
		{
			// Testing a range filter on a categorical attribute, should fail
			regex: rangeFilterRegex,
			inputFilter: "name:[2002-09-22T00:00:00ZTO*]",
			expectedGroupLen: 6,
			continuousFilter: true,
			expectedResult: nil,
			shouldFail: true,
		},
		{
			// Testing creating a valid filter with an invalid column name
			regex: categoryFilterRegex,
			inputFilter: "wrong_name:CLEAR*",
			expectedGroupLen: 3,
			continuousFilter: false,
			expectedResult: []string{},
			shouldFail: true,
		},
		{
			// Testing creating a valid filter with an invalid column name
			regex: rangeFilterRegex,
			inputFilter: "wrong_name:[2002-09-22T00:00:00ZTO*]",
			expectedGroupLen: 6,
			continuousFilter: true,
			expectedResult: nil,
			shouldFail: true,
		},
	}

	for i, test := range testCases {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			res, err := checkFilter(test.regex, test.inputFilter, test.expectedGroupLen, test.continuousFilter)
			if test.shouldFail {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedResult, res)
			}
		})
	}
}

func Test_getPaginationQueryParams(t *testing.T) {
	var testCases = []struct {
		queryParams url.Values
		expectedPage int
		expectedPageSize int
		shouldFail bool
	}{
		{
			queryParams: map[string][]string{pageQueryParam: {"1"}, pageSizeQueryParam: {"10"}},
			expectedPage: 1,
			expectedPageSize: 10,
			shouldFail: false,
		},
		{
			queryParams: map[string][]string{pageQueryParam: {"3"}, pageSizeQueryParam: {"15"}},
			expectedPage: 3,
			expectedPageSize: 15,
			shouldFail: false,
		},
		{
			// Testing page default value
			queryParams: map[string][]string{pageSizeQueryParam: {"15"}},
			expectedPage: 1,
			expectedPageSize: 15,
			shouldFail: false,
		},
		{
			// Testing pageSize default value
			queryParams: map[string][]string{pageQueryParam: {"3"}},
			expectedPage: 3,
			expectedPageSize: 20,
			shouldFail: false,
		},
		{
			// Testing default page param values
			queryParams: map[string][]string{},
			expectedPage: 1,
			expectedPageSize: 20,
			shouldFail: false,
		},
		{
			// Testing invalid page param values
			queryParams: map[string][]string{pageQueryParam: {"not a number"}, pageSizeQueryParam: {"15"}},
			expectedPage: 3,
			expectedPageSize: 15,
			shouldFail: true,
		},
		{
			// Testing invalid page param values
			queryParams: map[string][]string{pageQueryParam: {"1"}, pageSizeQueryParam: {"not a number"}},
			expectedPage: 3,
			expectedPageSize: 15,
			shouldFail: true,
		},
	}

	for i, test := range testCases {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			page, pageSize, err := getPaginationQueryParams(test.queryParams)
			if test.shouldFail {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedPage, page)
				assert.Equal(t, test.expectedPageSize, pageSize)
			}
		})
	}
}
