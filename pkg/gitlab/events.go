// Package gitlab for gitlab events
package gitlab

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// NullSHA null SHA
const NullSHA = "0000000000000000000000000000000000000000"

// Gitlab header
const (
	HeaderEvent = "X-Gitlab-Event"
	HeaderToken = "X-Gitlab-Token" // nolint: gosec
)

// Ci/CD status
const (
	StatusCreated  = "created"
	StatusPending  = "pending"
	StatusRunning  = "running"
	StatusCanceled = "canceled"
	StatusFailed   = "failed"
	StatusSuccess  = "success"
)

// EventType represents a Gitlab event type.
type EventType string

// List of available event types.
const (
	EventTypePush              EventType = "Push Hook"
	EventTypeTagPush           EventType = "Tag Push Hook"
	EventTypeIssue             EventType = "Issue Hook"
	EventTypeConfidentialIssue EventType = "Confidential Issue Hook"
	EventTypeNote              EventType = "Note Hook"
	EventConfidentialTypeNote  EventType = "Confidential Note Hook"
	EventTypeMergeRequest      EventType = "Merge Request Hook"
	EventTypeJob               EventType = "Job Hook"
	EventTypePipeline          EventType = "Pipeline Hook"
	EventTypeBuild             EventType = "Build Hook"
	EventTypeWikiPage          EventType = "Wiki Page Hook"
	EventTypeSystemHook        EventType = "System Hook"
)

// Noteable const
const (
	NoteableTypeCommit       = "Commit"
	NoteableTypeMergeRequest = "MergeRequest"
	NoteableTypeIssue        = "Issue"
	NoteableTypeSnippet      = "Snippet"
)

// EventNote represents a comments event.
type EventNote struct {
	ObjectKind       string     `json:"object_kind"`
	User             User       `json:"user"`
	ProjectID        int        `json:"project_id"`
	Project          Project    `json:"project"`
	Repository       Repository `json:"repository"`
	ObjectAttributes struct {
		ID           int         `json:"id"`
		Note         string      `json:"note"`
		NoteableType string      `json:"noteable_type"`
		AuthorID     int         `json:"author_id"`
		CreatedAt    string      `json:"created_at"`
		UpdatedAt    string      `json:"updated_at"`
		ProjectID    int         `json:"project_id"`
		Attachment   interface{} `json:"attachment"`
		LineCode     string      `json:"line_code"`
		CommitID     string      `json:"commit_id"`
		NoteableID   interface{} `json:"noteable_id"`
		System       bool        `json:"system"`
		StDiff       Diff        `json:"st_diff"`
		URL          string      `json:"url"`
	} `json:"object_attributes"`
	Commit struct {
		ID        string    `json:"id"`
		Message   string    `json:"message"`
		Timestamp time.Time `json:"timestamp"`
		URL       string    `json:"url"`
		Author    struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"author"`
	} `json:"commit"`
	Issue   Issue `json:"issue"`
	Snippet struct {
		ID                     int         `json:"id"`
		Title                  string      `json:"title"`
		Content                string      `json:"content"`
		AuthorID               int         `json:"author_id"`
		ProjectID              int         `json:"project_id"`
		CreatedAt              string      `json:"created_at"`
		UpdatedAt              string      `json:"updated_at"`
		FileName               string      `json:"file_name"`
		Type                   string      `json:"type"`
		VisibilityLevel        int         `json:"visibility_level"`
		Description            string      `json:"description"`
		EncryptedSecretToken   interface{} `json:"encrypted_secret_token"`
		EncryptedSecretTokenIv interface{} `json:"encrypted_secret_token_iv"`
		Secret                 bool        `json:"secret"`
		SecretToken            interface{} `json:"secret_token"`
	} `json:"snippet"`
	MergeRequest struct {
		ID                        int         `json:"id"`
		TargetBranch              string      `json:"target_branch"`
		SourceBranch              string      `json:"source_branch"`
		SourceProjectID           int         `json:"source_project_id"`
		AuthorID                  int         `json:"author_id"`
		AssigneeID                int         `json:"assignee_id"`
		Title                     string      `json:"title"`
		CreatedAt                 string      `json:"created_at"`
		UpdatedAt                 string      `json:"updated_at"`
		MilestoneID               int         `json:"milestone_id"`
		State                     string      `json:"state"`
		MergeStatus               string      `json:"merge_status"`
		TargetProjectID           int         `json:"target_project_id"`
		IID                       int         `json:"iid"`
		Description               string      `json:"description"`
		Position                  int         `json:"position"`
		LockedAt                  string      `json:"locked_at"`
		UpdatedByID               int         `json:"updated_by_id"`
		MergeError                string      `json:"merge_error"`
		MergeParams               MergeParams `json:"merge_params"`
		MergeWhenPipelineSucceeds bool        `json:"merge_when_pipeline_succeeds"`
		Squash                    bool        `json:"squash"`
		WorkInProgress            bool        `json:"work_in_progress"`
		MergeUserID               int         `json:"merge_user_id"`
		MergeCommitSHA            string      `json:"merge_commit_sha"`
		DeletedAt                 string      `json:"deleted_at"`
		InProgressMergeCommitSHA  string      `json:"in_progress_merge_commit_sha"`
		LockVersion               int         `json:"lock_version"`
		ApprovalsBeforeMerge      string      `json:"approvals_before_merge"`
		RebaseCommitSHA           string      `json:"rebase_commit_sha"`
		TimeEstimate              int         `json:"time_estimate"`
		LastEditedAt              string      `json:"last_edited_at"`
		LastEditedByID            int         `json:"last_edited_by_id"`
		Source                    Repository  `json:"source"`
		Target                    Repository  `json:"target"`
		LastCommit                struct {
			ID        string    `json:"id"`
			Message   string    `json:"message"`
			Timestamp time.Time `json:"timestamp"`
			URL       string    `json:"url"`
			Author    struct {
				Name  string `json:"name"`
				Email string `json:"email"`
			} `json:"author"`
		} `json:"last_commit"`
		TotalTimeSpent int `json:"total_time_spent"`
		HeadPipelineID int `json:"head_pipeline_id"`
	} `json:"merge_request"`
}

// EventPush represents a push event.
//
// GitLab API docs:
// https://docs.gitlab.com/ce/user/project/integrations/webhooks.html#push-events
type EventPush struct {
	ObjectKind   string     `json:"object_kind"`
	Before       string     `json:"before"`
	After        string     `json:"after"`
	Ref          string     `json:"ref"`
	CheckoutSHA  string     `json:"checkout_sha"`
	UserID       int        `json:"user_id"`
	UserName     string     `json:"user_name"`
	UserUsername string     `json:"user_username"`
	UserEmail    string     `json:"user_email"`
	UserAvatar   string     `json:"user_avatar"`
	ProjectID    int        `json:"project_id"`
	Project      Project    `json:"project"`
	Repository   Repository `json:"repository"`
	Commits      []struct {
		ID        string    `json:"id"`
		Message   string    `json:"message"`
		Timestamp time.Time `json:"timestamp"`
		URL       string    `json:"url"`
		Author    struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"author"`
		Added    []string `json:"added"`
		Modified []string `json:"modified"`
		Removed  []string `json:"removed"`
	} `json:"commits"`
	TotalCommitsCount int `json:"total_commits_count"`
}

// EventTagPush represents a tag event.
//
// GitLab API docs:
// https://docs.gitlab.com/ce/user/project/integrations/webhooks.html#tag-events
type EventTagPush struct {
	ObjectKind  string     `json:"object_kind"`
	Before      string     `json:"before"`
	After       string     `json:"after"`
	Ref         string     `json:"ref"`
	CheckoutSHA string     `json:"checkout_sha"`
	UserID      int        `json:"user_id"`
	UserName    string     `json:"user_name"`
	UserAvatar  string     `json:"user_avatar"`
	ProjectID   int        `json:"project_id"`
	Message     string     `json:"message"`
	Project     Project    `json:"project"`
	Repository  Repository `json:"repository"`
	Commits     []struct {
		ID        string    `json:"id"`
		Message   string    `json:"message"`
		Timestamp time.Time `json:"timestamp"`
		URL       string    `json:"url"`
		Author    struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"author"`
		Added    []string `json:"added"`
		Modified []string `json:"modified"`
		Removed  []string `json:"removed"`
	} `json:"commits"`
	TotalCommitsCount int `json:"total_commits_count"`
}

// EventIssue represents a issue event.
//
// GitLab API docs:
// https://docs.gitlab.com/ce/user/project/integrations/webhooks.html#issues-events
type EventIssue struct {
	ObjectKind       string     `json:"object_kind"`
	User             User       `json:"user"`
	Project          Project    `json:"project"`
	Repository       Repository `json:"repository"`
	ObjectAttributes Issue      `json:"object_attributes"`
	Assignee         struct {
		Name      string `json:"name"`
		Username  string `json:"username"`
		AvatarURL string `json:"avatar_url"`
	} `json:"assignee"`
	Assignees []struct {
		Name      string `json:"name"`
		Username  string `json:"username"`
		AvatarURL string `json:"avatar_url"`
	} `json:"assignees"`
	Labels  []Label `json:"labels"`
	Changes struct {
		AuthorID struct {
			Previous int `json:"previous"`
			Current  int `json:"current"`
		} `json:"author_id"`
		CreatedAt struct {
			Previous string `json:"previous"`
			Current  string `json:"current"`
		} `json:"created_at"`
		Description struct {
			Previous string `json:"previous"`
			Current  string `json:"current"`
		} `json:"description"`
		DueDate struct {
			Previous string `json:"previous"`
			Current  string `json:"current"`
		} `json:"due_date"`
		ID struct {
			Previous int `json:"previous"`
			Current  int `json:"current"`
		} `json:"id"`
		Iid struct {
			Previous int `json:"previous"`
			Current  int `json:"current"`
		} `json:"iid"`
		MilestoneID struct {
			Previous int `json:"previous"`
			Current  int `json:"current"`
		} `json:"milestone_id"`
		ProjectID struct {
			Previous int `json:"previous"`
			Current  int `json:"current"`
		} `json:"project_id"`
		RelativePosition struct {
			Previous int `json:"previous"`
			Current  int `json:"current"`
		} `json:"relative_position"`
		Title struct {
			Previous string `json:"previous"`
			Current  string `json:"current"`
		} `json:"title"`
		UpdatedAt struct {
			Previous string `json:"previous"`
			Current  string `json:"current"`
		} `json:"updated_at"`
		Weight struct {
			Previous int `json:"previous"`
			Current  int `json:"current"`
		} `json:"weight"`
		Assignees struct {
			Previous []User `json:"previous"`
			Current  []User `json:"current"`
		} `json:"assignees"`
		Labels struct {
			Previous []Label `json:"previous"`
			Current  []Label `json:"current"`
		} `json:"labels"`
		UpdatedByID struct {
			Previous int `json:"previous"`
			Current  int `json:"current"`
		} `json:"updated_by_id"`
		TotalTimeSpent struct {
			Previous int `json:"previous"`
			Current  int `json:"current"`
		} `json:"total_time_spent"`
	} `json:"changes"`
}

// Issue type
type Issue struct {
	ID                  int           `json:"id"`
	Title               string        `json:"title"`
	AssigneeID          int           `json:"assignee_id"`
	AuthorID            int           `json:"author_id"`
	ProjectID           int           `json:"project_id"`
	CreatedAt           string        `json:"created_at"` // Should be time.Time (see Gitlab issue #21468)
	UpdatedAt           string        `json:"updated_at"` // Should be time.Time (see Gitlab issue #21468)
	Position            int           `json:"position"`
	BranchName          string        `json:"branch_name"`
	Description         string        `json:"description"`
	MilestoneID         int           `json:"milestone_id"`
	State               IssueState    `json:"state"`
	IID                 int           `json:"iid"`
	URL                 string        `json:"url"`
	Action              IssueAction   `json:"action"`
	ClosedAt            interface{}   `json:"closed_at"`
	Confidential        bool          `json:"confidential"`
	DueDate             interface{}   `json:"due_date"`
	LastEditedAt        interface{}   `json:"last_edited_at"`
	LastEditedByID      interface{}   `json:"last_edited_by_id"`
	MovedToID           interface{}   `json:"moved_to_id"`
	DuplicatedToID      interface{}   `json:"duplicated_to_id"`
	RelativePosition    int           `json:"relative_position"`
	StateID             int           `json:"state_id"`
	TimeEstimate        int           `json:"time_estimate"`
	UpdatedByID         interface{}   `json:"updated_by_id"`
	Weight              int           `json:"weight"`
	TotalTimeSpent      int           `json:"total_time_spent"`
	HumanTotalTimeSpent interface{}   `json:"human_total_time_spent"`
	HumanTimeEstimate   interface{}   `json:"human_time_estimate"`
	AssigneeIDs         []int         `json:"assignee_ids"`
	Labels              []interface{} `json:"labels"`
}

// IssueState type
type IssueState string

// IssueState const
const (
	IssueStateOpened IssueState = "opened"
	IssueStateClosed IssueState = "closed"
)

// IssueAction type
type IssueAction string

// IssueAction const
const (
	IssueActionOpen   IssueAction = "open"
	IssueActionUpdate IssueAction = "update"
	IssueActionClose  IssueAction = "close"
	IssueActionReopen IssueAction = "reopen"
)

// EventJob represents a job event.
//
// GitLab API docs:
// https://gitlab.com/help/user/project/integrations/webhooks#job-events
type EventJob struct {
	ObjectKind        string  `json:"object_kind"`
	Ref               string  `json:"ref"`
	BeforeSHA         string  `json:"before_sha"`
	SHA               string  `json:"sha"`
	BuildID           int     `json:"build_id"`
	BuildName         string  `json:"build_name"`
	BuildStage        string  `json:"build_stage"`
	BuildStatus       string  `json:"build_status"`
	BuildStartedAt    string  `json:"build_started_at"`
	BuildFinishedAt   string  `json:"build_finished_at"`
	BuildDuration     float64 `json:"build_duration"`
	BuildAllowFailure bool    `json:"build_allow_failure"`
	Tag               bool    `json:"tag"`
	PipelineID        int     `json:"pipeline_id"`
	ProjectID         int     `json:"project_id"`
	ProjectName       string  `json:"project_name"`
	User              struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"user"`
	Commit struct {
		ID          int    `json:"id"`
		SHA         string `json:"sha"`
		Message     string `json:"message"`
		AuthorName  string `json:"author_name"`
		AuthorEmail string `json:"author_email"`
		AuthorURL   string `json:"author_url"`
		Status      string `json:"status"`
		Duration    int    `json:"duration"`
		StartedAt   string `json:"started_at"`
		FinishedAt  string `json:"finished_at"`
	} `json:"commit"`
	Repository Repository `json:"repository"`
}

// EventCommitComment represents a comment on a commit event.
//
// GitLab API docs:
// https://docs.gitlab.com/ce/user/project/integrations/webhooks.html#comment-on-commit
type EventCommitComment struct {
	ObjectKind       string     `json:"object_kind"`
	User             User       `json:"user"`
	ProjectID        int        `json:"project_id"`
	Project          Project    `json:"project"`
	Repository       Repository `json:"repository"`
	ObjectAttributes struct {
		ID           int    `json:"id"`
		Note         string `json:"note"`
		NoteableType string `json:"noteable_type"`
		AuthorID     int    `json:"author_id"`
		CreatedAt    string `json:"created_at"`
		UpdatedAt    string `json:"updated_at"`
		ProjectID    int    `json:"project_id"`
		Attachment   string `json:"attachment"`
		LineCode     string `json:"line_code"`
		CommitID     string `json:"commit_id"`
		NoteableID   int    `json:"noteable_id"`
		System       bool   `json:"system"`
		StDiff       struct {
			Diff        string `json:"diff"`
			NewPath     string `json:"new_path"`
			OldPath     string `json:"old_path"`
			AMode       string `json:"a_mode"`
			BMode       string `json:"b_mode"`
			NewFile     bool   `json:"new_file"`
			RenamedFile bool   `json:"renamed_file"`
			DeletedFile bool   `json:"deleted_file"`
		} `json:"st_diff"`
	} `json:"object_attributes"`
	Commit struct {
		ID        string    `json:"id"`
		Message   string    `json:"message"`
		Timestamp time.Time `json:"timestamp"`
		URL       string    `json:"url"`
		Author    struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"author"`
	} `json:"commit"`
}

// EventMergeComment represents a comment on a merge event.
//
// GitLab API docs:
// https://docs.gitlab.com/ce/user/project/integrations/webhooks.html#comment-on-merge-request
type EventMergeComment struct {
	ObjectKind       string  `json:"object_kind"`
	User             User    `json:"user"`
	ProjectID        int     `json:"project_id"`
	Project          Project `json:"project"`
	ObjectAttributes struct {
		ID           int    `json:"id"`
		DiscussionID string `json:"discussion_id"`
		Note         string `json:"note"`
		NoteableType string `json:"noteable_type"`
		AuthorID     int    `json:"author_id"`
		CreatedAt    string `json:"created_at"`
		UpdatedAt    string `json:"updated_at"`
		ProjectID    int    `json:"project_id"`
		Attachment   string `json:"attachment"`
		LineCode     string `json:"line_code"`
		CommitID     string `json:"commit_id"`
		NoteableID   int    `json:"noteable_id"`
		System       bool   `json:"system"`
		StDiff       Diff   `json:"st_diff"`
		URL          string `json:"url"`
	} `json:"object_attributes"`
	Repository   Repository `json:"repository"`
	MergeRequest struct {
		ID                        int         `json:"id"`
		TargetBranch              string      `json:"target_branch"`
		SourceBranch              string      `json:"source_branch"`
		SourceProjectID           int         `json:"source_project_id"`
		AuthorID                  int         `json:"author_id"`
		AssigneeID                int         `json:"assignee_id"`
		Title                     string      `json:"title"`
		CreatedAt                 string      `json:"created_at"`
		UpdatedAt                 string      `json:"updated_at"`
		MilestoneID               int         `json:"milestone_id"`
		State                     string      `json:"state"`
		MergeStatus               string      `json:"merge_status"`
		TargetProjectID           int         `json:"target_project_id"`
		IID                       int         `json:"iid"`
		Description               string      `json:"description"`
		Position                  int         `json:"position"`
		LockedAt                  string      `json:"locked_at"`
		UpdatedByID               int         `json:"updated_by_id"`
		MergeError                string      `json:"merge_error"`
		MergeParams               MergeParams `json:"merge_params"`
		MergeWhenPipelineSucceeds bool        `json:"merge_when_pipeline_succeeds"`
		Squash                    bool        `json:"squash"`
		WorkInProgress            bool        `json:"work_in_progress"`
		MergeUserID               int         `json:"merge_user_id"`
		MergeCommitSHA            string      `json:"merge_commit_sha"`
		DeletedAt                 string      `json:"deleted_at"`
		InProgressMergeCommitSHA  string      `json:"in_progress_merge_commit_sha"`
		LockVersion               int         `json:"lock_version"`
		ApprovalsBeforeMerge      string      `json:"approvals_before_merge"`
		RebaseCommitSHA           string      `json:"rebase_commit_sha"`
		TimeEstimate              int         `json:"time_estimate"`
		LastEditedAt              string      `json:"last_edited_at"`
		LastEditedByID            int         `json:"last_edited_by_id"`
		Source                    Repository  `json:"source"`
		Target                    Repository  `json:"target"`
		LastCommit                struct {
			ID        string    `json:"id"`
			Message   string    `json:"message"`
			Timestamp time.Time `json:"timestamp"`
			URL       string    `json:"url"`
			Author    struct {
				Name  string `json:"name"`
				Email string `json:"email"`
			} `json:"author"`
		} `json:"last_commit"`
		TotalTimeSpent int `json:"total_time_spent"`
		HeadPipelineID int `json:"head_pipeline_id"`
	} `json:"merge_request"`
}

// EventIssueComment represents a comment on an issue event.
//
// GitLab API docs:
// https://docs.gitlab.com/ce/user/project/integrations/webhooks.html#comment-on-issue
type EventIssueComment struct {
	ObjectKind       string     `json:"object_kind"`
	User             User       `json:"user"`
	ProjectID        int        `json:"project_id"`
	Project          Project    `json:"project"`
	Repository       Repository `json:"repository"`
	ObjectAttributes struct {
		ID           int    `json:"id"`
		Note         string `json:"note"`
		NoteableType string `json:"noteable_type"`
		AuthorID     int    `json:"author_id"`
		CreatedAt    string `json:"created_at"`
		UpdatedAt    string `json:"updated_at"`
		ProjectID    int    `json:"project_id"`
		Attachment   string `json:"attachment"`
		LineCode     string `json:"line_code"`
		CommitID     string `json:"commit_id"`
		NoteableID   int    `json:"noteable_id"`
		System       bool   `json:"system"`
		StDiff       []Diff `json:"st_diff"`
		URL          string `json:"url"`
	} `json:"object_attributes"`
	Issue struct {
		ID             int    `json:"id"`
		IID            int    `json:"iid"`
		ProjectID      int    `json:"project_id"`
		MilestoneID    int    `json:"milestone_id"`
		AuthorID       int    `json:"author_id"`
		Description    string `json:"description"`
		State          string `json:"state"`
		Title          string `json:"title"`
		LastEditedAt   string `json:"last_edit_at"`
		LastEditedByID int    `json:"last_edited_by_id"`
		UpdatedAt      string `json:"updated_at"`
		UpdatedByID    int    `json:"updated_by_id"`
		CreatedAt      string `json:"created_at"`
		ClosedAt       string `json:"closed_at"`
		// NOTE: check this
		// DueDate             ISOTime `json:"due_date"`
		URL                 string `json:"url"`
		TimeEstimate        int    `json:"time_estimate"`
		Confidential        bool   `json:"confidential"`
		TotalTimeSpent      int    `json:"total_time_spent"`
		HumanTotalTimeSpent int    `json:"human_total_time_spent"`
		HumanTimeEstimate   int    `json:"human_time_estimate"`
		AssigneeIDs         []int  `json:"assignee_ids"`
		AssigneeID          int    `json:"assignee_id"`
	} `json:"issue"`
}

// EventSnippetComment represents a comment on a snippet event.
//
// GitLab API docs:
// https://docs.gitlab.com/ce/user/project/integrations/webhooks.html#comment-on-code-snippet
type EventSnippetComment struct {
	ObjectKind       string     `json:"object_kind"`
	User             User       `json:"user"`
	ProjectID        int        `json:"project_id"`
	Project          Project    `json:"project"`
	Repository       Repository `json:"repository"`
	ObjectAttributes struct {
		ID           int    `json:"id"`
		Note         string `json:"note"`
		NoteableType string `json:"noteable_type"`
		AuthorID     int    `json:"author_id"`
		CreatedAt    string `json:"created_at"`
		UpdatedAt    string `json:"updated_at"`
		ProjectID    int    `json:"project_id"`
		Attachment   string `json:"attachment"`
		LineCode     string `json:"line_code"`
		CommitID     string `json:"commit_id"`
		NoteableID   int    `json:"noteable_id"`
		System       bool   `json:"system"`
		StDiff       Diff   `json:"st_diff"`
		URL          string `json:"url"`
	} `json:"object_attributes"`
	Snippet Snippet `json:"snippet"`
}

// EventMergeRequest represents a merge event.
//
// GitLab API docs:
// https://docs.gitlab.com/ce/user/project/integrations/webhooks.html#merge-request-events
type EventMergeRequest struct {
	ObjectKind       string  `json:"object_kind"`
	User             User    `json:"user"`
	Project          Project `json:"project"`
	ObjectAttributes struct {
		ID              int    `json:"id"`
		TargetBranch    string `json:"target_branch"`
		SourceBranch    string `json:"source_branch"`
		SourceProjectID int    `json:"source_project_id"`
		AuthorID        int    `json:"author_id"`
		AssigneeID      int    `json:"assignee_id"`
		AssigneeIDs     []int  `json:"assignee_ids"`
		Title           string `json:"title"`
		CreatedAt       string `json:"created_at"` // Should be time.Time (see Gitlab issue #21468)
		UpdatedAt       string `json:"updated_at"` // Should be time.Time (see Gitlab issue #21468)
		// NOTE: check this:
		// StCommits                []Commit    `json:"st_commits"`
		StDiffs                  []Diff      `json:"st_diffs"`
		MilestoneID              int         `json:"milestone_id"`
		State                    string      `json:"state"`
		MergeStatus              string      `json:"merge_status"`
		TargetProjectID          int         `json:"target_project_id"`
		IID                      int         `json:"iid"`
		Description              string      `json:"description"`
		Position                 int         `json:"position"`
		LockedAt                 string      `json:"locked_at"`
		UpdatedByID              int         `json:"updated_by_id"`
		MergeError               string      `json:"merge_error"`
		MergeUserID              int         `json:"merge_user_id"`
		MergeCommitSHA           string      `json:"merge_commit_sha"`
		MergeParams              MergeParams `json:"merge_params"`
		MergeWhenBuildSucceeds   bool        `json:"merge_when_build_succeeds"`
		WorkInProgress           bool        `json:"work_in_progress"`
		DeletedAt                string      `json:"deleted_at"`
		ApprovalsBeforeMerge     string      `json:"approvals_before_merge"`
		RebaseCommitSHA          string      `json:"rebase_commit_sha"`
		InProgressMergeCommitSHA string      `json:"in_progress_merge_commit_sha"`
		LockVersion              int         `json:"lock_version"`
		TimeEstimate             int         `json:"time_estimate"`
		Source                   Repository  `json:"source"`
		Target                   Repository  `json:"target"`
		LastCommit               struct {
			ID        string    `json:"id"`
			Message   string    `json:"message"`
			Timestamp time.Time `json:"timestamp"`
			URL       string    `json:"url"`
			Author    struct {
				Name  string `json:"name"`
				Email string `json:"email"`
			} `json:"author"`
		} `json:"last_commit"`
		URL      string        `json:"url"`
		Action   string        `json:"action"`
		OldRev   string        `json:"oldrev"`
		Assignee MergeAssignee `json:"assignee"`
	} `json:"object_attributes"`
	Repository Repository    `json:"repository"`
	Assignee   MergeAssignee `json:"assignee"`
	Labels     []Label       `json:"labels"`
	Changes    struct {
		Assignees struct {
			Previous []MergeAssignee `json:"previous"`
			Current  []MergeAssignee `json:"current"`
		} `json:"assignees"`
		Description struct {
			Previous string `json:"previous"`
			Current  string `json:"current"`
		} `json:"description"`
		Labels struct {
			Previous []Label `json:"previous"`
			Current  []Label `json:"current"`
		} `json:"labels"`
		SourceBranch struct {
			Previous string `json:"previous"`
			Current  string `json:"current"`
		} `json:"source_branch"`
		SourceProjectID struct {
			Previous int `json:"previous"`
			Current  int `json:"current"`
		} `json:"source_project_id"`
		TargetBranch struct {
			Previous string `json:"previous"`
			Current  string `json:"current"`
		} `json:"target_branch"`
		TargetProjectID struct {
			Previous int `json:"previous"`
			Current  int `json:"current"`
		} `json:"target_project_id"`
		Title struct {
			Previous string `json:"previous"`
			Current  string `json:"current"`
		} `json:"title"`
		UpdatedByID struct {
			Previous int `json:"previous"`
			Current  int `json:"current"`
		} `json:"updated_by_id"`
	} `json:"changes"`
}

// MergeAssignee represents a merge assignee.
type MergeAssignee struct {
	Name      string `json:"name"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
}

// MergeParams represents the merge params.
type MergeParams struct {
	ForceRemoveSourceBranch bool `json:"force_remove_source_branch"`
}

// UnmarshalJSON decodes the merge parameters
//
// This allows support of ForceRemoveSourceBranch for both type bool (>11.9) and string (<11.9)
func (p MergeParams) UnmarshalJSON(b []byte) error {
	type Alias MergeParams

	raw := struct {
		Alias
		ForceRemoveSourceBranch interface{} `json:"force_remove_source_branch"`
	}{
		Alias: (Alias)(p),
	}

	err := json.Unmarshal(b, &raw)
	if err != nil {
		return err
	}

	switch v := raw.ForceRemoveSourceBranch.(type) {
	case nil:
		// No action needed.
	case bool:
		p.ForceRemoveSourceBranch = v
	case string:
		p.ForceRemoveSourceBranch, err = strconv.ParseBool(v)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("failed to unmarshal ForceRemoveSourceBranch of type: %T", v)
	}

	return nil
}

// EventWikiPage represents a wiki page event.
//
// GitLab API docs:
// https://docs.gitlab.com/ce/user/project/integrations/webhooks.html#wiki-page-events
type EventWikiPage struct {
	ObjectKind string  `json:"object_kind"`
	User       User    `json:"user"`
	Project    Project `json:"project"`
	Wiki       struct {
		WebURL            string `json:"web_url"`
		GitSSHURL         string `json:"git_ssh_url"`
		GitHTTPURL        string `json:"git_http_url"`
		PathWithNamespace string `json:"path_with_namespace"`
		DefaultBranch     string `json:"default_branch"`
	} `json:"wiki"`
	ObjectAttributes struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		Format  string `json:"format"`
		Message string `json:"message"`
		Slug    string `json:"slug"`
		URL     string `json:"url"`
		Action  string `json:"action"`
	} `json:"object_attributes"`
}

// EventPipeline represents a pipeline event.
//
// GitLab API docs:
// https://docs.gitlab.com/ce/user/project/integrations/webhooks.html#pipeline-events
type EventPipeline struct {
	ObjectKind       string `json:"object_kind"`
	ObjectAttributes struct {
		ID         int      `json:"id"`
		Ref        string   `json:"ref"`
		Tag        bool     `json:"tag"`
		SHA        string   `json:"sha"`
		BeforeSHA  string   `json:"before_sha"`
		Status     string   `json:"status"`
		Stages     []string `json:"stages"`
		CreatedAt  string   `json:"created_at"`
		FinishedAt string   `json:"finished_at"`
		Duration   int      `json:"duration"`
	} `json:"object_attributes"`
	MergeRequest struct {
		ID                 int    `json:"id"`
		IID                int    `json:"iid"`
		Title              string `json:"title"`
		SourceBranch       string `json:"source_branch"`
		SourceProjectID    int    `json:"source_project_id"`
		TargetBranch       string `json:"target_branch"`
		TargetProjectID    int    `json:"target_project_id"`
		State              string `json:"state"`
		MergeRequestStatus string `json:"merge_status"`
		URL                string `json:"url"`
	} `json:"merge_request"`
	User    User    `json:"user"`
	Project Project `json:"project"`
	Commit  struct {
		ID        string    `json:"id"`
		Message   string    `json:"message"`
		Timestamp time.Time `json:"timestamp"`
		URL       string    `json:"url"`
		Author    struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"author"`
	} `json:"commit"`
	Builds []struct {
		ID         int    `json:"id"`
		Stage      string `json:"stage"`
		Name       string `json:"name"`
		Status     string `json:"status"`
		CreatedAt  string `json:"created_at"`
		StartedAt  string `json:"started_at"`
		FinishedAt string `json:"finished_at"`
		When       string `json:"when"`
		Manual     bool   `json:"manual"`
		User       User   `json:"user"`
		Runner     struct {
			ID          int    `json:"id"`
			Description string `json:"description"`
			Active      bool   `json:"active"`
			IsShared    bool   `json:"is_shared"`
		} `json:"runner"`
		ArtifactsFile struct {
			Filename string `json:"filename"`
			Size     int    `json:"size"`
		} `json:"artifacts_file"`
	} `json:"builds"`
}

//EventBuild represents a build event
//
// GitLab API docs:
// https://docs.gitlab.com/ce/user/project/integrations/webhooks.html#build-events
type EventBuild struct {
	ObjectKind        string  `json:"object_kind"`
	Ref               string  `json:"ref"`
	BeforeSHA         string  `json:"before_sha"`
	SHA               string  `json:"sha"`
	BuildID           int     `json:"build_id"`
	Tag               bool    `json:"tag"`
	BuildAllowFailure bool    `json:"build_allow_failure"`
	BuildName         string  `json:"build_name"`
	BuildStage        string  `json:"build_stage"`
	BuildStatus       string  `json:"build_status"`
	BuildStartedAt    string  `json:"build_started_at"`
	BuildFinishedAt   string  `json:"build_finished_at"`
	BuildDuration     float64 `json:"build_duration"`
	ProjectID         int     `json:"project_id"`
	ProjectName       string  `json:"project_name"`
	User              struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"user"`
	Commit struct {
		ID          int    `json:"id"`
		SHA         string `json:"sha"`
		Message     string `json:"message"`
		AuthorName  string `json:"author_name"`
		AuthorEmail string `json:"author_email"`
		Status      string `json:"status"`
		Duration    int    `json:"duration"`
		StartedAt   string `json:"started_at"`
		FinishedAt  string `json:"finished_at"`
	} `json:"commit"`
	Repository Repository `json:"repository"`
}

// User represents a user
type User struct {
	Name      string `json:"name"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
	Email     string `json:"email"`
}

// Repository represents a repository.
type Repository struct {
	Name              string          `json:"name"`
	Description       string          `json:"description"`
	WebURL            string          `json:"web_url"`
	AvatarURL         string          `json:"avatar_url"`
	GitSSHURL         string          `json:"git_ssh_url"`
	GitHTTPURL        string          `json:"git_http_url"`
	Namespace         string          `json:"namespace"`
	Visibility        VisibilityValue `json:"visibility"`
	PathWithNamespace string          `json:"path_with_namespace"`
	DefaultBranch     string          `json:"default_branch"`
	Homepage          string          `json:"homepage"`
	URL               string          `json:"url"`
	SSHURL            string          `json:"ssh_url"`
	HTTPURL           string          `json:"http_url"`
}

// Project represents a Project
type Project struct {
	ID                int             `json:"id"`
	Name              string          `json:"name"`
	Description       string          `json:"description"`
	AvatarURL         string          `json:"avatar_url"`
	GitSSHURL         string          `json:"git_ssh_url"`
	GitHTTPURL        string          `json:"git_http_url"`
	Namespace         string          `json:"namespace"`
	PathWithNamespace string          `json:"path_with_namespace"`
	DefaultBranch     string          `json:"default_branch"`
	Homepage          string          `json:"homepage"`
	URL               string          `json:"url"`
	SSHURL            string          `json:"ssh_url"`
	HTTPURL           string          `json:"http_url"`
	WebURL            string          `json:"web_url"`
	Visibility        VisibilityValue `json:"visibility"`
	VisibilityLevel   int             `json:"visibility_level"`
	CiConfigPath      interface{}     `json:"ci_config_path"`
}

// Diff represents a GitLab diff.
//
// GitLab API docs: https://docs.gitlab.com/ce/api/commits.html
type Diff struct {
	Diff        string `json:"diff"`
	NewPath     string `json:"new_path"`
	OldPath     string `json:"old_path"`
	AMode       string `json:"a_mode"`
	BMode       string `json:"b_mode"`
	NewFile     bool   `json:"new_file"`
	RenamedFile bool   `json:"renamed_file"`
	DeletedFile bool   `json:"deleted_file"`
}

// Label represents a GitLab label.
//
// GitLab API docs: https://docs.gitlab.com/ce/api/labels.html
type Label struct {
	ID                     int    `json:"id"`
	Name                   string `json:"name"`
	Color                  string `json:"color"`
	TextColor              string `json:"text_color"`
	Description            string `json:"description"`
	OpenIssuesCount        int    `json:"open_issues_count"`
	ClosedIssuesCount      int    `json:"closed_issues_count"`
	OpenMergeRequestsCount int    `json:"open_merge_requests_count"`
	Priority               int    `json:"priority"`
	Subscribed             bool   `json:"subscribed"`
	IsProjectLabel         bool   `json:"is_project_label"`
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (l *Label) UnmarshalJSON(data []byte) error {
	type alias Label

	if err := json.Unmarshal(data, (*alias)(l)); err != nil {
		return err
	}

	if l.Name == "" {
		var raw map[string]interface{}
		if err := json.Unmarshal(data, &raw); err != nil {
			return err
		}

		if title, ok := raw["title"].(string); ok {
			l.Name = title
		}
	}

	return nil
}

// VisibilityValue represents a visibility level within GitLab.
//
// GitLab API docs: https://docs.gitlab.com/ce/api/
type VisibilityValue string

// List of available visibility levels.
//
// GitLab API docs: https://docs.gitlab.com/ce/api/
const (
	PrivateVisibility  VisibilityValue = "private"
	InternalVisibility VisibilityValue = "internal"
	PublicVisibility   VisibilityValue = "public"
)

// Snippet represents a GitLab snippet.
//
// GitLab API docs: https://docs.gitlab.com/ce/api/snippets.html
type Snippet struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	FileName    string `json:"file_name"`
	Description string `json:"description"`
	Author      struct {
		ID        int       `json:"id"`
		Username  string    `json:"username"`
		Email     string    `json:"email"`
		Name      string    `json:"name"`
		State     string    `json:"state"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"author"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
	WebURL    string    `json:"web_url"`
	RawURL    string    `json:"raw_url"`
}
