package services

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io"
	"math"
	"net/http"
	"net/url"
	"organization_manager/pkg/database/models"
	"regexp"
	"strconv"
	"strings"
)

const (
	pageQueryParam               = "page"
	defaultPage                  = 1
	defaultPageSize              = 20
	pageSizeQueryParam           = "page_size"
	filterQueryParam             = "filter"
	rangeFilterQueryParam        = "range_filter"
	startInclusiveRangeDelimiter = "["
	endInclusiveRangeDelimiter   = "]"
	rangeFilterRegex             = `(.*):(\(|\[)(.*)TO(.*)(\)|\])`
	categoryFilterRegex          = `(.*):(.*)`
)

type PaginatedOrganizationResponse struct {
	Organizations []models.Organization `json:"organizations"`
	Page          int                   `json:"page"`
	PageSize      int                   `json:"page_size"`
	TotalPages    int                   `json:"total_pages"`
	TotalCount    int                   `json:"total_count"`
}

// SaveNewOrganization deserializes POST request and saves and returns a new organization object
// Will return an error and associated http response code as well
func SaveNewOrganization(requestContent io.ReadCloser) (*models.Organization, int, error) {
	var orgRequestObject models.Organization
	err := json.NewDecoder(requestContent).Decode(&orgRequestObject)
	if orgRequestObject.ID != uuid.Nil {
		log.Error("organization request content already contained ID value")
		return nil, http.StatusBadRequest, errors.New("invalid request body")
	}
	if err != nil {
		log.Errorf("error deserializing organization POST request body: %v", err)
		return nil, http.StatusBadRequest, errors.Wrap(err, "invalid request body")
	}

	err = orgRequestObject.Save()
	if err != nil {
		log.Errorf("error saving new organization: %v", err)
		return nil, http.StatusInternalServerError, err
	}

	return &orgRequestObject, http.StatusCreated, nil
}

// GetOrganizations parses query parameters from GET request to create database query and returns paginated result
// Will return an error and associated http response code as well
func GetOrganizations(queryParams url.Values) (*PaginatedOrganizationResponse, int, error) {
	categoryQueryFilters, _ := queryParams[filterQueryParam]
	rangeQueryFilters, _ := queryParams[rangeFilterQueryParam]

	page, pageSize, err := getPaginationQueryParams(queryParams)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	// Process categorical filters from query parameters
	var categoryDBFilters = make([]models.CategoryQueryFilter, len(categoryQueryFilters))
	for i, filter := range categoryQueryFilters {
		matchedGroups, err := checkFilter(categoryFilterRegex, filter, 3, false)
		if err != nil {
			return nil, http.StatusBadRequest, err
		}

		categoryDBFilters[i].DBfield = matchedGroups[1]
		// determines if the filtering value has a wildcard character
		if strings.Contains(matchedGroups[2], "*") {
			queryMatcher := strings.ReplaceAll(matchedGroups[2], "*", "%")
			categoryDBFilters[i].LikeFilter = queryMatcher
		} else {
			categoryDBFilters[i].ExactFilter = matchedGroups[2]
		}
	}

	// Process range filters from query parameters
	var rangeDBFilters = make([]models.RangeQueryFilter, len(rangeQueryFilters))
	for i, filter := range rangeQueryFilters {
		matchedGroups, err := checkFilter(rangeFilterRegex, filter, 6, true)
		if err != nil {
			return nil, http.StatusBadRequest, err
		}

		rangeDBFilters[i].DBfield = matchedGroups[1]
		rangeDBFilters[i].StartRange = matchedGroups[3]
		rangeDBFilters[i].EndRange = matchedGroups[4]
		if matchedGroups[2] == startInclusiveRangeDelimiter {
			rangeDBFilters[i].StartComparator = models.GTE
		} else {
			rangeDBFilters[i].StartComparator = models.GT
		}
		if matchedGroups[5] == endInclusiveRangeDelimiter {
			rangeDBFilters[i].EndComparator = models.LTE
		} else {
			rangeDBFilters[i].EndComparator = models.LT
		}
	}

	//sends parsed query params from request to query the database
	orgs, totalCount, err := models.SearchForOrganizations(categoryDBFilters, rangeDBFilters, page, pageSize)
	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))
	respObj := PaginatedOrganizationResponse{
		Organizations: orgs,
		Page:          page,
		PageSize:      pageSize,
		TotalPages:    totalPages,
		TotalCount:    int(totalCount),
	}
	var returnHTTPStatus = http.StatusOK
	if err != nil {
		// database query errors should be returned as an internal server error
		returnHTTPStatus = http.StatusInternalServerError
	}
	return &respObj, returnHTTPStatus, err
}

// getPaginationQueryParams parses request query parameters for pagination values
func getPaginationQueryParams(queryParams url.Values) (int, int, error) {
	var err error

	p := queryParams.Get(pageQueryParam)
	page := defaultPage
	if p != "" {
		page, err = strconv.Atoi(p)
		if err != nil {
			log.Errorf("error parsing page query param: %v", err)
			return -1, -1, errors.Errorf("invalid page query parameter '%s'", p)
		}
	}

	pageSize := defaultPageSize
	ps := queryParams.Get(pageSizeQueryParam)
	if ps != "" {
		pageSize, err = strconv.Atoi(ps)
		if err != nil {
			log.Errorf("error parsing page_size query param: %v", err)
			return -1, -1, errors.Errorf("invalid page_size query parameter '%s'", ps)
		}
	}
	return page, pageSize, nil
}

// checkFilter parses and validates the value from a filter or filter_range query, will return components of the parsed
// filter or an error if the query parameter was not provided correctly
func checkFilter(regexStrMatcher, filter string, expectedGroupLength int, continuousFilter bool) ([]string, error) {
	r, err := regexp.Compile(regexStrMatcher)
	if err != nil {
		return nil, err
	}
	isValidFilter := r.MatchString(filter)
	if !isValidFilter {
		return nil, errors.Errorf("invalid filter '%s'", filter)
	}

	groups := r.FindStringSubmatch(filter)
	if len(groups) != expectedGroupLength {
		return nil, errors.Errorf("invalid filter '%s'", filter)
	}

	// Checks to ensure a continuous attribute was provided with a range_filter, categorical attributes cannot be
	// used in range filters
	isAttrContinuous, colExists := models.OrganizationColumnNamesContinuousMap[groups[1]]
	if !colExists {
		return nil, errors.Errorf("invalid column name '%s'", groups[1])
	}
	if continuousFilter && !isAttrContinuous {
		return nil, errors.Errorf("cannot supply range filter for categorical column '%s'", groups[1])
	}

	return groups, nil
}
