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
	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
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
	req *api_adapter_v1.GetPageRequest,
) (adapterRequest *framework.Request[Config], reverseMapping *entityReverseIdMapping, adapterErr *api_adapter_v1.Error) {
	var errMsg string

	switch {
	case req == nil:
		errMsg = api_adapter_v1.ErrorNoDatasourceConfig
	case req.Datasource == nil:
		errMsg = api_adapter_v1.ErrorNoDatasourceConfig
	case req.Entity == nil:
		errMsg = api_adapter_v1.ErrorNoEntityConfig
	case req.PageSize <= 0:
		errMsg = api_adapter_v1.ErrorMsgEntityPageSizeTooSmall
	}

	if errMsg != "" {
		adapterErr = &api_adapter_v1.Error{
			Message: errMsg,
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
		}

		return
	}

	adapterRequest = &framework.Request[Config]{}

	if len(req.Datasource.Config) > 0 {
		config, err := ParseConfig[Config](req.Datasource.Config)

		if err != nil {
			return
		}

		adapterRequest.Config = config
	}

	var entityConfig *framework.EntityConfig

	entityConfig, reverseMapping, adapterErr = getEntity(req.Entity)

	if adapterErr != nil {
		return
	}

	adapterRequest.Address = req.Datasource.Address
	adapterRequest.Auth = getAdapterAuth(req.Datasource.Auth)
	adapterRequest.Entity = *entityConfig
	adapterRequest.Ordered = req.Entity.Ordered
	adapterRequest.PageSize = req.PageSize
	adapterRequest.Cursor = req.Cursor

	return
}

// getAdapterAuth converts a request DatasourceAuthCredentials into an adapter
// DatasourceAuthCredentials.
func getAdapterAuth(
	auth *api_adapter_v1.DatasourceAuthCredentials,
) *framework.DatasourceAuthCredentials {
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
		errMsg = api_adapter_v1.ErrorMsgNoIdProvidedInEntity
	case entity.ExternalId == "":
		errMsg = api_adapter_v1.ErrorMsgNoExternalIdProvidedInEntity
	case len(entity.Attributes) == 0:
		errMsg = api_adapter_v1.ErrorMsgEntityProvidedWithNoAttributes
	}

	if errMsg != "" {
		adapterErr = &api_adapter_v1.Error{
			Message: errMsg,
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
		}

		return
	}

	reverseMapping = &entityReverseIdMapping{}
	adapterEntity = &framework.EntityConfig{}

	reverseMapping.Id = entity.Id
	adapterEntity.ExternalId = entity.ExternalId

	reverseMapping.Attributes = make(map[string]*api_adapter_v1.AttributeConfig, len(entity.Attributes))
	adapterEntity.Attributes = make([]*framework.AttributeConfig, 0, len(entity.Attributes))
	for _, attribute := range entity.Attributes {
		switch {
		case attribute.Id == "":
			errMsg = api_adapter_v1.ErrorMsgNoIdProvidedInAttribute
		case attribute.ExternalId == "":
			errMsg = api_adapter_v1.ErrorMsgNoExternalIdProvidedInAttribute
		case api_adapter_v1.AttributeType_name[int32(attribute.Type)] == "":
			errMsg = api_adapter_v1.ErrorMsgAttributeInvalidType
		case attribute.Type == api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_UNSPECIFIED:
			errMsg = api_adapter_v1.ErrorMsgAttributeInvalidType
		case reverseMapping.Attributes[attribute.ExternalId] != nil:
			errMsg = api_adapter_v1.ErrorMsgEntityDuplicateAttributeExternalId
		}

		if errMsg != "" {
			adapterErr = &api_adapter_v1.Error{
				Message: errMsg,
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			}

			return
		}

		reverseMapping.Attributes[attribute.ExternalId] = attribute

		adapterEntity.Attributes = append(adapterEntity.Attributes, &framework.AttributeConfig{
			ExternalId: attribute.ExternalId,
			Type:       framework.AttributeType(attribute.Type),
			List:       attribute.List,
		})
	}

	if len(entity.ChildEntities) > 0 {
		reverseMapping.ChildEntities = make(map[string]*entityReverseIdMapping, len(entity.ChildEntities))
		adapterEntity.ChildEntities = make([]*framework.EntityConfig, 0, len(entity.ChildEntities))

		for _, childEntity := range entity.ChildEntities {
			switch {
			case reverseMapping.ChildEntities[childEntity.ExternalId] != nil:
				errMsg = api_adapter_v1.ErrorMsgEntityDuplicateChildEntityExternalId
			case reverseMapping.Attributes[childEntity.ExternalId] != nil:
				errMsg = api_adapter_v1.ErrorMsgEntityChildEntityExternalIdSameAsAttribute
			}

			if errMsg != "" {
				adapterErr = &api_adapter_v1.Error{
					Message: errMsg,
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
				}

				return
			}

			var adapterChildEntity *framework.EntityConfig
			var childReverseMapping *entityReverseIdMapping

			adapterChildEntity, childReverseMapping, adapterErr = getEntity(childEntity)

			if adapterErr != nil {
				return
			}

			reverseMapping.ChildEntities[childEntity.ExternalId] = childReverseMapping
			adapterEntity.ChildEntities = append(adapterEntity.ChildEntities, adapterChildEntity)
		}
	}

	return
}
