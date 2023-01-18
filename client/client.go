package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"projectThrift/gen-go/user_item"

	"github.com/apache/thrift/lib/go/thrift"
)

var ctx = context.Background()

func handleClient(client user_item.Item) (err error) {
	/*	item := user_item.TItem{
		ID:                1,
		ImageUrl:          "aaaaaaaaaaaaaa",
		URL:               "",
		Name:              "",
		Status:            0,
		CollectionId:      0,
		CreateTime:        0,
		UpdateTime:        0,
		ArrOwnerAddress:   nil,
		ArrOwnerUid:       nil,
		CreatorAddress:    "",
		CreatorUid:        0,
		ExternalLink:      "",
		EventIds:          nil,
		TotalView:         0,
		TotalLike:         0,
		Edition:           0,
		MetaData:          "",
		MediaType:         "",
		CategoryType:      0,
		ProductNo:         "",
		TotalEdition:      0,
		ExtendData:        nil,
		IsSensitive:       false,
		Hidden:            false,
		Price:             "",
		TokenId:           "",
		ContractAddress:   "",
		Chain:             "",
		CurrencyAddress:   nil,
		TypeNft:           0,
		TotalClaim:        0,
		TotalLimit:        0,
		StartAt:           0,
		EndAt:             0,
		UnlockableContent: "",
		IsLootbox:         false,
	}*/

	client.GetItem(ctx, 0, 3)
	/*client.AddItem(ctx, int16(item.ID), &item)
	fmt.Println("thêm item thành công")*/
	/*client.DeleteItem(ctx, 3)
	client.EditItems(ctx, 0, &item)*/
	//client.GetItemByUID(ctx, 0)
	return err
}

func runClient(transportFactory thrift.TTransportFactory, protocolFactory thrift.TProtocolFactory, addr string, secure bool) error {
	var transport thrift.TTransport
	var err error
	if secure {
		cfg := new(tls.Config)
		cfg.InsecureSkipVerify = true
		transport, err = thrift.NewTSSLSocket(addr, cfg)
	} else {
		transport, err = thrift.NewTSocket(addr)
	}
	if err != nil {
		fmt.Println("Error opening socket:", err)
		return err
	}
	if transport == nil {
		return fmt.Errorf("Error opening socket, got nil transport. Is server available?")
	}
	transport, _ = transportFactory.GetTransport(transport)
	if transport == nil {
		return fmt.Errorf("Error from transportFactory.GetTransport(), got nil transport. Is server available?")
	}

	err = transport.Open()
	if err != nil {
		return err
	}
	defer transport.Close()

	return handleClient(user_item.NewItemClientFactory(transport, protocolFactory))
}
