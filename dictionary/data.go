package dictionary

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

type localFS struct{}

var local localFS

type staticFS struct{}

var static staticFS

type Directory struct {
	fs   http.FileSystem
	name string
}

type file struct {
	compressed string
	size       int64
	local      string
	isDir      bool

	once sync.Once
	data []byte
	name string
}

func (fs localFS) Open(name string) (http.File, error) {
	if f, present := data[name]; present {
		return os.Open(f.local)
	} else {
		_, err := os.Stat(name)
		if os.IsNotExist(err) {
			return nil, os.ErrNotExist
		} else {
			return os.Open(name)
		}
	}
}

func (staticFS) prepare(name string) (*file, error) {
	f, present := data[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	var err error
	f.once.Do(func() {
		f.name = path.Base(name)
		if f.size == 0 {
			return
		}
		var gr *gzip.Reader
		b64 := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(f.compressed))
		gr, err = gzip.NewReader(b64)
		if err != nil {
			return
		}
		f.data, err = ioutil.ReadAll(gr)
	})
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (fs staticFS) Open(name string) (http.File, error) {
	f, err := fs.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.file()
}

func (dir Directory) Open(name string) (http.File, error) {
	return dir.fs.Open(dir.name + name)
}

func (f *file) file() (http.File, error) {
	return &httpfile{
		Reader: bytes.NewReader(f.data),
		file:   f,
	}, nil
}

type httpfile struct {
	*bytes.Reader
	*file
}

func (f *file) Stat() (os.FileInfo, error) {
	return f, nil
}

func (f *file) Close() error {
	return nil
}

func (f *file) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

func (f *file) Name() string {
	return f.name
}

func (f *file) Size() int64 {
	return f.size
}

func (f *file) Mode() os.FileMode {
	return 0
}

func (f *file) IsDir() bool {
	return f.isDir
}

func (f *file) ModTime() time.Time {
	return time.Now()
}

func (f *file) Sys() interface{} {
	return f
}

// FS returns a http.Filesystem for the embedded assets. If uselocal is true,
// the filesystem's contents are instead used.
func FS(uselocal bool) http.FileSystem {
	if uselocal {
		return local
	}
	return static
}

// Dir returns a http.Filesystem for the embedded assets on a given prefix dir.
// If uselocal is true, the filesystem's contents are instead used.
func Dir(uselocal bool, name string) http.FileSystem {
	if uselocal {
		return Directory{fs: local, name: name}
	}
	return Directory{fs: static, name: name}
}

// FSByte returns the named file from the embedded assets. If uselocal is
// true, the filesystem's contents are instead used.
func FSByte(uselocal bool, name string) ([]byte, error) {
	if uselocal {
		f, err := local.Open(name)
		if err != nil {
			return nil, err
		}
		b, err := ioutil.ReadAll(f)
		_ = f.Close()
		return b, err
	}
	f, err := static.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.data, nil
}

// FSMustByte is the same as FSByte, but panics if name is not present.
func FSMustByte(uselocal bool, name string) []byte {
	b, err := FSByte(uselocal, name)
	if err != nil {
		panic(err)
	}
	return b
}

// FSString is the string version of FSByte.
func FSString(uselocal bool, name string) (string, error) {
	b, err := FSByte(uselocal, name)
	return string(b), err
}

// FSMustString is the string version of FSMustByte.
func FSMustString(uselocal bool, name string) string {
	return string(FSMustByte(uselocal, name))
}

var data = map[string]*file{

	"/data/en/adjectives": {
		local: "data/en/adjectives",
		size:  119,
		compressed: `
H4sIAAAJbogA/xTLwQrCMBBF0f37q0CkujJYcR8nzzKQdmSc4O+XrC53cZI6FnakEYZifzqedUMaTd+d
uGYsZUUxjzo/609w6ZRwFRT6x3yvhxCrNuKljYbb8R2B+4iZrE4JPLhbTL9p1I4zAAD//9gYba13AAAA
`,
	},

	"/data/en/characters": {
		local: "data/en/characters",
		size:  72,
		compressed: `
H4sIAAAJbogA/wTAwwHEAAAAsH+mOWOc2nanbwKhSCyRyuQKpUqt0er0BqPJbLHa7A6nm7uHp5e3j6+f
vysAAP//cadMlkgAAAA=
`,
	},

	"/data/en/cities": {
		local: "data/en/cities",
		size:  4837,
		compressed: `
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
AA==
`,
	},

	"/data/en/colors": {
		local: "data/en/colors",
		size:  128,
		compressed: `
H4sIAAAJbogA/wTA0QqCYAwF4PvzVhUkEZKEBF0O/4MN55az1ev33dlwtCK6JB1PmsUPtxSfiaHybcRD
w/jBxZvOgZFiGNQXnGufVNCFNXpGQy/1xWErWSXVibFyq9CdOKWuezh6yQjHUBNxfcmi/wAAAP///hYI
poAAAAA=
`,
	},

	"/data/en/companies": {
		local: "data/en/companies",
		size:  3157,
		compressed: `
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
`,
	},

	"/data/en/continents": {
		local: "data/en/continents",
		size:  68,
		compressed: `
H4sIAAAJbogA/3Iszkzk8ssvKslQcMxNLcpMTuQKzi9F4jmWFpcUJeZkJnI5poEFXEuL8gtSuRzzShKL
kktAIoAAAAD///9NGMlEAAAA
`,
	},

	"/data/en/countries": {
		local: "data/en/countries",
		size:  2774,
		compressed: `
H4sIAAAJbogA/1xWS5Lbus6eYxWa/UlVziLa7mfadvy33H3qZAZJiISYAhyQtI88vyu7G7sFyu5OZUKA
AEXi8QHQzY9+QOGYUOC//wkoHdyEBoURbkJP5nQk4xalqnFUhBvp1Mxpr6GQzKEwCa1N3BaW+4wVSlct
0JrcIdxYT5JYnBup3G+5QbiJLUlkleopzs/nmAyDH3DO6ZmsQf6JAgsccMTo1JB9L33AjuIA/hB2GmFB
AS0X2nMenfKZYEHiH5CNbs5iyO7yQgMfGWGhUXg2+JHsTL0e3dSFpnjCwuQjpauJC8MzB1gYJ45D9cbW
89X++C5+ko5Rqm8toVQ7MuOkNsHCshBXt25jxIAjLHLo0R1dZNuzYHWPUX2TpWNY4thox+gMmarAEgU7
3x+oeiPrCJY4jfhhwZLEQ1jd/JgT90KH3ARuYTlgB8uBA/kqCMvBOKYR49W1pbYaq0/PRIGl//whDjo2
boSOahphqdLrl/ebK/1xFd3SqK1h4vYPre4/DNSYsHpxrCw1UdX939NR2QiWppj8GYfGcjp4HpdnaocP
H25JRrQ93P7kRnNiuNWRxa+6Mr85fNdm7NTgrp8OCe5CVWM4zpJfGZMaY6geMgsh3BkncxqTOj7v0sB6
YIR7DHu3+2p+9WmNwfERP8M9mtK7X/f8k+GepcTs3lBagnsjaQd/xIF02W01TEKR3wW15jSQfQCFKcID
NirwgCXyD6TWF2ojygQPg9/3wI1hSGjwYEQtFTK//2BUcPKQsaOg+UDOjr4kGjG4hkwiTXCJwEz+WnCM
mOEhT/7CI3JieCS0rtTHur1VwfAeDXhU6bJhYfrqWaWHxyw92gRPLRVTSin4qrPTT4biyy94ssuJGKjS
H9XaFdGQAjwlDBN8xRE9uV/xgAJfydzer2odCjzjGfdD6V3PJBPCMxs3mBie1Qi/VBu1NFw3JcTwnE/I
CZ4n66dz+XSFGmGFyfvAihoUFVhR1DQorLgpPXDFzeQrtUMiiYlYYMVpyKVTrvK/NDaarYc1tph9pa6A
aI0d9hhbNFhjwBMXMsWiCh0fKTpTxMllFgcM4T24a7TEwr8ywRqzccL5Vmc5u37SlAjW9C+3Cmtu7RLj
tYZOj04FXePVGYpcEgn1VmQpkhkmWKtpW46dHW7lvQm90mCDIzsAN5gtw4YOGGBDjtbZxN/4yjt/COTC
U7XEcImC774TlkxvuEXDPiNsOBNsuCebVz+o9kPDtVFAyZ9Xxdr7o+B7WDZqJ5zg24gCW9zP82uLAbOv
FD1mv3fdL9W3ts0Hpg62KDgibPGQsXLDLuDfzlZNsCXLsB048OHAQhG2nFpke59QW70QS7nHANtMltTb
mcL/o9fiC2VhFXjRseTrJUdPycsJpUOokSX5YEwDBRqni+CRAslV+8wpxVJvGzpyvEhXueXriRkZl82W
yYzm+vTkBb1q3lhaklRUaaBq7grFr3me1ziHVxRq1Gqn43zP1lhaPhDUmDvG6sbQUVA7djBATTZvp3ag
kvLabcBqRSoENUuPBzWCOugR93406LEM/tqHiX6Mq1rHMu9LgV6G1mVzaXrv1s/SGqU7cTt83HDwn4Ha
uFqh7BHq7O2hzsaCI0F9xNBc+9fX4u9Ev71/8nHuKa1P1JFAfeJ0niEN9eS43CGfUGCHP/kCth3KueR2
N+D88Y5Htb9WFBPBTnuFne7JEblT6RF2xsIdzkbstEE/kYUdGLtse5oK8X+j+YFs+xkBS2Qfyldzd/mI
IcNrX8D0uvcfIYJX4URdyVJ1N7JhongVPrP0nY7XbZ1c6d328mv3h+KP/5lXy6UuXutqzaJWfcspTCz9
x4m6foHXc0OX2LyhZEwZ3vwvAKVacprmy6tPjxqmqib6DG8kdM4UEN6YkuAIf2MIPPt8n1P+rd7/9qia
/4IOaAj/0EgC/+ReY0Dv29/nGfmdxwabE/0vAAD///Qw7PzWCgAA
`,
	},

	"/data/en/currencies": {
		local: "data/en/currencies",
		size:  1800,
		compressed: `
H4sIAAAJbogA/3RV3XLbOA+951PgMnkM/6ap49Sf5W932rtjCZWwpkAvSMZ1nn6HsqVkZ7o3GhEEgYOD
Q3D2s+2gEhOUxn8380eogF74xG7mWzYBLUVh0c16NqlBD/9XSdxQlZA4PtIyeD/sW8uaREE7jiG6WY7J
4EuA0aNYBFTVnXgv2kY3R4cecXKZozOIjjnn0Naj4djRASe4OeyIJnzyZ99K7mlt0Losrc/NR8a54V08
7RkS3Tz7FjZU9wa3gOKT52I9o/liNfs2hhoMq9liWnfi+V7aoitlfs9Q2rP2okdx++2cHn6z8egWwYf+
KCMvi/VuihliAu0Lq8VJObqFBSQBbbLCLa5ny5F2IWsT3eKd6472fM5HLzVtgmW9uiVrDzvRxoKyuSXn
FOuO6eGJrYdeH2kLO0W3DL2o1AOye4QbnmVOdUcPr5w6Ng9t4iM9ZfENW3QrxMSmtIDJ8cjQibBVez2n
EdoqplB0s7EQVNwqW4huLX/J5L4WLbEHMCfADQzwSMQdK03ob5ifgm/oW9aao3sy5pppaai7Huy+lMDa
/BfwL0Fb2pTPiOBL1hZ2pXUw0eSeax7OF96yueftmqoz1wJfclxEW9pL2yX3rI2A9vnMHMuitOlmEHTR
PVuhVOCH379H5T7bLfydoOdoYE+vfKGq4xMX7wR/pRcxdl/Ro4hghPoVZyh9Z3VfgzWYrsOG9Qqqptuz
Ccaghyrk1D3Sn0HdJl8gafR/4SM06AjiJf/i/hiytSPvW3hch2pE21bSYEkoqMRtkU2S5DgWv+VfUoe7
bLbBQl0HWop16KP71IaPLpR6fzAGJsbiXoNdcB31+q2f6NvhdBtH93Q7eI+m3O67BHZsmV4zv4VIVfDF
0omX81mU4x3WziOJfjoTvMSOfviQrm4XLOUWnlaxzk2I7n9IMNrLteTfh34YfgX0C8u0Hv5zvDX9WNJW
yI3QzHAcqBtOV+Lf2Ma0lWiLczCeyq58eMOp3JHbvS1rLuEP4bZfmkizn8OM3UObu+XW49Lb6lwm444j
J0T3W7FGenjerh9dZUIv0NOk2yoXGd11UF24YR06EMxVF0nvt86NsjhALtCBiBH+oYMMLnN0yR1MVBo0
VCyHcEQbPjyzSvx4Ng7ZTjzoHO7+chTiaNWLlQdkEtB9cyPaNqEfsf7rsZly/MHK75k9aB68vMH4tzZa
Z7ZU9oSToqdl0Nb9wDCNNxfUHdw/AQAA//+8vc6ZCAcAAA==
`,
	},

	"/data/en/currency_codes": {
		local: "data/en/currency_codes",
		size:  519,
		compressed: `
H4sIAAAJbogA/ySRW7KEIAxE/8+uwIA6PMwNOKPufyG3wA8qVVRDd5+E0zibsHplcYI7hY/e7NWoj7Bs
kccZQ+dixeWMPDLfOGtT73rDN8FvgpeO94IPEV8Ebxm/1vn3dUQuF1mystT7PYdyaWSxhc0Sy60sT0JS
QkJBDqXmlWsRwqqEkGaW+BHiXogWp+5yJ6vJ1G5J2M7I3hKX2Oyyj2nG/ifsQdlzY++ZT3n7fg4hhUay
H+knZK/kM1Juo/RMOY1yVYp7PQabeiSOYmgyLhU0VHQbfTqaK9oDf86wo2JHxk5Pc8blVtoqtJRoe598
h29oOvPmZDRZaSFN/v0n9M3Tu9Cr0O3GhXdnYw/f4PmGyLcKT3n5DM6Tt8aX3bgTm5yG/5zaZ+5w2n8A
AAD///dPkqcHAgAA
`,
	},

	"/data/en/domain_zones": {
		local: "data/en/domain_zones",
		size:  753,
		compressed: `
H4sIAAAJbogA/wTAUYKFIAhA0f+7SzQjU8wnWlOrnyMR2ZCE7IgiGamIIQ25kB8yEEcmspAH+UM+ghAC
YSMkwk5QwkHIhJNQCUZohIswCE6YhJvwEF7CRxRiJG7EnajEg5iJhViJRmzEiziIi3gT/4gv8WNLbCdb
YTO2i+0jRVIiKekgDZKTJmmxZ/aTvbAb+8U+UEEDuqEJ3VFFDzSjFTW0oR39oQN1dKILfdCXo3AYR+MY
HJNjkTdyIleykRv5Iv/Ig+zkyZk4jfPi7JREUcpByRSjNEqnDMpDeSkfVaiBGqmZWqiD6tRJXdSb+mKC
RWzDEqbYgRWsYoY17MI69sMG5tjEFnZjD/aHvdhHE1qkJdpOU1qmVdpF67RBW7SPy+hCT/SdrvSDXuiV
bvRGH3SnT/pDf/kJIzEuhjMW48EFD3jENzzhih94xk+84BU3vOEXPvCJL/zGX/xjRubG3JnKPJgnszAr
05iNeTE7czAn82Y+zI8lLGUVlrGc9bI+buGO3IlbuTN34148O4/zJt7Ju/iEz/ge/gMAAP//iM1ldPEC
AAA=
`,
	},

	"/data/en/first_names": {
		local: "data/en/first_names",
		size:  1353,
		compressed: `
H4sIAAAJbogA/zxU3W7zNgy9P29lJ1k7J/0Q2MHXa8ZmI64ylVFyC+3pB8rd7hLB5PnhId/IKq5UTGYh
XEQXQk92JyOcovxDdy4BA6vKBxveyIQwbZnUfz/IuOCYLJVQcZFM+EU6V5zJWNFzKRWvHFkxkS5GOCZV
woEsRYxbCZgCWVK8yRw4RsaFNiNMZBRwlvXOFiuOfE/+MHDOMhOmIBa54lC1BCF0+uBIeOMoORN6Y9fR
rRWdw41853km/BZ7iArhTCVEZsWV1lZIVgI5jHkZefVU+BlIhXe2VXEIJrmIcrOBMZBywYFKYPPXP4x0
5uyYGFKdGUchZXRRZsawRWG8cvscNzbObodJxktMbuvpix1mYHIsthrxJnExXhrhHWRIpOhycPnDtkgJ
GFNuZBzlzDFW/JI5RYdc6v+sd9UVt7BD9/z1461KZtxoXSv+NNZdGC7JBCP5WJpep9bpYrw3sqq4pM1L
fVgu2gvnvzeOzvO9mdgndQddO2Hc7hWXJBk3p3MNNUbJ+JVsJVxpi9QMa72EcZEYhRSnVWLFmO6iGGjl
jCEF9Qe20oJDHPHevl5xpC9ZMPqrLTgEstgqMj8DbiGtlH8sSU8fxNFHHBu8y/z0hFJc8MLJHu6nqm/A
VPiLFafl29v25sTG/ctOS0haceavxjB7nKmUwN948fW6ydr2w0ngQmYVA398GNcWmU9McyoFJ5N5Tx3/
OP2NkeqadMGL8SNZaxE2wsDe5OhbmfFOsbD9LPEnruz/Xj2zC45pe0TKeGW16kGO6KyEzTDWpuDBhiH5
hMi5z5/oYvN1SEol+OPm6fHAWsULm0s+s+duonX7z3nGSPEZcKFvY53ZIxiSQ4+pomf9i1ZR9LbNjN5I
l6ToFlqdqbkRvOCdqjJ68Qg3w/eEuVz2XW6XYiRdKl5TG8Rpe3hgXVjKGLecOUb06X6v+C1zSe1kOf+T
KeeCa3C2T9zSsrRzwjgYyQNdJPVj9K04RNo1TL6JXiHPPTMteVpxciMHWduBKUkl4d8AAAD//8A/yUNJ
BQAA
`,
	},

	"/data/en/genders": {
		local: "data/en/genders",
		size:  44,
		compressed: `
H4sIAAAJbogA//JNzEnlckvNBVGO6al5KalFXH75ebpOmXmJRZVc7mARBbec0swULkAAAAD//03w7XQs
AAAA
`,
	},

	"/data/en/industries": {
		local: "data/en/industries",
		size:  4922,
		compressed: `
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
`,
	},

	"/data/en/jobs": {
		local: "data/en/jobs",
		size:  2246,
		compressed: `
H4sIAAAJbogA/2RV0W7juA5911cIuB9wv6GbZNpZTFtvXGyfWZm2ObUlg6IyEyz23xeyZFt2X4LokLLI
o8OjB2NcsKJPznFDFsSxWrDLbzRB6IYrcsWJ0aMViLD+3z8v/y4xsLJbk+30g/fky4gE4Lt6aEay5IXT
V45Z++hr25JBVg8WBtfpE7EJJPqMnjqrR7DQ5ejdi67YdQzjGKH1s89L0oZgQ6CrAazd4RW7Fr13M+YM
gWCJBXEjCDmr6wkNwUA+F/0HOR9J8UKGwGYwNB3K/0tCcpk57Mmi97q+e8HR6zPecHDTiFaWTHXqcSQD
g77Yjiwiq1NP2C79byjdqEwayM7btkLVyY1jsCR3/RqEEUx/iE5ofWpvPd2NUxDktcRdAyfnRW/Xr84g
sBNSBD7Aoy7u1HHafI4qGvTTvSOMIXXs6Iz+U9yk6zBNjkW/oentzK7KROHyKWI08cOu1TUM6NWloXj+
ZUAjvKdv+ZO2XuyN2NkxFVPQsQ/Es9U6DrpGwyhRyt/Igo179ENzoyiSAskkfnOu0fM9elGP3+ulk5jC
pidBI+oRLTIMq1Yf0Q2u25eesYX+R4apJ5OlgKyeEAbp9cnFq51TnnCYYvxTv044k19ABaFPYQSrr+hd
YIP+OJLH8FLkd9s6XiYiK2SLCbKNPYZ0GX8GS44LU/mBXQwvR6kf9MHAsZxn4E/cO0iBLSc8o+khqzwf
vlKVBjxnpiZeghkQWFfuF/KW+RLYY/rVFYOJA+wim6/GhGnuLQqgR4YpXmDyoyNDK70VMAyxMVXBPU1y
8zN4mf+Ws1H1wCOY+Mmqv/u5jXTKvei6YtcEI1u5mxttVpdKKNa7Ka3Y/USzueBfAYZoAidnhd1O9Eto
Pe2KhgPJ9q/k/4pRi8jY6ETiFT0Cm/5IToknU92gtLOGFuVeKDJtnIe52JXW+0dI1Tjr6mhGGV6tYgGy
NeTV12nNgS9UZPxLSW7e/e74Mya5Vn4BYyTXhyEVsmB751nhN/RyjE2IptcVSJ8nXtUCbauPb21Ca0No
JWcd3qBaOBgJvHZIxZQUwQ1bnP6rZb/B75Lht/iEIKvNzt551srfld7evLhaZzcuFkknGcW5iOhBlRFK
Vv6OH6vDpTISsnsB/gsAAP//sTEEWMYIAAA=
`,
	},

	"/data/en/languages": {
		local: "data/en/languages",
		size:  821,
		compressed: `
H4sIAAAJbogA/yxSW5IaOwz91yp6K8AMMJfHJaGrUpm/A63YArdMyfYQ+MuCsopkYSm756f1aD2OzvHs
h8kV0ESzcIIKlGajh8mZZoZTMyNP+ZQwcmKaPUYYaPZkE5pzgJVUC+asDkFoLilgBM1jao3zEhxs8qxN
WCAjQGnhRVtsEbkWLJ589vQClVSNCb14+WAv9FJy/fOM6q4e9Kou1JrXlGNbspTLZILcRCMtRduQpbGe
PX3jlLulSQO6gkN50oqjuRazjdUY85VWBQb985tW5QJDFlpDKrhuYRwD05pPxndaiw5CazHpdjEXWhf9
vPLtzAE6yJnedIjKbeebSfLdChxqPiPU5H+4oRGwgSoG0AbJj2JCGzxx9bTxIxttojGUNsWGetLmYe7x
pC0ibZE/6qCtZF8m+bblJ4+nWMzV2h3OPEwM7RDgkB7N+fwiYKxeriB2f39FE9rVq73QLoYhftTGqC42
vLuomZWdidJ+4BMHpj3fquh7vnfvXFUduqM47bZQV+CY9tGyZ9PuGLOPNbxzo/1/kwfogJtgZM2RDki+
GrZG2SE2iQ/RcnGlIjwUveAk9KXw2RfQ1zhORx/jWEFMG4636f0c7/BSs/eq4vHOjb0eF7lSj1EC9RyK
K9RzLiP1HkJ9vHYHSaLUp6gO1Kc7FPRdhtb9XkKhfwEAAP//cenqNzUDAAA=
`,
	},

	"/data/en/last_names": {
		local: "data/en/last_names",
		size:  1764,
		compressed: `
H4sIAAAJbogA/0SV33LiSg/E7/utgCxQCd5QJvVR36Viaz1TjCWOZoxjP/0pDdmcm6TwH01L/VP7MsYS
8KpBsgquMaVIY8arCmdsTWfBCz1iRhNTYvMn/MFG1RgftCQ1bKRn86sfQUfKeKXuVquFWBhHMvP3yUp8
PnL3mweyLtL3ZV7R6mesGnaJ7IZWe4vDxCtOPMeMEzOulG5sOFJK2KTEgv/rJAOObELS84q3KAOuFodQ
cNI7rzjGlHDptBQcjFmw6b2/LXmlg8pKiVf85trVjqywoYmlC5wSzmxPZWwl42MyYcM5+JDuGTsa75/1
MTKv9utBkvGrn8n6jJ2mFCXjUngmK7iQdIFXNFrH0erAltEy99ip3vz6QIKtF2wmu4cFW4qJF7TxwUb+
1J0NbeyCH+By9QtH9eNw9T8fasYZZy5POw5GC1oao3fxSiNnXKn4na2p3jLeOKXFlbl/OFvsGFsW4VJw
Ve2xJXMOWs0ZR/5r804TjyR4Zbl5i2c2W3DW2bWfVAacqXxrOE5D4Ix9UqvH5xBlKC5hKs7TJY6jSsZe
c/nPEUdvISnYJP6q8tBOOXv9g8U/f6LgJdKKIy2c0Swuf6/W40hjTOXZe6ARlyml+CBxdhJ1XNuq5jCu
nAte1XoSvM8s7sYimvqMfczBHU0p5m+AvWb8rOx3vQqlHjubVgc4Byfy3UpccdDRTZ7MZ783roO6ckoZ
V/789H4r/pfCDz/yY+qcnbNW8o6T1H+xu2XsjOY/tSkWW7DVpUdDz+2zOqQ396pf3H1jH8mXihuutRWf
4mSScVDr/chAM46anIPWrf4Gu7okBdtE3Q0vJJFTxpnSWHfBpf+OXdCUfapS8CZ1wfZsw+Rvt5oZl6Li
yz5XJF4mEeei/jhOFdfLnaVzj8l6X6RL4XvwGZxpEcY5slX+HKaGSgk8Z2zMHcGVBnmmjzvSkvdcavH3
v5trmpKf25HgIktfk8IKdpNIlMFp2Br1vlEnEvbUMj+hnaKDZL5ce/1CW5duY2Mu5izv/Fa15cr0cP2e
I4wTzeYNVUo8X3aBHry6wxmbKXvaPVex7hk7DiS3FMVfddH7yA7bYSqRrSbN4tq7MMa+1JbwP8r/eAju
KJeYkuIamH1vdoHujtZ78nRAo1IGHdmWn4D4CXOVZ77XrSO5ZTS8sGEbc9A7mq7TxXOEa6w+yH5Syskg
ySzu9XfEHshW8nk9fJCxlMSOWalMk3ifv4dpYcGB1Qb2r4F+OpCxx1scsZ/qh+S0SBfwwuR7lRxEL20e
Pa03or41XcCJao7sjdbooie7cVWl/mVa0LD0upIrZlFsdfaxNNxHIex1Tmz/BgAA//9sQKZ55AYAAA==
`,
	},

	"/data/en/months": {
		local: "data/en/months",
		size:  86,
		compressed: `
H4sIAAAJbogA//JKzCtNLKrkcktNKgIzfBOLkjO4HAuKMnO4fBMrubxK81K5vEpzKrkcS9NLi0u4glML
SlJzk1KLuPyTS/JBtF9+GUTAJTUZwgAEAAD//6SVRvJWAAAA
`,
	},

	"/data/en/months_short": {
		local: "data/en/months_short",
		size:  48,
		compressed: `
H4sIAAAJbogA//JKzONyS03i8k0s4nIsKOLyTazk8irN4/IqzeFyLE3nCk4t4PJPLuHyyy/jcklN5gIE
AAD//8Xt5DkwAAAA
`,
	},

	"/data/en/name_prefixes": {
		local: "data/en/name_prefixes",
		size:  22,
		compressed: `
H4sIAAAJbogA//ItKtZT8AXhzOJiBd8iPQWXIj0uQAAAAP//vYg46BYAAAA=
`,
	},

	"/data/en/name_suffixes": {
		local: "data/en/name_suffixes",
		size:  29,
		compressed: `
H4sIAAAJbogA//JU8AQhTwXPMIUwBV8XBReXYIWADBcFlzBfLkAAAAD//5N+98YdAAAA
`,
	},

	"/data/en/nouns": {
		local: "data/en/nouns",
		size:  128,
		compressed: `
H4sIAAAJbogA/yzNwQrCMBAE0Pv8VVvwIl60eA/tKIvJJky2in8vld4fvJPloDDV0sTeq3D79mDB3fih
cKmbB6bUiWFNbbdDadkeRmGUrU9iVFpeDJwtMCt5LxY7vHKhvSnMm/8PD9Wcj646PfALAAD//2is2huA
AAAA
`,
	},

	"/data/en/patronymics": {
		local: "data/en/patronymics",
		size:  666,
		compressed: `
H4sIAAAJbogA/zSSUW/bMAyE3+9f2UmWwm2HwAnWZ9ZiLS4yaVByAv37QQ72KELkfXfkQAtnDBYVo32z
F3zKFIkTviQloQVHekjA2KoecIjkae/IvEbcoi2UcYguudga2XEkFU640JbwSX7H0ZRSwJnNZ8Y7q3KJ
uBZ+sOIUnm1s70KNYP/ZaYmmFe/8EMVA2RSfVErkJ87kFTdZrMS6Q+CD3CsG/vlxrvjlpHdcJysFJ5ep
6ayRFZ0G5ydGqotpwNl5Nt9HxI0wcBtyZFXJ+KJU2HGh4jLdceH2eiO3FHC0bU6U8cbqFQfyhM5L3Bxj
3R3M7BiMMWzU2Kc7urTnOphSia245SKK2y55Zm+W31laJrRs/5NnjJTWiA96OuvE+C1TtCY9WkXP+pcW
UfS+TYzeSYMpukBLI/UWBAd8UVVGLynVV+D4sE1ys8tLRUduipE0VLzZvojTNrPybswyxi1nTgm9fX9X
/JGpmLelNv6TK+eCS2y0K24WAgbOmXFwkhldIsU10lNxSPTycGXSvUPW183sl6cVpxbkIEuD0mIqhn8B
AAD//2GM+GaaAgAA
`,
	},

	"/data/en/phone_numbers_format": {
		local: "data/en/phones_format",
		size:  26,
		compressed: `
H4sIAAAJbogA/1LWVVaGY2UuFA4gAAD//9dpslkcAAAA
`,
	},

	"/data/en/state_abbrevs": {
		local: "data/en/state_abbrevs",
		size:  149,
		compressed: `
H4sIAAAJbogA/wTAURLCIAxF0f+7qzdQbYQEp0Sw7n8hHnXU0A9dFFEGJakHj85TnIZVrGOBiTZpN134
gVdcuOGBT3zgSRzEIk7iRThxE4WojJPRGBdvcRmzMCsZ5JdPspIlttiLbez7HwAA///Jy0DGlQAAAA==
`,
	},

	"/data/en/states": {
		local: "data/en/states",
		size:  471,
		compressed: `
H4sIAAAJbogA/0yQ4WrbQBCE/+9T+FWMQxs1lV3aNKY/J9JWt+i0a3ZPUdSnL3cuNHCwMwx8w9wx4xUL
6JgRM+jo8se03hkaCDohy29zFdDJsjlGo5Op8lBkWAs9cMYGZ/qUzWUEfWbzSUCP2CBC3Yhk1OUsahLU
6ShQUGcb6Ole8cRa1mHe6autEi3uIcrUw/cMHalHBIa0BpcS1MuQZIJSL6ocVkC9RNR3u0nTtrpQb1oq
7Myv3sad+Q1jPdvhEcstkjg394U9eG+y53cZrMlf5jOdzUs6nOCWpbKafcBcWy9JjC5zRrIFdHGeTOkb
q8ae31C/7HuykQ9dtBU/bP2Iutt/qGeuU4KZnvkdQT8LEr2wL6aFXsQnqbgrIolOxZSuHOXwP5AYTEOU
rrstotPfAAAA//+e33Im1wEAAA==
`,
	},

	"/data/en/street_suffixes": {
		local: "data/en/street_suffixes",
		size:  132,
		compressed: `
H4sIAAAJbogA/xzJsQoCMQwG4P1/q6MOIg5FD5zD+SPBkEKantSnl3P5lm8x48Sy0wdR6MlA0diMKG1E
okTrXf2FU+hOnNUMl+FbanNcxYkq8f7zkYkqvaOabDz8CmpTT9yaPHHPIBMrI45fQ9TwkPkLAAD//7Ck
MA2EAAAA
`,
	},

	"/data/en/streets": {
		local: "data/en/streets",
		size:  4162,
		compressed: `
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
AP//VN39s0IQAAA=
`,
	},

	"/data/en/top_level_domains": {
		local: "data/en/top_level_domains",
		size:  37,
		compressed: `
H4sIAAAJbogA/0rKrOJKzs/lysxLy+fKS8xN5cpLLeHKL0rnSs8v40pNKeXKzcwBBAAA//8EGszrJQAA
AA==
`,
	},

	"/data/en/weekdays": {
		local: "data/en/weekdays",
		size:  57,
		compressed: `
H4sIAAAJbogA/wouzUtJrOTyzQdTIaWpxSA6PDUlD8IKySgtAjPcijJBVHBiSWkRiAEIAAD//38mK3c5
AAAA
`,
	},

	"/data/en/weekdays_short": {
		local: "data/en/weekdays_short",
		size:  28,
		compressed: `
H4sIAAAJbogA/wouzePyzc/jCilN5QpPTeEKySjlcivK5ApOLOECBAAA///cuZyeHAAAAA==
`,
	},

	"/data/en/words": {
		local: "data/en/words",
		size:  1685,
		compressed: `
H4sIAAAJbogA/1RVUY4rMQj755ZMwu5aSkgaoOrxn0imVd/PamcmAWNslxvYqAw1eQR7LOJwmrJ+ZIlW
GBmcnqPFdHbpxKWEsTqiUx1tLOmPEOIpC9xJOJ8wjekRLMQXobVB0KeojyX0lAVnh5F4njEQr/IHl+KD
LmFnoSfyb0VxJgt1ktdsKHwNYpuy9AN1VPj+5yd+4WRS6RHg90ihebDzr+KGayTD8gwtdgyV7+nyCkjF
CrKrSk5zj5lTRT/V96tNDXfx06y4bFAVE1ZAT2nfgHQoafRHJEsIoz4qyKXPsZigBTVbhlPjK4kSvztv
+NyJG/btJHbxf1sJJ1F04kodis70FMU+CyMd5is6yUtW2eQPzUutcacy1hwrT0kfp0qur3+X3wPcz0cR
FlYwcYPFyAsKQyLZOFFJXiS5iL4n/dbYRsWRpZ/SSKITYgktmUv+RGtqZO8ImjU+UG5SxUxoc6H4Q6M+
mpiDN2kI85G6GFTxqzBDH0Y1CnqczV+NtcIBo7lYTI6eG0e9pe3DuVNu/d079S5NFA5iT12UsVZMBz0i
67/VdVT9wWQkryLTY4EM6jRKYSnsoBITFXuuVMdc44kqekRzbGgTBdstoVVodIURLG+wh5GstYXYkaTL
MQuUSrTJe9bx85MlqIrJyq99tAYHEys6KBdl9x5TI7fqoqejOJ//eB3dV0lK9icuaGlh8zwhryk5CFWY
Q4tjUCq24ZI1bpELlehko4WnG6739YbfTBoaM6+V2Hly9oq+6ybBNMdaYyvnWMiOQDo0LId4xKjU+YUu
NBsXYU+YsoTmSA2E3fR9hMxm0UUr72rn43H2kimtidawGz2ueEv2EflQeZOV0fMt7dQzWjuEfaJDDp3g
M8sHgaYHafLCCTTPeBw7obL2ztoNfC8x4ckFP+a5V6FSxGy7OiEayxSStL/sNJGvoLAcLCpYK8tRo/iX
e1KFS3bA5wH4TnPZ+z/d/lDIRU/SkfGEqMt2Rkk9ht/QUHB+P97NN7rO2CbZvBxKM87Py0M6O/0LAAD/
/13OYw+VBgAA
`,
	},

	"/data/en/full_names_format": {
		local: "data/en/full_names_format",
		size:  23,
		compressed: `
H4sIAAAJbogA/0rLLCouic9LzE0trlGoyUmEcbgAAQAA///p/eKdGQAAAA==
`,
	},

	"/data/en/street_address_format": {
		local: "data/en/street_address_format",
		size:  32,
		compressed: `
H4sIAAAJbogA/youKUpNLSmuUaiBsOKLS9PSMitSQSLKyspcgAAAAP//JhCp7iAAAAA=
`,
	},

	"/data/en/email_address_format": {
		local: "data/en/email_address_format",
		size:  89,
		compressed: `
H4sIAAAJbogA/0rLLCouic9LzE0trslJhDMdalLycxMz8yDc+LT8otzEEq7kjMSixOSS1CIi1AICAAD/
/ybd/wlZAAAA
`,
	},

	"/data/en/domain_names_format": {
		local: "data/en/domain_names_format",
		size:  30,
		compressed: `
H4sIAAAJbogA/0rOzy1IzMtMLa7RqynJL4jPSS1LzYlPyc9NzMwr5gIEAAD//7EGO2MeAAAA
`,
	},

	"/data/en/zip_codes_format": {
		local: "data/en/zips_format",
		size:  8,
		compressed: `
H4sIAAAJbogA/1KGAC5AAAAA//+H3Sc9CAAAAA==
`,
	},

	"/": {
		isDir: true,
		local: "",
	},

	"/data": {
		isDir: true,
		local: "data",
	},

	"/data/en": {
		isDir: true,
		local: "data/en",
	},
}
