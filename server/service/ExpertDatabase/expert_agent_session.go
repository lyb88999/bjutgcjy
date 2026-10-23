package ExpertDatabase

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"

	ExpertDatabaseRes "github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase/response"
)

const (
	agentSessionTTL         = 30 * time.Minute
	agentMaxSessions        = 500
	agentMaxSessionsPerUser = 20
)

type expertAgentSession struct {
	owner        uint
	updated      time.Time
	busy         bool
	plan         ExpertDatabaseRes.ExpertAgentPlan
	candidateIDs []uint
	selectedIDs  []uint
	hints        []expertAgentCandidateHint
	browseAll    bool
}

type expertAgentCandidateHint struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	UnitName string `json:"unitName"`
}

type expertAgentSessionStore struct {
	mu      sync.Mutex
	entries map[string]*expertAgentSession
}

func newExpertAgentSessionStore() *expertAgentSessionStore {
	return &expertAgentSessionStore{entries: make(map[string]*expertAgentSession)}
}

var defaultExpertAgentSessions = newExpertAgentSessionStore()

// acquire 不等待另一个同会话请求，也不创建后台goroutine；整轮调用持有busy标记。
func (s *expertAgentSessionStore) acquire(operatorID uint, id string) (string, *expertAgentSession, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	owned := 0
	for key, session := range s.entries {
		if !session.busy && now.Sub(session.updated) >= agentSessionTTL {
			delete(s.entries, key)
			continue
		}
		if session.owner == operatorID {
			owned++
		}
	}
	if id != "" {
		session, ok := s.entries[id]
		if !ok || session.owner != operatorID {
			return "", nil, false, ErrExpertAgentSessionNotFound
		}
		if session.busy {
			return "", nil, false, ErrExpertAgentSessionBusy
		}
		session.busy = true
		return id, session, false, nil
	}
	if len(s.entries) >= agentMaxSessions || owned >= agentMaxSessionsPerUser {
		return "", nil, false, ErrExpertAgentSessionLimit
	}
	var token [16]byte
	if _, err := rand.Read(token[:]); err != nil {
		return "", nil, false, err
	}
	id = hex.EncodeToString(token[:])
	session := &expertAgentSession{owner: operatorID, updated: now, busy: true, plan: ExpertDatabaseRes.ExpertAgentPlan{
		Keywords: []string{}, ExcludedExpertIDs: []uint{}, Count: 5, MaxPerUnit: 2, Diversify: true,
	}}
	s.entries[id] = session
	return id, session, true, nil
}

func (s *expertAgentSessionStore) release(id string, discard bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if discard {
		delete(s.entries, id)
		return
	}
	if session := s.entries[id]; session != nil {
		session.busy = false
	}
}

func cloneExpertAgentPlan(p ExpertDatabaseRes.ExpertAgentPlan) ExpertDatabaseRes.ExpertAgentPlan {
	p.Keywords = append([]string{}, p.Keywords...)
	p.ExcludedExpertIDs = append([]uint{}, p.ExcludedExpertIDs...)
	return p
}
