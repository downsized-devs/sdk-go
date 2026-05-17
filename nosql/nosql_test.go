package nosql

import (
	"context"
	"testing"

	"github.com/downsized-devs/sdk-go/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

// newMockDB constructs a *mongoDB backed by mtest's mock client so we can drive
// the operation methods without a real MongoDB server.
func newMockDB(mt *mtest.T) *mongoDB {
	return &mongoDB{
		client: mt.Client,
		cfg:    Config{DB: "test"},
		log:    logger.Init(logger.Config{}),
	}
}

func TestFind(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("decodes returned documents", func(mt *mtest.T) {
		db := newMockDB(mt)
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.items", mtest.FirstBatch,
			bson.D{{Key: "name", Value: "alpha"}},
			bson.D{{Key: "name", Value: "beta"}},
		))

		type item struct {
			Name string `bson:"name"`
		}
		var got []item
		err := db.Find(context.Background(), "items", &got, bson.D{})
		require.NoError(t, err)
		assert.Equal(t, []item{{"alpha"}, {"beta"}}, got)
	})

	mt.Run("wraps server errors", func(mt *mtest.T) {
		db := newMockDB(mt)
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    11000,
			Message: "boom",
			Name:    "BoomError",
		}))

		var got []bson.M
		err := db.Find(context.Background(), "items", &got, bson.D{})
		assert.Error(t, err)
	})
}

func TestFindOne(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("decodes one document", func(mt *mtest.T) {
		db := newMockDB(mt)
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.items", mtest.FirstBatch,
			bson.D{{Key: "name", Value: "alpha"}},
		))

		type item struct {
			Name string `bson:"name"`
		}
		var got item
		err := db.FindOne(context.Background(), "items", &got, bson.D{})
		require.NoError(t, err)
		assert.Equal(t, "alpha", got.Name)
	})

	mt.Run("returns error for not-found", func(mt *mtest.T) {
		db := newMockDB(mt)
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.items", mtest.FirstBatch))
		var got bson.M
		err := db.FindOne(context.Background(), "items", &got, bson.D{})
		assert.Error(t, err)
	})
}

func TestInsertOne(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("returns InsertOneResult", func(mt *mtest.T) {
		db := newMockDB(mt)
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		res, err := db.InsertOne(context.Background(), "items", bson.M{"name": "alpha"})
		require.NoError(t, err)
		assert.NotNil(t, res)
	})

	mt.Run("wraps server errors", func(mt *mtest.T) {
		db := newMockDB(mt)
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    11000,
			Message: "duplicate",
			Name:    "DuplicateKeyError",
		}))
		_, err := db.InsertOne(context.Background(), "items", bson.M{"name": "alpha"})
		assert.Error(t, err)
	})
}

func TestUpdateOne(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("returns UpdateResult", func(mt *mtest.T) {
		db := newMockDB(mt)
		mt.AddMockResponses(mtest.CreateSuccessResponse(
			bson.E{Key: "n", Value: 1},
			bson.E{Key: "nModified", Value: 1},
		))
		res, err := db.UpdateOne(context.Background(), "items",
			bson.D{{Key: "name", Value: "alpha"}},
			bson.D{{Key: "$set", Value: bson.D{{Key: "name", Value: "alpha2"}}}})
		require.NoError(t, err)
		require.NotNil(t, res)
		assert.EqualValues(t, 1, res.MatchedCount)
	})

	mt.Run("wraps server errors", func(mt *mtest.T) {
		db := newMockDB(mt)
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    1,
			Message: "fail",
			Name:    "GenericError",
		}))
		_, err := db.UpdateOne(context.Background(), "items", bson.D{}, bson.D{})
		assert.Error(t, err)
	})
}

func TestUpdateMany(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("returns UpdateResult", func(mt *mtest.T) {
		db := newMockDB(mt)
		mt.AddMockResponses(mtest.CreateSuccessResponse(
			bson.E{Key: "n", Value: 3},
			bson.E{Key: "nModified", Value: 3},
		))
		res, err := db.UpdateMany(context.Background(), "items",
			bson.D{{Key: "active", Value: false}},
			bson.D{{Key: "$set", Value: bson.D{{Key: "active", Value: true}}}})
		require.NoError(t, err)
		require.NotNil(t, res)
		assert.EqualValues(t, 3, res.MatchedCount)
	})

	mt.Run("wraps server errors", func(mt *mtest.T) {
		db := newMockDB(mt)
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code: 1, Message: "fail", Name: "GenericError",
		}))
		_, err := db.UpdateMany(context.Background(), "items", bson.D{}, bson.D{})
		assert.Error(t, err)
	})
}

// Close is exercised indirectly by mtest cleanup; testing it directly with the
// mtest mock client double-disconnects the mock and panics, so we cover the
// error-wrap branch via a real client check in integration tests instead.

func TestReplaceBindvarsWithArgs(t *testing.T) {
	t.Run("primitive ? placeholders are replaced positionally", func(t *testing.T) {
		got := replaceBindvarsWithArgs("SELECT ?, ?", 1, "two")
		assert.Contains(t, got, "1")
		assert.Contains(t, got, "two")
	})

	t.Run("named :tag placeholders are replaced from struct tags", func(t *testing.T) {
		type row struct {
			Name string `db:"name"`
			Age  int    `db:"age"`
		}
		got := replaceBindvarsWithArgs("name=:name age=:age", row{Name: "alpha", Age: 7})
		assert.Contains(t, got, "'alpha'")
		assert.Contains(t, got, "7")
	})
}
