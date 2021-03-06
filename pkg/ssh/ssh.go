package ssh

import (
    "os"
    "os/exec"
    "log"
    "sync"
    "text/template"
	"regexp"
    "net"
    "bytes"
    "strings"
    "io"
    "io/ioutil"
)

type AuthKey struct {
    Port string
    Key string
}

var sshdMutx = &sync.Mutex{}
var tmpl = &template.Template{}

func init() {
    tmpl, _ = template.New("authKeyTempl").Parse("command=" +
                "\"echo 'This account can only be used for [reverse tunneling]'\"," +
                "no-agent-forwarding,no-X11-forwarding," +
                "permitlisten=\"localhost:{{.Port}}\",permitopen=\"localhost:{{.Port}}\" {{.Key}} deviceKey\n")
    // Start the openssh server (No RC in docker image)
    //_ = exe_cmd(sshdCmd)
}


func exe_cmd(cmd string, args ...string) ([]byte) {
    out, err := exec.Command(cmd, args...).Output()
    if err != nil {
      log.Println("Exec Failed")
      log.Println(err)
    }
    return out
}

func genRandomPort() string {
    rexp, _ := regexp.Compile("[\\d]+")
    iface, _ := net.Listen("tcp", ":0")
    defer iface.Close()
    return rexp.FindString(iface.Addr().String())
}

func AddDeviceKey(email_tbl string,sshKey string) string {
    sshdMutx.Lock()
    port := genRandomPort()
    cfg := &AuthKey{port, sshKey}
    entry := new(bytes.Buffer)
    err := tmpl.Execute(entry, cfg)
    if err != nil {
        log.Println("Failed to execute template")
    }
    authKeysFile := "/home/"+email_tbl+"/.ssh/authorized_keys"

    f, err := os.OpenFile(authKeysFile, os.O_APPEND|os.O_WRONLY, 0600)
    if err != nil {
        panic(err)
    }
    defer f.Close()
    if _, err = f.WriteString(entry.String()); err != nil {
        panic(err)
      }
    sshdMutx.Unlock()
    return port
}

func AddUserKey(email_tbl string, sshKey string) {
    sshdMutx.Lock()
    authKeysFile := "/home/"+email_tbl+"/.ssh/authorized_keys"
    f, err := os.OpenFile(authKeysFile, os.O_APPEND|os.O_WRONLY, 0600)
    if err != nil {
        panic(err)
    }
    defer f.Close()
    if _, err = f.WriteString(sshKey + " userKey" + "\n"); err != nil {
        panic(err)
      }
    sshdMutx.Unlock()
}

func DelDeviceKey(email_tbl string, sshKey string) (error) {
    sshdMutx.Lock()
    authKeysFile := "/home/"+email_tbl+"/.ssh/authorized_keys"
    input, err := ioutil.ReadFile(authKeysFile)
    if err != nil {
        log.Println("Auth file not found")
    }
    authKeys := strings.Split(string(input), "\n")
    for i, key := range authKeys {
        if strings.Contains(key, sshKey) {
            if strings.Contains(key,"deviceKey") {
                authKeys = append(authKeys[:i], authKeys[i+1:]...)
            }
        }
    }
    output := strings.Join(authKeys, "\n")
    err = ioutil.WriteFile(authKeysFile, []byte(output), 0600)
    if err != nil {
        log.Println("Failed deleting key")
    }
    sshdMutx.Unlock()
    return err
}

func AddUser(email string, password string, sshKey string) {
    var b2 bytes.Buffer
    email_tbl := strings.Replace(email,"@","_",-1)
    email_tbl = strings.Replace(email_tbl,".","_",-1)
    // Pipe commands, make function
    c1 := exec.Command("echo","-e",password+"\n"+password)
    c2 := exec.Command("adduser", email_tbl)
    r, w := io.Pipe()
    c1.Stdout = w
    c2.Stdin = r
    c2.Stdout = &b2
    err1 := c1.Start()
    err2 := c2.Start()
    c1.Wait()
    w.Close()
    c2.Wait()
    if err1 != nil {
        log.Println(err1)
    }
    if err2 != nil {
        log.Println(err2)
    }
    exe_cmd("mkdir", "/home/"+email_tbl+"/.ssh")
    exe_cmd("touch", "/home/"+email_tbl+"/.ssh/authorized_keys")
    AddUserKey(email_tbl, sshKey)
}

func StatConnection() {
}

