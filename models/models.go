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

// Package models is for loading and updating documentation files.
package models

import (
	"encoding/json"
	"errors"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Unknwon/com"
	"github.com/Unknwon/goconfig"
	"github.com/Unknwon/macaron"
	"github.com/robfig/cron"
	"github.com/slene/blackfriday"

	"github.com/gogits/gogsweb/modules/log"
	"github.com/gogits/gogsweb/modules/setting"
)

var (
	docs = make(map[string]map[string]*DocRoot)
)

type oldDocNode struct {
	Sha  string
	Path string
	Type string
}

// docTree descriables a documentation file structure tree.
var docTree struct {
	Tree []oldDocNode
}

var productTree struct {
	Tree []oldDocNode
}

type docFile struct {
	Title string
	Data  []byte
}

var (
	docLock = &sync.RWMutex{}
	docMaps = map[string]map[string]*docFile{}
)

func GetDocByLocale(name, lang string) *DocRoot {
	return docs[name][lang]
}

func InitModels() {
	for _, name := range setting.Apps {
		parseDocs(name)
		initDocMap(name)
	}

	return

	c := cron.New()
	c.AddFunc("0 */5 * * * *", checkFileUpdates)
	c.Start()

	if needCheckUpdate() {
		checkFileUpdates()
		setting.Cfg.SetValue("app", "update_check_time", strconv.Itoa(int(time.Now().Unix())))
		goconfig.SaveConfigFile(setting.Cfg, setting.CFG_CUSTOM_PATH)
	}
}

func parseDocs(name string) {
	if docs[name] == nil {
		docs[name] = make(map[string]*DocRoot)
	}

	root, err := ParseDocs(path.Join("docs", name, "zh-CN"))
	if err != nil {
		log.Error("Fail to parse docs: %v", err)
	}

	if root != nil {
		docs[name]["zh-CN"] = root
	}

	root, err = ParseDocs(path.Join("docs", name, "en-US"))
	if err != nil {
		log.Error("Fail to parse docs: %v", err)
	}

	if root != nil {
		docs[name]["en-US"] = root
	}
}

func needCheckUpdate() bool {
	// Does not have record for check update.
	stamp, err := setting.Cfg.Int64("app", "update_check_time")
	if err != nil {
		return true
	}

	for _, name := range setting.Apps {
		if !com.IsFile("conf/docTree_" + name + ".json") {
			return true
		}
	}

	return time.Unix(stamp, 0).Add(5 * time.Minute).Before(time.Now())
}

func initDocMap(name string) {
	treeName := "conf/docTree_" + name + ".json"
	isConfExist := com.IsFile(treeName)
	if isConfExist {
		f, err := os.Open(treeName)
		if err != nil {
			log.Error("Fail to open '%s': %v", treeName, err)
			return
		}
		defer f.Close()

		d := json.NewDecoder(f)
		if err = d.Decode(&docTree); err != nil {
			log.Error("Fail to decode '%s': %v", treeName, err)
			return
		}
	} else {
		// Generate 'docTree'.
		docTree.Tree = append(docTree.Tree, oldDocNode{Path: ""})
	}

	docLock.Lock()
	defer docLock.Unlock()

	docMap := make(map[string]*docFile)

	for _, l := range setting.Langs {
		os.MkdirAll(path.Join("docs", name, l), os.ModePerm)
		for _, v := range docTree.Tree {
			var fullName string
			if isConfExist {
				fullName = v.Path
			} else {
				fullName = l + "/" + v.Path
			}

			docMap[fullName] = getFile(path.Join("docs", name, fullName))
		}
	}

	docMaps[name] = docMap
}

// loadFile returns []byte of file data by given path.
func loadFile(filePath string) ([]byte, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return []byte(""), errors.New("Fail to open file: " + err.Error())
	}

	fi, err := f.Stat()
	if err != nil {
		return []byte(""), errors.New("Fail to get file information: " + err.Error())
	}

	d := make([]byte, fi.Size())
	f.Read(d)
	return d, nil
}

func markdown(raw []byte) []byte {
	htmlFlags := 0
	htmlFlags |= blackfriday.HTML_USE_XHTML
	htmlFlags |= blackfriday.HTML_USE_SMARTYPANTS
	htmlFlags |= blackfriday.HTML_SMARTYPANTS_FRACTIONS
	htmlFlags |= blackfriday.HTML_SMARTYPANTS_LATEX_DASHES
	htmlFlags |= blackfriday.HTML_GITHUB_BLOCKCODE
	htmlFlags |= blackfriday.HTML_OMIT_CONTENTS
	htmlFlags |= blackfriday.HTML_COMPLETE_PAGE
	renderer := blackfriday.HtmlRenderer(htmlFlags, "", "")

	// set up the parser
	extensions := 0
	extensions |= blackfriday.EXTENSION_NO_INTRA_EMPHASIS
	extensions |= blackfriday.EXTENSION_TABLES
	extensions |= blackfriday.EXTENSION_FENCED_CODE
	extensions |= blackfriday.EXTENSION_AUTOLINK
	extensions |= blackfriday.EXTENSION_STRIKETHROUGH
	extensions |= blackfriday.EXTENSION_HARD_LINE_BREAK
	extensions |= blackfriday.EXTENSION_SPACE_HEADERS
	extensions |= blackfriday.EXTENSION_NO_EMPTY_LINE_BEFORE_BLOCK

	body := blackfriday.Markdown(raw, renderer, extensions)
	return body
}

func getFile(filePath string) *docFile {
	if strings.Contains(filePath, "images") {
		return nil
	}

	df := &docFile{}
	p, err := loadFile(filePath + ".md")
	if err != nil {
		log.Error("Fail to load MD file: %v", err)
		return nil
	}

	// Parse and render.
	s := string(p)
	i := strings.Index(s, "\n")
	if i > -1 {
		// Has title.
		df.Title = strings.TrimSpace(
			strings.Replace(s[:i+1], "#", "", -1))
		if len(s) >= i+2 {
			df.Data = []byte(strings.TrimSpace(s[i+2:]))
		}
	} else {
		df.Data = p
	}

	df.Data = markdown(df.Data)
	return df
}

// GetDoc returns 'docFile' by given name and language version.
func GetDoc(app, fullName, lang string) *docFile {
	filePath := path.Join("docs", app, lang, fullName)

	if macaron.Env == macaron.DEV {
		return getFile(filePath)
	}

	docLock.RLock()
	defer docLock.RUnlock()
	return docMaps[app][lang+"/"+fullName]
}

var checkTicker *time.Ticker

func checkTickerTimer(checkChan <-chan time.Time) {
	for {
		<-checkChan
		checkFileUpdates()
	}
}

type rawFile struct {
	name   string
	rawURL string
	data   []byte
}

func (rf *rawFile) Name() string {
	return rf.name
}

func (rf *rawFile) RawUrl() string {
	return rf.rawURL
}

func (rf *rawFile) Data() []byte {
	return rf.data
}

func (rf *rawFile) SetData(p []byte) {
	rf.data = p
}

func checkFileUpdates() {
	log.Debug("Checking file updates")

	type tree struct {
		ApiUrl, RawUrl, TreeName, Prefix string
	}

	// FIXME: make this configurable.
	var trees = []*tree{
		{
			ApiUrl:   "https://api.github.com/repos/gogits/docs/git/trees/master?recursive=1&" + setting.GithubCred,
			RawUrl:   "https://raw.github.com/gogits/docs/master/",
			TreeName: "conf/docTree_gogs.json",
			Prefix:   "docs/gogs/",
		},
		{
			ApiUrl:   "https://api.github.com/repos/macaron-contrib/docs/git/trees/master?recursive=1&" + setting.GithubCred,
			RawUrl:   "https://raw.github.com/macaron-contrib/docs/master/",
			TreeName: "conf/docTree_macaron.json",
			Prefix:   "docs/macaron/",
		},
	}

	for _, tree := range trees {
		var tmpTree struct {
			Tree []*oldDocNode
		}

		if err := com.HttpGetJSON(httpClient, tree.ApiUrl, &tmpTree); err != nil {
			log.Error("Fail to get trees: %v", err)
			return
		}

		var saveTree struct {
			Tree []*oldDocNode
		}
		saveTree.Tree = make([]*oldDocNode, 0, len(tmpTree.Tree))

		// Compare SHA.
		files := make([]com.RawFile, 0, len(tmpTree.Tree))
		for _, node := range tmpTree.Tree {
			// Skip non-md files and "README.md".
			if node.Type != "blob" || (!strings.HasSuffix(node.Path, ".md") &&
				!strings.Contains(node.Path, "images") &&
				!strings.HasSuffix(node.Path, ".json")) ||
				strings.HasPrefix(strings.ToLower(node.Path), "readme") {
				continue
			}

			name := strings.TrimSuffix(node.Path, ".md")

			if checkSHA(name, node.Sha, tree.Prefix) {
				log.Info("Need to update: %s", name)
				files = append(files, &rawFile{
					name:   name,
					rawURL: tree.RawUrl + node.Path,
				})
			}

			saveTree.Tree = append(saveTree.Tree, &oldDocNode{
				Path: name,
				Sha:  node.Sha,
			})
			// For save purpose, reset name.
			node.Path = name
		}

		// Fetch files.
		if err := com.FetchFiles(httpClient, files, nil); err != nil {
			log.Error("Fail to fetch files: %v", err)
			return
		}

		// Update data.
		for _, f := range files {
			os.MkdirAll(path.Join(tree.Prefix, path.Dir(f.Name())), os.ModePerm)
			suf := ".md"
			if strings.Contains(f.Name(), "images") ||
				strings.HasSuffix(f.Name(), ".json") {
				suf = ""
			}
			fw, err := os.Create(tree.Prefix + f.Name() + suf)
			if err != nil {
				log.Error("Fail to open file: %v", err)
				continue
			}

			_, err = fw.Write(f.Data())
			fw.Close()
			if err != nil {
				log.Error("Fail to write data: %v", err)
				continue
			}
		}

		// Save documentation information.
		f, err := os.Create(tree.TreeName)
		if err != nil {
			log.Error("Fail to save data: %v", err)
			return
		}

		e := json.NewEncoder(f)
		err = e.Encode(&saveTree)
		if err != nil {
			log.Error("Fail to encode data: %v", err)
			return
		}
		f.Close()
	}

	log.Debug("Finish check file updates")
	for _, name := range setting.Apps {
		parseDocs(name)
		initDocMap(name)
	}
}

// checkSHA returns true if the documentation file need to update.
func checkSHA(name, sha, prefix string) bool {
	var tree struct {
		Tree []oldDocNode
	}

	if strings.HasPrefix(prefix, "docs/") {
		tree = docTree
	}

	for _, v := range tree.Tree {
		if v.Path == name {
			// Found.
			if v.Sha != sha {
				// Need to update.
				return true
			}
			return false
		}
	}
	// Not found.
	return true
}
