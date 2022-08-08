package errors

import (
	"fmt"
	"net/http"
	"sync"
)

var (
	unknownCoder defaultCoder = defaultCoder{
		0,
		http.StatusInternalServerError,
		"An internal server error occurred",
		"http://github.com/chsir-zy/errors/README.md",
	}
)

func init() {
	codes[unknownCoder.Code()] = unknownCoder
}

type Coder interface {
	// http 状态错误码
	HTTPStatus() int

	// 用户的错误文本
	String() string

	// 为用户返回详细文档
	Reference() string

	//错误码
	Code() int
}

type defaultCoder struct {
	C    int
	HTTP int
	Ext  string
	Ref  string
}

func (dc defaultCoder) Code() int {
	return dc.C
}

func (dc defaultCoder) HTTPStatus() int {
	if dc.HTTP == 0 {
		return 500
	}
	return dc.HTTP
}

func (dc defaultCoder) String() string {
	return dc.Ext
}

func (dc defaultCoder) Reference() string {
	return dc.Ref
}

var codes = map[int]Coder{}
var codeMux = &sync.Mutex{}

func Register(coder Coder) {
	if coder.Code() == 0 {
		panic("code `0` is reserved by unkowncode")
	}

	codeMux.Lock()
	defer codeMux.Unlock()

	codes[coder.Code()] = coder
}

func MustRegister(coder Coder) {
	if coder.Code() == 0 {
		panic("code `0` is reserved by unkowncode")
	}

	codeMux.Lock()
	defer codeMux.Unlock()

	if _, ok := codes[coder.Code()]; ok {
		panic(fmt.Sprintf("code: %d already exits", coder.Code()))
	}

	codes[coder.Code()] = coder
}

//ParseCoder parse any error into *withCode.
func ParseCoder(err error) Coder {
	if err == nil {
		return nil
	}

	if v, ok := err.(*withCode); ok {
		if coder, ok := codes[v.code]; ok {
			return coder
		}
	}

	return unknownCoder
}

func IsCode(err error, code int) bool {
	if v, ok := err.(*withCode); ok {
		if v.code == code {
			return true
		}

		if v.cause != nil {
			return IsCode(v.cause, code)
		}

		return false
	}

	return true
}
