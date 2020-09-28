package fileList

type repository struct {
	db         *Mongo
	collection Collection
}

func NewFileListRepository(config Configuration) (*repository, error) {
	database, err := Connect(config)

	files := database.AddCollection(config.GetCollectionFiles())

	return &repository{
		db:         database,
		collection: Collection(files),
	}, err
}
