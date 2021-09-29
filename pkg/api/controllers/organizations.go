package controllers

import (
	"net/http"
	"organization_manager/pkg/api/services"
)

func CreateOrganization(w http.ResponseWriter, r *http.Request) {
	newOrg, httpRespCode, err := services.SaveNewOrganization(r.Body)
	if err != nil {
		JsonResponse(w, httpRespCode, ErrorResponse{err.Error()})
		return
	}

	JsonResponse(w, httpRespCode, newOrg)
}

func GetOrganizations(w http.ResponseWriter, r *http.Request) {
	// TODO: think about also returning http status code for 400 instead of 500
	resp, err := services.GetOrganizations(r.URL.Query())
	if err != nil {
		JsonResponse(w, http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	} else if len(resp.Organizations) == 0 {
		JsonResponse(w, http.StatusNotFound, ErrorResponse{"No organizations found"})
		return
	}
	JsonResponse(w, http.StatusOK, resp)
}
