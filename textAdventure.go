package main

import "fmt"

type choices struct {
	cmd         string
	description string
	nextNode    *storyNode
	nextChoice  *choices
}

type storyNode struct {
	text    string
	choices *choices
}

func (node *storyNode) addChoice(cmd string, descrption string, nextNode *storyNode) {
	choice := &choices{cmd, descrption, nextNode, nil}
	if node.choices == nil {
		node.choices = choice
	} else {
		currentChoice := node.choices
		for currentChoice.nextChoice != nil {
			currentChoice = currentChoice.nextChoice
		}
		currentChoice.nextChoice = choice
	}
}

func (node *storyNode) render() {
	currentChoice := node.choices
	for currentChoice != nil {
		fmt.Println(cmd, ":", description)
	}
}

func main() {

}
