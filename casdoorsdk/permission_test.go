package casdoorsdk

import (
	_ "embed"
	"fmt"
	"testing"
)

func TestAddPermission(t *testing.T) {
	InitConfigTest()

	// 定义测试用例
	testCases := []struct {
		permission    *Permission
		expected      bool
		expectedError error
	}{
		{
			permission: &Permission{
				Owner:       "casbin-forum",
				Name:        "test-perm1",
				DisplayName: "Permission-Test1",
			},
			expected:      true,
			expectedError: nil,
		},
		{
			permission: &Permission{
				Owner:       "casbin-forum",
				Name:        "test-perm2",
				DisplayName: "Permission-Test2",
			},
			expected:      true,
			expectedError: nil,
		},
	}

	// 循环遍历测试用例
	for _, tc := range testCases {
		// 调用被测试的方法
		ok, err := localClient.AddPermission(tc.permission)

		if !ok {
			t.Errorf("For owner %s and name %s, permission add failed", tc.permission.Owner, tc.permission.Name)
		}
		// 断言返回的错误是否符合预期
		if err != tc.expectedError {
			t.Errorf("Expected error %v, but got %v", tc.expectedError, err)
		}
	}
}

func TestGetPermission(t *testing.T) {
	InitConfigTest()

	testCases := []struct {
		name          string
		expected      *Permission
		expectedError error
	}{
		{
			name:          "test-perm1",
			expected:      &Permission{Owner: "casbin-forum", Name: "test-perm1", DisplayName: "Permission-Test1"},
			expectedError: nil,
		},
		{
			name:          "test-perm2",
			expected:      &Permission{Owner: "casbin-forum", Name: "test-perm2", DisplayName: "Permission-Test2"},
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
		permission, err := localClient.GetPermission(tc.name)

		fieldsToCompare := []field{
			{"Owner", tc.expected.Owner, permission.Owner},
			{"Name", tc.expected.Name, permission.Name},
			{"DisplayName", tc.expected.DisplayName, permission.DisplayName},
		}

		for _, f := range fieldsToCompare {
			if f.expected != f.actual {
				t.Errorf("For field %s, expected %v, but got %v", f.name, f.expected, f.actual)
			}
		}
		if err != tc.expectedError {
			t.Errorf("For permission %s, expected error %v, but got %v", tc.name, tc.expectedError, err)
		}
	}
}

func TestGetPermissions(t *testing.T) {
	InitConfigTest()

	InitConfigTest()

	permissions, err := localClient.GetPermissions()
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
		return
	}

	//
	fmt.Println("Number of permissions:", len(permissions))
}

func TestUpdatePermission(t *testing.T) {
	InitConfigTest()

	testCases := []struct {
		name          string
		expected      *Permission
		expectedError error
	}{
		{
			name:          "test-perm1",
			expected:      &Permission{Owner: "casbin-forum", Name: "test-perm1", DisplayName: "Permission-Test11"},
			expectedError: nil,
		},

		// 添加其他测试用例和期望值
	}

	// 循环遍历测试用例
	for _, tc := range testCases {
		// 调用被测试的方法
		// permission, err := localClient.GetPermission(tc.name)
		ok, err := localClient.UpdatePermission(tc.expected)

		if !ok {
			t.Errorf("For permission %s, update failed", tc.name)
		}

		if err != tc.expectedError {
			t.Errorf("For permission %s, expected error %v, but got %v", tc.name, tc.expectedError, err)
		}
	}
}

func TestDeletePermission(t *testing.T) {
	InitConfigTest()

	testCases := []struct {
		name          string
		expected      bool
		permission    *Permission
		expectedError error
	}{
		{
			name:          "test-perm1",
			permission:    &Permission{Owner: "casbin-forum", Name: "test-perm1", DisplayName: "Permission-Test2"},
			expected:      true,
			expectedError: nil,
		},
		{
			name:          "test-perm2",
			permission:    &Permission{Owner: "casbin-forum", Name: "test-perm2", DisplayName: "Permission-Test2"},
			expected:      true,
			expectedError: nil,
		},
	}

	// 循环遍历测试用例
	for _, tc := range testCases {

		ok, err := localClient.DeletePermission(tc.permission)

		if !ok {
			t.Errorf("For permission %s, fail to delete", tc.name)
		}

		if err != tc.expectedError {
			t.Errorf("For permission %s, expected error %v, but got %v", tc.name, tc.expectedError, err)
		}
	}
}
