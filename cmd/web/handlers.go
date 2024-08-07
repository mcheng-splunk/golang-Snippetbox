package main

import (
	"errors"
	"fmt" // New import
	"net/http"
	"strconv"
	"snippetbox.mcheng.net/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
    w.Header().Add("Server", "Go")

    snippets, err := app.snippets.Latest()
    if err != nil {
        app.serverError(w, r, err)
        return
    }

    // Use the new render helper.
    app.render(w, r, http.StatusOK, "home.tmpl", templateData{
        Snippets: snippets,
    })
    // for _, snippet := range snippets {
    //     fmt.Fprintf(w, "%+v\n", snippet)
    // }
    // Use the template.ParseFiles() function to read the template file into a
    // template set. If there's an error, we log the detailed error message, use
    // the http.Error() function to send an Internal Server Error response to the
    // user, and then return from the handler so no subsequent code is executed.
    // Initialize a slice containing the paths to the two files. It's important
    // to note that the file containing our base template must be the *first*
    // file in the slice.
    // files := []string{
    //     "./ui/html/base.tmpl",
    //     "./ui/html/partials/nav.tmpl",
    //     "./ui/html/pages/home.tmpl",
    // }

    // //     // Use the template.ParseFiles() function to read the files and store the
    // //     // templates in a template set. Notice that we use ... to pass the contents 
    // //     // of the files slice as variadic arguments.
    // ts, err := template.ParseFiles(files...)
    // if err != nil {
    //     // app.logger.Error(err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
    //     //http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    //     app.serverError(w, r, err)
    //     return
    // }

    // data := templateData{
    //     Snippets: snippets,
    // }

    // // Pass in the templateData struct when executing the template.
    // err = ts.ExecuteTemplate(w, "base", data)
    // if err != nil {
    //     app.serverError(w, r, err)
    // }
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(r.PathValue("id"))
    if err != nil || id < 1 {
        http.NotFound(w, r)
        return
    }
    snippet, err := app.snippets.Get(id)
    if err != nil {
        if errors.Is(err, models.ErrNoRecord) {
            http.NotFound(w, r)
        } else {
            app.serverError(w, r, err)
        }
        return
    }
    fmt.Printf("id is %d\n ", snippet.ID)
    // Use the new render helper.
    app.render(w, r, http.StatusOK, "view.tmpl", templateData{
        Snippet: snippet,
    })

    // files := []string{
    // "./ui/html/base.tmpl",
    // "./ui/html/partials/nav.tmpl",
    // "./ui/html/pages/view.tmpl",
    // }

    // // // Use the template.ParseFiles() function to read the files and store the
    // // // templates in a template set. Notice that we use ... to pass the contents 
    // // // of the files slice as variadic arguments.
    // ts, err := template.ParseFiles(files...)
    // if err != nil {
    //     // app.logger.Error(err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
    //     //http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    //     app.serverError(w, r, err)
	// 	return
    // }

    // // Create an instance of a templateData struct holding the snippet data.
    // data := templateData{
    //     Snippet: snippet,
    // }

    // // Pass in the templateData struct when executing the template.
    // err = ts.ExecuteTemplate(w, "base", data)
    // if err != nil {
    //     app.serverError(w, r, err)
    // }
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Display a form for creating a new snippet..."))
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
// Create some variables holding dummy data. We'll remove these later on
    // during the build.
    title := "O snail"
    content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n- Kobayashi Issa"
    expires := 7

    // Pass the data to the SnippetModel.Insert() method, receiving the
    // ID of the new record back.
    id, err := app.snippets.Insert(title, content, expires)
    if err != nil {
        app.serverError(w, r, err)
        return
    }

    // Redirect the user to the relevant page for the snippet.
    http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}