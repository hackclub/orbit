package router

import "github.com/gorilla/mux"

func Git() *mux.Router {
	m := mux.NewRouter()
	m.Path("(.*?)/git-upload-pack$").Methods("POST").Name(GitUploadPack)
	m.Path("(.*?)/git-receive-pack$").Methods("POST").Name(GitReceivePack)
	m.Path("(.*?)/info/refs$").Methods("GET").Name(GetInfoRefs)
	m.Path("(.*?)/HEAD$").Methods("GET").Name(GetHead)
	m.Path("(.*?)/objects/info/alternates$").Methods("GET").Name(GetInfoAlternates)
	m.Path("(.*?)/objects/info/http-alternates$").Methods("GET").Name(GetInfoHttpAlternates)
	m.Path("(.*?)/objects/info/packs$").Methods("GET").Name(GetInfoPacks)
	m.Path("(.*?)/objects/info/[^/]*$").Methods("GET").Name(GetInfoFile)
	m.Path("(.*?)/objects/[0-9a-f]{2}/[0-9a-f]{38}$").Methods("GET").Name(GetLooseObject)
	m.Path("(.*?)/objects/pack/pack-[0-9a-f]{40}\\.pack$").Methods("GET").Name(GetPackFile)
	m.Path("(.*?)/objects/pack/pack-[0-9a-f]{40}\\.idx$").Methods("GET").Name(GetIdxFile)
	return m
}
