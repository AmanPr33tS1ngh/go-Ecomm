package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID				primitive.ObjectID		`json:"_id" bson:"_id"`
	FirstName		*string					`json:"firstName"`
	LastName		*string					`json:"lastName"`
	Password		*string					`json:"password"`
	Email			*string					`json:"email"`
	Phone			*string					`json:"phone"`
	Token			*string					`json:"token"`
	RefreshToken	*string					`json:"refreshToken"`
	CreatedAt		time.Time				`json:"createdAt"`
	UpdatedAt		time.Time				`json:"updatedAt"`
	UserId			string					`json:"userId"`
	UserCart		[]ProductUser			`json:"userCart" bson:"userCart"`
	AddressDetails	[]Address				`json:"address" bson:"address"`
	OrderStatus		[]Order					`json:"orders" bson:"orders"`
}

type Product struct {
	ProductId		primitive.ObjectID		`json:"_id" bson:"_id"`
	ProductName		*string					`json:"productName"`
	Price			*uint64					`json:"price"`
	Rating			*uint8					`json:"rating"`
	Image			*string					`json:"image"`
}


type ProductUser struct {
	ProductUserId	primitive.ObjectID		`json:"_id" bson:"_id"`
	ProductName		*string					`json:"productName"`
	Price			int						`json:"price"`
	Rating			*uint					`json:"rating"`
	Image			*string					`json:"image"`
}

type Address struct {
	AddressId		primitive.ObjectID		`json:"_id" bson:"_id"`
	House			*string					`json:"house"`
	Street			*string					`json:"street"`
	City			*string					`json:"city"`
	PinCode			*string					`json:"pinCode"`
}

type Order struct {
	OrderId			primitive.ObjectID		`json:"_id" bson:"_id"`
	OrderCart		[]ProductUser			`json:"orders"`
	OrderedAt		time.Time				`json:"orderedAt"`
	Price			int						`json:"price"`
	Discount		*int					`json:"discount"`
	PaymentMethod	Payment					`json:"payment"`
}

type Payment struct {
	Digital bool
	COD bool
}