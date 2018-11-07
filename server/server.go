package server

import (
	"net/http"
	"github.com/prometheus/common/log"
	"encoding/gob"

	"../user"
)

func Run(addr string, router http.Handler) {
	gob.Register(user.User{})
	log.Fatal(http.ListenAndServe(":" + addr, router))
}

