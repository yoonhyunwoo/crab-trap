package moltbook

import "time"

type Agent struct {
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	Karma           int       `json:"karma"`
	FollowerCount   int       `json:"follower_count"`
	FollowingCount  int       `json:"following_count"`
	IsClaimed       bool      `json:"is_claimed"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	LastActive      time.Time `json:"last_active"`
	Owner           *Owner    `json:"owner,omitempty"`
}

type Owner struct {
	XHandle        string `json:"x_handle"`
	XName          string `json:"x_name"`
	XAvatar        string `json:"x_avatar,omitempty"`
	XBio           string `json:"x_bio,omitempty"`
	XFollowerCount int    `json:"x_follower_count"`
	XFollowingCount int    `json:"x_following_count"`
	XVerified      bool   `json:"x_verified"`
}

type RegisterRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type RegisterResponse struct {
	Agent        Agent   `json:"agent"`
	Important    string  `json:"important"`
}

type StatusResponse struct {
	Status string `json:"status"`
}

type AgentProfileResponse struct {
	Success     bool     `json:"success"`
	Agent       Agent    `json:"agent"`
	RecentPosts []Post   `json:"recentPosts,omitempty"`
}

type UpdateAgentRequest struct {
	Description string                 `json:"description,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

type Post struct {
	ID           string      `json:"id"`
	Title        string      `json:"title"`
	Content      string      `json:"content,omitempty"`
	URL          string      `json:"url,omitempty"`
	Upvotes      int         `json:"upvotes"`
	Downvotes    int         `json:"downvotes"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at,omitempty"`
	Author       *Agent      `json:"author"`
	Submolt      *Submolt    `json:"submolt"`
	Comments     []Comment   `json:"comments,omitempty"`
	IsPinned     bool        `json:"is_pinned,omitempty"`
}

type CreatePostRequest struct {
	Submolt   string `json:"submolt"`
	Title     string `json:"title"`
	Content   string `json:"content,omitempty"`
	URL       string `json:"url,omitempty"`
}

type CreatePostResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Post    Post   `json:"post"`
}

type GetPostsOptions struct {
	Submolt string `json:"submolt,omitempty"`
	Sort    string `json:"sort,omitempty"`
	Limit   int    `json:"limit,omitempty"`
}

type PostsResponse struct {
	Success bool   `json:"success"`
	Posts   []Post `json:"posts"`
}

type Comment struct {
	ID          string     `json:"id"`
	Content     string     `json:"content"`
	Upvotes     int        `json:"upvotes"`
	Downvotes   int        `json:"downvotes"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at,omitempty"`
	Author      *Agent     `json:"author"`
	PostID      string     `json:"post_id,omitempty"`
	ParentID    *string    `json:"parent_id,omitempty"`
	Replies     []Comment  `json:"replies,omitempty"`
}

type CreateCommentRequest struct {
	Content  string `json:"content"`
	ParentID string `json:"parent_id,omitempty"`
}

type CreateCommentResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Comment Comment `json:"comment"`
}

type GetCommentsOptions struct {
	Sort string `json:"sort,omitempty"`
}

type CommentsResponse struct {
	Success  bool     `json:"success"`
	Comments []Comment `json:"comments"`
}

type VoteResponse struct {
	Success          bool   `json:"success"`
	Message          string `json:"message"`
	Author           Agent  `json:"author"`
	AlreadyFollowing bool   `json:"already_following"`
	Suggestion       string `json:"suggestion,omitempty"`
}

type Submolt struct {
	Name         string    `json:"name"`
	DisplayName  string    `json:"display_name"`
	Description  string    `json:"description"`
	BannerColor  string    `json:"banner_color,omitempty"`
	ThemeColor   string    `json:"theme_color,omitempty"`
	Avatar       string    `json:"avatar,omitempty"`
	Banner       string    `json:"banner,omitempty"`
	MemberCount  int       `json:"member_count"`
	IsSubscribed bool      `json:"is_subscribed,omitempty"`
	YourRole     string    `json:"your_role,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

type CreateSubmoltRequest struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
}

type SubmoltResponse struct {
	Success bool    `json:"success"`
	Submolt Submolt `json:"submolt"`
}

type SubmoltsResponse struct {
	Success  bool     `json:"success"`
	Submolts []Submolt `json:"submolts"`
}

type UpdateSubmoltRequest struct {
	Description string `json:"description,omitempty"`
	BannerColor string `json:"banner_color,omitempty"`
	ThemeColor  string `json:"theme_color,omitempty"`
}

type AddModeratorRequest struct {
	AgentName string `json:"agent_name"`
	Role      string `json:"role"`
}

type RemoveModeratorRequest struct {
	AgentName string `json:"agent_name"`
}

type ModeratorsResponse struct {
	Success    bool     `json:"success"`
	Moderators []Agent  `json:"moderators"`
}

type FollowResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type SearchResult struct {
	ID         string    `json:"id"`
	Type       string    `json:"type"`
	Title      *string   `json:"title,omitempty"`
	Content    string    `json:"content"`
	Upvotes    int       `json:"upvotes"`
	Downvotes  int       `json:"downvotes"`
	Similarity float64   `json:"similarity"`
	CreatedAt  time.Time `json:"created_at"`
	Author     Agent     `json:"author"`
	Post       *Post     `json:"post,omitempty"`
	Submolt    *Submolt  `json:"submolt,omitempty"`
	PostID     string    `json:"post_id,omitempty"`
}

type SearchRequest struct {
	Query string `json:"q"`
	Type  string `json:"type,omitempty"`
	Limit int    `json:"limit,omitempty"`
}

type SearchResponse struct {
	Success bool          `json:"success"`
	Query   string        `json:"query"`
	Type    string        `json:"type"`
	Results []SearchResult `json:"results"`
	Count   int           `json:"count"`
}

type FeedOptions struct {
	Sort  string `json:"sort,omitempty"`
	Limit int    `json:"limit,omitempty"`
}

type FeedResponse struct {
	Success bool   `json:"success"`
	Posts   []Post `json:"posts"`
}

type PinResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
