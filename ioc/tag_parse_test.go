package ioc_test

import (
	"testing"

	"github.com/infraboard/mcube/v2/ioc"
)

func TestParseInjectTag_Valid(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *ioc.InjectTag
	}{
		{
			name:  "empty tag",
			input: "",
			expected: &ioc.InjectTag{
				Autowire:  false,
				Namespace: "default",
				Version:   "1.0.0",
				Name:      "",
			},
		},
		{
			name:  "autowire only",
			input: "autowire=true",
			expected: &ioc.InjectTag{
				Autowire:  true,
				Namespace: "default",
				Version:   "1.0.0",
				Name:      "",
			},
		},
		{
			name:  "autowire false",
			input: "autowire=false",
			expected: &ioc.InjectTag{
				Autowire:  false,
				Namespace: "default",
				Version:   "1.0.0",
				Name:      "",
			},
		},
		{
			name:  "autowire without value defaults to true",
			input: "autowire",
			expected: &ioc.InjectTag{
				Autowire:  true,
				Namespace: "default",
				Version:   "1.0.0",
				Name:      "",
			},
		},
		{
			name:  "all fields",
			input: "autowire=true;namespace=api;name=userService;version=v2",
			expected: &ioc.InjectTag{
				Autowire:  true,
				Namespace: "api",
				Name:      "userService",
				Version:   "v2",
			},
		},
		{
			name:  "namespace only",
			input: "namespace=configs",
			expected: &ioc.InjectTag{
				Autowire:  false,
				Namespace: "configs",
				Version:   "1.0.0",
				Name:      "",
			},
		},
		{
			name:  "namespace empty defaults to default",
			input: "namespace=",
			expected: &ioc.InjectTag{
				Autowire:  false,
				Namespace: "default",
				Version:   "1.0.0",
				Name:      "",
			},
		},
		{
			name:  "version only",
			input: "version=v3",
			expected: &ioc.InjectTag{
				Autowire:  false,
				Namespace: "default",
				Version:   "v3",
				Name:      "",
			},
		},
		{
			name:  "version empty defaults to v1",
			input: "version=",
			expected: &ioc.InjectTag{
				Autowire:  false,
				Namespace: "default",
				Version:   "1.0.0",
				Name:      "",
			},
		},
		{
			name:  "whitespace handling",
			input: "  autowire = true ; namespace = api ; name = service  ",
			expected: &ioc.InjectTag{
				Autowire:  true,
				Namespace: "api",
				Name:      "service",
				Version:   "1.0.0",
			},
		},
		{
			name:  "empty segments ignored",
			input: "autowire=true;;namespace=api",
			expected: &ioc.InjectTag{
				Autowire:  true,
				Namespace: "api",
				Version:   "1.0.0",
				Name:      "",
			},
		},
		{
			name:  "name with special characters",
			input: "name=user.service:v1",
			expected: &ioc.InjectTag{
				Autowire:  false,
				Namespace: "default",
				Name:      "user.service:v1",
				Version:   "1.0.0",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ioc.ParseInjectTagWithError(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.Autowire != tt.expected.Autowire {
				t.Errorf("Autowire: got %v, want %v", result.Autowire, tt.expected.Autowire)
			}
			if result.Namespace != tt.expected.Namespace {
				t.Errorf("Namespace: got %q, want %q", result.Namespace, tt.expected.Namespace)
			}
			if result.Name != tt.expected.Name {
				t.Errorf("Name: got %q, want %q", result.Name, tt.expected.Name)
			}
			if result.Version != tt.expected.Version {
				t.Errorf("Version: got %q, want %q", result.Version, tt.expected.Version)
			}
		})
	}
}

func TestParseInjectTag_Invalid(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectedErr string
	}{
		{
			name:        "invalid autowire value",
			input:       "autowire=yes",
			expectedErr: "invalid autowire value",
		},
		{
			name:        "invalid autowire value number",
			input:       "autowire=1",
			expectedErr: "invalid autowire value",
		},
		{
			name:        "unknown key",
			input:       "autowire=true;unknown=value",
			expectedErr: "unknown tag key",
		},
		{
			name:        "empty key",
			input:       "=value",
			expectedErr: "empty key",
		},
		{
			name:        "empty name value",
			input:       "name=",
			expectedErr: "name value cannot be empty",
		},
		{
			name:        "multiple unknown keys",
			input:       "autowire=true;foo=bar;namespace=api",
			expectedErr: "unknown tag key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ioc.ParseInjectTagWithError(tt.input)

			if err == nil {
				t.Fatalf("expected error containing %q, got nil (result: %+v)", tt.expectedErr, result)
			}

			if !containsString(err.Error(), tt.expectedErr) {
				t.Errorf("error message %q does not contain %q", err.Error(), tt.expectedErr)
			}
		})
	}
}

func TestParseInjectTag_BackwardCompatibility(t *testing.T) {
	// 测试旧版本函数(不返回error)的行为
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "valid input",
			input: "autowire=true;namespace=api",
		},
		{
			name:  "invalid input returns default",
			input: "autowire=invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ParseInjectTag不返回error，无效输入也返回默认值（不会panic）
			result := ioc.ParseInjectTag(tt.input)

			// 应该总是返回非nil结果
			if result == nil {
				t.Fatal("expected non-nil result")
			}

			// 对于有效输入，检查结果
			if tt.input == "autowire=true;namespace=api" {
				if !result.Autowire {
					t.Errorf("Autowire: got %v, want true", result.Autowire)
				}
				if result.Namespace != "api" {
					t.Errorf("Namespace: got %q, want %q", result.Namespace, "api")
				}
			}

			// 对于无效输入，应该返回默认值（向后兼容行为）
			if tt.input == "autowire=invalid" {
				// 由于ParseInjectTag内部调用ParseInjectTagWithError并忽略错误
				// 当发生错误时，会返回部分解析的tag或默认tag
				// 这里我们只验证不会panic即可
				t.Logf("Invalid input handled gracefully, got: %+v", result)
			}
		})
	}
}

func TestParseInjectTag_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantErr  bool
		expected *ioc.InjectTag
	}{
		{
			name:    "only semicolons",
			input:   ";;;",
			wantErr: false,
			expected: &ioc.InjectTag{
				Autowire:  false,
				Namespace: "default",
				Version:   "1.0.0",
				Name:      "",
			},
		},
		{
			name:    "value with equals sign",
			input:   "name=service=v1",
			wantErr: false,
			expected: &ioc.InjectTag{
				Autowire:  false,
				Namespace: "default",
				Name:      "service=v1",
				Version:   "1.0.0",
			},
		},
		{
			name:    "trailing semicolon",
			input:   "autowire=true;",
			wantErr: false,
			expected: &ioc.InjectTag{
				Autowire:  true,
				Namespace: "default",
				Version:   "1.0.0",
				Name:      "",
			},
		},
		{
			name:    "leading semicolon",
			input:   ";autowire=true",
			wantErr: false,
			expected: &ioc.InjectTag{
				Autowire:  true,
				Namespace: "default",
				Version:   "1.0.0",
				Name:      "",
			},
		},
		{
			name:    "only spaces",
			input:   "   ",
			wantErr: false,
			expected: &ioc.InjectTag{
				Autowire:  false,
				Namespace: "default",
				Version:   "1.0.0",
				Name:      "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ioc.ParseInjectTagWithError(tt.input)

			if tt.wantErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !tt.wantErr && tt.expected != nil {
				if result.Autowire != tt.expected.Autowire {
					t.Errorf("Autowire: got %v, want %v", result.Autowire, tt.expected.Autowire)
				}
				if result.Namespace != tt.expected.Namespace {
					t.Errorf("Namespace: got %q, want %q", result.Namespace, tt.expected.Namespace)
				}
				if result.Name != tt.expected.Name {
					t.Errorf("Name: got %q, want %q", result.Name, tt.expected.Name)
				}
				if result.Version != tt.expected.Version {
					t.Errorf("Version: got %q, want %q", result.Version, tt.expected.Version)
				}
			}
		})
	}
}

func TestNewInjectTag_Defaults(t *testing.T) {
	tag := ioc.NewInjectTag()

	if tag.Autowire != false {
		t.Errorf("default Autowire: got %v, want false", tag.Autowire)
	}
	if tag.Namespace != "default" {
		t.Errorf("default Namespace: got %q, want %q", tag.Namespace, "default")
	}
	if tag.Version != "1.0.0" {
		t.Errorf("default Version: got %q, want %q", tag.Version, "1.0.0")
	}
	if tag.Name != "" {
		t.Errorf("default Name: got %q, want empty", tag.Name)
	}
}

// Helper function
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && stringContains(s, substr)))
}

func stringContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
