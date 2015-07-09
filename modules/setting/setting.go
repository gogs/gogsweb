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

package setting

import (
	"fmt"

	"github.com/Unknwon/com"
	"github.com/Unknwon/macaron"
	"gopkg.in/ini.v1"
)

const (
	CFG_PATH        = "conf/app.ini"
	CFG_CUSTOM_PATH = "conf/custom.ini"
)

type App struct {
	Name     string
	RepoName string
}

var (
	Cfg *ini.File

	Apps         []App
	HttpPort     int
	Https        bool
	HttpsCert    string
	HttpsKey     string
	Langs, Names []string
	GithubCred   string
)

func setGithubCredentials(id, secret string) {
	GithubCred = "client_id=" + id + "&client_secret=" + secret
}

func init() {
	var err error
	Cfg, err = ini.Load(CFG_PATH)
	if err != nil {
		panic(fmt.Errorf("fail to load config file '%s': %v", CFG_PATH, err))
	}
	if com.IsFile(CFG_CUSTOM_PATH) {
		if err = Cfg.Append(CFG_CUSTOM_PATH); err != nil {
			panic(fmt.Errorf("fail to load config file '%s': %v", CFG_CUSTOM_PATH, err))
		}
	}

	appNames := Cfg.Section("").Key("apps").Strings(",")
	Apps = make([]App, len(appNames))
	for i, name := range appNames {
		Apps[i] = App{name, Cfg.Section(name).Key("repo_name").String()}
	}

	sec := Cfg.Section("app")
	if sec.Key("run_mode").MustString("dev") == "prod" {
		macaron.Env = macaron.PROD
	}
	HttpPort = sec.Key("http_port").MustInt(8091)
	Https = sec.Key("https").MustBool()
	HttpsCert = sec.Key("https_cert").String()
	HttpsKey = sec.Key("https_key").String()

	Langs = Cfg.Section("i18n").Key("langs").Strings(",")
	Names = Cfg.Section("i18n").Key("names").Strings(",")

	setGithubCredentials(Cfg.Section("github").Key("client_id").String(), Cfg.Section("github").Key("client_secret").String())
}
