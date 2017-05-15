package accessKey

import (
	"github.com/MessageDream/goby/core/tokenService"
	"github.com/MessageDream/goby/module/context"
)

func Add(ctx *context.APIContext) {
	createdBy := ctx.Query("createdBy")
	friendlyName := ctx.Query("friendlyName")
	ttl := ctx.QueryInt64("ttl")
	description := ctx.Query("description")

	token, err := tokenService.Create(ctx.User, createdBy, friendlyName, description, ttl)

	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(200, map[string]interface{}{
		"accessKey": token.Convert(),
	})
}

func Get(ctx *context.APIContext) {
	tokens, err := ctx.User.Tokens()
	if err != nil {
		ctx.Error(err)
		return
	}
	outputTokens := make([]interface{}, 0, len(tokens))
	for _, token := range tokens {
		outputTokens = append(outputTokens, token.Convert())
	}

	ctx.JSON(200, map[string]interface{}{
		"accessKeys": outputTokens,
	})
}

func Delete(ctx *context.APIContext) {
	name := ctx.Params("name")
	uid := ctx.User.ID

	if err := tokenService.DeleteByUIDAndName(uid, name); err != nil {
		ctx.Error(err)
	}

	ctx.JSON(200, map[string]interface{}{
		"friendlyName": name,
	})
}

func Patch(ctx *context.APIContext) {
	name := ctx.Params("name")
	friendlyName := ctx.Query("friendlyName")
	ttl := ctx.QueryInt64("ttl")
	uid := ctx.User.ID

	token, err := tokenService.Update(uid, name, friendlyName, ttl)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, map[string]interface{}{
		"accessKey": token.Convert(),
	})
}
