package main

import (
	"io"
	"net/http"
	"os"
	"strconv"
	"text/template"
)

var out string

var check1, check2 bool

// Starting to find the html page
func StartPage(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" || req.Method != "GET" {
		error404(res)
		return
	} else {
		pageTemplate, err := template.ParseFiles("./templates/home.page.tmpl")
		if err != nil {
			error500(res)
			return
		}
		res.WriteHeader(200)
		pageTemplate.Execute(res, nil)
	}
}

// Form data after submiting

func SubmitTing(res http.ResponseWriter, req *http.Request) {

	if req.FormValue("sub") == "sub" {
		check1 = true
	}
	if req.Method != "POST" {
		error405(res)
		return
	}
	if req.URL.Path != "/ascii-art" {
		error404(res)
		return
	}
	tpl, err := template.ParseFiles("./templates/home.page.tmpl")
	req.ParseForm()
	if req.FormValue("banner") != "standard" &&
		req.FormValue("banner") != "thinkertoy" &&
		req.FormValue("banner") != "shadow" || err != nil {
		error500(res)
		return
	}

	txt := req.FormValue("input")
	ban := req.FormValue("banner")
	i, _ := os.ReadFile(ban + ".txt")
	if ban == "shadow" {
		if len(i) != 7463 {
			error500(res)
			return
		}
	}
	if ban == "standard" {
		if len(i) != 26488 {
			error500(res)
			return
		}
	}
	if ban == "thinkertoy" {
		if len(i) != 4703 {
			error500(res)
			return
		}
	}

	_, err = os.ReadFile(ban + ".txt")
	if err != nil {
		error500(res)
		return
	}
	// isValid find if there an anpritable caracter
	if !isValid(txt) {
		error400(res)
		return
	}
	if err != nil {
		// log.Fatalln(err)
		error500(res)
		return
	}
	res.WriteHeader(200)
	tpl.Execute(res, ConvertStr(txt, ban))

	out = ConvertStr(txt, ban)

}

func download(res http.ResponseWriter, req *http.Request) {

	data := os.MkdirAll("data", 0755)
	if data != nil {
		error500(res)
		return
	}
	// Enregistrer le texte ASCII dans un fichier
	if req.FormValue("fileformat") == "txt" {
		data = os.WriteFile("./data/result.txt", []byte(out), 0644)
		if data != nil {
			error500(res)
			return
		}
	}

	if req.FormValue("fileformat") == "doc" {
		data = os.WriteFile("./data/result.doc", []byte(out), 0644)
		if data != nil {
			error500(res)
			return
		}

	}

	if req.FormValue("telecharge") == "telecharge" {
		check2 = true
	}
	if check1 && check2 {

		formatType := req.FormValue("fileformat")

		filename := "./data/result." + formatType

		// Ouvrir le fichier à télécharger
		f, err := os.Open(filename)
		if err != nil {
			error500(res)
			return
		}
		defer f.Close()

		// Récupérer les informations sur le fichier
		fileInfo, err := f.Stat()
		if err != nil {
			error500(res)
			return
		}
		fileSize := fileInfo.Size()

		// Définition des en-têtes
		res.Header().Set("Content-Disposition", "attachment; filename=result."+formatType)
		res.Header().Set("Content-Type", "text/plain")
		res.Header().Set("Content-Length", strconv.FormatInt(fileSize, 10))
		f.Seek(0, io.SeekStart)
		// Transmission du contenu du fichier
		io.Copy(res, f)
	} else {
		error400(res)
	}
	check1 = false
	check2 = false
}
