package consts

const NODE_TYPE_LOCAL = "local"
const NODE_TYPE_SSH = "ssh"
const NODE_TYPE_WINRM = "winrm"

const USER_AUTHTYPE_KEY = "key"           // SSH Auth ssh key string
const USER_AUTHTYPE_KEYFILE = "keyfile"   // SSH Auth key file
const USER_AUTHTYPE_PASSWORD = "password" // SSH Auth password string

const HOST_AUTHTYPE_KEY = "key"               // SSH Host Auth key string
const HOST_AUTHTYPE_IGNORE = "insecureIgnore" // SSH Host Auth No check(insecure)
