package chromem_test

import (
	"context"
	"testing"

	"github.com/philippgille/chromem-go"
)

func TestDB_ListCollections(t *testing.T) {
	// Values in the collection
	name := "test"
	metadata := map[string]string{"foo": "bar"}
	embeddingFunc := func(_ context.Context, _ string) ([]float32, error) {
		return []float32{-0.1, 0.1, 0.2}, nil
	}

	// Create initial collection
	db := chromem.NewDB()
	// We ignore the return value. CreateCollection is tested elsewhere.
	_ = db.CreateCollection(name, metadata, embeddingFunc)

	// List collections
	res := db.ListCollections()

	// Check expectations
	if len(res) != 1 {
		t.Error("expected 1 collection, got", len(res))
	}
	c, ok := res[name]
	if !ok {
		t.Error("expected collection", name, "not found")
	}
	if c.Name != name {
		t.Error("expected name", name, "got", c.Name)
	}
	if len(c.Metadata) != 1 {
		t.Error("expected 1 metadata, got", len(c.Metadata))
	}
	if c.Metadata["foo"] != "bar" {
		t.Error("expected metadata", metadata, "got", c.Metadata)
	}

	// And it should be a copy
	res["foo"] = &chromem.Collection{}
	if len(db.ListCollections()) != 1 {
		t.Error("expected 1 collection, got", len(db.ListCollections()))
	}
}

func TestDB_GetCollection(t *testing.T) {
	// Values in the collection
	name := "test"
	metadata := map[string]string{"foo": "bar"}
	embeddingFunc := func(_ context.Context, _ string) ([]float32, error) {
		return []float32{-0.1, 0.1, 0.2}, nil
	}

	// Create initial collection
	db := chromem.NewDB()
	// We ignore the return value. CreateCollection is tested elsewhere.
	_ = db.CreateCollection(name, metadata, embeddingFunc)

	// Get collection
	c := db.GetCollection(name)

	// Check expectations
	if c.Name != name {
		t.Error("expected name", name, "got", c.Name)
	}
	if len(c.Metadata) != 1 {
		t.Error("expected 1 metadata, got", len(c.Metadata))
	}
	if c.Metadata["foo"] != "bar" {
		t.Error("expected metadata", metadata, "got", c.Metadata)
	}
	// TODO: Check documents content as soon as we have access to them
	// TODO: Same for the EmbeddingFunc
	// TODO: Check documents map being a copy as soon as we have access to it
}

func TestDB_DeleteCollection(t *testing.T) {
	// Values in the collection
	name := "test"
	metadata := map[string]string{"foo": "bar"}
	embeddingFunc := func(_ context.Context, _ string) ([]float32, error) {
		return []float32{-0.1, 0.1, 0.2}, nil
	}

	// Create initial collection
	db := chromem.NewDB()
	// We ignore the return value. CreateCollection is tested elsewhere.
	_ = db.CreateCollection(name, metadata, embeddingFunc)

	// Delete collection
	db.DeleteCollection(name)

	// Check expectations
	// We don't have access to the documents field, but we can rely on DB.ListCollections()
	// because it's tested elsewhere.
	if len(db.ListCollections()) != 0 {
		t.Error("expected 0 collections, got", len(db.ListCollections()))
	}
}