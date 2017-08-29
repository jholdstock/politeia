// Copyright (c) 2017 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package backend

import (
	"crypto/sha256"
	"errors"
)

var (
	// ErrProposalNotFound is emitted when a proposal could not be found
	ErrProposalNotFound = errors.New("proposal not found")

	// ErrShutdown is emitted when the backend is shutting down.
	ErrShutdown = errors.New("backend is shutting down")

	// ErrInvalidTransition is emitted when an invalid status transition
	// occurs.  The only valid transitions are from unvetted -> vetted and
	// unvetted to censored.
	ErrInvalidTransition = errors.New("invalid proposal status transition")
)

// ContentVerificationError is returned when a submitted proposal contains
// unacceptable file formats or corrupt data.
type ContentVerificationError struct {
	Err error
}

func (c ContentVerificationError) Error() string {
	return c.Err.Error()
}

type File struct {
	Name    string // Basename of the file
	MIME    string // MIME type
	Digest  string // SHA256 of decoded Payload
	Payload string // base64 encoded file
}

type PSRStatusT int

const (
	// All possible PSR status codes
	PSRStatusInvalid  PSRStatusT = 0
	PSRStatusUnvetted PSRStatusT = 1
	PSRStatusVetted   PSRStatusT = 2
	PSRStatusCensored PSRStatusT = 3
)

var (
	// PSRStatus converts a status code to a human readable error.
	PSRStatus = map[PSRStatusT]string{
		PSRStatusInvalid:  "invalid",
		PSRStatusUnvetted: "unvetted",
		PSRStatusVetted:   "vetted",
		PSRStatusCensored: "censored",
	}
)

type ProposalStorageRecord struct {
	Version uint              // Iteration count of proposal
	Status  PSRStatusT        // Current status of the proposal
	Merkle  [sha256.Size]byte // Merkle root of all files in proposal
	Name    string            // Short name of proposal
	Token   []byte            // Proposal authentication token
}

type Backend interface {
	// Create new proposal
	New(string, []File) (*ProposalStorageRecord, error)

	// Get unvetted proposal
	GetUnvetted([]byte) ([]File, *ProposalStorageRecord, error)

	// Get vetted proposal
	GetVetted([]byte) ([]File, *ProposalStorageRecord, error)

	// Set unvetted proposal status
	SetUnvettedStatus([]byte, PSRStatusT) (PSRStatusT, error)

	// Close performs cleanup of the backend.
	Close()
}