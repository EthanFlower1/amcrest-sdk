package amcrest

import (
	"context"
	"fmt"
	"net/url"
)

// FaceService handles face detection and recognition related API calls.
// PDF Reference: docs/HTTP_API_V3.26.pdf pp. 334-361 (Section 9.2)
type FaceService struct {
	client *Client
}

// CreateGroup creates a new face recognition group with the given name and detail.
// Returns the groupID assigned by the camera.
// CGI: faceRecognitionServer.cgi?action=createGroup&groupName=X&groupDetail=Y
func (s *FaceService) CreateGroup(ctx context.Context, name, detail string) (string, error) {
	params := url.Values{
		"groupName":   {name},
		"groupDetail": {detail},
	}
	body, err := s.client.cgiGet(ctx, "faceRecognitionServer.cgi", "createGroup", params)
	if err != nil {
		return "", fmt.Errorf("FaceService.CreateGroup: %w", err)
	}
	kv := parseKV(body)
	id, ok := kv["groupID"]
	if !ok {
		return "", fmt.Errorf("FaceService.CreateGroup: groupID not found in response: %s", body)
	}
	return id, nil
}

// DeleteGroup deletes the face recognition group with the given groupID.
// CGI: faceRecognitionServer.cgi?action=deleteGroup&groupID=X
func (s *FaceService) DeleteGroup(ctx context.Context, groupID string) error {
	params := url.Values{
		"groupID": {groupID},
	}
	return s.client.cgiAction(ctx, "faceRecognitionServer.cgi", "deleteGroup", params)
}

// FindGroup returns the raw response listing all face recognition groups.
// CGI: faceRecognitionServer.cgi?action=findGroup
func (s *FaceService) FindGroup(ctx context.Context) (string, error) {
	body, err := s.client.cgiGet(ctx, "faceRecognitionServer.cgi", "findGroup", nil)
	if err != nil {
		return "", fmt.Errorf("FaceService.FindGroup: %w", err)
	}
	return body, nil
}

// GetGroupForChannel returns the raw response for the face group assigned to
// the given video channel.
// CGI: faceRecognitionServer.cgi?action=getGroup&channel=N
func (s *FaceService) GetGroupForChannel(ctx context.Context, channel int) (string, error) {
	params := url.Values{
		"channel": {fmt.Sprintf("%d", channel)},
	}
	body, err := s.client.cgiGet(ctx, "faceRecognitionServer.cgi", "getGroup", params)
	if err != nil {
		return "", fmt.Errorf("FaceService.GetGroupForChannel: %w", err)
	}
	return body, nil
}
