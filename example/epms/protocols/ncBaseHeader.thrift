/*******************************************************************************************
 * ncBaseHeader.thrift
 * 	    Copyright (c) Eisoo Software, Inc.(2012 - ), All rights reserved
 * 
 * Purpose:
 * 	    定义 EPMS Header
 * 
 * Author:
 * 	    xing.baosong@eisoo.com
 * 	    xu.caihua@eisoo.com
 * 
 * Creating Time:
 *      2016-10-26
 ******************************************************************************************/


namespace cpp ncThriftNS
namespace py ncBaseHeader


/**
 * COMMON EPMS MSG
 **/
const string NC_EPMS_CONNECT_MSG = "msg://epms/system/connect"
const string NC_EPMS_HEARTBEAT_MSG = "msg://epms/system/heartbeat"
const string NC_EPMS_DISCONNECT_MSG = "msg://epms/system/disconnect"


/**
 * 连接消息（IP, PORT 直接从连接状态中获取）
 **/
struct ncConnectionInfo {
    1: i32     osVersion,                  // 系统版本
    2: string  hostName,                   // 机器名称
    3: string  processName,                // 进程名称
    4: string  machineCode,                // 机器码
    5: string  ipAddr,                     // ip地址
}


/**
 * 连接请求
 **/
struct ncConnectRequest {
    1: ncConnectionInfo connInfo
    2: i32 connectType,                // 连接类型
    3: i32 detecttime                  // 断网检测时间
    4: i64 reconnectId = -1,           // 重连Id
    5: string guid = "",               // 重连guid
}


/**
 * 连接请求回复
 **/
struct ncConnectReply {
    1: ncConnectionInfo connInfo
    2: i64 reconnectId,           // 重连id
    3: string guid,               // 重连guid
}


/**
 * 消息类型
 **/
enum ncEPMSMsgType {
    NC_EPMS_SEND_MSG                        = 0,  // 发送消息
    NC_EPMS_REPLY_MSG                       = 1,  // 回复消息
    NC_EPMS_SEND_SUCCESS                    = 2,  // 发送结果 - 成功
    NC_EPMS_SEND_FAILED                     = 3,  // 发送结果 - 失败
    NC_EPMS_NO_SUBSCRIBER                   = 4,  // 没有订阅对象
    NC_EPMS_CONNECT                         = 5,  // 连接消息
    NC_EPMS_DISCONNECT                      = 6,  // 断开消息
    NC_EPMS_HEARTBEAT                       = 7,  // 心跳消息
}


/**
 * 消息配置选项
 **/
enum ncEPMSMsgOpt {
    NC_EPMS_ENABLE_COMPRESS = 0x00000001, // 启用压缩
    NC_EPMS_ENABLE_ENCRYPT  = 0x00000002, // 启用加密
}


/**
 * EPMS 消息头
 *
 * EPMS 根据  【type + msgName】 判断接收到的消息是哪种类型
 *            【bufLength + buffer】 为消息的实际内容，将 proto 对象转换为二进制数据块后得出
 *
 *     1. NC_EPMS_SEND_MSG + msgName：     由发送端发送过来的消息，EPMS 接收到后将消息通知给订阅函数
 *     2. NC_EPMS_REPLY_MSG + msgName：    由发送端发送过来的回复消息，EPMS 接收到后将消息通知给发送结果回调函数
 *     3. NC_EPMS_SEND_SUCCESS + msgName： 由发送端发送过来的发送成功，EPMS 接收到后将成功结果通知给发送结果回调函数
 *     4. NC_EPMS_SEND_FAILED + msgName：  由发送端发送过来的发送失败，EPMS 接收到后将失败结果及错误内容通知给发送结果回调函数
 */
struct ncEPMSMsgHeader {
    1: ncEPMSMsgType  msgType,                // 消息类型
    2: string         msgName,                // 消息名称
    3: i64            sourceId,               // 回复消息和发送结果回复，需要带之前的sourceId
    4: string         protoName,              // 消息类型名，用于消息类型校验
    5: i32            bufLength,              // 缓冲块长度 - 【消息内容长度】
    6: binary         buffer,                 // 缓存 buf   - 【消息内容二进制块】
    7: i32            option = 0,             // 消息选项
}


/**
 * 异常类型
 **/
enum ncEPMSExceptionType {
    NC_ROOT_EXCEPTION = 0,            // 根异常
    NC_ABORT_EXCEPTION = 1,           // 中断性异常
    NC_WARN_EXCEPTION = 2,            // 警告异常
    NC_INFO_EXCEPTION = 3,            // 提示性异常
    NC_IGNORE_EXCEPTION = 4,          // 忽略性异常
    NC_NON_CORE_EXCEPTION = 5,        // 非 ncCoreRootException 类型异常
}


/**
 * 系统异常协议
 **/
struct ncEPMSException {
    1: ncEPMSExceptionType expType,     // 异常类型
    2: i32 codeLine,                    // 异常发生的代码行数
    3: i32 errID,                       // 错误 id
    4: string fileName,                 // 发生异常的文件名
    5: string errmsg,                   // 错误消息
    6: string errProvider,              // 错误提供者
    7: list<string> stackInfo,          // 堆栈信息
}


/**
 * 系统异常协议(链式)
 **/
struct ncEPMSExceptionTProto {
    1: ncEPMSException excp,
    2: list<ncEPMSException> nextexcp,
}
