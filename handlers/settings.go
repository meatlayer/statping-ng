package handlers

import (
	"fmt"
	"github.com/hunterlong/statup/core"
	"github.com/hunterlong/statup/notifiers"
	"github.com/hunterlong/statup/utils"
	"net/http"
)

func PluginsHandler(w http.ResponseWriter, r *http.Request) {
	if !IsAuthenticated(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	//CoreApp.FetchPluginRepo()

	//var pluginFields []PluginSelect
	//
	//for _, p := range allPlugins {
	//	fields := structs.Map(p.GetInfo())
	//
	//	pluginFields = append(pluginFields, PluginSelect{p.GetInfo().Name, p.GetForm(), fields})
	//}

	//CoreApp.PluginFields = pluginFields

	ExecuteResponse(w, r, "settings.html", core.CoreApp)
}

func SaveSettingsHandler(w http.ResponseWriter, r *http.Request) {
	if !IsAuthenticated(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	r.ParseForm()
	name := r.PostForm.Get("project")
	if name != "" {
		core.CoreApp.Name = name
	}
	description := r.PostForm.Get("description")
	if description != core.CoreApp.Description {
		core.CoreApp.Description = description
	}
	style := r.PostForm.Get("style")
	if style != core.CoreApp.Style {
		core.CoreApp.Style = style
	}
	footer := r.PostForm.Get("footer")
	if footer != core.CoreApp.Footer {
		core.CoreApp.Footer = footer
	}
	domain := r.PostForm.Get("domain")
	if domain != core.CoreApp.Domain {
		core.CoreApp.Domain = domain
	}
	core.CoreApp.UseCdn = (r.PostForm.Get("enable_cdn") == "on")
	core.CoreApp.Update()
	core.OnSettingsSaved(core.CoreApp)
	http.Redirect(w, r, "/settings", http.StatusSeeOther)
}

func SaveSASSHandler(w http.ResponseWriter, r *http.Request) {
	if !IsAuthenticated(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	r.ParseForm()
	theme := r.PostForm.Get("theme")
	variables := r.PostForm.Get("variables")
	core.SaveAsset(theme, "scss/base.scss")
	core.SaveAsset(variables, "scss/variables.scss")
	core.CompileSASS()
	http.Redirect(w, r, "/settings", http.StatusSeeOther)
}

func SaveAssetsHandler(w http.ResponseWriter, r *http.Request) {
	if !IsAuthenticated(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	core.CreateAllAssets()
	core.UsingAssets = true
	http.Redirect(w, r, "/settings", http.StatusSeeOther)
}

func SaveNotificationHandler(w http.ResponseWriter, r *http.Request) {
	if !IsAuthenticated(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	r.ParseForm()

	notifierId := r.PostForm.Get("id")
	enabled := r.PostForm.Get("enable")

	host := r.PostForm.Get("host")
	port := int(utils.StringInt(r.PostForm.Get("port")))
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")
	var1 := r.PostForm.Get("var1")
	var2 := r.PostForm.Get("var2")
	apiKey := r.PostForm.Get("api_key")
	apiSecret := r.PostForm.Get("api_secret")
	limits := int64(utils.StringInt(r.PostForm.Get("limits")))
	notifer := notifiers.Select(utils.StringInt(notifierId))
	if host != "" {
		notifer.Host = host
	}
	if port != 0 {
		notifer.Port = port
	}
	if username != "" {
		notifer.Username = username
	}
	if password != "" && password != "##########" {
		notifer.Password = password
	}
	if var1 != "" {
		notifer.Var1 = var1
	}
	if var2 != "" {
		notifer.Var2 = var2
	}
	if apiKey != "" {
		notifer.ApiKey = apiKey
	}
	if apiSecret != "" {
		notifer.ApiSecret = apiSecret
	}
	if limits != 0 {
		notifer.Limits = limits
	}
	if enabled == "on" {
		notifer.Enabled = true
	} else {
		notifer.Enabled = false
	}
	notifer, err := notifer.Update()
	if err != nil {
		utils.Log(3, err)
	}
	msg := fmt.Sprintf("%v - %v - %v", notifierId, notifer, enabled)
	w.Write([]byte(msg))
}
