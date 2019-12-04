#include "event.h"
#include "export.h"

void initGlobalEvent()
{
}

//当文档ready的时候
void WKE_CALL_TYPE onDocumentReady2Callback(wkeWebView window, void *param, wkeWebFrameHandle frameId)
{
    //只触发main frame 的 ready
    if (wkeWebFrameGetMainFrame(window) == frameId)
    {
        goOnDocumentReadyCallback(window);
    }
}

//当网页标题(title)改变的时候
void WKE_CALL_TYPE onTitleChangedCallback(wkeWebView window, void *param, const wkeString title)
{
    goOnTitleChangedCallback(window, wkeGetString(title));
}

bool WKE_CALL_TYPE onUrlLoadBegin(wkeWebView window, void *param, const char *url, wkeNetJob job)
{
    struct goOnUrlLoadBeginCheck_Return checkReturn  = goOnUrlLoadBeginCheck(window, url);
    if (checkReturn.checkFailed) {
        wkeNetCancelRequest(job);
        return true;
    }
    // 设置了回调才hook 因为很影响性能
    if (checkReturn.urlEndCbDefined) {
        wkeNetHookRequest(job);
    }
    return false;
}

void WKE_CALL_TYPE onUrlLoadEnd(wkeWebView window, void* param, const char *url, void *job, void *buf, int len) {
    char * data;
    char * databuf;
    const char * mime;
    databuf = (char *)buf;
    mime = wkeNetGetMIMEType(job, NULL);
    data = goOnUrlLoadEndHandle(window, mime, url, databuf, &len);
    // 0值不要调setdata 否则blink可能会崩溃
    if (len > 0) {
        wkeNetSetData(job, data, len);
    }
}

void initWebViewEvent(wkeWebView window)
{
    //窗口被销毁
    wkeOnWindowDestroy(window, goOnWindowDestroyCallback, NULL);
    //JS引擎初始化完毕
    wkeOnDidCreateScriptContext(window, onDidCreateScriptContextCallback, NULL);
    //document ready
    wkeOnDocumentReady2(window, onDocumentReady2Callback, NULL);
    //title changed
    wkeOnTitleChanged(window, onTitleChangedCallback, NULL);
    // load url begin
    wkeOnLoadUrlBegin(window, onUrlLoadBegin, NULL);
    // load url end
    wkeOnLoadUrlEnd(window, onUrlLoadEnd, NULL);
}