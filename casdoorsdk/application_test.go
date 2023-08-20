package casdoorsdk

import (
	_ "embed"
	"fmt"
	"testing"
)

func TestAddApplication(t *testing.T) {
	InitConfigTest()

	testCases := []struct {
		application   *Application
		expected      bool
		expectedError error
	}{
		{
			application: &Application{
				Owner:        "admin",
				Name:         "test-app1",
				DisplayName:  "Application-Test1",
				Organization: "casbin-forum",
			},
			expected:      true,
			expectedError: nil,
		},
		{
			application: &Application{
				Owner:        "admin",
				Name:         "test-app2",
				DisplayName:  "Application-Test2",
				Organization: "casbin-forum",
			},
			expected:      true,
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		ok, err := localClient.AddApplication(tc.application)

		if !ok {
			t.Errorf("For owner %s and name %s, application add failed", tc.application.Owner, tc.application.Name)
		}
		if err != tc.expectedError {
			t.Errorf("Expected error %v, but got %v", tc.expectedError, err)
		}
	}
}

func TestGetApplication(t *testing.T) {
	InitConfigTest()

	testCases := []struct {
		name          string
		expected      *Application
		expectedError error
	}{
		{
			name:          "test-app1",
			expected:      &Application{Owner: "admin", Name: "test-app1", DisplayName: "Application-Test1", Organization: "casbin-forum"},
			expectedError: nil,
		},
		{
			name:          "test-app2",
			expected:      &Application{Owner: "admin", Name: "test-app2", DisplayName: "Application-Test2", Organization: "casbin-forum"},
			expectedError: nil,
		},
	}

	type field struct {
		name     string
		expected interface{}
		actual   interface{}
	}

	for _, tc := range testCases {
		application, err := localClient.GetApplication(tc.name)

		fieldsToCompare := []field{
			{"Owner", tc.expected.Owner, application.Owner},
			{"Name", tc.expected.Name, application.Name},
			{"DisplayName", tc.expected.DisplayName, application.DisplayName},
		}

		for _, f := range fieldsToCompare {
			if f.expected != f.actual {
				t.Errorf("For field %s, expected %v, but got %v", f.name, f.expected, f.actual)
			}
		}
		if err != tc.expectedError {
			t.Errorf("For application %s, expected error %v, but got %v", tc.name, tc.expectedError, err)
		}
	}
}

func TestGetApplications(t *testing.T) {
	InitConfigTest()

	InitConfigTest()

	applications, err := localClient.GetApplications()
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
		return
	}

	fmt.Println("Number of applications:", len(applications))
}

func TestUpdateApplication(t *testing.T) {
	InitConfigTest()

	testCases := []struct {
		name          string
		expected      *Application
		expectedError error
	}{
		{
			name:          "test-app1",
			expected:      &Application{Owner: "admin", Name: "test-app1", DisplayName: "Application-Test11"},
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		ok, err := localClient.UpdateApplication(tc.expected)

		if !ok {
			t.Errorf("For application %s, update failed", tc.name)
		}

		if err != tc.expectedError {
			t.Errorf("For application %s, expected error %v, but got %v", tc.name, tc.expectedError, err)
		}
	}
}

func TestDeleteApplication(t *testing.T) {
	InitConfigTest()

	testCases := []struct {
		name          string
		expected      bool
		expectedError error
	}{
		{
			name:          "test-app1",
			expected:      true,
			expectedError: nil,
		},
		{
			name:          "test-app2",
			expected:      true,
			expectedError: nil,
		},
	}

	for _, tc := range testCases {

		ok, err := localClient.DeleteApplication(tc.name)

		if !ok {
			t.Errorf("For application %s, fail to delete", tc.name)
		}

		if err != tc.expectedError {
			t.Errorf("For application %s, expected error %v, but got %v", tc.name, tc.expectedError, err)
		}
	}
}
