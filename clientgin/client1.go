package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"projectThrift/clientgin/routes/docs"

	"github.com/go-playground/validator/v10"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"net/http"
	"projectThrift/gen-go/user_item"
	"strconv"
)

var validate *validator.Validate
var ctx = context.Background()
var client user_item.Item

type Response struct {
	Data  map[string]interface{}
	Error map[string]interface{}
}
type RequestItem struct {
	ImageUrl          string                     `json:"imageUrl"`
	Url               string                     `json:"url"`
	Name              string                     `json:"name"`
	Status            int16                      `json:"status"`
	CollectionId      int16                      `json:"collectionId"`
	ArrOwnerAddress   map[string]int             `json:"arrOwnerAddress"`
	ArrOwnerUid       map[string]int             `json:"arrOwnerUid"`
	CreatorAddress    string                     `json:"creatorAddress"`
	OwnerAddress      string                     `json:"ownerAddress"`
	OwnerUID          int16                      `json:"ownerUID"`
	CreatorUid        int16                      `json:"creatorUid"`
	ExternalLink      string                     `json:"externalLink"`
	EventIds          []int16                    `json:"eventIds"`
	TotalView         int16                      `json:"totalView"`
	TotalLike         int16                      `json:"totalLike"`
	Edition           int16                      `json:"edition"`
	MetaData          string                     `json:"metaData"`
	MediaType         string                     `json:"mediaType"`
	CategoryType      int16                      `json:"categoryType"`
	ProductNo         string                     `json:"productNo"`
	TotalEdition      int16                      `json:"totalEdition"`
	ExtendData        user_item.TExtendData      `json:"extendData"`
	IsSensitive       bool                       `json:"isSensitive"`
	Hidden            bool                       `json:"hidden"`
	Price             string                     `json:"price"`
	TokenId           string                     `json:"tokenId"`
	ContractAddress   string                     `json:"contractAddress"`
	Chain             string                     `json:"chain"`
	CurrencyAddress   user_item.TCurrencyAddress `json:"currencyAddress"`
	TypeNft           int16                      `json:"typeNft"`
	TotalClaim        int32                      `json:"totalClaim"`
	TotalLimit        int32                      `json:"totalLimit"`
	StartAt           int16                      `json:"startAt"`
	EndAt             int16                      `json:"endAt"`
	UnlockableContent string                     `json:"unlockableContent"`
	IsLootbox         bool                       `validate:"boolean" json:"isLootbox"`
}

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
	validate = validator.New()
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
	err := r.Run("localhost:8090")
	if err != nil {
		return
	} //listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
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

		b := map[string]interface{}{"code": http.StatusBadRequest, "message": err.Error()}
		res := Response{Data: nil, Error: b}
		c.IndentedJSON(http.StatusOK, res)
	}
	b := map[string]interface{}{"code": http.StatusOK, "message": "Success!"}
	a := map[string]interface{}{"data": val, "total": len(val)}
	res := Response{Data: a, Error: b}
	c.IndentedJSON(http.StatusOK, res)
}

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags Item
// @Accept json
// @Param NewItem body RequestItem true "Th??ng tin Item m???i "
// @Produce json
// @Success 200 {object} user_item.TItem
// @Router /item [post]
func newItems(c *gin.Context) {
	var reqItem RequestItem
	err := c.BindJSON(&reqItem)
	if err != nil {
		fmt.Println(err)
	}
	err = validate.Struct(reqItem)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
			return
		}

		for _, err := range err.(validator.ValidationErrors) {

			fmt.Println(err.Namespace())
			fmt.Println(err.Field())
			fmt.Println(err.StructNamespace())
			fmt.Println(err.StructField())
			fmt.Println(err.Tag())
			fmt.Println(err.ActualTag())
			fmt.Println(err.Kind())
			fmt.Println(err.Type())
			fmt.Println(err.Value())
			fmt.Println(err.Param())
			fmt.Println()
		}
		b := map[string]interface{}{"code": http.StatusBadRequest, "message": err.Error()}
		res := Response{Data: nil, Error: b}
		c.IndentedJSON(http.StatusOK, res)
		// from here you can create your own error messages in whatever language you wish
		return
	}
	var NewItem user_item.TItem
	d, _ := json.Marshal(reqItem)
	err = json.Unmarshal(d, &NewItem)

	if err != nil {
		fmt.Println(err)
		return
	}

	val, err := client.AddItem(ctx, int16(NewItem.ID), &NewItem)

	if err != nil {
		b := map[string]interface{}{"code": http.StatusBadRequest, "message": err.Error()}
		res := Response{Data: nil, Error: b}
		c.IndentedJSON(http.StatusOK, res)
		return
	}
	b := map[string]interface{}{"code": http.StatusOK, "message": "Success!"}
	a := map[string]interface{}{"data": val}
	res := Response{Data: a, Error: b}

	c.IndentedJSON(http.StatusCreated, res)

}

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags Item
// @Accept json
// @Param id path int true "Id c???n x??a"
// @Produce json
// @Success 200 {object} user_item.TItem
// @Router /item/{id} [delete]
func deleteItems(c *gin.Context) {
	id := c.Param("id")
	d, _ := strconv.Atoi(id)
	_, err := client.DeleteItem(ctx, int16(d))
	if err != nil {
		b := map[string]interface{}{"code": http.StatusBadRequest, "message": err.Error()}
		res := Response{Data: nil, Error: b}
		c.IndentedJSON(http.StatusOK, res)
		return
	}
	b := map[string]interface{}{"code": http.StatusOK, "message": "Success!"}
	a := map[string]interface{}{"data": "Xo?? th??nh c??ng"}
	res := Response{Data: a, Error: b}
	c.IndentedJSON(http.StatusOK, res)
}

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags Item
// @Accept json
// @Param id path int true "ID c???n s???a"
// @Param ItemEdit body user_item.TItem true "Th??ng tin c???n s???a"
// @Produce json
// @Success 200 {object} user_item.TItem
// @Router /item/{id} [put]
func editItems(c *gin.Context) {
	id := c.Param("id")
	d, _ := strconv.Atoi(id)
	var ItemEdit user_item.TItem
	if err := c.BindJSON(&ItemEdit); err != nil {
		c.IndentedJSON(http.StatusNotFound, err)
		return
	}
	val, err := client.EditItems(ctx, int16(d), &ItemEdit)
	if err != nil {
		b := map[string]interface{}{"code": http.StatusBadRequest, "message": err.Error()}
		a := map[string]interface{}{"data": nil}
		res := Response{Data: a, Error: b}
		c.IndentedJSON(http.StatusOK, res)
		return
	}
	b := map[string]interface{}{"code": http.StatusOK, "message": "Success!"}
	a := map[string]interface{}{"data": val}
	res := Response{Data: a, Error: b}
	c.IndentedJSON(http.StatusOK, res)
}

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags User
// @Accept json
// @Param NewUser body user_item.TUser true "Th??ng tin user m???i "
// @Produce json
// @Success 200 {object} user_item.TUser
// @Router /user [post]
func addUser(c *gin.Context) {
	var NewUser user_item.TUser
	if err := c.BindJSON(&NewUser); err != nil {
		c.IndentedJSON(http.StatusNotFound, err)
		return
	}
	_, err := client.AddUser(ctx, &NewUser)
	if err != nil {
		b := map[string]interface{}{"code": http.StatusBadRequest, "message": err.Error()}
		a := map[string]interface{}{"data": nil}
		res := Response{Data: a, Error: b}
		c.IndentedJSON(http.StatusOK, res)
		return
	}
	b := map[string]interface{}{"code": http.StatusOK, "message": "Success!"}
	a := map[string]interface{}{"data": NewUser}
	res := Response{Data: a, Error: b}
	c.IndentedJSON(http.StatusOK, res)

}

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags Item
// @Accept json
// @Param uid path int true "ID c???a ng?????i d??ng c???n t??m ki???m "
// @Param pos query int64 false "position"
// @Param count query int64 false "count"
// @Produce json
// @Success 200 {object} user_item.TItem
// @Router /item/{uid} [get]
func GetItemByUserID(c *gin.Context) {
	uid, _ := strconv.Atoi(c.Param("uid"))
	val, err := client.GetItemByUID(ctx, int16(uid))
	if err != nil {
		b := map[string]interface{}{"code": http.StatusBadRequest, "message": err.Error()}
		a := map[string]interface{}{"data": nil}
		res := Response{Data: a, Error: b}
		c.IndentedJSON(http.StatusOK, res)
		return
	}
	b := map[string]interface{}{"code": http.StatusOK, "message": "Success!"}
	a := map[string]interface{}{"data": val, "total": len(val)}
	res := Response{Data: a, Error: b}
	c.IndentedJSON(http.StatusOK, res)
}
