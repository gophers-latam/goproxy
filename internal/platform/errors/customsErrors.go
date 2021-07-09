package errors

var (
	// errors configs
	ErrReadConfig      = NewError("config: error to load yaml file")
	ErrUnmarshalConfig = NewError("config: error to unmarsahl yaml file")
	// errors security
	ErrApiKeyGenerator     = NewError("security: error on apikey generator method")
	ErrCreatingSettingFile = NewError("security: error on create setting file")
	ErrWritingSettingFile  = NewError("security: error on write setting file")
	// loadbalancer
	ErrIsBackendAlive = NewError("lb: Site unreachcable dial tcp")
	// repositoryDB
	ErrSavekeyUpdateTX    = NewError("badgerdb: error executing TX to save apikey")
	ErrSavekeyUpdate      = NewError("badgerdb: error to save apikey")
	ErrSavekeyCreateLocal = NewError("localfile: error to save apikey")
)