package akismet

import "net/url"

type Comment struct {
	Blog               string // required: if not provided then use Akismet.Blog
	UserIp             string // required: IP address of the comment submitter
	UserAgent          string // required: User agent string of the web browser submitting the comment - typically the HTTP_USER_AGENT cgi variable. Not to be confused with the user agent of your Akismet library.
	Referrer           string // optional: The content of the HTTP_REFERER header should be sent here.
	Permalink          string // optional: The permanent location of the entry the comment was submitted to.
	CommentType        string // optional: May be blank, comment, trackback, pingback, or a made up value like "registration".
	CommentAuthor      string // optional: Name submitted with the comment
	CommentAuthorEmail string // optional: Email address submitted with the comment
	CommentAuthorUrl   string // optional: URL submitted with comment
	CommentContent     string // optional: The content that was submitted.
}

func (c Comment) ToValues(vals url.Values) {
	// Required fields
	if c.Blog != "" {
		vals.Set("blog", c.Blog)
	} else {
		vals.Set("blog", c.Blog)
	}
	vals.Set("user_ip", c.UserIp)
	vals.Set("user_agent", c.UserAgent)
	// Optional fields
	if c.Referrer != "" {
		vals.Set("referrer", c.Referrer)
	}
	if c.Permalink != "" {
		vals.Set("permalink", c.Permalink)
	}
	if c.CommentType != "" {
		vals.Set("comment_type", c.CommentType)
	}
	if c.CommentAuthor != "" {
		vals.Set("comment_author", c.CommentAuthor)
	}
	if c.CommentAuthorEmail != "" {
		vals.Set("comment_author_email", c.CommentAuthorEmail)
	}
	if c.CommentAuthorUrl != "" {
		vals.Set("comment_author_url", c.CommentAuthorUrl)
	}
	if c.CommentContent != "" {
		vals.Set("comment_content", c.CommentContent)
	}
}
