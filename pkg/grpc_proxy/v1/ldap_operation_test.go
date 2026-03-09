// Tests for LDAP operation proto message types (add, modify, delete, modifyDN).
// Verifies construction, oneof routing, marshal/unmarshal roundtrip, and enum value alignment.

package v1

import (
	"testing"

	"google.golang.org/protobuf/proto"
)

func TestLDAPOperationRequest_GivenAddRequest_WhenConstructed_ThenFieldsAreCorrect(t *testing.T) {
	tests := []struct {
		name       string
		dn         string
		attributes []*LDAPAttribute
	}{
		{
			name: "single_attribute",
			dn:   "cn=John Doe,ou=Users,dc=example,dc=com",
			attributes: []*LDAPAttribute{
				{Type: "objectClass", Values: []string{"user"}},
			},
		},
		{
			name: "multiple_attributes",
			dn:   "cn=Jane Doe,ou=Users,dc=example,dc=com",
			attributes: []*LDAPAttribute{
				{Type: "objectClass", Values: []string{"user", "person"}},
				{Type: "sAMAccountName", Values: []string{"jdoe"}},
				{Type: "cn", Values: []string{"Jane Doe"}},
			},
		},
		{
			name:       "empty_attributes",
			dn:         "cn=Empty,ou=Users,dc=example,dc=com",
			attributes: []*LDAPAttribute{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			addReq := &LDAPAddRequest{
				Dn:         tt.dn,
				Attributes: tt.attributes,
			}

			// Act
			opReq := &LDAPOperationRequest{
				Url:          "ldaps://ad.example.com:636",
				BindDn:       "cn=admin,dc=example,dc=com",
				BindPassword: "secret",
				Operation:    &LDAPOperationRequest_AddRequest{AddRequest: addReq},
			}

			// Assert
			if opReq.GetAddRequest() == nil {
				t.Fatal("expected add_request to be set")
			}
			if opReq.GetAddRequest().GetDn() != tt.dn {
				t.Errorf("got dn %q, want %q", opReq.GetAddRequest().GetDn(), tt.dn)
			}
			if len(opReq.GetAddRequest().GetAttributes()) != len(tt.attributes) {
				t.Errorf("got %d attributes, want %d", len(opReq.GetAddRequest().GetAttributes()), len(tt.attributes))
			}
			if opReq.GetModifyRequest() != nil {
				t.Error("expected modify_request to be nil")
			}
			if opReq.GetDeleteRequest() != nil {
				t.Error("expected delete_request to be nil")
			}
			if opReq.GetModifyDnRequest() != nil {
				t.Error("expected modify_dn_request to be nil")
			}
		})
	}
}

func TestLDAPOperationRequest_GivenModifyRequest_WhenConstructed_ThenFieldsAreCorrect(t *testing.T) {
	tests := []struct {
		name    string
		dn      string
		changes []*LDAPModifyChange
	}{
		{
			name: "replace_single_attribute",
			dn:   "cn=John Doe,ou=Users,dc=example,dc=com",
			changes: []*LDAPModifyChange{
				{
					Operation:    LDAPModifyOperation_LDAP_MODIFY_OPERATION_REPLACE,
					Modification: &LDAPAttribute{Type: "userAccountControl", Values: []string{"512"}},
				},
			},
		},
		{
			name: "multiple_changes",
			dn:   "cn=John Doe,ou=Users,dc=example,dc=com",
			changes: []*LDAPModifyChange{
				{
					Operation:    LDAPModifyOperation_LDAP_MODIFY_OPERATION_REPLACE,
					Modification: &LDAPAttribute{Type: "title", Values: []string{"Engineer"}},
				},
				{
					Operation:    LDAPModifyOperation_LDAP_MODIFY_OPERATION_DELETE,
					Modification: &LDAPAttribute{Type: "description", Values: []string{}},
				},
				{
					Operation:    LDAPModifyOperation_LDAP_MODIFY_OPERATION_ADD,
					Modification: &LDAPAttribute{Type: "mail", Values: []string{"john@example.com"}},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			modReq := &LDAPModifyRequest{
				Dn:      tt.dn,
				Changes: tt.changes,
			}

			// Act
			opReq := &LDAPOperationRequest{
				Url:          "ldaps://ad.example.com:636",
				BindDn:       "cn=admin,dc=example,dc=com",
				BindPassword: "secret",
				Operation:    &LDAPOperationRequest_ModifyRequest{ModifyRequest: modReq},
			}

			// Assert
			if opReq.GetModifyRequest() == nil {
				t.Fatal("expected modify_request to be set")
			}
			if opReq.GetModifyRequest().GetDn() != tt.dn {
				t.Errorf("got dn %q, want %q", opReq.GetModifyRequest().GetDn(), tt.dn)
			}
			if len(opReq.GetModifyRequest().GetChanges()) != len(tt.changes) {
				t.Errorf("got %d changes, want %d", len(opReq.GetModifyRequest().GetChanges()), len(tt.changes))
			}
			if opReq.GetAddRequest() != nil {
				t.Error("expected add_request to be nil")
			}
		})
	}
}

func TestLDAPOperationRequest_GivenDeleteRequest_WhenConstructed_ThenFieldsAreCorrect(t *testing.T) {
	// Arrange
	delReq := &LDAPDeleteRequest{
		Dn: "cn=John Doe,ou=Users,dc=example,dc=com",
	}

	// Act
	opReq := &LDAPOperationRequest{
		Url:          "ldaps://ad.example.com:636",
		BindDn:       "cn=admin,dc=example,dc=com",
		BindPassword: "secret",
		Operation:    &LDAPOperationRequest_DeleteRequest{DeleteRequest: delReq},
	}

	// Assert
	if opReq.GetDeleteRequest() == nil {
		t.Fatal("expected delete_request to be set")
	}
	if opReq.GetDeleteRequest().GetDn() != "cn=John Doe,ou=Users,dc=example,dc=com" {
		t.Errorf("got dn %q, want %q", opReq.GetDeleteRequest().GetDn(), "cn=John Doe,ou=Users,dc=example,dc=com")
	}
	if opReq.GetAddRequest() != nil {
		t.Error("expected add_request to be nil")
	}
	if opReq.GetModifyRequest() != nil {
		t.Error("expected modify_request to be nil")
	}
	if opReq.GetModifyDnRequest() != nil {
		t.Error("expected modify_dn_request to be nil")
	}
}

func TestLDAPOperationRequest_GivenModifyDNRequest_WhenConstructed_ThenFieldsAreCorrect(t *testing.T) {
	tests := []struct {
		name         string
		dn           string
		newRDN       string
		deleteOldRDN bool
		newSuperior  string
	}{
		{
			name:         "rename_only",
			dn:           "cn=John Doe,ou=Users,dc=example,dc=com",
			newRDN:       "cn=Jane Doe",
			deleteOldRDN: true,
			newSuperior:  "",
		},
		{
			name:         "move_to_new_parent",
			dn:           "cn=John Doe,ou=Users,dc=example,dc=com",
			newRDN:       "cn=John Doe",
			deleteOldRDN: true,
			newSuperior:  "ou=ArchivedUsers,dc=example,dc=com",
		},
		{
			name:         "rename_and_move",
			dn:           "cn=John Doe,ou=Users,dc=example,dc=com",
			newRDN:       "cn=Jane Doe",
			deleteOldRDN: false,
			newSuperior:  "ou=Managers,dc=example,dc=com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			modDNReq := &LDAPModifyDNRequest{
				Dn:           tt.dn,
				NewRdn:       tt.newRDN,
				DeleteOldRdn: tt.deleteOldRDN,
				NewSuperior:  tt.newSuperior,
			}

			// Act
			opReq := &LDAPOperationRequest{
				Url:          "ldaps://ad.example.com:636",
				BindDn:       "cn=admin,dc=example,dc=com",
				BindPassword: "secret",
				Operation:    &LDAPOperationRequest_ModifyDnRequest{ModifyDnRequest: modDNReq},
			}

			// Assert
			if opReq.GetModifyDnRequest() == nil {
				t.Fatal("expected modify_dn_request to be set")
			}
			got := opReq.GetModifyDnRequest()
			if got.GetDn() != tt.dn {
				t.Errorf("got dn %q, want %q", got.GetDn(), tt.dn)
			}
			if got.GetNewRdn() != tt.newRDN {
				t.Errorf("got new_rdn %q, want %q", got.GetNewRdn(), tt.newRDN)
			}
			if got.GetDeleteOldRdn() != tt.deleteOldRDN {
				t.Errorf("got delete_old_rdn %v, want %v", got.GetDeleteOldRdn(), tt.deleteOldRDN)
			}
			if got.GetNewSuperior() != tt.newSuperior {
				t.Errorf("got new_superior %q, want %q", got.GetNewSuperior(), tt.newSuperior)
			}
		})
	}
}

func TestLDAPOperationResponse_GivenResultFields_WhenConstructed_ThenFieldsAreCorrect(t *testing.T) {
	tests := []struct {
		name       string
		errMsg     string
		resultCode int32
		matchedDN  string
	}{
		{
			name:       "success_result_code_zero",
			errMsg:     "",
			resultCode: 0,
			matchedDN:  "",
		},
		{
			name:       "error_invalid_credentials",
			errMsg:     "invalid credentials",
			resultCode: 49,
			matchedDN:  "",
		},
		{
			name:       "error_no_such_object",
			errMsg:     "no such object",
			resultCode: 32,
			matchedDN:  "dc=example,dc=com",
		},
		{
			name:       "error_insufficient_access",
			errMsg:     "insufficient access rights",
			resultCode: 50,
			matchedDN:  "cn=John,ou=Users,dc=example,dc=com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange & Act
			resp := &LDAPOperationResponse{
				Error:      tt.errMsg,
				ResultCode: tt.resultCode,
				MatchedDn:  tt.matchedDN,
			}

			// Assert - success is derived from result_code == 0, no explicit success field
			if resp.GetError() != tt.errMsg {
				t.Errorf("got error %q, want %q", resp.GetError(), tt.errMsg)
			}
			if resp.GetResultCode() != tt.resultCode {
				t.Errorf("got result_code %d, want %d", resp.GetResultCode(), tt.resultCode)
			}
			if resp.GetMatchedDn() != tt.matchedDN {
				t.Errorf("got matched_dn %q, want %q", resp.GetMatchedDn(), tt.matchedDN)
			}
		})
	}
}

func TestLDAPModifyOperation_GivenEnumValues_WhenChecked_ThenNumericValuesAreStable(t *testing.T) {
	// Verify enum numeric values remain stable across regenerations.
	// These values are chosen to align with go-ldap/ldap/v3 constants
	// (AddAttribute=0, DeleteAttribute=1, ReplaceAttribute=2) by convention,
	// but this test validates the proto-generated values, not the go-ldap library directly.

	tests := []struct {
		name     string
		op       LDAPModifyOperation
		expected int32
	}{
		{
			name:     "add_is_zero",
			op:       LDAPModifyOperation_LDAP_MODIFY_OPERATION_ADD,
			expected: 0,
		},
		{
			name:     "delete_is_one",
			op:       LDAPModifyOperation_LDAP_MODIFY_OPERATION_DELETE,
			expected: 1,
		},
		{
			name:     "replace_is_two",
			op:       LDAPModifyOperation_LDAP_MODIFY_OPERATION_REPLACE,
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			got := int32(tt.op)

			// Assert
			if got != tt.expected {
				t.Errorf("got %d, want %d", got, tt.expected)
			}
		})
	}
}

func TestRequest_GivenLDAPOperationRequest_WhenOneofSet_ThenGetReturnsCorrectType(t *testing.T) {
	// Arrange
	addReq := &LDAPAddRequest{
		Dn: "cn=Test,dc=example,dc=com",
		Attributes: []*LDAPAttribute{
			{Type: "objectClass", Values: []string{"user"}},
		},
	}
	opReq := &LDAPOperationRequest{
		Url:          "ldaps://ad.example.com:636",
		BindDn:       "cn=admin,dc=example,dc=com",
		BindPassword: "secret",
		Operation:    &LDAPOperationRequest_AddRequest{AddRequest: addReq},
	}

	// Act
	req := &Request{
		RequestType: &Request_LdapOperationRequest{LdapOperationRequest: opReq},
	}

	// Assert
	if req.GetLdapOperationRequest() == nil {
		t.Fatal("expected ldap_operation_request to be set")
	}
	if req.GetHttpRequest() != nil {
		t.Error("expected http_request to be nil")
	}
	if req.GetSqlQueryReq() != nil {
		t.Error("expected sql_query_req to be nil")
	}
	if req.GetLdapSearchRequest() != nil {
		t.Error("expected ldap_search_request to be nil")
	}
	if req.GetLdapOperationRequest().GetUrl() != "ldaps://ad.example.com:636" {
		t.Errorf("got url %q, want %q", req.GetLdapOperationRequest().GetUrl(), "ldaps://ad.example.com:636")
	}
}

func TestResponse_GivenLDAPOperationResponse_WhenOneofSet_ThenGetReturnsCorrectType(t *testing.T) {
	// Arrange
	opResp := &LDAPOperationResponse{
		ResultCode: 0,
	}

	// Act
	resp := &Response{
		ResponseType: &Response_LdapOperationResponse{LdapOperationResponse: opResp},
	}

	// Assert
	if resp.GetLdapOperationResponse() == nil {
		t.Fatal("expected ldap_operation_response to be set")
	}
	if resp.GetHttpResponse() != nil {
		t.Error("expected http_response to be nil")
	}
	if resp.GetSqlQueryResponse() != nil {
		t.Error("expected sql_query_response to be nil")
	}
	if resp.GetLdapSearchResponse() != nil {
		t.Error("expected ldap_search_response to be nil")
	}
	if resp.GetLdapOperationResponse().GetResultCode() != 0 {
		t.Errorf("expected result_code 0, got %d", resp.GetLdapOperationResponse().GetResultCode())
	}
}

func TestLDAPOperationRequest_GivenAddRequest_WhenMarshalUnmarshal_ThenRoundTripsCorrectly(t *testing.T) {
	// Arrange
	original := &LDAPOperationRequest{
		Url:          "ldaps://ad.example.com:636",
		BindDn:       "cn=admin,dc=example,dc=com",
		BindPassword: "secret",
		Operation: &LDAPOperationRequest_AddRequest{
			AddRequest: &LDAPAddRequest{
				Dn: "cn=John Doe,ou=Users,dc=example,dc=com",
				Attributes: []*LDAPAttribute{
					{Type: "objectClass", Values: []string{"user", "person"}},
					{Type: "sAMAccountName", Values: []string{"jdoe"}},
					{Type: "cn", Values: []string{"John Doe"}},
				},
			},
		},
	}

	// Act
	data, err := proto.Marshal(original)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	decoded := &LDAPOperationRequest{}
	if err := proto.Unmarshal(data, decoded); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	// Assert
	if decoded.GetUrl() != original.GetUrl() {
		t.Errorf("url: got %q, want %q", decoded.GetUrl(), original.GetUrl())
	}
	if decoded.GetBindDn() != original.GetBindDn() {
		t.Errorf("bind_dn: got %q, want %q", decoded.GetBindDn(), original.GetBindDn())
	}
	if decoded.GetBindPassword() != original.GetBindPassword() {
		t.Errorf("bind_password: got %q, want %q", decoded.GetBindPassword(), original.GetBindPassword())
	}
	if decoded.GetAddRequest() == nil {
		t.Fatal("expected add_request to be set after unmarshal")
	}
	if decoded.GetAddRequest().GetDn() != original.GetAddRequest().GetDn() {
		t.Errorf("dn: got %q, want %q", decoded.GetAddRequest().GetDn(), original.GetAddRequest().GetDn())
	}
	if len(decoded.GetAddRequest().GetAttributes()) != len(original.GetAddRequest().GetAttributes()) {
		t.Fatalf("attributes count: got %d, want %d", len(decoded.GetAddRequest().GetAttributes()), len(original.GetAddRequest().GetAttributes()))
	}
	for i, attr := range decoded.GetAddRequest().GetAttributes() {
		origAttr := original.GetAddRequest().GetAttributes()[i]
		if attr.GetType() != origAttr.GetType() {
			t.Errorf("attribute[%d].type: got %q, want %q", i, attr.GetType(), origAttr.GetType())
		}
		if len(attr.GetValues()) != len(origAttr.GetValues()) {
			t.Errorf("attribute[%d].values count: got %d, want %d", i, len(attr.GetValues()), len(origAttr.GetValues()))
		}
	}
}

func TestLDAPOperationRequest_GivenModifyRequest_WhenMarshalUnmarshal_ThenRoundTripsCorrectly(t *testing.T) {
	// Arrange
	original := &LDAPOperationRequest{
		Url:          "ldaps://ad.example.com:636",
		BindDn:       "cn=admin,dc=example,dc=com",
		BindPassword: "secret",
		Operation: &LDAPOperationRequest_ModifyRequest{
			ModifyRequest: &LDAPModifyRequest{
				Dn: "cn=John Doe,ou=Users,dc=example,dc=com",
				Changes: []*LDAPModifyChange{
					{
						Operation:    LDAPModifyOperation_LDAP_MODIFY_OPERATION_REPLACE,
						Modification: &LDAPAttribute{Type: "userAccountControl", Values: []string{"512"}},
					},
					{
						Operation:    LDAPModifyOperation_LDAP_MODIFY_OPERATION_ADD,
						Modification: &LDAPAttribute{Type: "mail", Values: []string{"john@example.com"}},
					},
				},
			},
		},
	}

	// Act
	data, err := proto.Marshal(original)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	decoded := &LDAPOperationRequest{}
	if err := proto.Unmarshal(data, decoded); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	// Assert
	if decoded.GetModifyRequest() == nil {
		t.Fatal("expected modify_request to be set after unmarshal")
	}
	if decoded.GetModifyRequest().GetDn() != original.GetModifyRequest().GetDn() {
		t.Errorf("dn: got %q, want %q", decoded.GetModifyRequest().GetDn(), original.GetModifyRequest().GetDn())
	}
	changes := decoded.GetModifyRequest().GetChanges()
	origChanges := original.GetModifyRequest().GetChanges()
	if len(changes) != len(origChanges) {
		t.Fatalf("changes count: got %d, want %d", len(changes), len(origChanges))
	}
	for i, change := range changes {
		if change.GetOperation() != origChanges[i].GetOperation() {
			t.Errorf("change[%d].operation: got %v, want %v", i, change.GetOperation(), origChanges[i].GetOperation())
		}
		if change.GetModification().GetType() != origChanges[i].GetModification().GetType() {
			t.Errorf("change[%d].modification.type: got %q, want %q", i, change.GetModification().GetType(), origChanges[i].GetModification().GetType())
		}
	}
}

func TestLDAPOperationRequest_GivenModifyDNRequest_WhenMarshalUnmarshal_ThenRoundTripsCorrectly(t *testing.T) {
	// Arrange
	original := &LDAPOperationRequest{
		Url:          "ldaps://ad.example.com:636",
		BindDn:       "cn=admin,dc=example,dc=com",
		BindPassword: "secret",
		Operation: &LDAPOperationRequest_ModifyDnRequest{
			ModifyDnRequest: &LDAPModifyDNRequest{
				Dn:           "cn=John Doe,ou=Users,dc=example,dc=com",
				NewRdn:       "cn=John Doe",
				DeleteOldRdn: true,
				NewSuperior:  "ou=ArchivedUsers,dc=example,dc=com",
			},
		},
	}

	// Act
	data, err := proto.Marshal(original)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	decoded := &LDAPOperationRequest{}
	if err := proto.Unmarshal(data, decoded); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	// Assert
	if decoded.GetModifyDnRequest() == nil {
		t.Fatal("expected modify_dn_request to be set after unmarshal")
	}
	got := decoded.GetModifyDnRequest()
	orig := original.GetModifyDnRequest()
	if got.GetDn() != orig.GetDn() {
		t.Errorf("dn: got %q, want %q", got.GetDn(), orig.GetDn())
	}
	if got.GetNewRdn() != orig.GetNewRdn() {
		t.Errorf("new_rdn: got %q, want %q", got.GetNewRdn(), orig.GetNewRdn())
	}
	if got.GetDeleteOldRdn() != orig.GetDeleteOldRdn() {
		t.Errorf("delete_old_rdn: got %v, want %v", got.GetDeleteOldRdn(), orig.GetDeleteOldRdn())
	}
	if got.GetNewSuperior() != orig.GetNewSuperior() {
		t.Errorf("new_superior: got %q, want %q", got.GetNewSuperior(), orig.GetNewSuperior())
	}
}

func TestLDAPOperationResponse_GivenResponse_WhenMarshalUnmarshal_ThenRoundTripsCorrectly(t *testing.T) {
	tests := []struct {
		name     string
		original *LDAPOperationResponse
	}{
		{
			name: "success_response",
			original: &LDAPOperationResponse{
				ResultCode: 0,
			},
		},
		{
			name: "error_response",
			original: &LDAPOperationResponse{
				Error:      "no such object",
				ResultCode: 32,
				MatchedDn:  "dc=example,dc=com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			data, err := proto.Marshal(tt.original)
			if err != nil {
				t.Fatalf("marshal failed: %v", err)
			}

			decoded := &LDAPOperationResponse{}
			if err := proto.Unmarshal(data, decoded); err != nil {
				t.Fatalf("unmarshal failed: %v", err)
			}

			// Assert
			if decoded.GetError() != tt.original.GetError() {
				t.Errorf("error: got %q, want %q", decoded.GetError(), tt.original.GetError())
			}
			if decoded.GetResultCode() != tt.original.GetResultCode() {
				t.Errorf("result_code: got %d, want %d", decoded.GetResultCode(), tt.original.GetResultCode())
			}
			if decoded.GetMatchedDn() != tt.original.GetMatchedDn() {
				t.Errorf("matched_dn: got %q, want %q", decoded.GetMatchedDn(), tt.original.GetMatchedDn())
			}
		})
	}
}

func TestRequest_GivenLDAPOperationRequest_WhenMarshalUnmarshal_ThenOneofPreserved(t *testing.T) {
	// Arrange
	original := &Request{
		RequestType: &Request_LdapOperationRequest{
			LdapOperationRequest: &LDAPOperationRequest{
				Url:          "ldaps://ad.example.com:636",
				BindDn:       "cn=admin,dc=example,dc=com",
				BindPassword: "secret",
				Operation: &LDAPOperationRequest_DeleteRequest{
					DeleteRequest: &LDAPDeleteRequest{
						Dn: "cn=John Doe,ou=Users,dc=example,dc=com",
					},
				},
			},
		},
	}

	// Act
	data, err := proto.Marshal(original)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	decoded := &Request{}
	if err := proto.Unmarshal(data, decoded); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	// Assert
	if decoded.GetLdapOperationRequest() == nil {
		t.Fatal("expected ldap_operation_request to be set after unmarshal")
	}
	if decoded.GetLdapOperationRequest().GetDeleteRequest() == nil {
		t.Fatal("expected delete_request to be set after unmarshal")
	}
	if decoded.GetLdapOperationRequest().GetDeleteRequest().GetDn() != "cn=John Doe,ou=Users,dc=example,dc=com" {
		t.Errorf("got dn %q, want %q",
			decoded.GetLdapOperationRequest().GetDeleteRequest().GetDn(),
			"cn=John Doe,ou=Users,dc=example,dc=com")
	}
}

func TestResponse_GivenLDAPOperationResponse_WhenMarshalUnmarshal_ThenOneofPreserved(t *testing.T) {
	// Arrange
	original := &Response{
		ResponseType: &Response_LdapOperationResponse{
			LdapOperationResponse: &LDAPOperationResponse{
				Error:      "insufficient access rights",
				ResultCode: 50,
				MatchedDn:  "cn=John,ou=Users,dc=example,dc=com",
			},
		},
	}

	// Act
	data, err := proto.Marshal(original)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	decoded := &Response{}
	if err := proto.Unmarshal(data, decoded); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	// Assert
	if decoded.GetLdapOperationResponse() == nil {
		t.Fatal("expected ldap_operation_response to be set after unmarshal")
	}
	got := decoded.GetLdapOperationResponse()
	if got.GetResultCode() != 50 {
		t.Errorf("got result_code %d, want %d", got.GetResultCode(), 50)
	}
	if got.GetError() != "insufficient access rights" {
		t.Errorf("got error %q, want %q", got.GetError(), "insufficient access rights")
	}
	if got.GetMatchedDn() != "cn=John,ou=Users,dc=example,dc=com" {
		t.Errorf("got matched_dn %q, want %q", got.GetMatchedDn(), "cn=John,ou=Users,dc=example,dc=com")
	}
}
