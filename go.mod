module cryptopener

go 1.15


require (
    github.com/larikuistio/cryptopener v0.0.0
    github.com/larikuistio/cryptopener/testServer v0.0.0
)

replace (
    github.com/larikuistio/cryptopener => ./cryptopener
    github.com/larikuistio/cryptopener/testServer => ./testServer
)