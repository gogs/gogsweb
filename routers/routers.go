// Copyright 2014 Unknown
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package routers

import (
	"io"
	"os"
	"path"
	"strings"

	"github.com/Unknwon/macaron"
	"github.com/macaron-contrib/i18n"

	"github.com/gogits/gogsweb/models"
)

func Home(ctx *macaron.Context) {
	ctx.Data["IsPageHome"] = true
	ctx.HTML(200, "home", ctx.Data)
}

func Docs(ctx *macaron.Context, locale i18n.Locale) {
	docRoot := models.GetDocByLocale(locale.Lang)
	if docRoot == nil {
		ctx.Error(404)
		return
	}

	link := strings.TrimPrefix(ctx.Params("*"), "/")
	link = strings.TrimSuffix(link, ".html")
	link = strings.TrimSuffix(link, ".md")
	ctx.Data["Link"] = "/docs/" + link

	var doc *models.DocNode
	if len(link) == 0 {
		ctx.Redirect("/docs/intro/")
		return
	}

	doc, _ = docRoot.GetNodeByLink(link)
	if doc == nil {
		doc, _ = docRoot.GetNodeByLink(link + "/")
	}
	if doc == nil {
		ctx.Error(404)
		return
	}

	ctx.Data["DocRoot"] = docRoot
	ctx.Data["Doc"] = doc
	ctx.Data["Title"] = doc.Name
	ctx.Data["Data"] = doc.GetContent()
	ctx.HTML(200, "document", ctx.Data)
}

func DocsStatic(ctx *macaron.Context) {
	if len(ctx.Params(":all")) > 0 {
		f, err := os.Open(path.Join("docs/images", ctx.Params(":all")))
		if err != nil {
			ctx.JSON(500, map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		defer f.Close()

		_, err = io.Copy(ctx.ResponseWriter, f)
		if err != nil {
			ctx.JSON(500, map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
	}
	ctx.Error(404)
}

func About(ctx *macaron.Context, locale i18n.Locale) {
	ctx.Data["Link"] = "/about"
	df := models.GetDoc("about", locale.Lang)
	ctx.Data["Data"] = string(df.Data)
	ctx.HTML(200, "page", ctx.Data)
}

func Team(ctx *macaron.Context, locale i18n.Locale) {
	ctx.Data["Link"] = "/team"
	df := models.GetDoc("team", locale.Lang)
	ctx.Data["Data"] = string(df.Data)
	ctx.HTML(200, "page", ctx.Data)
}

func Donate(ctx *macaron.Context, locale i18n.Locale) {
	ctx.Data["Link"] = "/donate"
	df := models.GetDoc("donate", locale.Lang)
	ctx.Data["Data"] = string(df.Data)
	ctx.HTML(200, "page", ctx.Data)
}
