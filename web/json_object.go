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

package web

import (
	"fmt"

	framework "github.com/sgnl-ai/adapter-framework"
)

// ConvertJSONObjectList parses and converts a list of JSON objects received
// from the given requested entity.
func ConvertJSONObjectList(entity *framework.EntityConfig, objects []map[string]any, opts ...JSONOption) ([]framework.Object, error) {
	options := defaultJSONOptions()
	for _, opt := range opts {
		opt.apply(options)
	}
	return convertJSONObjectList(entity, objects, options)
}

// convertJSONObjectList parses and converts a list of JSON objects received
// from the given requested entity.
func convertJSONObjectList(entity *framework.EntityConfig, objects []map[string]any, opts *jsonOptions) ([]framework.Object, error) {
	if len(objects) == 0 {
		return nil, nil
	}

	parsedObjects := make([]framework.Object, 0, len(objects))

	for _, object := range objects {
		parsedObject, err := convertJSONObject(entity, object, opts)

		if err != nil {
			return nil, err
		}

		parsedObjects = append(parsedObjects, parsedObject)
	}

	return parsedObjects, nil
}

// convertJSONObject parses and converts a JSON object received from the given
//requested entity.
//
// If
func convertJSONObject(entity *framework.EntityConfig, object map[string]any, opts *jsonOptions) (framework.Object, error) {
	parsedObject := make(framework.Object)

	// Parse attributes.
	for _, attribute := range entity.Attributes {
		externalId := attribute.ExternalId

		// TODO: Support flattening single-valued complex attributes,
		// e.g. flatten this JSON object:
		//
		// {
		//   "manager": {
		//     "id": "1234",
		//     "email": "john@example.com"
		//   }
		// }
		//
		// into:
		//
		// {
		//   "manager__id": "1234",
		//   "manager__email": "john@example.com"
		// }

		value, found := object[externalId]
		if !found {
			continue
		}

		parsedValue, err := convertJSONAttributeValue(attribute, value, opts)
		if err != nil {
			return nil, err
		}

		// Do not return null attribute values.
		if parsedValue == nil {
			continue
		}

		parsedObject[externalId] = parsedValue
	}

	// Parse child entities.
	for _, childEntity := range entity.ChildEntities {
		externalId := childEntity.ExternalId

		childObjectsRaw, found := object[externalId]
		if !found {
			continue
		}

		childObjects, ok := childObjectsRaw.([]map[string]any)
		if !ok {
			return nil, fmt.Errorf("child entity %s is not associated with a list of JSON objects", externalId)
		}

		if len(childObjects) == 0 {
			continue
		}

		parsedChildObjects, err := convertJSONObjectList(childEntity, childObjects, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to parse objects for child entity %s: %w", externalId, err)
		}

		parsedObject[externalId] = parsedChildObjects
	}

	return parsedObject, nil
}
