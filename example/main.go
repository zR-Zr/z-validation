package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/zR-Zr/z-validation/validation"
	"github.com/zR-Zr/z-validation/zvr"
)

type User struct {
	Name       *string `json:"name"`
	Password   *string `json:"password"`
	RePassword *string `json:"rePassword"`
	Email      string  `json:"email"`
}

func main() {

	eng := gin.Default()
	eng.GET("/", func(c *gin.Context) {
		err := V()
		if err != nil {
			c.JSON(200, gin.H{
				"msg":     "参数校验错误",
				"details": err.(*validation.ValidationError).Details(),
			})
		} else {
			c.JSON(200, gin.H{
				"msg": "ok",
			})
		}
	})

	eng.GET("/find", func(ctx *gin.Context) {
		name := ctx.Query("name")
		err := validation.ValidateValue(name, validation.RV("username").Required("用户名不能为空").MinLength(6, "用户名不能小于6位").MaxLength(16, "用户名不能大于16位"))
		if err != nil {
			ctx.JSON(200, gin.H{
				"msg":     "参数校验错误",
				"details": err.(*validation.ValidationError).Details(),
			})
		} else {
			ctx.JSON(200, gin.H{
				"msg": "ok",
			})

		}
	})

	eng.GET("/map", func(ctx *gin.Context) {
		values := map[string]any{
			"name":       "zrcoder",
			"password":   "123456789",
			"rePassword": "12345678",
		}
		err := validation.ValidateMap(values, &validation.Rules{
			validation.R("name").ConvertField("name").Required("用户名不能为空"),
			validation.R("password").ConvertField("password").Required("密码不能为空").MinLength(6, "密码不能小于6位").MaxLength(16, "密码不能大于16位"),
			validation.R("rePassword").ConvertField("repassword").Required("确认密码不能为空").EqualField("password", "两次密码不一致"),
		})
		if err != nil {
			ctx.JSON(200, gin.H{
				"msg":     "参数校验错误",
				"details": err.(*validation.ValidationError).Details(),
			})
		} else {
			ctx.JSON(200, gin.H{
				"msg": "ok",
			})

		}
	})

	eng.POST("/email", func(ctx *gin.Context) {
		var req struct {
			Email string `json:"email"`
		}

		var rules = zvr.Rs().Add("Email", func(r *validation.Rule) {
			r.ConvertField("email").Required("邮箱不能为空").Email("邮箱格式不正确")
		})()
		err := ctx.BindJSON(&req)
		if err != nil {
			ctx.JSON(200, gin.H{
				"msg":     "参数格式错误",
				"details": err.Error(),
			})
			return
		}

		verr := validation.ValidateStruct(req, &rules)
		if verr != nil {
			ctx.JSON(200, gin.H{
				"msg":     "参数校验错误",
				"details": verr.(*validation.ValidationError).Details(),
			})
			return
		}

		ctx.JSON(200, gin.H{
			"msg":  "ok",
			"data": req,
		})
	})

	eng.Run(":8081")

}
func V() error {
	user := User{}
	// user.Name = (&[]string{"zrcoder"}[0])
	user.Password = (&[]string{"123456789"}[0])
	// user.RePassword = (&[]string{"12345678"}[0])

	rules := validation.Rules{
		validation.R("Name").ConvertField("name").Required(), //.ConvertField("name").Required("用户名不能为空"),
		validation.R("Password").ConvertField("password").Required("密码不能为空").MinLength(6, "密码不能小于6位").MaxLength(16, "密码不能大于16位"),
		validation.R("RePassword").ConvertField("repassword").Required("确认密码不能为空").EqualField("Password", "两次密码不一致"),
	}

	errors := validation.ValidateStruct(user, &rules)
	fmt.Println(errors)

	if errors != nil {
		return errors
	}

	return nil
}
