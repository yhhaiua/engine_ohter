package httputil

import (
	"github.com/yhhaiua/engine/log"
	"github.com/yhhaiua/engine/util/treemap"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

var logger  = log.GetLogger()


func Post(url string,body io.Reader) []byte {

	transport := http.Transport{
		DisableKeepAlives: true,
	}
	client := &http.Client{
		Transport:&transport,
		Timeout: 5 * time.Second,
	}
	req, err := client.Post(url,"application/x-www-form-urlencoded",body)
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

func PostForm(url string,data url.Values) []byte {

	transport := http.Transport{
		DisableKeepAlives: true,
	}
	client := &http.Client{
		Transport:&transport,
		Timeout: 5 * time.Second,
	}
	req, err :=client.PostForm(url,data)
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

func PostFormTree(url string, tree *treemap.Map) []byte {
	u := MapToUrlValue1(tree)

	transport := http.Transport{
		DisableKeepAlives: true,
	}
	client := &http.Client{
		Transport:&transport,
		Timeout: 5 * time.Second,
	}
	req, err := client.PostForm(url,u)
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
func PostFormMap(url string, tree map[string]string) []byte {
	u := MapToUrlValue2(tree)

	transport := http.Transport{
		DisableKeepAlives: true,
	}
	client := &http.Client{
		Transport:&transport,
		Timeout: 5 * time.Second,
	}
	req, err := client.PostForm(url,u)
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
//MapToUrlValue1 生成url
func MapToUrlValue1(tree *treemap.Map)url.Values  {
	it := tree.Iterator()
	u := make(url.Values)
	for it.Next() {
		key := it.Key()
		value := it.Value()
		u.Set(key.(string),value.(string))
	}
	return u
}

//MapToUrlValue2 生成url
func MapToUrlValue2(tree map[string]string)url.Values  {
	u := make(url.Values)
	for key,value := range tree{
		u.Set(key,value)
	}
	return u
}