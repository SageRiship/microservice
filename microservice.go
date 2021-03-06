package microservice

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var UserCollection *mongo.Collection
var Client *mongo.Client
var WalletCollection *mongo.Collection
var WalletTransactionCollection *mongo.Collection

type User struct {
	Id          primitive.ObjectID `json:"_id" bson:"_id"`
	UserId      int                `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Uname       string             `json:"uname,omitempty" bson:"uname,omitempty"`
	DisplayName string             `json:"display_name,omitempty" bson:"display_name,omitempty"`
	UserRole    string             `json:"user_role,omitempty" bson:"user_role,omitempty"`
	Password    string             `json:"password,omitempty" bson:"password,omitempty"`
	Phone       []Phone            `json:"phone,omitempty" bson:"phone,omitempty"`
	Address     []Address          `json:"address,omitempty" bson:"address,omitempty"`
	FriendsList []FriendsList      `json:"friends_list,omitempty" bson:"friends_list,omitempty"`
	CreatedBy   string             `json:"created_by,omitempty" bson:"created_by,omitempty"`
	CreatedOn   primitive.DateTime `json:"created_on,omitempty" bson:"created_on,omitempty"`
	UpdatedBy   string             `json:"updated_by,omitempty" bson:"updated_by,omitempty"`
	UpdatedOn   primitive.DateTime `json:"updated_on,omitempty" bson:"updated_on,omitempty"`
}

type Book struct {
	Name      string `json:"name" bson:"name"`
	Author    string `json:"author" bson:"author"`
	PageCount int    `json:"page_count" bson:"page_count"`
}

type Phone struct {
	Number  string `json:"num,omitempty" bson:"num,omitempty"`
	Primary bool   `json:"primary" bson:"primary"`
}

type Address struct {
	Street  string `json:"street,omitempty" bson:"street,omitempty"`
	City    string `json:"city,omitempty" bson:"city,omitempty"`
	State   string `json:"state,omitempty" bson:"state,omitempty"`
	Country string `json:"country,omitempty" bson:"country,omitempty"`
}

type FriendsList struct {
	UserId          int      `json:"user_id" bson:"user_id"`
	Uname           string   `json:"uname" bson:"uname"`
	Blocked         bool     `json:"blocked" bson:"blocked"`
	BlockedForGames []string `json:"blocked_for_games" bson:"blocked_for_games"`
}

var (
	Ctx = context.TODO()
	Db  *mongo.Database
)

func init() {
	host := "10.102.78.95"
	//	host := "localhost"
	port := "27017"
	connectionURI := "mongodb://" + host + ":" + port + "/"
	clientOptions := options.Client().ApplyURI(connectionURI)
	Client, _ = mongo.Connect(Ctx, clientOptions)

	fmt.Println("Database : ", connectionURI)
	fmt.Println("MongoDB connected Successfully...")

	//Db = client.Database("UserData")              //Local database on PC
	Db = Client.Database("MVRDB")
	UserCollection = Db.Collection("User")
	WalletCollection = Db.Collection("Wallet")
	WalletTransactionCollection = Db.Collection("Wallet_transaction")

	//fmt.Println(UserCollection)

}
func GetAllUsersService() ([]User, error) {

	var users []User
	//collection = client.Database(dbname).Collection(colname)
	//ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := UserCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}

	//defer cursor.Close(ctx)
	for cursor.Next(context.Background()) {
		var user User
		cursor.Decode(&user)
		users = append(users, user)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return users, nil

}

func GetAllUsers(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var users []User
	//collection = client.Database(dbname).Collection(colname)
	//ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	//cursor, err := collection.UserCollection.Find(context.Background(), bson.M{})
	users, err := GetAllUsersService()
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	json.NewEncoder(response).Encode(users)

}

// func main() {
// 	router := mux.NewRouter()

// 	//	User Endpoints
// 	router.HandleFunc("/api/user", GetAllUsers).Methods("GET")

// 	fmt.Println("Server is getting started...")
// 	fmt.Println("Listening at port 4000 ...")
// 	log.Fatal(http.ListenAndServe(":4000", router))
// }
