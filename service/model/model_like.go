/*
 * WASAPhoto
 *
 * Each user will be presented with a stream of photos (images) in reverse chronological order, with information about when each photo was uploaded (date and time  and how many likes and comments it has. The stream is composed by photos from “following” (other users that the user follows). Users can place (and later remove) a “like” to photos from other users. Also, users can add comments to any image (even those uploaded by themself).  Only authors can remove their comments.Users can ban other users. If user Alice bans user Eve, Eve won’t be able to see any information about Alice. Alice can decide to remove the ban at any moment.Users will have their profiles.   The personal profile page for the user shows: the user’s photos (in reverse chronological order), how many photos have been uploaded, and the user’s followers and following. Users can change their usernames, upload photos, remove photos, and follow/unfollow other users.  Removal of an image will also remove likes and comments.A user can search other user profiles via username. A user can log in just by specifying the username. See the “Simplified login” section for details. operations availables:  (Mandatory) - doLogin (see simplified login) - setMyUserName - uploadPhoto - followUser - unfollowUser - banUser - unbanUser - likePhoto - unlikePhoto - commentPhoto - uncommentPhoto - deletePhoto - getUserProfile - getMyStream  Some useful links: - [project spec ](http://gamificationlab.uniroma1.it/notes/Project.pdf)
 *
 * API version: 0.0.1
 * Contact: merlini.1166162@studenti.uniroma1.it
 * Generated by: model Codegen (https://github.com/model-api/model-codegen.git)
 */
package model

// the like object
type Like struct {
	// the like id
	Id int64 `json:"id,omitempty"`

	User *User `json:"user,omitempty"`
	// the photo id
	PhotoId int64 `json:"photoId,omitempty"`
}