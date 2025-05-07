# z-validation

## 特性

- ✅ 支持 **结构体**（Struct）、**Map**（map[string]any）、**单个参数**（如字符串、数字）等多种类型校验；
- ✅ 链式调用 API（如 `R("name").Required().MinLength(6)`），代码简洁易读；
- ✅ 自定义错误消息：每个校验规则支持传入自定义提示（如 `Required("用户名不能为空")`）；
- ✅ 内置常用规则：`Required`（必填）、`MinLength`（最小长度）、`MaxLength`（最大长度）、`EqualField`（字段相等）等；
- ✅ 错误信息结构化：返回包含字段名和具体错误的列表，方便前端或日志解析。

## 安装

```bash
go get github.com/zR-Zr/z-validation
```

#### 4. **快速开始**

提供 2-3 个典型场景的示例代码（结合你的 `example/main.go` 中的用例）：

**示例 1：校验结构体（User 类型）**

### 校验结构体

```go
package main

import (
	"fmt"
	"github.com/zR-Zr/z-validation/validation"
)

type User struct {
	Name       *string `json:"name"`
	Password   *string `json:"password"`
	RePassword *string `json:"rePassword"`
}

func main() {
	user := User{
		Password:   &[]string{"123456789"}[0],
		RePassword: &[]string{"12345678"}[0],
	}

	rules := validation.Rules{
		validation.R("Name").ConvertField("name").Required("用户名不能为空"),
		validation.R("Password").ConvertField("password").Required("密码不能为空").MinLength(6, "密码不能小于6位"),
		validation.R("RePassword").ConvertField("repassword").Required("确认密码不能为空").EqualField("Password", "两次密码不一致"),
	}

	errors := validation.ValidateStruct(user, &rules)
	if errors != nil {
		fmt.Println("校验失败:", errors)
	} else {
		fmt.Println("校验通过")
	}
}
```

**示例 2：校验 Map**

### 校验 Map

```go
func main() {
	values := map[string]any{
		"name":       "zrcoder",
		"password":   "123456789",
		"rePassword": "12345678",
	}

	err := validation.ValidateMap(values, &validation.Rules{
		validation.R("name").ConvertField("name").Required("用户名不能为空"),
		validation.R("password").ConvertField("password").Required("密码不能为空").MinLength(6, "密码不能小于6位"),
		validation.R("rePassword").ConvertField("repassword").Required("确认密码不能为空").EqualField("password", "两次密码不一致"),
	})

	if err != nil {
		fmt.Println("校验失败:", err)
	} else {
		fmt.Println("校验通过")
	}
}
```

#### 5. **API 文档（简要）**

列出核心类型和方法的作用（无需详细参数，保持简洁）：

## API 概览

| 类型/方法             | 说明                                                                                   |
| --------------------- | -------------------------------------------------------------------------------------- |
| `Rule`                | 单个校验规则的封装，支持链式调用添加 `Required`、`MinLength` 等规则                    |
| `Rules`               | 校验规则的集合，提供 `ValidateStruct`（结构体校验）、`ValidateValue`（通用值校验）方法 |
| `R(field string)`     | 初始化一个针对结构体字段的规则构建器                                                   |
| `RV(name string)`     | 初始化一个针对非结构体（如单个参数、Map）的规则构建器                                  |
| `ConvertField(field)` | 自定义返回的字段名（如将结构体字段 `Name` 映射为 `"username"`）                        |

## 许可证
