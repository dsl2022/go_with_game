package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type choice struct {
	cmd         string
	description string
	nextNode    *storyNode
}

type storyNode struct {
	text    string
	choices []*choice
}

func (node *storyNode) addChoice(cmd string, descrption string, nextNode *storyNode) {
	choice := &choice{cmd, descrption, nextNode}
	node.choices = append(node.choices, choice)
	fmt.Println(node, "node test")
}

func (node *storyNode) render() {
	fmt.Println(node.text)
	if node.choices != nil {
		for _, choice := range node.choices {
			fmt.Println(choice.cmd, choice.description)
		}
	}
}

func (node *storyNode) executeCmd(cmd string) *storyNode {
	for _, choice := range node.choices {
		if strings.ToLower(choice.cmd) == strings.ToLower(cmd) {
			return choice.nextNode
		}
	}
	return node
}

var scanner *bufio.Scanner

func (node *storyNode) play() {
	node.render()
	if node.choices != nil {
		scanner.Scan()
		node.executeCmd(scanner.Text()).play()
	}
}

func main() {
	scanner = bufio.NewScanner(os.Stdin)
	start := storyNode{text: `
	You are in large chamber, deep underground, You see three passages leading out.
	 A north passage leads into darkness. To thee south, a passage appears to head upward. 
	 the eastern passages appears flat and well travelled. 
	`}

	darkRoom := storyNode{text: `It is pitch black, You cannot see a thing`}

	darkRoomLit := storyNode{text: `The dark passage is now lit by your lantern.
	Youu can continue north or head bac south`}

	grue := storyNode{text: `While stumbling around in the darkness, you are eaten by a grue`}

	trap := storyNode{text: `You head down the well travelled path when suddenly a trap door opens and you fall into a pit`}

	treasure := storyNode{text: `You arrive at a small chamber, filled with treasure`}

	start.addChoice("N", "Go North", &darkRoom)
	start.addChoice("S", "Go South", &darkRoom)
	start.addChoice("E", "Go East", &trap)

	darkRoom.addChoice("S", "Go South", &grue)
	darkRoom.addChoice("O", "Turn on lantern", &darkRoomLit)

	darkRoomLit.addChoice("N", "Go North", &treasure)
	darkRoomLit.addChoice("S", "Go South", &start)

	start.play()
	fmt.Println(start, "test start")
	fmt.Println("The End.")
}
