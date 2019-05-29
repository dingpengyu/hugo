// Copyright 2017 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cast

import (
	"html/template"
	"time"

	"github.com/gohugoio/hugo/common/hugo"
	"github.com/gohugoio/hugo/langs"
	"github.com/gohugoio/hugo/navigation"
	"github.com/gohugoio/hugo/resources/page"

	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/docshelper"
	"github.com/gohugoio/hugo/tpl/internal"
	"github.com/spf13/viper"
)

// This file provides documentation support and is randomly put into this package.
func init() {
	docsProvider := func() map[string]interface{} {
		docs := make(map[string]interface{})
		d := &deps.Deps{
			Cfg:                 viper.New(),
			Log:                 loggers.NewErrorLogger(),
			BuildStartListeners: &deps.Listeners{},
			Site:                newTestHugoSite(),
		}

		var namespaces internal.TemplateFuncsNamespaces

		for _, nsf := range internal.TemplateFuncsNamespaceRegistry {
			nf := nsf(d)
			namespaces = append(namespaces, nf)

		}

		docs["funcs"] = namespaces
		return docs
	}

	docshelper.AddDocProvider("tpl", docsProvider)
}

type testSite struct {
	h hugo.Info
	l *langs.Language
}

func (t testSite) Hugo() hugo.Info {
	return t.h
}

func (t testSite) ServerPort() int {
	return 1313
}

func (testSite) LastChange() (t time.Time) {
	return
}

func (t testSite) Title() string {
	return "foo"
}

func (t testSite) Sites() page.Sites {
	return nil
}

func (t testSite) IsServer() bool {
	return false
}

func (t testSite) Language() *langs.Language {
	return t.l
}

func (t testSite) Pages() page.Pages {
	return nil
}

func (t testSite) RegularPages() page.Pages {
	return nil
}

func (t testSite) Menus() navigation.Menus {
	return nil
}

func (t testSite) Taxonomies() interface{} {
	return nil
}

func (t testSite) BaseURL() template.URL {
	return ""
}

func (t testSite) Params() map[string]interface{} {
	return nil
}

func (t testSite) Data() map[string]interface{} {
	return nil
}

// newTestHugoSite creates a new minimal test site.
func newTestHugoSite() page.Site {
	return testSite{
		h: hugo.NewInfo(hugo.EnvironmentProduction),
		l: langs.NewLanguage("en", newTestConfig()),
	}
}

func newTestConfig() *viper.Viper {
	v := viper.New()
	v.Set("contentDir", "content")
	return v
}
