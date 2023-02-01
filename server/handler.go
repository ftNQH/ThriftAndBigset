package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/sirupsen/logrus"
	"projectThrift/gen-go/OpenStars/Core/BigSetKV"
	"projectThrift/gen-go/user_item"
	"strconv"
	"time"
)

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements. See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership. The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License. You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

type thriftHandler struct {
	log map[int]*user_item.TItem
}

func (t thriftHandler) GetItem(ctx context.Context, post int16, count int16) (r user_item.ItemList, err error) {
	var host = "127.0.0.1"
	var port = "18823"
	var connTimeout = time.Second * 3600
	client, _ := GetConnect(host, port, connTimeout)
	result, _ := client.BsGetSlice(ctx, "item", int32(post), int32(count))
	logrus.Error(client.GetTotalCount(ctx, "item"))
	var items user_item.ItemList

	for _, a := range result.Items.Items {
		var item user_item.TItem
		err := json.Unmarshal(a.Value, &item)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	logrus.Error(items)
	return items, nil
}

func (t thriftHandler) DeleteItem(ctx context.Context, id int16) (r *user_item.TItem, err error) {
	var host = "127.0.0.1"
	var port = "18823"
	var connTimeout = time.Second * 3600
	client, _ := GetConnect(host, port, connTimeout)
	bData, _ := json.Marshal(id)
	result, err := client.BsGetItem(ctx, "item", bData)
	if err != nil {
		return nil, err
	}
	if result.Item == nil {
		fmt.Println("Xóa - Item không tồn tại")
		return nil, err
	}

	if _, err := client.BsRemoveItem(ctx, "item", bData); err != nil {
		return nil, err
	}
	fmt.Println("Xóa thành công")
	return nil, err
}

func (t thriftHandler) AddItem(ctx context.Context, id int16, item *user_item.TItem) (r *user_item.TItem, err error) {
	var host = "127.0.0.1"
	var port = "18823"
	var connTimeout = time.Second * 3600
	client, _ := GetConnect(host, port, connTimeout)
	bDataId, _ := json.Marshal(id)
	result, err := client.BsGetItem(ctx, "item", bDataId)
	if err != nil {
		return nil, err
	}
	if result.Item != nil {
		fmt.Println("Add - Item đã tồn tại")
		return nil, err
	}
	bData, _ := json.Marshal(item)
	TItem := BigSetKV.NewTItem()
	TItem.Value = bData
	strKey := strconv.FormatInt(int64(item.ID), 10)
	TItem.Key = []byte(strKey)
	_, err = client.BsPutItem(ctx, "item", TItem)
	if err != nil {
		return nil, err

	}
	fmt.Println("Add - Thêm item thành công")
	return item, nil
}

func (t thriftHandler) EditItems(ctx context.Context, id int16, item *user_item.TItem) (r *user_item.TItem, err error) {
	var host = "127.0.0.1"
	var port = "18823"
	var connTimeout = time.Second * 3600
	client, _ := GetConnect(host, port, connTimeout)
	c, _ := json.Marshal(id)
	result, err := client.BsGetItem(ctx, "item", c)
	if result == nil {
		fmt.Println("Edit - Item chưa tồn tại")
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	var b BigSetKV.TItem
	a, _ := json.Marshal(item)
	b.Key = c
	b.Value = a
	_, err = client.BsPutItem(ctx, "item", &b)
	if err != nil {
		return nil, err
	}
	fmt.Println("Edit - Sửa item thành công")

	return item, nil
}

func (t thriftHandler) AddUser(ctx context.Context, user *user_item.TUser) (r *user_item.TUser, err error) {
	var host = "127.0.0.1"
	var port = "18823"
	var connTimeout = time.Second * 3600
	client, _ := GetConnect(host, port, connTimeout)
	TItem := BigSetKV.NewTItem()
	bData, _ := json.Marshal(user)
	TItem.Value = bData
	TItem.Key, _ = json.Marshal(user.UID)
	_, err = client.BsPutItem(ctx, "user", TItem)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (t thriftHandler) GetItemByUID(ctx context.Context, id int16) (r *user_item.TItem, err error) {
	var host = "127.0.0.1"
	var port = "18823"
	var connTimeout = time.Second * 3600
	client, _ := GetConnect(host, port, connTimeout)
	length, _ := client.GetTotalCount(ctx, "item")
	result, _ := client.BsGetSlice(ctx, "item", 0, int32(length))
	var b user_item.TItem
	var c = BigSetKV.TStringKey(strconv.Itoa(int(id)) + "itembyUID")
	for _, a := range result.Items.Items {

		err := json.Unmarshal(a.Value, &b)
		if err != nil {
			return nil, err
		}
		if b.OwnerUID == id {
			_, err = client.BsPutItem(ctx, c, a)
			if err != nil {
				return nil, err
			}
		}
	}

	return &b, nil
}

func NewThriftHandler() *thriftHandler {
	return &thriftHandler{log: make(map[int]*user_item.TItem)}
}

func GetConnect(host string, port string, connTimeout time.Duration) (*BigSetKV.TStringBigSetKVServiceClient, error) {
	socket, err := thrift.NewTSocketTimeout(fmt.Sprintf("%s:%s", host, port), connTimeout)
	if err != nil {
		return nil, err
	}
	protocolFactory := thrift.NewTBinaryProtocolFactory(true, true)
	var transportFactory thrift.TTransportFactory
	transportFactory = thrift.NewTBufferedTransportFactory(8192)
	transportFactory = thrift.NewTFramedTransportFactory(transportFactory)
	transport, _ := transportFactory.GetTransport(socket)
	c := thrift.NewTStandardClient(protocolFactory.GetProtocol(transport), protocolFactory.GetProtocol(transport))
	var client = BigSetKV.NewTStringBigSetKVServiceClient(c)
	err = transport.Open()
	if err != nil {
		fmt.Println("GetThriftClientCreatorFunc", err)
		return nil, err
	}
	return client, nil
}
