package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

type Person struct {
	Name     string
	Nickname string
	Position string
	Email    string
	Phone    string
}
type AboutMe struct {
	Label string
	Desc  string
}

// Define struct to represent Education section
type Education struct {
	Label string
	Desc  string
}

func main() {
	c := colly.NewCollector()

	var personData Person
	var aboutMeData []AboutMe
	var educationData []Education

	c.OnHTML("div.info", func(info *colly.HTMLElement) {
		info.ForEach("div.name.es, small.es, div.position.es, div.email.es, div.contact.clearfix", func(i int, child *colly.HTMLElement) {
			// Read the text of each child element
			childText := strings.TrimSpace(child.Text)
			switch i {
			case 0:
				personData.Name = childText
			case 1:
				personData.Nickname = childText
			case 2:
				personData.Position = childText
			case 3:
				personData.Email = childText
			case 4:
				personData.Phone = childText
			}
		})
	})

	c.OnHTML("div.items", func(items *colly.HTMLElement) {
		// Loop through each child item
		items.ForEach("div.item", func(_ int, item *colly.HTMLElement) {
			// Read the label and desc elements of each item
			label := strings.TrimSpace(item.ChildText("div.label"))
			desc := strings.TrimSpace(item.ChildText("div.desc.es"))

			// Determine if it's an "About Me" section or "Education" section
			// and populate the respective struct
			if label == "About Me" {
				aboutMeData = append(aboutMeData, AboutMe{
					Label: label,
					Desc:  desc,
				})
			} else if label == "Education" {
				educationData = append(educationData, Education{
					Label: label,
					Desc:  desc,
				})
			}
		})
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Printf("Visiting %s\n", r.URL)

	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Printf("Visited URL: %s\n", r.Request.URL)
		fmt.Printf("Response Status: %d\n", r.StatusCode)
	})

	c.OnError(func(r *colly.Response, e error) {
		fmt.Printf("Error while scraping: %s\n", e.Error())
	})

	err := c.Visit("https://people.njit.edu/faculty/malek#about")
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create("person.txt")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Write the data to the text file
	fmt.Fprintf(file, "Name: %s\n", personData.Name)
	fmt.Fprintf(file, "Nickname: %s\n", personData.Nickname)
	fmt.Fprintf(file, "Position: %s\n", personData.Position)
	fmt.Fprintf(file, "Email: %s\n", personData.Email)
	fmt.Fprintf(file, "Phone: %s\n", personData.Phone)
}
