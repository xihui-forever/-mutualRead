package types

type Topic string

func (p Topic) String() string {
	return string(p)
}

const (
	EventAppealChangedTopic Topic = "appeal_changed"
)
