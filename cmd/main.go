package main

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Template struct {
	tmpl *template.Template
}

func newTemplate() *Template {
	return &Template{
		tmpl: template.Must(template.ParseGlob("views/*.html")),
	}
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.tmpl.ExecuteTemplate(w, name, data)
}

type Contact struct {
	Name  string
	Email string
}

type Contacts = []Contact

type ContactData struct {
	Contacts Contacts
}

func createNewData() ContactData {
	return ContactData{
		Contacts: []Contact{
			createNewContact("John Doe", "johndoe@gmail.com"),
			createNewContact("Isopropyl Alcohol", "ipa@gmail.com"),
		},
	}
}

func createNewContact(name, email string) Contact {
	return Contact{
		Name:  name,
		Email: email,
	}
}

func main() {
	e := echo.New()

	contacts := createNewData()

	e.Renderer = newTemplate()
	e.Use(middleware.Logger())

	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index", contacts)
	})

	e.POST("/contacts", func(c echo.Context) error {
		name, email := c.FormValue("name"), c.FormValue("email")

		contacts.Contacts = append(contacts.Contacts, createNewContact(name, email))

		return c.Render(200, "contacts-list", contacts)
	})

	e.Logger.Fatal(e.Start(":3001"))
}
