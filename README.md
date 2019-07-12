# Photo Registry Contract

An Orbs smart contract managing a photo registry over the blockchain.

Supports the following public API:

### Registration
Registers a photo with metadata

`register(phash string, meta string)`

`phash` is [pHash](https://www.phash.org) string of the registering photo

`meta` is the registering photo JSON metadata

### Verification
Verifies if photo(s) were registered with a given pHash 

`verify(phash string) string`

`phash` is [pHash](https://www.phash.org) string of the registering photo

Returns the match photo's metadata

### Search
Searches photos similar to photo(s) with a given pHash (including phto(s) themselves)

`search(phash string, minScore uint64) string`

`phash` is [pHash](https://www.phash.org) string of the registering photo

`minScore` is the minimum [Hamming Distance](https://en.wikipedia.org/wiki/Hamming_distance) from 0 to 100

Returns JSON array of matching photos' metadata with the scores
