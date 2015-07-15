// Copyright 2014 Unknwon
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

func GogsHome(ctx *macaron.Context) {
	ctx.Data["IsPageHome"] = true
	ctx.HTML(200, "home_gogs", ctx.Data)
}

func About(ctx *macaron.Context, locale i18n.Locale) {
	ctx.Data["Link"] = "/about"
	df := models.GetDoc("gogs", "about", locale.Lang)
	ctx.Data["Data"] = string(df.Data)
	ctx.HTML(200, "page", ctx.Data)
}

func Team(ctx *macaron.Context, locale i18n.Locale) {
	ctx.Data["Link"] = "/team"
	df := models.GetDoc("gogs", "team", locale.Lang)
	ctx.Data["Data"] = string(df.Data)
	ctx.HTML(200, "page", ctx.Data)
}

func Donate(ctx *macaron.Context, locale i18n.Locale) {
	ctx.Data["Link"] = "/donate"
	df := models.GetDoc("gogs", "donate", locale.Lang)
	ctx.Data["Data"] = string(df.Data)
	ctx.HTML(200, "page", ctx.Data)
}

func docs(ctx *macaron.Context, locale i18n.Locale, name string) {
	docRoot := models.GetDocByLocale(name, locale.Lang)
	if docRoot == nil {
		docRoot = models.GetDocByLocale(name, "en-US")
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
	ctx.HTML(200, "document_"+name, ctx.Data)
}

func GogsDocs(ctx *macaron.Context, locale i18n.Locale) {
	docs(ctx, locale, "gogs")
}

func docsStatic(ctx *macaron.Context, name string) {
	if len(ctx.Params(":all")) > 0 {
		f, err := os.Open(path.Join("docs", name, "images", ctx.Params(":all")))
		if err != nil {
			ctx.JSON(500, map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		defer f.Close()

		_, err = io.Copy(ctx.RW(), f)
		if err != nil {
			ctx.JSON(500, map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		return
	}
	ctx.Error(404)
}

func GogsStatic(ctx *macaron.Context) {
	docsStatic(ctx, "gogs")
}
