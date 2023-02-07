package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"projectThrift/gen-go/OpenStars/Core/BigSetKV"
	"projectThrift/gen-go/user_item"
	"strconv"
	"time"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/sirupsen/logrus"
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
	client, _, transport := GetConnect(host, port, connTimeout)
	defer transport.Close()
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
	client, _, transport := GetConnect(host, port, connTimeout)
	defer transport.Close()
	bData, _ := json.Marshal(id)
	result, err := client.BsGetItem(ctx, "item", bData)
	if err != nil {
		return nil, err
	}
	if result.Item == nil {
		err = errors.New("Item không tồn tại")
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
	client, _, transport := GetConnect(host, port, connTimeout)
	defer transport.Close()
	bDataId, _ := json.Marshal(1)
	max, err := client.BsGetItem(ctx, "MaxID", bDataId)
	if err != nil {
		return nil, err
	}
	if max.Item == nil {
		err = errors.New("Không thấy giá trị trong bigset MaxID")
		return nil, err
	}
	var b int32
	err = json.Unmarshal(max.Item.Value, &b)
	if err != nil {
		return nil, err
	}
	item.ID = b + 1
	item.CreateTime = time.Now().Unix()
	item.UpdateTime = time.Now().Unix()
	bData, _ := json.Marshal(item)
	tItem := BigSetKV.NewTItem()
	tItem.Value = bData
	strKey := strconv.FormatInt(int64(item.ID), 10)
	tItem.Key = []byte(strKey)
	_, err = client.BsPutItem(ctx, "item", tItem)
	if err != nil {
		return nil, err
	}
	tID := BigSetKV.NewTItem()
	cData, err := json.Marshal(item.ID)
	tID.Value = cData
	strKey = strconv.FormatInt(1, 10)
	tID.Key = []byte(strKey)
	_, err = client.BsPutItem(ctx, "MaxID", tID)
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
	client, _, transport := GetConnect(host, port, connTimeout)
	defer transport.Close()
	c, _ := json.Marshal(id)
	result, err := client.BsGetItem(ctx, "item", c)
	if result.Item == nil {
		err = errors.New("Item chưa tồn tại")
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	var b BigSetKV.TItem
	a, _ := json.Marshal(item)
	if id != int16(item.ID) {
		err = errors.New("Không thể sửa id của Item")
		return nil, err
	}
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
	client, _, transport := GetConnect(host, port, connTimeout)
	defer transport.Close()
	tItem := BigSetKV.NewTItem()
	bData, _ := json.Marshal(user)
	tItem.Value = bData
	tItem.Key, _ = json.Marshal(user.UID)
	result, err := client.BsGetItem(ctx, "user", tItem.Key)
	if err != nil {
		return nil, err
	}
	if result.Item != nil {
		err = errors.New("User đã tồn tại")
		return nil, err
	}
	_, err = client.BsPutItem(ctx, "user", tItem)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (t thriftHandler) GetItemByUID(ctx context.Context, id int16) (r user_item.ItemListById, err error) {
	var host = "127.0.0.1"
	var port = "18823"
	var connTimeout = time.Second * 3600
	client, _, transport := GetConnect(host, port, connTimeout)
	defer transport.Close()
	var c = BigSetKV.TStringKey(strconv.Itoa(int(id)) + "_itembyUID")
	var b user_item.TItem
	var list user_item.ItemListById
	bDataId, _ := json.Marshal(id)
	checkExist, err := client.BsGetItem(ctx, "user", bDataId)
	if err != nil {
		return nil, err
	}
	if checkExist.Item == nil {
		err = errors.New("User này chưa tồn tại")
		return nil, err
	}

	length, _ := client.GetTotalCount(ctx, "item")
	result, _ := client.BsGetSlice(ctx, "item", 0, int32(length))
	_, err = client.RemoveAll(ctx, c)
	for _, a := range result.Items.Items {
		err := json.Unmarshal(a.Value, &b)
		if err != nil {
			return nil, err
		}
		if b.OwnerUID == id {
			temp := b
			list = append(list, &temp)
			_, err = client.BsPutItem(ctx, c, a)
			if err != nil {
				return nil, err
			}
		}
	}
	return list, nil
}

func NewThriftHandler() *thriftHandler {
	return &thriftHandler{log: make(map[int]*user_item.TItem)}
}

func GetConnect(host string, port string, connTimeout time.Duration) (*BigSetKV.TStringBigSetKVServiceClient, error, thrift.TTransport) {
	socket, err := thrift.NewTSocketTimeout(fmt.Sprintf("%s:%s", host, port), connTimeout)
	if err != nil {
		return nil, err, nil
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
		return nil, err, nil
	}
	return client, nil, transport
}
