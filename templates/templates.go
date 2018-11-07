package templates

import "github.com/gin-contrib/multitemplate"

func LoadTemplates() multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	//auth app
	r.AddFromFiles("login", "templates/layout/base.html", "templates/auth/login.html")
	r.AddFromFiles("signup", "templates/layout/base.html", "templates/auth/registration.html")
	r.AddFromFiles("recovery", "templates/layout/base.html", "templates/auth/recovery.html")
	r.AddFromFiles("confirmation", "templates/layout/base.html", "templates/auth/confirmation.html")

	//contest app
	r.AddFromFiles("index", "templates/layout/base.html", "templates/common/index.html")
	r.AddFromFiles("contest-create", "templates/layout/base.html", "templates/contest/contest-create.html")
	r.AddFromFiles("contest-list", "templates/layout/base.html", "templates/contest/contest-list.html")

	//user app
	r.AddFromFiles("profile", "templates/layout/base.html", "templates/user/profile.html")
	r.AddFromFiles("judge-list", "templates/layout/base.html", "templates/user/judge-list.html")

	return r
}

