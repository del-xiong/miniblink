//网络文件系统
//拦截webview所有的请求,实现特定域的伪文件系统
#include "netfs.h"
#include "export.h"

//初始化指定webview的网络文件系统
void initNetFS(wkeWebView window)
{
    wkeOnLoadUrlBegin(window, handleLoadUrlBegin, NULL);
    wkeOnLoadUrlEnd(window, handleLoadUrlEnd, NULL);
}

//url加载开始,回调
bool handleLoadUrlBegin(wkeWebView window, void *param, const char *url, wkeNetJob job)
{
    //从golang获取网络文件系统数据
    struct goGetNetFSData_Return returnValue = goGetNetFSData(window, url);
    if (returnValue.result == 1)
    {
        // 返回1,表示网络文件系统不处理
        // 判断黑白名单/urlEndCb处理
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

    if (returnValue.result == 0)
    {
        //设置mimetype
        wkeNetSetMIMEType(job, returnValue.mineType);
        free(returnValue.mineType);
        //设置返回的数据
        wkeNetSetData(job, returnValue.data, returnValue.length);
        free(returnValue.data);
        return true;
    }
    else
    {
        //TODO:暂时返回不处理,交由上层,因为不知道怎么返回404
        return false;
    }
}

//url加载完毕,回调
void WKE_CALL_TYPE handleLoadUrlEnd(wkeWebView window, void* param, const char *url, void *job, void *buf, int len) {
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