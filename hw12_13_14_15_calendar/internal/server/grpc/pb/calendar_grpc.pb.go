// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: internal/server/grpc/calendar.proto

package calendarpb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// CalendarClient is the client API for Calendar service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CalendarClient interface {
	Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateResult, error)
	Update(ctx context.Context, in *UpdateRequest, opts ...grpc.CallOption) (*UpdateResult, error)
	Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteResult, error)
	Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResult, error)
	GetEventsListByDates(ctx context.Context, in *GetEventsListByDatesRequest, opts ...grpc.CallOption) (*GetEventsListByDatesResult, error)
	GetEventsForNotify(ctx context.Context, in *GetEventsForNotifyRequest, opts ...grpc.CallOption) (*GetEventsForNotifyResult, error)
	GetEventsListOnDate(ctx context.Context, in *GetEventsListOnDateRequest, opts ...grpc.CallOption) (*GetEventsListOnDateResult, error)
	GetEventsListOnWeek(ctx context.Context, in *GetEventsListOnWeekRequest, opts ...grpc.CallOption) (*GetEventsListOnWeekResult, error)
	GetEventsListOnMonth(ctx context.Context, in *GetEventsListOnMonthRequest, opts ...grpc.CallOption) (*GetEventsListOnMonthResult, error)
}

type calendarClient struct {
	cc grpc.ClientConnInterface
}

func NewCalendarClient(cc grpc.ClientConnInterface) CalendarClient {
	return &calendarClient{cc}
}

func (c *calendarClient) Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateResult, error) {
	out := new(CreateResult)
	err := c.cc.Invoke(ctx, "/calendar.Calendar/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calendarClient) Update(ctx context.Context, in *UpdateRequest, opts ...grpc.CallOption) (*UpdateResult, error) {
	out := new(UpdateResult)
	err := c.cc.Invoke(ctx, "/calendar.Calendar/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calendarClient) Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteResult, error) {
	out := new(DeleteResult)
	err := c.cc.Invoke(ctx, "/calendar.Calendar/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calendarClient) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResult, error) {
	out := new(GetResult)
	err := c.cc.Invoke(ctx, "/calendar.Calendar/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calendarClient) GetEventsListByDates(ctx context.Context, in *GetEventsListByDatesRequest, opts ...grpc.CallOption) (*GetEventsListByDatesResult, error) {
	out := new(GetEventsListByDatesResult)
	err := c.cc.Invoke(ctx, "/calendar.Calendar/GetEventsListByDates", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calendarClient) GetEventsForNotify(ctx context.Context, in *GetEventsForNotifyRequest, opts ...grpc.CallOption) (*GetEventsForNotifyResult, error) {
	out := new(GetEventsForNotifyResult)
	err := c.cc.Invoke(ctx, "/calendar.Calendar/GetEventsForNotify", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calendarClient) GetEventsListOnDate(ctx context.Context, in *GetEventsListOnDateRequest, opts ...grpc.CallOption) (*GetEventsListOnDateResult, error) {
	out := new(GetEventsListOnDateResult)
	err := c.cc.Invoke(ctx, "/calendar.Calendar/GetEventsListOnDate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calendarClient) GetEventsListOnWeek(ctx context.Context, in *GetEventsListOnWeekRequest, opts ...grpc.CallOption) (*GetEventsListOnWeekResult, error) {
	out := new(GetEventsListOnWeekResult)
	err := c.cc.Invoke(ctx, "/calendar.Calendar/GetEventsListOnWeek", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calendarClient) GetEventsListOnMonth(ctx context.Context, in *GetEventsListOnMonthRequest, opts ...grpc.CallOption) (*GetEventsListOnMonthResult, error) {
	out := new(GetEventsListOnMonthResult)
	err := c.cc.Invoke(ctx, "/calendar.Calendar/GetEventsListOnMonth", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CalendarServer is the server API for Calendar service.
// All implementations must embed UnimplementedCalendarServer
// for forward compatibility
type CalendarServer interface {
	Create(context.Context, *CreateRequest) (*CreateResult, error)
	Update(context.Context, *UpdateRequest) (*UpdateResult, error)
	Delete(context.Context, *DeleteRequest) (*DeleteResult, error)
	Get(context.Context, *GetRequest) (*GetResult, error)
	GetEventsListByDates(context.Context, *GetEventsListByDatesRequest) (*GetEventsListByDatesResult, error)
	GetEventsForNotify(context.Context, *GetEventsForNotifyRequest) (*GetEventsForNotifyResult, error)
	GetEventsListOnDate(context.Context, *GetEventsListOnDateRequest) (*GetEventsListOnDateResult, error)
	GetEventsListOnWeek(context.Context, *GetEventsListOnWeekRequest) (*GetEventsListOnWeekResult, error)
	GetEventsListOnMonth(context.Context, *GetEventsListOnMonthRequest) (*GetEventsListOnMonthResult, error)
	mustEmbedUnimplementedCalendarServer()
}

// UnimplementedCalendarServer must be embedded to have forward compatible implementations.
type UnimplementedCalendarServer struct {
}

func (UnimplementedCalendarServer) Create(context.Context, *CreateRequest) (*CreateResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedCalendarServer) Update(context.Context, *UpdateRequest) (*UpdateResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedCalendarServer) Delete(context.Context, *DeleteRequest) (*DeleteResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedCalendarServer) Get(context.Context, *GetRequest) (*GetResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedCalendarServer) GetEventsListByDates(context.Context, *GetEventsListByDatesRequest) (*GetEventsListByDatesResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEventsListByDates not implemented")
}
func (UnimplementedCalendarServer) GetEventsForNotify(context.Context, *GetEventsForNotifyRequest) (*GetEventsForNotifyResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEventsForNotify not implemented")
}
func (UnimplementedCalendarServer) GetEventsListOnDate(context.Context, *GetEventsListOnDateRequest) (*GetEventsListOnDateResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEventsListOnDate not implemented")
}
func (UnimplementedCalendarServer) GetEventsListOnWeek(context.Context, *GetEventsListOnWeekRequest) (*GetEventsListOnWeekResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEventsListOnWeek not implemented")
}
func (UnimplementedCalendarServer) GetEventsListOnMonth(context.Context, *GetEventsListOnMonthRequest) (*GetEventsListOnMonthResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEventsListOnMonth not implemented")
}
func (UnimplementedCalendarServer) mustEmbedUnimplementedCalendarServer() {}

// UnsafeCalendarServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CalendarServer will
// result in compilation errors.
type UnsafeCalendarServer interface {
	mustEmbedUnimplementedCalendarServer()
}

func RegisterCalendarServer(s grpc.ServiceRegistrar, srv CalendarServer) {
	s.RegisterService(&Calendar_ServiceDesc, srv)
}

func _Calendar_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalendarServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/calendar.Calendar/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalendarServer).Create(ctx, req.(*CreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Calendar_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalendarServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/calendar.Calendar/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalendarServer).Update(ctx, req.(*UpdateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Calendar_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalendarServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/calendar.Calendar/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalendarServer).Delete(ctx, req.(*DeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Calendar_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalendarServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/calendar.Calendar/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalendarServer).Get(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Calendar_GetEventsListByDates_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetEventsListByDatesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalendarServer).GetEventsListByDates(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/calendar.Calendar/GetEventsListByDates",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalendarServer).GetEventsListByDates(ctx, req.(*GetEventsListByDatesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Calendar_GetEventsForNotify_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetEventsForNotifyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalendarServer).GetEventsForNotify(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/calendar.Calendar/GetEventsForNotify",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalendarServer).GetEventsForNotify(ctx, req.(*GetEventsForNotifyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Calendar_GetEventsListOnDate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetEventsListOnDateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalendarServer).GetEventsListOnDate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/calendar.Calendar/GetEventsListOnDate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalendarServer).GetEventsListOnDate(ctx, req.(*GetEventsListOnDateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Calendar_GetEventsListOnWeek_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetEventsListOnWeekRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalendarServer).GetEventsListOnWeek(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/calendar.Calendar/GetEventsListOnWeek",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalendarServer).GetEventsListOnWeek(ctx, req.(*GetEventsListOnWeekRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Calendar_GetEventsListOnMonth_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetEventsListOnMonthRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalendarServer).GetEventsListOnMonth(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/calendar.Calendar/GetEventsListOnMonth",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalendarServer).GetEventsListOnMonth(ctx, req.(*GetEventsListOnMonthRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Calendar_ServiceDesc is the grpc.ServiceDesc for Calendar service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Calendar_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "calendar.Calendar",
	HandlerType: (*CalendarServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _Calendar_Create_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _Calendar_Update_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _Calendar_Delete_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _Calendar_Get_Handler,
		},
		{
			MethodName: "GetEventsListByDates",
			Handler:    _Calendar_GetEventsListByDates_Handler,
		},
		{
			MethodName: "GetEventsForNotify",
			Handler:    _Calendar_GetEventsForNotify_Handler,
		},
		{
			MethodName: "GetEventsListOnDate",
			Handler:    _Calendar_GetEventsListOnDate_Handler,
		},
		{
			MethodName: "GetEventsListOnWeek",
			Handler:    _Calendar_GetEventsListOnWeek_Handler,
		},
		{
			MethodName: "GetEventsListOnMonth",
			Handler:    _Calendar_GetEventsListOnMonth_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/server/grpc/calendar.proto",
}