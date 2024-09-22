package main

import (
	"fmt"
	"html/template"
	"io"
	"strconv"

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

type FormData struct {
	Values map[string]string
	Errors map[string]string
}

var id int = 0

type Contact struct {
	Id    int
	Name  string
	Email string
}

type Contacts = []Contact

type PageData struct {
	Form     FormData
	Data ContactData
}

type ContactData struct {
	Contacts Contacts
}

func createDefaultContacts() ContactData {
	return ContactData{
		Contacts: []Contact{
			createNewContact("John Doe", "johndoe@gmail.com"),
			createNewContact("Isopropyl Alcohol", "ipa@gmail.com"),
		},
	}
}

func createEmptyFormData() FormData {
	return FormData{
		Values: map[string]string{},
		Errors: map[string]string{},
	}
}

func createDefaultPageData() PageData {
	return PageData{
		Form:     createEmptyFormData(),
		Data: createDefaultContacts(),
	}
}

func createNewContact(name, email string) Contact {
	id++

	return Contact{
		Name:  name,
		Email: email,
		Id:    id,
	}
}

func (d *ContactData) checkContactExists(email string) bool {
	for _, contact := range d.Contacts {
		if contact.Email == email {
			return true
		}
	}

	return false
}

func (d *ContactData) indexOf(contactId int) int {
	for i, contact := range d.Contacts {
		if contact.Id == contactId {
			return i
		}
	}

	return -1
}

func main() {
	e := echo.New()

	pageData := createDefaultPageData()

	e.Renderer = newTemplate()
	e.Use(middleware.Logger())

	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index", pageData)
	})

	e.POST("/contacts", func(c echo.Context) error {
		name, email := c.FormValue("name"), c.FormValue("email")

		if pageData.Data.checkContactExists(email) {
			formData := FormData{
				Values: map[string]string{
					"name":  name,
					"email": email,
				},
				Errors: map[string]string{
					"email": "Email already exists",
				},
			}

			return c.Render(422, "form", formData)
		}

		newContact := createNewContact(name, email)
		pageData.Data.Contacts = append(pageData.Data.Contacts, newContact)

		fmt.Println(pageData.Data)

		// render a clean form
		err := c.Render(200, "form", createEmptyFormData())

		if err != nil {
			return err
		}

		return c.Render(200, "contact-list-oob", newContact)
	})

	e.DELETE("/contacts/:id", func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			return c.String(400, "Invalid id")
		}

		// confirm this contact exists
		deletedIndex := pageData.Data.indexOf(id)

		if deletedIndex == -1 {
			return c.String(404, "Contact not found")
		}

		pageData.Data.Contacts = append(pageData.Data.Contacts[:deletedIndex], pageData.Data.Contacts[deletedIndex+1:]...)

		return c.NoContent(200)
	})

	e.Logger.Fatal(e.Start(":3000"))
}
