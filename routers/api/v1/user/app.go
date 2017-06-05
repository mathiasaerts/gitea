// Copyright 2014 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package user

import (
	api "code.gitea.io/sdk/gitea"

	"code.gitea.io/gitea/models"
	"code.gitea.io/gitea/modules/context"
)

// ListAccessTokens list all the access tokens
func ListAccessTokens(ctx *context.APIContext) {
	// swagger:route GET /user/tokens userGetTokens
	//
	//     Produces:
	//     - application/json
	//
	//     Responses:
	//       200: AccessTokenList
	//       500: error

	tokens, err := models.ListAccessTokens(ctx.User.ID)
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
	// swagger:route POST /user/tokens userCreateToken
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

	t := &models.AccessToken{
		UID:  ctx.User.ID,
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
	// swagger:route DELETE /user/tokens/{id} userDeleteAccessToken
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

	if err := models.DeleteAccessTokenByID(ctx.ParamsInt64(":id"), ctx.User.ID); err != nil {
		if models.IsErrAccessTokenAccessDenied(err) {
			ctx.Error(403, "", "You do not have access to this token")
		} else {
			ctx.Error(500, "DeleteAccessToken", err)
		}
		return
	}

	ctx.Status(204)
}
