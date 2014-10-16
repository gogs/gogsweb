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
	"github.com/Unknwon/goconfig"
	"github.com/Unknwon/macaron"
)

const (
	CFG_PATH        = "conf/app.ini"
	CFG_CUSTOM_PATH = "conf/custom.ini"
)

var (
	Cfg *goconfig.ConfigFile

	Apps         []string
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
	Cfg, err = goconfig.LoadConfigFile(CFG_PATH)
	if err != nil {
		panic(fmt.Errorf("fail to load config file '%s': %v", CFG_PATH, err))
	}
	if com.IsFile(CFG_CUSTOM_PATH) {
		if err = Cfg.AppendFiles(CFG_CUSTOM_PATH); err != nil {
			panic(fmt.Errorf("fail to load config file '%s': %v", CFG_CUSTOM_PATH, err))
		}
	}

	Apps = Cfg.MustValueArray("", "apps", ",")

	if Cfg.MustValue("app", "run_mode", "dev") == "prod" {
		macaron.Env = macaron.PROD
	}
	HttpPort = Cfg.MustInt("app", "http_port", 8091)
	Https = Cfg.MustBool("app", "https")
	HttpsCert = Cfg.MustValue("app", "https_cert")
	HttpsKey = Cfg.MustValue("app", "https_key")

	Langs = Cfg.MustValueArray("i18n", "langs", ",")
	Names = Cfg.MustValueArray("i18n", "names", ",")

	setGithubCredentials(Cfg.MustValue("github", "client_id"), Cfg.MustValue("github", "client_secret"))
}
