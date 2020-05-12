package web

import (

	"fmt"
	"log"
)

func createUser(username string, password, email string) bool {
	preparedStatement, err := mysqlConn.Prepare("INSERT INTO user_info VALUES(NULL, ?, ?, ?)")

	if err != nil {
		fmt.Println(err)
		return false
	}
	defer preparedStatement.Close()
	//Thanks bcrypt
	_, err = preparedStatement.Exec(username, hashPassword(password), email)

	if err != nil {
		fmt.Println(err)
		return false
	}

	fmt.Println("Created user.")

	return true
}

func userExists(username string) bool {
	var tmp string
	err := mysqlConn.QueryRow("SELECT `username` FROM user_info WHERE `username` = ?", username).Scan(&tmp)
	return err == nil
}

func emailExists(email string) bool {
	var tmp string
	err := mysqlConn.QueryRow("SELECT `email` FROM user_info WHERE `email` = ?", email).Scan(&tmp)
	return err == nil
}
func getEmail(username string) string {
	var email string
	err := mysqlConn.QueryRow("SELECT `email` FROM user_info WHERE `username` = ?", username).Scan(&email)
	if err != nil {
		return ""
	}

	return email
}

func getUsername(email string) string {
	var username string
	err := mysqlConn.QueryRow("SELECT `username` FROM user_info WHERE `email` = ?", email).Scan(&username)
	if err != nil {
		return ""
	}

	return username
}

func getUserID(username string) int {
	var userID int
	err := mysqlConn.QueryRow("SELECT `id` FROM user_info WHERE `username` = ?", username).Scan(&userID)
	if err != nil {
		log.Println(err)
		return -1
	}
	return userID
}

func CheckUserPassword(username, password string) bool {
	var dbPwd string
	err := mysqlConn.QueryRow("SELECT `password` FROM user_info WHERE `username` = ?", username).Scan(&dbPwd)
	if err != nil {
		return false
	}
	return checkPassword(password, dbPwd)
}

func createCenter(centerName string, channelID string, guildID string, ggLink string) bool {
	preparedStatement, err := mysqlConn.Prepare("INSERT INTO centers" +
		"(center_name, discord_guild_id, discord_channel_id, gg_leap_link, twitter_access_token, twitter_secret, facebook_page_name, facebook_page_id, facebook_page_authtoken)" +
		"VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer preparedStatement.Close()
	_, err = preparedStatement.Exec(centerName, channelID, guildID, ggLink, "", "", "", "", "")

	if err != nil {
		fmt.Println(err)
		return false
	}

	fmt.Println("Created center.")
	return true
}

func updateCenter(centerName string, channelID string, guildID string, ggLink string, id int) bool {
	preparedStatement, err := mysqlConn.Prepare("UPDATE centers " +
		"SET center_name = ?, discord_channel_id = ?, discord_guild_id =?, gg_leap_link = ? " +
		"WHERE id = ?;")
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer preparedStatement.Close()
	_, err = preparedStatement.Exec(centerName, channelID, guildID, ggLink, id)

	if err != nil {
		fmt.Println(err)
		return false
	}

	fmt.Println("Updated center.")
	return true
}

func insertNewSchedule(timetoPost string, dayOfWeek string, game string, centerId int) bool {
	preparedStatement, err := mysqlConn.Prepare("INSERT INTO schedules" +
		"(time_to_post, day_of_week, game, center_id)" +
		"VALUES(?, ?, ?, ?)")
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer preparedStatement.Close()
	_, err = preparedStatement.Exec(timetoPost, dayOfWeek, game, centerId)

	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func updateSchedule(timetoPost string, dayOfWeek string, game string, id int) bool {
	preparedStatement, err := mysqlConn.Prepare("UPDATE schedules " +
		"SET time_to_post = ?, day_of_week = ?, game =? " +
		"WHERE id = ?;")
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer preparedStatement.Close()
	_, err = preparedStatement.Exec(timetoPost, dayOfWeek, game, id)

	if err != nil {
		fmt.Println(err)
		return false
	}

	fmt.Println("Updated schedule.")
	return true
}

func updateTwitterOAUTH(accessToken, accessSecret string, id int) bool {
	preparedStatement, err := mysqlConn.Prepare("UPDATE centers " +
		"SET twitter_access_token = ?, twitter_secret = ? " +
		"WHERE id = ?;")
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer preparedStatement.Close()
	_, err = preparedStatement.Exec(accessToken, accessSecret, id)

	if err != nil {
		fmt.Println(err)
		return false
	}

	fmt.Println("Updated twitter oauth tokens")
	return true
}


func updateFacebookOAUTH(facebook_page_name, facebook_page_id, facebook_page_authtoken string, id int) bool {
	preparedStatement, err := mysqlConn.Prepare("UPDATE centers " +
		"SET facebook_page_name = ?, facebook_page_id = ?,  facebook_page_authtoken = ? " +
		"WHERE id = ?;")
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer preparedStatement.Close()
	_, err = preparedStatement.Exec(facebook_page_name, facebook_page_id, facebook_page_authtoken, id)

	if err != nil {
		fmt.Println(err)
		return false
	}

	fmt.Println("Updated facebook oauth tokens")
	return true
}

func deleteSchedules(centerID int) bool {
	preparedStatement, err := mysqlConn.Prepare("DELETE FROM schedules where center_id = ?")
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer preparedStatement.Close()
	_, err = preparedStatement.Exec(centerID)

	if err != nil {
		fmt.Println(err)
		return false
	}

	fmt.Println("Cleared schedules for: " + string(centerID))
	return true
}

type ScheduleInfo struct {
	ID         int
	TimeToPost string
	DayOfWeek  string
	Game       string
	CenterID   int
}

func getAllSchedules() []ScheduleInfo {
	rows, err := mysqlConn.Query("SELECT * FROM schedules")
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()

	scheduleInfo := make([]ScheduleInfo, 0)

	for rows.Next() {
		var (
			id           int
			time_to_post string
			day_of_week  string
			game         string
			center_id    int
		)
		if err := rows.Scan(&id, &time_to_post, &day_of_week, &game, &center_id); err != nil {
			log.Println(err)
			return nil
		}

		scheduleInfo = append(scheduleInfo, ScheduleInfo{ID: id, TimeToPost: time_to_post, DayOfWeek: day_of_week, Game: game, CenterID: center_id})
	}
	return scheduleInfo
}

func getAllSchedulesForCenter(centerID int) []ScheduleInfo {
	rows, err := mysqlConn.Query("SELECT * FROM schedules WHERE center_id = ?", centerID)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()

	scheduleInfo := make([]ScheduleInfo, 0)

	for rows.Next() {
		var (
			id           int
			time_to_post string
			day_of_week  string
			game         string
			center_id    int
		)
		if err := rows.Scan(&id, &time_to_post, &day_of_week, &game, &center_id); err != nil {
			log.Println(err)
			return nil
		}

		scheduleInfo = append(scheduleInfo, ScheduleInfo{ID: id, TimeToPost: time_to_post, DayOfWeek: day_of_week, Game: game, CenterID: center_id})
	}
	return scheduleInfo
}



type CenterInfo struct {
	CenterName string
	TwitterAccessToken string
	TwitterSecret     string
	FacebookPageName       string
	FacebookPageID       string
	FacebookPageAuthtoken       string
	DiscordChannelID string
	DiscordGuildID string
	GGLeapLink string
}

type CenterWithID struct {
	CenterID int   `json:"centerID"`
	CenterName string  `json:"centerName"`
	DiscordChannelID string `json:"discordChannelID"`
	DiscordGuildID string `json:"discordGuildID"`
	GGLeapLink string`json:"ggLeapLink"`
	Schedules []ScheduleInfo
}

func getAllCenters() []CenterWithID {
	rows, err := mysqlConn.Query("SELECT id, center_name, discord_channel_id, discord_guild_id, gg_leap_link FROM centers")
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()

	centers := make([]CenterWithID, 0)

	for rows.Next() {
		var (
			id           int
			center_name string
			discord_channel_id string
			discord_guild_id string
			gg_leap_link string
		)

		if err := rows.Scan(&id, &center_name, &discord_channel_id, &discord_guild_id, &gg_leap_link); err != nil {
			log.Println(err)
			return nil
		}
		schedules := getAllSchedulesForCenter(id)
		centers = append(centers, CenterWithID{CenterID: id, CenterName: center_name, DiscordChannelID:
			discord_channel_id, DiscordGuildID: discord_guild_id,GGLeapLink: gg_leap_link, Schedules: schedules})
	}

	return centers
}

func getCenterInfo(centerId int) CenterInfo {
	var center CenterInfo
	err := mysqlConn.QueryRow("SELECT center_name, twitter_access_token, twitter_secret, facebook_page_name, facebook_page_id, facebook_page_authtoken, discord_channel_id, discord_guild_id, gg_leap_link FROM centers WHERE id = ?", centerId).
		Scan(&center.CenterName, &center.TwitterAccessToken, &center.TwitterSecret, &center.FacebookPageName, &center.FacebookPageID, &center.FacebookPageAuthtoken,  &center.DiscordChannelID, &center.DiscordGuildID, &center.GGLeapLink)
	if err != nil {
		log.Println(err)
		return center
	}

	return center
}
