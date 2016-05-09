package controller

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/hkohlsaat/vtr/model"
	"github.com/julienschmidt/httprouter"
)

func PostPlan(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseMultipartForm(65536)
	file, _, err := r.FormFile("vertretungsplan")
	if err != nil {
		log.Printf("error recieving file from upload: %v\n", err)
		return
	}
	defer file.Close()

	now := time.Now()
	f, err := os.OpenFile(now.Format("2006-01-02-15-04-05"), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Printf("can't open file to save plan: %v\n", err)
	} else {
		defer f.Close()
		io.Copy(f, file)
	}

	_ = "breakpoint"
	file.Seek(0, 0)
	plan, err := model.ToPlan(file)
	log.Printf("%#v", *plan)
}
