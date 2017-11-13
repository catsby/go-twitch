# go-twitch

`go-twitch` is a SDK for the [Twitch.tv](https://twitch.tv) platform API written
in golang, currently under development. It supports v5 of the Twitch API, with
plans to support the new Twitch API as it is developed.

The current API, V5, is deprecated and will be removed on December 31st, 2018.
The new Twitch API is live but not 100% feature complete, and is under active
development. 

- V5 docs: https://dev.twitch.tv/docs/v5
- New API docs: https://dev.twitch.tv/docs/api ("new" API, not complete yet)

# Usage 

The Current API service, V5, is nicknamed `kraken` (or so I assume, based on the API base url
or `https://api.twitch.tv/kraken/`). The new API is nicknamed `helix` (same
guess off of base url, `https://api.twitch.tv/helix/`). To support both, this
SDK offers a client for each, `kraken` and `helix`, with functionality tied to
their corresponding endpoint/service. `Kraken` is more compmlete in both terms
of API functionality and SDK coverage, while `Helix` is a work in progress on
both parts. I plan to continue to implement things missing in `Kraken` while
also building out `Helix` as it's developed. This SDK will support both until
`Kraken` is turned off.

To include `go-twitch` in your project, first get it with `go get`:

    $ go get -u github.com/catsby/go-twitch

You can then import and use it in your code:

    # main.go
    package main
    
    import (
    	"fmt"
    	"log"
    
    	"github.com/catsby/go-twitch/service/kraken"
    )
    
    func main() {
    	client := kraken.DefaultClient(nil)
    
    	me, err := client.GetUser(nil)
    	if err != nil {
    		log.Fatalf("Error finding me: %s", err)
    	}
    
    	fmt.Println("My name is", me.Name)
    }

See `examples/streaming/main.go` in this repository for an example.

# Development

*Note:* This is considered alpha software. It should work as described without
crash or error, however method/function names and signatures may change during
the alpha period.  I'm still trying to finalize some design decisions, for
example: 1:1 mapping of API endpoints by name, or intent, and returning raw
response or not. Currently I'm using https://github.com/dnaeon/go-vcr to record
interactions and test against for development. The test fixtures are stored in
the `fixtures` directory. 

To get started with development, checkout this repository and change
directories.

    $ mkdir -p $GOPATH/src/github.com/catsby
    $ cd $GOPATH/src/github.com/catsby
    $ git clone git@github.com:catsby/go-twitch.git
    $ cd go-twitch

From there, run `make bootstrap` to install the required libraries.

    $ make bootstrap
    ==> Bootstrapping github.com/catsby/go-twitch...
    --> Installing github.com/ajg/form
    --> Installing github.com/dnaeon/go-vcr/cassette
    --> Installing github.com/dnaeon/go-vcr/recorder
    --> Installing github.com/hashicorp/go-cleanhttp
    --> Installing github.com/mitchellh/mapstructure
    --> Installing gopkg.in/yaml.v2
    --> Installing gopkg.in/check.v1


The tests expect an environment variable `TWITCH_ACCESS_TOKEN` to be set. If you
want to work within the fixtures provided, you should be able to set this to
anything:

    $ export TWITCH_ACCESS_TOKEN="anything"

In order to modify or add new features, you'll need to create your own OAuth
Access token. See the [Authentication Guide on
dev.twitch.tv](https://dev.twitch.tv/docs/v5/guides/authentication) for details.
After you have your access token, you can put it into your env, similar to
above.

Verify the tests pass before you start working:

    $ make test

You can also run specific tests by matching a naming pattern:

    $ make test TESTARGS="-v -run=TestUser_Get_self"


# Pre-commit hook, fuzzing

Included in this repository is a `scripts` directory, containing 3 files:

    - scripts/
      - fuzz.sh
      - pre-commit
      - unfuzz.sh

The pre-commit hook will check for any files that contain the contents of
`$TWITCH_ACCESS_TOKEN` environment variable. The purpose of this pre-commit hook
is to search through the files and detect if your personal access token is in
the source, to prevent you from committing it. It depends on the environment
variable being set, and `pre-commit` hook being in the correct place. You can
install it with a `make` command:

    $ make pre-commit

## Fuzzing

If you've added an endpoint and as a result, `go-vcr` has added your API key to
a fixture file, you can fuzz(?) it with `scripts/fuzz.sh`:

    $ ./scripts/fuzz.sh

This will replace any occurrence of `$TWITCH_ACCESS_TOKEN` with `xxxxxxxxxxxxx`.
The `unfuzz.sh` script will reverse the process if you want. The auth token is
not required when tests that use the fixtures as they are.

# Support

This library is provided as is, with no guarantees. I will do my best to ensure
quality but always assume your own risk. See the included License as well.

# Contributions

Contributions are welcome! Please open a GitHub issue or Pull Request and I'll
review them as soon as I'm able to. Thanks!

The "core" of this library is mostly complete and functional, being borrowed
from github.com/sethvargo/go-fastly. Majority of the remaining work needed is
getting coverage of Twitch's v5 API:

- https://dev.twitch.tv/docs/v5
