package md

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/yhhaiua/engine/util/treemap"
	"strings"
)

//GetMD5 获取MD5
func GetMD5(src string) string  {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(src))
	cipherStr := md5Ctx.Sum(nil)
	sign := hex.EncodeToString(cipherStr)
	return sign
}

//CheckMd5Value 检测Map值组成的MD5
func CheckMd5Value(tree *treemap.Map,sign,key string) (string,bool)  {
	md := GenerateMd5Value(tree,key)
	if md != sign{
		return md,false
	}
	return md,true
}
//GenerateMd5Value 生成Md5(只有Value)
func GenerateMd5Value(tree *treemap.Map,key string)string  {
	it := tree.Iterator()
	var b strings.Builder
	b.Grow(128)
	for it.Next() {
		value := it.Value()
		b.WriteString(value.(string))
	}
	b.WriteString(key)
	return GetMD5(b.String())
}

//GenerateMd5Map 生成Md5(Key = Value)
func GenerateMd5Map(tree *treemap.Map,key string)string  {
	it := tree.Iterator()
	var b strings.Builder
	b.Grow(128)
	i := 0
	count := tree.Size()
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
	b.WriteString(key)
	return GetMD5(b.String())
}
