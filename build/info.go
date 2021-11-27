package build

import "gitlab.com/golibs-starter/golib/actuator"

type Version string
type CommitHash string
type Time string

type InfoResponse struct {
	Version    Version    `json:"version,omitempty"`
	CommitHash CommitHash `json:"commit_hash,omitempty"`
	Time       Time       `json:"time,omitempty"`
}

type Informer struct {
	version    Version
	commitHash CommitHash
	buildTime  Time
}

func NewInformer(version Version, commitHash CommitHash, buildTime Time) actuator.Informer {
	return &Informer{
		version:    version,
		commitHash: commitHash,
		buildTime:  buildTime,
	}
}

func (b Informer) Key() string {
	return "build"
}

func (b Informer) Value() interface{} {
	return InfoResponse{
		Version:    b.version,
		CommitHash: b.commitHash,
		Time:       b.buildTime,
	}
}
