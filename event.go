package miniblink

//#include "event.h"
import "C"
import (
	"reflect"
	"strings"
	"unsafe"
)

//export goOnWindowDestroyCallback
func goOnWindowDestroyCallback(window C.wkeWebView, param unsafe.Pointer) {
	go func() {
		view := getWebViewByWindow(window)
		view.Emit("destroy", view)
	}()
}

//export goOnDocumentReadyCallback
func goOnDocumentReadyCallback(window C.wkeWebView) {
	go func() {
		view := getWebViewByWindow(window)
		view.Emit("documentReady", view)
	}()
}

//export goOnTitleChangedCallback
func goOnTitleChangedCallback(window C.wkeWebView, titleString *C.char) {
	//把C过来的字符串转化为golang的
	titleGoString := C.GoString(titleString)

	go func() {
		view := getWebViewByWindow(window)
		view.Emit("titleChanged", view, titleGoString)
	}()
}
//export goOnUrlLoadBeginCheck
func goOnUrlLoadBeginCheck(window C.wkeWebView, url *C.char) (bool, bool) {
	urlGoString := C.GoString(url)
	view := getWebViewByWindow(window)
	if !view.isRequestAllowed(urlGoString) {
		view.Emit("requestBlocked", view, urlGoString)
		return true, false
	}
	return false, reflect.ValueOf(view.urlEndHandler).IsValid()
}
//export goOnUrlLoadEndHandle
func goOnUrlLoadEndHandle(window C.wkeWebView, mime *C.char,url *C.char, buf *C.char, charlen *int) unsafe.Pointer  {
	// 如果为0长度 直接返回
	if *charlen == 0 {
		*charlen = 0
		return C.CBytes([]byte(""))
	}
	goContent := C.GoBytes(unsafe.Pointer(buf), C.int(*charlen))
	goUrl := C.GoString(url)
	goMime := C.GoString(mime)
	view := getWebViewByWindow(window)
	// 判断文档类型是否允许处理 不允许则直接返回
	var shouldHandleDoc bool = false
	for _, v := range view.urlEndHandlerMimeTypes {
		if strings.Contains(goMime, v) {
			shouldHandleDoc = true
		}
	}
	if !shouldHandleDoc {
		return C.CBytes(goContent)
	}
	// 执行任务
	v := reflect.ValueOf(view.urlEndHandler)
	t := v.Type()
	var params = make([]reflect.Value, t.NumIn())
	for index := 0; index < t.NumIn(); index++ {
		if index == 1 {
			params[1] = reflect.ValueOf(goMime)
		} else if index == 2 {
			params[2] = reflect.ValueOf(goUrl)
		} else if index ==3 {
			params[3] = reflect.ValueOf(goContent)
		} else {
			paramType := t.In(index)
			//判断参数是引用还是值,取到正确的类型
			if paramType.Kind() == reflect.Ptr {
				params[index] = reflect.New(paramType.Elem())
			} else {
				params[index] = reflect.New(paramType).Elem()
			}
		}
	}
	// 调用函数 获取返回
	handleResult := v.Call(params)
	if len(handleResult) == 0 {
		*charlen = 0
		return C.CBytes([]byte(""))
	}
	returnData := handleResult[0].Bytes()
	// 修改返回给c的数据长度
	*charlen = len([]byte(returnData))
	return C.CBytes(returnData)
}
