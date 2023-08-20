package casdoorsdk

import (
	_ "embed"
	"fmt"
	"testing"
)

// crypto/rsa: too few primes of given length to generate an RSA key
func TestAddCert(t *testing.T) {
	InitConfigTest()

	// 定义测试用例
	testCases := []struct {
		Cert          *Cert
		expected      bool
		expectedError error
	}{
		{
			Cert:          &Cert{Owner: "admin", Name: "test-cert1", DisplayName: "Cert-Test1", BitSize: 64},
			expected:      true,
			expectedError: nil,
		},
		{
			Cert:          &Cert{Owner: "casbin-forum", Name: "test-cert2", DisplayName: "Cert-Test2"},
			expected:      true,
			expectedError: nil,
		},
	}

	// 循环遍历测试用例
	for _, tc := range testCases {
		// 调用被测试的方法
		ok, err := localClient.AddCert(tc.Cert)

		if !ok {
			t.Errorf("For owner %s and name %s, Cert add failed", tc.Cert.Owner, tc.Cert.Name)
		}
		// 断言返回的错误是否符合预期
		if err != tc.expectedError {
			t.Errorf("Expected error %v, but got %v", tc.expectedError, err)
		}
	}
}

func TestGetCert(t *testing.T) {
	InitConfigTest()

	testCases := []struct {
		name          string
		expected      *Cert
		expectedError error
	}{
		{
			name:          "test-cert2",
			expected:      &Cert{Owner: "casbin-forum", Name: "test-cert2", DisplayName: "Cert-Test2"},
			expectedError: nil,
		},
		//{
		//	name:          "test-cert2",
		//	expected:      &Cert{Owner: "casbin-forum", Name: "test-cert1", DisplayName: "Cert-Test2"},
		//	expectedError: nil,
		//},
		//{
		//	name:          "cert-built-in",
		//	expected:      &Cert{Owner: "admin", Name: "cert-built-in", DisplayName: "Built-in Cert"},
		//	expectedError: nil,
		//},
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
		Cert, err := localClient.GetCert(tc.name)

		fieldsToCompare := []field{
			{"Owner", tc.expected.Owner, Cert.Owner},
			{"Name", tc.expected.Name, Cert.Name},
			{"DisplayName", tc.expected.DisplayName, Cert.DisplayName},
		}

		for _, f := range fieldsToCompare {
			if f.expected != f.actual {
				t.Errorf("For field %s, expected %v, but got %v", f.name, f.expected, f.actual)
			}
		}
		if err != tc.expectedError {
			t.Errorf("For Cert %s, expected error %v, but got %v", tc.name, tc.expectedError, err)
		}
	}
}

func TestGetGlobalCerts(t *testing.T) {
	InitConfigTest()

	Certs, err := localClient.GetGlobalCerts()
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
		return
	}

	//
	fmt.Println("Number of Certs:", len(Certs))
}

func TestGetCerts(t *testing.T) {
	InitConfigTest()

	Certs, err := localClient.GetCerts()
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
		return
	}

	//
	fmt.Println("Number of Certs:", len(Certs))
}

func TestUpdateCert(t *testing.T) {
	InitConfigTest()

	testCases := []struct {
		name          string
		expected      *Cert
		expectedError error
	}{
		{
			name:          "test-perm11",
			expected:      &Cert{Owner: "casbin-forum", Name: "test-cert1", DisplayName: "Cert-Test11"},
			expectedError: nil,
		},

		// 添加其他测试用例和期望值
	}

	// 循环遍历测试用例
	for _, tc := range testCases {
		// 调用被测试的方法
		// Cert, err := localClient.GetCert(tc.name)
		ok, err := localClient.UpdateCert(tc.expected)

		if !ok {
			t.Errorf("For Cert %s, update failed", tc.name)
		}

		if err != tc.expectedError {
			t.Errorf("For Cert %s, expected error %v, but got %v", tc.name, tc.expectedError, err)
		}
	}
}

func TestDeleteCert(t *testing.T) {
	InitConfigTest()

	testCases := []struct {
		name          string
		expected      bool
		Cert          *Cert
		expectedError error
	}{
		{
			name:          "test-perm1",
			Cert:          &Cert{Owner: "casbin-forum", Name: "test-perm1", DisplayName: "Cert-Test2"},
			expected:      true,
			expectedError: nil,
		},
		{
			name:          "test-perm2",
			Cert:          &Cert{Owner: "casbin-forum", Name: "test-perm2", DisplayName: "Cert-Test2"},
			expected:      true,
			expectedError: nil,
		},
	}

	// 循环遍历测试用例
	for _, tc := range testCases {

		ok, err := localClient.DeleteCert(tc.Cert)

		if !ok {
			t.Errorf("For Cert %s, fail to delete", tc.name)
		}

		if err != tc.expectedError {
			t.Errorf("For Cert %s, expected error %v, but got %v", tc.name, tc.expectedError, err)
		}
	}
}
