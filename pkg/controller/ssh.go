package controller

import (
    "os"
    "os/exec"
    "sync"
    "log"
    "text/template"
	"regexp"
    "net"
    "bytes"
)

type AuthKey struct {
    Port string
    Key string
}

var sshdMutx = &sync.Mutex{}
var tmpl = &template.Template{}
var sshdConf string
var sshdCmd string
var authKeysFile string

func init() {
    tmpl, _ = template.New("authKeyTempl").Parse("command=" +
                "\"echo 'This account can only be used for [reverse tunneling]'\"," +
                "no-agent-forwarding,no-X11-forwarding," +
                "permitlisten=\"localhost:{{.Port}}\",permitopen=\"localhost:{{.Port}}\" {{.Key}}\n")
    sshdCmd =  os.Getenv("SSHD_CMD")
    authKeysFile = os.Getenv("AUTHORIZED_KEYS")
    log.Println("sshdCmd - ", sshdCmd)
    log.Println("authKeysFile - ", authKeysFile)
    // Start the openssh server (No RC in docker image)
    //_ = exe_cmd(sshdCmd)

}


func exe_cmd(cmd string) ([]byte) {
    sshdMutx.Lock()
    out, err := exec.Command(cmd).Output()
    if err != nil {
      log.Println("Exec Failed")
    }
    return out
}

func genRandomPort() string {
    rexp, _ := regexp.Compile("[\\d]+")
    iface, _ := net.Listen("tcp", ":0")
    defer iface.Close()
    return rexp.FindString(iface.Addr().String())
}

func AddKey(sshKey string) {
    port := genRandomPort()
    cfg := &AuthKey{port, sshKey}
    entry := new(bytes.Buffer)
    err := tmpl.Execute(entry, cfg)
    if err != nil {
        log.Println("Failed to execute template")
    }

    log.Println("Reached ssh.go")
    f, err := os.OpenFile(authKeysFile, os.O_APPEND|os.O_WRONLY, 0600)
    if err != nil {
        panic(err)
    }
    defer f.Close()
    if _, err = f.WriteString(entry.String()); err != nil {
        panic(err)
      }

}
