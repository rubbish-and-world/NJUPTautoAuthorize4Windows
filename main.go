package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"golang.org/x/text/encoding/simplifiedchinese"
)

/*
中国电信 : njxy
中国移动 : cmcc
*/
const ISP string = ""  
const username = ""
const password = ""
const SSID string = ""
/*
中国电信 : NJUPT-CHINANET
中国移动 : NJUPT-CMCC
*/


// Test if currently connected to WIFI
func IsWifiConnected() (bool , error) {
	status , ssid := false , ""
	cmd := exec.Command("netsh" , "wlan" , "show" , "interfaces")
	var out bytes.Buffer
	cmd.Stdout = &out 
	err := cmd.Run()
	if err != nil { return status , err }
	decoder := simplifiedchinese.GBK.NewDecoder()
	decodedOutput ,err := decoder.String(out.String())
	lines := strings.Split(decodedOutput, "\n")

	for _ , line := range lines {
		if strings.Contains(line, "状态") && strings.Contains(line, "已连接") {
			status = true
		}
		if strings.Contains(line, "SSID") && strings.Contains(line, SSID) {
			ssid = strings.TrimSpace(strings.Split(line, ":")[1])
		}
	}

	fmt.Println(status , ssid)

	return status , nil
}

//if not connected, connect to the ISP wifi
func connectToWiFi(ssid string) error {
	cmd := exec.Command("netsh", "wlan", "connect", "name="+ssid)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error: %v\noutput: %s", err, string(out))
	}
	fmt.Println("Command output:", string(out))
	return nil
}

//get user IP after ISP WIFI connected
func getIP() (string, error) {
	cmd := exec.Command("ipconfig")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", err
	}

	// Match line containing "10.x.x.x"
	re := regexp.MustCompile(`(10\.\d+\.\d+\.\d+)`)
	matches := re.FindStringSubmatch(out.String())
	if len(matches) >= 2 {
		return matches[1], nil
	}
	return "", fmt.Errorf("IP not found")
}


func sendRequest(url string) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return 
	}

	// Set headers
	req.Header.Set("Sec-Ch-Ua", `"Not:A-Brand";v="99", "Chromium";v="112"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.5615.50 Safari/537.36")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Windows"`)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "no-cors")
	req.Header.Set("Sec-Fetch-Dest", "script")
	req.Header.Set("Referer", "https://p.njupt.edu.cn/")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Login Request Sent")
}


func testHTTPConnection(url string) bool {
	client := http.Client{
		Timeout: 3 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}


func main() {
	for {
		time.Sleep(5 * time.Second)
		status , err := IsWifiConnected()
		if err != nil { log.Fatal(err) }
		if status == false {
			err := connectToWiFi(SSID)
			if err != nil { log.Fatal(err) }
			continue
		}
		//connected
		myIP , err  := getIP()
		if err != nil { log.Fatal(err) }
		fmt.Println("IP :" , myIP)
		
		url := fmt.Sprintf("https://p.njupt.edu.cn:802/eportal/portal/login?callback=dr1003&login_method=1&user_account=%%2C0%%2C%s%%40%s&user_password=%s&wlan_user_ip=%s&wlan_user_ipv6=&wlan_user_mac=000000000000&wlan_ac_ip=&wlan_ac_name=&jsVersion=4.1.3&terminal_type=1&lang=en&v=8426&lang=en", username, ISP, password, myIP)
		sendRequest(url)

		//test if connected
		connected := testHTTPConnection("https://www.baidu.com")
		if connected {
			break
		}
		//if failed, redo all the job
		fmt.Println("Trying to reconnect ...")
	}
}