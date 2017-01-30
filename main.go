
package main

import (
    "encoding/csv"
    "fmt"
    "log"
    "net/http"
    "os"
    "strings"
    "strconv"
    "time"
)

func main() {

    // var err error
    hubotToken := os.Getenv("C3PO_SLACK_TOKEN")

    if hubotToken == "" {
        log.Fatal("$C3PO_SLACK_TOKEN must be set")
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
		if len(parts) == 2 && parts[1] == "help" {
			go func(m Message) {
				helpText:="'" +prefix + " duty' - shows who is on duty this week \n"
				helpText += "'" +prefix + " duty who' - shows team members (login, name)\n"
				helpText += "'" +prefix + " duty weekNumber morningLogin eveningLogin' - assign persons on duty, complete week\n"
				helpText += "(e.g. '" +prefix +" duty 6 iivanov ppetrov' - assign iivanov to morning shift and ppetrov to evening on week number 6 ) \n"
				helpText += "'" +prefix + " duty yyyy-mm-dd morningLogin eveningLogin' - assign persons on duty, single day\n"
				helpText += "'" +prefix + " week' - shows current week number\n"
				m.Text = fmt.Sprintf(helpText)

				postMessage(ws, m)
			}(m)
		} else if len(parts) == 3 && parts[1] == "stock" {
			// looks good, get the quote and reply with the result
			go func(m Message) {
				m.Text = getQuote(parts[2])
				postMessage(ws, m)
			}(m)
			// NOTE: the Message object is copied, this is intentional
		} else if len(parts) == 2 && parts[1] == "week" {
			// show current week number
			go func(m Message) {
				m.Text = "The current week number is " + getWeekNumber()
				postMessage(ws, m)
			}(m)
		} else if len(parts) == 3 && parts[1] == "duty" && parts[2] == "who" {
			// show team members
			go func(m Message) {
				m.Text = "Members: " + who()
				postMessage(ws, m)
			}(m)

		} else if len(parts) == 2 && parts[1] == "duty" {
			morningText, eveningText:= getDutyString()
			go func(m Message) {
				m.Text = morningText
				postMessage(ws, m)
				m.Text = eveningText
				postMessage(ws, m)
			}(m)
		} else if len(parts) == 5 && parts[1] == "duty" {
			// get all devops team members
			membersList:=getMembersLogin()
			// check if the duty persons is from the devops team
			if stringInSlice(parts[3], membersList) && stringInSlice(parts[4], membersList) {
				// parse attribute week
				weekNum, err := strconv.Atoi(parts[2])
				if err!=nil {
					// parse attribute day
					_, err := time.Parse("2006-01-2", parts[2])
					if err != nil {
						m.Text = fmt.Sprintf("sorry, the second attribute should be a week number or a day (e.g. 2016-02-23) \n")
						postMessage(ws, m)
					} else {
						insertDutyForADay(parts[2], parts[3], parts[4])
						m.Text = fmt.Sprintf("Done. Duty scheduled for "+parts[2]+" day.\n")
						postMessage(ws, m)
					}
				} else {
					insertDutyForAWeek(weekNum, parts[3], parts[4])
					m.Text = fmt.Sprintf("Done. Duty scheduled for "+parts[2]+" week.\n")
					postMessage(ws, m)
				}
			} else {
				m.Text = fmt.Sprintf("sorry, not a team member login \n")
				postMessage(ws, m)
			}
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

func getDutyString() (string, string) {
	var dutyWeek [7]duty = getDutyNow()
	var morningEq bool = true
	var eveningEq bool = true

	for i:=0; i<7 ; i++  {
		if (dutyWeek[i].morningLogin != dutyWeek[0].morningLogin) {
			morningEq = false
		}
		if (dutyWeek[i].eveningLogin != dutyWeek[0].eveningLogin) {
			eveningEq = false
		}
	}
	var morningText string = "DevOps (1 shift 4:00-14:00) : "
	if (morningEq) {
		morningText += dutyWeek[0].morningName
	} else {
		for i:=0; i<7 ; i++  {
			morningText += dutyWeek[0].date +" : "+ dutyWeek[0].morningName
		}
	}

	var eveningText string = "DevOps (2 shift 16:00-2:00) : "
	if (eveningEq) {
		eveningText += dutyWeek[0].eveningName
	} else {
		for i:=0; i<7 ; i++  {
			eveningText += dutyWeek[0].date +" : "+ dutyWeek[0].eveningName
		}
	}
	return morningText, eveningText
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