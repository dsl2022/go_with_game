package main

import "fmt"

type storyPage struct {
	text     string
	nextPage *storyPage
}

func (page *storyPage) playStory() {
	for page != nil {
		fmt.Println(page.text)
		page = page.nextPage
	}

	// recursive method
	// if page == nil {
	// 	return
	// }
	// fmt.Println(page.text)
	// page.nextPage.playStory()
}

func (page *storyPage) addToEnd(text string) {
	pageToAdd := &storyPage{text, nil}
	for page.nextPage != nil {
		page = page.nextPage
	}
	page.nextPage = pageToAdd
}

func main() {
	//scanner := bufio.NewScanner(os.Stdin)
	// much simpler
	page1 := storyPage{"It was a dark and storm night.", nil}
	page1.addToEnd("You are alone, you need to find the sacred helmet before the bad guys do")
	page1.addToEnd("You see a troll ahead")

	// page2 := storyPage{"You are alone, you need to find the sacred helmet before the bad guys do", nil}
	// page3 := storyPage{"You see a troll ahead", nil}
	// page1.nextPage = &page2
	// page2.nextPage = &page3
	// method way of calling
	page1.playStory()
	// function
	//playStory(&page1)

	// Functions - has return value - may also excute commands.
	// Procedures - has no return value, just excutes commands.
	// Methods - functions that attached to a struct/object/etc

}
