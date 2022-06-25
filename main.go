package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

func main() {
	dates := []interface{}{}
	for len(dates) == 0 {
		fmt.Println(time.Now())
		url2 := "https://terminvereinbarung.muenchen.de/fs/termin/index.php?loc=FS&ct=1071898"
		resp, err := http.Post(url2, "application/x-www-form-urlencoded", nil)
		if err != nil {
			log.Fatalln(err)
		}
		cookies := resp.Cookies()
		data := url.Values{}
		data.Set("CASETYPES[FS Umschreibung Ausländischer FS]", "1")
		data.Set("step", "WEB_APPOINT_SEARCH_BY_CASETYPES")
		r, err := http.NewRequest(http.MethodPost, url2, strings.NewReader(data.Encode()))
		if err != nil {
			log.Fatal(err)
		}
		r.AddCookie(cookies[0])
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		resp.Body.Close()
		resp, err = http.DefaultClient.Do(r)
		if err != nil {
			log.Fatalln(err)
		}
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		p, err := regexp.Compile(`jsonAppoints = \'(.*?)\'`)
		if err != nil {
			log.Fatalln(err)
		}
		m := p.FindSubmatch(b)[0]
		j := make(map[string]interface{})
		s := strings.TrimPrefix(string(m), `jsonAppoints = '`)
		s = strings.TrimSuffix(s, `'`)
		json.Unmarshal([]byte(s), &j)
		j = j["Termin FS Allgemeinschalter_G"].(map[string]interface{})
		l := j["appoints"].(map[string]interface{})
		for k, v := range l {
			if len(v.([]interface{})) > 0 {
				fmt.Println(k)
				dates = append(dates, v)
			}
		}
		fmt.Println(dates)
		time.Sleep(3 * time.Minute)
		resp.Body.Close()
	}
}
