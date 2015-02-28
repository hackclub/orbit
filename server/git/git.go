// Adapted from https://code.csdn.net/flycutter/git-http-backend
package git

import "net/http"

type Service struct {
	Method  string
	Handler func(handlerReq)
	RPC     string
}

type Config struct {
	ProjectRoot string
	GitBinPath  string
	UploadPack  bool
	ReceivePack bool
}

type handlerReq struct {
	w    http.ResponseWriter
	r    *http.Request
	RPC  string
	Dir  string
	File string
}
