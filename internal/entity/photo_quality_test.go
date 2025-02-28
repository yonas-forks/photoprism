package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPhoto_QualityScore(t *testing.T) {
	t.Run("PhotoFixture19800101_000002_D640C559", func(t *testing.T) {
		assert.Equal(t, 3, PhotoFixtures.Pointer("19800101_000002_D640C559").QualityScore())
	})
	t.Run("PhotoFixturePhoto01 - favorite true - taken at before 2008", func(t *testing.T) {
		assert.Equal(t, 7, PhotoFixtures.Pointer("Photo01").QualityScore())
	})
	t.Run("PhotoFixturePhoto06 - taken at after 2012 - resolution 2", func(t *testing.T) {
		assert.Equal(t, 3, PhotoFixtures.Pointer("Photo06").QualityScore())
	})
	t.Run("PhotoFixturePhoto07 - score < 3 bit edited", func(t *testing.T) {
		assert.Equal(t, 3, PhotoFixtures.Pointer("Photo07").QualityScore())
	})
	t.Run("PhotoFixturePhoto15 - description with non-photographic", func(t *testing.T) {
		assert.Equal(t, 2, PhotoFixtures.Pointer("Photo15").QualityScore())
	})
}

func TestPhoto_UpdateQuality(t *testing.T) {
	t.Run("Hidden", func(t *testing.T) {
		p := &Photo{PhotoQuality: -1}
		err := p.UpdateQuality()
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, -1, p.PhotoQuality)
	})
	t.Run("Favorite", func(t *testing.T) {
		p := &Photo{PhotoQuality: 0, PhotoFavorite: true}
		err := p.UpdateQuality()
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, 4, p.PhotoQuality)
	})
}

func TestPhoto_IsNonPhotographic(t *testing.T) {
	t.Run("Raw", func(t *testing.T) {
		m := PhotoFixtures.Get("Photo01")
		assert.False(t, m.IsNonPhotographic())
	})
	t.Run("Image", func(t *testing.T) {
		m := PhotoFixtures.Get("Photo04")
		assert.False(t, m.IsNonPhotographic())
	})
	t.Run("Video", func(t *testing.T) {
		m := PhotoFixtures.Get("Photo10")
		assert.False(t, m.IsNonPhotographic())
	})
	t.Run("Animated", func(t *testing.T) {
		m := PhotoFixtures.Get("Photo52")
		assert.True(t, m.IsNonPhotographic())
	})
}
