package constants

type APIKey int

const (
	GetBookError APIKey = iota
	CreateBookError
	UpdateBookByIdError
	UpdateBooksByAuthorNameError
	DeleteBookError
	KeyspaceFailed
	TableFailed
	FailedToInsetData
)

var APIFailureMessages = map[APIKey]string{
	GetBookError:                 "Error while fetching book",
	CreateBookError:              "Error while creating books",
	UpdateBookByIdError:          "Error while updating books by ID",
	UpdateBooksByAuthorNameError: "Error while updating books by author name",
	DeleteBookError:              "Error while deleting book",
	KeyspaceFailed:               "Failed to create keyspace",
	TableFailed:                  "Failed to create table",
	FailedToInsetData:            "Failed to insert data",
}

const (
	BookFetched APIKey = iota
	BooksFetched
	BookUpdatedById
	BookUpdatedByAuthorName
	BookDeleted
	BookCreated
)

var APISuccessMessages = map[APIKey]string{
	BookFetched:             "Book fetched successfully",
	BooksFetched:            "Books fetched successfully",
	BookUpdatedById:         "Book successfully updated by Id",
	BookUpdatedByAuthorName: "Book successfully updated by authorName",
	BookDeleted:             "Book deleted successfully",
	BookCreated:             "Book created successfully",
}

const (
	ServerURL APIKey = iota
	DatabaseConnected
	DatabaseConnectionFailed
)

var StringConstantsMessages = map[APIKey]string{
	ServerURL:                "mongodb://localhost:27017",
	DatabaseConnected:        "Database connected successfully",
	DatabaseConnectionFailed: "failed to connect to the database:",
}

const (
	DatabaseName APIKey = iota
	CollectionName
)

var DatabaseInfo = map[APIKey]string{
	DatabaseName:   "book",
	CollectionName: "booklist",
}
