package text

import "image/color"

type Message struct {
	FirstLine string
	SecondLine string // use "" if empty
	Color color.RGBA
	IsDialogue bool
	IsSkippable bool
}

func (self *Message) HasTwoLines() bool {
	return self.SecondLine != ""
}

func NewMsg1(line string, clr color.RGBA) *Message {
	return &Message{
		FirstLine: line,
		SecondLine: "",
		Color: clr,
		IsDialogue: false,
		IsSkippable: false,
	}
}

func NewMsg2(line1, line2 string, clr color.RGBA) *Message {
	return &Message{
		FirstLine: line1,
		SecondLine: line2,
		Color: clr,
		IsDialogue: false,
		IsSkippable: false,
	}
}

// One line of dialogue message.
func NewDialogueMsg1(line string, clr color.RGBA) *Message {
	return &Message{
		FirstLine: line,
		SecondLine: "",
		Color: clr,
		IsDialogue: true,
		IsSkippable: true,
	}
}

func NewDialogueMsg2(line1, line2 string, clr color.RGBA) *Message {
	return &Message{
		FirstLine: line1,
		SecondLine: line2,
		Color: clr,
		IsDialogue: true,
		IsSkippable: true,
	}
}

func NewSkippableMsg1(line string, clr color.RGBA) *Message {
	return &Message{
		FirstLine: line,
		SecondLine: "",
		Color: clr,
		IsDialogue: false,
		IsSkippable: true,
	}
}

func NewSkippableMsg2(line1, line2 string, clr color.RGBA) *Message {
	return &Message{
		FirstLine: line1,
		SecondLine: line2,
		Color: clr,
		IsDialogue: false,
		IsSkippable: true,
	}
}
