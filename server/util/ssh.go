package util

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/yakumo-saki/phantasma-flow/global/consts"
	"golang.org/x/crypto/ssh"
)

func GetSignerFromKeyAndPass(keystring, passphrase string) ssh.Signer {
	keyBytes := []byte(keystring)

	var signer ssh.Signer
	var err error
	if passphrase == "" {
		signer, err = ssh.ParsePrivateKey(keyBytes)
		if err != nil {
			panic(err)
		}
	} else {
		signer, err = ssh.ParsePrivateKeyWithPassphrase(keyBytes, []byte(passphrase))
		if err != nil {
			log.Println(err)
			panic(err)
		}
	}

	return signer
}

func GetHostKeyCallback(hostKeyAuthType, hostKey string) ssh.HostKeyCallback {
	switch hostKeyAuthType {
	case consts.HOST_AUTHTYPE_IGNORE:
		return ssh.InsecureIgnoreHostKey()
	case consts.HOST_AUTHTYPE_KEY:
		_, _, pubKey, _, _, err := ssh.ParseKnownHosts([]byte(hostKey))
		if err != nil {
			panic("error " + err.Error())
		}
		return ssh.FixedHostKey(ssh.PublicKey(pubKey))
	default:
		panic("unknown hostKeyAuthType " + hostKeyAuthType)
	}

}

func ReadPublicKeyfile(path string) string {
	realpath := ApplyPhFlowPath(path)

	privKey, err := ioutil.ReadFile(realpath)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	return string(privKey)
}

func ApplyPhFlowPath(path string) string {
	newPath := path
	envs := []string{
		consts.ENV_HOME_DIR,
		consts.ENV_DATA_DIR,
		consts.ENV_DEF_DIR}

	for _, env := range envs {
		envstr := fmt.Sprintf("${%s}", env)
		newPath = strings.ReplaceAll(newPath, envstr, os.Getenv(env))
	}

	// warning
	if strings.Contains(newPath, "${") {
		log.Println("path still contains ${ =>" + newPath)
	}
	return newPath
}
