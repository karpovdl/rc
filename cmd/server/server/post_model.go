package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	mgo "gopkg.in/mgo.v2"
	bson "gopkg.in/mgo.v2/bson"
)

//PostModel type
type PostModel struct {
	DB *mgo.Database
}

// Create new post
func (postModel *PostModel) create(post *Post, author *Author) (int, error) {
	if postModel == nil || postModel.DB == nil {
		return http.StatusInternalServerError, errors.New("post model type not initialized")
	}

	if post == nil {
		return http.StatusInternalServerError, errors.New("post not initialized")
	}

	if author == nil {
		return http.StatusInternalServerError, errors.New("author not initialized")
	}

	post.ID = bson.NewObjectId()
	post.Score = 1
	post.Views = 0
	post.UpvotePercentage = 100

	now := time.Now()
	post.Created = fmt.Sprintf("%sT%sZ", now.Format("2006-01-02"), now.Format("03:04:05.000"))

	post.Votes = []Vote{{
		User: author.ID,
		Vote: 1},
	}

	post.Author = Author{
		ID:       author.ID,
		Username: author.Username,
	}

	post.Comments = []Comment{}

	if err := postModel.DB.C("posts").Insert(post); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// Check post
func (postModel *PostModel) checkExistsByID(post *Post) (int, error) {
	if postModel == nil || postModel.DB == nil {
		return http.StatusInternalServerError, errors.New("post model type not initialized")
	}

	if post == nil {
		return http.StatusInternalServerError, errors.New("post not initialized")
	}

	if count, err := postModel.DB.C("posts").Find(
		bson.M{
			"_id": bson.ObjectIdHex(string(post.ID)),
		}).Count(); count == 0 || err != nil {
		if count == 0 {
			return http.StatusNotFound, errors.New("post not found")
		} else if count > 1 {
			return http.StatusInternalServerError, errors.New("post many found")
		}
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// Find post by id
func (postModel *PostModel) findByID(post *Post) (int, error) {
	if postModel == nil || postModel.DB == nil {
		return http.StatusInternalServerError, errors.New("post model type not initialized")
	}

	if post == nil {
		return http.StatusInternalServerError, errors.New("post not initialized")
	}

	if err := postModel.DB.C("posts").Find(
		bson.M{
			"_id": bson.ObjectIdHex(string(post.ID)),
		}).One(&post); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// Update post by id
func (postModel *PostModel) updateByID(post *Post) (int, error) {
	if postModel == nil || postModel.DB == nil {
		return http.StatusInternalServerError, errors.New("post model type not initialized")
	}

	if post == nil {
		return http.StatusInternalServerError, errors.New("post not initialized")
	}

	if err := postModel.DB.C("posts").Update(
		bson.M{
			"_id": post.ID,
		},
		post); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// Delete exists post by id
func (postModel *PostModel) deleteByID(post *Post, author *Author) (int, error) {
	if postModel == nil || postModel.DB == nil {
		return http.StatusInternalServerError, errors.New("post model type not initialized")
	}

	if post == nil {
		return http.StatusInternalServerError, errors.New("post not initialized")
	}

	if status, err := postModel.checkExistsByID(post); err != nil {
		return status, err
	}

	if err := postModel.DB.C("posts").Remove(
		bson.M{
			"_id":        bson.ObjectIdHex(string(post.ID)),
			"author._id": bson.ObjectId(author.ID),
		}); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// Get all posts
func (postModel *PostModel) getAll() ([]Post, int, error) {
	if postModel == nil || postModel.DB == nil {
		return nil, http.StatusInternalServerError, errors.New("post model type not initialized")
	}

	var p []Post
	if err := postModel.DB.C("posts").Find(nil).All(&p); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return p, http.StatusOK, nil
}

// Get all posts only auth user
func (postModel *PostModel) getAllByUser(author *Author) ([]Post, int, error) {
	if postModel == nil || postModel.DB == nil {
		return nil, http.StatusInternalServerError, errors.New("post model type not initialized")
	}

	if author == nil {
		return nil, http.StatusInternalServerError, errors.New("author not initialized")
	}

	var p []Post
	if err := postModel.DB.C("posts").Find(
		bson.M{
			"author.username": author.Username,
		}).All(&p); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return p, http.StatusOK, nil
}

// Get all posts only category
func (postModel *PostModel) getAllByCategory(category string) ([]Post, int, error) {
	if postModel == nil || postModel.DB == nil {
		return nil, http.StatusInternalServerError, errors.New("post model type not initialized")
	}

	if category == "" {
		return nil, http.StatusInternalServerError, errors.New("category not initialized")
	}

	var p []Post
	if err := postModel.DB.C("posts").Find(
		bson.M{
			"category": category,
		}).All(&p); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return p, http.StatusOK, nil
}

// Get post from id
func (postModel *PostModel) getByID(post *Post) (int, error) {
	if postModel == nil || postModel.DB == nil {
		return http.StatusInternalServerError, errors.New("post model type not initialized")
	}

	if post == nil {
		return http.StatusInternalServerError, errors.New("post not initialized")
	}

	if status, err := postModel.checkExistsByID(post); err != nil {
		return status, err
	}

	if status, err := postModel.findByID(post); err != nil {
		return status, err
	}

	post.Views++

	return postModel.updateByID(post)
}

// Upvote post by id
func (postModel *PostModel) upvoteByID(post *Post, author *Author) (int, error) {
	if postModel == nil || postModel.DB == nil {
		return http.StatusInternalServerError, errors.New("post model type not initialized")
	}

	if post == nil {
		return http.StatusInternalServerError, errors.New("post not initialized")
	}

	if status, err := postModel.checkExistsByID(post); err != nil {
		return status, err
	}

	if status, err := postModel.findByID(post); err != nil {
		return status, err
	}

	setVote(post, author, true)

	setScoreAndUpvotePercentage(post)

	return postModel.updateByID(post)
}

// Downvote post by id
func (postModel *PostModel) downvoteByID(post *Post, author *Author) (int, error) {
	if postModel == nil || postModel.DB == nil {
		return http.StatusInternalServerError, errors.New("post model type not initialized")
	}

	if post == nil {
		return http.StatusInternalServerError, errors.New("post not initialized")
	}

	if status, err := postModel.checkExistsByID(post); err != nil {
		return status, err
	}

	if status, err := postModel.findByID(post); err != nil {
		return status, err
	}

	setVote(post, author, false)

	setScoreAndUpvotePercentage(post)

	return postModel.updateByID(post)
}

// Add new comment to post by id
func (postModel *PostModel) addCommentByID(post *Post, author *Author, comment *Comment) (int, error) {
	if postModel == nil || postModel.DB == nil {
		return http.StatusInternalServerError, errors.New("post model type not initialized")
	}

	if post == nil {
		return http.StatusInternalServerError, errors.New("post not initialized")
	}

	if status, err := postModel.checkExistsByID(post); err != nil {
		return status, err
	}

	if status, err := postModel.findByID(post); err != nil {
		return status, err
	}

	now := time.Now()
	post.Comments = append(post.Comments, Comment{
		ID:      bson.NewObjectId(),
		Created: fmt.Sprintf("%sT%sZ", now.Format("2006-01-02"), now.Format("03:04:05.000")),
		Author:  *author,
		Body:    comment.Comment,
	})

	return postModel.updateByID(post)
}

// Delete exists comment to post by id
func (postModel *PostModel) deleteCommentByID(post *Post, author *Author, comment *Comment) (int, error) {
	if postModel == nil || postModel.DB == nil {
		return http.StatusInternalServerError, errors.New("post model type not initialized")
	}

	if post == nil {
		return http.StatusInternalServerError, errors.New("post not initialized")
	}

	if status, err := postModel.checkExistsByID(post); err != nil {
		return status, err
	}

	if status, err := postModel.findByID(post); err != nil {
		return status, err
	}

	for indx, c := range post.Comments {
		if c.Author.ID == author.ID && c.ID == bson.ObjectIdHex(string(comment.ID)) {
			post.Comments = append(post.Comments[:indx], post.Comments[indx+1:]...)
		}
	}

	return postModel.updateByID(post)
}

func setVote(post *Post, author *Author, isUpvote bool) (int, error) {
	if post == nil {
		return http.StatusInternalServerError, errors.New("post not initialized")
	}

	upAdd := 0
	if isUpvote {
		upAdd = 1
	} else {
		upAdd = -1
	}

	if len(post.Votes) == 0 {
		post.Votes = append(post.Votes, Vote{User: author.ID, Vote: upAdd})
	} else {
		isExistsAuthor := false
		for indx, vote := range post.Votes {
			if vote.User == author.ID {
				isExistsAuthor = true
				if vote.Vote == (1 * upAdd) {
					post.Votes = append(post.Votes[:indx], post.Votes[indx+1:]...)
				} else if vote.Vote == (-1 * upAdd) {
					post.Votes[indx].Vote = upAdd
				} else {
					post.Votes = append(post.Votes, Vote{User: author.ID, Vote: upAdd})
				}
				break
			}
		}
		if !isExistsAuthor {
			post.Votes = append(post.Votes, Vote{User: author.ID, Vote: upAdd})
		}
	}

	return http.StatusOK, nil
}

func setScoreAndUpvotePercentage(post *Post) (int, error) {
	if post == nil {
		return http.StatusInternalServerError, errors.New("post not initialized")
	}

	countUp := 0
	countDown := 0
	countUpDown := 0
	for _, vote := range post.Votes {
		if vote.Vote == 1 {
			countUp++
			countUpDown++
		}
		if vote.Vote == -1 {
			countDown--
			countUpDown++
		}
	}

	post.Score = countUp + countDown

	if countUpDown > 0 {
		post.UpvotePercentage = (100 * countUp) / countUpDown
	} else {
		post.UpvotePercentage = 0
	}

	return http.StatusOK, nil
}
