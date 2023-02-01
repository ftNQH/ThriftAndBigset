package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"os"
	"projectThrift/clientgin/routes/docs"

	"net/http"
	"projectThrift/gen-go/user_item"
	"strconv"
)

var ctx = context.Background()
var client user_item.Item

func runClient(transportFactory thrift.TTransportFactory, protocolFactory thrift.TProtocolFactory, addr string, secure bool) (thrift.TTransport, error) {
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
		return nil, err
	}
	if transport == nil {
		return nil, fmt.Errorf("Error opening socket, got nil transport. Is server available?")
	}
	transport, _ = transportFactory.GetTransport(transport)
	if transport == nil {
		return nil, fmt.Errorf("Error from transportFactory.GetTransport(), got nil transport. Is server available?")
	}

	err = transport.Open()
	if err != nil {
		return nil, err
	}
	//defer transport.Close()

	return transport, nil
}
func Usage() {
	fmt.Fprint(os.Stderr, "Usage of ", os.Args[0], ":\n")
	flag.PrintDefaults()
	fmt.Fprint(os.Stderr, "\n")
}
func Get() (thrift.TTransport, thrift.TProtocolFactory) {
	flag.Usage = Usage
	// always be a client in this copy
	//server := flag.Bool("server", false, "Run server")
	protocol := flag.String("P", "binary", "Specify the protocol (binary, compact, json, simplejson)")
	framed := flag.Bool("framed", false, "Use framed transport")
	buffered := flag.Bool("buffered", false, "Use buffered transport")
	addr := flag.String("addr", "localhost:18723", "Address to listen to")
	secure := flag.Bool("secure", false, "Use tls secure transport")

	flag.Parse()

	var protocolFactory thrift.TProtocolFactory
	switch *protocol {
	case "compact":
		protocolFactory = thrift.NewTCompactProtocolFactory()
	case "simplejson":
		protocolFactory = thrift.NewTSimpleJSONProtocolFactory()
	case "json":
		protocolFactory = thrift.NewTJSONProtocolFactory()
	case "binary", "":
		protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()
	default:
		fmt.Fprint(os.Stderr, "Invalid protocol specified", protocol, "\n")
		Usage()
		os.Exit(1)
	}

	var transportFactory thrift.TTransportFactory
	if *buffered {
		transportFactory = thrift.NewTBufferedTransportFactory(8192)
	} else {
		transportFactory = thrift.NewTTransportFactory()
	}

	if *framed {
		transportFactory = thrift.NewTFramedTransportFactory(transportFactory)
	}

	// always be client
	fmt.Printf("*secure = '%v'\n", *secure)
	fmt.Printf("*addr   = '%v'\n", *addr)
	transport, err := runClient(transportFactory, protocolFactory, *addr, *secure)
	if err != nil {
		fmt.Println("error running client:", err)
	}
	return transport, protocolFactory
}
func main() {
	transport, protocolFactory := Get()
	fmt.Println(transport, protocolFactory)
	client = user_item.NewItemClientFactory(transport, protocolFactory)
	println(client, ": clientgin/client1.go:110")
	r := gin.Default()
	docs.SwaggerInfo.BasePath = ""

	r.POST("/item", newItems)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.GET("/item", getItems)
	r.DELETE("/item/:id", deleteItems)
	r.PUT("/item/:id", editItems)
	r.POST("/user", addUser)
	r.GET("/item/:uid", GetItemByUserID)
	err := r.Run()
	if err != nil {
		return
	} // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags Item
// @Accept json
// @Produce json
// @Param pos query int64 false "position"
// @Param count query int64 false "count"
// @Success 200 {object} user_item.TItem
// @Router /item [get]
func getItems(c *gin.Context) {

	pos, _ := c.GetQuery("pos")
	cou, _ := c.GetQuery("count")
	position, _ := strconv.Atoi(pos)
	count, _ := strconv.Atoi(cou)
	val, err := client.GetItem(ctx, int16(position), int16(count))
	if err != nil {
		logrus.Error(err)
	}
	c.IndentedJSON(http.StatusOK, val)
}

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags Item
// @Accept json
// @Param NewItem body user_item.TItem true "Thông tin item mới "
// @Produce json
// @Success 200 {object} user_item.TItem
// @Router /item [post]
func newItems(c *gin.Context) {
	var NewItem user_item.TItem
	if err := c.BindJSON(&NewItem); err != nil {
		return
	}
	val, _ := client.AddItem(ctx, int16(NewItem.ID), &NewItem)
	c.IndentedJSON(http.StatusCreated, val)

}

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags Item
// @Accept json
// @Param id path int true "Id cần xóa"
// @Produce json
// @Success 200 {object} user_item.TItem
// @Router /item/{id} [delete]
func deleteItems(c *gin.Context) {
	id := c.Param("id")
	d, _ := strconv.Atoi(id)
	_, err := client.DeleteItem(ctx, int16(d))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
	}
	c.IndentedJSON(http.StatusOK, "xóa thành công")
}

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags Item
// @Accept json
// @Param id path int true "ID cần sửa"
// @Param ItemEdit body user_item.TItem true "Thông tin cần sửa"
// @Produce json
// @Success 200 {object} user_item.TItem
// @Router /item/{id} [put]
func editItems(c *gin.Context) {
	id := c.Param("id")
	d, _ := strconv.Atoi(id)
	var ItemEdit user_item.TItem
	if err := c.BindJSON(&ItemEdit); err != nil {
		c.IndentedJSON(http.StatusNotFound, err)
	}
	val, err := client.EditItems(ctx, int16(d), &ItemEdit)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, err)
	}
	c.IndentedJSON(http.StatusOK, &val)
}

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags User
// @Accept json
// @Param NewUser body user_item.TUser true "Thông tin user mới "
// @Produce json
// @Success 200 {object} user_item.TUser
// @Router /user [post]
func addUser(c *gin.Context) {
	var NewUser user_item.TUser
	if err := c.BindJSON(&NewUser); err != nil {
		c.IndentedJSON(http.StatusNotFound, err)
	}
	val, _ := client.AddUser(ctx, &NewUser)
	c.IndentedJSON(http.StatusCreated, val)

}

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags Item
// @Accept json
// @Param uid path int true "ID của người dùng cần tìm kiếm "
// @Param pos query int64 false "position"
// @Param count query int64 false "count"
// @Produce json
// @Success 200 {object} user_item.TItem
// @Router /item/{uid} [get]
func GetItemByUserID(c *gin.Context) {
	uid, _ := strconv.Atoi(c.Param("uid"))
	val, _ := client.GetItemByUID(ctx, int16(uid))
	c.IndentedJSON(http.StatusOK, val)
}
