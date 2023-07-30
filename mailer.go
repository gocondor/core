// Copyright 2023 Harran Ali <harran.m@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"crypto/tls"
	"fmt"
	"net/mail"
	"os"
	"strconv"

	"github.com/harranali/mailing"
)

type Mailer struct {
	mailer     *mailing.Mailer
	sender     mail.Address
	receiver   mail.Address
	cc         []mail.Address
	bcc        []mail.Address
	subject    string
	htmlBody   string
	plainText  string
	attachment string
}

type EmailAddress struct {
	Name    string // the name can be empty
	Address string // ex: john@example.com
}
type EmailAttachment struct {
	Name string // name of the file
	Path string // full path to the file
}

func initiateMailerWithSMTP() *Mailer {
	portStr := os.Getenv("SMTP_PORT")
	if portStr == "" {
		panic("error reading smtp port env var")
	}
	port, err := strconv.ParseInt(portStr, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("error parsing smtp port env var: %v", err))
	}
	skipTlsVerifyStr := os.Getenv("SMTP_TLS_SKIP_VERIFY_HOST")
	if skipTlsVerifyStr == "" {
		panic("error reading smtp tls verify env var")
	}
	skipTlsVerify, err := strconv.ParseBool(skipTlsVerifyStr)
	if err != nil {
		panic(fmt.Sprintf("error parsing smtp tls verify env var: %v", err))
	}

	return &Mailer{
		mailer: mailing.NewMailerWithSMTP(&mailing.SMTPConfig{
			Host:     os.Getenv("SMTP_HOST"),
			Port:     int(port),
			Username: os.Getenv("SMTP_USERNAME"),
			Password: os.Getenv("SMTP_PASSWORD"),
			TLSConfig: tls.Config{
				ServerName:         os.Getenv("SMTP_HOST"),
				InsecureSkipVerify: skipTlsVerify,
			},
		}),
	}

}

func initiateMailerWithSparkPost() *Mailer {
	apiVersionStr := os.Getenv("SPARKPOST_API_VERSION")
	if apiVersionStr == "" {
		panic("error reading sparkpost base url env var")
	}
	apiVersion, err := strconv.ParseInt(apiVersionStr, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("error parsing sparkpost base url env var: %v", apiVersion))
	}
	return &Mailer{
		mailer: mailing.NewMailerWithSparkPost(&mailing.SparkPostConfig{
			BaseUrl:    os.Getenv("SPARKPOST_BASE_URL"),
			ApiKey:     os.Getenv("SPARKPOST_API_KEY"),
			ApiVersion: int(apiVersion),
		}),
	}
}

func initiateMailerWithSendGrid() *Mailer {
	return &Mailer{
		mailer: mailing.NewMailerWithSendGrid(&mailing.SendGridConfig{
			Host:     os.Getenv("SENDGRID_HOST"),
			Endpoint: os.Getenv("SENDGRID_ENDPOINT"),
			ApiKey:   os.Getenv("SENDGRID_API_KEY"),
		}),
	}
}

func initiateMailerWithMailGun() *Mailer {
	skipTlsVerifyStr := os.Getenv("MAILGUN_TLS_SKIP_VERIFY_HOST")
	if skipTlsVerifyStr == "" {
		panic("error reading mailgun tls verify env var")
	}
	skipTlsVerify, err := strconv.ParseBool(skipTlsVerifyStr)
	if err != nil {
		panic(fmt.Sprintf("error parsing mailgun tls verify env var: %v", err))
	}
	return &Mailer{
		mailer: mailing.NewMailerWithMailGun(&mailing.MailGunConfig{
			Domain:              os.Getenv("MAILGUN_DOMAIN"),
			APIKey:              os.Getenv("MAILGUN_API_KEY"),
			SkipTLSVerification: skipTlsVerify,
		}),
	}
}

func (m *Mailer) SetFrom(emailAddresses EmailAddress) *Mailer {
	e := mailing.EmailAddress{
		Name:    emailAddresses.Name,
		Address: emailAddresses.Address,
	}
	m.mailer.SetFrom(e)
	return m
}

func (m *Mailer) SetTo(emailAddresses []EmailAddress) *Mailer {
	var addressesList []mailing.EmailAddress
	for _, v := range emailAddresses {
		addressesList = append(addressesList, mailing.EmailAddress{Name: v.Name, Address: v.Address})
	}

	m.mailer.SetTo(addressesList)
	return m
}

func (m *Mailer) SetCC(emailAddresses []EmailAddress) *Mailer {
	var addressesList []mailing.EmailAddress
	for _, v := range emailAddresses {
		addressesList = append(addressesList, mailing.EmailAddress{Name: v.Name, Address: v.Address})
	}

	m.mailer.SetCC(addressesList)
	return m
}

func (m *Mailer) SetBCC(emailAddresses []EmailAddress) *Mailer {
	var addressesList []mailing.EmailAddress
	for _, v := range emailAddresses {
		addressesList = append(addressesList, mailing.EmailAddress{Name: v.Name, Address: v.Address})
	}

	m.mailer.SetBCC(addressesList)
	return m
}

func (m *Mailer) SetSubject(subject string) *Mailer {
	m.mailer.SetSubject(subject)
	return m
}

func (m *Mailer) SetHTMLBody(body string) *Mailer {
	m.mailer.SetHTMLBody(body)
	return m
}

func (m *Mailer) SetPlainTextBody(body string) *Mailer {
	m.mailer.SetPlainTextBody(body)
	return m
}

func (m *Mailer) SetAttachments(attachments []EmailAttachment) *Mailer {
	var aList []mailing.Attachment
	for _, v := range attachments {
		aList = append(aList, mailing.Attachment{Name: v.Name, Path: v.Path})
	}
	m.mailer.SetAttachments(aList)
	return m
}

func (m *Mailer) Send() error {
	return m.mailer.Send()
}
