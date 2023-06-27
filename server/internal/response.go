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
	"fmt"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
)

// getResponse converts an adapter Response into a GetPageResponse.
func getResponse(
	reverseMapping *entityReverseIdMapping,
	resp *framework.Response,
) (rpcResponse *api_adapter_v1.GetPageResponse) {
	if resp.Error != nil {
		return api_adapter_v1.NewGetPageResponseError(&api_adapter_v1.Error{
			Message: resp.Error.Message,
			Code:    resp.Error.Code,
			// TODO: RetryAfter
		})
	}

	if resp.Success == nil {
		return api_adapter_v1.NewGetPageResponseError(&api_adapter_v1.Error{
			Message: "Adapter returned an empty response. This is always indicative of a bug within the Adapter implementation.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		})
	}

	entityObjects, adapterErr := getEntityObjects(reverseMapping, resp.Success.Objects)

	if adapterErr != nil {
		return api_adapter_v1.NewGetPageResponseError(adapterErr)
	}

	page := &api_adapter_v1.Page{
		NextCursor: resp.Success.NextCursor,
		Objects:    entityObjects.Objects,
	}

	return api_adapter_v1.NewGetPageResponseSuccess(page)
}

// getEntityObjects converts an adapter list of objects for an entity into an
// EntityObject.
func getEntityObjects(
	reverseMapping *entityReverseIdMapping,
	objects []framework.Object,
) (entityObjects api_adapter_v1.EntityObjects, adapterErr *api_adapter_v1.Error) {
	entityObjects.EntityId = reverseMapping.Id

	if len(objects) > 0 {
		entityObjects.Objects = make([]*api_adapter_v1.Object, 0, len(objects))

		for _, object := range objects {
			var entityObject api_adapter_v1.Object
			entityObject, adapterErr = getEntityObject(reverseMapping, object)

			if adapterErr != nil {
				return
			}

			entityObjects.Objects = append(entityObjects.Objects, &entityObject)
		}

	}

	return
}

// getEntityObject converts an adapter Object into an RPC object.
func getEntityObject(
	reverseMapping *entityReverseIdMapping,
	object framework.Object,
) (entityObject api_adapter_v1.Object, adapterErr *api_adapter_v1.Error) {
	for externalId, value := range object {
		switch v := value.(type) {

		case []framework.Object: // Child objects for a child entity.
			childReverseMapping, validExternalId := reverseMapping.ChildEntities[externalId]

			if !validExternalId {
				adapterErr = &api_adapter_v1.Error{
					Message: fmt.Sprintf("Adapter returned child objects with an invalid entity external ID: %s. This is always indicative of a bug within the Adapter implementation.", externalId),
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				}

				return
			}

			var childEntityObjects api_adapter_v1.EntityObjects

			childEntityObjects, adapterErr = getEntityObjects(childReverseMapping, v)

			if adapterErr != nil {
				return
			}

			entityObject.ChildObjects = append(entityObject.ChildObjects, &childEntityObjects)

		default: // Attribute.
			attributeMetadata, validExternalId := reverseMapping.Attributes[externalId]

			if !validExternalId {
				adapterErr = &api_adapter_v1.Error{
					Message: fmt.Sprintf("Adapter returned an attribute with an invalid external ID: %s. This is always indicative of a bug within the Adapter implementation.", externalId),
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				}

				return
			}

			var values []*api_adapter_v1.AttributeValue

			// TODO: Validate that the value has the correct type re: attributeMetadata.
			values, adapterErr = getAttributeValues(value)

			if adapterErr != nil {
				return
			}
			entityObject.Attributes = append(entityObject.Attributes, &api_adapter_v1.Attribute{
				Id:     attributeMetadata.Id,
				Values: values,
			})
		}
	}

	return
}
