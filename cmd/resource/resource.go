package resource

import (
	"errors"
	"fmt"

	"github.com/momentohq/cfn-cache/internal/utility"

	"github.com/aws-cloudformation/cloudformation-cli-go-plugin/cfn/handler"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/momentohq/client-sdk-go/momento"
)

const secretName = "/momento/authToken"

// Create handles the Create event from the Cloudformation service.
func Create(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	client, err := getMomentoClient(req)
	if err != nil {
		return handler.ProgressEvent{}, errors.New("error initializing client")
	}
	err = client.CreateCache(&momento.CreateCacheRequest{
		CacheName: *currentModel.Name,
	})
	if err != nil {
		if momentoErr, ok := err.(momento.MomentoError); ok {
			if momentoErr.Code() != momento.AlreadyExistsError {
				return handleGeneralError(fmt.Sprintf("error occurred creating cache %+v", err))
			} else {
				return handler.ProgressEvent{
					OperationStatus:  handler.Failed,
					HandlerErrorCode: cloudformation.HandlerErrorCodeAlreadyExists,
					Message:          fmt.Sprintf("cache with name %s already exists", *currentModel.Name),
				}, nil
			}
		}
	}
	response := handler.ProgressEvent{
		OperationStatus: handler.Success,
		ResourceModel:   currentModel,
	}

	return response, nil

}

// Read handles the Read event from the Cloudformation service.
func Read(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	client, err := getMomentoClient(req)
	if err != nil {
		return handleGeneralError(fmt.Sprintf("error initializing momento client %+v", err))
	}

	// List caches
	foundCache, err := findCache(client, *currentModel.Name)
	if err != nil {
		return handleGeneralError(fmt.Sprintf("error occurred inspecting your cache %+v", err))
	}

	if !foundCache {
		return handler.ProgressEvent{
			OperationStatus:  handler.Failed,
			Message:          "Cache NotFound",
			HandlerErrorCode: cloudformation.HandlerErrorCodeNotFound,
			ResourceModel:    currentModel,
		}, nil
	}
	return handler.ProgressEvent{
		OperationStatus: handler.Success,
		ResourceModel:   currentModel,
	}, nil
}

// Update handles the Update event from the Cloudformation service.
func Update(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	client, err := getMomentoClient(req)
	if err != nil {
		return handleGeneralError(fmt.Sprintf("error initializing momento client %+v", err))
	}

	// List caches
	foundCache, err := findCache(client, *currentModel.Name)
	if err != nil {
		return handleGeneralError(fmt.Sprintf("error occurred inspecting your cache for update %+v", err))
	}

	if !foundCache {
		return handler.ProgressEvent{
			OperationStatus:  handler.Failed,
			Message:          "Cache NotFound cant update",
			HandlerErrorCode: cloudformation.HandlerErrorCodeNotFound,
			ResourceModel:    currentModel,
		}, nil
	}
	return handler.ProgressEvent{
		OperationStatus: handler.Success,
		ResourceModel:   currentModel,
	}, nil
}

// Delete handles the Delete event from the Cloudformation service.
func Delete(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	client, err := getMomentoClient(req)
	if err != nil {
		return handleGeneralError(fmt.Sprintf("error initializing momento client %+v", err))
	}

	foundCache, err := findCache(client, *currentModel.Name)
	if err != nil {
		return handleGeneralError(fmt.Sprintf("error occurred inspecting your cache %+v", err))
	}

	if !foundCache {
		return handler.ProgressEvent{
			OperationStatus:  handler.Failed,
			Message:          "Cache NotFound",
			HandlerErrorCode: cloudformation.HandlerErrorCodeNotFound,
			ResourceModel:    currentModel,
		}, nil
	}

	err = client.DeleteCache(&momento.DeleteCacheRequest{
		CacheName: *currentModel.Name,
	})
	if err != nil {
		return handleGeneralError(fmt.Sprintf("error occurred deleting cache %+v", err))
	}

	return handler.ProgressEvent{
		OperationStatus: handler.Success,
	}, nil
}

// List handles the List event from the Cloudformation service.
func List(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	client, err := getMomentoClient(req)
	if err != nil {
		return handleGeneralError(fmt.Sprintf("error initializing momento client %+v", err))
	}

	// List caches
	foundCache, err := findCache(client, *currentModel.Name)

	if err != nil {
		return handleGeneralError(fmt.Sprintf("error occurred inspecting your cache %+v", err))
	}

	if !foundCache {
		return handler.ProgressEvent{
			OperationStatus: handler.Success,
			ResourceModels:  []interface{}{}, // Empty list for not found on list
		}, nil
	}
	return handler.ProgressEvent{
		OperationStatus: handler.Success,
		ResourceModels:  []interface{}{currentModel},
	}, nil
}

func getMomentoClient(req handler.Request) (momento.ScsClient, error) {
	authToken, err := utility.GetSecret(secretsmanager.New(req.Session), secretName)
	if err != nil {
		return nil, err
	}
	client, err := momento.NewSimpleCacheClient(
		authToken,
		600, // TTL is irrelevant we are just doing control plane actions in this resource
	)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func findCache(client momento.ScsClient, name string) (bool, error) {
	token := ""
	foundCache := false
	for {
		listCacheResp, err := client.ListCaches(&momento.ListCachesRequest{NextToken: token})
		if err != nil {
			return false, err
		}
		for _, cacheInfo := range listCacheResp.Caches() {
			if cacheInfo.Name() == name {
				foundCache = true
				break
			}
		}
		token = listCacheResp.NextToken()
		if token == "" {
			break
		}
	}
	return foundCache, nil
}

func handleGeneralError(msg string) (handler.ProgressEvent, error) {
	return handler.ProgressEvent{
		OperationStatus:  handler.Failed,
		Message:          msg,
		HandlerErrorCode: cloudformation.HandlerErrorCodeGeneralServiceException,
	}, nil
}
