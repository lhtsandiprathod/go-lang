package database

import (
	"bookname/constants"
	"bookname/graph/model"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type DB struct{ client *mongo.Client }
type BookWithID struct {
	ID          primitive.ObjectID `json:"_id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Author      string             `json:"author"`
	AddedOn     float64            `json:"addedOn"`
	// Other fields...
}

func Connect() *DB {

	client, err := mongo.NewClient(options.Client().ApplyURI(constants.StringConstantsMessages[constants.ServerURL]))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	log.Println(constants.StringConstantsMessages[constants.DatabaseConnected])
	return &DB{
		client: client,
	}
}

func (db *DB) GetBookById(id string) *model.BookListing {
	bookCollec := db.client.Database(constants.DatabaseInfo[constants.DatabaseName]).Collection(constants.DatabaseInfo[constants.CollectionName])
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_id, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": _id}
	var bookListing model.BookListing
	err := bookCollec.FindOne(ctx, filter).Decode(&bookListing)

	if err != nil {
		log.Println(constants.APIFailureMessages[constants.GetBookError])
	}
	bookListing.ID = _id.Hex()

	log.Println(constants.APISuccessMessages[constants.BookFetched])
	return &bookListing

}

func (db *DB) GetBooks() []*model.BookListing {
	bookCollec := db.client.Database("book").Collection("booklist")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var bookListings []*model.BookListing
	var bookListingsDB []*model.BookListingDB

	cursor, err := bookCollec.Find(ctx, bson.D{})
	if err != nil {
		log.Println(constants.APIFailureMessages[constants.GetBookError], err)
	}

	if err = cursor.All(context.TODO(), &bookListingsDB); err != nil {
		log.Println(constants.APIFailureMessages[constants.GetBookError], err)
	}

	for _, bookDB := range bookListingsDB {
		bookListings = append(bookListings, &model.BookListing{
			ID:          bookDB.ID.Hex(),
			Title:       bookDB.Title,
			Description: bookDB.Description,
			Author:      bookDB.Author,
			AddedOn:     bookDB.AddedOn,
		})
	}
	log.Println(constants.APISuccessMessages[constants.BooksFetched])
	return bookListings
}

func (db *DB) CreateBookListing(bookInfo model.CreateBookListingInput) *model.BookListing {
	bookCollec := db.client.Database("book").Collection("booklist")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	//currentTime := strconv.FormatInt(time.Now().Unix(), 10)
	currentTime := time.Now()
	currentTimeFloat := float64(currentTime.Unix())
	inserg, err := bookCollec.InsertOne(ctx, bson.M{"title": bookInfo.Title, "description": bookInfo.Description, "author": bookInfo.Author, "addedOn": currentTimeFloat})

	if err != nil {
		log.Println(constants.APIFailureMessages[constants.CreateBookError], err)
	}

	// createdTime := strconv.FormatInt(currentTime, 10)
	insertedID := inserg.InsertedID.(primitive.ObjectID).Hex()
	returnJobListing := model.BookListing{ID: insertedID, Title: bookInfo.Title, Author: bookInfo.Author, Description: bookInfo.Description, AddedOn: currentTimeFloat}

	log.Println(constants.APISuccessMessages[constants.BookCreated])
	return &returnJobListing
}

func (db *DB) UpdateBooks(bookId string, bookInfo model.UpdateBookListingInput) *model.BookListing {

	jobCollec := db.client.Database("book").Collection("booklist")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	currentTime := time.Now()
	currentTimeFloat := float64(currentTime.Unix())
	updateBookInfo := bson.M{}

	if bookInfo.Title != nil {
		updateBookInfo["title"] = bookInfo.Title
	}
	if bookInfo.Description != nil {
		updateBookInfo["description"] = bookInfo.Description
	}
	if bookInfo.Author != nil {
		updateBookInfo["author"] = bookInfo.Author
	}
	updateBookInfo["addedOn"] = currentTimeFloat

	_id, _ := primitive.ObjectIDFromHex(bookId)
	filter := bson.M{"_id": _id}
	update := bson.M{"$set": updateBookInfo}

	results := jobCollec.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(1))

	var bookListing model.BookListing

	bookListing.ID = _id.Hex()

	if err := results.Decode(&bookListing); err != nil {
		log.Println(constants.APIFailureMessages[constants.UpdateBookByIdError], err)
	}

	log.Println(constants.APISuccessMessages[constants.BookUpdatedById])
	return &bookListing
}

func (db *DB) DeleteBookByID(BookId string) *model.DeleteBook {
	BookCollec := db.client.Database("book").Collection("booklist")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_id, _ := primitive.ObjectIDFromHex(BookId)
	filter := bson.M{"_id": _id}
	_, err := BookCollec.DeleteOne(ctx, filter)
	if err != nil {
		log.Println(constants.APIFailureMessages[constants.DeleteBookError], err)
	}

	log.Println(constants.APISuccessMessages[constants.BookDeleted])
	return &model.DeleteBook{DeletedBookID: BookId}
}

// func (db *DB) GetBookByAuthorName(author string) []*model.BookListing {
// 	bookCollec := db.client.Database("book").Collection("booklist")
// 	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
// 	defer cancel()

// 	filter := bson.M{"author": author}
// 	cursor, err := bookCollec.Find(ctx, filter)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer cursor.Close(ctx)

// 	var bookListings []*model.BookListing

// 	for cursor.Next(ctx) {
// 		var bookListing model.BookListing
// 		err := cursor.Decode(&bookListing)
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		// Convert ObjectID to string representation
// 		bookListing.ID = bookListing.ID

// 		bookListings = append(bookListings, &bookListing)
// 	}

// 	if err := cursor.Err(); err != nil {
// 		log.Fatal(err)
// 	}

// 	return bookListings
// }

func (db *DB) GetBookByAuthorName(author string) []*model.BookListing {
	bookCollec := db.client.Database("book").Collection("booklist")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	filter := bson.M{"author": author}
	cursor, err := bookCollec.Find(ctx, filter)
	if err != nil {
		log.Println(constants.APIFailureMessages[constants.GetBookError], err)
	}
	defer cursor.Close(ctx)

	var bookListings []*model.BookListing
	for cursor.Next(ctx) {
		var bookListing model.BookListing
		err := cursor.Decode(&bookListing)
		if err != nil {
			log.Println(constants.APIFailureMessages[constants.GetBookError], err)
		}
		bookListings = append(bookListings, &bookListing)
	}
	if err := cursor.Err(); err != nil {
		log.Println(constants.APIFailureMessages[constants.GetBookError], err)
	}

	log.Println(constants.APISuccessMessages[constants.BookFetched])
	return bookListings
}
