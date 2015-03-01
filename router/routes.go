package router

const (
	GitUploadPack         = "git-upload-pack"
	GitReceivePack        = "git-receive-pack"
	GetInfoRefs           = "get-info-refs"
	GetHead               = "get-head"
	GetInfoAlternates     = "get-info-alternates"
	GetInfoHttpAlternates = "get-info-http-alternates"
	GetInfoPacks          = "get-info-packs"
	GetInfoFile           = "get-info-file"
	GetLooseObject        = "get-loose-object"
	GetPackFile           = "get-pack-file"
	GetIdxFile            = "get-idx-file"

	Project           = "project"
	CreateProject     = "project:create"
	UpdateProject     = "project:update"
	RunProjectCommand = "project:runcommand"

	Service       = "service"
	Services      = "services"
	CreateService = "service:create"
)
