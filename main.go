package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/smtp"
	"os"
	"strings"

	"github.com/aymerick/raymond"
	"github.com/robfig/cron"
	"github.com/spf13/viper"
)

type IP struct {
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	ISP         string  `json:"isp"`
	Org         string  `json:"org"`
	As          string  `json:"as"`
	Query       string  `json:"query"`
}

func initConfig(cfgFile, defaultCron, defaultEmails string) {
	viper.SetConfigFile(cfgFile)
	viper.AutomaticEnv()

	viper.SetDefault("cron.schedule", defaultCron)
	viper.SetDefault("emails", defaultEmails)

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s\n", err)
		os.Exit(1)
	}

	viper.BindEnv("emails")
	viper.BindEnv("cron.schedule")
}

func getIPInfo() (IP, error) {
	var ip IP
	resp, err := http.Get("http://ip-api.com/json/")
	if err != nil {
		return ip, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ip, err
	}

	err = json.Unmarshal(body, &ip)
	return ip, err
}

func renderTemplate(data IP) (string, error) {
	template := `
<!DOCTYPE html>
<html>
<head>
    <title>IP Information</title>
</head>
<body>
    <h1>IP Information</h1>
    <p><strong>Query:</strong> {{query}}</p>
    <p><strong>ISP:</strong> {{isp}}</p>
    <p><strong>Organization:</strong> {{org}}</p>
    <p><strong>Country:</mark> {{country}} ({{countryCode}})</strong></p>
    <p><strong>Region:</strong> {{regionName}} ({{region}})</p>
    <p><strong>City:</mark> {{city}}, {{zip}}</strong></p>
    <p><strong>Latitude:</strong> {{lat}}</p>
    <p><strong>Longitude:</strong> {{lon}}</p>
    <p><strong>Timezone:</strong> {{timezone}}</p>
    <p><strong>AS:</strong> {{as}}</p>
</body>
</html>
`
	return raymond.Render(template, data)
}

func sendEmail(body, to string) error {
	from := viper.GetString("email.from")
	pass := viper.GetString("email.password")
	host := viper.GetString("smtp.host")
	port := viper.GetString("smtp.port")
	address := host + ":" + port

	auth := smtp.PlainAuth("", from, pass, host)

	msg := []byte("To: " + to + "\r\n" +
		"Subject: Detailed IP Information\r\n" +
		"MIME-version: 1.0;\r\n" +
		"Content-Type: text/html; charset=UTF-8;\r\n" +
		"\r\n" +
		body + "\r\n")

	toList := strings.Split(to, ",")
	err := smtp.SendMail(address, auth, from, toList, msg)
	return err
}

func executeTask() {
	fmt.Println("Fetching IP information and sending emails...")
	emailList := viper.GetString("emails")
	if emailList == "" {
		fmt.Println("No email addresses provided")
		return
	}

	ipInfo, err := getIPInfo()
	if err != nil {
		fmt.Println("Failed to get IP information:", err)
		return
	}

	body, err := renderTemplate(ipInfo)
	if err != nil {
		fmt.Println("Error rendering email template:", err)
		return
	}

	if err := sendEmail(body, emailList); err != nil {
		fmt.Println("Failed to send email:", err)
	} else {
		fmt.Println("Email sent successfully!")
	}
}

func main() {
	configPath := flag.String("config", "config.yaml", "Path to config file")
	cronSpec := flag.String("cron", "0 6,12,18,0 * * *", "Cron specification for how often to fetch IP info and send emails")
	emailList := flag.String("emails", "", "List of email addresses to send IP information")

	flag.Parse()

	initConfig(*configPath, *cronSpec, *emailList)

	executeTask() // run the task immediately for gettting feedback

	c := cron.New()
	err := c.AddFunc(*cronSpec, executeTask)
	if err != nil {
		fmt.Println("Error scheduling the task:", err)
		return
	}

	fmt.Println("Service started successfully. IP information will be fetched and emails sent according to the provided cron specification.")
	c.Start()

	// Block the main thread from exiting to keep the scheduler running.
	select {}
}
