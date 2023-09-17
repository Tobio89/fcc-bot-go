package msg

import (
	"math/rand"
	"time"
)

type greeting []string

func (g *greeting) GetRandom() string {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	return (*g)[r.Intn(len((*g)))]
}

var Opening greeting = []string{"Hello", "Welcome", "Sup", "Nice to meet you", "안녕하세요", "반가워요", "Hey there", "Oh hi", "Heyyyyyyyy"}
var Suggestion greeting = []string{"Feel free to", "Why don't you", "Please", "We'd love it if you'd", "It'd be great if you could"}
var Closing greeting = []string{"Happy Coding!", "Good luck, have fun!", "Keep it real, y'all", "Good luck with your coding!"}

const (
	LogOnReady      string = "[Init]          "
	LogNewMember    string = "[Joined]        "
	LogMemberLeave  string = "[Left]          "
	LogRFR          string = "[RFR]           "
	LogVerification string = "[Verified]      "
	LogError        string = "[ERROR]         "
	LogLearning     string = "[Learning]      "
	LogShutdown     string = "[Shutdown]      "
)

const (
	CommandErase       string = "<Erase>         "
	CommandForceLog    string = "<Force Log>     "
	CommandRemind      string = "<Remind>        "
	CommandVerify      string = "<Verify>        "
	CommandDeVerify    string = "<DeVerify>      "
	CommandClearIntros string = "<Clear Intros>  "
)
