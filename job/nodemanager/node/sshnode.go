package node

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/pkg/sftp"
	"github.com/yakumo-saki/phantasma-flow/global/consts"
	"github.com/yakumo-saki/phantasma-flow/job/jobparser"
	"github.com/yakumo-saki/phantasma-flow/messagehub"
	"github.com/yakumo-saki/phantasma-flow/pkg/objects"
	"github.com/yakumo-saki/phantasma-flow/util"
	"golang.org/x/crypto/ssh"
)

type sshExecNode struct {
	nodeDef objects.NodeDefinition
	jobStep jobparser.ExecutableJobStep

	seqNo      uint64 // log sequence no (use atomic.Add)
	scriptPath string // created script file path.
	sshClient  *ssh.Client
}

func (n *sshExecNode) GetName() string {
	return "sshExecNode"
}

func (n *sshExecNode) Initialize(def objects.NodeDefinition, jobStep jobparser.ExecutableJobStep) error {

	log := util.GetLoggerWithSource(n.GetName(), "Init")
	log.Info().Msgf("%s", def)

	n.nodeDef = def
	n.jobStep = jobStep

	n.connectSSH()

	// create script. if jobStep is SCRIPT
	if jobStep.ExecType == objects.JOB_EXEC_TYPE_SCRIPT {
		var err error
		n.scriptPath, err = n.createScriptFile(jobStep)
		if err != nil {
			panic(err) // XXX job fail
		}
		n.doSftp()
	}

	return nil
}

func (n *sshExecNode) connectSSH() error {
	log := util.GetLoggerWithSource(n.GetName(), "connectSSH")

	config := &ssh.ClientConfig{
		User: n.nodeDef.Ssh.User,
		Auth: []ssh.AuthMethod{
			n.getAuthMethod(),
		},
		HostKeyCallback: util.GetHostKeyCallback(n.nodeDef.Ssh.HostAuthType, n.nodeDef.Ssh.HostKey),
	}

	host := fmt.Sprintf("%s:%v", n.nodeDef.Ssh.Host, n.nodeDef.Ssh.Port)
	client, err := ssh.Dial("tcp", host, config)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to dial:%s %s", host)
	}

	n.sshClient = client
	return nil
}

func (n *sshExecNode) getAuthMethod() ssh.AuthMethod {
	switch n.nodeDef.Ssh.AuthType {
	case consts.USER_AUTHTYPE_KEY:
		signer := util.GetSignerFromKeyAndPass(n.nodeDef.Ssh.Key, n.nodeDef.Ssh.KeyPassphrase)
		return ssh.PublicKeys(signer)
	case consts.USER_AUTHTYPE_KEYFILE:
		sshkey := util.ReadPublicKeyfile(n.nodeDef.Ssh.Keyfile)
		signer := util.GetSignerFromKeyAndPass(sshkey, n.nodeDef.Ssh.Password)
		return ssh.PublicKeys(signer)
	case consts.USER_AUTHTYPE_PASSWORD:
		return ssh.Password(n.nodeDef.Ssh.Password)
	default:
		panic("unknown ssh authtype " + n.nodeDef.Ssh.AuthType)
	}
}

//
func (n *sshExecNode) createScriptFile(jobStep jobparser.ExecutableJobStep) (string, error) {
	tempFilename := fmt.Sprintf("%s_%s_*", jobStep.JobId, jobStep.Name)
	tempfile, err := os.CreateTemp("", tempFilename)
	if err != nil {
		return "", err
	}

	// if script has not shebang, /bin/bash assumed
	if !strings.HasPrefix(jobStep.Script, "#!") {
		tempfile.WriteString("#!/bin/bash\n") // XXX #50
	}
	_, err = tempfile.WriteString(jobStep.Script)
	if err != nil {
		return "", err
	}
	err = tempfile.Chmod(os.FileMode(int(0700)))
	if err != nil {
		return "", err
	}

	tempfile.Close()
	if err != nil {
		return "", err
	}

	return tempfile.Name(), nil

}

func (n *sshExecNode) Run(ctx context.Context) {

	jobStep := n.jobStep

	log := util.GetLoggerWithSource(n.GetName(), "run").With().
		Str("jobId", n.jobStep.JobId).Str("runId", jobStep.RunId).
		Str("node", n.nodeDef.Id).Str("step", jobStep.Name).Logger()

	switch n.jobStep.ExecType {
	case objects.JOB_EXEC_TYPE_COMMAND:
		log.Trace().Msgf("Run command %s", jobStep.Command)

		n.doCommand(ctx, jobStep.Command)
		// cmd = exec.CommandContext(ctx, "sh", "-c", jobStep.Command)
		// n.doCommand(ctx, "sh -c "+jobStep.Command)
	case objects.JOB_EXEC_TYPE_SCRIPT:
		// Run script created on initialize #25
		log.Trace().Msgf("Run script %s", n.scriptPath)
		// cmd = exec.CommandContext(ctx, n.scriptPath)
	default:
		panic(fmt.Sprintf("Unknown execType %s on %s/%s",
			jobStep.ExecType, jobStep.JobId, jobStep.Name))
	}

	// teardown:
	if n.scriptPath != "" {
		n.doSftpDelete()
	}

	if n.sshClient != nil {
		err := n.sshClient.Close()
		if err != nil {
			log.Warn().Err(err).Msgf("SSH Close error")
		}
	}
}

func (n *sshExecNode) doCommand(ctx context.Context, cmd string) {
	log := util.GetLoggerWithSource(n.GetName(), "doCommand")

	session, err := n.sshClient.NewSession()
	if err != nil {
		log.Error().Err(err).Msg("Failed to create session: ")
	}
	defer session.Close()

	stderr, err := session.StdoutPipe()
	if err == nil {
		go n.pipeToLog("stderr", stderr)
	} else {
		log.Err(err)
	}

	stdout, err := session.StdoutPipe()
	if err == nil {
		go n.pipeToLog("stdout", stdout)
	} else {
		log.Err(err)
	}

	log.Debug().Msgf("start")

	if err := session.Start(cmd); err != nil {
		log.Err(err).Msgf("Failed to run: %s", cmd)
	}
	err = session.Wait()
	if err != nil {
		log.Err(err).Msgf("wait err: %s", cmd)
	}

}

func (n *sshExecNode) pipeToLog(name string, pipe io.Reader) {
	// log := util.GetLoggerWithSource(n.GetName(), "pipeToLog", name)

	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		seq := atomic.AddUint64(&n.seqNo, 1)
		logmsg := scanner.Text()

		msg := createJobLogMsg(seq, n.jobStep)
		msg.Source = name
		msg.Message = logmsg
		messagehub.Post(messagehub.TOPIC_JOB_LOG, msg)
	}

}

func (n *sshExecNode) doSftp() {

	// open an SFTP session over an existing ssh connection.
	client, err := sftp.NewClient(n.sshClient)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// leave your mark
	f, err := client.Create("hello.txt")
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.Write([]byte("Hello world! " + time.Now().String())); err != nil {
		log.Fatal(err)
	}
	f.Close()

	// check it's there
	fi, err := client.Lstat("hello.txt")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fi)
}

func (n *sshExecNode) doSftpDelete() {

	// open an SFTP session over an existing ssh connection.
	client, err := sftp.NewClient(n.sshClient)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// leave your mark
	err = client.Remove("hello.txt")
	if err != nil {
		log.Fatal(err)
	}
	client.Close()
}
