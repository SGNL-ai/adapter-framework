// Copyright 2023 SGNL.ai, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package internal

import (
	"context"
	"fmt"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	grpcMetadata "google.golang.org/grpc/metadata"
)

// entityReverseMapping maps external IDs to IDs.
type entityReverseIdMapping struct {
	// Id is the entity's ID.
	Id string

	// AttributeIds maps the entity's attributes' external IDs to their IDs
	// and types.
	Attributes map[string]*api_adapter_v1.AttributeConfig

	// ChildEntities maps the entity's child entities' external IDs to their
	// entityReverseIdMapping.
	ChildEntities map[string]*entityReverseIdMapping
}

// getAdapterRequest converts a GetPageRequest into an adapter Request.
func getAdapterRequest[Config any](
	ctx context.Context,
	req *api_adapter_v1.GetPageRequest,
) (adapterRequest *framework.Request[Config], reverseMapping *entityReverseIdMapping, adapterErr *api_adapter_v1.Error) {
	var errMsg string

	switch {
	case req == nil:
		errMsg = "Request is nil."
	case req.Datasource == nil:
		errMsg = "Request contains no datasource config."
	case req.Entity == nil:
		errMsg = "Request contains no entity config."
	case req.PageSize <= 0:
		errMsg = fmt.Sprintf("Request contains an invalid page size: %d. Must be greater than 0.", req.PageSize)
	}

	if errMsg != "" {
		adapterErr = &api_adapter_v1.Error{
			Message: errMsg,
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
		}

		return nil, nil, adapterErr
	}

	// TODO [sc-11900]: Move this check to the start of the function and remove the check for
	// `Type != ""` in the if statement below once this should be enforced on all adapters.
	// Currently, this will only be enforced on adapters that have been upgraded to use discrete
	// adapters (e.g. they have `Type` set).
	if req.Datasource.Type != "" && !canAccessAdapter(ctx) {
		adapterErr = &api_adapter_v1.Error{
			Message: "Forbidden.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_FORBIDDEN,
		}

		return nil, nil, adapterErr
	}

	switch {
	case req.Datasource.Id == "":
		errMsg = "Datasource config contains no ID."
	}

	if errMsg != "" {
		adapterErr = &api_adapter_v1.Error{
			Message: errMsg,
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}

		return nil, nil, adapterErr
	}

	adapterRequest = &framework.Request[Config]{}

	if len(req.Datasource.Config) > 0 {
		config, err := ParseConfig[Config](req.Datasource.Config)

		if err != nil {
			adapterErr = &api_adapter_v1.Error{
				Message: fmt.Sprintf("Config in datasource config could not parsed as JSON: %s.", err),
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			}

			return nil, nil, adapterErr
		}

		adapterRequest.Config = config
	}

	var entityConfig *framework.EntityConfig

	entityConfig, reverseMapping, adapterErr = getEntity(req.Entity)

	if adapterErr != nil {
		return nil, nil, adapterErr
	}

	adapterRequest.Address = req.Datasource.Address
	adapterRequest.Auth = getAdapterAuth(req.Datasource.Auth)
	adapterRequest.Entity = *entityConfig
	adapterRequest.Ordered = req.Entity.Ordered
	adapterRequest.PageSize = req.PageSize
	adapterRequest.Cursor = req.Cursor
	adapterRequest.Type = req.Datasource.Type

	return
}

// canAccessAdapter verifies the request has the correct token to access the
// adapter. Will return true if the provided token matches any of the tokens
// specified in TODO. Otherwise, will return false.
func canAccessAdapter(ctx context.Context) bool {
	metadata, ok := grpcMetadata.FromIncomingContext(ctx)
	if !ok {
		return false
	}

	tokens := metadata.Get("token")
	if len(tokens) != 1 {
		return false
	}

	// TODO: Get from file
	x := []string{"TODO GET FROM FILE"}

	// TODO: Once upgrading go to 1.21+, replace with the `Contains` method
	for _, y := range x {
		if y == tokens[0] {
			return true
		}
	}

	return false
}

// getAdapterAuth converts a request DatasourceAuthCredentials into an adapter
// DatasourceAuthCredentials.
func getAdapterAuth(
	auth *api_adapter_v1.DatasourceAuthCredentials,
) *framework.DatasourceAuthCredentials {
	if auth == nil || auth.AuthMechanism == nil {
		return nil
	}

	switch a := auth.AuthMechanism.(type) {
	case *api_adapter_v1.DatasourceAuthCredentials_Basic_:
		return &framework.DatasourceAuthCredentials{
			Basic: &framework.BasicAuthCredentials{
				Username: a.Basic.Username,
				Password: a.Basic.Password,
			},
		}

	case *api_adapter_v1.DatasourceAuthCredentials_HttpAuthorization:
		return &framework.DatasourceAuthCredentials{
			HTTPAuthorization: a.HttpAuthorization,
		}

	default:
		return nil
	}
}

// getEntity converts a request EntityConfig into an adapter EntityConfig.
func getEntity(
	entity *api_adapter_v1.EntityConfig,
) (adapterEntity *framework.EntityConfig, reverseMapping *entityReverseIdMapping, adapterErr *api_adapter_v1.Error) {
	var errMsg string

	switch {
	case entity.Id == "":
		errMsg = "Entity config contains no ID."
	case entity.ExternalId == "":
		errMsg = fmt.Sprintf("Entity config %s contains no external ID.", entity.Id)
	case len(entity.Attributes) == 0:
		errMsg = fmt.Sprintf("Entity config %s (%s) contains no attributes.", entity.Id, entity.ExternalId)
	}

	if errMsg != "" {
		adapterErr = &api_adapter_v1.Error{
			Message: errMsg,
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
		}

		return nil, nil, adapterErr
	}

	adapterEntity = &framework.EntityConfig{}
	reverseMapping = &entityReverseIdMapping{}

	adapterEntity.ExternalId = entity.ExternalId
	reverseMapping.Id = entity.Id

	attributeIds := make(map[string]bool, len(entity.Attributes))
	adapterEntity.Attributes = make([]*framework.AttributeConfig, 0, len(entity.Attributes))
	reverseMapping.Attributes = make(map[string]*api_adapter_v1.AttributeConfig, len(entity.Attributes))
	for _, attribute := range entity.Attributes {
		switch {
		case attribute.Id == "":
			errMsg = fmt.Sprintf("Attribute in entity config %s (%s) contains no ID.", entity.Id, entity.ExternalId)
		case attribute.ExternalId == "":
			errMsg = fmt.Sprintf("Attribute in entity config %s (%s) contains no external ID: %s.", entity.Id, entity.ExternalId, attribute.Id)
		case attributeIds[attribute.Id]:
			errMsg = fmt.Sprintf("Attribute in entity config %s (%s) contains duplicate ID: %s (%s).", entity.Id, entity.ExternalId, attribute.Id, attribute.ExternalId)
		case reverseMapping.Attributes[attribute.ExternalId] != nil:
			errMsg = fmt.Sprintf("Attribute in entity config %s (%s) contains duplicate external ID: %s (%s).", entity.Id, entity.ExternalId, attribute.Id, attribute.ExternalId)
		case attribute.Type == api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_UNSPECIFIED:
			errMsg = fmt.Sprintf("Attribute in entity config %s (%s) contains unspecified type %s: %s (%s).", entity.Id, entity.ExternalId, attribute.Type, attribute.Id, attribute.ExternalId)
		case api_adapter_v1.AttributeType_name[int32(attribute.Type)] == "":
			errMsg = fmt.Sprintf("Attribute in entity config %s (%s) contains invalid type %s: %s (%s).", entity.Id, entity.ExternalId, attribute.Type, attribute.Id, attribute.ExternalId)
		}

		if errMsg != "" {
			adapterErr = &api_adapter_v1.Error{
				Message: errMsg,
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			}

			return nil, nil, adapterErr
		}

		attributeIds[attribute.Id] = true

		adapterEntity.Attributes = append(adapterEntity.Attributes, &framework.AttributeConfig{
			ExternalId: attribute.ExternalId,
			Type:       framework.AttributeType(attribute.Type),
			List:       attribute.List,
		})

		reverseMapping.Attributes[attribute.ExternalId] = attribute
	}

	if len(entity.ChildEntities) > 0 {
		childEntityIds := make(map[string]bool, len(entity.ChildEntities))
		adapterEntity.ChildEntities = make([]*framework.EntityConfig, 0, len(entity.ChildEntities))
		reverseMapping.ChildEntities = make(map[string]*entityReverseIdMapping, len(entity.ChildEntities))

		for _, childEntity := range entity.ChildEntities {
			switch {
			case childEntity.Id == "":
				errMsg = fmt.Sprintf("Child entity in entity config %s (%s) contains no ID.", entity.Id, entity.ExternalId)
			case childEntity.ExternalId == "":
				errMsg = fmt.Sprintf("Child entity in entity config %s (%s) contains no external ID: %s.", entity.Id, entity.ExternalId, childEntity.Id)
			case childEntityIds[childEntity.Id]:
				errMsg = fmt.Sprintf("Child entity in entity config %s (%s) contains duplicate ID: %s (%s).", entity.Id, entity.ExternalId, childEntity.Id, childEntity.ExternalId)
			case reverseMapping.ChildEntities[childEntity.ExternalId] != nil:
				errMsg = fmt.Sprintf("Child entity in entity config %s (%s) contains duplicate external ID: %s (%s).", entity.Id, entity.ExternalId, childEntity.Id, childEntity.ExternalId)
			case reverseMapping.Attributes[childEntity.ExternalId] != nil:
				errMsg = fmt.Sprintf("Child entity in entity config %s (%s) contains same external ID as attribute: %s (%s).", entity.Id, entity.ExternalId, childEntity.Id, childEntity.ExternalId)
			}

			if errMsg != "" {
				adapterErr = &api_adapter_v1.Error{
					Message: errMsg,
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
				}

				return nil, nil, adapterErr
			}

			childEntityIds[childEntity.Id] = true

			var adapterChildEntity *framework.EntityConfig
			var childReverseMapping *entityReverseIdMapping

			adapterChildEntity, childReverseMapping, adapterErr = getEntity(childEntity)

			if adapterErr != nil {
				return nil, nil, adapterErr
			}

			adapterEntity.ChildEntities = append(adapterEntity.ChildEntities, adapterChildEntity)
			reverseMapping.ChildEntities[childEntity.ExternalId] = childReverseMapping
		}
	}

	return
}
