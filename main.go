
package main

import (
    "encoding/csv"
    "fmt"
    "log"
    "net/http"
    "os"
    "strings"
    "strconv"
)

func main() {
    /*if len(os.Args) != 2 {
        fmt.Fprintf(os.Stderr, "usage: mybot slack-bot-token\n")
        os.Exit(1)
    }*/


    // var err error
    hubotToken := os.Getenv("HUBOT_SLACK_TOKEN")

    if hubotToken == "" {
        log.Fatal("$HUBOT_SLACK_TOKEN must be set")
    }

    // start a websocket-based Real Time API session
    ws, id := slackConnect(hubotToken)

	development := os.Getenv("DEVELOPMENT_NAME")
	prefix:=""
	if development != "" {
		id = development
		log.Print("DEV mode: Bot name changed to " + id)
		prefix="@"+id
	}else
	{
		prefix="<@"+id+">"
	}

	fmt.Println(id + " ready, ^C exits")


    for {
        // read each incoming message
        m, err := getMessage(ws)
        if err != nil {
            log.Fatal(err)
        }

        // see if we're mentioned
        if m.Type == "message" && strings.HasPrefix(m.Text, prefix) {
            // if so try to parse if
            parts := strings.Fields(m.Text)
            if len(parts) == 3 && parts[1] == "stock" {
                // looks good, get the quote and reply with the result
                go func(m Message) {
                    m.Text = getQuote(parts[2])
                    postMessage(ws, m)
                }(m)
                // NOTE: the Message object is copied, this is intentional
            } else if len(parts) == 3 && parts[1] == "create" {
		    // looks good, get the quote and reply with the result
		    createInstance(parts[2], ws)
		    // NOTE: the Message object is copied, this is intentional
	    } else {
                // huh?
                m.Text = fmt.Sprintf("sorry, that does not compute\n")
                postMessage(ws, m)
            }
        }
    }
}

// Get the quote via Yahoo. You should replace this method to something
// relevant to your team!
func getQuote(sym string) string {
    sym = strings.ToUpper(sym)
    url := fmt.Sprintf("http://download.finance.yahoo.com/d/quotes.csv?s=%s&f=nsl1op&e=.csv", sym)
    resp, err := http.Get(url)
    if err != nil {
        return fmt.Sprintf("error: %v", err)
    }
    rows, err := csv.NewReader(resp.Body).ReadAll()
    if err != nil {
        return fmt.Sprintf("error: %v", err)
    }
    if len(rows) >= 1 && len(rows[0]) == 5 {
	    //var previousClose int :=rows[0][4]
	    previousClose, err := strconv.ParseFloat(rows[0][4], 64)
	    if err != nil {
		    return fmt.Sprintf("error: %v", err)
	    }
	    open, err := strconv.ParseFloat(rows[0][3], 64)
	    if err != nil {
		    return fmt.Sprintf("error: %v", err)
	    }
	    nowValue, err := strconv.ParseFloat(rows[0][2], 64)
	    if err != nil {
		    return fmt.Sprintf("error: %v", err)
	    }
	    previousClose = nowValue - previousClose

	    open = nowValue - open

        return fmt.Sprintf("%s (%s) is trading at $%.2f. Since open: %.2f, Since yesterday: %.2f", rows[0][0], rows[0][1], nowValue, open, previousClose)
    }
    return fmt.Sprintf("unknown response format (symbol was \"%s\")", sym)
}