package controllers

import (
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"../db"
	"../session-mgr"
	"crypto/rsa"
)

//Data Access Object Pattern (Store) Handler <-> Store <-> Database
//Program to Interface (Store), No to Implementation Principle
//Injecting dependencies to Handler, so that we don't need to have global variables
type Handler struct {
	Store        db.Store
	SessionStore session_mgr.SessionStore
	Type         string
	SignKey      *rsa.PrivateKey
	VerifyKey    *rsa.PublicKey
	ContextData  interface{}
}

/*
package controller has a struct Handler that wraps third party dependencies for handlers.

This way, dependencies can be injected to the controller function/methods rather using public
global variable.

Also the dependencies are not direct, rather comes through interfaces. So we can change the dependencies
any time without breaking the code.

 */
