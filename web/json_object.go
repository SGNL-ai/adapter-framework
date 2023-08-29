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
	"context"
	"fmt"
	"strings"

	"github.com/PaesslerAG/gval"
	"github.com/PaesslerAG/jsonpath"
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

	// Parse and validate all JSONPaths.
	// Create a map indexed by the paths themselves.
	var jsonPaths map[string]gval.Evaluable
	if opts.enableJSONPath {
		jsonPaths = make(map[string]gval.Evaluable)
		if err := parseJSONPaths(entity, jsonPaths); err != nil {
			return nil, err
		}
	}

	parsedObjects := make([]framework.Object, 0, len(objects))

	for _, object := range objects {
		if object == nil {
			continue
		}

		parsedObject, err := convertJSONObject(entity, object, opts, jsonPaths)

		if err != nil {
			return nil, err
		}

		if len(parsedObject) == 0 {
			continue
		}

		parsedObjects = append(parsedObjects, parsedObject)
	}

	return parsedObjects, nil
}

// parseJSONPaths parses all JSONPaths from attribute external IDs starting
// with '$' into the given map.
func parseJSONPaths(entity *framework.EntityConfig, out map[string]gval.Evaluable) error {
	for _, attribute := range entity.Attributes {
		externalId := attribute.ExternalId
		if externalId[0] == '$' {
			if _, found := out[externalId]; found { // Already parsed.
				continue
			}

			jsonPath, err := jsonpath.New(externalId)
			if err != nil {
				return fmt.Errorf("failed to parse attribute external ID %q as JSONPath: %w", externalId, err)
			}

			out[externalId] = jsonPath
		}
	}

	// Recursively parse JSONPaths in child entities' attributes.
	for _, childEntity := range entity.ChildEntities {
		externalId := childEntity.ExternalId

		if externalId[0] == '$' {
			if _, found := out[externalId]; found { // Already parsed.
				continue
			}

			jsonPath, err := jsonpath.New(externalId)
			if err != nil {
				return fmt.Errorf("failed to parse child entity external ID %q as JSONPath: %w", externalId, err)
			}

			out[externalId] = jsonPath
		}

		if err := parseJSONPaths(childEntity, out); err != nil {
			return err
		}
	}

	return nil
}

// convertJSONObject parses and converts a JSON object received from the given
// requested entity.
func convertJSONObject(entity *framework.EntityConfig, object map[string]any, opts *jsonOptions, jsonPaths map[string]gval.Evaluable) (framework.Object, error) {
	parsedObject := make(framework.Object)

	// If the flattening of single-valued complex attributes is enabled,
	// parse single-valued complex attributes that are required to parse
	// attributes and child objects, recursively.
	var complexAttributes map[string]framework.Object
	if opts.complexAttributeNameDelimiter != "" {
		// Map of each single-valued complex attribute external ID to a
		// pseudo entity config that can be used to parse that attribute.
		complexAttributeFakeEntities := make(map[string]*framework.EntityConfig)

		// Identify single-valued complex attributes needed to map attributes
		// with delimiters in their external IDs.
		for _, attribute := range entity.Attributes {
			externalId := attribute.ExternalId

			if _, found := jsonPaths[externalId]; found {
				// The attribute external ID is a JSONPath.
				// Ignore it here.
				continue
			}

			externalIdComponents := strings.SplitN(externalId, opts.complexAttributeNameDelimiter, 2)

			if len(externalIdComponents) != 2 {
				// The external ID doesn't contain the delimiter, so it doesn't
				// identify a single-valued complex attribute. Ignore it.
				continue
			}

			localExternalId := externalIdComponents[0]
			subExternalId := externalIdComponents[1]

			rawValue, found := object[localExternalId]
			if !found || rawValue == nil {
				continue
			}

			var fakeEntity *framework.EntityConfig
			fakeEntity, wasCached := complexAttributeFakeEntities[localExternalId]
			if !wasCached {
				fakeEntity = &framework.EntityConfig{}
			}
			fakeEntity.Attributes = append(fakeEntity.Attributes, &framework.AttributeConfig{
				ExternalId: subExternalId,
				Type:       attribute.Type,
				List:       attribute.List,
			})
			complexAttributeFakeEntities[localExternalId] = fakeEntity
		}

		// Identify single-valued complex attributes needed to map child entities.
		for _, childEntity := range entity.ChildEntities {
			externalId := childEntity.ExternalId

			if _, found := jsonPaths[externalId]; found {
				// The child entity external ID is a JSONPath.
				// Ignore it here.
				continue
			}

			externalIdComponents := strings.SplitN(externalId, opts.complexAttributeNameDelimiter, 2)

			if len(externalIdComponents) != 2 {
				// The external ID doesn't contain the delimiter, so it doesn't
				// identify a single-valued complex attribute. Ignore it.
				continue
			}

			localExternalId := externalIdComponents[0]
			subExternalId := externalIdComponents[1]

			_, found := object[localExternalId]
			if !found {
				continue
			}

			var fakeEntity *framework.EntityConfig
			fakeEntity, wasCached := complexAttributeFakeEntities[localExternalId]
			if !wasCached {
				fakeEntity = &framework.EntityConfig{}
			}
			fakeEntity.ChildEntities = append(fakeEntity.ChildEntities, &framework.EntityConfig{
				ExternalId:    subExternalId,
				Attributes:    childEntity.Attributes,
				ChildEntities: childEntity.ChildEntities,
			})
			complexAttributeFakeEntities[localExternalId] = fakeEntity
		}

		complexAttributes = make(map[string]framework.Object, len(complexAttributeFakeEntities))

		for localExternalId, fakeEntity := range complexAttributeFakeEntities {
			complexValue, ok := object[localExternalId].(map[string]any)
			if !ok {
				return nil, fmt.Errorf("attribute %s is not a single-valued complex attribute", localExternalId)
			}

			parsedComplexValue, err := convertJSONObject(fakeEntity, complexValue, opts, jsonPaths)
			if err != nil {
				return nil, fmt.Errorf("failed to parse single-valued complex attribute %s: %w", localExternalId, err)
			}

			complexAttributes[localExternalId] = parsedComplexValue
		}
	}

	// Parse attributes.
	for _, attribute := range entity.Attributes {
		externalId := attribute.ExternalId

		var parsedValue any

		if jsonPath, found := jsonPaths[externalId]; found { // The attribute external ID is a JSONPath.
			value, err := jsonPath(context.Background(), object)
			if err != nil { // The JSONPath didn't match. Evaluate to null.
				continue
			}

			// A JSONPath that matches multi-valued complex attributes, such as
			// $.emails[?(@.primary==true)].value, always returns a list, even
			// when that list contains only one value.
			// In the case the attribute is configured with a non-list type and
			// the JSONPath evaluation returned a list, convert just that one
			// value.
			if !attribute.List {
				if rawList, isList := value.([]any); isList {
					switch len(rawList) {
					case 0: // Nothing matched? Ignore.
						continue
					case 1:
						value = rawList[0]
					default:
						return nil, fmt.Errorf("non-list attribute %s matched multiple values", externalId)
					}
				}
			}

			parsedValue, err = convertJSONAttributeValue(attribute, value, opts)
			if err != nil {
				return nil, err
			}
		} else if opts.complexAttributeNameDelimiter != "" {
			// The flattening of single-valued complex attributes is enabled,
			// so the attribute's parent complex attribute, if it exists, has
			// been parsed in the loop above. Look up the attribute directly in
			// the parsed object for that complex attribute.

			externalIdComponents := strings.SplitN(externalId, opts.complexAttributeNameDelimiter, 2)
			if len(externalIdComponents) == 2 {
				localExternalId := externalIdComponents[0]
				subExternalId := externalIdComponents[1]

				complexAttribute, found := complexAttributes[localExternalId]
				if !found {
					continue
				}

				parsedValue, found = complexAttribute[subExternalId]
				if !found {
					continue
				}
			}
		} else {
			value, found := object[externalId]
			if !found {
				continue
			}

			var err error
			parsedValue, err = convertJSONAttributeValue(attribute, value, opts)
			if err != nil {
				return nil, err
			}
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

		var parsedChildObjects []framework.Object

		if jsonPath, found := jsonPaths[externalId]; found { // The entity external ID is a JSONPath.
			childObjectsRaw, err := jsonPath(context.Background(), object)
			if err != nil { // The JSONPath didn't match. Evaluate to null.
				continue
			}

			childObjectsRawList, ok := childObjectsRaw.([]any)
			if !ok {
				return nil, fmt.Errorf("child entity %s is not associated with a list", externalId)
			}

			if len(childObjectsRawList) == 0 {
				continue
			}

			childObjects := make([]map[string]any, 0, len(childObjectsRawList))
			for _, childObjectRaw := range childObjectsRawList {
				childObject, ok := childObjectRaw.(map[string]any)
				if !ok {
					return nil, fmt.Errorf("child entity %s is not associated with a list of JSON objects", externalId)
				}
				childObjects = append(childObjects, childObject)
			}

			parsedChildObjects, err = convertJSONObjectList(childEntity, childObjects, opts)
			if err != nil {
				return nil, fmt.Errorf("failed to parse objects for child entity %s: %w", externalId, err)
			}
		} else if opts.complexAttributeNameDelimiter != "" {
			// The flattening of single-valued complex attributes is enabled,
			// so the attribute's parent complex attribute, if it exists, has
			// been parsed in the loop above. Look up the attribute directly in
			// the parsed object for that complex attribute.

			externalIdComponents := strings.SplitN(externalId, opts.complexAttributeNameDelimiter, 2)
			if len(externalIdComponents) == 2 {
				localExternalId := externalIdComponents[0]
				subExternalId := externalIdComponents[1]

				complexAttribute, found := complexAttributes[localExternalId]
				if !found {
					continue
				}

				parsedChildObjectsAny, found := complexAttribute[subExternalId]
				if !found {
					continue
				}

				var ok bool
				parsedChildObjects, ok = parsedChildObjectsAny.([]framework.Object)
				if !ok {
					panic(fmt.Sprintf("list of objects for child entity %s is not of type []framework.Object", externalId))
				}
			}
		} else {
			var childObjectsRaw any

			childObjectsRaw, found := object[externalId]
			if !found {
				continue
			}

			childObjectsRawList, ok := childObjectsRaw.([]any)
			if !ok {
				return nil, fmt.Errorf("child entity %s is not associated with a list", externalId)
			}

			if len(childObjectsRawList) == 0 {
				continue
			}

			childObjects := make([]map[string]any, 0, len(childObjectsRawList))
			for _, childObjectRaw := range childObjectsRawList {
				childObject, ok := childObjectRaw.(map[string]any)
				if !ok {
					return nil, fmt.Errorf("child entity %s is not associated with a list of JSON objects", externalId)
				}
				childObjects = append(childObjects, childObject)
			}

			var err error
			parsedChildObjects, err = convertJSONObjectList(childEntity, childObjects, opts)
			if err != nil {
				return nil, fmt.Errorf("failed to parse objects for child entity %s: %w", externalId, err)
			}
		}

		// Do not return an empty list.
		if len(parsedChildObjects) == 0 {
			continue
		}

		parsedObject[externalId] = parsedChildObjects
	}

	return parsedObject, nil
}
