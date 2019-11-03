package main

import (
	"fmt"
	"github.com/goodsign/monday"
	"github.com/julienbayle/listedeseleves/pointage"
	"github.com/zserge/webview"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

// Input filename
var studentsFileName = ""

func startServer() string {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		defer ln.Close()
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`
			<!doctype html>
			<html>
				<head>
					<meta http-equiv="X-UA-Compatible" content="IE=edge">
				</head>
				<body>
					<p>1 - Sélectionner un fichier source : <button onclick="external.invoke('open')">Ouvrir</button></p>
					<p>2 - Choisir un mois (Exemple : 10/2019) <input id="month" type="text" />
					<p>3 - Lancer le générateur <button onclick="external.invoke('generate:'+document.getElementById('month').value)">
						Générer la fiche de pointage
					</button>
					<p> (Penser à utiliser un nom de fichier avec l'extension xlsx)
				</body>
			</html>
			`))
		})
		log.Fatal(http.Serve(ln, nil))
	}()
	return "http://" + ln.Addr().String()
}

func handleRPC(w webview.WebView, data string) {
	switch {
	case data == "open":
		studentsFileName = w.Dialog(webview.DialogTypeOpen, webview.DialogFlagFile, "Fichier source", "")
		log.Println("open ", studentsFileName)
	case strings.HasPrefix(data, "generate:"):
		defer func() {
			if r := recover(); r != nil {
				w.Dialog(webview.DialogTypeAlert, webview.DialogFlagError, "Erreur", r.(error).Error())
			}
		}()
		month, err := time.Parse("01/2006", strings.TrimPrefix(data, "generate:"))
		if err != nil {
			w.Dialog(webview.DialogTypeAlert, webview.DialogFlagError, "Erreur", err.Error())
		}
		monthFR := monday.Format(month, "January 2006", monday.LocaleFrFR)
		var exportFileName = w.Dialog(webview.DialogTypeSave, webview.DialogFlagFile, "Nom du fichier de sortie", fmt.Sprintf("pointage-%s.xslx", monthFR))
		log.Println("save to ", exportFileName)
		classesOfStudents := pointage.Load(studentsFileName)
		pointage.Export(classesOfStudents, month, exportFileName)
		w.Dialog(webview.DialogTypeAlert, webview.DialogFlagInfo, "Information", "Fichier généré avec succès")
	}
}

func main() {
	url := startServer()
	w := webview.New(webview.Settings{
		Width:                  550,
		Height:                 200,
		Title:                  "Pointage OGEC",
		Resizable:              true,
		URL:                    url,
		ExternalInvokeCallback: handleRPC,
	})
	defer w.Exit()
	w.Run()
}
