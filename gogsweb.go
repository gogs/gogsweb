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

// An open source project for official documentation website of Gogs.
package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/Unknwon/macaron"
	"github.com/macaron-contrib/i18n"
	"github.com/macaron-contrib/switcher"

	"github.com/gogits/gogsweb/models"
	"github.com/gogits/gogsweb/modules/base"
	"github.com/gogits/gogsweb/modules/log"
	"github.com/gogits/gogsweb/modules/setting"
	"github.com/gogits/gogsweb/routers"
)

const APP_VER = "0.2.7.0202"

var funcMap = map[string]interface{}{
	"dict":     base.Dict,
	"str2html": base.Str2html,
}

func newGogsInstance() *macaron.Macaron {
	m := macaron.Classic()

	// Middlewares.
	m.Use(macaron.Renderer(macaron.RenderOptions{
		Funcs: []template.FuncMap{funcMap},
	}))
	m.Use(i18n.I18n(i18n.Options{
		Langs:    setting.Langs,
		Names:    setting.Names,
		Redirect: true,
	}))

	// Routers.
	m.Get("/", routers.GogsHome)
	m.Get("/docs", routers.GogsDocs)
	m.Get("/docs/images/:all", routers.GogsStatic)
	m.Get("/docs/*", routers.GogsDocs)
	m.Get("/about", routers.About)
	m.Get("/team", routers.Team)
	m.Get("/donate", routers.Donate)

	return m
}

func newMacaronInstance() *macaron.Macaron {
	m := macaron.Classic()

	// Middlewares.
	m.Use(macaron.Renderer(macaron.RenderOptions{
		Funcs: []template.FuncMap{funcMap},
	}))
	m.Use(i18n.I18n(i18n.Options{
		Langs:    setting.Langs,
		Names:    setting.Names,
		Redirect: true,
	}))

	// Routers.
	m.Get("/", routers.MacaronDocs)
	m.Get("/docs", routers.MacaronDocs)
	m.Get("/docs/images/:all", routers.MacaronStatic)
	m.Get("/docs/*", routers.MacaronDocs)

	return m
}

func main() {
	log.Info("Gogs Web %s", APP_VER)
	log.Info("Run Mode: %s", strings.Title(macaron.Env))

	models.InitModels()

	m1 := newGogsInstance()
	m2 := newMacaronInstance()
	hs := switcher.NewHostSwitcher()
	hs.Set("gogs.io", m1)
	hs.Set("macaron.gogs.io", m2)

	var err error

	// In dev mode, listen two ports just for convenience.
	if macaron.Env == macaron.DEV {
		// Gogs.
		listenAddr := fmt.Sprintf("0.0.0.0:%d", setting.HttpPort)
		log.Info("Listen: http://%s", listenAddr)
		go http.ListenAndServe(listenAddr, m1)

		// Macaron.
		listenAddr = fmt.Sprintf("0.0.0.0:%d", setting.HttpPort+1)
		log.Info("Listen: http://%s", listenAddr)
		err = http.ListenAndServe(listenAddr, m2)
	} else {
		schema := "http"
		if setting.Https {
			schema = "https"
		}
		listenAddr := fmt.Sprintf("0.0.0.0:%d", setting.HttpPort)
		log.Info("Listen: %v://%s", schema, listenAddr)
		if setting.Https {
			err = http.ListenAndServeTLS(listenAddr, setting.HttpsCert, setting.HttpsKey, hs)
		} else {
			err = http.ListenAndServe(listenAddr, hs)
		}
	}
	if err != nil {
		log.Fatal("Fail to start server: %v", err)
	}
}
