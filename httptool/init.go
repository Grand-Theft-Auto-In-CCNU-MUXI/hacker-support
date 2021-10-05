package httptool

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"unicode"

	"github.com/blackfireio/osinfo"
)

func init() {
	fmt.Println("CHECKING YOUR NETWORKS")
	checkNetworks()
	fmt.Println()
	fmt.Println("CHECKING YOUR MACHINES")
	checkMachines()
	fmt.Println()
}

func checkNetworks() {
	url := "https://ip.tool.lu"
	method := "GET"

	client := &http.Client{
	}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		panic(err)
		return
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println("PLEASE CHECK YOUR NETWORKS", err.Error())
		os.Exit(1)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
		return
	}
	r := strings.FieldsFunc(string(body), func(c rune) bool {
		if c == '.' {
			return false
		}

		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	})
	if len(r) >= 2 {
		fmt.Println("YOUR IP:", r[1])
	}
	fmt.Println("STATUS: OK")

	return
}

func checkMachines() {
	info, err := osinfo.GetOSInfo()
	if err != nil {
		fmt.Println("SKIP OS CHECK")
		return
	}

	fmt.Printf("Family:       %v\n", info.Family)
	fmt.Printf("Architecture: %v\n", info.Architecture)
	fmt.Printf("OS:           %v\n", fmt.Sprintf("%s %s %s", info.ID, info.Name, info.Build))
	fmt.Printf("Codename:     %v\n", info.Codename)
}
