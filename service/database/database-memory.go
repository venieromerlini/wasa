package database

import (
	"errors"
	"sort"
	"strconv"
	"time"
	"wasa/service/model"
)

// AppDatabaseMemory is the high level interface for the DB
type AppDatabaseMemory interface {
	FindAllUsers() []model.User
	SaveUser(username string) model.User

	FindAllBans(username string) []model.Ban
	SaveBan(ban model.BanRequest) model.Ban
	DeleteBan(requestorUser string, id int64)

	FindAllFollow(username string) []model.Follow
	SaveFollow(follow model.FollowRequest) model.Follow
	DeleteFollow(requestorUser string, id int64)

	FindAllComments(photoId int64) []model.Comment
	SaveComment(commentRequest model.CommentRequest) model.Comment
	DeleteComment(requestorUser string, id int64)

	FindAllLikes(photoId int64) []model.Like
	SaveLike(likeRequest model.LikeRequest) model.Like
	DeleteLike(requestorUser string, inInt64 int64)

	SavePhoto(username string, bytes []byte) model.Photo
	DeletePhoto(requestorUser string, id int64)
	FindAllPhotos(username string) []model.Photo
	FindPhoto(photoId int64) (model.Photo, error)

	UpdateUsername(requestorUser string, oldUsername string, newUsername string) model.User
	FindUserHomePageByUsername(username string) model.UserHomePage
	FindUserProfileByUsername(username string) model.Profile
}

type appdbmemimpl struct {
	usersMap    map[int64]model.User
	userIdsMap  map[string]int64
	bansMap     map[int64]model.Ban
	followsMap  map[int64]model.Follow
	commentsMap map[int64]model.Comment
	likesMap    map[int64]model.Like
	photosMaps  map[int64]model.Photo
	sequence    map[int64]int64
}

func (m appdbmemimpl) FindUserProfileByUsername(username string) model.Profile {
	profile := new(model.Profile)
	user := m.findUserByUsername(username)
	profile.User = &user
	profile.Photos = m.FindAllPhotos(username)
	profile.Followees = m.FindAllFollow(username)
	profile.Followers = m.findAllFollowers(username)
	return *profile
}

func (m appdbmemimpl) UpdateUsername(_ string, oldUsername string, newUsername string) model.User {
	user := new(model.User)
	user.Id = m.userIdsMap[oldUsername]
	user.Username = newUsername
	m.usersMap[user.Id] = *user
	m.userIdsMap[newUsername] = user.Id
	delete(m.userIdsMap, oldUsername)
	return *user
}

/*
UserStream // Each user will be presented with a stream of photos (images) in reverse chronological order,
with information about when each photo was uploaded (date and time  and how many likes and comments it has.
The stream is composed by photos from “following” (other users that the user follows).
*/
func (m appdbmemimpl) subtract(from []string, what []string) []string {
	var final []string
	found := false
	for _, f := range from {
		elem := f
		for _, w := range what {
			if f == w {
				found = true
				break
			}
		}
		if !found {
			final = append(final, elem)
		}
		found = false

	}
	return final
}

func (m appdbmemimpl) FindUserHomePageByUsername(username string) model.UserHomePage {
	follows := m.FindAllFollow(username)
	bans := m.findAllBanners(username)

	var followUsernames []string
	for _, follow := range follows {
		followUsernames = append(followUsernames, follow.Followee.Username)
	}

	var bansUsernames []string
	for _, ban := range bans {
		bansUsernames = append(bansUsernames, ban.User.Username)
	}

	var photos []model.Photo
	for _, usernames := range m.subtract(followUsernames, bansUsernames) {
		photos = append(photos, m.FindAllPhotos(usernames)...)
	}

	sort.Slice(photos, func(i, j int) bool {
		return photos[i].UploadDate.Before(photos[j].UploadDate)
	})
	homepage := new(model.UserHomePage)
	homepage.Id = 0
	user := m.findUserByUsername(username)
	homepage.User = &user
	homepage.Followees = follows
	homepage.Photos = photos
	return *homepage
}

func (m appdbmemimpl) findUserByUsername(username string) model.User {
	user := new(model.User)
	user.Id = m.userIdsMap[username]
	user.Username = username
	return *user
}

func (m appdbmemimpl) FindPhoto(photoId int64) (model.Photo, error) {
	photo, ok := m.photosMaps[photoId]
	if ok {
		photo.User.Username = m.usersMap[photo.User.Id].Username
		photo.Likes = m.FindAllLikes(photo.Id)
		photo.Comments = m.FindAllComments(photo.Id)
	} else {
		return photo, errors.New("empty name")
	}
	return photo, nil
}

func (m appdbmemimpl) FindAllPhotos(username string) []model.Photo {
	var photos []model.Photo

	if username != "" {
		for _, photo := range m.photosMaps {
			if photo.User.Id == m.userIdsMap[username] {
				photo.User.Username = m.usersMap[photo.User.Id].Username
				photo.Likes = m.FindAllLikes(photo.Id)
				photo.Comments = m.FindAllComments(photo.Id)
				photos = append(photos, photo)
			}
		}
	} else {
		panic("username is mandatory")
	}
	if len(photos) == 0 {
		photos = make([]model.Photo, 0)
	}
	return photos

}

func (m appdbmemimpl) DeletePhoto(requestorUser string, id int64) {
	elem, ok := m.photosMaps[id]
	if ok && elem.User.Username == requestorUser {
		likes := m.FindAllLikes(id)
		for _, like := range likes {
			delete(m.likesMap, like.Id)
		}
		comments := m.FindAllComments(id)
		for _, comment := range comments {
			delete(m.commentsMap, comment.Id)
		}
		delete(m.photosMaps, id)
	}
}

func (m appdbmemimpl) SavePhoto(username string, bytes []byte) model.Photo {
	photo := new(model.Photo)
	user := m.usersMap[m.userIdsMap[username]]
	photo.User = new(model.User)
	photo.User.Username = user.Username
	photo.User.Id = user.Id
	photo.Data = bytes
	photo.UploadDate = time.Now()
	photo.Id = m.incrementAndGet()
	photo.Link = "/photos/" + strconv.Itoa(int(photo.Id))
	m.photosMaps[photo.Id] = *photo
	return *photo
}

func (m appdbmemimpl) FindAllLikes(photoId int64) []model.Like {
	var likes []model.Like
	if photoId > 0 {
		for _, like := range m.likesMap {
			if like.PhotoId == photoId {
				like.User.Username = m.usersMap[like.User.Id].Username
				likes = append(likes, like)
			}
		}
	}
	if len(likes) == 0 {
		likes = make([]model.Like, 0)
	}
	return likes
}

func (m appdbmemimpl) existsLike(likeRequest model.LikeRequest) (bool, model.Like) {
	for _, e := range m.likesMap {
		if e.User.Username == likeRequest.User.Username && e.PhotoId == likeRequest.PhotoId {
			return true, e
		}
	}
	return false, model.Like{}
}

func (m appdbmemimpl) SaveLike(likeRequest model.LikeRequest) model.Like {
	ok, like := m.existsLike(likeRequest)
	if ok {
		return like
	}
	like.Id = m.incrementAndGet()
	like.User = likeRequest.User
	like.PhotoId = likeRequest.PhotoId
	m.likesMap[like.Id] = like
	return like
}

func (m appdbmemimpl) DeleteLike(requestorUser string, id int64) {
	elem, ok := m.likesMap[id]
	if ok && elem.User.Username == requestorUser {
		delete(m.likesMap, id)
	}
}

func (m appdbmemimpl) DeleteComment(requestorUser string, id int64) {
	elem, ok := m.commentsMap[id]
	if ok && elem.User.Username == requestorUser {
		delete(m.commentsMap, id)
	}
}

func (m appdbmemimpl) SaveComment(commentRequest model.CommentRequest) model.Comment {
	comment := new(model.Comment)
	comment.Id = m.incrementAndGet()
	comment.User = commentRequest.User
	comment.PhotoId = commentRequest.PhotoId
	comment.Text = commentRequest.Text
	m.commentsMap[comment.Id] = *comment
	return *comment
}

func (m appdbmemimpl) FindAllComments(photoId int64) []model.Comment {
	var comments []model.Comment
	if photoId > 0 {
		for _, comment := range m.commentsMap {
			if comment.PhotoId == photoId {
				comment.User.Username = m.usersMap[comment.User.Id].Username
				comments = append(comments, comment)
			}
		}
	}
	if len(comments) == 0 {
		comments = make([]model.Comment, 0)
	}
	return comments

}

func (m appdbmemimpl) DeleteFollow(requestorUser string, id int64) {
	elem, ok := m.followsMap[id]
	if ok && elem.User.Username == requestorUser {
		delete(m.followsMap, id)
	}
}

func (m appdbmemimpl) existsFollow(followRequest model.FollowRequest) (bool, model.Follow) {
	for _, e := range m.followsMap {
		if e.User.Username == followRequest.User.Username && e.Followee.Username == followRequest.Followee.Username {
			return true, e
		}
	}
	return false, model.Follow{}
}

func (m appdbmemimpl) SaveFollow(followRequest model.FollowRequest) model.Follow {
	ok, follow := m.existsFollow(followRequest)
	if ok {
		return follow
	}
	follow.Id = m.incrementAndGet()
	follow.User = followRequest.User
	follow.User.Id = m.userIdsMap[followRequest.User.Username]
	follow.Followee = followRequest.Followee
	follow.Followee.Id = m.userIdsMap[followRequest.Followee.Username]
	m.followsMap[follow.Id] = follow
	return follow
}

func (m appdbmemimpl) FindAllFollow(username string) []model.Follow {
	var follows []model.Follow
	if username != "" {
		for _, follow := range m.followsMap {
			if follow.User.Id == m.userIdsMap[username] {
				follow.User.Username = m.usersMap[follow.User.Id].Username
				follow.Followee.Username = m.usersMap[follow.Followee.Id].Username
				follows = append(follows, follow)
			}
		}
	} else {
		for _, follow := range m.followsMap {
			follows = append(follows, follow)
		}
	}
	if len(follows) == 0 {
		follows = make([]model.Follow, 0)
	}
	return follows
}

func (m appdbmemimpl) findAllFollowers(username string) []model.Follow {
	var follows []model.Follow
	if username != "" {
		for _, follow := range m.followsMap {
			if follow.Followee.Id == m.userIdsMap[username] {
				follow.User.Username = m.usersMap[follow.User.Id].Username
				follow.Followee.Username = m.usersMap[follow.Followee.Id].Username
				follows = append(follows, follow)
			}
		}
	}
	if len(follows) == 0 {
		follows = make([]model.Follow, 0)
	}
	return follows
}

func (m appdbmemimpl) DeleteBan(requestorUser string, id int64) {
	elem, ok := m.bansMap[id]
	if ok && elem.User.Username == requestorUser {
		delete(m.bansMap, id)
	}
}

func (m appdbmemimpl) incrementAndGet() int64 {
	m.sequence[0] = m.sequence[0] + 1
	return m.sequence[0]
}

func (m appdbmemimpl) SaveBan(banRequest model.BanRequest) model.Ban {
	ok, ban := m.existsBan(banRequest)
	if ok {
		return ban
	}
	ban.Id = m.incrementAndGet()
	ban.User = banRequest.User
	ban.Banned = banRequest.Banned
	m.bansMap[ban.Id] = ban
	return ban
}

func (m appdbmemimpl) existsBan(banRequest model.BanRequest) (bool, model.Ban) {
	for _, e := range m.bansMap {
		if e.User.Username == banRequest.User.Username && e.Banned.Username == banRequest.Banned.Username {
			return true, e
		}
	}
	return false, model.Ban{}
}

func (m appdbmemimpl) FindAllBans(username string) []model.Ban {
	var bans []model.Ban
	if username != "" {
		for _, ban := range m.bansMap {
			if ban.User.Id == m.userIdsMap[username] {
				ban.User.Username = m.usersMap[ban.User.Id].Username
				ban.Banned.Username = m.usersMap[ban.Banned.Id].Username
				bans = append(bans, ban)
			}
		}
	}
	if len(bans) == 0 {
		bans = make([]model.Ban, 0)
	}
	return bans
}

func (m appdbmemimpl) findAllBanners(username string) []model.Ban {
	var bans []model.Ban
	if username != "" {
		for _, ban := range m.bansMap {
			if ban.Banned.Id == m.userIdsMap[username] {
				ban.User.Username = m.usersMap[ban.User.Id].Username
				ban.Banned.Username = m.usersMap[ban.Banned.Id].Username
				bans = append(bans, ban)
			}
		}
	}
	if len(bans) == 0 {
		bans = make([]model.Ban, 0)
	}
	return bans
}

func (m appdbmemimpl) FindAllUsers() []model.User {
	var users []model.User
	for e := range m.usersMap {
		users = append(users, m.usersMap[e])
	}
	return users
}

func (m appdbmemimpl) SaveUser(username string) model.User {
	e, ok := m.userIdsMap[username]
	if ok {
		return m.usersMap[e]
	} else {
		user := new(model.User)
		user.Id = m.incrementAndGet()
		user.Username = username
		m.usersMap[user.Id] = *user
		m.userIdsMap[username] = user.Id
		return *user
	}
}

func NewMem() (AppDatabaseMemory, error) {

	return &appdbmemimpl{
		usersMap:    make(map[int64]model.User),
		userIdsMap:  make(map[string]int64),
		bansMap:     make(map[int64]model.Ban),
		followsMap:  make(map[int64]model.Follow),
		likesMap:    make(map[int64]model.Like),
		photosMaps:  make(map[int64]model.Photo),
		commentsMap: make(map[int64]model.Comment),
		sequence:    map[int64]int64{0: 1},
	}, nil
}
