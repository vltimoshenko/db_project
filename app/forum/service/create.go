package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"

	"github.com/db_project/pkg/messages"
	. "github.com/db_project/pkg/models"
)

func (s Service) CreatePosts(body io.ReadCloser, slugOrId string) ([]Post, error) {
	bytes, err := ioutil.ReadAll(body)
	if err != nil {
		//return uuid.UUID{}, errors.New(BadRequestMsg)
	}

	var posts []NewPost
	err = json.Unmarshal(bytes, &posts)
	if err != nil {
		//return uuid.UUID{}, errors.New(InvalidJSONMsg)
	}

	threadID, err := strconv.Atoi(slugOrId)

	var thread Thread
	if err != nil {
		thread, err = s.Repository.GetThreadBySlug(slugOrId)
	} else {
		thread, err = s.Repository.GetThreadByID(threadID)
	}
	if err != nil {
		return []Post{}, errors.New(messages.ThreadDoesNotExist)
	}

	returnPosts, err := s.Repository.CreatePosts(posts, threadID, thread.Forum)

	return returnPosts, err
}

func (s Service) Vote(vote Vote, slugOrId string) (Thread, error) {
	threadID, err := strconv.Atoi(slugOrId)

	//should remove
	if vote.Voice != 1 && vote.Voice != -1 {
		return Thread{}, fmt.Errorf("Invalid value")
	}

	var thread Thread
	if err != nil {
		thread, err = s.Repository.GetThreadBySlug(slugOrId)
	} else {
		thread, err = s.Repository.GetThreadByID(threadID)
	}
	if err != nil {
		return thread, errors.New(messages.ThreadDoesNotExist)
	}

	checkVote, err := s.Repository.GetVote(vote.Nickname, thread.ID)

	var dif int
	if err != nil {
		err = s.Repository.CreateVote(vote, thread.ID)
		if err != nil {
			return thread, err
		}

		if vote.Voice == 1 {
			dif = 1
		} else {
			dif = -1
		}
	} else {
		err = s.Repository.ChangeVote(vote, thread.ID)
		if err != nil {
			fmt.Println(err)
			return thread, err
		}

		if checkVote.Voice == vote.Voice {

		} else if checkVote.Voice == 1 {
			dif = -2
		} else {
			dif = 2
		}
	}

	thread.Votes += dif
	err = s.Repository.ChangeThreadRate(dif, thread.ID)
	return thread, err
}
