// Copyright 2018 cg33.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package echo

import (
	"bytes"
	"errors"
	"github.com/chenhg5/go-admin/context"
	"github.com/chenhg5/go-admin/engine"
	"github.com/chenhg5/go-admin/modules/auth"
	"github.com/chenhg5/go-admin/modules/config"
	"github.com/chenhg5/go-admin/modules/logger"
	"github.com/chenhg5/go-admin/modules/menu"
	"github.com/chenhg5/go-admin/plugins"
	"github.com/chenhg5/go-admin/plugins/admin/modules/constant"
	"github.com/chenhg5/go-admin/template"
	"github.com/chenhg5/go-admin/template/types"
	"github.com/labstack/echo"
	template2 "html/template"
	"net/http"
	"strings"
)

type Echo struct {
}

func init() {
	engine.Register(new(Echo))
}

func (e *Echo) Use(router interface{}, plugin []plugins.Plugin) error {
	var (
		eng *echo.Echo
		ok  bool
	)
	if eng, ok = router.(*echo.Echo); !ok {
		return errors.New("wrong parameter")
	}

	for _, plug := range plugin {
		var plugCopy plugins.Plugin
		plugCopy = plug
		for _, req := range plug.GetRequest() {
			eng.Add(strings.ToUpper(req.Method), req.URL, func(c echo.Context) error {
				ctx := context.NewContext(c.Request())

				for _, key := range c.ParamNames() {
					if c.Request().URL.RawQuery == "" {
						c.Request().URL.RawQuery += strings.Replace(key, ":", "", -1) + "=" + c.Param(key)
					} else {
						c.Request().URL.RawQuery += "&" + strings.Replace(key, ":", "", -1) + "=" + c.Param(key)
					}
				}

				ctx.SetHandlers(plugCopy.GetHandler(c.Request().URL.Path, strings.ToLower(c.Request().Method))).Next()
				for key, head := range ctx.Response.Header {
					c.Response().Header().Set(key, head[0])
				}
				if ctx.Response.Body != nil {
					buf := new(bytes.Buffer)
					_, _ = buf.ReadFrom(ctx.Response.Body)
					_ = c.String(ctx.Response.StatusCode, buf.String())
				}
				return nil
			})
		}
	}

	return nil
}

func (e *Echo) Content(contextInterface interface{}, c types.GetPanel) {

	var (
		ctx echo.Context
		ok  bool
	)
	if ctx, ok = contextInterface.(echo.Context); !ok {
		panic("wrong parameter")
	}

	globalConfig := config.Get()

	sesKey, err := ctx.Cookie("go_admin_session")

	if err != nil || sesKey == nil {
		_ = ctx.Redirect(http.StatusFound, globalConfig.Url("/login"))
		return
	}

	userId, ok := auth.Driver.Load(sesKey.Value)["user_id"]

	if !ok {
		_ = ctx.Redirect(http.StatusFound, globalConfig.Url("/login"))
		return
	}

	user, ok := auth.GetCurUserById(int64(userId.(float64)))

	if !ok {
		_ = ctx.Redirect(http.StatusFound, globalConfig.Url("/login"))
		return
	}

	var panel types.Panel

	if !auth.CheckPermissions(user, ctx.Path(), ctx.Request().Method) {
		alert := template.Get(globalConfig.Theme).Alert().SetTitle(template2.HTML(`<i class="icon fa fa-warning"></i> Error!`)).
			SetTheme("warning").SetContent(template2.HTML("no permission")).GetContent()

		panel = types.Panel{
			Content:     alert,
			Description: "Error",
			Title:       "Error",
		}
	} else {
		panel = c()
	}

	tmpl, tmplName := template.Get(globalConfig.Theme).GetTemplate(ctx.Request().Header.Get(constant.PjaxHeader) == "true")

	ctx.Response().Header().Set("Content-Type", "text/html; charset=utf-8")

	buf := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(buf, tmplName, types.Page{
		User: user,
		Menu: *(menu.GetGlobalMenu(user).SetActiveClass(globalConfig.UrlRemovePrefix(ctx.Request().URL.String()))),
		System: types.SystemInfo{
			Version: "0.0.1",
		},
		Panel:       panel,
		UrlPrefix:   globalConfig.Prefix(),
		Title:       globalConfig.Title,
		Logo:        globalConfig.Logo,
		MiniLogo:    globalConfig.MiniLogo,
		ColorScheme: globalConfig.ColorScheme,
	})
	if err != nil {
		logger.Error("Echo Content", err)
	}
	_ = ctx.String(http.StatusOK, buf.String())
}
