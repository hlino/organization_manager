package models

import (
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"organization_manager/pkg/database"
	"time"
)

type Organization struct {
	ID            uuid.UUID `gorm:"primary_key;column:id"`
	Name          string    `gorm:"column:name" json:"name"`
	CreationDate  time.Time `gorm:"column:creation_date" json:"creation_date"`
	EmployeeCount int       `gorm:"employee_count" json:"employee_count"`
	IsPublic      bool      `gorm:"column:is_public" json:"is_public"`
}

type RangeQueryFilter struct {
	DBfield         string
	StartRange      string
	StartComparator string
	EndRange        string
	EndComparator   string
}

type CategoryQueryFilter struct {
	DBfield     string
	LikeFilter  string
	ExactFilter string
}

const (
	GTE = ">="
	GT = ">"
	LTE = "<="
	LT = "<"
	OpenRangeDelimiter = "*"
)

// Map of organization column name to boolean value determining if the field is continuous or not
var OrganizationColumnNamesContinuousMap = map[string]bool {
	"name": false,
	"creation_date": true,
	"employee_count": true,
	"is_public": false,
}

func (o *Organization) Save() error {
	o.ID = uuid.New()
	return database.DB.Create(o).Error
}

func SearchForOrganizations(categoryFilters []CategoryQueryFilter, rangeFilters []RangeQueryFilter, page,
	pageSize int) ([]Organization, int64, error) {

	query := database.DB
	for _, categoryFilter := range categoryFilters {
		if categoryFilter.LikeFilter != "" {
			query = query.Where(fmt.Sprintf("%s LIKE ?", categoryFilter.DBfield), categoryFilter.LikeFilter)
		} else {
			query = query.Where(fmt.Sprintf("%s = ?", categoryFilter.DBfield), categoryFilter.ExactFilter)
		}
	}

	for _, rangeFilter := range rangeFilters {
		if rangeFilter.StartRange != OpenRangeDelimiter {
			query = query.Where(fmt.Sprintf("%s %s ?", rangeFilter.DBfield, rangeFilter.StartComparator),
				rangeFilter.StartRange)
		}
		if rangeFilter.EndRange != OpenRangeDelimiter {
			query = query.Where(fmt.Sprintf("%s %s ?", rangeFilter.DBfield, rangeFilter.EndComparator),
				rangeFilter.EndRange)
		}
	}

	// Setting pagination parameters on query
	offset := (page - 1) * pageSize
	query = query.Limit(pageSize)
	query = query.Offset(offset)

	var totalCount int64
	var organizationsFound []Organization
	err := query.Model(&Organization{}).Debug().Order("id").Count(&totalCount).Find(&organizationsFound).Error
	log.Infof("Count: %d", totalCount)
	return organizationsFound, totalCount, err
}
