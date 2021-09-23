package httputil

import (
	"github.com/yhhaiua/engine/util/treemap"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func Get(url string) []byte {

	transport := http.Transport{
		DisableKeepAlives: true,
	}
	client := &http.Client{
		Transport:&transport,
		Timeout: 5 * time.Second,
	}
	req, err := client.Get(url)
	if err != nil {
		logger.Errorf("Post error : %s",err.Error())
		return nil
	}
	defer req.Body.Close()
	bodys, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Errorf("Post ReadAll error : %s",err.Error())
		return nil
	}
	return bodys
}

func GetTree(url string, tree *treemap.Map) []byte {

	var b strings.Builder
	b.Grow(128)
	b.WriteString(url + "?")

	i := 0
	count := tree.Size()

	it := tree.Iterator()
	for it.Next() {
		i++
		key := it.Key()
		value := it.Value()
		b.WriteString(key.(string))
		b.WriteString("=")
		b.WriteString(value.(string))
		if i != count{
			b.WriteString("&")
		}
	}

	transport := http.Transport{
		DisableKeepAlives: true,
	}
	client := &http.Client{
		Transport:&transport,
		Timeout: 5 * time.Second,
	}
	req, err := client.Get(b.String())
	if err != nil {
		logger.Errorf("Post error : %s",err.Error())
		return nil
	}
	defer req.Body.Close()
	bodys, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Errorf("Post ReadAll error : %s",err.Error())
		return nil
	}
	return bodys
}
func GetMap(url string, tree map[string]string) []byte {
	var b strings.Builder
	b.Grow(128)
	b.WriteString(url + "?")
	i := 0
	count := len(tree)
	for k,v := range tree{
		i++
		b.WriteString(k)
		b.WriteString("=")
		b.WriteString(v)
		if i != count{
			b.WriteString("&")
		}
	}
	transport := http.Transport{
		DisableKeepAlives: true,
	}
	client := &http.Client{
		Transport:&transport,
		Timeout: 5 * time.Second,
	}
	req, err := client.Get(b.String())
	if err != nil {
		logger.Errorf("Post error : %s",err.Error())
		return nil
	}
	defer req.Body.Close()
	bodys, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Errorf("Post ReadAll error : %s",err.Error())
		return nil
	}
	return bodys
}
