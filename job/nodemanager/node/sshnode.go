package node

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"sync/atomic"

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
		n.scriptPath, err = createScriptFile(jobStep)
		if err != nil {
			panic(err) // XXX job fail
		}
		err = n.createScriptOnRemote(n.scriptPath)
		if err != nil {
			panic(err) // XXX job fail
		}

		// We dont need script file on local now.
		err = os.Remove(n.scriptPath)
		if err != nil {
			log.Warn().Err(err).Msgf("Failed to delete temporary file %s", n.scriptPath)
		}

		// we dont need local script path. but still need filename.
		n.scriptPath = path.Base(n.scriptPath)
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

func (n *sshExecNode) Run(ctx context.Context) int {

	exitcode := 0
	jobStep := n.jobStep

	log := util.GetLoggerWithSource(n.GetName(), "run").With().
		Str("jobId", n.jobStep.JobId).Str("runId", jobStep.RunId).
		Str("node", n.nodeDef.Id).Str("step", jobStep.Name).Logger()

	switch n.jobStep.ExecType {
	case objects.JOB_EXEC_TYPE_COMMAND:
		log.Trace().Msgf("Run command %s", jobStep.Command)

		exitcode = n.doCommand(ctx, jobStep.Command)
	case objects.JOB_EXEC_TYPE_SCRIPT:
		log.Trace().Msgf("Run script %s", n.scriptPath)
		exitcode = n.doCommand(ctx, "~/"+n.scriptPath)
	default:
		log.Error().Msgf("Unknown execType %s on %s/%s", jobStep.ExecType, jobStep.JobId, jobStep.Name)
		return EC_DEF_ERR
	}

	// teardown:
	if n.scriptPath != "" {
		n.deleteScriptOnRemote(n.scriptPath)
	}

	if n.sshClient != nil {
		err := n.sshClient.Close()
		if err != nil {
			log.Warn().Err(err).Msgf("SSH Close error")
		}
	}

	return exitcode
}

// doCommand exec cmd.
//  Returns exitcode. exitcode will be negative value. (exec failed etc)
func (n *sshExecNode) doCommand(ctx context.Context, cmd string) int {
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
		return -1
	}
	err = session.Wait()
	code, msg := exitCodeFromError(err)
	log.Debug().Err(err).Msgf("Exitcode: %v msg: %s", code, msg)

	return code
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

// transfer script file to target node and set chmod 700
func (n *sshExecNode) createScriptOnRemote(filepath string) error {
	log := util.GetLoggerWithSource(n.GetName(), "doSftp")

	basename := path.Base(filepath)
	script, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Err(err).Msg("Read temp script failed")
		return err
	}

	// open an SFTP session over an existing ssh connection.
	client, err := sftp.NewClient(n.sshClient)
	if err != nil {
		log.Err(err)
		return err
	}
	defer client.Close()

	f, err := client.Create(basename)
	if err != nil {
		log.Err(err).Msgf("Failed to create file %s", basename)
		return err
	}
	if _, err := f.Write(script); err != nil {
		log.Err(err).Msgf("Failed to write %s", basename)
		return err
	}
	f.Close()

	client.Chmod(basename, 0700)
	if err != nil {
		log.Err(err).Msgf("Failed to change permission %s", basename)
		return err
	}

	return nil
}

func (n *sshExecNode) deleteScriptOnRemote(filepath string) error {
	log := util.GetLoggerWithSource(n.GetName(), "doSftp")

	basename := path.Base(filepath)

	// open an SFTP session over an existing ssh connection.
	client, err := sftp.NewClient(n.sshClient)
	if err != nil {
		log.Err(err)
		return err
	}
	defer client.Close()

	// leave your mark
	err = client.Remove(basename)
	if err != nil {
		log.Err(err)
	}

	return nil
}
