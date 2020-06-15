package web

import (
	"contender/discord"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/carlescere/scheduler"
	"github.com/dlion/goImgur"
)

var scheduledJobs []*scheduler.Job

func setupScheduler() {

	log.Println("Reloading scheduler...")

	allSchedules := getAllSchedules()

	fmt.Println(time.Now().UTC())

	for _, schedule := range allSchedules {
		schedule.TimeToPost = schedule.TimeToPost + ":05"
		fmt.Println("Registering event for: " + schedule.TimeToPost)
		switch schedule.DayOfWeek {
		case "monday":
			job, _ := scheduler.Every(1).Monday().At(schedule.TimeToPost).Run(func() {
				runDistrubition(schedule)
			})
			scheduledJobs = append(scheduledJobs, job)
		case "tuesday":
			job, _ := scheduler.Every(1).Tuesday().At(schedule.TimeToPost).Run(func() {
				runDistrubition(schedule)
			})
			scheduledJobs = append(scheduledJobs, job)
		case "wednesday":
			job, _ := scheduler.Every(1).Wednesday().At(schedule.TimeToPost).Run(func() {
				runDistrubition(schedule)
			})
			scheduledJobs = append(scheduledJobs, job)
		case "thursday":
			job, _ := scheduler.Every(1).Thursday().At(schedule.TimeToPost).Run(func() {
				runDistrubition(schedule)
			})
			scheduledJobs = append(scheduledJobs, job)
		case "friday":
			job, _ := scheduler.Every(1).Friday().At(schedule.TimeToPost).Run(func() {
				runDistrubition(schedule)
			})
			scheduledJobs = append(scheduledJobs, job)
		case "saturday":
			job, _ := scheduler.Every(1).Saturday().At(schedule.TimeToPost).Run(func() {
				runDistrubition(schedule)
			})
			scheduledJobs = append(scheduledJobs, job)
		case "sunday":
			job, _ := scheduler.Every(1).Sunday().At(schedule.TimeToPost).Run(func() {
				runDistrubition(schedule)
			})
			scheduledJobs = append(scheduledJobs, job)
		}
	}
}

func clearScheduler() {
	for _, val := range scheduledJobs {
		val.Quit <- true
	}
}

func runDistrubition(info ScheduleInfo) {

	fmt.Println("Running distribution for: ", info)

	center := getCenterInfo(info.CenterID)

	game := info.Game
	//url := "https://slider.ggleap.com/?center=476fc5ba-6114-4445-94f8-b3734e7f770d&screen=main"
	url := center.GGLeapLink
	filename := game + "-" + strconv.Itoa(info.CenterID) + ".png"

	cmd := exec.Command("node", "screenshot-puller/main.js", game, filename, url)
	_, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	f, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}

	content := "Here are the top players this week in " + game + "! Congrats to these players!"

	ms := &discordgo.MessageSend{
		Content: content,
		Files: []*discordgo.File{
			&discordgo.File{
				Name:   filename,
				Reader: f,
			},
		},
	}

	messageInfo := discord.ChannelMessageInfo{ms, center.DiscordChannelID}

	DiscordChannel <- messageInfo

	str, err := goImgur.Upload(filename, "3e7a4deb7ac67da")
	if err != nil {
		log.Println("Failed uploading image to imgur")
		// do something
	}
	uploadRes := parseImgurResult(*str)
	if uploadRes == nil {
		log.Println("Failed to parse imgur result")
	}
	imgurLink := uploadRes.Data.Link

	//Post to facebook
	err = postToFacebook(center, Post{
		Message:        content,
		AttachmentType: PhotoAttachment,
		Attachment:     imgurLink,
	})
	if err != nil {
		log.Println("Failed to post to facebook:", err)
	}

	//Face to twitter
	err = sendTwitterTweeterTweet(center, content, filename)
	if err != nil {
		log.Println("Failed to post to twitter:", err)
	}

	//TODO: Wait on channel.
	time.Sleep(2 * time.Second)

	f.Close()

	fmt.Println("Finished shit.")
}
