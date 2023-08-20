package casdoorsdk

import (
	_ "embed"
	"fmt"
	"testing"
)

func TestAddEnforcer(t *testing.T) {
	InitConfigTest()

	// 定义测试用例
	testCases := []struct {
		enforcer      *Enforcer
		expected      bool
		expectedError error
	}{
		{
			enforcer: &Enforcer{
				Owner:       "casbin-forum",
				Name:        "test-enforcer1",
				DisplayName: "Enforcer-Test1",
			},
			expected:      true,
			expectedError: nil,
		},
		{
			enforcer: &Enforcer{
				Owner:       "casbin-forum",
				Name:        "test-enforcer2",
				DisplayName: "Enforcer-Test2",
			},
			expected:      true,
			expectedError: nil,
		},
	}

	// 循环遍历测试用例
	for _, tc := range testCases {
		// 调用被测试的方法
		ok, err := localClient.AddEnforcer(tc.enforcer)

		if !ok {
			t.Errorf("For owner %s and name %s, enforcer add failed", tc.enforcer.Owner, tc.enforcer.Name)
		}
		// 断言返回的错误是否符合预期
		if err != tc.expectedError {
			t.Errorf("Expected error %v, but got %v", tc.expectedError, err)
		}
	}
}

func TestGetEnforcer(t *testing.T) {
	InitConfigTest()

	testCases := []struct {
		name          string
		expected      *Enforcer
		expectedError error
	}{
		{
			name:          "test-enforcer1",
			expected:      &Enforcer{Owner: "casbin-forum", Name: "test-enforcer1", DisplayName: "Enforcer-Test1"},
			expectedError: nil,
		},
		{
			name:          "test-enforcer2",
			expected:      &Enforcer{Owner: "casbin-forum", Name: "test-enforcer2", DisplayName: "Enforcer-Test2"},
			expectedError: nil,
		},
		// 添加其他测试用例和期望值
	}

	// 使用循环比较字段值
	type field struct {
		name     string
		expected interface{}
		actual   interface{}
	}

	// 循环遍历测试用例
	for _, tc := range testCases {
		// 调用被测试的方法
		enforcer, err := localClient.GetEnforcer(tc.name)

		fieldsToCompare := []field{
			{"Owner", tc.expected.Owner, enforcer.Owner},
			{"Name", tc.expected.Name, enforcer.Name},
			{"DisplayName", tc.expected.DisplayName, enforcer.DisplayName},
		}

		for _, f := range fieldsToCompare {
			if f.expected != f.actual {
				t.Errorf("For field %s, expected %v, but got %v", f.name, f.expected, f.actual)
			}
		}
		if err != tc.expectedError {
			t.Errorf("For enforcer %s, expected error %v, but got %v", tc.name, tc.expectedError, err)
		}
	}
}

func TestGetEnforcers(t *testing.T) {
	InitConfigTest()

	InitConfigTest()

	enforcers, err := localClient.GetEnforcers()
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
		return
	}

	//
	fmt.Println("Number of enforcers:", len(enforcers))
}

func TestUpdateEnforcer(t *testing.T) {
	InitConfigTest()

	testCases := []struct {
		name          string
		expected      *Enforcer
		expectedError error
	}{
		{
			name:          "test-enforcer1",
			expected:      &Enforcer{Owner: "casbin-forum", Name: "test-enforcer1", DisplayName: "Enforcer-Test11"},
			expectedError: nil,
		},

		// 添加其他测试用例和期望值
	}

	// 循环遍历测试用例
	for _, tc := range testCases {
		// 调用被测试的方法
		// enforcer, err := localClient.GetEnforcer(tc.name)
		ok, err := localClient.UpdateEnforcer(tc.expected)

		if !ok {
			t.Errorf("For enforcer %s, update failed", tc.name)
		}

		if err != tc.expectedError {
			t.Errorf("For enforcer %s, expected error %v, but got %v", tc.name, tc.expectedError, err)
		}
	}
}

func TestDeleteEnforcer(t *testing.T) {
	InitConfigTest()

	testCases := []struct {
		name          string
		expected      bool
		enforcer      *Enforcer
		expectedError error
	}{
		{
			name:          "test-enforcer1",
			enforcer:      &Enforcer{Owner: "casbin-forum", Name: "test-enforcer1", DisplayName: "Enforcer-Test2"},
			expected:      true,
			expectedError: nil,
		},
		{
			name:          "test-enforcer2",
			enforcer:      &Enforcer{Owner: "casbin-forum", Name: "test-enforcer2", DisplayName: "Enforcer-Test2"},
			expected:      true,
			expectedError: nil,
		},
	}

	// 循环遍历测试用例
	for _, tc := range testCases {

		ok, err := localClient.DeleteEnforcer(tc.enforcer)

		if !ok {
			t.Errorf("For enforcer %s, fail to delete", tc.name)
		}

		if err != tc.expectedError {
			t.Errorf("For enforcer %s, expected error %v, but got %v", tc.name, tc.expectedError, err)
		}
	}
}
