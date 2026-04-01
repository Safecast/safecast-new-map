// Local text embedding using the feature hashing trick.
// No external API or API key required.
//
// Each text is tokenized into unigrams + bigrams, each term is hashed into a
// 512-dimensional float32 vector via FNV-32a, and the vector is L2-normalised.
// Cosine similarity on the resulting vectors gives a good approximation of
// semantic overlap for domain-specific short queries (radiation, locations, etc.).
//
// If query quality proves insufficient, swap getEmbedding() for a Voyage AI
// call (voyageai.com, model: voyage-3-lite) — same []float32 return type,
// no other changes required.

package main

import (
	"context"
	"math"
	"regexp"
	"strings"
)

const embDims = 512

// stopWords contains common English words that carry no topical signal.
var stopWords = map[string]bool{
	"a": true, "an": true, "the": true, "and": true, "or": true, "but": true,
	"in": true, "on": true, "at": true, "to": true, "for": true, "of": true,
	"with": true, "is": true, "are": true, "was": true, "were": true,
	"be": true, "been": true, "being": true, "have": true, "has": true,
	"had": true, "do": true, "does": true, "did": true, "will": true,
	"would": true, "could": true, "should": true, "may": true, "might": true,
	"can": true, "this": true, "that": true, "these": true, "those": true,
	"it": true, "its": true, "from": true, "by": true, "about": true,
	"what": true, "why": true, "how": true, "when": true, "where": true,
	"which": true, "who": true, "i": true, "me": true, "my": true,
	"you": true, "your": true, "we": true, "us": true, "our": true,
}

var nonAlnum = regexp.MustCompile(`[^a-z0-9]+`)

// tokenize lowercases text, splits on non-alphanumeric characters, and
// returns unigrams + bigrams, skipping stop words for unigrams.
func tokenize(text string) []string {
	lower := strings.ToLower(text)
	parts := nonAlnum.Split(lower, -1)

	var words []string
	for _, w := range parts {
		if len(w) > 1 {
			words = append(words, w)
		}
	}

	var tokens []string
	for i, w := range words {
		if !stopWords[w] {
			tokens = append(tokens, w)
		}
		// Bigrams include stop words so phrases like "near the reactor" →
		// "the_reactor" still captures the noun.
		if i+1 < len(words) {
			tokens = append(tokens, w+"_"+words[i+1])
		}
	}
	return tokens
}

// fnv32a computes the FNV-1a 32-bit hash of a string.
func fnv32a(s string) uint32 {
	h := uint32(2166136261)
	for i := 0; i < len(s); i++ {
		h ^= uint32(s[i])
		h *= 16777619
	}
	return h
}

// getEmbedding returns a normalised embDims-dimensional feature-hash vector
// for the given text. ctx is accepted for API-compatible signature; it is
// unused in this local implementation. Never returns an error.
func getEmbedding(_ context.Context, text string) ([]float32, error) {
	tokens := tokenize(text)
	vec := make([]float32, embDims)
	for _, tok := range tokens {
		idx := fnv32a(tok) % embDims
		vec[idx] += 1.0
	}
	l2Normalize(vec)
	return vec, nil
}

// l2Normalize divides the vector by its L2 norm in-place.
func l2Normalize(v []float32) {
	var sum float64
	for _, x := range v {
		sum += float64(x) * float64(x)
	}
	if sum == 0 {
		return
	}
	norm := float32(math.Sqrt(sum))
	for i := range v {
		v[i] /= norm
	}
}

// cosineSimilarity computes cosine similarity between two float32 slices.
// Both vectors should already be L2-normalised (dot product == cosine similarity).
func cosineSimilarity(a, b []float32) float32 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}
	var dot float64
	for i := range a {
		dot += float64(a[i]) * float64(b[i])
	}
	return float32(dot)
}
