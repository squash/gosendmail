package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/mail"
	"net/smtp"
	"os"
	"os/user"
	"github.com/spf13/viper"
	"github.com/spf13/pflag"
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
	err=viper.ReadInConfig()
	if err != nil {
			log.Fatal("Config file not found.")
	}
	var Log *log.Logger
	logfile:=viper.GetString("config.logfile")
	if logfile != "" {
		file, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
		    log.Fatalln("Failed to open log file:", err)
		}		

		Log = log.New(file,
		    "Info: ",
		    log.Ldate|log.Ltime|log.Lshortfile)

	} else {
			Log=log.New(os.Stderr, "Info:", log.Ldate|log.Ltime|log.Lshortfile)
			
	}
	
	var recip []string

    smtpaddr:=viper.GetString("config.smtpaddr")
    fromaddr=viper.GetString("config.fromaddr")
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
		Log.Fatal( "error reading stdin")
	}

	msg, err := mail.ReadMessage(bytes.NewReader(body))
	if err != nil {
		Log.Fatal("error parsing message body")
	}

	if len(recip) == 0 {
		// We only need to parse the message to get a recipient if none where
		// provided on the command line.
		tmp,err:=mail.ParseAddress(msg.Header.Get("To"))
		if err != nil {
			Log.Fatal("No recipient specified")
		}
		recip = append(recip, tmp.String())
	}

	Log.Println("Starting: From ", fromaddr, " To ", recip)

	err = smtp.SendMail(smtpaddr, nil, fromaddr, recip, body)
	if err != nil {
		Log.Fatal("Error sending mail: ", err)
	}
	Log.Println("No errors.")

}
