package api

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// User represents a Linear user
type User struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Email       string     `json:"email"`
	AvatarURL   string     `json:"avatarUrl"`
	DisplayName string     `json:"displayName"`
	IsMe        bool       `json:"isMe"`
	Active      bool       `json:"active"`
	Admin       bool       `json:"admin"`
	CreatedAt   *time.Time `json:"createdAt"`
}

// Team represents a Linear team
type Team struct {
	ID                 string  `json:"id"`
	Key                string  `json:"key"`
	Name               string  `json:"name"`
	Description        string  `json:"description"`
	Icon               *string `json:"icon"`
	Color              string  `json:"color"`
	Private            bool    `json:"private"`
	IssueCount         int     `json:"issueCount"`
	CyclesEnabled      bool    `json:"cyclesEnabled"`
	CycleStartDay      int     `json:"cycleStartDay"`
	CycleDuration      int     `json:"cycleDuration"`
	UpcomingCycleCount int     `json:"upcomingCycleCount"`
}

// Issue represents a Linear issue
type Issue struct {
	ID                  string            `json:"id"`
	Identifier          string            `json:"identifier"`
	Title               string            `json:"title"`
	Description         string            `json:"description"`
	Priority            int               `json:"priority"`
	Estimate            *float64          `json:"estimate"`
	CreatedAt           time.Time         `json:"createdAt"`
	UpdatedAt           time.Time         `json:"updatedAt"`
	DueDate             *string           `json:"dueDate"`
	State               *State            `json:"state"`
	Assignee            *User             `json:"assignee"`
	Team                *Team             `json:"team"`
	Labels              *Labels           `json:"labels"`
	Children            *Issues           `json:"children"`
	Parent              *Issue            `json:"parent"`
	URL                 string            `json:"url"`
	BranchName          string            `json:"branchName"`
	Cycle               *Cycle            `json:"cycle"`
	Project             *Project          `json:"project"`
	ProjectMilestone    *ProjectMilestone `json:"projectMilestone"`
	Attachments         *Attachments      `json:"attachments"`
	Comments            *Comments         `json:"comments"`
	SnoozedUntilAt      *time.Time        `json:"snoozedUntilAt"`
	CompletedAt         *time.Time        `json:"completedAt"`
	CanceledAt          *time.Time        `json:"canceledAt"`
	ArchivedAt          *time.Time        `json:"archivedAt"`
	TriagedAt           *time.Time        `json:"triagedAt"`
	CustomerTicketCount int               `json:"customerTicketCount"`
	PreviousIdentifiers []string          `json:"previousIdentifiers"`
	// Additional fields
	Number                int              `json:"number"`
	BoardOrder            float64          `json:"boardOrder"`
	SubIssueSortOrder     float64          `json:"subIssueSortOrder"`
	PriorityLabel         string           `json:"priorityLabel"`
	IntegrationSourceType *string          `json:"integrationSourceType"`
	Creator               *User            `json:"creator"`
	Subscribers           *Users           `json:"subscribers"`
	Relations             *IssueRelations  `json:"relations"`
	History               *IssueHistory    `json:"history"`
	Reactions             []Reaction       `json:"reactions"`
	SlackIssueComments    []SlackComment   `json:"slackIssueComments"`
	ExternalUserCreator   *ExternalUser    `json:"externalUserCreator"`
	CustomerTickets       []CustomerTicket `json:"customerTickets"`
}

// State represents an issue state
type State struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Color       string  `json:"color"`
	Description *string `json:"description"`
	Position    float64 `json:"position"`
}

// Project represents a Linear project
type Project struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	State       string     `json:"state"`
	Progress    float64    `json:"progress"`
	StartDate   *string    `json:"startDate"`
	TargetDate  *string    `json:"targetDate"`
	Lead        *User      `json:"lead"`
	Teams       *Teams     `json:"teams"`
	URL         string     `json:"url"`
	Icon        *string    `json:"icon"`
	Color       string     `json:"color"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	CompletedAt *time.Time `json:"completedAt"`
	CanceledAt  *time.Time `json:"canceledAt"`
	ArchivedAt  *time.Time `json:"archivedAt"`
	Creator     *User      `json:"creator"`
	Members     *Users     `json:"members"`
	Issues      *Issues    `json:"issues"`
	// Additional fields
	SlugId              string             `json:"slugId"`
	Content             string             `json:"content"`
	ConvertedFromIssue  *Issue             `json:"convertedFromIssue"`
	LastAppliedTemplate *Template          `json:"lastAppliedTemplate"`
	ProjectUpdates      *ProjectUpdates    `json:"projectUpdates"`
	Documents           *Documents         `json:"documents"`
	ProjectMilestones   *ProjectMilestones `json:"projectMilestones"`
	Health              string             `json:"health"`
	Scope               int                `json:"scope"`
	SlackNewIssue       bool               `json:"slackNewIssue"`
	SlackIssueComments  bool               `json:"slackIssueComments"`
	SlackIssueStatuses  bool               `json:"slackIssueStatuses"`
}

// Paginated collections
type Issues struct {
	Nodes    []Issue  `json:"nodes"`
	PageInfo PageInfo `json:"pageInfo"`
}

type Teams struct {
	Nodes    []Team   `json:"nodes"`
	PageInfo PageInfo `json:"pageInfo"`
}

type Projects struct {
	Nodes    []Project `json:"nodes"`
	PageInfo PageInfo  `json:"pageInfo"`
}

type Users struct {
	Nodes    []User   `json:"nodes"`
	PageInfo PageInfo `json:"pageInfo"`
}

type Labels struct {
	Nodes []Label `json:"nodes"`
}

type Label struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Color       string  `json:"color"`
	Description *string `json:"description"`
	IsGroup     bool    `json:"isGroup"`
	Parent      *Label  `json:"parent"`
	Children    *Labels `json:"children"`
}

// Cycle represents a Linear cycle (sprint)
type Cycle struct {
	ID           string     `json:"id"`
	Number       int        `json:"number"`
	Name         string     `json:"name"`
	Description  *string    `json:"description"`
	StartsAt     string     `json:"startsAt"`
	EndsAt       string     `json:"endsAt"`
	Progress     float64    `json:"progress"`
	CompletedAt  *time.Time `json:"completedAt"`
	ScopeHistory []float64  `json:"scopeHistory"`
}

// Attachment represents a file attachment or link
type Attachment struct {
	ID        string                 `json:"id"`
	Title     string                 `json:"title"`
	Subtitle  *string                `json:"subtitle"`
	URL       string                 `json:"url"`
	Metadata  map[string]interface{} `json:"metadata"`
	CreatedAt time.Time              `json:"createdAt"`
	Creator   *User                  `json:"creator"`

	// Use a map to capture any extra fields Linear might return
	Extra map[string]interface{} `json:"-"`
}

// UnmarshalJSON implements custom unmarshaling to handle unexpected fields from Linear API
func (a *Attachment) UnmarshalJSON(data []byte) error {
	// Create an alias to avoid infinite recursion
	type Alias Attachment
	aux := &struct {
		*Alias
		// Capture extra fields that might come from Linear
		Source     interface{} `json:"source,omitempty"`
		SourceType interface{} `json:"sourceType,omitempty"`
	}{
		Alias: (*Alias)(a),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Store unexpected fields in Extra map if needed
	if aux.Source != nil || aux.SourceType != nil {
		a.Extra = make(map[string]interface{})
		if aux.Source != nil {
			a.Extra["source"] = aux.Source
		}
		if aux.SourceType != nil {
			a.Extra["sourceType"] = aux.SourceType
		}
	}

	return nil
}

// Attachments represents a paginated list of attachments
type Attachments struct {
	Nodes []Attachment `json:"nodes"`
}

// Initiative represents a Linear initiative
type Initiative struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Content     string     `json:"content"`
	Status      string     `json:"status"`
	Health      *string    `json:"health"`
	TargetDate  *string    `json:"targetDate"`
	URL         string     `json:"url"`
	SlugId      string     `json:"slugId"`
	Owner       *User      `json:"owner"`
	Creator     *User      `json:"creator"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	ArchivedAt  *time.Time `json:"archivedAt"`
	Projects    *Projects  `json:"projects"`
}

// Initiatives represents a paginated list of initiatives
type Initiatives struct {
	Nodes    []Initiative `json:"nodes"`
	PageInfo PageInfo     `json:"pageInfo"`
}

// InitiativeToProject represents the join record between an initiative and a project
type InitiativeToProject struct {
	ID         string      `json:"id"`
	SortOrder  float64     `json:"sortOrder"`
	Project    *Project    `json:"project"`
	Initiative *Initiative `json:"initiative"`
}

// InitiativeToProjects represents a paginated list of initiative-project joins
type InitiativeToProjects struct {
	Nodes    []InitiativeToProject `json:"nodes"`
	PageInfo PageInfo              `json:"pageInfo"`
}

type PageInfo struct {
	HasNextPage bool   `json:"hasNextPage"`
	EndCursor   string `json:"endCursor"`
}

// Additional types for expanded fields
type IssueRelations struct {
	Nodes []IssueRelation `json:"nodes"`
}

type IssueRelation struct {
	ID           string `json:"id"`
	Type         string `json:"type"`
	Issue        *Issue `json:"issue"`
	RelatedIssue *Issue `json:"relatedIssue"`
}

type IssueHistory struct {
	Nodes []IssueHistoryEntry `json:"nodes"`
}

type IssueHistoryEntry struct {
	ID              string    `json:"id"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
	Changes         string    `json:"changes"`
	Actor           *User     `json:"actor"`
	FromAssignee    *User     `json:"fromAssignee"`
	ToAssignee      *User     `json:"toAssignee"`
	FromState       *State    `json:"fromState"`
	ToState         *State    `json:"toState"`
	FromPriority    *int      `json:"fromPriority"`
	ToPriority      *int      `json:"toPriority"`
	FromTitle       *string   `json:"fromTitle"`
	ToTitle         *string   `json:"toTitle"`
	FromCycle       *Cycle    `json:"fromCycle"`
	ToCycle         *Cycle    `json:"toCycle"`
	FromProject     *Project  `json:"fromProject"`
	ToProject       *Project  `json:"toProject"`
	AddedLabelIds   []string  `json:"addedLabelIds"`
	RemovedLabelIds []string  `json:"removedLabelIds"`
}

type Reaction struct {
	ID        string    `json:"id"`
	Emoji     string    `json:"emoji"`
	User      *User     `json:"user"`
	CreatedAt time.Time `json:"createdAt"`
}

type SlackComment struct {
	ID   string `json:"id"`
	Body string `json:"body"`
}

type ExternalUser struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CustomerTicket struct {
	ID         string    `json:"id"`
	Title      string    `json:"title"`
	CreatedAt  time.Time `json:"createdAt"`
	ExternalId string    `json:"externalId"`
}

type Template struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ProjectMilestone represents a Linear project milestone
type ProjectMilestone struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	TargetDate  *string    `json:"targetDate"`
	SortOrder   float64    `json:"sortOrder"`
	Status      string     `json:"status"`
	Progress    float64    `json:"progress"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	ArchivedAt  *time.Time `json:"archivedAt"`
}

// ProjectMilestones represents a paginated list of project milestones
type ProjectMilestones struct {
	Nodes    []ProjectMilestone `json:"nodes"`
	PageInfo PageInfo           `json:"pageInfo"`
}

type Roadmaps struct {
	Nodes []Roadmap `json:"nodes"`
}

type Roadmap struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Creator     *User     `json:"creator"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type ProjectUpdates struct {
	Nodes []ProjectUpdate `json:"nodes"`
}

type ProjectUpdate struct {
	ID        string     `json:"id"`
	Body      string     `json:"body"`
	User      *User      `json:"user"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	EditedAt  *time.Time `json:"editedAt"`
	Health    string     `json:"health"`
}

type Documents struct {
	Nodes    []Document `json:"nodes"`
	PageInfo PageInfo   `json:"pageInfo"`
}

type Document struct {
	ID         string      `json:"id"`
	Title      string      `json:"title"`
	Content    string      `json:"content"`
	Icon       *string     `json:"icon"`
	Color      string      `json:"color"`
	SlugId     string      `json:"slugId"`
	URL        string      `json:"url"`
	Creator    *User       `json:"creator"`
	UpdatedBy  *User       `json:"updatedBy"`
	Project    *Project    `json:"project"`
	Initiative *Initiative `json:"initiative"`
	CreatedAt  time.Time   `json:"createdAt"`
	UpdatedAt  time.Time   `json:"updatedAt"`
}

type ProjectLinks struct {
	Nodes []ProjectLink `json:"nodes"`
}

type ProjectLink struct {
	ID        string    `json:"id"`
	URL       string    `json:"url"`
	Label     string    `json:"label"`
	Creator   *User     `json:"creator"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// GetViewer returns the current authenticated user
func (c *Client) GetViewer(ctx context.Context) (*User, error) {
	query := `
		query Me {
			viewer {
				id
				name
				email
				avatarUrl
				isMe
				active
				admin
			}
		}
	`

	var response struct {
		Viewer User `json:"viewer"`
	}

	err := c.Execute(ctx, query, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response.Viewer, nil
}

// GetIssues returns a list of issues with optional filtering
func (c *Client) GetIssues(ctx context.Context, filter map[string]interface{}, first int, after string, orderBy string) (*Issues, error) {
	query := `
		query Issues($filter: IssueFilter, $first: Int, $after: String, $orderBy: PaginationOrderBy) {
			issues(filter: $filter, first: $first, after: $after, orderBy: $orderBy) {
				nodes {
					id
					identifier
					title
					description
					priority
					estimate
					createdAt
					updatedAt
					dueDate
					url
					state {
						id
						name
						type
						color
					}
					assignee {
						id
						name
						email
					}
					creator {
						id
						name
						email
					}
					team {
						id
						key
						name
					}
					labels {
						nodes {
							id
							name
							color
						}
					}
				}
				pageInfo {
					hasNextPage
					endCursor
				}
			}
		}
	`

	variables := map[string]interface{}{
		"first": first,
	}
	if filter != nil {
		variables["filter"] = filter
	}
	if after != "" {
		variables["after"] = after
	}
	if orderBy != "" {
		variables["orderBy"] = orderBy
	}

	var response struct {
		Issues Issues `json:"issues"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.Issues, nil
}

// IssueSearch returns issues that match a full-text query
func (c *Client) IssueSearch(ctx context.Context, term string, filter map[string]interface{}, first int, after string, orderBy string, includeArchived bool) (*Issues, error) {
	query := `
		query IssueSearch($term: String!, $filter: IssueFilter, $first: Int, $after: String, $orderBy: PaginationOrderBy, $includeArchived: Boolean) {
			searchIssues(term: $term, filter: $filter, first: $first, after: $after, orderBy: $orderBy, includeArchived: $includeArchived) {
				nodes {
					id
					identifier
					title
					description
					priority
					estimate
					createdAt
					updatedAt
					dueDate
					url
					state {
						id
						name
						type
						color
					}
					assignee {
						id
						name
						email
					}
					creator {
						id
						name
						email
					}
					team {
						id
						key
						name
					}
					labels {
						nodes {
							id
							name
							color
						}
					}
				}
				pageInfo {
					hasNextPage
					endCursor
				}
			}
		}
	`

	variables := map[string]interface{}{
		"term":            term,
		"first":           first,
		"includeArchived": includeArchived,
	}
	if filter != nil {
		variables["filter"] = filter
	}
	if after != "" {
		variables["after"] = after
	}
	if orderBy != "" {
		variables["orderBy"] = orderBy
	}

	var response struct {
		SearchIssues struct {
			Nodes    []Issue  `json:"nodes"`
			PageInfo PageInfo `json:"pageInfo"`
		} `json:"searchIssues"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &Issues{
		Nodes:    response.SearchIssues.Nodes,
		PageInfo: response.SearchIssues.PageInfo,
	}, nil
}

// GetIssue returns a single issue by ID
func (c *Client) GetIssue(ctx context.Context, id string) (*Issue, error) {
	query := `
		query Issue($id: String!) {
			issue(id: $id) {
				id
				identifier
				number
				title
				description
				priority
				priorityLabel
				estimate
				boardOrder
				subIssueSortOrder
				createdAt
				updatedAt
				dueDate
				url
				branchName
				snoozedUntilAt
				completedAt
				canceledAt
				archivedAt
				triagedAt
				customerTicketCount
				previousIdentifiers
				integrationSourceType
				state {
					id
					name
					type
					color
					description
					position
				}
				assignee {
					id
					name
					email
					avatarUrl
					displayName
					active
					admin
					createdAt
				}
				creator {
					id
					name
					email
					avatarUrl
					displayName
					active
				}
				team {
					id
					key
					name
					description
					icon
					color
					cyclesEnabled
					cycleStartDay
					cycleDuration
					upcomingCycleCount
				}
				labels {
					nodes {
						id
						name
						color
						description
						parent {
							id
							name
						}
					}
				}
				parent {
					id
					identifier
					title
					state {
						name
						type
					}
				}
				children {
					nodes {
						id
						identifier
						title
						priority
						createdAt
						state {
							name
							type
							color
						}
						assignee {
							name
							email
						}
					}
				}
				cycle {
					id
					number
					name
					description
					startsAt
					endsAt
					progress
					completedAt
					scopeHistory
				}
				project {
					id
					name
					description
					state
					progress
					startDate
					targetDate
					health
					lead {
						name
						email
					}
				}
				projectMilestone {
					id
					name
					description
					targetDate
					status
					progress
				}
				attachments(first: 20) {
					nodes {
						id
						title
						subtitle
						url
						metadata
						createdAt
						creator {
							name
							email
						}
					}
				}
				comments(first: 10) {
					nodes {
						id
						body
						createdAt
						updatedAt
						editedAt
						user {
							name
							email
							avatarUrl
						}
						parent {
							id
						}
						children {
							nodes {
								id
								body
								user {
									name
								}
							}
						}
					}
				}
				subscribers {
					nodes {
						id
						name
						email
						avatarUrl
					}
				}
				relations {
					nodes {
						id
						type
						relatedIssue {
							id
							identifier
							title
							state {
								name
								type
							}
						}
					}
				}
				history(first: 10) {
					nodes {
						id
						createdAt
						updatedAt
						actor {
							name
							email
						}
						fromAssignee {
							name
						}
						toAssignee {
							name
						}
						fromState {
							name
						}
						toState {
							name
						}
						fromPriority
						toPriority
						fromTitle
						toTitle
						fromCycle {
							name
						}
						toCycle {
							name
						}
						fromProject {
							name
						}
						toProject {
							name
						}
						addedLabelIds
						removedLabelIds
					}
				}
				reactions {
					id
					emoji
					user {
						name
						email
					}
					createdAt
				}
				externalUserCreator {
					id
					name
					email
					avatarUrl
				}
			}
		}
	`

	variables := map[string]interface{}{
		"id": id,
	}

	var response struct {
		Issue Issue `json:"issue"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.Issue, nil
}

// GetTeams returns a list of teams
func (c *Client) GetTeams(ctx context.Context, first int, after string, orderBy string) (*Teams, error) {
	query := `
		query Teams($first: Int, $after: String, $orderBy: PaginationOrderBy) {
			teams(first: $first, after: $after, orderBy: $orderBy) {
				nodes {
					id
					key
					name
					description
					private
					issueCount
				}
				pageInfo {
					hasNextPage
					endCursor
				}
			}
		}
	`

	variables := map[string]interface{}{
		"first": first,
	}
	if after != "" {
		variables["after"] = after
	}
	if orderBy != "" {
		variables["orderBy"] = orderBy
	}

	var response struct {
		Teams Teams `json:"teams"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.Teams, nil
}

// GetProjects returns a list of projects
func (c *Client) GetProjects(ctx context.Context, filter map[string]interface{}, first int, after string, orderBy string) (*Projects, error) {
	query := `
		query Projects($filter: ProjectFilter, $first: Int, $after: String, $orderBy: PaginationOrderBy) {
			projects(filter: $filter, first: $first, after: $after, orderBy: $orderBy) {
				nodes {
					id
					name
					description
					state
					progress
					startDate
					targetDate
					url
					createdAt
					updatedAt
					lead {
						id
						name
						email
					}
					teams {
						nodes {
							id
							key
							name
						}
					}
				}
				pageInfo {
					hasNextPage
					endCursor
				}
			}
		}
	`

	variables := map[string]interface{}{
		"first": first,
	}
	if filter != nil {
		variables["filter"] = filter
	}
	if after != "" {
		variables["after"] = after
	}
	if orderBy != "" {
		variables["orderBy"] = orderBy
	}

	var response struct {
		Projects Projects `json:"projects"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.Projects, nil
}

// GetProject returns a single project by ID
func (c *Client) GetProject(ctx context.Context, id string) (*Project, error) {
	query := `
		query Project($id: String!) {
			project(id: $id) {
				id
				slugId
				name
				description
				content
				state
				progress
				health
				scope
				startDate
				targetDate
				url
				icon
				color
				createdAt
				updatedAt
				completedAt
				canceledAt
				archivedAt
				slackNewIssue
				slackIssueComments
				slackIssueStatuses
				lead {
					id
					name
					email
					avatarUrl
					displayName
					active
				}
				creator {
					id
					name
					email
					avatarUrl
					active
				}
				convertedFromIssue {
					id
					identifier
					title
				}
				lastAppliedTemplate {
					id
					name
					description
				}
				teams {
					nodes {
						id
						key
						name
						description
						icon
						color
						cyclesEnabled
					}
				}
				members {
					nodes {
						id
						name
						email
						avatarUrl
						displayName
						active
						admin
					}
				}
				issues(first: 50, orderBy: updatedAt) {
					nodes {
						id
						identifier
						number
						title
						description
						priority
						estimate
						createdAt
						updatedAt
						completedAt
						state {
							name
							type
							color
						}
						assignee {
							name
							email
						}
						labels {
							nodes {
								name
								color
							}
						}
					}
				}
				projectUpdates(first: 10) {
					nodes {
						id
						body
						health
						createdAt
						updatedAt
						editedAt
						user {
							name
							email
							avatarUrl
						}
					}
				}
				documents(first: 20) {
					nodes {
						id
						title
						content
						icon
						color
						createdAt
						updatedAt
						creator {
							name
							email
						}
						updatedBy {
							name
							email
						}
					}
				}
				projectMilestones(first: 250) {
					nodes {
						id
						name
						description
						targetDate
						sortOrder
						status
						progress
						createdAt
						updatedAt
						archivedAt
					}
					pageInfo {
						hasNextPage
						endCursor
					}
				}
			}
		}
	`

	variables := map[string]interface{}{
		"id": id,
	}

	var response struct {
		Project Project `json:"project"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.Project, nil
}

// GetProjectMilestones returns milestones for a project
func (c *Client) GetProjectMilestones(ctx context.Context, projectID string, first int, after string) (*ProjectMilestones, error) {
	query := `
		query ProjectMilestones($id: String!, $first: Int, $after: String) {
			project(id: $id) {
				projectMilestones(first: $first, after: $after) {
					nodes {
						id
						name
						description
						targetDate
						sortOrder
						status
						progress
						createdAt
						updatedAt
						archivedAt
					}
					pageInfo {
						hasNextPage
						endCursor
					}
				}
			}
		}
	`

	variables := map[string]interface{}{
		"id":    projectID,
		"first": first,
	}
	if after != "" {
		variables["after"] = after
	}

	var response struct {
		Project struct {
			ProjectMilestones ProjectMilestones `json:"projectMilestones"`
		} `json:"project"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.Project.ProjectMilestones, nil
}

// UpdateIssue updates an issue's fields
func (c *Client) UpdateIssue(ctx context.Context, id string, input map[string]interface{}) (*Issue, error) {
	query := `
		mutation UpdateIssue($id: String!, $input: IssueUpdateInput!) {
			issueUpdate(id: $id, input: $input) {
				issue {
					id
					identifier
					title
					description
					priority
					estimate
					createdAt
					updatedAt
					dueDate
					state {
						id
						name
						type
						color
					}
					assignee {
						id
						name
						email
					}
					team {
						id
						key
						name
					}
					labels {
						nodes {
							id
							name
							color
						}
					}
					projectMilestone {
						id
						name
					}
				}
			}
		}
	`

	variables := map[string]interface{}{
		"id":    id,
		"input": input,
	}

	var response struct {
		IssueUpdate struct {
			Issue Issue `json:"issue"`
		} `json:"issueUpdate"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.IssueUpdate.Issue, nil
}

// CreateIssue creates a new issue
func (c *Client) CreateIssue(ctx context.Context, input map[string]interface{}) (*Issue, error) {
	query := `
		mutation CreateIssue($input: IssueCreateInput!) {
			issueCreate(input: $input) {
				issue {
					id
					identifier
					title
					description
					priority
					estimate
					createdAt
					updatedAt
					dueDate
					state {
						id
						name
						type
						color
					}
					assignee {
						id
						name
						email
					}
					team {
						id
						key
						name
					}
					labels {
						nodes {
							id
							name
							color
						}
					}
				}
			}
		}
	`

	variables := map[string]interface{}{
		"input": input,
	}

	var response struct {
		IssueCreate struct {
			Issue Issue `json:"issue"`
		} `json:"issueCreate"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.IssueCreate.Issue, nil
}

// GetTeam returns a single team by key
func (c *Client) GetTeam(ctx context.Context, key string) (*Team, error) {
	query := `
		query Team($key: String!) {
			team(id: $key) {
				id
				key
				name
				description
				private
				issueCount
			}
		}
	`

	variables := map[string]interface{}{
		"key": key,
	}

	var response struct {
		Team Team `json:"team"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.Team, nil
}

// Comment represents a Linear comment
type Comment struct {
	ID        string     `json:"id"`
	Body      string     `json:"body"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	EditedAt  *time.Time `json:"editedAt"`
	User      *User      `json:"user"`
	Parent    *Comment   `json:"parent"`
	Children  *Comments  `json:"children"`
}

// Comments represents a paginated list of comments
type Comments struct {
	Nodes    []Comment `json:"nodes"`
	PageInfo PageInfo  `json:"pageInfo"`
}

// WorkflowState represents a Linear workflow state
type WorkflowState struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Color       string  `json:"color"`
	Description string  `json:"description"`
	Position    float64 `json:"position"`
}

// GetTeamStates returns workflow states for a team
func (c *Client) GetTeamStates(ctx context.Context, teamKey string) ([]WorkflowState, error) {
	query := `
		query TeamStates($key: String!) {
			team(id: $key) {
				states {
					nodes {
						id
						name
						type
						color
						description
						position
					}
				}
			}
		}
	`

	variables := map[string]interface{}{
		"key": teamKey,
	}

	var response struct {
		Team struct {
			States struct {
				Nodes []WorkflowState `json:"nodes"`
			} `json:"states"`
		} `json:"team"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return response.Team.States.Nodes, nil
}

// GetTeamMembers returns members of a specific team
func (c *Client) GetTeamMembers(ctx context.Context, teamKey string) (*Users, error) {
	query := `
		query TeamMembers($key: String!) {
			team(id: $key) {
				members {
					nodes {
						id
						name
						email
						avatarUrl
						isMe
						active
						admin
					}
					pageInfo {
						hasNextPage
						endCursor
					}
				}
			}
		}
	`

	variables := map[string]interface{}{
		"key": teamKey,
	}

	var response struct {
		Team struct {
			Members Users `json:"members"`
		} `json:"team"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.Team.Members, nil
}

// GetUsers returns a list of all users
func (c *Client) GetUsers(ctx context.Context, first int, after string, orderBy string) (*Users, error) {
	query := `
		query Users($first: Int, $after: String, $orderBy: PaginationOrderBy) {
			users(first: $first, after: $after, orderBy: $orderBy) {
				nodes {
					id
					name
					email
					avatarUrl
					isMe
					active
					admin
				}
				pageInfo {
					hasNextPage
					endCursor
				}
			}
		}
	`

	variables := map[string]interface{}{
		"first": first,
	}
	if after != "" {
		variables["after"] = after
	}
	if orderBy != "" {
		variables["orderBy"] = orderBy
	}

	var response struct {
		Users Users `json:"users"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.Users, nil
}

// GetUser returns a specific user by email
func (c *Client) GetUser(ctx context.Context, email string) (*User, error) {
	query := `
		query User($email: String!) {
			user(email: $email) {
				id
				name
				email
				avatarUrl
				isMe
				active
				admin
			}
		}
	`

	variables := map[string]interface{}{
		"email": email,
	}

	var response struct {
		User User `json:"user"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.User, nil
}

// GetIssueComments returns comments for a specific issue
func (c *Client) GetIssueComments(ctx context.Context, issueID string, first int, after string, orderBy string) (*Comments, error) {
	query := `
		query IssueComments($id: String!, $first: Int, $after: String, $orderBy: PaginationOrderBy) {
			issue(id: $id) {
				comments(first: $first, after: $after, orderBy: $orderBy) {
					nodes {
						id
						body
						createdAt
						updatedAt
						user {
							id
							name
							email
						}
					}
					pageInfo {
						hasNextPage
						endCursor
					}
				}
			}
		}
	`

	variables := map[string]interface{}{
		"id":    issueID,
		"first": first,
	}
	if after != "" {
		variables["after"] = after
	}
	if orderBy != "" {
		variables["orderBy"] = orderBy
	}

	var response struct {
		Issue struct {
			Comments Comments `json:"comments"`
		} `json:"issue"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.Issue.Comments, nil
}

// CreateComment creates a new comment on an issue
func (c *Client) CreateComment(ctx context.Context, issueID string, body string) (*Comment, error) {
	query := `
		mutation CreateComment($input: CommentCreateInput!) {
			commentCreate(input: $input) {
				comment {
					id
					body
					createdAt
					updatedAt
					user {
						id
						name
						email
					}
				}
			}
		}
	`

	input := map[string]interface{}{
		"issueId": issueID,
		"body":    body,
	}

	variables := map[string]interface{}{
		"input": input,
	}

	var response struct {
		CommentCreate struct {
			Comment Comment `json:"comment"`
		} `json:"commentCreate"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.CommentCreate.Comment, nil
}

// GetLabels returns all labels in the workspace, optionally filtered by team.
// Paginates automatically to collect all results.
func (c *Client) GetLabels(ctx context.Context, teamKey string) ([]Label, error) {
	query := `
		query IssueLabels($first: Int, $after: String, $filter: IssueLabelFilter) {
			issueLabels(first: $first, after: $after, filter: $filter) {
				nodes {
					id
					name
					color
					description
					isGroup
					parent {
						id
						name
					}
					children {
						nodes {
							id
							name
							color
							description
							isGroup
						}
					}
				}
				pageInfo {
					hasNextPage
					endCursor
				}
			}
		}
	`

	var allLabels []Label
	var after string

	for {
		variables := map[string]interface{}{
			"first": 100,
		}
		if after != "" {
			variables["after"] = after
		}

		if teamKey != "" {
			variables["filter"] = map[string]interface{}{
				"team": map[string]interface{}{
					"key": map[string]interface{}{"eq": teamKey},
				},
			}
		}

		var response struct {
			IssueLabels struct {
				Nodes    []Label  `json:"nodes"`
				PageInfo PageInfo `json:"pageInfo"`
			} `json:"issueLabels"`
		}

		err := c.Execute(ctx, query, variables, &response)
		if err != nil {
			return nil, err
		}

		allLabels = append(allLabels, response.IssueLabels.Nodes...)

		if !response.IssueLabels.PageInfo.HasNextPage {
			break
		}
		after = response.IssueLabels.PageInfo.EndCursor
	}

	return allLabels, nil
}

// ResolveLabelNames resolves label names to IDs, returning an error if any are not found
func (c *Client) ResolveLabelNames(ctx context.Context, names []string) ([]string, error) {
	labels, err := c.GetLabels(ctx, "")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch labels: %w", err)
	}

	// Build a case-insensitive lookup map
	labelMap := make(map[string]string)
	for _, label := range labels {
		labelMap[strings.ToLower(label.Name)] = label.ID
	}

	var ids []string
	var missing []string
	for _, name := range names {
		id, ok := labelMap[strings.ToLower(name)]
		if !ok {
			missing = append(missing, name)
		} else {
			ids = append(ids, id)
		}
	}

	if len(missing) > 0 {
		// List available labels for a helpful error
		var available []string
		for _, label := range labels {
			available = append(available, label.Name)
		}
		return nil, fmt.Errorf("labels not found: %s. Available: %s",
			strings.Join(missing, ", "), strings.Join(available, ", "))
	}

	return ids, nil
}

// CreateLabel creates a new workspace label. If parentID is set, the label is
// created as a child of that group. Set isGroup=true to create a label group.
func (c *Client) CreateLabel(ctx context.Context, name, color, description, parentID string, isGroup bool) (*Label, error) {
	query := `
		mutation CreateLabel($input: IssueLabelCreateInput!) {
			issueLabelCreate(input: $input) {
				success
				issueLabel {
					id
					name
					color
					description
					isGroup
					parent {
						id
						name
					}
				}
			}
		}
	`

	input := map[string]interface{}{
		"name": name,
	}
	if color != "" {
		input["color"] = color
	}
	if description != "" {
		input["description"] = description
	}
	if parentID != "" {
		input["parentId"] = parentID
	}
	if isGroup {
		input["isGroup"] = true
	}

	variables := map[string]interface{}{
		"input": input,
	}

	var response struct {
		IssueLabelCreate struct {
			Success    bool  `json:"success"`
			IssueLabel Label `json:"issueLabel"`
		} `json:"issueLabelCreate"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	if !response.IssueLabelCreate.Success {
		return nil, fmt.Errorf("failed to create label")
	}

	return &response.IssueLabelCreate.IssueLabel, nil
}

// UpdateLabel updates an existing label's fields.
func (c *Client) UpdateLabel(ctx context.Context, id string, input map[string]interface{}) (*Label, error) {
	query := `
		mutation UpdateLabel($id: String!, $input: IssueLabelUpdateInput!) {
			issueLabelUpdate(id: $id, input: $input) {
				success
				issueLabel {
					id
					name
					color
					description
					isGroup
					parent {
						id
						name
					}
				}
			}
		}
	`

	variables := map[string]interface{}{
		"id":    id,
		"input": input,
	}

	var response struct {
		IssueLabelUpdate struct {
			Success    bool  `json:"success"`
			IssueLabel Label `json:"issueLabel"`
		} `json:"issueLabelUpdate"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	if !response.IssueLabelUpdate.Success {
		return nil, fmt.Errorf("failed to update label")
	}

	return &response.IssueLabelUpdate.IssueLabel, nil
}

// DeleteLabel deletes a label by ID.
func (c *Client) DeleteLabel(ctx context.Context, id string) error {
	query := `
		mutation DeleteLabel($id: String!) {
			issueLabelDelete(id: $id) {
				success
			}
		}
	`

	variables := map[string]interface{}{
		"id": id,
	}

	var response struct {
		IssueLabelDelete struct {
			Success bool `json:"success"`
		} `json:"issueLabelDelete"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return err
	}

	if !response.IssueLabelDelete.Success {
		return fmt.Errorf("failed to delete label")
	}

	return nil
}

// CreateIssueRelation creates a relation between two issues
func (c *Client) CreateIssueRelation(ctx context.Context, issueId, relatedIssueId, relationType string) (*IssueRelation, error) {
	query := `
		mutation CreateIssueRelation($input: IssueRelationCreateInput!) {
			issueRelationCreate(input: $input) {
				issueRelation {
					id
					type
					issue {
						id
						identifier
						title
					}
					relatedIssue {
						id
						identifier
						title
					}
				}
			}
		}
	`

	input := map[string]interface{}{
		"issueId":        issueId,
		"relatedIssueId": relatedIssueId,
		"type":           relationType,
	}

	variables := map[string]interface{}{
		"input": input,
	}

	var response struct {
		IssueRelationCreate struct {
			IssueRelation IssueRelation `json:"issueRelation"`
		} `json:"issueRelationCreate"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.IssueRelationCreate.IssueRelation, nil
}

// DeleteIssueRelation deletes an issue relation by ID
func (c *Client) DeleteIssueRelation(ctx context.Context, relationId string) error {
	query := `
		mutation DeleteIssueRelation($id: String!) {
			issueRelationDelete(id: $id) {
				success
			}
		}
	`

	variables := map[string]interface{}{
		"id": relationId,
	}

	var response struct {
		IssueRelationDelete struct {
			Success bool `json:"success"`
		} `json:"issueRelationDelete"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return err
	}

	if !response.IssueRelationDelete.Success {
		return fmt.Errorf("failed to delete relation")
	}

	return nil
}

// documentFields is the shared GraphQL selection for a Document node.
const documentFields = `
	id
	title
	content
	icon
	color
	slugId
	url
	createdAt
	updatedAt
	creator { id name email }
	updatedBy { id name email }
	project { id name url }
	initiative { id name url }
`

// GetDocuments returns a list of documents with optional filtering.
func (c *Client) GetDocuments(ctx context.Context, filter map[string]interface{}, first int, after string, orderBy string) (*Documents, error) {
	query := fmt.Sprintf(`
		query Documents($filter: DocumentFilter, $first: Int, $after: String, $orderBy: PaginationOrderBy) {
			documents(filter: $filter, first: $first, after: $after, orderBy: $orderBy) {
				nodes {
					%s
				}
				pageInfo {
					hasNextPage
					endCursor
				}
			}
		}
	`, documentFields)

	variables := map[string]interface{}{
		"first": first,
	}
	if filter != nil {
		variables["filter"] = filter
	}
	if after != "" {
		variables["after"] = after
	}
	if orderBy != "" {
		variables["orderBy"] = orderBy
	}

	var response struct {
		Documents Documents `json:"documents"`
	}

	if err := c.Execute(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return &response.Documents, nil
}

// GetDocument returns a single document by ID.
func (c *Client) GetDocument(ctx context.Context, id string) (*Document, error) {
	query := fmt.Sprintf(`
		query Document($id: String!) {
			document(id: $id) {
				%s
			}
		}
	`, documentFields)

	variables := map[string]interface{}{
		"id": id,
	}

	var response struct {
		Document Document `json:"document"`
	}

	if err := c.Execute(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return &response.Document, nil
}

// CreateDocument creates a new document.
func (c *Client) CreateDocument(ctx context.Context, input map[string]interface{}) (*Document, error) {
	query := fmt.Sprintf(`
		mutation CreateDocument($input: DocumentCreateInput!) {
			documentCreate(input: $input) {
				success
				document {
					%s
				}
			}
		}
	`, documentFields)

	variables := map[string]interface{}{
		"input": input,
	}

	var response struct {
		DocumentCreate struct {
			Success  bool     `json:"success"`
			Document Document `json:"document"`
		} `json:"documentCreate"`
	}

	if err := c.Execute(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	if !response.DocumentCreate.Success {
		return nil, fmt.Errorf("failed to create document")
	}

	return &response.DocumentCreate.Document, nil
}

// UpdateDocument updates an existing document.
func (c *Client) UpdateDocument(ctx context.Context, id string, input map[string]interface{}) (*Document, error) {
	query := fmt.Sprintf(`
		mutation UpdateDocument($id: String!, $input: DocumentUpdateInput!) {
			documentUpdate(id: $id, input: $input) {
				success
				document {
					%s
				}
			}
		}
	`, documentFields)

	variables := map[string]interface{}{
		"id":    id,
		"input": input,
	}

	var response struct {
		DocumentUpdate struct {
			Success  bool     `json:"success"`
			Document Document `json:"document"`
		} `json:"documentUpdate"`
	}

	if err := c.Execute(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	if !response.DocumentUpdate.Success {
		return nil, fmt.Errorf("failed to update document")
	}

	return &response.DocumentUpdate.Document, nil
}

// DeleteDocument deletes (moves to trash) a document by ID.
func (c *Client) DeleteDocument(ctx context.Context, id string) error {
	query := `
		mutation DeleteDocument($id: String!) {
			documentDelete(id: $id) {
				success
			}
		}
	`

	variables := map[string]interface{}{
		"id": id,
	}

	var response struct {
		DocumentDelete struct {
			Success bool `json:"success"`
		} `json:"documentDelete"`
	}

	if err := c.Execute(ctx, query, variables, &response); err != nil {
		return err
	}

	if !response.DocumentDelete.Success {
		return fmt.Errorf("failed to delete document")
	}

	return nil
}

// projectWriteFields is a lean Project selection returned by write mutations.
const projectWriteFields = `
	id
	name
	description
	state
	url
	startDate
	targetDate
	lead { id name email }
`

// CreateProject creates a new project.
func (c *Client) CreateProject(ctx context.Context, input map[string]interface{}) (*Project, error) {
	query := fmt.Sprintf(`
		mutation CreateProject($input: ProjectCreateInput!) {
			projectCreate(input: $input) {
				success
				project {
					%s
				}
			}
		}
	`, projectWriteFields)

	variables := map[string]interface{}{
		"input": input,
	}

	var response struct {
		ProjectCreate struct {
			Success bool    `json:"success"`
			Project Project `json:"project"`
		} `json:"projectCreate"`
	}

	if err := c.Execute(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	if !response.ProjectCreate.Success {
		return nil, fmt.Errorf("failed to create project")
	}

	return &response.ProjectCreate.Project, nil
}

// UpdateProject updates an existing project.
func (c *Client) UpdateProject(ctx context.Context, id string, input map[string]interface{}) (*Project, error) {
	query := fmt.Sprintf(`
		mutation UpdateProject($id: String!, $input: ProjectUpdateInput!) {
			projectUpdate(id: $id, input: $input) {
				success
				project {
					%s
				}
			}
		}
	`, projectWriteFields)

	variables := map[string]interface{}{
		"id":    id,
		"input": input,
	}

	var response struct {
		ProjectUpdate struct {
			Success bool    `json:"success"`
			Project Project `json:"project"`
		} `json:"projectUpdate"`
	}

	if err := c.Execute(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	if !response.ProjectUpdate.Success {
		return nil, fmt.Errorf("failed to update project")
	}

	return &response.ProjectUpdate.Project, nil
}

// CreateProjectUpdate posts a project status update (health + body).
func (c *Client) CreateProjectUpdate(ctx context.Context, input map[string]interface{}) (*ProjectUpdate, error) {
	query := `
		mutation CreateProjectUpdate($input: ProjectUpdateCreateInput!) {
			projectUpdateCreate(input: $input) {
				success
				projectUpdate {
					id
					body
					health
					createdAt
					user { id name email }
				}
			}
		}
	`

	variables := map[string]interface{}{
		"input": input,
	}

	var response struct {
		ProjectUpdateCreate struct {
			Success       bool          `json:"success"`
			ProjectUpdate ProjectUpdate `json:"projectUpdate"`
		} `json:"projectUpdateCreate"`
	}

	if err := c.Execute(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	if !response.ProjectUpdateCreate.Success {
		return nil, fmt.Errorf("failed to create project update")
	}

	return &response.ProjectUpdateCreate.ProjectUpdate, nil
}

// initiativeFields is the shared GraphQL selection for an Initiative node.
const initiativeFields = `
	id
	name
	description
	content
	status
	health
	targetDate
	url
	slugId
	createdAt
	updatedAt
	archivedAt
	owner { id name email }
	creator { id name email }
`

// GetInitiatives returns a list of initiatives with optional filtering.
func (c *Client) GetInitiatives(ctx context.Context, filter map[string]interface{}, first int, after string, orderBy string) (*Initiatives, error) {
	query := fmt.Sprintf(`
		query Initiatives($filter: InitiativeFilter, $first: Int, $after: String, $orderBy: PaginationOrderBy) {
			initiatives(filter: $filter, first: $first, after: $after, orderBy: $orderBy) {
				nodes {
					%s
				}
				pageInfo {
					hasNextPage
					endCursor
				}
			}
		}
	`, initiativeFields)

	variables := map[string]interface{}{
		"first": first,
	}
	if filter != nil {
		variables["filter"] = filter
	}
	if after != "" {
		variables["after"] = after
	}
	if orderBy != "" {
		variables["orderBy"] = orderBy
	}

	var response struct {
		Initiatives Initiatives `json:"initiatives"`
	}

	if err := c.Execute(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return &response.Initiatives, nil
}

// GetInitiative returns a single initiative by ID, including its projects.
func (c *Client) GetInitiative(ctx context.Context, id string) (*Initiative, error) {
	query := fmt.Sprintf(`
		query Initiative($id: String!) {
			initiative(id: $id) {
				%s
				projects(first: 100) {
					nodes {
						id
						name
						state
						url
					}
				}
			}
		}
	`, initiativeFields)

	variables := map[string]interface{}{
		"id": id,
	}

	var response struct {
		Initiative Initiative `json:"initiative"`
	}

	if err := c.Execute(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return &response.Initiative, nil
}

// CreateInitiative creates a new initiative.
func (c *Client) CreateInitiative(ctx context.Context, input map[string]interface{}) (*Initiative, error) {
	query := fmt.Sprintf(`
		mutation CreateInitiative($input: InitiativeCreateInput!) {
			initiativeCreate(input: $input) {
				success
				initiative {
					%s
				}
			}
		}
	`, initiativeFields)

	variables := map[string]interface{}{
		"input": input,
	}

	var response struct {
		InitiativeCreate struct {
			Success    bool       `json:"success"`
			Initiative Initiative `json:"initiative"`
		} `json:"initiativeCreate"`
	}

	if err := c.Execute(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	if !response.InitiativeCreate.Success {
		return nil, fmt.Errorf("failed to create initiative")
	}

	return &response.InitiativeCreate.Initiative, nil
}

// UpdateInitiative updates an existing initiative.
func (c *Client) UpdateInitiative(ctx context.Context, id string, input map[string]interface{}) (*Initiative, error) {
	query := fmt.Sprintf(`
		mutation UpdateInitiative($id: String!, $input: InitiativeUpdateInput!) {
			initiativeUpdate(id: $id, input: $input) {
				success
				initiative {
					%s
				}
			}
		}
	`, initiativeFields)

	variables := map[string]interface{}{
		"id":    id,
		"input": input,
	}

	var response struct {
		InitiativeUpdate struct {
			Success    bool       `json:"success"`
			Initiative Initiative `json:"initiative"`
		} `json:"initiativeUpdate"`
	}

	if err := c.Execute(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	if !response.InitiativeUpdate.Success {
		return nil, fmt.Errorf("failed to update initiative")
	}

	return &response.InitiativeUpdate.Initiative, nil
}

// AddProjectToInitiative attaches a project to an initiative.
func (c *Client) AddProjectToInitiative(ctx context.Context, initiativeID, projectID string) error {
	query := `
		mutation AddProjectToInitiative($input: InitiativeToProjectCreateInput!) {
			initiativeToProjectCreate(input: $input) {
				success
			}
		}
	`

	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"initiativeId": initiativeID,
			"projectId":    projectID,
		},
	}

	var response struct {
		InitiativeToProjectCreate struct {
			Success bool `json:"success"`
		} `json:"initiativeToProjectCreate"`
	}

	if err := c.Execute(ctx, query, variables, &response); err != nil {
		return err
	}

	if !response.InitiativeToProjectCreate.Success {
		return fmt.Errorf("failed to attach project to initiative")
	}

	return nil
}

// findInitiativeToProjectID resolves the join-record ID for an
// (initiative, project) pair by paging the unfilterable initiativeToProjects
// query. Returns an empty string if no such join exists.
func (c *Client) findInitiativeToProjectID(ctx context.Context, initiativeID, projectID string) (string, error) {
	query := `
		query InitiativeToProjects($first: Int, $after: String) {
			initiativeToProjects(first: $first, after: $after) {
				nodes {
					id
					initiative { id }
					project { id }
				}
				pageInfo {
					hasNextPage
					endCursor
				}
			}
		}
	`

	after := ""
	for {
		variables := map[string]interface{}{"first": 250}
		if after != "" {
			variables["after"] = after
		}

		var response struct {
			InitiativeToProjects InitiativeToProjects `json:"initiativeToProjects"`
		}
		if err := c.Execute(ctx, query, variables, &response); err != nil {
			return "", err
		}

		for _, node := range response.InitiativeToProjects.Nodes {
			if node.Initiative != nil && node.Project != nil &&
				node.Initiative.ID == initiativeID && node.Project.ID == projectID {
				return node.ID, nil
			}
		}

		if !response.InitiativeToProjects.PageInfo.HasNextPage {
			return "", nil
		}
		after = response.InitiativeToProjects.PageInfo.EndCursor
	}
}

// RemoveProjectFromInitiative detaches a project from an initiative by
// resolving and deleting the join record.
func (c *Client) RemoveProjectFromInitiative(ctx context.Context, initiativeID, projectID string) error {
	joinID, err := c.findInitiativeToProjectID(ctx, initiativeID, projectID)
	if err != nil {
		return err
	}
	if joinID == "" {
		return fmt.Errorf("project %s is not attached to initiative %s", projectID, initiativeID)
	}

	query := `
		mutation RemoveProjectFromInitiative($id: String!) {
			initiativeToProjectDelete(id: $id) {
				success
			}
		}
	`

	variables := map[string]interface{}{"id": joinID}

	var response struct {
		InitiativeToProjectDelete struct {
			Success bool `json:"success"`
		} `json:"initiativeToProjectDelete"`
	}

	if err := c.Execute(ctx, query, variables, &response); err != nil {
		return err
	}

	if !response.InitiativeToProjectDelete.Success {
		return fmt.Errorf("failed to detach project from initiative")
	}

	return nil
}
