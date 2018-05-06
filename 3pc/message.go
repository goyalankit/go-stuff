package main

// VoteReqMsg is the message sent by coordinator
// to ask for votes. This is vote-req message.
type VoteReqMsg struct{}

// VoteResponseType  is the response to
// vote-req message.
type VoteResponseType bool

const (
	// Yes means cohort votes yes.
	Yes VoteResponseType = true

	// No means cohort votes no.
	No VoteResponseType = false
)

type VoteReqRespMsg struct {
	Vote VoteResponseType
}

type PrecommitMsg struct{}

type AbortMsg struct{}

type AckMsg struct{}
