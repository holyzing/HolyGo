package api

import (
	"context"
)

type LukeAPI interface {
	CreateJob(ctx context.Context, input *CreateJobInput) (*CreateJobOutput, error)
	GetJob(ctx context.Context, input *GetJobInput) (*GetJobOutput, error)
}

type CreateJobInput struct {
	User         string `json:"user"`
	TenantName   string `json:"tenant_name"`
	Organization string `json:"organization"`

	Sync       bool              `json:"sync,omitempty"` // add
	Cores      int32             `json:"cores,omitempty"`
	SysPrio    float32           `json:"sys_prio,omitempty"`
	Entrypoint []string          `json:"entrypoint"` // add
	Options    map[string]string `json:"options,omitempty"`
}

type JobInfoResponse struct {
	JobHandle string `json:"job_id,omitempty"`
	JobID     int64  `json:"job_handle,omitempty"`
	Info      string `json:"info,omitempty"`
}

type CreateJobOutput struct {
	RequestID string           `json:"requestID,omitempty"`
	Body      *JobInfoResponse `json:"jobInfo,omitempty"`
}

// ------------------------------------------------------------------------

type GetJobInput struct {
	User         string `json:"user"`
	TenantName   string `json:"tenant_name"`
	Organization string `json:"organization"`

	ID       int64    `json:"id"`
	Handle   string   `json:"handle"`
	Fields   []string `json:"fields"`
	Combined bool     `json:"combined"`
}

type JobDetailResponse struct {
	JobDetails []*CreateJobInput `protobuf:"bytes,1,rep,name=job_details,json=jobDetails,proto3" json:"job_details,omitempty"`
}

type GetJobOutput struct {
	RequestID string             `json:"requestID,omitempty"`
	Body      *JobDetailResponse `json:"jobDetail,omitempty"`
}
