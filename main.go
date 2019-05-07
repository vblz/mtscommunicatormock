package main

import (
	"github.com/go-pkgz/lgr"
	"github.com/jessevdk/go-flags"
	"github.com/vblz/mtscommunicatormock/handlers"
	"github.com/vblz/mtscommunicatormock/mtsMock"
	"github.com/vblz/mtscommunicatormock/mtsWsdl"
	"github.com/vblz/mtscommunicatormock/store/inMemory"
	"net/http"
	"os"
	"strconv"
	"time"
)

type options struct {
	Port      uint16 `long:"http-port" env:"HTTP_PORT" default:"9000" description:"port for http connection"`
	UtcOffset int8   `long:"utc-offset" env:"UTC_OFFSET" default:"-4" description:"UTC offset as set in MTS settings"`
	Dbg       bool   `long:"dbg" env:"DEBUG" description:"debug mode"`
}

func main() {
	var opts options
	p := flags.NewParser(&opts, flags.Default)

	if _, err := p.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}
	setupLog(opts.Dbg)

	mtsWsdl.UtcOffset = time.Hour * time.Duration(opts.UtcOffset)

	store := inMemory.NewInMemory()
	mts := mtsMock.NewMtsMock("my_login", "my_password", "my_naming", store)
	handler := handlers.NewHandler(mts, store)
	http.HandleFunc("/test.svc", handler.SoapHandler)
	http.HandleFunc("/ui/send", handler.SendHandler)
	http.HandleFunc("/ui/list", handler.ListHandler)
	lgr.Printf("[INFO] Starting")
	err := http.ListenAndServe(":"+strconv.Itoa(int(opts.Port)), nil)
	if err != nil {
		lgr.Fatalf("ListenAndServe error: %s", err)
	}
}

func setupLog(dbg bool) {
	if dbg {
		lgr.Setup(lgr.Debug, lgr.CallerFile, lgr.Msec, lgr.LevelBraces)
		return
	}
	lgr.Setup(lgr.Msec, lgr.CallerFile, lgr.LevelBraces, lgr.CallerPkg)
}
