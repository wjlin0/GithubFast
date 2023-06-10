package main

import (
	"fmt"
	"gopkg.in/xmlpath.v2"
	"io"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	GITHUB_URLS = []string{
		"alive.github.com", "api.github.com", "assets-cdn.github.com",
		"avatars.githubusercontent.com", "avatars0.githubusercontent.com",
		"avatars1.githubusercontent.com", "avatars2.githubusercontent.com",
		"avatars3.githubusercontent.com", "avatars4.githubusercontent.com",
		"avatars5.githubusercontent.com", "camo.githubusercontent.com",
		"central.github.com", "cloud.githubusercontent.com", "codeload.github.com",
		"collector.github.com", "desktop.githubusercontent.com",
		"favicons.githubusercontent.com", "gist.github.com",
		"github-cloud.s3.amazonaws.com", "github-com.s3.amazonaws.com",
		"github-production-release-asset-2e65be.s3.amazonaws.com",
		"github-production-repository-file-5c1aeb.s3.amazonaws.com",
		"github-production-user-asset-6210df.s3.amazonaws.com", "github.blog",
		"github.com", "github.community", "github.githubassets.com",
		"github.global.ssl.fastly.net", "github.io", "github.map.fastly.net",
		"githubstatus.com", "live.github.com", "media.githubusercontent.com",
		"objects.githubusercontent.com", "pipelines.actions.githubusercontent.com",
		"raw.githubusercontent.com", "user-images.githubusercontent.com",
		"vscode.dev", "education.github.com",
	}
	ipsMapMutex        = sync.Mutex{}
	ipsMap             = make(map[string]string)
	osVersion          = runtime.GOOS
	currentDateAndTime = time.Now()
	year               = currentDateAndTime.Year()
	month              = currentDateAndTime.Month()
	day                = currentDateAndTime.Day()
	hour               = currentDateAndTime.Hour()

	writeList []string
	hostPath  string
	PingBack  func([]string) string
	cmd       string
)

func Ping(ipList []string) string {
	var minLatency int64 = -1
	var minLatencyIP string

	for _, ip := range ipList {
		// 执行 ping 命令
		cmd := exec.Command("ping", "-n", "1", ip)
		output, err := cmd.CombinedOutput()
		if err != nil {
			continue
		}

		// 解析 ping 输出，提取延迟时间
		outputStr := string(output)
		fmt.Println(outputStr)
		if strings.Contains(outputStr, "time=") {
			latencyStr := strings.Split(outputStr, "time=")[1]
			latencyStr = strings.Split(latencyStr, " ")[0]
			latency := strings.TrimSuffix(latencyStr, "ms")

			// 转换延迟时间为整数
			var latencyInt int64
			_, err := fmt.Sscanf(latency, "%d", &latencyInt)
			if err != nil {
				continue
			}

			// 更新最低延迟和对应的 IP 地址
			if minLatency == -1 || latencyInt < minLatency {
				minLatency = latencyInt
				minLatencyIP = ip
			}
		}
	}

	// 如果有延迟最低的 IP 地址，则返回它
	if minLatency != -1 {
		return minLatencyIP
	}

	// 随机返回一个 Ping 成功的 IP 地址
	//if len(ipList) > 0 {
	//	randomIndex := rand.Intn(len(ipList))
	//	return ipList[randomIndex]
	//}

	// 如果没有可用的 IP 地址，则返回空字符串
	return ""
}

func PingLinux(ipList []string) string {
	var bestIP string
	var bestRtt time.Duration

	for _, ip := range ipList {
		rtt, err := pingIP(ip)
		if err == nil {
			if bestIP == "" || rtt < bestRtt {
				bestIP = ip
				bestRtt = rtt
			}
		} else {
			fmt.Printf("[Err] %s\n", err.Error())
		}
	}

	if bestIP != "" {
		return bestIP
	}

	// 如果没有延迟最低的 IP，随机返回一个 ping 的 IP
	if len(ipList) > 0 {
		rand.Seed(time.Now().UnixNano())
		randomIndex := rand.Intn(len(ipList))
		return ipList[randomIndex]
	}

	return ""
}

func pingIP(ip string) (time.Duration, error) {

	pinger, err := net.Dial("ip4:icmp", ip)
	if err != nil {
		return -1, err
	}
	defer pinger.Close()

	// 发送ICMP回显请求（Ping）
	msg := []byte("Ping")
	deadline := time.Now().Add(time.Second * 4) // 设置超时时间为4秒
	err = pinger.SetDeadline(deadline)
	if err != nil {
		return -1, err
	}

	start := time.Now()
	_, err = pinger.Write(msg)
	if err != nil {
		return -1, err
	}

	// 接收ICMP回显回应（Pong）
	recv := make([]byte, 1024)
	_, err = pinger.Read(recv)
	if err != nil {
		return -1, err
	}

	rtt := time.Since(start)
	return rtt, nil
}
func stringInArray(str string, arr []string) bool {
	for _, s := range arr {
		if s == str {
			return true
		}
	}
	return false
}
func getIP(session *http.Client, githubURL string) string {
	url := fmt.Sprintf("https://www.ipaddress.com/site/%s", githubURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("[Err]", err)
		return ""
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36")

	resp, err := session.Do(req)
	if err != nil {
		fmt.Println("[Err]", err)
		return ""
	}
	defer resp.Body.Close()

	root, err := xmlpath.ParseHTML(resp.Body)
	if err != nil {
		fmt.Println("[Err]", err)
		return ""
	}

	// 使用XPath查找指定的元素路径
	path := xmlpath.MustCompile("//*[@id=\"tabpanel-dns-a\"]/pre/a[1]")
	if value, ok := path.String(root); ok {
		// 打印找到的元素内容
		return value
	} else {
		return ""
	}
	//body, _ := io.ReadAll(resp.Body)
	//pattern := `(?:[0-9]{1,3}\.){3}[0-9]{1,3}`
	//re := regexp.MustCompile(pattern)
	//ipList := re.FindAllString(string(body), -1)
	////var TList []string
	//var Tlist string
	//for _, ip := range ipList {
	//	if strings.Contains(Tlist, ip) {
	//		continue
	//	}
	//	Tlist += ip + ","
	//}
	//Tlist = strings.TrimRight(Tlist, ",")
	//bestIP := PingBack(strings.Split(Tlist, ","))
	//return bestIP
}

func copyFile(sourcePath, destinationPath string) error {
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(destinationPath)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	err = destinationFile.Sync()
	if err != nil {
		return err
	}

	return nil
}

func writeHost(path string, ipsMap map[string]string) {
	nowDir := filepath.Dir(path)
	bakDir := filepath.Join(nowDir, "bak")
	if _, err := os.Stat(bakDir); os.IsNotExist(err) {
		_ = os.Mkdir(bakDir, 0755)
	}

	for k, v := range ipsMap {
		now := fmt.Sprintf("%d.%d.%d_%d", year, month, day, hour)
		p := fmt.Sprintf("%s %s # %s", v, k, now)
		writeList = append(writeList, p)
	}

	err := os.WriteFile(path, []byte(strings.Join(writeList, "\n")), 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func deleteOldFile(path string) {
	nowDir := filepath.Dir(path)
	bakDir := filepath.Join(nowDir, "bak")
	if _, err := os.Stat(bakDir); os.IsNotExist(err) {
		return
	}

	err := filepath.Walk(bakDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			oldPathInclude := fmt.Sprintf("%d.%d.%d", year, month, day-1)
			if strings.Contains(path, oldPathInclude) {
				_ = os.Remove(path)
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
}

func getHost(path string) {
	ips, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return
	}

	lines := strings.Split(string(ips), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "#") {
			writeList = append(writeList, line)
			continue
		}
		ipDomain := strings.Split(line, " ")
		found := false
		for _, ip := range ipDomain {
			if stringInArray(ip, GITHUB_URLS) {
				found = true
				break
			}
		}
		if found {
			continue
		}

		writeList = append(writeList, line)

	}
}

func check(path string) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		fmt.Println(err)
		return
	}

	if fileInfo.Mode().Perm()&0200 == 0 {
		fmt.Printf("\t%s\n没有写入权限,请用以下方式打开运行文件\ne.g. sudo ./GithubFast(windows 请用管理员运行)", path)
		os.Exit(0)
	}

	dirInfo, err := os.Stat(filepath.Dir(path))
	if err != nil {
		fmt.Println(err)
		return
	}

	if dirInfo.Mode().Perm()&0200 == 0 {
		fmt.Printf("\t%s\n没有写入权限,请用以下方式打开运行文件\ne.g. sudo ./GithubFast(windows 请用管理员运行)", filepath.Dir(path))
		os.Exit(0)
	}
	bakDir := filepath.Join(filepath.Dir(path), "bak")
	if _, err := os.Stat(bakDir); os.IsNotExist(err) {
		err := os.Mkdir(bakDir, 0755)
		if err != nil {
			fmt.Printf("\t%s\n无法创建文件  \n e.g. sudo ./GithubFast (windows 请用管理员运行)", bakDir)
			os.Exit(0)
		}
	}
	backName := filepath.Join(bakDir, fmt.Sprintf("%d.%d.%d_%d.txt", year, month, day, hour))
	if err := copyFile(path, backName); err != nil {
		fmt.Println(err)
		os.Exit(0)

	}

}

func main() {
	fmt.Println("开始运行....")
	switch osVersion {
	case "linux":
		hostPath = "/etc/hosts"
		PingBack = PingLinux
		cmd = "sudo nscd restart"
	case "windows":
		hostPath = "C:\\Windows\\System32\\drivers\\etc\\hosts"
		PingBack = Ping
		cmd = "ipconfig /flushdns"
	case "darwin":
		hostPath = "/etc/hosts"
		PingBack = PingLinux
		cmd = "sudo killall -HUP mDNSResponder"
	}
	check(hostPath)
	fmt.Println("-> 检查权限通过....")
	getHost(hostPath)
	deleteOldFile(hostPath)
	fmt.Println("-> 删除昨日备份....")
	client := &http.Client{}
	var wg sync.WaitGroup
	sleepTime := 1
	for _, uri := range GITHUB_URLS {
		wg.Add(1)
		sleepTime += 1
		go func(url string, sleepTime int) {
			defer wg.Done()
			attempts := 0
			time.Sleep(time.Duration(sleepTime) * time.Second)
			for attempts < 3 {
				ip := getIP(client, url)
				ipsMapMutex.Lock()
				if ip != "" {
					ipsMap[url] = ip
					fmt.Printf("[Run] %s -> %s \n", url, ip)
					ipsMapMutex.Unlock()
					return
				}
				ipsMapMutex.Unlock()
				attempts++
			}
			fmt.Printf("[Exit] Unable to fetch IP for %s after %d attempts.\n", url, attempts)
		}(uri, sleepTime)
	}
	wg.Wait()
	writeHost(hostPath, ipsMap)
	var err error
	switch osVersion {
	case "linux":
	default:
		c := strings.Split(cmd, " ")
		cmd := exec.Command(c[0], c[1:]...)
		err = cmd.Run()
	}
	if err != nil {
		fmt.Printf("[Err]cmd.Run() failed with %s\n", err)
	}
}
