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
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/sosodev/duration"
)

// Object is an object returned for the top entity or a child entity.
// Each key is the external ID of an attribute or of a child entity.
//
// Entries should be added into an Object by calling functions AddAttribute
// and AddChildObjects to ensure that the correct types are used for values.
type Object map[string]any

// Duration is the value of an attribute of type duration.
// It is the sum of all the fields' durations.
type Duration struct {
	// Nanos is a duration as a number of nanoseconds.
	Nanos int32 `json:"nanos,omitempty"`

	// Seconds is a duration as a number of seconds.
	Seconds int64 `json:"seconds,omitempty"`

	// Days is a duration as a number of days.
	Days int64 `json:"days,omitempty"`

	// Months is a duration as a number of months.
	Months int64 `json:"months,omitempty"`
}

// ParseDuration parses a valid ISO8601 duration string into a Duration.
// Years as duration is not supported.
// Fractional components of durations are not supported.
func ParseISO8601Duration(durationStr string) (*Duration, error) {
	d, err := duration.Parse(durationStr)
	if err != nil {
		return nil, errors.New("failed to parse the duration string: " + durationStr)
	}

	// Convert years into months, weeks into days, and minutes and hours into seconds.
	d.Months += d.Years * 12.0
	d.Days += d.Weeks * 7.0
	d.Seconds += d.Hours * 3_600.0
	d.Seconds += d.Minutes * 60.0

	// Round the numbers of months, days, seconds, and nanoseconds.
	months, monthsFraction := math.Modf(d.Months)
	d.Days += monthsFraction * 30.0

	days, daysFraction := math.Modf(d.Days)
	d.Seconds += daysFraction * 24.0 * 3_600.0

	seconds, secondsFraction := math.Modf(d.Seconds)
	nanos := secondsFraction * 1_000_000_000.0

	return &Duration{
		Months:  int64(months),
		Days:    int64(days),
		Seconds: int64(seconds),
		Nanos:   int32(nanos),
	}, nil
}

// AttributeValue is the set of types allowed for values of non-list attributes
// in an Object.
type AttributeValue interface {
	// Types of non-list attribute values.
	bool | time.Time | Duration | float64 | int64 | string |
		// Types of list attribute values.
		[]bool | []time.Time | []Duration | []float64 | []int64 | []string
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
