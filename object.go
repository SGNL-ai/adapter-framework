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

package framework

import (
	"fmt"
	"time"
)

// Object is an object returned for the top entity or a child entity.
// Each key is the external ID of an attribute or of a child entity.
//
// Entries should be added into an Object by calling functions AddAttribute
// and AddChildObjects to ensure that the correct types are used for values.
type Object map[string]any

// AttributeValue is the set of types allowed for values of non-list attributes
// in an Object.
type AttributeValue interface {
	// Types of non-list attribute values.
	bool | time.Time | time.Duration | float64 | int64 | string |
		// Types of list attribute values.
		[]bool | []time.Time | []time.Duration | []float64 | []int64 | []string
}

// AddAttribute adds a attribute into the given object.
// Returns an error if an attribute or child object has already been added with
// the same external ID.
func AddAttribute[Value AttributeValue](object Object, attributeExternalId string, value Value) error {
	_, found := object[attributeExternalId]
	if found {
		return fmt.Errorf("external ID already exists in object: %s", attributeExternalId)
	}

	object[attributeExternalId] = value

	return nil
}

// AddChildObjects adds child objects into the given object.
// If child objects have already been added for the same entity external ID,
// the given objects are appended to the list of child objects in that entity.
// Returns an error if an attribute has already been added with the same
// external ID.
func AddChildObjects(object Object, entityExternalId string, childObjects ...Object) error {
	if len(childObjects) == 0 {
		return nil
	}

	value, found := object[entityExternalId]
	if found {
		currentChildObjects, isListOfObjects := value.([]Object)
		if !isListOfObjects {
			return fmt.Errorf("attribute already exists with that external ID in object: %s", entityExternalId)
		}

		object[entityExternalId] = append(currentChildObjects, childObjects...)
	} else {
		object[entityExternalId] = childObjects
	}

	return nil
}
