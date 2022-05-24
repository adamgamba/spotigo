## Spotigo - COS 316 Final Project

# By Adam Gamba and Adam Ziff

This project implements a Go wrapper for the Spotify API to make it as
simple as possible for developers to access this robust API in Go while
preserving all of the powerful functionality it enables. Our primary
goal was to abstract away as much complexity as possible to enable
developers to get up and running with the API as quickly as possible,
while simultaneously giving developers the ability to seamlessly access
as much of the API’s powerful functionality as possible.

# Authentication

There are two distinct sets of API calls in the Spotify API: those that
require a user to log in to their account, and those that do not. These
different sets of API calls require different authentication paths, but
we have simplified these processes to one line of code each.

To create a Query struct (not requiring browser authentication), one
would call:

```go
client := "API client key..."
secret := "API secret key..."
query, ok := spotigo.NewUser(client, secret)
```

Creating a Query struct with the newQuery(client, secret string) method
performs the necessary authentication to enable the developer to access
all endpoints that don’t require this login flow.

Similarly, to create a User struct (requiring browser authentication),
one would call:

```go
client := "API client key..."
secret := "API secret key..."
user, ok := spotigo.NewUser(client, secret)
```

Creating a User struct with the newUser(client, secret string) method
prompts the user to login to their account; once they have done so, the
developer can then use this struct to access all API calls implemented
that require this user-level authentication.

# URI Abstraction

Each of these structs are created in one line of code, providing the
developer’s application Client Key and Secret Key (as created in the
Spotify Developer Dashboard); this format builds on the existing Go
Spotify API wrappers to further simplify the authentication process and
improve ease of use for developers. One other important element we
sought to abstract away from developers were URIs. URIs are Spotify’s
unique IDs for any piece of content on its platform: songs, artists,
albums, etc. Existing libraries only enable developers to get these
pieces of content from the API by URI; however, we leveraged Spotify’s
powerful general search endpoint to enable developers to search for any
type of content in a plain text query as well. Notably, the implicit
assumption of these functions is that the user provides enough
information in the search query string and that Spotify’s search
algorithm finds the element they’re looking for; without the precision
of the exact URI, we rely on the accuracy of this search engine in our
results.

# Functionality

With the Query struct, developers can search for albums, artists,
tracks, and playlists, all of which can be accessed either by their URI
or through a plain text search query. We’ve created structs for the
return values of each of these endpoints: the Album, Artist, Track, and
Playlist structs exactly correspond to the JSON data returned by the
API. We have implemented a few Get methods for each struct to simplify
access to what we believe to be the most important and frequently
accessed fields, while simultaneously exporting every field of the
struct so that developers are free to access any field they need for
their application. With the User struct, developers have access to a few
distinct sets of functionality. They can control playback from a user’s
account; this includes features like playing and pausing music, adding
tracks to the queue, switching playback between devices, and more. They
can get and set information on a user’s profile; this includes getting a
list of saved tracks, albums, or playlists, saving or unsaving any of
these elements, and following or unfollowing artists. Finally, they can
access all of the detailed information that Spotify saves for a given
track through the audio-analysis and audio-features endpoints.

# Safety

An important final note to add: through all these design decisions, we
sought to prioritize safety and code security. We made intentional
design choices around which fields and methods should be exported,
trying to balance code security with maximizing the amount of
flexibility developers have. This is most apparent in the User.go, where
the User struct is exported but none of its fields are, and only a few
methods are exported.

# Outcome

For this project, we feel that we were able to successfully build a
robust wrapper for the Spotify API to enable applications developers to
access this API in a simple and effective manner in Go. We implemented a
significant subset of the API’s functionality, choosing to focus on the
music-related functionality as deeply as possible. Future work for this
project would be to implement the rest of the API’s endpoints in methods
in our library, which would entail delving into the podcast-related
functionality that the API offers in addition to the music-focused
endpoints. Barring unforeseen errors that our tests have not caught, we
are confident that our library works and is sufficiently documented for
other developers to start building with Spotigo.
