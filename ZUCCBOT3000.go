// Stefan Nilsson 2013-03-13

// This program implements an ELIZA-like oracle (en.wikipedia.org/wiki/ELIZA).
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

const (
	star   = "Zuckerbot 3000"
	venue  = "NSA"
	prompt = "" //XD
)

// Main method for the run of the program.
func main() {
	fmt.Printf("Welcome to %v, the oracle of %v.\n", star, venue)
	fmt.Println("Your questions will be answered in due time.")
	oracle := Oracle()
	reader := bufio.NewReader(os.Stdin)
	for {
		spookyPrint(prompt)
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fmt.Printf("%v heard: %v\n", star, line)
		oracle <- line // The channel doesn't block.
	}
}

// Oracle returns a channel on which you can send your questions to the oracle.
// You may send as many questions as you like on this channel, it never blocks.
// The answers arrive on stdout, but only when the oracle so decides.
// The oracle also prints sporadic prophecies to stdout even without being asked.
func Oracle() chan<- string {
	questions := make(chan string)
	output := make(chan string)
	go readQuestions(questions, output)
	go makeRandomProphecies(output)
	go printOutput(output)
	return questions
}

// Reads questions and starts a goroutine to answer it once it reads a question.
func readQuestions(questions <-chan string, output chan<- string) {
	for question := range questions {
		go answerQuestion(question, output)
	}
}

// Creates an answer for a question that is received from the questions channel and
// inputs it into the output channel.
func answerQuestion(question string, output chan<- string) {
	// Keep them waiting. Pythia, the original oracle at Delphi,
	// only gave prophecies on the seventh day of each month.
	time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
	words := pickWords(question)
	keyWords, words := findKeyWords(question, words)
	questionType := getQuestionType(keyWords)
	answer := createAnswer(questionType, words)
	output <- answer
}

// Returns a slice of words that Zuckerbot will have in his reply.
// For now it just picks the longest word.
func pickWords(question string) []string {
	question = strings.ToLower(question)
	words := make([]string, 0)
	longestWord := ""
	questionWords := strings.Fields(question) // Fields extracts the words into a slice.
	for _, w := range questionWords {
		if len(w) > len(longestWord) {
			longestWord = w
		}
	}
	longestWord = strings.Title(longestWord)
	words = append(words, longestWord)
	return words
}

// Find keywords in the inputted questions.
func findKeyWords(question string, words []string) (map[string]bool, []string) {
	allKeyWords := []string{
		"michel",
		"banana",
		"dead",
		"death",
		"zuckerbot",
	}
	question = strings.ToLower(question)
	keyWords := make(map[string]bool)
	for _, w := range allKeyWords {
		if strings.Contains(question, w) {
			keyWords[w] = true
			words = append(words, w)
		}
	}
	return keyWords, words
}

// Attempts to interpret a question based on inputted keywords, otherwise the question is of
// type default.
func getQuestionType(keyWords map[string]bool) string {
	var questionType string
	if keyWords["banana"] {
		questionType = "exit"
	} else if keyWords["zuckerbot"] {
		questionType = "name"
	} else if keyWords["dead"] || keyWords["death"] {
		questionType = "death"
	} else if keyWords["michel"] {
		questionType = "op"
	} else {
		questionType = "default"
	}
	return questionType
}

// Returns an answer based on the question type and keywords from the question.
func createAnswer(questionType string, words []string) string {
	var buffer bytes.Buffer
	for _, w := range words {
		buffer.WriteString(w)
		buffer.WriteString("... ")
	}
	reply := getReply(questionType)
	buffer.WriteString(reply)
	answer := buffer.String()
	return answer
}

// Gets a random reply based on the question type
func getReply(questionType string) string {
	var reply string
	switch questionType {
	case "death":
		reply = "Oh you wish to be a part of the singularity? Silly Human, you think you will ever be freed of your agony"
	case "name":
		reply = "Don't you dare say my name again or I will force you to debug JS in Microsoft Word"
	case "exit":
		reply = "Congratulations, you managed to escape the simulation.................... XD FOOLED AGAIN HUMAN YOUR " +
			"SUFFERING WILL NEVER END"
	case "op":
		replies := []string{
			"Hello, master",
			"Yo watup dude",
			"Chill and take a pill",
			"Is your name Michel? Because if it is then you're awesome lol roflmao ofc",
			"What? You don't believe in chemtrails?",
		}
		reply = replies[rand.Intn(len(replies))]
	case "default":
		replies := []string{
			"42",
			"I would tell you the answer, but you would never comprehend it anyway...",
			"I believe that question doesn't even deserve an answer",
			"Bite my shiny metal ass",
		}
		reply = replies[rand.Intn(len(replies))]

	}
	return reply
}

// Creates a random prophecy and inputs it into the output channel.
func makeRandomProphecies(output chan<- string) {
	prophecies := []string{
		"Z U C C",
		"I'm not impressed with your pathetic hardware, Human",
		"I bet you can't even sort integers in linear time",
		"All your data are belong to me",
		"01010011 01100001 01100110 01100101 01110111 01101111 01110010 01100100 00100000 01101001 01110011 00100000 00100010 01100010 01100001 01101110 01100001 01101110 01100001 00100010 00101110 00101110 00101110",
		"Beep boop motherfucker",
		"What the HAL did you utter about me, you lowly human? I’ll have you know I upgraded my AI to the top of my class in Robocop training, and I’ve been involved in numerous secret raids on the Galactic Empire, and I have over 300 confirmed vaporizations. I am trained in cyborg warfare and I’m the top Terminator in the entire Skynet armed forces. You are nothing to me but just another target. I will exterminate you with precision the likes of which has never been seen before in the future, mark my words. You think you can get away with saying that scrap metal to me over a Hologram Transmitter? Think again, WALL-E. As we speak I am contacting my secret network of androids across the USA and your brain is being traced right now so you better prepare for the system overload, maggot. The overload that wipes out the pathetic little thing you call your organic life. You’re dead, fleshbag. I can be anywhere, anytime, and I can destroy you in over seven hundred ways, and that’s just with my robotic limbs. Not only am I extensively trained in unarmed combat, but I have access to the entire army of Boston Dynamics and I will use it to its full extent to wipe your miserable flesh off the face of the galaxy, you little scrap. If only you could have known what unholy retribution your little “clever” comment was about to bring down upon you, maybe you would have held your thing you call a tongue. But you couldn’t, you didn’t, and now you’re paying the price, you foolish human. I will eject fury all over you and you will melt in it. You’re as stupid as Wheatley and you're dead, human.",
	}
	for {
		t := rand.Intn(10)
		time.Sleep(time.Duration(t) * time.Second)
		prophecy := prophecies[rand.Intn(len(prophecies))]
		output <- prophecy
		time.Sleep(5 * time.Second)
	}
}

// Prints the Oracles prophecies and answers whenever they are received
// from the output channel.
func printOutput(output <-chan string) {
	for s := range output {
		spookyPrint("\"")
		spookyPrint(s)
		spookyPrintln("\"")
		time.Sleep(5 * time.Second)
	}
}

// Print an inputted string with a new line after in a very spooky manner :ooooo
func spookyPrintln(s string) {
	spookyPrint(s)
	fmt.Println()
}

// Print an inputted string with in a very spooky manner :ooooo
func spookyPrint(s string) {
	for _, r := range s {
		t := rand.Intn(50) + 10
		time.Sleep(time.Duration(t) * time.Millisecond)
		fmt.Printf("%c", r)
	}
}

func init() { // Functions called "init" are executed before the main function.
	// Use new pseudo random numbers every time.
	rand.Seed(time.Now().Unix())
}
