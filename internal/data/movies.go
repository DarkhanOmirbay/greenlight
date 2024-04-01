package data

import (
	"time"
)

type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"` // - (hyphen)directive use for hiding field
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`
	Runtime   Runtime   `json:"runtime,omitempty,string"` // string directive use for represent field in JSON STRING
	// key-word omitempty uses for hide empty field
	Genres  []string `json:"genres,omitempty"`
	Version int32    `json:"version"`
}

//type Movie struct {
//	ID        int64
//	CreatedAt time.Time
//	Title     string
//	Year      int32
//	Runtime   int32
//	Genres    []string
//	Version   int32
//}

//func (m Movie) MarshalJSON() ([]byte, error) {
//
//	var runtime string
//
//	if m.Runtime != 0 {
//		runtime = fmt.Sprintf("%d mins", m.Runtime)
//	}
//	aux := struct {
//		ID      int64    `json:"id"`
//		Title   string   `json:"title"`
//		Year    int32    `json:"year,omitempty"`
//		Runtime string   `json:"runtime,omitempty"` // This is a string.
//		Genres  []string `json:"genres,omitempty"`
//		Version int32    `json:"version"`
//	}{
//		ID:      m.ID,
//		Title:   m.Title,
//		Year:    m.Year,
//		Runtime: runtime, // Note that we assign the value from the runtime variable here.
//		Genres:  m.Genres,
//		Version: m.Version,
//	}
//	return json.Marshal(aux)
//}

//func (m Movie) MarshalJSON() ([]byte, error) {
//	// Create a variable holding the custom runtime string, just like before.
//	var runtime string
//	if m.Runtime != 0 {
//		runtime = fmt.Sprintf("%d mins", m.Runtime)
//	}
//	// Define a MovieAlias type which has the underlying type Movie. Due to the way that
//	// Go handles type definitions (https://golang.org/ref/spec#Type_definitions) the
//	// MovieAlias type will contain all the fields that our Movie struct has but,
//	// importantly, none of the methods.
//	type MovieAlias Movie
//	// Embed the MovieAlias type inside the anonymous struct, along with a Runtime field
//	// that has the type string and the necessary struct tags. It's important that we
//	// embed the MovieAlias type here, rather than the Movie type directly, to avoid
//	// inheriting the MarshalJSON() method of the Movie type (which would result in an
//	// infinite loop during encoding).
//	aux := struct {
//		MovieAlias
//		Runtime string `json:"runtime,omitempty"`
//	}{
//		MovieAlias: MovieAlias(m),
//		Runtime: runtime,
//	}
//	return json.Marshal(aux)
//}
