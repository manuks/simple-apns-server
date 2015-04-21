package main

import (
	"errors"
	"flag"
	"fmt"
	apns "github.com/corneldamian/APNs"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Payload struct {
	message string
	badge   int
}

func (p Payload) ToJson() ([]byte, error) {
	var alert, badge, badgeSep string
	var err error

	if p.message != "" {
		alert = fmt.Sprintf(`"alert" : {"body": "%s"}`, p.message)
		badgeSep = ", "
	}

	if p.badge != 0 {
		badge = fmt.Sprintf(`%s"badge" : %d`, badgeSep, p.badge)
	}

	if p.message == "" && p.badge == 0 {
		err = errors.New("Either message or badge should be present in the payload!")
	}

	aps := fmt.Sprintf(`{"aps": {%s}}`, alert+badge)

	return []byte(aps), err
}

func (p Payload) Config() (expiration uint32, priority uint8) {
	expiration = 300
	priority = 1
	return
}

func main() {
	var cert, key, ip, port string
	var sandbox, debug bool
	//openssl pkcs12 -in Certificates.p12 -out cert.pem -clcerts -nokeys
	//openssl pkcs12 -in Certificates.p12 -out key.pem -nocerts -nodes
	flag.StringVar(&cert, "cert", "", "The certificate file name")
	flag.StringVar(&key, "key", "", "The key file name")
	flag.StringVar(&ip, "ip", "127.0.0.1", "The ip address that it lisents to")
	flag.StringVar(&port, "port", "8080", "The port number")
	flag.BoolVar(&sandbox, "sandbox", false, "Sandbox mode or not")
	flag.BoolVar(&debug, "log", false, "Show more log messages")
	flag.Parse()

	if _, err := os.Stat(cert); os.IsNotExist(err) {
		// path does not exist
		log.Fatal("Could not find certificate file!")
		return
	}

	if _, err := os.Stat(key); os.IsNotExist(err) {
		// path does not exist
		log.Fatal("Could not find certificate key file!")
		return
	}

	if debug {
		logger := log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lshortfile)
		apns.SetGoLogger(logger)
	}

	cfg := apns.Config{
		IsProduction:           !sandbox,
		MaxPoolSize:            10,
		SuccessTimeout:         time.Millisecond * 300,
		NotificationExpiration: time.Duration(-1),
		Certificate:            cert,
		CertificateKey:         key,
	}

	if err := apns.Init(cfg); err != nil {
		log.Fatal("Initialization err: ", err)
		panic(err)
	}

	log.Println("Initialized " + ip + ":" + port)

	http.HandleFunc("/apn", handler)

	err := http.ListenAndServe(ip+":"+port, nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("message: %s, badge: %s, device: %s", r.FormValue("message"), r.FormValue("badge"), r.FormValue("device"))
	p := Payload{}
	p.message = r.FormValue("message")
	badge, _ := strconv.Atoi(r.FormValue("badge"))

	p.badge = badge
	deviceToken := r.FormValue("device")

	id := apns.Send(p, deviceToken)
	pushNotificationStatus := <-apns.Confirm()

	if pushNotificationStatus.Error != nil {
		log.Printf("Error %s", pushNotificationStatus.Error)
		fmt.Fprintf(w, "failed")
	} else {
		log.Printf("Sent %d", id)
		fmt.Fprintf(w, "success")
	}

}
