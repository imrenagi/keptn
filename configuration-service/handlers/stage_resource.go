package handlers

import (
	"io/ioutil"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/keptn/keptn/configuration-service/common"
	"github.com/keptn/keptn/configuration-service/config"
	"github.com/keptn/keptn/configuration-service/models"
	"github.com/keptn/keptn/configuration-service/restapi/operations/stage_resource"
)

// GetProjectProjectNameStageStageNameResourceHandlerFunc get list of stage resources
func GetProjectProjectNameStageStageNameResourceHandlerFunc(params stage_resource.GetProjectProjectNameStageStageNameResourceParams) middleware.Responder {
	if !common.StageExists(params.ProjectName, params.StageName) {
		return stage_resource.NewGetProjectProjectNameStageStageNameResourceNotFound().WithPayload(&models.Error{Code: 404, Message: swag.String("Stage does not exist")})
	}

	err := common.CheckoutBranch(params.ProjectName, params.StageName)
	if err != nil {
		return stage_resource.NewGetProjectProjectNameStageStageNameResourceDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
	}

	projectConfigPath := config.ConfigDir + "/" + params.ProjectName
	result := common.GetPaginatedResources(projectConfigPath, params.PageSize, params.NextPageKey)
	return stage_resource.NewGetProjectProjectNameStageStageNameResourceOK().WithPayload(result)
}

// GetProjectProjectNameStageStageNameResourceResourceURIHandlerFunc get the specified resource
func GetProjectProjectNameStageStageNameResourceResourceURIHandlerFunc(params stage_resource.GetProjectProjectNameStageStageNameResourceResourceURIParams) middleware.Responder {
	projectConfigPath := config.ConfigDir + "/" + params.ProjectName
	resourcePath := projectConfigPath + "/" + params.ResourceURI
	if !common.StageExists(params.ProjectName, params.StageName) {
		return stage_resource.NewGetProjectProjectNameStageStageNameResourceResourceURINotFound().WithPayload(&models.Error{Code: 404, Message: swag.String("Project not found")})
	}
	// utils.Debug("", "Checking out "+params.StageName+" branch")
	err := common.CheckoutBranch(params.ProjectName, params.StageName)
	if err != nil {
		return stage_resource.NewGetProjectProjectNameStageStageNameResourceResourceURIDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
	}

	if !common.FileExists(resourcePath) {
		return stage_resource.NewGetProjectProjectNameStageStageNameResourceResourceURINotFound().WithPayload(&models.Error{Code: 404, Message: swag.String("Stage resource not found")})
	}

	dat, err := ioutil.ReadFile(resourcePath)
	if err != nil {
		return stage_resource.NewGetProjectProjectNameStageStageNameResourceResourceURIDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
	}

	resourceContent := strfmt.Base64(dat)
	return stage_resource.NewGetProjectProjectNameStageStageNameResourceResourceURIOK().WithPayload(
		&models.Resource{
			ResourceURI:     &params.ResourceURI,
			ResourceContent: resourceContent,
		})
}

// PostProjectProjectNameStageStageNameResourceHandlerFunc creates list of new resources in a stage
func PostProjectProjectNameStageStageNameResourceHandlerFunc(params stage_resource.PostProjectProjectNameStageStageNameResourceParams) middleware.Responder {
	if !common.StageExists(params.ProjectName, params.StageName) {
		return stage_resource.NewPostProjectProjectNameStageStageNameResourceBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Stage does not exist")})
	}

	projectConfigPath := config.ConfigDir + "/" + params.ProjectName

	// utils.Debug("", "Creating new resource(s) in: "+projectConfigPath+" in stage "+params.StageName)
	// utils.Debug("", "Checking out branch: "+params.StageName)
	err := common.CheckoutBranch(params.ProjectName, params.StageName)
	if err != nil {
		return stage_resource.NewPostProjectProjectNameStageStageNameResourceBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String(err.Error())})
	}

	for _, res := range params.Resources.Resources {
		filePath := projectConfigPath + "/" + *res.ResourceURI
		// don't overwrite existing files
		if !common.FileExists(filePath) {
			// utils.Debug("", "Adding resource: "+filePath)
			common.WriteFile(filePath, res.ResourceContent)
		}
	}

	// utils.Debug("", "Staging Changes")
	err = common.StageAndCommitAll(params.ProjectName, "Added resources")
	if err != nil {
		return stage_resource.NewPostProjectProjectNameStageStageNameResourceBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String(err.Error())})
	}
	// utils.Debug("", "Successfully added resources")

	newVersion, err := common.GetCurrentVersion(params.ProjectName)
	if err != nil {
		return stage_resource.NewPostProjectProjectNameStageStageNameResourceBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String(err.Error())})
	}
	return stage_resource.NewPostProjectProjectNameStageStageNameResourceCreated().WithPayload(&models.Version{
		Version: newVersion,
	})
}

// PutProjectProjectNameStageStageNameResourceHandlerFunc updates list of stage resources
func PutProjectProjectNameStageStageNameResourceHandlerFunc(params stage_resource.PutProjectProjectNameStageStageNameResourceParams) middleware.Responder {
	if !common.StageExists(params.ProjectName, params.StageName) {
		return stage_resource.NewPutProjectProjectNameStageStageNameResourceBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Stage does not exist")})
	}
	projectConfigPath := config.ConfigDir + "/" + params.ProjectName

	// utils.Debug("", "Creating new resource(s) in: "+projectConfigPath+" in stage "+params.StageName)
	// utils.Debug("", "Checking out branch: "+params.StageName)
	err := common.CheckoutBranch(params.ProjectName, params.StageName)
	if err != nil {
		return stage_resource.NewPutProjectProjectNameStageStageNameResourceBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String(err.Error())})
	}

	for _, res := range params.Resources.Resources {
		filePath := projectConfigPath + "/" + *res.ResourceURI
		common.WriteFile(filePath, res.ResourceContent)
	}

	// utils.Debug("", "Staging Changes")
	err = common.StageAndCommitAll(params.ProjectName, "Updated resources")
	if err != nil {
		return stage_resource.NewPutProjectProjectNameStageStageNameResourceBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String(err.Error())})
	}
	// utils.Debug("", "Successfully updated resources")

	newVersion, err := common.GetCurrentVersion(params.ProjectName)
	if err != nil {
		return stage_resource.NewPutProjectProjectNameStageStageNameResourceBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String(err.Error())})
	}
	return stage_resource.NewPutProjectProjectNameStageStageNameResourceCreated().WithPayload(&models.Version{
		Version: newVersion,
	})
}

// PutProjectProjectNameStageStageNameResourceResourceURIHandlerFunc updates the specified stage resource
func PutProjectProjectNameStageStageNameResourceResourceURIHandlerFunc(params stage_resource.PutProjectProjectNameStageStageNameResourceResourceURIParams) middleware.Responder {
	if !common.StageExists(params.ProjectName, params.StageName) {
		return stage_resource.NewPutProjectProjectNameStageStageNameResourceResourceURIBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Stage does not exist")})
	}
	projectConfigPath := config.ConfigDir + "/" + params.ProjectName

	// utils.Debug("", "Creating new resource(s) in: "+projectConfigPath+" in stage "+params.StageName)
	// utils.Debug("", "Checking out branch: "+params.StageName)
	err := common.CheckoutBranch(params.ProjectName, params.StageName)
	if err != nil {
		return stage_resource.NewPutProjectProjectNameStageStageNameResourceResourceURIBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String(err.Error())})
	}

	filePath := projectConfigPath + "/" + params.ResourceURI
	common.WriteFile(filePath, params.Resource.ResourceContent)

	// utils.Debug("", "Staging Changes")
	err = common.StageAndCommitAll(params.ProjectName, "Updated resource: "+params.ResourceURI)
	if err != nil {
		return stage_resource.NewPutProjectProjectNameStageStageNameResourceResourceURIBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String(err.Error())})
	}
	// utils.Debug("", "Successfully updated resource: "+params.ResourceURI)

	newVersion, err := common.GetCurrentVersion(params.ProjectName)
	if err != nil {
		return stage_resource.NewPutProjectProjectNameStageStageNameResourceResourceURIBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String(err.Error())})
	}
	return stage_resource.NewPutProjectProjectNameStageStageNameResourceResourceURICreated().WithPayload(&models.Version{
		Version: newVersion,
	})
}

// DeleteProjectProjectNameStageStageNameResourceResourceURIHandlerFunc deletes the specified stage resource
func DeleteProjectProjectNameStageStageNameResourceResourceURIHandlerFunc(params stage_resource.DeleteProjectProjectNameStageStageNameResourceResourceURIParams) middleware.Responder {
	return middleware.NotImplemented("operation stage_resource.DeleteProjectProjectNameStageStageNameResourceResourceURI has not yet been implemented")
}
