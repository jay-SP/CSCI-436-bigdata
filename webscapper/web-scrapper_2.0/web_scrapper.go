package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/gocolly/colly"
)

func main() {
	// Create a new Collector
	c := colly.NewCollector()

	var hrefValues []string
	var wg sync.WaitGroup // to sync between main and go routine so the writing processes finishes

	baseURL := "https://www.stevens.edu/"

	uniqueURLs := make(map[string]bool)

	c.OnHTML("a", func(e *colly.HTMLElement) { // parent page "https://www.stevens.edu/school-engineering-science/faculty
		// Extract the href attribute value from the anchor tag
		href := e.Attr("href") // contains identifer
		if strings.Contains(href, "/profile/") {
			// Append the base URL to the href value
			fullURL := strings.TrimSuffix(baseURL, "/") + "/" + strings.TrimPrefix(href, "/")

			if !uniqueURLs[fullURL] {
				// Mark the URL as seen by adding it to the uniqueURLs map
				uniqueURLs[fullURL] = true

				// Store the full URL in the hrefValues slice
				hrefValues = append(hrefValues, fullURL)
				// Create a new Collector for each URL and set up the scraping logic
				wg.Add(1)
				go scrapeURL(fullURL, &wg)
			}
		}
	})

	/* 	c.OnRequest(func(r *colly.Request) {
	   		fmt.Printf("Visiting %s\n", r.URL)
	   	})

	   	c.OnResponse(func(r *colly.Response) {
	   		fmt.Printf("Visited URL: %s\n", r.Request.URL)
	   		fmt.Printf("Response Status: %d\n", r.StatusCode)
	   	})
	   	c.OnError(func(r *colly.Response, e error) {
	   		fmt.Printf("Error while scraping: %s\n", e.Error())
	   	}) */

	err := c.Visit("https://www.stevens.edu/school-engineering-science/faculty")
	if err != nil {
		log.Printf("Error visiting")
	}
	/* 	 for _, href := range hrefValues {
		fmt.Println(href)
	} */
	wg.Wait()
	saveDataToile("bios.txt", professorsBios)
	saveDataToFile("courses_taught.txt", professorsCourses)
	saveHrefValuesToFile("bio_urls.txt", hrefValues)
}

var professorsCourses = make(map[string][]string)
var mu sync.Mutex
var professorsBios = make(map[string][]string)

func scrapeURL(url string, wg *sync.WaitGroup) {
	// Create a new Collector
	defer wg.Done()
	c := colly.NewCollector()
	var name string
	var courses []string
	var educations []string
	var research string
	var experience string
	var honors string

	//Targeting courses
	c.OnHTML("section#courses", func(e *colly.HTMLElement) {
		// Finding the div containing course descriptions
		courseDiv := e.ChildText("div.profile-section_f--description__cefgX")

		// Split the course descriptions by line breaks to get individual courses
		courseList := strings.Split(strings.TrimSpace(courseDiv), "\n")

		// Append each course to the courses slice
		for _, course := range courseList {
			courses = append(courses, strings.TrimSpace(course))
		}
	})

	//Targeting edu for bio
	c.OnHTML("section#education", func(e *colly.HTMLElement) {
		// Find the div containing course descriptions
		education := e.ChildText("div.profile-section_f--description__cefgX")
		educations = append(educations, strings.TrimSpace(education))

	})

	//Targeting reasearch for bio
	c.OnHTML("section#research", func(e *colly.HTMLElement) {
		// Find the div containing course descriptions
		research = e.ChildText("div.profile-section_f--description__cefgX")
		fmt.Println(research)

	})

	//Targeting experience for bio
	c.OnHTML("section#experience", func(e *colly.HTMLElement) {
		// Find the div containing course descriptions
		experience = e.ChildText("div.profile-section_f--description__cefgX")
		fmt.Println(experience)

	})

	//Targeting honors for bio
	c.OnHTML("section#honors-and-awards", func(e *colly.HTMLElement) {
		// Find the div containing course descriptions
		honors = e.ChildText("div.profile-section_f--description__cefgX")
		fmt.Println(honors)

	})

	//Targeting prof name
	c.OnHTML("div.hero-profile_f--page-title__Nh2t_", func(e *colly.HTMLElement) {
		// Extract the name from the h1 tag inside the div
		name = strings.TrimSpace(e.ChildText("h1"))
		//fmt.Println("Name:", name)
	})

	/* 	// Set up the callback to be executed before making the request
	   	c.OnRequest(func(r *colly.Request) {
	   		fmt.Printf("Visiting %s\n", r.URL)
	   	})

	   	// Set up the callback to be executed after receiving the response
	   	c.OnResponse(func(r *colly.Response) {
	   		fmt.Printf("Visited URL: %s\n", r.Request.URL)
	   		fmt.Printf("Response Status: %d\n", r.StatusCode)
	   	})

	   	// Set up the callback to be executed when an error occurs
	   	c.OnError(func(r *colly.Response, e error) {
	   		fmt.Printf("Error while scraping: %s\n", e.Error())
	   	}) */

	// Start the scraping process for the given URL
	err := c.Visit(url)
	if err != nil {
		log.Printf("Error visiting %s: %s", url, err)
	}

	mu.Lock()
	professorsCourses[name] = courses
	professorsBios[name] = []string{name, strings.Join(educations, "\n"), research, experience, honors}

	mu.Unlock()
}

func saveDataToFile(fileName string, professorsCourses map[string][]string) {
	file, err := os.Create(fileName)
	if err != nil {
		log.Printf("Error creating file: %s", err)
		return
	}
	defer file.Close()

	// Write the data to the file
	for professor, courses := range professorsCourses {
		fmt.Fprintf(file, "Name: %s\n", professor)
		if len(courses) == 0 {
			fmt.Fprintln(file, "No courses taught")
		} else {
			fmt.Fprintln(file, "Courses:")
			coursesString := strings.Join(courses, "\n")
			fmt.Fprintln(file, coursesString)
		}
		fmt.Fprintln(file, "")
	}

	fmt.Printf("Scraped data saved to %s\n", fileName)

}

func saveDataToile(fileName string, professorsBios map[string][]string) {
	file, err := os.Create(fileName)
	if err != nil {
		log.Printf("Error creating file: %s", err)
		return
	}
	defer file.Close()

	// Write the data to the file
	for _, data := range professorsBios {
		// Extract the name from the data
		name := data[0]
		fmt.Fprintf(file, "Name: %s\n", name)

		// Extract the education from the data
		education := data[1]
		if education != "" {
			fmt.Fprintf(file, "Education:\n%s\n", education)
		} else {
			fmt.Fprintln(file, "Education: No education information available")
		}

		// Extract the research from the data
		research := data[2]
		if research != "" {
			fmt.Fprintf(file, "Research:\n%s\n", research)
		} else {
			fmt.Fprintln(file, "Research: No research information available")
		}

		// Extract the experience from the data
		experience := data[3]
		if experience != "" {
			fmt.Fprintf(file, "Experience:\n%s\n", experience)
		} else {
			fmt.Fprintln(file, "Experience: No experience information available")
		}

		// Extract the honors from the data
		honors := data[4]
		if honors != "" {
			fmt.Fprintf(file, "Honors and Awards:\n%s\n", honors)
		} else {
			fmt.Fprintln(file, "Honors and Awards: No honors and awards information available")
		}

		fmt.Fprintln(file, "") // Add a blank line to separate each professor's data
	}

	fmt.Printf("Data saved to %s\n", fileName)
}

func saveHrefValuesToFile(fileName string, hrefValues []string) {
	file, err := os.Create(fileName)
	if err != nil {
		log.Printf("Error creating file: %s", err)
		return
	}
	defer file.Close()

	// Write the href values to the file
	for _, href := range hrefValues {
		fmt.Fprintln(file, href)
	}

	fmt.Printf("Href values saved to %s\n", fileName)
}
