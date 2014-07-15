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

// An open source project for official documentation website of Gogs.
package main

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/Unknwon/macaron"
	"github.com/macaron-contrib/i18n"

	"github.com/gogits/gogsweb/models"
	"github.com/gogits/gogsweb/modules/base"
	"github.com/gogits/gogsweb/modules/log"
	"github.com/gogits/gogsweb/modules/setting"
	"github.com/gogits/gogsweb/routers"
)

const (
	APP_VER = "0.1.0.0715"
)

func main() {
	log.Info("Gogs Web %s", APP_VER)
	log.Info("Run Mode: %s", strings.Title(macaron.Env))

	m := macaron.Classic()

	// Middlewares.
	m.Use(macaron.Renderer(macaron.RenderOptions{
		Funcs: []template.FuncMap{map[string]interface{}{
			"dict":     base.Dict,
			"str2html": base.Str2html,
		}},
	}))
	m.Use(i18n.I18n(i18n.LocaleOptions{
		Langs:    setting.Langs,
		Names:    setting.Names,
		Redirect: true,
	}))

	// Routers.
	m.Get("/", routers.Home)
	m.Get("/docs", routers.Docs)
	m.Get("/docs/*all", routers.Docs)
	m.Get("/about", routers.About)
	m.Get("/team", routers.Team)
	m.Get("/donate", routers.Donate)

	models.InitModels()
	http.ListenAndServe("0.0.0.0:"+setting.HttpPort, m)
}
