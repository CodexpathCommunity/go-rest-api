package mailgun

import (
	"fmt"
	"net/http"
	"github.com/mailgun/mailgun-go/v3"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

func SendSimpleMessageHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	httpc := urlfetch.Client(ctx)

	mg := mailgun.NewMailgun(
		//"YOUR_DOMAIN_NAME", // Domain name
		"sandbox25224b1d21a0489d823328614e0bf07c.mailgun.org",
		//"YOUR_API_KEY",     // API Key
		"3b802a0169e414f8991b84f250717afb-dbdfb8ff-9e66809e",
	)
	mg.SetClient(httpc)

	msg, id, err := mg.Send(appengine.NewContext(r), mg.NewMessage(
		/* From */ "Excited User <mailgun@sandbox25224b1d21a0489d823328614e0bf07c.mailgun.org>",
		/* Subject */ "Message sent from golang application!!",
		/* Body */ "Hey, Congrats on this message!!",
		/* To */ "2014eeb1078@iitrpr.ac.in",
	))
	if err != nil {
		msg := fmt.Sprintf("Could not send message: %v, ID %v, %+v", err, id, msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Message sent!"))
}

func SendComplexMessageHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	httpc := urlfetch.Client(ctx)

	mg := mailgun.NewMailgun(
		"YOUR_DOMAIN_NAME", // Domain name
		"YOUR_API_KEY",     // API Key
	)
	mg.SetClient(httpc)

	message := mg.NewMessage(
		/* From */ "Excited User <mailgun@YOUR_DOMAIN_NAME>",
		/* Subject */ "Hello",
		/* Body */ "Testing some Mailgun awesomness!",
		/* To */ "foo@example.com",
	)
	message.AddCC("baz@example.com")
	message.AddBCC("bar@example.com")
	message.SetHtml("<html>HTML version of the body</html>")
	message.AddAttachment("files/test.jpg")
	message.AddAttachment("files/test.txt")

	msg, id, err := mg.Send(appengine.NewContext(r), message)
	if err != nil {
		msg := fmt.Sprintf("Could not send message: %v, ID %v, %+v", err, id, msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Message sent!"))
}