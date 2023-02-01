# ThriftAndBigset

#Introduce 
* This is a project created to manage NFT items

#Prerequisites
   * At least Go 1.19 is required to run the tutorial code.
   * The GOPATH may need to be adjusted, alternatively manually put the Go Thrift library files into a suitable location.
   * The Thrift library and compiler must be the same version. Applications might work even with a version mismatch, but this unsupported. To use a specific version of the library, either clone the repository for that version, or use a package manager like dep or Go modules.


#Functions
* AddItem
* EditItem
* DeleteItem
* GetItemByOwnerID
* AddUser

#Additional Information

   * Try using one of the other available protocols, default is binary.
   * Try using the buffered and/or framed transport options.
   * Note that both server and client must use the exact same protocol and transport stack

#Build 
* Using build.sh to autobuild this project

#HowToRun 
*Run Binary server in directory server by using: ./server (run in terrminal cd to server)
*Run Binary clientgin in directory clientgin by using: ./clientgin (run in another terminal cd to clientgin)
*Test this project in http://localhost:8080/swagger/index.html
