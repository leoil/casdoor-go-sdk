package casdoorsdk

import (
	_ "embed"
	"fmt"
	"testing"
)

func TestAddGroup(t *testing.T) {
	InitConfigTest()

	// 定义测试用例
	testCases := []struct {
		group         *Group
		expected      bool
		expectedError error
	}{
		{
			group: &Group{
				Owner:       "casbin-forum",
				Name:        "test-group1",
				DisplayName: "Group-Test1",
			},
			expected:      true,
			expectedError: nil,
		},
		{
			group: &Group{
				Owner:       "casbin-forum",
				Name:        "test-group2",
				DisplayName: "Group-Test2",
			},
			expected:      true,
			expectedError: nil,
		},
	}

	// 循环遍历测试用例
	for _, tc := range testCases {
		// 调用被测试的方法
		ok, err := localClient.AddGroup(tc.group)

		if !ok {
			t.Errorf("For owner %s and name %s, group add failed", tc.group.Owner, tc.group.Name)
		}
		// 断言返回的错误是否符合预期
		if err != tc.expectedError {
			t.Errorf("Expected error %v, but got %v", tc.expectedError, err)
		}
	}
}

func TestGetGroup(t *testing.T) {
	InitConfigTest()

	testCases := []struct {
		name          string
		expected      *Group
		expectedError error
	}{
		{
			name:          "test-group1",
			expected:      &Group{Owner: "casbin-forum", Name: "test-group1", DisplayName: "Group-Test1"},
			expectedError: nil,
		},
		{
			name:          "test-group2",
			expected:      &Group{Owner: "casbin-forum", Name: "test-group2", DisplayName: "Group-Test2"},
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
		group, err := localClient.GetGroup(tc.name)

		fieldsToCompare := []field{
			{"Owner", tc.expected.Owner, group.Owner},
			{"Name", tc.expected.Name, group.Name},
			{"DisplayName", tc.expected.DisplayName, group.DisplayName},
		}

		for _, f := range fieldsToCompare {
			if f.expected != f.actual {
				t.Errorf("For field %s, expected %v, but got %v", f.name, f.expected, f.actual)
			}
		}
		if err != tc.expectedError {
			t.Errorf("For group %s, expected error %v, but got %v", tc.name, tc.expectedError, err)
		}
	}
}

func TestGetGroups(t *testing.T) {
	InitConfigTest()

	InitConfigTest()

	groups, err := localClient.GetGroups()
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
		return
	}

	//
	fmt.Println("Number of groups:", len(groups))
}

func TestUpdateGroup(t *testing.T) {
	InitConfigTest()

	testCases := []struct {
		name          string
		expected      *Group
		expectedError error
	}{
		{
			name:          "test-group1",
			expected:      &Group{Owner: "casbin-forum", Name: "test-group1", DisplayName: "Group-Test11"},
			expectedError: nil,
		},

		// 添加其他测试用例和期望值
	}

	// 循环遍历测试用例
	for _, tc := range testCases {
		// 调用被测试的方法
		// group, err := localClient.GetGroup(tc.name)
		ok, err := localClient.UpdateGroup(tc.expected)

		if !ok {
			t.Errorf("For group %s, update failed", tc.name)
		}

		if err != tc.expectedError {
			t.Errorf("For group %s, expected error %v, but got %v", tc.name, tc.expectedError, err)
		}
	}
}

func TestDeleteGroup(t *testing.T) {
	InitConfigTest()

	testCases := []struct {
		name          string
		expected      bool
		group         *Group
		expectedError error
	}{
		{
			name:          "test-group1",
			group:         &Group{Owner: "casbin-forum", Name: "test-group1", DisplayName: "Group-Test2"},
			expected:      true,
			expectedError: nil,
		},
		{
			name:          "test-group2",
			group:         &Group{Owner: "casbin-forum", Name: "test-group2", DisplayName: "Group-Test2"},
			expected:      true,
			expectedError: nil,
		},
	}

	// 循环遍历测试用例
	for _, tc := range testCases {

		ok, err := localClient.DeleteGroup(tc.group)

		if !ok {
			t.Errorf("For group %s, fail to delete", tc.name)
		}

		if err != tc.expectedError {
			t.Errorf("For group %s, expected error %v, but got %v", tc.name, tc.expectedError, err)
		}
	}
}
