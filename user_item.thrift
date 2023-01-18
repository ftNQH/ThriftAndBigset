namespace go user_item
struct TExtendData{
	1:string Category
	2:string Collection
	3:string Create
	4:string Like
	5:string Owner
	6:string TotalOffer
}
struct TCurrencyAddress{
	1:string Name
	2:string Symbol
	3:string Address
	4:string Chain
	5:i16 Decimal
	7:string Img
	8:bool Status	
}
struct TUser{
	1:required i32 Uid,
	2: string Name
	3: string Username
	4: string Pubkey
	5: string Bio
	6: string Avatar
	7: string Background 
	8: string Email
	9: string Facebook
	10:string Instagram
	11:string Twitter
	12:string Youtube 
	13:string CountryCode
	14:string Language
	15:string BlockchainAddress
	16:i16 Vip
	17:i16 TypeUser
	18:bool CanMint
	19:bool Verify
	20:bool IsPartner
	21:string App
	22:i32 TotalFollowing 
	23:i32 TotalFollers
	24:TExtendData ExtendData
	25:bool Block
	26:bool Deleted
	27:i16 CreateTime
	28:i16 UpdateTime
}
struct TItem{
	1: required i32 Id
	2: string ImageUrl
	3: string Url
	4: string Name
	5: i16 Status
	6: i16 CollectionId
	7: i16 CreateTime 
	8: i16 UpdateTime
	9: map<string,i16> ArrOwnerAddress
	10:map<string,i16> ArrOwnerUid
	11:string CreatorAddress
	12:string OwnerAddress
	13:i16 OwnerUID
	14:i16 CreatorUid
	15:string ExternalLink
	16:list<i16> EventIds
	17:i16 TotalView
	18:i16 TotalLike
	19:i16 Edition
	20:string MetaData
	21:string MediaType
	22:i16 CategoryType
	23:string ProductNo
	24:i16 TotalEdition
	25:TExtendData ExtendData
	26:bool IsSensitive
	27:bool Hidden
	28:string Price
	29:string TokenId
	30:string ContractAddress
	31:string Chain
	32:TCurrencyAddress CurrencyAddress
	33:i16 TypeNft
	34:i32 TotalClaim
	35:i32 TotalLimit
	36:i16 StartAt
	37:i16 EndAt
	38:string UnlockableContent
	39:bool IsLootbox
}


typedef string TStringKey
typedef binary TItemKey //key - index of an item, with simple set, itemkey is equivalent to item
typedef list<TItem> ItemList

typedef list<TItem> ItemListById

exception ItemUnavailable {
    1: string message;
}

service Item {
	ItemList getItem(1:i16 pos,2:i16 count),
	TItem deleteItem(1:i16 id ),
	TItem addItem(1:i16 id, 2:TItem item),
	TItem editItems(1: i16 id, 2:TItem item),
	TUser addUser(1:TUser user),
	TItem getItemByUID(1:i16 id)
	
	
}



