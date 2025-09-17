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
	"encoding/json"
	"errors"
	"testing"

	"github.com/PaesslerAG/gval"
	framework "github.com/sgnl-ai/adapter-framework"
)

func TestConvertJSONObject_NoFlattening(t *testing.T) {
	testJSONOptions := defaultJSONOptions()

	tests := map[string]struct {
		entity     *framework.EntityConfig
		objectJSON string
		opts       *jsonOptions
		wantObject framework.Object
		wantError  error
	}{
		"empty_object_no_attributes": {
			entity: &framework.EntityConfig{
				ExternalId: "test",
			},
			objectJSON: `{}`,
			opts:       testJSONOptions,
			wantObject: framework.Object{},
			wantError:  nil,
		},
		"one_attribute_matching": {
			entity: &framework.EntityConfig{
				ExternalId: "test",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "a",
						Type:       framework.AttributeTypeString,
					},
				},
			},
			objectJSON: `{
				"a": "a value"
			}`,
			opts: testJSONOptions,
			wantObject: framework.Object{
				"a": "a value",
			},
			wantError: nil,
		},
		"one_attribute_not_matching": {
			entity: &framework.EntityConfig{
				ExternalId: "test",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "b",
						Type:       framework.AttributeTypeString,
					},
				},
			},
			objectJSON: `{
				"a": "a value"
			}`,
			opts:       testJSONOptions,
			wantObject: framework.Object{},
			wantError:  nil,
		},
		"multiple_attributes": {
			entity: &framework.EntityConfig{
				ExternalId: "test",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "a",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "c",
						Type:       framework.AttributeTypeString,
					},
				},
			},
			objectJSON: `{
				"a": "a value",
				"b": "b value",
				"c": "c value",
				"d": "d value"
			}`,
			opts: testJSONOptions,
			wantObject: framework.Object{
				"a": "a value",
				"c": "c value",
			},
			wantError: nil,
		},
		"one_child_entity_matching": {
			entity: &framework.EntityConfig{
				ExternalId: "test",
				ChildEntities: []*framework.EntityConfig{
					{
						ExternalId: "childEntity1",
						Attributes: []*framework.AttributeConfig{
							{
								ExternalId: "a",
								Type:       framework.AttributeTypeString,
							},
							{
								ExternalId: "c",
								Type:       framework.AttributeTypeString,
							},
						},
					},
				},
			},
			objectJSON: `{
				"childEntity1": [
					{
						"a": "childEntity1-object1-a",
						"b": "childEntity1-object1-b"
					},
					{
						"c": "childEntity1-object2-c",
						"d": "childEntity1-object2-d"
					}
				]
			}`,
			opts: testJSONOptions,
			wantObject: framework.Object{
				"childEntity1": []framework.Object{
					{
						"a": "childEntity1-object1-a",
					},
					{
						"c": "childEntity1-object2-c",
					},
				},
			},
			wantError: nil,
		},
		"recursive_child_entities": {
			entity: &framework.EntityConfig{
				ExternalId: "test",
				ChildEntities: []*framework.EntityConfig{
					{
						ExternalId: "childEntity1",
						ChildEntities: []*framework.EntityConfig{
							{
								ExternalId: "childEntity1.1",
								Attributes: []*framework.AttributeConfig{
									{
										ExternalId: "a",
										Type:       framework.AttributeTypeString,
									},
								},
							},
							{
								ExternalId: "childEntity1.2",
								Attributes: []*framework.AttributeConfig{
									{
										ExternalId: "c",
										Type:       framework.AttributeTypeString,
									},
								}},
						},
					},
				},
			},
			objectJSON: `{
				"childEntity1": [
					{
						"childEntity1.1": [
							{
								"a": "childEntity1.1-object1-a"
							},
							{
								"b": "childEntity1.1-object2-b"
							}
						]
					},
					{
						"childEntity1.2": [
							{
								"a": "childEntity1.2-object1-a"
							},
							{
								"c": "childEntity1.2-object2-c"
							}
						]
					}
				]
			}`,
			opts: testJSONOptions,
			wantObject: framework.Object{
				"childEntity1": []framework.Object{
					{
						"childEntity1.1": []framework.Object{
							{
								"a": "childEntity1.1-object1-a",
							},
						},
					},
					{
						"childEntity1.2": []framework.Object{
							{
								"c": "childEntity1.2-object2-c",
							},
						},
					},
				},
			},
			wantError: nil,
		},
		"attributes_and_recursive_child_entities": {
			entity: &framework.EntityConfig{
				ExternalId: "users",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "id",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "name",
						Type:       framework.AttributeTypeString,
					},
				},
				ChildEntities: []*framework.EntityConfig{
					{
						ExternalId: "emails",
						Attributes: []*framework.AttributeConfig{
							{
								ExternalId: "email",
								Type:       framework.AttributeTypeString,
							},
							{
								ExternalId: "primary",
								Type:       framework.AttributeTypeBool,
							},
						},
					},
					{
						ExternalId: "addresses",
						Attributes: []*framework.AttributeConfig{
							{
								ExternalId: "streetLines",
								Type:       framework.AttributeTypeString,
								List:       true,
							},
							{
								ExternalId: "postalCode",
								Type:       framework.AttributeTypeString,
							},
							{
								ExternalId: "region",
								Type:       framework.AttributeTypeString,
							},
							{
								ExternalId: "country",
								Type:       framework.AttributeTypeString,
							},
						},
					},
				},
			},
			objectJSON: `{
				"id": "1234",
				"name": "John Doe",
				"emails": [
					{
						"email": "john@doe.com",
						"primary": true
					},
					{
						"email": "john@doe.org"
					}
				],
				"addresses": [
					{
						"streetLines": [
							"1234 Somewhere St"
						],
						"region": "CA",
						"country": "USA"
					}
				]
			}`,
			opts: testJSONOptions,
			wantObject: framework.Object{
				"id":   "1234",
				"name": "John Doe",
				"emails": []framework.Object{
					{
						"email":   "john@doe.com",
						"primary": true,
					},
					{
						"email": "john@doe.org",
					},
				},
				"addresses": []framework.Object{
					{
						"streetLines": []string{
							"1234 Somewhere St",
						},
						"region":  "CA",
						"country": "USA",
					},
				},
			},
			wantError: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var object map[string]any
			err := json.Unmarshal([]byte(tc.objectJSON), &object)
			if err != nil {
				t.Fatalf("Failed to unmarshal test input JSON object: %v", err)
			}

			gotObject, gotError := convertJSONObject(tc.entity, object, tc.opts, nil)

			AssertDeepEqual(t, tc.wantError, gotError)
			AssertDeepEqual(t, tc.wantObject, gotObject)
		})
	}
}

func TestConvertJSONObject_Flattening(t *testing.T) {
	testJSONOptions := defaultJSONOptions()
	testJSONOptions.complexAttributeNameDelimiter = "__"

	tests := map[string]struct {
		entity     *framework.EntityConfig
		objectJSON string
		opts       *jsonOptions
		wantObject framework.Object
		wantError  error
	}{
		"one_complex_attribute_matching": {
			entity: &framework.EntityConfig{
				ExternalId: "test",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "a__b",
						Type:       framework.AttributeTypeString,
					},
				},
			},
			objectJSON: `{
				"a": {
					"a": "a__a value",
					"b": "a__b value",
					"c": "a__c value"
				}
			}`,
			opts: testJSONOptions,
			wantObject: framework.Object{
				"a__b": "a__b value",
			},
			wantError: nil,
		},
		"one_complex_attribute_not_matching": {
			entity: &framework.EntityConfig{
				ExternalId: "test",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "a__b",
						Type:       framework.AttributeTypeString,
					},
				},
			},
			objectJSON: `{
				"c": {
					"b": "c__b value"
				}
			}`,
			opts:       testJSONOptions,
			wantObject: framework.Object{},
			wantError:  nil,
		},
		"one_complex_attribute_sub_not_matching": {
			entity: &framework.EntityConfig{
				ExternalId: "test",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "a__b",
						Type:       framework.AttributeTypeString,
					},
				},
			},
			objectJSON: `{
				"a": {
					"c": "a__c value"
				}
			}`,
			opts:       testJSONOptions,
			wantObject: framework.Object{},
			wantError:  nil,
		},
		"multiple_complex_attributes": {
			entity: &framework.EntityConfig{
				ExternalId: "test",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "e",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "a__b",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "c__d",
						Type:       framework.AttributeTypeString,
					},
				},
			},
			objectJSON: `{
				"a": {
					"a": "a__a value",
					"b": "a__b value"
				},
				"c": {
					"c": "c__c value",
					"d": "c__d value"
				},
				"e": "e value",
				"f": "f value"
			}`,
			opts: testJSONOptions,
			wantObject: framework.Object{
				"a__b": "a__b value",
				"c__d": "c__d value",
				"e":    "e value",
			},
			wantError: nil,
		},
		"recursive_complex_attributes": {
			entity: &framework.EntityConfig{
				ExternalId: "test",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "a__a__a",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "a__a__b",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "a__c__b",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "a__b",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "c__a__a",
						Type:       framework.AttributeTypeString,
					},
				},
			},
			objectJSON: `{
				"a": {
					"a": {
						"a": "a__a__a value",
						"b": "a__a__b value",
						"c": "a__a__c value"
					},
					"b": "a__b value",
					"c": {
						"a": "a__c__a value",
						"b": "a__c__b value"
					},
					"d": "a__d value"
				},
				"c": {
					"a": {
						"a": "c__a__a value",
						"b": "c__a__b value"
					},
					"b": "c__b value"
				}
			}`,
			opts: testJSONOptions,
			wantObject: framework.Object{
				"a__a__a": "a__a__a value",
				"a__a__b": "a__a__b value",
				"a__c__b": "a__c__b value",
				"a__b":    "a__b value",
				"c__a__a": "c__a__a value",
			},
			wantError: nil,
		},
		"null_complex_attribute": {
			entity: &framework.EntityConfig{
				ExternalId: "test",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "a__b__c",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "d__e",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "d__e__f",
						Type:       framework.AttributeTypeString,
					},
				},
			},
			// Suppose a__b is regularly a complex attribute, but in this case it is null. Same with d.
			// We should ignore these attributes.
			objectJSON: `{
				"a": {
					"b": null
				},
				"d": null
			}`,
			opts:       testJSONOptions,
			wantObject: framework.Object{},
			wantError:  nil,
		},
		"child_entities_in_complex_attributes": {
			entity: &framework.EntityConfig{
				ExternalId: "test",
				ChildEntities: []*framework.EntityConfig{
					{
						ExternalId: "a__b",
						Attributes: []*framework.AttributeConfig{
							{
								ExternalId: "a",
								Type:       framework.AttributeTypeString,
							},
						},
					},
					{
						ExternalId: "a__c__d",
						Attributes: []*framework.AttributeConfig{
							{
								ExternalId: "c",
								Type:       framework.AttributeTypeString,
							},
						},
					},
					{
						ExternalId: "a__d",
						Attributes: []*framework.AttributeConfig{
							{
								ExternalId: "a",
								Type:       framework.AttributeTypeString,
							},
						},
					},
				},
			},
			objectJSON: `{
				"a": {
					"a": "a__a value",
					"b": [
						{
							"a": "a__b child1 a value",
							"b": "a__b child1 b value"
						},
						{
							"a": "a__b child2 a value",
							"b": "a__b child2 b value"
						}
					],
					"c": {
						"a": "a__c__a value",
						"d": [
							{
								"c": "a__c__d child1 c value"
							},
							{
								"d": "a__c__d child2 d value"
							}
						]
					}
				}
			}`,
			opts: testJSONOptions,
			wantObject: framework.Object{
				"a__b": []framework.Object{
					{
						"a": "a__b child1 a value",
					},
					{
						"a": "a__b child2 a value",
					},
				},
				"a__c__d": []framework.Object{
					{
						"c": "a__c__d child1 c value",
					},
				},
				// "a__d" didn't match any attribute
			},
			wantError: nil,
		},
		"complex_attributes_in_child_entities": {
			entity: &framework.EntityConfig{
				ExternalId: "test",
				ChildEntities: []*framework.EntityConfig{
					{
						ExternalId: "a",
						Attributes: []*framework.AttributeConfig{
							{
								ExternalId: "a",
								Type:       framework.AttributeTypeString,
							},
							{
								ExternalId: "b__a",
								Type:       framework.AttributeTypeString,
							},
						},
					},
					{
						ExternalId: "b__a",
						Attributes: []*framework.AttributeConfig{
							{
								ExternalId: "a__a",
								Type:       framework.AttributeTypeString,
							},
						},
						ChildEntities: []*framework.EntityConfig{
							{
								ExternalId: "a__b",
								Attributes: []*framework.AttributeConfig{
									{
										ExternalId: "a__b",
										Type:       framework.AttributeTypeString,
									},
								},
							},
						},
					},
				},
			},
			objectJSON: `{
				"a": [
					{
						"a": "a child1 a value",
						"b": {
							"a": "a child1 b__a value",
							"b": "a child1 b__b value"
						},
						"c": "a child1 c value"
					},
					{
						"a": "a child2 a value",
						"b": {
							"z": "a child2 b__z value"
						},
						"c": "a child2 c value"
					},
					{
						"b": {
							"z": "a child3 b__z value"
						},
						"c": "a child3 c value"
					}
				],
				"b": {
					"a": [
						{
							"a": {
								"a": "b__a child1 a__a value",
								"b": [
									{
										"a": {
											"a": "b__a child1 a__b child1 a value",
											"b": "b__a child1 a__b child1 b value"
										}
									},
									{
										"a": {
											"a": "b__a child1 a__b child2 a value"
										}
									}
								]
							}
						},
						{
							"a": {
								"z": "b__a child1 a__z value"
							}
						}
					]
				},
				"c": "c value"
			}`,
			opts: testJSONOptions,
			wantObject: framework.Object{
				"a": []framework.Object{
					{
						"a":    "a child1 a value",
						"b__a": "a child1 b__a value",
					},
					{
						"a": "a child2 a value",
					},
				},
				"b__a": []framework.Object{
					{
						"a__a": "b__a child1 a__a value",
						"a__b": []framework.Object{
							{
								"a__b": "b__a child1 a__b child1 b value",
							},
						},
					},
				},
			},
			wantError: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var object map[string]any
			err := json.Unmarshal([]byte(tc.objectJSON), &object)
			if err != nil {
				t.Fatalf("Failed to unmarshal test input JSON object: %v", err)
			}

			gotObject, gotError := convertJSONObject(tc.entity, object, tc.opts, nil)

			AssertDeepEqual(t, tc.wantError, gotError)
			AssertDeepEqual(t, tc.wantObject, gotObject)
		})
	}
}

func TestConvertJSONObject_JSONPath(t *testing.T) {
	testJSONOptions := defaultJSONOptions()
	testJSONOptions.enableJSONPath = true

	tests := map[string]struct {
		entity     *framework.EntityConfig
		objectJSON string
		opts       *jsonOptions
		wantObject framework.Object
		wantError  error
	}{
		"empty_object_no_attributes": {
			entity: &framework.EntityConfig{
				ExternalId: "test",
			},
			objectJSON: `{}`,
			opts:       testJSONOptions,
			wantObject: framework.Object{},
			wantError:  nil,
		},
		"one_complex_attribute_matching": {
			entity: &framework.EntityConfig{
				ExternalId: "test",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "$.a.b",
						Type:       framework.AttributeTypeString,
					},
				},
			},
			objectJSON: `{
				"a": {
					"a": "a__a value",
					"b": "a__b value",
					"c": "a__c value"
				}
			}`,
			opts: testJSONOptions,
			wantObject: framework.Object{
				"$.a.b": "a__b value",
			},
			wantError: nil,
		},
		"one_complex_attribute_not_matching": {
			entity: &framework.EntityConfig{
				ExternalId: "test",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "$.a.b",
						Type:       framework.AttributeTypeString,
					},
				},
			},
			objectJSON: `{
				"c": {
					"b": "c__b value"
				}
			}`,
			opts:       testJSONOptions,
			wantObject: framework.Object{},
			wantError:  nil,
		},
		"one_complex_attribute_sub_not_matching": {
			entity: &framework.EntityConfig{
				ExternalId: "test",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "$.a.b",
						Type:       framework.AttributeTypeString,
					},
				},
			},
			objectJSON: `{
				"a": {
					"c": "a__c value"
				}
			}`,
			opts:       testJSONOptions,
			wantObject: framework.Object{},
			wantError:  nil,
		},
		"multiple_complex_attributes": {
			entity: &framework.EntityConfig{
				ExternalId: "test",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "e",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "$.a.b",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "$.c.d",
						Type:       framework.AttributeTypeString,
					},
				},
			},
			objectJSON: `{
				"a": {
					"a": "a__a value",
					"b": "a__b value"
				},
				"c": {
					"c": "c__c value",
					"d": "c__d value"
				},
				"e": "e value",
				"f": "f value"
			}`,
			opts: testJSONOptions,
			wantObject: framework.Object{
				"$.a.b": "a__b value",
				"$.c.d": "c__d value",
				"e":     "e value",
			},
			wantError: nil,
		},
		"recursive_complex_attributes": {
			entity: &framework.EntityConfig{
				ExternalId: "test",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "$.a.a.a",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "$.a.a.b",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "$.a.c.b",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "$.a.b",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "$.c.a.a",
						Type:       framework.AttributeTypeString,
					},
				},
			},
			objectJSON: `{
				"a": {
					"a": {
						"a": "a__a__a value",
						"b": "a__a__b value",
						"c": "a__a__c value"
					},
					"b": "a__b value",
					"c": {
						"a": "a__c__a value",
						"b": "a__c__b value"
					},
					"d": "a__d value"
				},
				"c": {
					"a": {
						"a": "c__a__a value",
						"b": "c__a__b value"
					},
					"b": "c__b value"
				}
			}`,
			opts: testJSONOptions,
			wantObject: framework.Object{
				"$.a.a.a": "a__a__a value",
				"$.a.a.b": "a__a__b value",
				"$.a.c.b": "a__c__b value",
				"$.a.b":   "a__b value",
				"$.c.a.a": "c__a__a value",
			},
			wantError: nil,
		},
		"null_complex_attribute": {
			entity: &framework.EntityConfig{
				ExternalId: "test",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "$.a.b.c",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "$.d.e",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "$.d.e.f",
						Type:       framework.AttributeTypeString,
					},
				},
			},
			// Suppose a__b is regularly a complex attribute, but in this case it is null. Same with d.
			// We should ignore these attributes.
			objectJSON: `{
				"a": {
					"b": null
				},
				"d": null
			}`,
			opts:       testJSONOptions,
			wantObject: framework.Object{},
			wantError:  nil,
		},
		"child_entities_in_complex_attributes": {
			entity: &framework.EntityConfig{
				ExternalId: "test",
				ChildEntities: []*framework.EntityConfig{
					{
						ExternalId: "$.a.b",
						Attributes: []*framework.AttributeConfig{
							{
								ExternalId: "a",
								Type:       framework.AttributeTypeString,
							},
						},
					},
					{
						ExternalId: "$.a.c.d",
						Attributes: []*framework.AttributeConfig{
							{
								ExternalId: "c",
								Type:       framework.AttributeTypeString,
							},
						},
					},
					{
						ExternalId: "$.a.d",
						Attributes: []*framework.AttributeConfig{
							{
								ExternalId: "a",
								Type:       framework.AttributeTypeString,
							},
						},
					},
				},
			},
			objectJSON: `{
				"a": {
					"a": "a__a value",
					"b": [
						{
							"a": "a__b child1 a value",
							"b": "a__b child1 b value"
						},
						{
							"a": "a__b child2 a value",
							"b": "a__b child2 b value"
						}
					],
					"c": {
						"a": "a__c__a value",
						"d": [
							{
								"c": "a__c__d child1 c value"
							},
							{
								"d": "a__c__d child2 d value"
							}
						]
					}
				}
			}`,
			opts: testJSONOptions,
			wantObject: framework.Object{
				"$.a.b": []framework.Object{
					{
						"a": "a__b child1 a value",
					},
					{
						"a": "a__b child2 a value",
					},
				},
				"$.a.c.d": []framework.Object{
					{
						"c": "a__c__d child1 c value",
					},
				},
				// "$.a.d" didn't match any attribute
			},
			wantError: nil,
		},
		"complex_attributes_in_child_entities": {
			entity: &framework.EntityConfig{
				ExternalId: "test",
				ChildEntities: []*framework.EntityConfig{
					{
						ExternalId: "a",
						Attributes: []*framework.AttributeConfig{
							{
								ExternalId: "a",
								Type:       framework.AttributeTypeString,
							},
							{
								ExternalId: "$.b.a",
								Type:       framework.AttributeTypeString,
							},
						},
					},
					{
						ExternalId: "$.b.a",
						Attributes: []*framework.AttributeConfig{
							{
								ExternalId: "$.a.a",
								Type:       framework.AttributeTypeString,
							},
						},
						ChildEntities: []*framework.EntityConfig{
							{
								ExternalId: "$.a.b",
								Attributes: []*framework.AttributeConfig{
									{
										ExternalId: "$.a.b",
										Type:       framework.AttributeTypeString,
									},
								},
							},
						},
					},
				},
			},
			objectJSON: `{
				"a": [
					{
						"a": "a child1 a value",
						"b": {
							"a": "a child1 b__a value",
							"b": "a child1 b__b value"
						},
						"c": "a child1 c value"
					},
					{
						"a": "a child2 a value",
						"b": {
							"z": "a child2 b__z value"
						},
						"c": "a child2 c value"
					},
					{
						"b": {
							"z": "a child3 b__z value"
						},
						"c": "a child3 c value"
					}
				],
				"b": {
					"a": [
						{
							"a": {
								"a": "b__a child1 a__a value",
								"b": [
									{
										"a": {
											"a": "b__a child1 a__b child1 a__a value",
											"b": "b__a child1 a__b child1 a__b value"
										}
									},
									{
										"a": {
											"a": "b__a child1 a__b child2 a__a value"
										}
									}
								]
							}
						},
						{
							"a": {
								"z": "b__a child1 a__z value"
							}
						}
					]
				},
				"c": "c value"
			}`,
			opts: testJSONOptions,
			wantObject: framework.Object{
				"a": []framework.Object{
					{
						"a":     "a child1 a value",
						"$.b.a": "a child1 b__a value",
					},
					{
						"a": "a child2 a value",
					},
				},
				"$.b.a": []framework.Object{
					{
						"$.a.a": "b__a child1 a__a value",
						"$.a.b": []framework.Object{
							{
								"$.a.b": "b__a child1 a__b child1 a__b value",
							},
						},
					},
				},
			},
			wantError: nil,
		},
		"multivalued_complex_attribute_into_list": {
			entity: &framework.EntityConfig{
				ExternalId: "test",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "$.emails[*].value",
						Type:       framework.AttributeTypeString,
						List:       true,
					},
					{
						ExternalId: "$.emails[?(@.primary==true)].value",
						Type:       framework.AttributeTypeString,
						List:       true,
					},
				},
			},
			objectJSON: `{
				"emails": [
					{
						"value": "primary@example.com",
						"primary": true
					},
					{
						"value": "secondary@example.com",
						"primary": false
					}
				]
			}`,
			opts: testJSONOptions,
			wantObject: framework.Object{
				"$.emails[*].value":                  []string{"primary@example.com", "secondary@example.com"},
				"$.emails[?(@.primary==true)].value": []string{"primary@example.com"},
			},
			wantError: nil,
		},
		"multivalued_complex_attribute_into_single_value": {
			entity: &framework.EntityConfig{
				ExternalId: "test",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "$.emails[?(@.primary==true)].value",
						Type:       framework.AttributeTypeString,
					},
				},
			},
			objectJSON: `{
				"emails": [
					{
						"value": "primary@example.com",
						"primary": true
					},
					{
						"value": "secondary@example.com",
						"primary": false
					}
				]
			}`,
			opts: testJSONOptions,
			wantObject: framework.Object{
				"$.emails[?(@.primary==true)].value": "primary@example.com",
			},
			wantError: nil,
		},
		"multivalued_complex_attribute_into_single_value_error": {
			entity: &framework.EntityConfig{
				ExternalId: "test",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "$.emails[*].value", // This JSONPath selects multiple values.
						Type:       framework.AttributeTypeString,
					},
				},
			},
			objectJSON: `{
				"emails": [
					{
						"value": "primary@example.com",
						"primary": true
					},
					{
						"value": "secondary@example.com",
						"primary": false
					}
				]
			}`,
			opts:       testJSONOptions,
			wantObject: nil,
			wantError:  errors.New("non-list attribute $.emails[*].value matched multiple values"),
		},
		"child_entities_jsonpath_returns_nil": {
			entity: &framework.EntityConfig{
				ExternalId: "test",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "id",
						Type:       framework.AttributeTypeString,
					},
				},
				ChildEntities: []*framework.EntityConfig{
					{
						ExternalId: "$.nullChildren", // This JSONPath will return nil (no error)
						Attributes: []*framework.AttributeConfig{
							{
								ExternalId: "name",
								Type:       framework.AttributeTypeString,
							},
						},
					},
					{
						ExternalId: "$.existing.children",
						Attributes: []*framework.AttributeConfig{
							{
								ExternalId: "name",
								Type:       framework.AttributeTypeString,
							},
						},
					},
				},
			},
			objectJSON: `{
				"id": "test123",
				"nullChildren": null,
				"existing": {
					"children": [
						{
							"name": "child1"
						}
					]
				}
			}`,
			opts: testJSONOptions,
			wantObject: framework.Object{
				"id": "test123",
				"$.existing.children": []framework.Object{
					{
						"name": "child1",
					},
				},
			},
			wantError: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var object map[string]any
			if err := json.Unmarshal([]byte(tc.objectJSON), &object); err != nil {
				t.Fatalf("Failed to unmarshal test input JSON object: %v", err)
			}

			jsonPaths := make(map[string]gval.Evaluable)
			if err := parseJSONPaths(tc.entity, jsonPaths); err != nil {
				t.Fatalf("parseJSONPaths failed unexpectedly: %v", err)
			}

			gotObject, gotError := convertJSONObject(tc.entity, object, tc.opts, jsonPaths)

			AssertDeepEqual(t, tc.wantError, gotError)
			AssertDeepEqual(t, tc.wantObject, gotObject)
		})
	}
}
