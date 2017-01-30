package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"log"
)


type duty struct {
	date string
	morningLogin string
	eveningLogin string
	morningName string
	eveningName string
}

func getDutyNow() [7]duty {

	weekArray := GetWeekDays()

	query := "SELECT a.login as MORNING, c.login as EVENING, a.name as MORNING_NAME, c.name as EVENING_NAME, d.dutyday as DDAY from duty d JOIN dashboard.person a ON d.morning_id=a.id JOIN dashboard.person c ON d.evening_id=c.id where d.dutyday in ("

	// todo: check that 7 days return

	for i := 0; i < 7; i++ {
		if (i < 6) {
			query += "'"+ weekArray[i] + "', "
		}else {
			query += "'" + weekArray[i] + "'"
		}
	}

	query += ")"

	mariaJdbc := os.Getenv("MARIA_DASHBOARD_JDBC")

	if mariaJdbc == "" {
		log.Fatal("$MARIA_DASHBOARD_JDBC must be set")
	}

	db, err := sql.Open("mysql", mariaJdbc)
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	var dutyWeek [7]duty

	// Execute the query
	rows, err := db.Query(query)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	// Make a slice for the values
	values := make([]sql.RawBytes, len(columns))

	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	var j int =0;
	// Fetch rows
	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		// Now do something with the data.
		// Here we just print each column as a string.
		var value string
		for i, col := range values {
			// Here we can check if the value is nil (NULL value)
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}

			//fmt.Println(columns[i], ": ", value)
			if (columns[i] == "MORNING") {
				dutyWeek[j].morningLogin = value
			}
			if (columns[i] == "MORNING_NAME") {
				dutyWeek[j].morningName = value
			}
			if (columns[i] == "EVENING") {
				dutyWeek[j].eveningLogin = value
			}
			if (columns[i] == "EVENING_NAME") {
				dutyWeek[j].eveningName = value
			}
			if (columns[i] == "DDAY") {
				dutyWeek[j].date = value
			}


		}
		j++
	}
	if err = rows.Err(); err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	return dutyWeek
}

func insertDutyForAWeek(weekNum int, morningLogin string, eveningLogin string){
	days:=GetWeekDaysForWeekNumber(weekNum)
	var insertQuery string = "insert into duty (dutyday, morning_id, evening_id) values "
	for _, day := range days {
		insertQuery = insertQuery + "('"+day+"', (select M.id from person M where M.login='"+ morningLogin+"'), (select E.id from person E where E.login='"+eveningLogin+"')), "
	}
	insertQuery = TrimSuffix(insertQuery, ", ")
	insertQuery = insertQuery + "on duplicate key update dutyday=VALUES(dutyday), morning_id=VALUES(morning_id), evening_id=VALUES(evening_id)"
	mariaJdbc := os.Getenv("MARIA_DASHBOARD_JDBC")
	if mariaJdbc == "" {
		log.Fatal("$MARIA_DASHBOARD_JDBC must be set")
	}
	db, err := sql.Open("mysql", mariaJdbc)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	resultInsert, err := db.Exec(insertQuery)
	if err != nil {
		panic(err.Error())
		panic(resultInsert)
	}
}

func insertDutyForADay(day string, morningLogin string, eveningLogin string){
	var insertQuery string = "insert into duty (dutyday, morning_id, evening_id) values "
	insertQuery = insertQuery + "('"+day+"', (select M.id from person M where M.login='"+ morningLogin+"'), (select E.id from person E where E.login='"+eveningLogin+"'))"
	insertQuery = insertQuery + "on duplicate key update dutyday=VALUES(dutyday), morning_id=VALUES(morning_id), evening_id=VALUES(evening_id)"
	mariaJdbc := os.Getenv("MARIA_DASHBOARD_JDBC")
	if mariaJdbc == "" {
		log.Fatal("$MARIA_DASHBOARD_JDBC must be set")
	}
	db, err := sql.Open("mysql", mariaJdbc)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	resultInsert, err := db.Exec(insertQuery)
	if err != nil {
		panic(err.Error())
		panic(resultInsert)
	}
}

func getMembersLogin() []string {
	var query string = "select PRS.login as P_LOGIN from person PRS, person_on_stream POS where PRS.id=POS.person_id and POS.stream_id=17 and POS.onboarding=1"

	var result []string


	mariaJdbc := os.Getenv("MARIA_DASHBOARD_JDBC")

	if mariaJdbc == "" {
		log.Fatal("$MARIA_DASHBOARD_JDBC must be set")
	}

	db, err := sql.Open("mysql", mariaJdbc)

	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	// Execute the query
	rows, err := db.Query(query)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	// Make a slice for the values
	values := make([]sql.RawBytes, len(columns))

	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	var j int = 0;
	// Fetch rows
	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		var value string
		for i, col := range values {
			// Here we can check if the value is nil (NULL value)
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			if (columns[i] == "P_LOGIN") {
				result = ExtendStringSlice(result, value)
			}
		}
		j++
	}
	if err = rows.Err(); err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	return result
}


func who() string {
	//show team members
	var query string = "select PRS.login as P_LOGIN, PRS.name as P_NAME from person PRS, person_on_stream POS where PRS.id=POS.person_id and POS.stream_id=17 and POS.onboarding=1"

	var result string


	mariaJdbc := os.Getenv("MARIA_DASHBOARD_JDBC")

	if mariaJdbc == "" {
		log.Fatal("$MARIA_DASHBOARD_JDBC must be set")
	}

	db, err := sql.Open("mysql", mariaJdbc)

	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	// Execute the query
	rows, err := db.Query(query)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	// Make a slice for the values
	values := make([]sql.RawBytes, len(columns))

	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	// Fetch rows
	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}

		// Now do something with the data.
		// Here we just print each column as a string.
		var value string
		for i, col := range values {
			// Here we can check if the value is nil (NULL value)
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			if (columns[i] == "P_LOGIN") {
				result += "Login: " +value
			}
			if (columns[i] == "P_NAME") {
				result += ", Name: " + value
			}

		}
		result+="\n"
	}
	if err = rows.Err(); err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	return result
}