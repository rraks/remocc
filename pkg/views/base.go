package views


import (
    "github.com/rraks/remocc/pkg/models"
    "net/http"
    "log"
)


var tableView *View

func init() {
  tableView = NewView("base", "web/views/templates/deviceTableRow.html")
}


// TODO : Create a more generic renderer function
func RenderTableRow(w http.ResponseWriter, Devs []*models.Device) {
    row := struct {
        Devices []*models.Device
    }{Devs}
    err := tableView.Render(w, row)
    if err != nil {
        log.Println(err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

