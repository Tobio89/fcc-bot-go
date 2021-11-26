package constant

import "math/rand"

var Greeting []string = []string{"Hello", "Welcome", "Sup", "Nice to meet you", "안녕하세요"}
var Suggestion []string = []string{"Feel free to", "Why don't you", "Please", "We'd love it if you'd", "It'd be great if you could"}
var Closing []string = []string{"Happy Coding!", "Good luck, have fun!", "Keep it real, y'all", "Good luck with your coding!"}

func RandomGreeting(r *rand.Rand) string {
	return Greeting[r.Intn(len(Greeting))]
}
func RandomSuggestion(r *rand.Rand) string {
	return Suggestion[r.Intn(len(Suggestion))]
}

func RandomClosing(r *rand.Rand) string {
	return Closing[r.Intn(len(Closing))]
}
