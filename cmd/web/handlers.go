package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/half-blood-prince-2710/snippetbox/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	
	snippets,err := app.snippets.Latest()
	if err!=nil{
		app.serverError(w,r,err)
		return
	}

	w.WriteHeader(http.StatusOK)
	for _,snippets := range snippets {
		fmt.Fprintf(w,"%v\n",snippets)
	}
	w.Write([]byte("Hello, Snippetbox"))
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	
	snippet ,err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err,models.ErrNoRecord){ 
			http.NotFound(w,r)
		} else {
			app.serverError(w,r,err)
		}
	}
	flash := app.sessionManager.PopString(r.Context(), "flash")
	app.logger.Info("Request for snippet", "snippet", snippet,"flash",flash)
	fmt.Fprintf(w, "Display a specific snippet with ID %+v", snippet)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating a new snippet..."))
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {

	slog.Info("enter create post")
	title,content,expires:="O snail","O beautiful snail, \nBUT SLOWLY<SLOWLY! \n - IDIOT SAN" , 7

	id,err := app.snippets.Insert(title,content,expires)
	if err!=nil {
		app.serverError(w,r,err)
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")
	http.Redirect(w,r,fmt.Sprintf("/snippet/view/%d",id),http.StatusSeeOther)
}
