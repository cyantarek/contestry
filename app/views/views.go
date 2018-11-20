package views

import (
	"github.com/gin-contrib/multitemplate"
)

func LoadTemplates() multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	//auth app
	r.AddFromFiles("login", "app/views/templates/layout/base.html", "app/views/templates/auth/login.html")
	r.AddFromFiles("signup", "app/views/templates/layout/base.html", "app/views/templates/auth/registration.html")
	r.AddFromFiles("recovery", "app/views/templates/layout/base.html", "app/views/templates/auth/recovery.html")
	r.AddFromFiles("confirmation", "app/views/templates/layout/base.html", "app/views/templates/auth/confirmation.html")

	r.AddFromFiles("index", "app/views/templates/layout/base.html", "app/views/templates/common/index.html")

	//contest app
	r.AddFromFiles("contest-create", "app/views/templates/layout/base.html", "app/views/templates/contest/contest-create.html")
	r.AddFromFiles("contest-list", "app/views/templates/layout/base.html", "app/views/templates/contest/contest-list.html")
	r.AddFromFiles("contest-add-judges", "app/views/templates/layout/base.html", "app/views/templates/contest/contest-add-judges.html")
	r.AddFromFiles("contest-detail", "app/views/templates/layout/base.html", "app/views/templates/contest/contest-detail.html")

	r.AddFromFiles("question-create", "app/views/templates/layout/base.html", "app/views/templates/question/question-create.html")
	r.AddFromFiles("question-detail", "app/views/templates/layout/base.html", "app/views/templates/question/question-detail.html")

	//user app
	r.AddFromFiles("profile", "app/views/templates/layout/base.html", "app/views/templates/user/profile.html")
	r.AddFromFiles("judge-list", "app/views/templates/layout/base.html", "app/views/templates/user/judge-list.html")
	r.AddFromFiles("rank-list", "app/views/templates/layout/base.html", "app/views/templates/user/rank-list.html")

	//team app
	//r.AddFromFiles("team-list", "app/views/templates/layout/base.html", "app/views/templates/user/team-list.html")
	//r.AddFromFiles("team-create", "app/views/templates/layout/base.html", "app/views/templates/user/team-create.html")

	return r
}


//func Render(h *controllers.Handler, ctx *gin.Context, tplName string, ctxData gin.H) {
//	//session := sessions.Default(ctx)
//	currentUser := h.SessionStore.GetSessionData(ctx, "user")
//	if currentUser != nil {
//		if ctxData != nil {
//			ctx.HTML(200, tplName, gin.H{"current_user": currentUser, "data": ctxData})
//		} else {
//			ctx.HTML(200, tplName, gin.H{"current_user": currentUser})
//		}
//	} else {
//		ctx.HTML(200, "index", nil)
//	}
//
//}
