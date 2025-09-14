package jsonrpc

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/v2/http/restful/response"
)

// RPC 方法处理器类型
type HandlerFunc func(ctx context.Context, params any) (any, error)

// 1. 把业务 注册给RPC
func Registry(methodName string, handler HandlerFunc) {
	j := Get()
	j.mu.Lock()
	defer j.mu.Unlock()

	// 验证 handler 的第二个参数必须是指针类型
	handlerType := reflect.TypeOf(handler)
	if handlerType.NumIn() != 2 {
		panic(fmt.Sprintf("handler %s must have exactly 2 parameters", methodName))
	}

	paramType := handlerType.In(1)
	if paramType.Kind() != reflect.Ptr {
		panic(fmt.Sprintf("handler %s second parameter must be a pointer, got %s", methodName, paramType.Kind()))
	}

	// 获取原始函数名
	funcName := getFunctionName(handler)

	j.methods[methodName] = &MethodInfo{
		Name:      methodName,
		Handler:   handler,
		FuncName:  funcName,
		ParamType: paramType,
	}
}

// 注册结构体方法（自动发现以 RPC 开头的方法）
func RegisterService(service any) {
	j := Get()
	j.mu.Lock()
	defer j.mu.Unlock()

	v := reflect.ValueOf(service)
	t := v.Type()
	serviceName := getTypeName(service)

	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		if strings.HasPrefix(method.Name, "RPC") {
			handler := j.createHandlerFromMethod(v, method)

			// 获取参数类型
			paramType := method.Type.In(2)
			funcName := fmt.Sprintf("%s.%s", serviceName, method.Name)

			j.methods[funcName] = &MethodInfo{
				Name:      funcName,
				Handler:   handler,
				FuncName:  funcName,
				ParamType: paramType,
			}
		}
	}
}

// 从结构体方法创建处理器
func (s *JsonRpc) createHandlerFromMethod(receiver reflect.Value, method reflect.Method) HandlerFunc {
	// 创建参数值
	methodType := method.Type

	// 验证方法签名: receiver, context, params
	if methodType.NumIn() != 3 {
		panic(fmt.Sprintf("method %s must have exactly 3 parameters", method.Name))
	}

	// 验证第二个参数是 context.Context
	contextType := methodType.In(1)
	if contextType.String() != "context.Context" {
		panic(fmt.Sprintf("method %s second parameter must be context.Context, got %s", method.Name, contextType))
	}

	// 验证第三个参数必须是指针
	paramType := methodType.In(2)
	if paramType.Kind() != reflect.Ptr {
		panic(fmt.Sprintf("method %s third parameter must be a pointer, got %s", method.Name, paramType.Kind()))
	}

	// 获取参数的实际类型（去掉指针）
	elemType := paramType.Elem()
	return func(ctx context.Context, params any) (any, error) {
		// 创建参数实例
		paramValue := reflect.New(elemType)

		// 如果传入了参数，进行反序列化
		if params != nil {
			// 将 params 转换为 JSON 再反序列化到目标结构
			paramsJSON, err := json.Marshal(params)
			if err != nil {
				return nil, ErrInvalidParams.WithMessagef("Invalid params, %s", err)
			}

			if len(paramsJSON) > 0 && string(paramsJSON) != "null" {
				if err := json.Unmarshal(paramsJSON, paramValue.Interface()); err != nil {
					return nil, ErrInvalidParams.WithMessagef("Invalid params: %s", err.Error())
				}
			}
		}

		// 调用方法
		results := method.Func.Call([]reflect.Value{
			receiver,
			reflect.ValueOf(ctx),
			paramValue,
		})

		// 处理返回结果
		if len(results) != 2 {
			return nil, ErrProtocalError.WithMessagef("Internal error: invalid return values")
		}

		// 处理错误
		errVal := results[1].Interface()
		if errVal != nil {
			return results[0].Interface(), errVal.(error)
		}

		return results[0].Interface(), nil
	}
}

// 处理 JSON-RPC 请求
func (j *JsonRpc) HandleRequest(r *restful.Request, w *restful.Response) {
	var rpcReq Request[json.RawMessage]
	if err := r.ReadEntity(&rpcReq); err != nil {
		response.Failed(w, err)
		return
	}

	// 验证 JSON-RPC 版本
	if rpcReq.JSONRPC != "2.0" {
		response.Failed(w, ErrProtocalError)
		return
	}

	// 获取方法处理器
	j.mu.RLock()
	handler, exists := j.methods[rpcReq.Method]
	j.mu.RUnlock()

	if !exists {
		response.Failed(w, ErrMethodNotFound.WithMessagef("method %s not found", rpcReq.Method))
		return
	}

	// 从 GoRestful 请求中获取上下文，并可以添加额外信息
	ctx := r.Request.Context()

	// 添加请求信息到上下文
	// ctx = context.WithValue(ctx, "rpcMethod", rpcReq.Method)
	// ctx = context.WithValue(ctx, "requestID", rpcReq.ID)
	// ctx = context.WithValue(ctx, "remoteAddr", r.Request.RemoteAddr)

	// 注册的时候拿到的参数的类型 反序列化参数
	req := reflect.New(handler.ParamType.Elem()).Interface()
	err := json.Unmarshal(rpcReq.Params, req)
	if err != nil {
		response.Failed(w, ErrInvalidParams.WithMessagef("unmarshal error, %s", err))
		return
	}

	// 调用处理器
	resp := NewResponse[any]().SetID(rpcReq.ID)
	result, err := handler.Handler(ctx, req)
	if err != nil {
		response.Failed(w, err)
		return
	}
	*resp.Result = result

	// 返回响应
	response.Success(w, resp)
}
