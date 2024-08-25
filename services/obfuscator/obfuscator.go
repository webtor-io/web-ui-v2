package obfuscator

// see https://github.com/SpectreH/hunter-js-obfuscator

import (
	"fmt"
	"strconv"
	"strings"
)

type Obfuscator struct {
	code     string
	mask     string
	interval int
	option   int
}

func NewObfuscator(code string) Obfuscator {
	return Obfuscator{
		code:     code,
		mask:     getMask(),
		interval: randomRange(1, 50),
		option:   randomRange(2, 8),
	}
}

func (o Obfuscator) Obfuscate() string {
	rand := randomRange(0, 99)
	rand1 := randomRange(0, 99)

	return fmt.Sprintf("var _0xc%de=[\"\",\"\x73\x70\x6C\x69\x74\",\"\x30\x31\x32\x33\x34\x35\x36\x37\x38\x39\x61\x62\x63\x64\x65\x66\x67\x68\x69\x6A\x6B\x6C\x6D\x6E\x6F\x70\x71\x72\x73\x74\x75\x76\x77\x78\x79\x7A\x41\x42\x43\x44\x45\x46\x47\x48\x49\x4A\x4B\x4C\x4D\x4E\x4F\x50\x51\x52\x53\x54\x55\x56\x57\x58\x59\x5A\x2B\x2F\",\"\x73\x6C\x69\x63\x65\",\"\x69\x6E\x64\x65\x78\x4F\x66\",\"\",\"\",\"\x2E\",\"\x70\x6F\x77\",\"\x72\x65\x64\x75\x63\x65\",\"\x72\x65\x76\x65\x72\x73\x65\",\"\x30\"];function _0xe%dc(d,e,f){var g=_0xc%de[2][_0xc%de[1]](_0xc%de[0]);var h=g[_0xc%de[3]](0,e);var i=g[_0xc%de[3]](0,f);var j=d[_0xc%de[1]](_0xc%de[0])[_0xc%de[10]]()[_0xc%de[9]](function(a,b,c){if(h[_0xc%de[4]](b)!==-1)return a+=h[_0xc%de[4]](b)*(Math[_0xc%de[8]](e,c))},0);var k=_0xc%de[0];while(j>0){k=i[j%%f]+k;j=(j-(j%%f))/f}return k||_0xc%de[11]}eval(function(h,u,n,t,e,r){r=\"\";for(var i=0,len=h.length;i<len;i++){var s=\"\";while(h[i]!==n[e]){s+=h[i];i++}for(var j=0;j<n.length;j++)s=s.replace(new RegExp(n[j],\"g\"),j);r+=String.fromCharCode(_0xe%dc(s,e,10)-t)}return decodeURIComponent(escape(r))}(\"%s\",%d,\"%s\",%d,%d,%d))", rand, rand1, rand, rand, rand, rand, rand, rand, rand, rand, rand, rand, rand, rand, rand, rand, rand1, o.encodeIt(), randomRange(1, 100), o.mask, o.interval, o.option, randomRange(1, 60))
}

func (o Obfuscator) hashIt(str string) string {
	for i, c := range o.mask {
		str = strings.Replace(str, fmt.Sprint(i), string(c), -1)
	}
	return str
}

func (o Obfuscator) encodeIt() string {
	str := ""

	for _, c := range o.code {
		parsed := strconv.FormatInt(int64(c)+int64(o.interval), o.option)
		str += fmt.Sprint(o.hashIt(parsed), string(o.mask[o.option]))
	}

	return str
}

func getMask() string {
	charSet := strShuffle("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	return charSet[0:9]
}

func Obfuscate(code string) string {
	o := NewObfuscator(code)
	return o.Obfuscate()
}
