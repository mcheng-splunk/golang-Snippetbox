package main

import (
	"errors"
	"fmt" // New import
	"net/http"
	"strconv"
	"snippetbox.mcheng.net/internal/models"
    // "strings"      // New import
    // "unicode/utf8" // New import
    "snippetbox.mcheng.net/internal/validator"
)
// Define a snippetCreateForm struct to represent the form data and validation
// errors for the form fields. Note that all the struct fields are deliberately
// exported (i.e. start with a capital letter). This is because struct fields
// must be exported in order to be read by the html/template package when
// rendering the template.
type snippetCreateForm struct {
    Title       string
    Content     string
    Expires     int
    // FieldErrors map[string]string
    validator.Validator
}


func (app *application) home(w http.ResponseWriter, r *http.Request) {
    // w.Header().Add("Server", "Go")

    snippets, err := app.snippets.Latest()
    if err != nil {
        app.serverError(w, r, err)
        return
    }

    // // Use the new render helper.
    // app.render(w, r, http.StatusOK, "home.tmpl", templateData{
    //     Snippets: snippets,
    // })
    // Debugging: Print the number of records fetched
    fmt.Printf("Number of snippets fetched: %d\n", len(snippets))

    // Call the newTemplateData() helper to get a templateData struct containing
    // the 'default' data (which for now is just the current year), and add the
    // snippets slice to it.
    data := app.newTemplateData(r)
    data.Snippets = snippets

    // Also print the fetched snippets for review
    fmt.Printf("Fetched snippets: %+v\n", snippets)

    // Pass the data to the render() helper as normal.
    app.render(w, r, http.StatusOK, "home.tmpl", data)
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
    // // Use the new render helper.
    // app.render(w, r, http.StatusOK, "view.tmpl", templateData{
    //     Snippet: snippet,
    // })
        // And do the same thing again here...
    data := app.newTemplateData(r)
    data.Snippet = snippet

    app.render(w, r, http.StatusOK, "view.tmpl", data)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
    // w.Write([]byte("Display a form for creating a new snippet..."))
    data := app.newTemplateData(r)

    // Initialize a new createSnippetForm instance and pass it to the template.
    // Notice how this is also a great opportunity to set any default or
    // 'initial' values for the form --- here we set the initial value for the 
    // snippet expiry to 365 days.
    data.Form = snippetCreateForm{
        Expires: 365,
    }

    app.render(w, r, http.StatusOK, "create.tmpl", data)
}



func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {

    
 // First we call r.ParseForm() which adds any data in POST request bodies
    // to the r.PostForm map. This also works in the same way for PUT and PATCH
    // requests. If there are any errors, we use our app.ClientError() helper to 
    // send a 400 Bad Request response to the user.
    err := r.ParseForm()
    if err != nil {
        app.clientError(w, http.StatusBadRequest)
        return
    }

    // // Use the r.PostForm.Get() method to retrieve the title and content
    // // from the r.PostForm map.
    // title := r.PostForm.Get("title")
    // content := r.PostForm.Get("content")

    // The r.PostForm.Get() method always returns the form data as a *string*.
    // However, we're expecting our expires value to be a number, and want to
    // represent it in our Go code as an integer. So we need to manually covert
    // the form data to an integer using strconv.Atoi(), and we send a 400 Bad
    // Request response if the conversion fails.
    expires, err := strconv.Atoi(r.PostForm.Get("expires"))
    if err != nil {
        app.clientError(w, http.StatusBadRequest)
        return
    }

    // Create an instance of the snippetCreateForm struct containing the values
    // from the form and an empty map for any validation errors.
    form := snippetCreateForm{
        Title:       r.PostForm.Get("title"),
        Content:     r.PostForm.Get("content"),
        Expires:     expires,
        // FieldErrors: map[string]string{},
    }

    // Because the Validator struct is embedded by the snippetCreateForm struct,
    // we can call CheckField() directly on it to execute our validation checks.
    // CheckField() will add the provided key and error message to the
    // FieldErrors map if the check does not evaluate to true. For example, in
    // the first line here we "check that the form.Title field is not blank". In
    // the second, we "check that the form.Title field has a maximum character
    // length of 100" and so on.
    form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
    form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
    form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
    form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

    // Use the Valid() method to see if any of the checks failed. If they did,
    // then re-render the template passing in the form in the same way as
    // before.
    if !form.Valid() {
        data := app.newTemplateData(r)
        data.Form = form
        app.render(w, r, http.StatusUnprocessableEntity, "create.tmpl", data)
        return
    }


    // // Check that the title value is not blank and is not more than 100
    // // characters long. If it fails either of those checks, add a message to the
    // // errors map using the field name as the key.
    // if strings.TrimSpace(form.Title) == "" {
    //     form.FieldErrors["title"] = "This field cannot be blank"
    // } else if utf8.RuneCountInString(form.Title) > 100 {
    //     form.FieldErrors["title"] = "This field cannot be more than 100 characters long"
    // }

    // // Check that the Content value isn't blank.
    // if strings.TrimSpace(form.Content) == "" {
    //     form.FieldErrors["content"] = "This field cannot be blank"
    // }

    // // Check the expires value matches one of the permitted values (1, 7 or
    // // 365).
    // if form.Expires != 1 && form.Expires != 7 && form.Expires != 365 {
    //     form.FieldErrors["expires"] = "This field must equal 1, 7 or 365"
    // }

    // // If there are any validation errors, then re-display the create.tmpl template,
    // // passing in the snippetCreateForm instance as dynamic data in the Form 
    // // field. Note that we use the HTTP status code 422 Unprocessable Entity 
    // // when sending the response to indicate that there was a validation error.
    // if len(form.FieldErrors) > 0 {
    //     data := app.newTemplateData(r)
    //     data.Form = form
    //     app.render(w, r, http.StatusUnprocessableEntity, "create.tmpl", data)
    //     return
    // }

    id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
    if err != nil {
        app.serverError(w, r, err)
        return
    }

    http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}