package service

import (
	"time"

	coretypes "github.com/vitelabs/vite-portal/internal/core/types"
	nodetypes "github.com/vitelabs/vite-portal/internal/node/types"
	roottypes "github.com/vitelabs/vite-portal/internal/types"
)

// HandleDispatch handles the session information for a client request
func (s *Service) HandleDispatch(header coretypes.SessionHeader) (*coretypes.DispatchResponse, roottypes.Error) {
	header.Timestamp = time.Now().UnixMilli()
	err := header.ValidateHeader()
	if err != nil {
		return nil, err
	}

	session, err := s.getSession(header)
	if err != nil {
		return nil, err
	}

	// get actual nodes, compare count and update session if needed
	if a := s.getActualNodes(session); len(session.Nodes) != len(a) {
		session.Nodes = a
		s.Cache.SetSession(session)
	}

	currentNodeCount := s.NodeService.GetNodeCount(header.Chain)
	minNodeCount := roottypes.GlobalConfig.ConsensusNodeCount
	// make sure session has sufficient nodes
	if currentNodeCount > minNodeCount && minNodeCount > len(session.Nodes) || len(session.Nodes) == 0 {
		// delete current session and create new
		s.Cache.DeleteSession(header)
		session, err = s.getSession(header)
		if err != nil {
			return nil, err
		}
	}

	return &coretypes.DispatchResponse{
		Session: session,
	}, nil
}

func (s *Service) getSession(header coretypes.SessionHeader) (coretypes.Session, roottypes.Error) {
	// check cache
	session, found := s.Cache.GetSession(header, roottypes.GlobalConfig.MaxSessionDuration)
	if !found {
		// create new session
		session, err := coretypes.NewSession(s.NodeService, header, roottypes.GlobalConfig.SessionNodeCount)
		if err != nil {
			return coretypes.Session{}, err
		}
		// add to cache
		s.Cache.SetSession(session)
	}
	return session, nil
}

func (s *Service) getActualNodes(session coretypes.Session) []nodetypes.Node {
	var actualNodes []nodetypes.Node
	for _, v := range session.Nodes {
		n, found := s.NodeService.GetNode(v.Id)
		if !found || n.Chain != session.Header.Chain {
			continue
		}
		actualNodes = append(actualNodes, n)
	}
	if len(actualNodes) == 0 {
		return []nodetypes.Node{}
	}
	return actualNodes
}