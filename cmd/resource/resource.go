package resource

import (
	"context"
	"fmt"
	"os"

	"github.com/aws-cloudformation/cloudformation-cli-go-plugin/cfn/handler"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/momento"
	"github.com/momentohq/client-sdk-go/responses"
)

// Create handles the Create event from the Cloudformation service.
func Create(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	client, err := getMomentoClient(currentModel)
	if err != nil {
		return handler.NewFailedEvent(fmt.Errorf("error initializing client %w", err)), nil
	}
	rsp, err := client.CreateCache(context.Background(), &momento.CreateCacheRequest{
		CacheName: *currentModel.Name,
	})
	if err != nil {
		return handleGeneralError(fmt.Sprintf("error occurred creating cache %+v", err))
	}
	switch rsp.(type) {
	case *responses.CreateCacheAlreadyExists:
		return handler.ProgressEvent{
			OperationStatus:  handler.Failed,
			HandlerErrorCode: cloudformation.HandlerErrorCodeAlreadyExists,
			Message:          fmt.Sprintf("cache with name %s already exists", *currentModel.Name),
		}, nil
	case *responses.CreateCacheSuccess:
		return handler.ProgressEvent{
			OperationStatus: handler.Success,
			ResourceModel:   currentModel,
		}, nil
	default:
		return handleGeneralError(fmt.Sprintf("unexpected response type from create cache api request %T", rsp))
	}
}

// Read handles the Read event from the Cloudformation service.
func Read(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	client, err := getMomentoClient(currentModel)
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
	client, err := getMomentoClient(currentModel)
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
	client, err := getMomentoClient(currentModel)
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

	_, err = client.DeleteCache(context.Background(), &momento.DeleteCacheRequest{
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
	client, err := getMomentoClient(currentModel)
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

func getMomentoClient(currentModel *Model) (momento.CacheClient, error) {
	var credProvider auth.CredentialProvider
	var err error
	if os.Getenv("MODE") == "TEST" {
		credProvider, err = auth.NewEnvMomentoTokenProvider("MOMENTO_AUTH_TOKEN")
	} else {
		credProvider, err = auth.NewStringMomentoTokenProvider(*currentModel.AuthToken)
	}
	if err != nil {
		return nil, err
	}
	return momento.NewCacheClient(config.LaptopLatest(), credProvider, 1)
}

func findCache(client momento.CacheClient, name string) (bool, error) {
	token := ""
	foundCache := false
	for {
		listCacheResp, err := client.ListCaches(context.Background(), &momento.ListCachesRequest{NextToken: token})
		if err != nil {
			return false, err
		}
		if r, ok := listCacheResp.(*responses.ListCachesSuccess); ok {
			r.Caches()
			for _, cacheInfo := range r.Caches() {
				if cacheInfo.Name() == name {
					foundCache = true
					break
				}
			}
			token = r.NextToken()
			if token == "" {
				break
			}
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
