package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

//---------------------------------HOW MANY THREADS TO RUN----------------------------------

var Threads int = 5

//-----------------------------------EDIT WITH RSI HANDLE-----------------------------------

var RSIHandle string = ""

//-----------------------------------EDIT THESE WITH YOUR COOKIE VALUES-----------------------------------

var RsiDevice string = ""
var RsiToken string = ""
var CookieConsent string = ""
var RsiXSRF string = ""

//---------------------------------------------------------------------------------------------------------

var working bool
var countGifts int

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomStringB(length int) string {
	chars := []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}

func pledgeGiftCode() {
	for working {
		Code := RandomStringB(13)
		req, err := http.NewRequest("GET", "https://robertsspaceindustries.com/pledge/claim-gift/"+Code, nil)

		req.AddCookie(&http.Cookie{Name: "wsc_hide", Value: "false", Path: "/", Domain: "robertsspaceindustries.com"})
		req.AddCookie(&http.Cookie{Name: "_rsi_device", Value: RsiDevice, Path: "/", Domain: ".robertsspaceindustries.com"})
		req.AddCookie(&http.Cookie{Name: "wsc_view_count", Value: "5", Path: "/", Domain: ".robertsspaceindustries.com"})
		req.AddCookie(&http.Cookie{Name: "CookieConsent", Value: CookieConsent, Path: "/", Domain: ".robertsspaceindustries.com"})
		req.AddCookie(&http.Cookie{Name: "Rsi-Token", Value: RsiToken, Path: "/", Domain: ".robertsspaceindustries.com"})
		req.AddCookie(&http.Cookie{Name: "moment_timezone", Value: "America%2FPort-au-Prince", Path: "/", Domain: ".robertsspaceindustries.com"})
		req.AddCookie(&http.Cookie{Name: "Rsi-XSRF", Value: RsiXSRF, Path: "/", Domain: ".robertsspaceindustries.com"})

		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.82 Safari/537.36")
		req.Header.Set("authority", "robertsspaceindustries.com")
		req.Header.Set("scheme", "https")
		req.Header.Set("path", "/pledge/claim-gift/"+Code)
		req.Header.Set("cache-control", "max-age=0")
		req.Header.Set("upgrade-insecure-requests", "1")
		req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
		req.Header.Set("sec-gpc", "1")
		req.Header.Set("sec-fetch-site", "none")
		req.Header.Set("sec-fetch-mode", "navigate")
		req.Header.Set("sec-fetch-user", "?1")
		req.Header.Set("sec-fetch-dest", "document")
		req.Header.Set("accept-encoding", "gzip, deflate, br")
		req.Header.Set("accept-language", "en-US,en;q=0.9")

		req.Header.Set("Content-Type", "text/html; charset=UTF-8")

		client := &http.Client{}

		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}

		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			var reader io.ReadCloser
			reader, err = gzip.NewReader(resp.Body)

			body, _ := ioutil.ReadAll(reader)

			if !strings.Contains(string(body), "This gift code is no longer valid") && strings.Contains(string(body), RSIHandle) && strings.Contains(string(body), "2012-2021 Cloud Imperium Rights LLC and Cloud Imperium Rights Ltd.") {
				fmt.Println("HIT!", Code, countGifts, resp.Request.URL.String())
				fmt.Println("https://robertsspaceindustries.com/connect?jumpto=" + resp.Request.URL.String())
				//Good Code, Should be added to account automatically, Next...
			} else if strings.Contains(string(body), "This gift code is no longer valid") && strings.Contains(string(body), RSIHandle) && strings.Contains(string(body), "2012-2021 Cloud Imperium Rights LLC and Cloud Imperium Rights Ltd.") {
				//Bad Code, Next...
			} else if !strings.Contains(string(body), RSIHandle) && !strings.Contains(string(body), "2012-2021 Cloud Imperium Rights LLC and Cloud Imperium Rights Ltd.") {
				fmt.Println("PAGE ERROR", countGifts, Code, resp.Request.URL.String())
				working = false
				break
			} else if strings.Contains(string(body), "Sign into RSI") {
				fmt.Println("BAD COOKIES - CHECK YOUR COOKIE SETTINGS!", countGifts, Code, resp.Request.URL.String())
				working = false
				break
			} else {
				//Unknown Error...
				fmt.Println("UNKNOWN ERROR", countGifts, Code, resp.Request.URL.String())
				fmt.Println(string(body))
				working = false
				break
			}
			resp.Body.Close()
			reader.Close()
			body = []byte("")
			countGifts++
		} else {
			//fmt.Println(" Bad Status Code: ", resp.StatusCode, Code)
			//working = false
			//break
		}

		time.Sleep(500 * time.Millisecond)
	}
}

func setTitle(text string) {
	os := runtime.GOOS
	if os == "windows" {
		syscall.MustLoadDLL("Kernel32.dll").MustFindProc("SetConsoleTitleW").Call(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(text)))) // Remove if compiling for any non-windows
	}
}

func cleanMemory() {
	for {
		debug.FreeOSMemory()
		time.Sleep(60 * time.Second)
	}
}

func main() {
	fmt.Println("Roberts Space Industries: Gift Code Bot")

	working = true

	for i := 0; i < Threads; i++ {
		time.Sleep(250 * time.Millisecond)
		go pledgeGiftCode()
	}

	go cleanMemory()

	for {
		setTitle("Checked Codes: " + string(strconv.Itoa(countGifts)))
		time.Sleep(125 * time.Millisecond)
	}
}
