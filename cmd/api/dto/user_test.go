package dto

import (
	"encoding/json"
	"testing"
)

func TestUserUpdateRequestDistinguishesOmittedAndEmptyFields(t *testing.T) {
	var request UserUpdateRequest
	if err := json.Unmarshal([]byte(`{"user":{"email":"new@example.com","bio":""}}`), &request); err != nil {
		t.Fatal(err)
	}

	if request.User.Email == nil || *request.User.Email != "new@example.com" {
		t.Fatalf("email = %v, want new@example.com", request.User.Email)
	}
	if request.User.Bio == nil || *request.User.Bio != "" {
		t.Fatalf("bio = %v, want pointer to empty string", request.User.Bio)
	}
	if request.User.Username != nil || request.User.Image != nil || request.User.Password != nil {
		t.Fatal("omitted fields should remain nil")
	}
}

func TestUserUpdateRejectsEmptyIdentityField(t *testing.T) {
	empty := ""
	request := UserUpdateRequest{}
	request.User.Email = &empty

	if validationErrors := ValidateUserUpdateRequest(&request); !validationErrors.HasAny() {
		t.Fatal("expected an empty email to fail validation")
	}
}
