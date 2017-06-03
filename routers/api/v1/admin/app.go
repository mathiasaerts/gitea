// Copyright 2014 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package admin

import (
	api "code.gitea.io/sdk/gitea"

	"code.gitea.io/gitea/models"
	"code.gitea.io/gitea/modules/context"
	"code.gitea.io/gitea/routers/user"
)

// ListAccessTokens list all the access tokens
func ListAccessTokens(ctx *context.APIContext) {
	// swagger:route GET /admin/users/{username}/tokens userGetTokens
	//
	//     Produces:
	//     - application/json
	//
	//     Responses:
	//       200: AccessTokenList
	//       500: error

	u, err := models.GetUserByName(ctx.Params(":username"))
	if err != nil {
		if models.IsErrUserNotExist(err) {
			ctx.Status(404)
		} else {
			ctx.Error(500, "GetUserByName", err)
		}
		return
	}

	tokens, err := models.ListAccessTokens(u.ID)
	if err != nil {
		ctx.Error(500, "ListAccessTokens", err)
		return
	}

	apiTokens := make([]*api.AccessToken, len(tokens))
	for i := range tokens {
		apiTokens[i] = &api.AccessToken{
			ID: tokens[i].ID,
			Name: tokens[i].Name,
			Sha1: tokens[i].Sha1,
		}
	}
	ctx.JSON(200, &apiTokens)
}

// CreateAccessToken create access tokens
func CreateAccessToken(ctx *context.APIContext, form api.CreateAccessTokenOption) {
	// swagger:route POST /admin/users/{username}/tokens userCreateToken
	//
	//     Consumes:
	//     - application/json
	//
	//     Produces:
	//     - application/json
	//
	//     Responses:
	//       200: AccessToken
	//       500: error

	u, err := models.GetUserByName(ctx.Params(":username"))
	if err != nil {
		if models.IsErrUserNotExist(err) {
			ctx.Status(404)
		} else {
			ctx.Error(500, "GetUserByName", err)
		}
		return
	}

	t := &models.AccessToken{
		UID:  u.ID,
		Name: form.Name,
	}
	if err := models.NewAccessToken(t); err != nil {
		ctx.Error(500, "NewAccessToken", err)
		return
	}
	ctx.JSON(201, &api.AccessToken{
		ID: t.ID,
		Name: t.Name,
		Sha1: t.Sha1,
	})
}

//DeleteAccessToken remove access tokens
func DeleteAccessToken(ctx *context.APIContext) {
	// swagger:route DELETE /admin/users/{username}/tokens/{id} userDeleteAccessToken
	//
	//     Consumes:
	//     - application/json
	//
	//     Produces:
	//     - application/json
	//
	//     Responses:
	//       204: empty
	//       403: forbidden
	//       500: error

	u, err := models.GetUserByName(ctx.Params(":username"))
	if err != nil {
		if models.IsErrUserNotExist(err) {
			ctx.Status(404)
		} else {
			ctx.Error(500, "GetUserByName", err)
		}
		return
	}

	if err := models.DeleteAccessTokenByID(ctx.ParamsInt64(":id"), u.ID); err != nil {
		if models.IsErrAccessTokenAccessDenied(err) {
			ctx.Error(403, "", "You do not have access to this token")
		} else {
			ctx.Error(500, "DeleteAccessToken", err)
		}
		return
	}

	ctx.Status(204)
}
