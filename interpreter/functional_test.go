package interpreter

import (
	. "github.com/ThoughtWorksStudios/bobcat/emitter"
	. "github.com/ThoughtWorksStudios/bobcat/test_helpers"
	"github.com/json-iterator/go"
	"testing"
)

const TEST_FILE = "testdata/reviews.lang"

func TestWithNestedEmitter(t *testing.T) {
	withBasename(func(basename string) {
		filename := basename + ".json"
		defer cleanup(filename)()

		InterpretExpectsSuccess(t, NestedEmitterForFile, filename)

		withParsed(t, filename, func(res jsoniter.Any, raw string) {
			AssertEqual(t, 3, res.Size(), "Should have created 3 top-level entities")
			AssertEqual(t, "Author", res.Get(0, "$type").ToString(), "First entity is an Author")
			AssertEqual(t, "Review", res.Get(1, "$type").ToString(), "Second entity is a Review")
			AssertEqual(t, "Review", res.Get(2, "$type").ToString(), "Third entity is a Review")

			authorId := res.Get(0, "$id").ToString()

			eachElem(res, func(elem jsoniter.Any, i int) {
				switch elem.Get("$type").ToString() {
				case "Author":
					AssertNotEqual(t, "", elem.Get("name").ToString(), "Should have generated Author.name")
				case "Review":
					AssertEqual(t, authorId, elem.Get("author").ToString(), "Should have reference to Author")
					AssertEqual(t, "Rating", elem.Get("rating", "$type").ToString(), "Should have nested rating object")
					AssertEqual(t, elem.Get("$id").ToString(), elem.Get("rating", "$parent").ToString(), "Nested objects are linked via $id <-> $parent")
				default:
					t.Errorf("Unexpected entity type %q!", elem.Get("$type").ToString())
				}
			})
		})
	})
}

func TestWithFlatEmitter(t *testing.T) {
	withBasename(func(basename string) {
		filename := basename + ".json"
		defer cleanup(filename)()

		InterpretExpectsSuccess(t, FlatEmitterForFile, filename)

		withParsed(t, filename, func(res jsoniter.Any, raw string) {
			AssertEqual(t, 5, res.Size(), "Should have created 4 entities in a flat array")

			AssertEqual(t, "Author", res.Get(0, "$type").ToString(), "Failed entity $type check")

			AssertEqual(t, "Rating", res.Get(1, "$type").ToString(), "Failed entity $type check")
			AssertEqual(t, "Rating", res.Get(3, "$type").ToString(), "Failed entity $type check")

			AssertEqual(t, "Review", res.Get(2, "$type").ToString(), "Failed entity $type check")
			AssertEqual(t, "Review", res.Get(4, "$type").ToString(), "Failed entity $type check")

			AssertEqual(t, res.Get(0, "$id").ToString(), res.Get(2, "author").ToString(), "Relationship expressed as Review.author -> Author.$id")
			AssertEqual(t, res.Get(0, "$id").ToString(), res.Get(4, "author").ToString(), "Relationship expressed as Review.author -> Author.$id")

			AssertEqual(t, res.Get(1, "$id").ToString(), res.Get(2, "rating").ToString(), "Relationship expressed as Review.rating -> Rating.$id")
			AssertEqual(t, res.Get(3, "$id").ToString(), res.Get(4, "rating").ToString(), "Relationship expressed as Review.rating -> Rating.$id")
		})
	})
}

func TestWithSplitEmitter(t *testing.T) {
	withBasename(func(basename string) {
		filename := basename + ".json"
		defer cleanup(basename + "-Rating.json")()
		defer cleanup(basename + "-Review.json")()

		InterpretExpectsSuccess(t, SplitEmitterForFile, filename)

		withParsed(t, basename+"-Rating.json", func(ratings jsoniter.Any, raw1 string) {
			withParsed(t, basename+"-Review.json", func(reviews jsoniter.Any, raw2 string) {
				AssertEqual(t, 2, ratings.Size(), "Should have created 2 ratings in a flat array")
				AssertEqual(t, 2, reviews.Size(), "Should have created 2 reviews in a flat array")

				eachElem(ratings, func(rating jsoniter.Any, i int) {
					AssertEqual(t, reviews.Get(i, "rating").ToString(), rating.Get("$id").ToString(), "Relationships expressed as Review.rating -> Rating.$id")
				})
			})
		})
	})
}

func InterpretExpectsSuccess(t *testing.T, factory EmitterFactory, outFilename string) {
	if emitter, err := factory(outFilename); err != nil {
		t.Fatalf("Failed to create emitter with factory %v; err: %v", factory, err)
		return
	} else {
		defer emitter.Finalize()
		emitter.Init()
		if _, err = New(emitter, false).LoadFile(TEST_FILE, NewRootScope()); err != nil {
			t.Fatalf("Should not have received an error interpreting %q, but got: %v", TEST_FILE, err)
		}
	}
}
