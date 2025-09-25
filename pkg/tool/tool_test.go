package tool

import "testing"

func TestConvertPathToArgsFormat(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"/api/{{user}}", "/api/{{.Args.user}}"},
		{"/api/{{user}}/profile/{{id}}", "/api/{{.Args.user}}/profile/{{.Args.id}}"},
		{"/static/path", "/static/path"},                     // 无占位符
		{"/api/{{ user }}/test", "/api/{{.Args.user}}/test"}, // 带空格
		{"{{param}}", "{{.Args.param}}"},
	}

	for _, test := range tests {
		result := ConvertPathToArgsFormatV2(test.input)
		if result != test.expected {
			t.Errorf("ConvertPathToArgsFormat(%q) = %q, expected %q",
				test.input, result, test.expected)
		}
	}
}

func TestConvertArgsToPathFormat(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "单个参数逆向转换",
			input:    "/api/{{.Args.user}}",
			expected: "/api/{{user}}",
		},
		{
			name:     "多个参数逆向转换",
			input:    "/api/{{.Args.user}}/profile/{{.Args.id}}",
			expected: "/api/{{user}}/profile/{{id}}",
		},
		{
			name:     "带空格的参数",
			input:    "/api/{{ .Args.user }}/test",
			expected: "/api/{{user}}/test",
		},
		{
			name:     "无参数路径保持不变",
			input:    "/static/path",
			expected: "/static/path",
		},
		{
			name:     "复杂多参数路径",
			input:    "/api/{{.Args.version}}/users/{{.Args.user_id}}/posts/{{.Args.post_id}}",
			expected: "/api/{{version}}/users/{{user_id}}/posts/{{post_id}}",
		},
		{
			name:     "类似格式不转换",
			input:    "/api/{{.Other.user}}/test",
			expected: "/api/{{.Other.user}}/test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 测试简单版本
			//result := ConvertArgsToPathFormat(tt.input)
			//if result != tt.expected {
			//	t.Errorf("ConvertArgsToPathFormat(%q) = %q, expected %q",
			//		tt.input, result, tt.expected)
			//}

			// 测试复杂版本
			resultV2 := ConvertArgsToPathFormatV2(tt.input)
			if resultV2 != tt.expected {
				t.Errorf("ConvertArgsToPathFormatV2(%q) = %q, expected %q",
					tt.input, resultV2, tt.expected)
			}
		})
	}
}
