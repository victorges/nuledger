// Package authorizer contains the higher-level components of the authorizer
// logic.
//
// It's main components are a Ledger which handles the core business
// logic of creating accounts and performing transactions and a Handler which
// translates between lower-level JSON messages and the Ledger interface.
package authorizer
