package app

import (
	"crud/pkg/crud/models"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
)

func (receiver *server) handleBurgersList() func(http.ResponseWriter, *http.Request) {
	tpl, err := template.ParseFiles(filepath.Join(receiver.templatesPath, "index.gohtml"))
	if err != nil {
		log.Printf("can't parse index page: %v", err)
		panic(err)
	}
	return func(writer http.ResponseWriter, request *http.Request) {
		list, err := receiver.burgersSvc.BurgersList()
		if err != nil {
			log.Printf("can't execute Burgers list sevice: %v", err)
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		data := struct {
			Title   string
			Burgers []models.Burger
		}{
			Title:   "BURGER KING",
			Burgers: list,
		}

		err = tpl.Execute(writer, data)
		if err != nil {
			log.Printf("can't execute print burgers data: %v", err)
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}

func (receiver *server) handleBurgersSave() func(responseWriter http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {

		name := request.FormValue("name")
		price := request.FormValue("price")

		parsedPrice, err := strconv.Atoi(price)
		if err != nil {
			log.Printf("incorect data from request: %v", err)
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		err = receiver.burgersSvc.Save(models.Burger{Name: name, Price: parsedPrice})
		if err != nil {
			log.Printf("error while saving burger: %v", err)
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}


		http.Redirect(writer, request, "/", http.StatusPermanentRedirect)
		return
	}
}

func (receiver *server) handleBurgersRemove() func(responseWriter http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {

		id := request.FormValue("id")

		parsedId, err := strconv.Atoi(id)
		if err != nil {
			log.Printf("incorect data from request: %v", err)
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		err = receiver.burgersSvc.RemoveById(parsedId)
		if err != nil {
			log.Printf("error while removing burger: %v", err)
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		http.Redirect(writer, request, "/", http.StatusPermanentRedirect)
		return
	}
}

func (receiver *server) handleFavicon() func(http.ResponseWriter, *http.Request) {
	file, err := ioutil.ReadFile(filepath.Join(receiver.assetsPath, "favicon.ico"))
	if err != nil {
		log.Printf("can't read favicon file: %v", err)
		panic(err)
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		_, err := writer.Write(file)
		if err != nil {
			log.Printf("error while sent favicon: %v", err)
		}
	}
}