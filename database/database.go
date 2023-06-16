package database

import (
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

func Connect() *DB {

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
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

	return &DB{
		client: client,
	}
}

func (db *DB) GetBookById(id string) *model.BookListing {
	bookCollec := db.client.Database("book").Collection("booklist")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_id, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": _id}
	var bookListing model.BookListing
	err := bookCollec.FindOne(ctx, filter).Decode(&bookListing)
	if err != nil {
		log.Fatal(err)
	}
	return &bookListing
}

func (db *DB) GetBooks() []*model.BookListing {
	bookCollec := db.client.Database("book").Collection("booklist")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var bookListings []*model.BookListing
	cursor, err := bookCollec.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	if err = cursor.All(context.TODO(), &bookListings); err != nil {
		panic(err)
	}

	return bookListings
}

func (db *DB) CreateBookListing(bookInfo model.CreateBookListingInput) *model.BookListing {
	bookCollec := db.client.Database("book").Collection("booklist")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	inserg, err := bookCollec.InsertOne(ctx, bson.M{"title": bookInfo.Title, "description": bookInfo.Description, "author": bookInfo.Author})

	if err != nil {
		log.Fatal(err)
	}

	insertedID := inserg.InsertedID.(primitive.ObjectID).Hex()
	returnJobListing := model.BookListing{ID: insertedID, Title: bookInfo.Title, Author: bookInfo.Author, Description: bookInfo.Description}
	return &returnJobListing
}

func (db *DB) UpdateBooks(bookId string, bookInfo model.UpdateBookListingInput) *model.BookListing {
	jobCollec := db.client.Database("book").Collection("booklist")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	updateBookInfo := bson.M{}

	if bookInfo.Title != nil {
		updateBookInfo["title"] = bookInfo.Title
	}
	if bookInfo.Description != nil {
		updateBookInfo["description"] = bookInfo.Description
	}
	if bookInfo.Author != nil {
		updateBookInfo["url"] = bookInfo.Author
	}

	_id, _ := primitive.ObjectIDFromHex(bookId)
	filter := bson.M{"_id": _id}
	update := bson.M{"$set": updateBookInfo}

	results := jobCollec.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(1))

	var bookListing model.BookListing

	if err := results.Decode(&bookListing); err != nil {
		log.Fatal(err)
	}

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
		log.Fatal(err)
	}
	return &model.DeleteBook{DeletedBookID: BookId}
}

func (db *DB) GetBookByAuthorName(author string) []*model.BookListing {
	bookCollec := db.client.Database("book").Collection("booklist")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	filter := bson.M{"author": author}
	cursor, err := bookCollec.Find(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	var bookListings []*model.BookListing
	for cursor.Next(ctx) {
		var bookListing model.BookListing
		err := cursor.Decode(&bookListing)
		if err != nil {
			log.Fatal(err)
		}
		bookListings = append(bookListings, &bookListing)
	}
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}
	return bookListings
}
