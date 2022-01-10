package node

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync/atomic"
	"time"

	"github.com/pkg/sftp"
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
	pk, err := ioutil.ReadFile(os.Getenv("HOME") + "/.ssh/id_rsa_nopass")
	if err != nil {
		panic("failed to read ssh key")
	}

	signer, err := ssh.ParsePrivateKey(pk)
	if err != nil {
		fmt.Println(err)
		panic("failed to parse ssh key")
	}

	config := &ssh.ClientConfig{
		User: "yakumo",
		Auth: []ssh.AuthMethod{
			ssh.Password("empty"),
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", "192.168.10.20:22", config)
	if err != nil {
		log.Fatal("Failed to dial: ", err)
	}

	n.sshClient = client
	return nil
}

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

	var err error
	var cmd *exec.Cmd
	switch n.jobStep.ExecType {
	case objects.JOB_EXEC_TYPE_COMMAND:
		log.Trace().Msgf("Run command %s", jobStep.Command)
		cmd = exec.CommandContext(ctx, "sh", "-c", jobStep.Command)
		n.doCommand(ctx, "sh -c "+jobStep.Command)
	case objects.JOB_EXEC_TYPE_SCRIPT:
		// Run script created on initialize #25
		log.Trace().Msgf("Run script %s", n.scriptPath)
		cmd = exec.CommandContext(ctx, n.scriptPath)
	default:
		panic(fmt.Sprintf("Unknown execType %s on %s/%s",
			jobStep.ExecType, jobStep.JobId, jobStep.Name))
	}
	stderr, err := cmd.StderrPipe()
	if err == nil {
		go n.PipeToLog(ctx, "stderr", stderr)
	} else {
		log.Err(err)
	}

	stdout, err := cmd.StdoutPipe()
	if err == nil {
		go n.PipeToLog(ctx, "stdout", stdout)
	} else {
		log.Err(err)
	}

	err = cmd.Run() // block until process exit
	if err != nil {
		log.Err(err)
	}

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

func (n *sshExecNode) PipeToLog(ctx context.Context, name string, pipe io.Reader) {
	// log := util.GetLoggerWithSource(n.GetName(), "run", name)

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

func (n *sshExecNode) doCommand(ctx context.Context, cmd string) {
	session, err := n.sshClient.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: ", err)
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Start(cmd); err != nil {
		log.Fatal("Failed to run: " + cmd + " " + err.Error())
	}
	err = session.Wait()
	if err != nil {
		log.Fatal("wait err: " + cmd + " " + err.Error())
	}
	fmt.Println(b.String())

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
