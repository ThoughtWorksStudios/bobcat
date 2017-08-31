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

		InterpretExpectsSuccess(t, NewNestedEmitter, filename)

		withParsed(t, filename, func(res jsoniter.Any, raw string) {
			AssertEqual(t, 2, res.Get("Review").Size(), "Should have created 2 reviews")

			eachElem(res.Get("Review"), func(review jsoniter.Any, i int) {
				AssertEqual(t, "Rating", review.Get("rating", "$type").ToString(), "Should have nested rating object")
				AssertEqual(t, review.Get("$id").ToString(), review.Get("rating", "$parent").ToString(), "Nested objects are linked via $id <-> $parent")
			})
		})
	})
}

func TestWithFlatEmitter(t *testing.T) {
	withBasename(func(basename string) {
		filename := basename + ".json"
		defer cleanup(filename)()

		InterpretExpectsSuccess(t, NewFlatEmitter, filename)

		withParsed(t, filename, func(res jsoniter.Any, raw string) {
			AssertEqual(t, 4, res.Size(), "Should have created 4 entities in a flat array")

			AssertEqual(t, res.Get(0, "$id").ToString(), res.Get(1, "rating").ToString(), "Relationship expressed as Review.rating -> Rating.$id")
			AssertEqual(t, res.Get(2, "$id").ToString(), res.Get(3, "rating").ToString(), "Relationship expressed as Review.rating -> Rating.$id")

			AssertEqual(t, "Rating", res.Get(0, "$type").ToString(), "Failed entity $type check")
			AssertEqual(t, "Rating", res.Get(2, "$type").ToString(), "Failed entity $type check")

			AssertEqual(t, "Review", res.Get(1, "$type").ToString(), "Failed entity $type check")
			AssertEqual(t, "Review", res.Get(3, "$type").ToString(), "Failed entity $type check")
		})
	})
}

func TestWithSplitEmitter(t *testing.T) {
	withBasename(func(basename string) {
		filename := basename + ".json"
		defer cleanup(basename + "-Rating.json")()
		defer cleanup(basename + "-Review.json")()

		InterpretExpectsSuccess(t, NewSplitEmitter, filename)

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
		t.Fatalf("Failed to create emitter with factory: %v", factory)
		return
	} else {
		_, err = New(emitter, false).LoadFile(TEST_FILE, NewRootScope())
		emitter.Finalize()
		AssertNil(t, err, "Should not have received an error interpreting %q", TEST_FILE)
	}
}
