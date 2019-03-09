package views


import (
    "github.com/rraks/remocc/pkg/models"
    "net/http"
    "log"
    "bytes"
    "html/template"
)


var tableView *View
var deviceView *View

func init() {
  tableView = NewView("base", "web/views/templates/deviceTableRow.html")
  tableView = NewView("base", "web/views/templates/deviceTableRow.html")
}


// TODO : Create a more generic renderer function
func RenderTableRow(w http.ResponseWriter, Devs []*models.Device) {
    rowHolder := struct {
        Devices []*models.Device
    }{Devs}
    err := tableView.Render(w, rowHolder)
    if err != nil {
        log.Println(err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}


func RenderDevicePreview(w http.ResponseWriter, logs []*models.DeviceLog, device *models.Device) ([]byte, error) {
    t := template.New("action")
    t, err := template.ParseFiles("web/views/templates/deviceView.html")
    if err != nil {
        return []byte("Failed to fetch template"), err
    }
    var tpl  bytes.Buffer
    logHolder := struct {
        Logs []*models.DeviceLog
        Device *models.Device
    }{logs, device}
    if err = t.Execute(&tpl, logHolder); err != nil {
        return []byte("Failed to execute template"), err
    }
    res := tpl.Bytes()
    return res, nil
}
