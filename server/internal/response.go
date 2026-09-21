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
	"encoding/json"
	"fmt"
	"log"
	"sort"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"google.golang.org/protobuf/types/known/durationpb"
)

// getResponse converts an adapter Response into a GetPageResponse.
func getResponse(
	reverseMapping *entityReverseIdMapping,
	resp *framework.Response,
) (rpcResponse *api_adapter_v1.GetPageResponse) {
	if resp == nil {
		return api_adapter_v1.NewGetPageResponseError(&api_adapter_v1.Error{
			Message: "Adapter returned nil response. This is always indicative of a bug within the Adapter implementation.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		})
	}

	if resp.Error != nil {
		err := &api_adapter_v1.Error{
			Message: resp.Error.Message,
			Code:    resp.Error.Code,
		}

		if resp.Error.RetryAfter != nil {
			err.RetryAfter = durationpb.New(*resp.Error.RetryAfter)
		}

		return api_adapter_v1.NewGetPageResponseError(err)
	}

	if resp.Success == nil {
		return api_adapter_v1.NewGetPageResponseError(&api_adapter_v1.Error{
			Message: "Adapter returned empty response. This is always indicative of a bug within the Adapter implementation.",
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
) (entityObjects *api_adapter_v1.EntityObjects, adapterErr *api_adapter_v1.Error) {
	entityObjects = &api_adapter_v1.EntityObjects{
		EntityId: reverseMapping.Id,
	}

	if len(objects) > 0 {
		entityObjects.Objects = make([]*api_adapter_v1.Object, 0, len(objects))

		for _, object := range objects {
			var entityObject *api_adapter_v1.Object
			entityObject, adapterErr = getEntityObject(reverseMapping, object)

			if adapterErr != nil {
				// Log the complete object and skip it instead of returning error
				objectJSON, jsonErr := json.Marshal(object)
				if jsonErr != nil {
					log.Printf("[ERROR] Failed to marshal object for logging: %v. Original error: %s", jsonErr, adapterErr.Message)
				} else {
					log.Printf("[ERROR] Skipping object due to validation error. Entity: %s, Error: %s, Object: %s",
						reverseMapping.Id, adapterErr.Message, string(objectJSON))
				}

				continue // Skip this object and process the next one
			}

			entityObjects.Objects = append(entityObjects.Objects, entityObject)
		}

	}

	return
}

// getEntityObject converts an adapter Object into an RPC object.
func getEntityObject(
	reverseMapping *entityReverseIdMapping,
	object framework.Object,
) (entityObject *api_adapter_v1.Object, adapterErr *api_adapter_v1.Error) {
	entityObject = new(api_adapter_v1.Object)

	// Iterate over the sorted externalIds, in order to always return
	// attributes and child objects in the same order.
	sortedExternalIds := make([]string, 0, len(object))
	for externalId := range object {
		sortedExternalIds = append(sortedExternalIds, externalId)
	}
	sort.Strings(sortedExternalIds)

	for _, externalId := range sortedExternalIds {
		value := object[externalId]
		switch v := value.(type) {

		case []framework.Object: // Child objects for a child entity.
			childReverseMapping, validExternalId := reverseMapping.ChildEntities[externalId]

			if !validExternalId {
				adapterErr = &api_adapter_v1.Error{
					Message: fmt.Sprintf("Adapter returned an object for entity %s which contains child objects with an invalid entity external ID: %s. This is always indicative of a bug within the Adapter implementation.", reverseMapping.Id, externalId),
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				}

				return nil, adapterErr
			}

			// As an optimization, ignore the child entity if there are no
			// objects to return.
			if len(v) == 0 {
				continue
			}

			var childEntityObjects *api_adapter_v1.EntityObjects
			childEntityObjects, adapterErr = getEntityObjects(childReverseMapping, v)

			if adapterErr != nil {
				return nil, adapterErr
			}

			entityObject.ChildObjects = append(entityObject.ChildObjects, childEntityObjects)

		default: // Attribute.
			attributeMetadata, validExternalId := reverseMapping.Attributes[externalId]

			if !validExternalId {
				adapterErr = &api_adapter_v1.Error{
					Message: fmt.Sprintf("Adapter returned an object for entity %s which contains an attribute with an invalid external ID: %s. This is always indicative of a bug within the Adapter implementation.", reverseMapping.Id, externalId),
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				}

				return nil, adapterErr
			}

			adapterErr = validateAttributeValue(attributeMetadata, value)

			if adapterErr != nil {
				return nil, adapterErr
			}

			var values []*api_adapter_v1.AttributeValue
			values, adapterErr = getAttributeValues(value)

			if adapterErr != nil {
				return nil, adapterErr
			}

			// As an optimization, ignore the attribute value if it is null, as
			// it is equivalent to returning null.
			if values == nil {
				continue
			}

			entityObject.Attributes = append(entityObject.Attributes, &api_adapter_v1.Attribute{
				Id:     attributeMetadata.Id,
				Values: values,
			})
		}
	}

	if len(entityObject.Attributes) == 0 {
		adapterErr = &api_adapter_v1.Error{
			Message: fmt.Sprintf("Adapter returned an object for entity %s which contains no non-null attributes. This is always indicative of a bug within the Adapter implementation.", reverseMapping.Id),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}

		return nil, adapterErr
	}

	return
}
