package microservice

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

var (
	//CtxACL is context key for acl
	CtxACL = contextKey("acl")
	//CtxUserID is context key for user id
	CtxUserID = contextKey("uid")
	//CtxUserID is context key for user id msb
	CtxUserIDMsb = contextKey("uid_msb")
	//CtxUserID is context key for user id lsb
	CtxUserIDLsb = contextKey("uid_lsb")
	//CtxUserID is context key for user id lsb
	CtxUserUUID = contextKey("uuid")
	//CtxDomain is context key for domain
	CtxDomain = contextKey("domain")
	//CtxPhone is context key for phone number
	CtxPhone = contextKey("phone_number")
	//CtxEmail is context key for email
	CtxEmail = contextKey("email")
	//CtxName is context key for user name
	CtxName = contextKey("name")
	//CtxTopic is context key for topic name
	CtxTopic = contextKey("topic")
)
