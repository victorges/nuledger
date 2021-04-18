# nuledger

This is an implementation of a simple "new" ledger, originally for a remote
interview. It's original features included the tracking of state of a single
account, with an available limit of currency to be spent and whether its card
was active or not. Also multiple rules to authorize transactions in the account.
The active state of the card cannot be changed, but the available limit should
reduce for every transaction that is authorized.

Implemented the [specification](#Spec) with some design decisions to make it easy to
expand in the future. Be it with new business logic for authorizing transactions
as well as incrementing the wire format in which the operations are received or
the actual ledger implementation.

After submitting the project for evaluation, also implemented support for
multiple accounts, which I felt would be a nice addition. The next nice feature
to add after that could be to support changing the active state of an account's
card, for example to (un)block the card of an existing account.

There are both extensive documentation and 100% coverage tests of the code,
which can be both inspected easily in the browser. The integrated server for
documentation (`godoc`) can be started with `make doc` while the tests server
(`goconvey`) with `make test_server`.

## Project Operation

There is a `Makefile` with hopefully all the commands that will be necessary for
inspecting, building, testing and running the project.

### Requirements

 `go version 1.13+` and/or `docker`

### Inspect

Instead of reading the documentation directly in the source files, you can read
it in the `godoc` interface in your browswer.

To start the `godoc` server:
```
make doc
```

Then read the project documentation in your browser at:
http://localhost:6060/pkg/nuledger

### Test

To generate the mock implementations:
```
make generate
```

Notice it is not always necessary to re-generate the mocks, only when the mocked
interfaces have changed and thus need some update. Otherwise, they're already
there and checked-in to version control.

To execute the tests:
```
make test
```

Alternatively, start the `goconvey` server for visualizing the test results:
```
make test_server
```

Then watch the test results display in your browser. If the browser doesn't open
automatically, click:
http://localhost:8080

### Build
To build the project locally:
```
make build
```

It outpts the executable file to the `./build` folder with `authorizer` name.
e.g. You can test it with one of the integraction test cases with:
```
make build && ./build/authorizer < testcases/base/in.jsonl
```

Alternatively, to build the project with Docker (also runs tests):
```
make docker
```

### Run
To run the application locally:
```
make run
```
That one doesn't generate build output as well so it leaves the folder clean.

To run the application in Docker:
```
make docker_run
```

That will run the docker build command again and then run the container. You can
pipe the input directly to make and it will work as expected, e.g.:
```
make docker_run < testcases/highFreqTransactions/in.jsonl
```

## Spec

The application should read all the input from `stdin` and write all the output
to `stdout`, both in the [JSON Lines](https://jsonlines.org) format. Each JSON
in the input represents an operation to be performed, either creating an account
(`account` field in input object) or performing – after authorizing – a
transaction (`transaction` field).

Each output should include the final state of the corresponding account, in the
`account` field of the output object. It could be either the transformed state
after the given operation was performed or the exactly the previous state in
case there was some business rule violation for performing the transaction.

The business rule violations are well-specified, being:
 - `account-already-initialized`: Account had already been initialized when
   another create account operation was requested (with the same account ID).
 - `account-not-initialized`: A perform transaction operation was attempted
   before the corresponding account was actually initialized.
 - `card-not-active`: A perform transaction operation was attempted in an
   account whose card is not active.
 - `insufficient-limit`: A transaction is attempted with an amount higher than
   the available limit of the account.
 - `high-frequency-small-interval`: Too many transactions are performed in the
   same account within a small interval. This limit is currently configured as a
   maximum of 3 transactions in the same account every 2 minutes.
 - `double-transaction`: A duplicate transaction was attempted. This means that
   the attempted transaction has the same account, amount and merchant of a
   recent transaction. A transaction is currently considered to be recent if
   performed at most 2 minutes ago.

Any violations of business logic must be included in the `violations` field of
the output object written to `stdout`. This means that if the operation was not
authorized then at least one violation will be present in the output.

Some examples of input/output combinations can be seen in the `testcases` folder
in the root of the project (run by the `main_test.go` test). A very simple one
(`base`) is:

 - Input:
```
{"account": {"active-card": true, "available-limit": 100}}
{"transaction": {"merchant": "Pizza Zagga", "amount": 20, "time": "2019-02-13T10:00:00.000Z"}}
{"transaction": {"merchant": "TT Burger", "amount": 90, "time": "2019-02-13T11:00:00.000Z"}}
```
 - Output:
```
{"account":{"active-card":true,"available-limit":100},"violations":[]}
{"account":{"active-card":true,"available-limit":80},"violations":[]}
{"account":{"active-card":true,"available-limit":80},"violations":["insufficient-limit"]}
```

Notice that every transaction has a timestamp, and it is a hard requirement by
the program that the timestamps must be received in order. Otherwise, we
wouldn't be able to process transactions one by one since many of the algorithms
depend on the chronological order. If the program didn't have that guarantee,
it'd also need to keep a buffer of some of the last seen transactions to be sure
to process each transaction only when no transaction before it could show up.

Also notice that this example does not have any reference to an account ID.
Since the multi-account was implemented as an additional feature, it is also
completely optional. So an account can specify no ID which has the same behavior
as an empty string ID. To use the multi-account feature, an `id` field has to be
specified in the create account operation and an `accountId` field has to be
specified in the perform transaction operation, and it will correspondingly
appear in the output objects.

## Design

Some design decisions were made, so some of the higher level ones will be
detailed here for easier understanding of the whole project. As mentioned in the
previous session, lower-level documentation is also available for all the
exported components in the code and can be inspected either directly in the code
or in the browser via `godoc`.

### I/O Processor

For the input/output processing part of the application I created a separate
`iop` package, which stands for exactly that.

The application was proposed with a very specific form of input/output, reading
line-separated JSONs from the standard input, processing them and writing them
to the standard output. With that in mind, I believed it made a lot of sense to
separate that specific logic, both this read-process-write pipeline and the
actual JSON format of the received and returned objects, into that separate
package.

The package also exports a single interface which is where any business logic
for processing the operations needs to plugin. The interface basically receives
the input JSON and returns the response JSON, so a slightly higher level
component can already be created without worrying about the I/O pipeline.

The `iop` is similar to an HTTP framework that handles the lower level protocol
and provides the higher level objects to a component to do any business logic.

### Rule Authorizers

The other core piece of the code architecture are the rule authorizers. Their
interface and some generic helpers are defined in the `authorizer/rule` package,
while the specific rule implementations are in `authorizer/rules`.

They exist so that the actual account-managing part of the application are as
extensible as possible, with authorization rules being easily created or removed
from the default set of rules. In the specific implementation, each violation
code is validated by a specific authorizer, but we can also combine multiple
of them in a single authorizer if helpful. The violations are returned as errors
in the authorization, later translated into an actual violations array in the
response.

These can also allow for flexible managing of accounts, and we could choose
different sets of authorizers depending on other specific rules. For example, an
account could have some overdraft feature to alow it to go below its limit, so
we could include the specific authorizer about that or not.

For the specific violation about maximum frequency of transactions, there is a
rate limiter utility in the `util` package which has the core frequency limiting
logic. It implements an "optimal response" algorithm for rate-limiting, in the
sense that it might become expensive for a lot of events but it guarantees that
the correct response will be given regarding the number of transactions in the
past observation interval. Decided to use it for the double transaction as well
to avoid re-implementing some custom logic, even though the double transactions
would be rather simpler to implement directly with just a timestamp.

### Authorizer

The `authorizer` is the package with the "most core" business logic of the
application apart from the rule authorizers mentioned above.

The first relevant component is an implementation of an I/O handler which
receives the JSON objects from the `iop` package, processes them internally and
translates the response back to the I/O pipeline. Following the same analogy of
HTTP frameworks, it would be an API router/controller which routes specific
requests to the corresponding API that should be called and then translates the
response to the protocol being used.

The final one is the ledger itself, which ends up being pretty simple given the
above abstractions. It provides explicit methods for each of the operations
supported by the system (creating an account and performing a transaction) and
is the component called by the handler from the paragraph above. Its
implementation is rather simple though, since it only needs to validate the
(non-)existence of the account, authorizer the operation with the configured
rule authorizer, and then actually perform the operation if all is fine.

### Tests

There are both integration and unit tests in the project.

The integration test is written in the root of the project, in the
`main_test.go` file. It goes through all the test cases in the `testcases`
folder, each represented by a sub-folder with an `in.jsonl` and `out.jsonl`
files for input and expected output respectively.

The unit tests are written across the project in the `*_test.go` files, at least
one in each (non-generated) package. These unit tests also make use of test
mocks generated by the `golang/mock` library. The mock generation makes use of
the `go generate` command that processes a couple of `//go:generate` directives
in some files, which leverage from the `gen_mocks.sh` script in the project root
as well to save some boilerplate.

All the tests are also written with the `goconvey` library, so one can run their
CLI/webserver to see a friendlier UI with all the executed tests, test coverage
and which automatically refreshes with changes in the source code as well.
