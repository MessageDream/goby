package form

import (
	"github.com/go-macaron/binding"
	macaron "gopkg.in/macaron.v1"
)

type AppAddForm struct {
	Name     string `form:"name" binding:"Required;MaxSize(35)"`
	Platform string `form:"platform" binding:"Required"`
}

func (f *AppAddForm) Validate(ctx *macaron.Context, errs binding.Errors) binding.Errors {
	return validate(errs, ctx.Data, f)
}
