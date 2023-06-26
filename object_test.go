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
	"reflect"
	"testing"
	"time"
)

func TestAddAttribute_ValidTypes(t *testing.T) {
	var boolValue bool = true
	var timeValue time.Time
	timeValue, _ = time.Parse(time.RFC3339, "2023-06-23T12:34:56Z")
	var durationValue time.Duration = 12 * time.Second
	var doubleValue float64 = 123.45
	var int64Value int64 = 123
	var stringValue string = "abcd"
	var boolListValue []bool = []bool{true, true, false, true}
	var timeListValue []time.Time = []time.Time{timeValue, timeValue.Add(10 * time.Second), timeValue.Add(20 * time.Second)}
	var durationListValue []time.Duration = []time.Duration{2 * time.Second, 4 * time.Second, 5 * time.Second}
	var doubleListValue []float64 = []float64{1.2, 3.4, 5.6}
	var int64ListValue []int64 = []int64{1, 2, 3, 5, 8}
	var stringListValue []string = []string{"a", "b", "c"}

	object := make(Object)

	AddAttribute(object, "bool", boolValue)
	AddAttribute(object, "time", timeValue)
	AddAttribute(object, "duration", durationValue)
	AddAttribute(object, "double", doubleValue)
	AddAttribute(object, "int64", int64Value)
	AddAttribute(object, "string", stringValue)
	AddAttribute(object, "boolList", boolListValue)
	AddAttribute(object, "timeList", timeListValue)
	AddAttribute(object, "durationList", durationListValue)
	AddAttribute(object, "doubleList", doubleListValue)
	AddAttribute(object, "int64List", int64ListValue)
	AddAttribute(object, "stringList", stringListValue)

	wantObject := Object{
		"bool":         boolValue,
		"time":         timeValue,
		"duration":     durationValue,
		"double":       doubleValue,
		"int64":        int64Value,
		"string":       stringValue,
		"boolList":     boolListValue,
		"timeList":     timeListValue,
		"durationList": durationListValue,
		"doubleList":   doubleListValue,
		"int64List":    int64ListValue,
		"stringList":   stringListValue,
	}

	if !reflect.DeepEqual(wantObject, object) {
		t.Errorf("Expected %#v, got %#v", wantObject, object)
	}
}

func TestAddAttribute_SameName(t *testing.T) {
	object := make(Object)

	var err error

	err = AddAttribute(object, "theName", "abcd")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	err = AddAttribute(object, "theName", "abcd")
	wantErr := errors.New("external ID already exists in object: theName")
	if err == nil || wantErr.Error() != err.Error() {
		t.Errorf("Expected %#v, got %#v", wantErr, err)
	}

	err = AddChildObjects(object, "theName", Object{})
	wantErr = errors.New("attribute already exists with that external ID in object: theName")
	if err == nil || wantErr.Error() != err.Error() {
		t.Errorf("Expected %#v, got %#v", wantErr, err)
	}
}

func TestAddChildObjects_Valid(t *testing.T) {
	object := make(Object)

	childEntity1Object1 := Object{"id": "entity1Object1"}
	childEntity1Object2 := Object{"id": "entity1Object2"}
	childEntity2Object1 := Object{"id": "entity2Object1"}
	childEntity2Object2 := Object{"id": "entity2Object2"}

	if err := AddChildObjects(object, "entity1"); err != nil { // No-op
		t.Errorf("Unexpected error: %v", err)
	}
	if err := AddChildObjects(object, "entity1", childEntity1Object1); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if err := AddChildObjects(object, "entity1", childEntity1Object1); err != nil { // Add same object twice
		t.Errorf("Unexpected error: %v", err)
	}
	if err := AddChildObjects(object, "entity1", childEntity1Object2); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if err := AddChildObjects(object, "entity2", childEntity2Object1, childEntity2Object2); err != nil { // Multiple objects
		t.Errorf("Unexpected error: %v", err)
	}

	wantObject := Object{
		"entity1": []Object{
			childEntity1Object1,
			childEntity1Object1,
			childEntity1Object2,
		},
		"entity2": []Object{
			childEntity2Object1,
			childEntity2Object2,
		},
	}

	if !reflect.DeepEqual(wantObject, object) {
		t.Errorf("Expected %#v, got %#v", wantObject, object)
	}
}

func TestAddChildObjects_SameName(t *testing.T) {
	object := make(Object)

	var err error

	err = AddChildObjects(object, "theName", Object{})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	err = AddAttribute(object, "theName", "abcd")
	wantErr := errors.New("external ID already exists in object: theName")
	if err == nil || wantErr.Error() != err.Error() {
		t.Errorf("Expected %#v, got %#v", wantErr, err)
	}

	err = AddChildObjects(object, "theName", Object{})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}
