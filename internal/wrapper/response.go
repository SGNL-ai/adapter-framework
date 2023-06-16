package wrapper

import (
	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
)

// getResponse converts an adapter Response into a GetPageResponse.
func getResponse(
	reverseMapping entityReverseIdMapping,
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
			Message: api_adapter_v1.ErrMsgAdapterEmptyResponse,
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
	reverseMapping entityReverseIdMapping,
	objects []*framework.Object,
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
	reverseMapping entityReverseIdMapping,
	object *framework.Object,
) (entityObject api_adapter_v1.Object, adapterErr *api_adapter_v1.Error) {
	if len(object.Attributes) > 0 {
		entityObject.Attributes = make([]*api_adapter_v1.Attribute, 0, len(object.Attributes))

		for attributeExternalId, attributeValue := range object.Attributes {
			attributeMetadata, validExternalId := reverseMapping.Attributes[attributeExternalId]

			if !validExternalId {
				adapterErr = &api_adapter_v1.Error{
					Message: api_adapter_v1.ErrMsgAdapterInvalidEntityExternalId,
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				}

				return
			}

			var values []*api_adapter_v1.AttributeValue

			values, adapterErr = getAttributeValues(attributeMetadata, attributeValue)

			if adapterErr != nil {
				return
			}

			entityObject.Attributes = append(entityObject.Attributes, &api_adapter_v1.Attribute{
				Id:     attributeMetadata.Id,
				Values: values,
			})
		}
	}

	if len(object.Children) > 0 {
		entityObject.ChildObjects = make([]*api_adapter_v1.EntityObjects, 0, len(object.Children))

		for entityExternalId, childObjects := range object.Children {
			childReverseMapping, validExternalId := reverseMapping.ChildEntities[entityExternalId]

			if !validExternalId {
				adapterErr = &api_adapter_v1.Error{
					Message: api_adapter_v1.ErrMsgAdapterInvalidEntityExternalId,
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				}

				return
			}

			var childEntityObjects api_adapter_v1.EntityObjects

			childEntityObjects, adapterErr = getEntityObjects(childReverseMapping, childObjects)

			if adapterErr != nil {
				return
			}

			entityObject.ChildObjects = append(entityObject.ChildObjects, &childEntityObjects)
		}
	}

	return
}
