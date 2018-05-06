package main

// VoteReq is the message sent by coordinator
// to ask for votes.
type VoteReqMsg struct{}

type VoteResponseType bool

const (
	Yes VoteResponseType = true
	No  VoteResponseType = false
)

type VoteReqRespMsg struct {
	Vote VoteResponseType
}

type PrecommitMsg struct{}

type AbortMsg struct{}

type AckMsg struct{}
