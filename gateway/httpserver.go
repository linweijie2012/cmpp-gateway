package gateway

import (
	"cmpp-gateway/pages"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

var pageSize = 5

// handler echoes the HTTP request.
func handler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Print(err)
	}
	src := r.Form.Get("src")
	content := r.Form.Get("cont")
	dest := r.Form.Get("dest")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if src == "" || content == "" || dest == "" {
		result, _ := json.Marshal(
			map[string]interface{}{"result": -1, "error": "请输入 参数'src' 'dest' 'const' 缺一不可"})
		fmt.Fprintf(w, string(result))
		return
	}
	mes := SmsMes{Src: src, Content: content, Dest: dest}
	Messages <- mes
	result, _ := json.Marshal(
		map[string]interface{}{"error": "", "result": 0})
	fmt.Fprintf(w, string(result))
}

func index(w http.ResponseWriter, r *http.Request) {
	findTemplate(w, r, "index.html")

}

func findTemplate(w http.ResponseWriter, r *http.Request, tpl string) {
	t, error := template.New(tpl).ParseFiles(tpl)
	if error != nil {
		fmt.Fprintf(w, "template error %v", error)
		return
	}

	err := t.Execute(w, struct{}{})
	if err != nil {
		fmt.Fprintf(w, "error %v", err)
		return
	}
}

func listMessage(w http.ResponseWriter, r *http.Request, listName string) {
	r.ParseForm()
	parameter := r.Form.Get("page")

	var c_page int
	if parameter == "" {
		c_page = 1
	} else {
		c_page, _ = strconv.Atoi(parameter)
	}
	count := SCache.Length(listName)
	page := pages.NewPage(c_page, pageSize, count)
	t, err := template.New(listName + ".html").ParseFiles(listName + ".html")
	if err != nil {
		fmt.Fprintf(w, "template error %v", err)
		return
	}
	v := SCache.GetList(listName, page.StartRow, page.EndRow)
	ret := map[string]interface{}{
		"data": v,
		"page": page,
	}
	err = t.Execute(w, ret)
	if err != nil {
		fmt.Fprintf(w, "error %v", err)
		return
	}
}

func messages(w http.ResponseWriter, r *http.Request, listName string) {

	v := SCache.GetList(listName, 0, 1000)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	result, _ := json.Marshal(
		map[string]interface{}{"data": v})

	fmt.Fprintf(w, string(result))

}

func messages_in(w http.ResponseWriter, r *http.Request) {
	messages(w, r, "list_mo")
}
func messages_out(w http.ResponseWriter, r *http.Request) {
	messages(w, r, "list_message")
}

func listSubmits(w http.ResponseWriter, r *http.Request) {
	listMessage(w, r, "list_message")
}

func listMo(w http.ResponseWriter, r *http.Request) {
	listMessage(w, r, "list_mo")
}

func Serve(config *Config) {

	http.HandleFunc("/send", handler)
	http.HandleFunc("/messages_in", messages_in)
	http.HandleFunc("/messages_out", messages_out)
	http.HandleFunc("/", index)
	http.HandleFunc("/list_message", listSubmits)
	http.HandleFunc("/list_mo", listMo)
	log.Fatal(http.ListenAndServe(config.HttpHost+":"+config.HttpPort, nil))
}
