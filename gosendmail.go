package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/mail"
	"net/smtp"
	"os"
	"os/user"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Go runs the MailHog sendmail replacement.
func main() {

	host, err := os.Hostname()
	if err != nil {
		host = "localhost"
	}

	username := "nobody"
	user, err := user.Current()
	if err == nil && user != nil && len(user.Username) > 0 {
		username = user.Username
	}

	fromaddr := username + "@" + host
	viper.SetDefault("config.fromaddr", fromaddr)
	viper.SetDefault("config.smtpaddr", "localhost:25")
	viper.SetDefault("config.logfile", "")
	viper.SetConfigName("gosendmail")
	viper.AddConfigPath("/etc")
	err = viper.ReadInConfig()
	if err != nil {
		log.Fatal("Error reading config file: " + err.Error())
	}
	var Log *log.Logger
	logfile := viper.GetString("config.logfile")
	if logfile != "" {
		file, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalln("Failed to open log file:", err)
		}

		Log = log.New(file,
			"Info: ",
			log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		Log = log.New(os.Stderr, "Info:", log.Ldate|log.Ltime|log.Lshortfile)
	}

	var recip []string

	smtpaddr := viper.GetString("config.smtpaddr")
	fromaddr = viper.GetString("config.fromaddr")
	// override defaults from cli flags
	pflag.StringVar(&smtpaddr, "smtp-addr", smtpaddr, "SMTP server address")
	pflag.StringVarP(&fromaddr, "from", "f", fromaddr, "SMTP sender")
	pflag.BoolP("long-i", "i", true, "Ignored. This flag exists for sendmail compatibility.")
	pflag.BoolP("long-t", "t", true, "Ignored. This flag exists for sendmail compatibility.")
	pflag.Parse()

	// allow recipient to be passed as an argument
	recip = pflag.Args()

	body, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		Log.Fatal("error reading stdin")
	}

	msg, err := mail.ReadMessage(bytes.NewReader(body))
	if err != nil {
		Log.Fatal("error parsing message body: " + err.Error())
	}

	if len(recip) == 0 {
		// We only need to parse the message to get a recipient if none where
		// provided on the command line.
		addresslist:=msg.Header.Get("To")
		addresses := strings.Split(addresslist, ",")
		for _, address := range addresses {
			if (!strings.HasPrefix(address, "<")) {
				address=fmt.Sprintf("<%s>", address)
			}
			address=strings.Replace(address, " ", "", -1)
			tmp, err := mail.ParseAddress(address)
			if err != nil {
				Log.Fatalf("Recipient missing or invalid: [%s]: %s", address, err.Error())
			}
			recip = append(recip, tmp.Address)
		}
	} else {
		for i, _ := range recip {

			tmp, err := mail.ParseAddress(recip[i])
			if err != nil {
				Log.Fatal("Invalid recipient specified: " + err.Error())
			}
			recip[i] = tmp.Address
		}
	}
	// Allow message headers to override default sender

	sender := msg.Header.Get("From")
	if sender != "" {
		tmp, err := mail.ParseAddress(sender)
		if err != nil {
			Log.Println("From header found in message but unable to parse. Using default.")
		} else {
			fromaddr = tmp.Address
		}
	}
	Log.Println("Ready to send email from ", fromaddr, " to ", recip)

	err = smtp.SendMail(smtpaddr, nil, fromaddr, recip, body)
	if err != nil {
		Log.Fatal("Error sending mail: ", err)
	}
	Log.Println("No errors.")

}
