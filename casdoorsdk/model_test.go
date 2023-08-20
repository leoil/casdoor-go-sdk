package casdoorsdk

import (
	_ "embed"
	"fmt"
	"testing"
)

func TestAddModel(t *testing.T) {
	InitConfigTest()

	// 定义测试用例
	testCases := []struct {
		model         *Model
		expected      bool
		expectedError error
	}{
		{
			model: &Model{
				Owner:       "casbin-forum",
				Name:        "test-model1",
				DisplayName: "Model-Test1",
			},
			expected:      true,
			expectedError: nil,
		},
		{
			model: &Model{
				Owner:       "casbin-forum",
				Name:        "test-model2",
				DisplayName: "Model-Test2",
			},
			expected:      true,
			expectedError: nil,
		},
	}

	// 循环遍历测试用例
	for _, tc := range testCases {
		// 调用被测试的方法
		ok, err := localClient.AddModel(tc.model)

		if !ok {
			t.Errorf("For owner %s and name %s, model add failed", tc.model.Owner, tc.model.Name)
		}
		// 断言返回的错误是否符合预期
		if err != tc.expectedError {
			t.Errorf("Expected error %v, but got %v", tc.expectedError, err)
		}
	}
}

func TestGetModel(t *testing.T) {
	InitConfigTest()

	testCases := []struct {
		name          string
		expected      *Model
		expectedError error
	}{
		{
			name:          "test-model1",
			expected:      &Model{Owner: "casbin-forum", Name: "test-model1", DisplayName: "Model-Test1"},
			expectedError: nil,
		},
		{
			name:          "test-model2",
			expected:      &Model{Owner: "casbin-forum", Name: "test-model2", DisplayName: "Model-Test2"},
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
		model, err := localClient.GetModel(tc.name)

		fieldsToCompare := []field{
			{"Owner", tc.expected.Owner, model.Owner},
			{"Name", tc.expected.Name, model.Name},
			{"DisplayName", tc.expected.DisplayName, model.DisplayName},
		}

		for _, f := range fieldsToCompare {
			if f.expected != f.actual {
				t.Errorf("For field %s, expected %v, but got %v", f.name, f.expected, f.actual)
			}
		}
		if err != tc.expectedError {
			t.Errorf("For model %s, expected error %v, but got %v", tc.name, tc.expectedError, err)
		}
	}
}

func TestGetModels(t *testing.T) {
	InitConfigTest()

	InitConfigTest()

	models, err := localClient.GetModels()
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
		return
	}

	//
	fmt.Println("Number of models:", len(models))
}

func TestUpdateModel(t *testing.T) {
	InitConfigTest()

	testCases := []struct {
		name          string
		expected      *Model
		expectedError error
	}{
		{
			name:          "test-model1",
			expected:      &Model{Owner: "casbin-forum", Name: "test-model1", DisplayName: "Model-Test11"},
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		// model, err := localClient.GetModel(tc.name)
		ok, err := localClient.UpdateModel(tc.expected)

		if !ok {
			t.Errorf("For model %s, update failed", tc.name)
		}

		if err != tc.expectedError {
			t.Errorf("For model %s, expected error %v, but got %v", tc.name, tc.expectedError, err)
		}
	}
}

func TestDeleteModel(t *testing.T) {
	InitConfigTest()

	testCases := []struct {
		name          string
		expected      bool
		model         *Model
		expectedError error
	}{
		{
			name:          "test-model1",
			model:         &Model{Owner: "casbin-forum", Name: "test-model1", DisplayName: "Model-Test2"},
			expected:      true,
			expectedError: nil,
		},
		{
			name:          "test-model2",
			model:         &Model{Owner: "casbin-forum", Name: "test-model2", DisplayName: "Model-Test2"},
			expected:      true,
			expectedError: nil,
		},
	}

	// 循环遍历测试用例
	for _, tc := range testCases {

		ok, err := localClient.DeleteModel(tc.model)

		if !ok {
			t.Errorf("For model %s, fail to delete", tc.name)
		}

		if err != tc.expectedError {
			t.Errorf("For model %s, expected error %v, but got %v", tc.name, tc.expectedError, err)
		}
	}
}
