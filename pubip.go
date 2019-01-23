package main

import (
	"log"
	"net/smtp"

	externalip "github.com/glendc/go-external-ip"
	cfgutils "github.com/mbarbita/golib-cfgutils"
	storestruct "github.com/mbarbita/golib-storestruct"
)

// PubIP is a struct to work with current ip and store old ip
type PubIP struct {
	IP string
}

func main() {

	var pubip = new(PubIP)

	// Create the default consensus,
	// using the default configuration and no logger.
	consensus := externalip.DefaultConsensus(nil, nil)

	// Get your IP,
	// which is never <nil> when err is <nil>.
	ip, err := consensus.ExternalIP()
	if err == nil {
		// fmt.Println(ip.String()) // print IPv4/IPv6 in string format
		pubip.IP = ip.String()
		log.Println("Actual public IP:", pubip.IP)

		// load old ip
		var oldpubip = new(PubIP)
		if err = storestruct.Load("pubip.txt", oldpubip); err != nil {
			log.Fatalln(err)
		}

		log.Println("Old public IP:   ", oldpubip.IP)
		if pubip.IP == oldpubip.IP {
			log.Println("Same IP, quitting.")
			return
		}

	}
	log.Println("Different IP, sending email:")

	//Save actual ip
	if err = storestruct.Save("pubip.txt", pubip); err != nil {
		log.Fatalln(err)
	}

	cfgMap := cfgutils.ReadCfgFile("cfg.ini")

	// Set up authentication information.
	auth := smtp.PlainAuth("",
		cfgMap["username"],
		cfgMap["pass"],
		cfgMap["hostname"])

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	msg := []byte("Subject: Public IP\r\n" +
		pubip.IP + "\r\n")

	err = smtp.SendMail(cfgMap["hostname"]+cfgMap["port"],
		auth,
		cfgMap["from"],
		[]string{cfgMap["to"]},
		msg)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("email sent!")

}
